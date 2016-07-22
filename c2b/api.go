/*
Copyright (C) 2016  Eric Ziscky

    This program is free software; you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation; either version 2 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License along
    with this program; if not, write to the Free Software Foundation, Inc.,
    51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
*/
package c2b

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/ziscky/mock-pesa/common"

	"github.com/gorilla/mux"
	server "github.com/icub3d/graceful"
)

//Ident struct used to identify mock-pesa
//transactions stored in memory
type Ident struct {
	MerchTrxID string
	SysTrxID   string
}

//C2B defines the procedures as defined in
//the mpesa online checkout documentation.
//address -> api listening address
//callback -> channel used to call user callbacks :)
//store -> thread safe inmemory transaction storage
//config -> application configuration
//lock -> refer to store
type C2B struct {
	address  string
	callback chan *Transaction
	store    map[*Ident]interface{}
	config   common.Config
	lock     sync.Mutex
}

//NewAPI returns a reference to a C2B api instance
//takes port and conf obj as params
func NewAPI(port string, conf common.Config) *C2B {
	return &C2B{
		address:  fmt.Sprintf(":%s", port),
		store:    make(map[*Ident]interface{}),
		callback: make(chan *Transaction),
		config:   conf,
	}
}

//Start starts the c2b api server
//starts the callback listener
func (c2b *C2B) Start() {

	r := mux.NewRouter()
	r.HandleFunc("/", c2b.useMiddleware(
		common.CheckHeader(),
	))
	r.HandleFunc("/{code}", c2b.useMiddleware(
		common.CheckHeader(),
	))
	r.HandleFunc("/wsdl/get", func(rw http.ResponseWriter, r *http.Request) {
		r.Header.Add("Content-Type", "text/xml")
		data := bytes.Trim([]byte(wsdl), "\n\t")
		rw.Write(data)
	})

	go func() {
		server.ListenAndServe(c2b.address, r)
	}()

	go func() {
		for {
			select {
			case trx := <-c2b.callback:
				if trx == nil {
					return
				}
				time.Sleep(time.Second * time.Duration(c2b.config.CallBackDelay))
				c2b.callClientCallBack(trx)
			}
		}
	}()

}

//callClientCallBack calls the user specified callback with the method specified
func (c2b *C2B) callClientCallBack(trx *Transaction) {
	var request *http.Request
	client := &http.Client{}
	if strings.ToUpper(trx.CallBackMethod) == "POST" {
		data := url.Values{}
		data.Add("MSISDN", trx.MSISDN)
		data.Add("AMOUNT", trx.Amount)
		data.Add("MPESA_TRX_DATE", time.Now().Format(time.ANSIC))
		data.Add("MPESA_TRX_ID", common.GenerateMpesaTrx())
		data.Add("TRX_STATUS", trx.TrxStatus)
		data.Add("RETURN_CODE", trx.ReturnCode)
		data.Add("DESCRIPTION", trx.Description)
		data.Add("MERCHANT_TRANSACTION_ID", trx.MerchantTrxID)
		data.Add("ENCODED_PARAMS", trx.EncodedParams)
		data.Add("TRX_ID", trx.TrxID)

		request, _ = http.NewRequest("POST", trx.CallBackURL, bytes.NewBufferString(data.Encode()))
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	if strings.ToUpper(trx.CallBackMethod) == "XML" {
		var res bytes.Buffer
		tpl, _ := template.New("response").Parse(callBackRespTPL)
		tpl.Execute(&res, trx)
		request, _ = http.NewRequest("POST", trx.CallBackURL, bytes.NewBuffer(res.Bytes()))
		request.Header.Add("Content-Type", "text/xml")
	}

	if strings.ToUpper(trx.CallBackMethod) == "GET" {
		request, _ = http.NewRequest("GET", trx.CallBackURL, nil)
		query := request.URL.Query()
		query.Add("MSISDN", trx.MSISDN)
		query.Add("AMOUNT", trx.Amount)
		query.Add("MPESA_TRX_DATE", time.Now().Format(time.ANSIC))
		query.Add("MPESA_TRX_ID", common.GenerateMpesaTrx())
		query.Add("TRX_STATUS", trx.TrxStatus)
		query.Add("RETURN_CODE", trx.ReturnCode)
		query.Add("DESCRIPTION", trx.Description)
		query.Add("MERCHANT_TRANSACTION_ID", trx.MerchantTrxID)
		query.Add("ENCODED_PARAMS", trx.EncodedParams)
		query.Add("TRX_ID", trx.TrxID)
		request.URL.RawQuery = query.Encode()
	}

	_, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return
	}

}

//parseRequest emulates the SOAP style rpc call by figuring out
//which "method" is being called
//fix Me: improve performance
func (c2b *C2B) parseRequest(r *http.Request) http.HandlerFunc {
	body, err := ioutil.ReadAll(r.Body)

	processCheckoutOp := &ProcessCheckoutRequest{}
	confirmTrxOp := &ConfirmTransactionRequest{}
	transStatusOp := &TransactionStatusRequest{}

	err = xml.Unmarshal(body, processCheckoutOp)
	err = xml.Unmarshal(body, confirmTrxOp)
	err = xml.Unmarshal(body, transStatusOp)

	if err != nil {
		log.Println(err)
		return func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(500)
			rw.Write([]byte(err.Error()))
		}
	}
	operation := processCheckoutOp.Body.ProcessCheckout.XMLName.Local
	if strings.Contains(operation, "processCheckOut") {
		return c2b.processCheckout(processCheckoutOp)
	}
	operation = confirmTrxOp.Body.ConfirmTransaction.XMLName.Local
	if strings.Contains(operation, "transactionConfirmRequest") {
		return c2b.confirmTransaction(confirmTrxOp)
	}
	operation = transStatusOp.Body.TransactionStatus.XMLName.Local
	if strings.Contains(operation, "transactionStatusRequest") {
		return c2b.transactionStatus(transStatusOp)
	}
	return c2b.unknownOperation(operation)

}

//useMiddleware applys all Middleware parsed to the request context
func (c2b *C2B) useMiddleware(middlewares ...common.Middleware) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		for _, m := range middlewares {
			if !m(rw, r) {
				rw.WriteHeader(400)
				return
			}
			c2b.parseRequest(r)(rw, r)
		}
	}
}

//unknownOperation returns a status Bad Request if a method not specified
//in the WSDL is called
func (c2b *C2B) unknownOperation(data interface{}) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(404)
		rw.Write([]byte("Unknown Operation: " + data.(string)))
	}
}

//processCheckout handles the processCheckOutRequest method specified in the WSDL
//stores the transaction in the internal storage
//
func (c2b *C2B) processCheckout(data interface{}) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		parsed := data.(*ProcessCheckoutRequest)

		if !validate(rw,
			validAuthDetails(c2b.config.SAGPasskey, parsed.Header.CheckoutHeader),
			validMSISDN(parsed.Body.ProcessCheckout.MSISDN),
			validCallBackURL(parsed.Body.ProcessCheckout.CallBackURL),
			validCallBackMethod(parsed.Body.ProcessCheckout.CallBackMethod),
			validMerchTrxID(parsed.Body.ProcessCheckout.MerchantTransID),
			validAmount(parsed.Body.ProcessCheckout.Amount, c2b.config, false),
		) {
			return
		}

		trx := new(Transaction)
		trx.Amount = parsed.Body.ProcessCheckout.Amount
		trx.MerchantTrxID = parsed.Body.ProcessCheckout.MerchantTransID
		trx.MSISDN = parsed.Body.ProcessCheckout.MSISDN
		trx.EncodedParams = parsed.Body.ProcessCheckout.EncodedParams
		trx.CallBackURL = parsed.Body.ProcessCheckout.CallBackURL
		trx.CallBackMethod = parsed.Body.ProcessCheckout.CallBackMethod
		trx.TrxStatus = "Success"
		trx.ReturnCode = "00"
		trx.TrxID = bson.NewObjectId().Hex()

		if c2b.idExists(trx.MerchantTrxID, trx.TrxID) != nil {
			resp := new(ProcessCheckoutResponse)
			resp.ReturnCode = duplicateRequest
			resp.Description = "Failed"
			resp.Message = "MERCHANT_TRANSACTION_ID already used"
			resp.TransactionID = trx.TrxID
			tpl, _ := template.New("response").Parse(processCheckOutRespTPL)
			tpl.Execute(rw, resp)
			return
		}

		c2b.lock.Lock()
		c2b.store[&Ident{trx.MerchantTrxID, trx.TrxID}] = trx
		c2b.lock.Unlock()

		resp := new(ProcessCheckoutResponse)
		resp.ReturnCode = success
		resp.Description = "Success"
		resp.Message = "To complete this transaction, enter your PIN on your handset. if you don't have one enter 0 and follow instructions"
		resp.TransactionID = trx.TrxID
		tpl, _ := template.New("response").Parse(processCheckOutRespTPL)
		tpl.Execute(rw, resp)
	}
}

//confirmTransaction handles the transactionConfirmRequest method as specified in the WSDL
//can take extra url params, to "simulate" certain scenarios, useful for unit tests
//and payment system hardening against all edge cases. Cow says moo
func (c2b *C2B) confirmTransaction(data interface{}) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		parsed := data.(*ConfirmTransactionRequest)

		vars := mux.Vars(r)
		code, ok := vars["code"]

		if !validate(rw,
			validAuthDetails(c2b.config.SAGPasskey, parsed.Header.CheckoutHeader),
			validPassedConfirmTrxID(parsed.Body.ConfirmTransaction.TransID, parsed.Body.ConfirmTransaction.MerchantTransID),
		) {
			return
		}

		trx := c2b.idExists(parsed.Body.ConfirmTransaction.MerchantTransID, parsed.Body.ConfirmTransaction.TransID)

		if trx == nil {
			resp := new(ProcessCheckoutResponse)
			resp.ReturnCode = transactionMismatch
			resp.Description = "transaction details are different from original captured request details."
			resp.TransactionID = ""
			resp.ConfirmTrx = true
			tpl, _ := template.New("response").Parse(processCheckOutRespTPL)
			tpl.Execute(rw, resp)
			return
		}

		obj := validAmount(trx.Amount, c2b.config, true)()
		if obj != nil {
			trx.Description = obj.Description
			trx.ReturnCode = obj.ReturnCode
		}

		if ok {
			trx.ReturnCode = code

			if code != "00" {
				trx.TrxStatus = "Failure"
				trx.Description = ResolveCode(code)
			}
		}
		c2b.callback <- trx

		resp := new(ProcessCheckoutResponse)
		resp.ReturnCode = success
		resp.Description = "Success"
		resp.TransactionID = trx.TrxID
		resp.MerchantTransID = trx.MerchantTrxID
		resp.ConfirmTrx = true
		tpl, _ := template.New("response").Parse(processCheckOutRespTPL)
		tpl.Execute(rw, resp)
	}
}

//transactionStatus handles the transactionStatusRequest as defined in the WSDL
//response is the XML version of the client callback
func (c2b *C2B) transactionStatus(data interface{}) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		parsed := data.(*TransactionStatusRequest)

		if !validate(rw,
			validAuthDetails(c2b.config.SAGPasskey, parsed.Header.CheckoutHeader),
			validPassedConfirmTrxID(parsed.Body.TransactionStatus.TransID, parsed.Body.TransactionStatus.MerchantTransID),
		) {
			return
		}

		trx := c2b.idExists(parsed.Body.TransactionStatus.MerchantTransID, parsed.Body.TransactionStatus.TransID)

		if trx == nil {
			resp := new(ProcessCheckoutResponse)
			resp.ReturnCode = transactionMismatch
			resp.Description = "transaction details are different from original captured request details."
			resp.TransactionID = ""
			resp.ConfirmTrx = true
			tpl, _ := template.New("response").Parse(processCheckOutRespTPL)
			tpl.Execute(rw, resp)
			return
		}
		tpl, _ := template.New("response").Parse(transactionRespTPL)
		tpl.Execute(rw, trx)
	}
}

//idExists checks if the merchant/system transaction id already exists
func (c2b *C2B) idExists(merchID, sysID string) *Transaction {
	c2b.lock.Lock()
	defer c2b.lock.Unlock()
	for k, v := range c2b.store {
		if k.MerchTrxID == merchID || k.SysTrxID == sysID {
			return v.(*Transaction)
		}
	}
	return nil
}

func (c2b *C2B) Clear() {
	c2b.lock.Lock()
	defer c2b.lock.Unlock()
	c2b.store = make(map[*Ident]interface{})
}

//Stop gracefully stops the API and the callback listener
func (c2b *C2B) Stop() {
	server.Close()
	c2b.store = make(map[*Ident]interface{}) //clear map
	c2b.callback <- nil
	// fmt.Println("C2B stopped")
}
