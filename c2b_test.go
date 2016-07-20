package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	"bytes"

	"time"

	"encoding/xml"

	"github.com/braintree/manners"
	"github.com/gorilla/mux"
	"github.com/ziscky/mock-pesa/c2b"
	"github.com/ziscky/mock-pesa/common"
)

type TestClient struct {
	c2bClient *c2b.C2B
}

func (tc *TestClient) Reset() {
	tc.c2bClient.Clear() //clears the c2b API transaction storage
}

//CallBack allows for a minimal callback server for the purpose of data retrieval
type CallBack struct {
	data c2b.CallBackResponseContent
	read bool //check if data is ready to be read
	err  error
}

//Get waits for the server to retrieve data and returns it
func (cb *CallBack) Get() c2b.CallBackResponseContent {
	for !cb.read {
		time.Sleep(time.Millisecond * 1)
	}
	return cb.data
}

//Reset called after every callback response retrieval
func (cb *CallBack) Reset() {
	cb.data = c2b.CallBackResponseContent{}
	cb.read = false
	cb.err = nil
}

//Start sets the 2 method routers
//runs the callback server
func (cb *CallBack) Start() {
	router := mux.NewRouter()
	router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		data := new(c2b.CallBackResponseContent)

		header := r.Header.Get("Content-Type")
		if header == "text/xml" {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				rw.WriteHeader(500)
				cb.err = err
			}
			d := new(c2b.CallBackResponse)
			if err := xml.Unmarshal(b, d); err != nil {
				rw.WriteHeader(500)
				cb.err = err
			}
			cb.data = d.Body.Content
			cb.read = true
			rw.WriteHeader(200)
			return
		}
		data.MSISDN = r.FormValue("MSISDN")
		data.Amount = r.FormValue("AMOUNT")
		data.Date = r.FormValue("MPESA_TRX_DATE")
		data.MpesaTrxID = r.FormValue("MPESA_TRX_ID")
		data.TrxStatus = r.FormValue("TRX_STATUS")
		data.ReturnCode = r.FormValue("RETURN_CODE")
		data.Description = r.FormValue("DESCRIPTION")
		data.MerchantTrxID = r.FormValue("MERCHANT_TRANSACTION_ID")
		data.EncodedParams = r.FormValue("ENCODED_PARAMS")
		data.TrxID = r.FormValue("TRX_ID")
		cb.data = *data
		cb.read = true
		rw.WriteHeader(200)
		return
	}).Methods("POST")

	router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		data := new(c2b.CallBackResponseContent)
		params := r.URL.Query()
		data.MSISDN = params.Get("MSISDN")
		data.Amount = params.Get("AMOUNT")
		data.Date = params.Get("MPESA_TRX_DATE")
		data.MpesaTrxID = params.Get("MPESA_TRX_ID")
		data.TrxStatus = params.Get("TRX_STATUS")
		data.ReturnCode = params.Get("RETURN_CODE")
		data.Description = params.Get("DESCRIPTION")
		data.MerchantTrxID = params.Get("MERCHANT_TRANSACTION_ID")
		data.EncodedParams = params.Get("ENCODED_PARAMS")
		data.TrxID = params.Get("TRX_ID")
		cb.data = *data
		cb.read = true
		rw.WriteHeader(200)
	}).Methods("GET")

	manners.ListenAndServe(fmt.Sprintf(":%s", callbackPort), router)
}

func init() {
	c2bClient = &TestClient{
		c2b.NewAPI(port, common.Config{
			MaxAmount:                    70000,
			MinAmount:                    10,
			MerchantID:                   merchantID,
			CallBackDelay:                0,
			SAGPasskey:                   passkey,
			MaxCustomerTransactionPerDay: 150000,
			EnabledAPIS:                  []string{"c2b"},
		},
		)}
	c2bClient.c2bClient.Start()
	callbackClient = &CallBack{}
	go callbackClient.Start()
}

//doRequest makes an api call to the global address
//header -> specifies whether to include the text/html header
func doRequest(t *testing.T, data []byte, header bool) *http.Response {
	req, err := http.NewRequest("POST", address, bytes.NewBuffer(data))
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	req.Close = true
	if header {
		req.Header.Add("Content-Type", "text/xml")
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	return resp
}

//doCustom makes an api call to a custom address
func doCustom(t *testing.T, data []byte, header bool, address string) *http.Response {
	req, err := http.NewRequest("POST", address, bytes.NewBuffer(data))
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	req.Close = true
	if header {
		req.Header.Add("Content-Type", "text/xml")
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	return resp
}

//testConfirmTransactionDryRun builds the request template for the processCheckOut operation
func testProcessCheckOutDryRun(t *testing.T, method, msisdn, amount, trxid, merchID, callback string) bytes.Buffer {
	var (
		err  error
		temp *template.Template
		res  bytes.Buffer
	)
	ts := int(time.Now().UnixNano() / int64(time.Millisecond))
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%s%s%d", merchantID, passkey, ts)))
	password := base64.StdEncoding.EncodeToString([]byte(hex.EncodeToString(hash.Sum(nil))))
	params := map[string]interface{}{
		"MerchantID":    merchID,
		"Password":      password,
		"Timestamp":     strconv.Itoa(ts),
		"TransactionID": trxid,
		"Amount":        amount,
		"MSISDN":        msisdn,
		"Callback":      callback,
		"Method":        method,
	}

	temp, err = template.New("test").Parse(processCheckoutRequest)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if err := temp.Execute(&res, params); err != nil {
		t.Error(err)
		t.Fail()
	}
	return res
}

//testConfirmTransactionDryRun parses the request template for the confirmTransaction operation
func testConfirmTransactionDryRun(t *testing.T, merchID, trxid, systrxid string) bytes.Buffer {
	var (
		err  error
		temp *template.Template
		res  bytes.Buffer
	)
	ts := int(time.Now().UnixNano() / int64(time.Millisecond))
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%s%s%d", merchantID, passkey, ts)))
	password := base64.StdEncoding.EncodeToString([]byte(hex.EncodeToString(hash.Sum(nil))))
	params := map[string]interface{}{
		"MerchantID":    merchID,
		"Password":      password,
		"Timestamp":     strconv.Itoa(ts),
		"TransactionID": trxid,
		"TrxID":         systrxid,
	}

	temp, err = template.New("test").Parse(confirmTransactionRequest)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if err := temp.Execute(&res, params); err != nil {
		t.Error(err)
		t.Fail()
	}
	return res
}

//testTransactionQueryDryRun parses the request template for the transactionStatusRequest operation
func testTransactionQueryDryRun(t *testing.T, trxid, systrxid string) bytes.Buffer {
	var (
		err  error
		temp *template.Template
		res  bytes.Buffer
	)
	ts := int(time.Now().UnixNano() / int64(time.Millisecond))
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%s%s%d", merchantID, passkey, ts)))
	password := base64.StdEncoding.EncodeToString([]byte(hex.EncodeToString(hash.Sum(nil))))
	params := map[string]interface{}{
		"MerchantID":    merchantID,
		"Password":      password,
		"Timestamp":     strconv.Itoa(ts),
		"TransactionID": trxid,
		"TrxID":         systrxid,
	}

	temp, err = template.New("test").Parse(transactionStatusRequest)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if err := temp.Execute(&res, params); err != nil {
		t.Error(err)
		t.Fail()
	}
	return res
}

//TestInvalidHeader expects a failure when the text/xml header is not added to the request
func TestInvalidHeader(t *testing.T) {
	defer c2bClient.Reset()
	var (
		resp *http.Response
	)
	resp = doRequest(t, []byte(invalidRequest), false)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected: %d Got: %d", http.StatusBadRequest, resp.StatusCode)
	}

}

//TestInvalidOperation expects a failure on trying to make an invalid SOAP call
func TestInvalidOperation(t *testing.T) {
	defer c2bClient.Reset()
	var (
		resp *http.Response
	)

	resp = doRequest(t, []byte(invalidRequest), true)
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected: %d Got: %d", http.StatusNotFound, resp.StatusCode)
	}
}

//TestProcessCheckOut tests the processCheckOut operation with valid details
func TestProcessCheckOut(t *testing.T) {
	defer c2bClient.Reset()
	var (
		resp *http.Response
		res  bytes.Buffer
	)
	res = testProcessCheckOutDryRun(t, "post", "254723200817", "5500", "trx101", merchantID, callback)
	resp = doRequest(t, res.Bytes(), true)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected: %d Got: %d", http.StatusOK, resp.StatusCode)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	p := new(c2b.ProcessCheckOutResponse)
	if err := xml.Unmarshal(b, p); err != nil {
		t.Error(err)
		return
	}
	if p.Body.Content.ReturnCode != "00" {
		t.Errorf("Expected: 00 Got: %s", p.Body.Content.ReturnCode)
		return
	}
}

//TestProcessCheckOutInvalidMSISDN expects a failure due to an invalid phonenumber
func TestProcessCheckOutInvalidMSISDN(t *testing.T) {
	defer c2bClient.Reset()
	var (
		resp *http.Response
		res  bytes.Buffer
	)
	invalid := []string{"723", "0723200182", "254724", "invalid"}

	for _, msisdn := range invalid {
		res = testProcessCheckOutDryRun(t, "post", msisdn, "5500", "trx101", merchantID, callback)
		resp = doRequest(t, res.Bytes(), true)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected: %d Got: %d", http.StatusOK, resp.StatusCode)
		}
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
			return
		}
		p := new(c2b.ProcessCheckOutResponse)
		if err := xml.Unmarshal(b, p); err != nil {
			t.Error(err)
			return
		}
		if p.Body.Content.ReturnCode != "41" {
			t.Errorf("Expected: 41 Got: %s", p.Body.Content.ReturnCode)
			return
		}
	}
}

//TestProcessCheckOutInvalidAuth uses wrong authentication details to test integrity of auth
func TestProcessCheckOutInvalidAuth(t *testing.T) {
	defer c2bClient.Reset()
	var (
		resp *http.Response
		res  bytes.Buffer
	)
	res = testProcessCheckOutDryRun(t, "post", "254723200817", "5500", "trx101", "invalid", callback)
	resp = doRequest(t, res.Bytes(), true)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected: %d Got: %d", http.StatusOK, resp.StatusCode)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	p := new(c2b.ProcessCheckOutResponse)
	if err := xml.Unmarshal(b, p); err != nil {
		t.Error(err)
		return
	}
	if p.Body.Content.ReturnCode != "36" {
		t.Errorf("Expected: 36 Got: %s", p.Body.Content.ReturnCode)
		return
	}
}

//TestProcessCheckOutInvalidCallBackURL tests if the URL validation works
func TestProcessCheckOutInvalidCallBackURL(t *testing.T) {
	defer c2bClient.Reset()
	var (
		resp *http.Response
		res  bytes.Buffer
	)
	res = testProcessCheckOutDryRun(t, "post", "254723200817", "5500", "trx101", merchantID, "invalid")
	resp = doRequest(t, res.Bytes(), true)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected: %d Got: %d", http.StatusOK, resp.StatusCode)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	p := new(c2b.ProcessCheckOutResponse)
	if err := xml.Unmarshal(b, p); err != nil {
		t.Error(err)
		return
	}
	if p.Body.Content.ReturnCode != "40" {
		t.Errorf("Expected: 40 Got: %s", p.Body.Content.ReturnCode)
		return
	}
}

//TestProcessCheckOutInvalidCallBackMethod tests if the callback method is invalid.
//allowed: xml,post,get
func TestProcessCheckOutInvalidCallBackMethod(t *testing.T) {
	defer c2bClient.Reset()
	var (
		resp *http.Response
		res  bytes.Buffer
	)
	res = testProcessCheckOutDryRun(t, "invalid", "254723200817", "5500", "trx101", merchantID, callback)
	resp = doRequest(t, res.Bytes(), true)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected: %d Got: %d", http.StatusOK, resp.StatusCode)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	p := new(c2b.ProcessCheckOutResponse)
	if err := xml.Unmarshal(b, p); err != nil {
		t.Error(err)
		return
	}
	if p.Body.Content.ReturnCode != "40" {
		t.Errorf("Expected: 40 Got: %s", p.Body.Content.ReturnCode)
		return
	}
}

//TestProcessCheckOutInvalidMerchTrxID the MERCHANT_TRANSACTION_ID cannot be empty or reused
func TestProcessCheckOutInvalidMerchTrxID(t *testing.T) {
	defer c2bClient.Reset()
	var (
		resp *http.Response
		res  bytes.Buffer
	)
	res = testProcessCheckOutDryRun(t, "post", "254723200817", "5500", "", merchantID, callback)
	resp = doRequest(t, res.Bytes(), true)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected: %d Got: %d", http.StatusOK, resp.StatusCode)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	p := new(c2b.ProcessCheckOutResponse)
	if err := xml.Unmarshal(b, p); err != nil {
		t.Error(err)
		return
	}
	if p.Body.Content.ReturnCode != "40" {
		t.Errorf("Expected: 40 Got: %s", p.Body.Content.ReturnCode)
		return
	}
}

//TestProcessCheckOutInvalidAmount amount cannot be empty, has to be a number
func TestProcessCheckOutInvalidAmount(t *testing.T) {
	defer c2bClient.Reset()
	var (
		resp *http.Response
		res  bytes.Buffer
	)
	amounts := []string{"", "invlaid"}
	for _, amount := range amounts {
		res = testProcessCheckOutDryRun(t, "post", "254723200817", amount, "trx101", merchantID, callback)
		resp = doRequest(t, res.Bytes(), true)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected: %d Got: %d", http.StatusOK, resp.StatusCode)
		}
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
			return
		}
		p := new(c2b.ProcessCheckOutResponse)
		if err := xml.Unmarshal(b, p); err != nil {
			t.Error(err)
			return
		}
		if p.Body.Content.ReturnCode != "31" {
			t.Errorf("Expected: 31 Got: %s", p.Body.Content.ReturnCode)
			return
		}
	}
}

//TestConfirmTrx tests if a transaction can be initiated and confirmed
func TestConfirmTrx(t *testing.T) {
	defer c2bClient.Reset()
	var (
		resp *http.Response
		res  bytes.Buffer
	)

	methods := []string{"post", "xml", "get"}
	trxid := 0
	p := new(c2b.ProcessCheckOutResponse)

	for _, method := range methods {
		id := startTrx(t, method, strconv.Itoa(trxid), "5500")
		res = testConfirmTransactionDryRun(t, merchantID, strconv.Itoa(trxid), id)

		resp = doRequest(t, res.Bytes(), true)
		if callbackClient.err != nil {
			t.Error(callbackClient.err.Error())
			return
		}
		data := callbackClient.Get()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
			t.Fail()
		}

		if err := xml.Unmarshal(b, p); err != nil {
			t.Errorf(err.Error())
			t.Fail()
		}
		if p.Body.Content.ReturnCode != "00" {
			t.Errorf("Expected: 00 Got: %s", p.Body.Content.ReturnCode)
			return
		}
		if data.ReturnCode != "00" {
			t.Errorf("Expected: 00 Got: %s", data.ReturnCode)
			return
		}
		trxid++
		callbackClient.Reset()
	}

}

//startTrx convinience method to initiate a transaction
func startTrx(t *testing.T, method, trxid, amount string) string {
	res := testProcessCheckOutDryRun(t, method, "254723200817", amount, trxid, merchantID, callback)
	resp := doRequest(t, res.Bytes(), true)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected: %d Got: %d", http.StatusOK, resp.StatusCode)
		t.Fail()
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	p := new(c2b.ProcessCheckOutResponse)
	if err := xml.Unmarshal(b, p); err != nil {
		t.Error(err)
		t.Fail()
	}
	return p.Body.Content.TransactionID
}

//TestConfirmTrxWithCode tests if custom scenarios can be used via a modiied api URL
func TestConfirmTrxWithCode(t *testing.T) {
	defer c2bClient.Reset()
	var (
		resp *http.Response
		res  bytes.Buffer
	)
	codes := []string{"33", "11", "29", "32", "34", "08", "10"}
	trxid := 0
	p := new(c2b.ProcessCheckOutResponse)

	for _, code := range codes {
		id := startTrx(t, "post", strconv.Itoa(trxid), "5500")
		res = testConfirmTransactionDryRun(t, merchantID, strconv.Itoa(trxid), id)

		resp = doCustom(t, res.Bytes(), true, address+code)
		if callbackClient.err != nil {
			t.Error(callbackClient.err.Error())
			return
		}
		data := callbackClient.Get()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
			t.Fail()
		}

		if err := xml.Unmarshal(b, p); err != nil {
			t.Errorf(err.Error())
			t.Fail()
		}
		if p.Body.Content.ReturnCode != "00" {
			t.Errorf("Expected: 00 Got: %s", p.Body.Content.ReturnCode)
			return
		}
		if data.ReturnCode != code {
			t.Errorf("Expected: %s Got: %s", code, data.ReturnCode)
			return
		}
		trxid++
		callbackClient.Reset()
	}

}

//TestConfirmTrxInvalidParams tests that either MERCHANT_TRANSACTION_ID or TRX_ID are passed in order
//to confirm a transaction
func TestConfirmTrxInvalidParams(t *testing.T) {
	defer c2bClient.Reset()
	var (
		resp *http.Response
		res  bytes.Buffer
	)
	trxid := 0
	p := new(c2b.ProcessCheckOutResponse)

	_ = startTrx(t, "post", strconv.Itoa(trxid), "5500")
	res = testConfirmTransactionDryRun(t, merchantID, "", "")

	resp = doRequest(t, res.Bytes(), true)
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if err := xml.Unmarshal(b, p); err != nil {
		t.Errorf(err.Error())
		t.Fail()
	}
	if p.Body.Content.ReturnCode != "40" {
		t.Errorf("Expected: 40 Got: %s", p.Body.Content.ReturnCode)
		return
	}
	trxid++
	callbackClient.Reset()

}

//TestConfirmTrxInvalidIDs checks if non existent IDS can be used
func TestConfirmTrxInvalidIDs(t *testing.T) {
	defer c2bClient.Reset()
	var (
		resp *http.Response
		res  bytes.Buffer
	)
	trxid := 0
	p := new(c2b.ProcessCheckOutResponse)

	_ = startTrx(t, "post", strconv.Itoa(trxid), "5500")
	res = testConfirmTransactionDryRun(t, merchantID, "invalid", "invalid")

	resp = doRequest(t, res.Bytes(), true)

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if err := xml.Unmarshal(b, p); err != nil {
		t.Errorf(err.Error())
		t.Fail()
	}
	if p.Body.Content.ReturnCode != "12" {
		t.Errorf("Expected: 12 Got: %s", p.Body.Content.ReturnCode)
		return
	}
}

//TestConfirmTrxInvalidAmount checks if the amount satisfies the configurable
//system criteria
func TestConfirmTrxInvalidAmount(t *testing.T) {
	defer c2bClient.Reset()
	var (
		resp *http.Response
		res  bytes.Buffer
	)
	//2d table of amount:expected return code
	scenarios := map[string]string{
		"71000": "04",
		"5":     "03",
	}

	trxid := 0
	p := new(c2b.ProcessCheckOutResponse)

	for amount, expected := range scenarios {
		id := startTrx(t, "post", strconv.Itoa(trxid), amount)
		res = testConfirmTransactionDryRun(t, merchantID, strconv.Itoa(trxid), id)

		resp = doRequest(t, res.Bytes(), true)
		if callbackClient.err != nil {
			t.Error(callbackClient.err.Error())
			return
		}
		data := callbackClient.Get()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
			t.Fail()
		}

		if err := xml.Unmarshal(b, p); err != nil {
			t.Errorf(err.Error())
			t.Fail()
		}
		if p.Body.Content.ReturnCode != "00" {
			t.Errorf("Expected: 00 Got: %s", p.Body.Content.ReturnCode)
			return
		}
		if data.ReturnCode != expected {
			t.Errorf("Expected: %s Got: %s", expected, data.ReturnCode)
			return
		}
		trxid++
		callbackClient.Reset()
	}
}

//TestTrxStatus checks if an initiated Trxs status can be retrieved
func TestTrxStatus(t *testing.T) {
	defer c2bClient.Reset()
	var (
		resp *http.Response
		res  bytes.Buffer
	)
	trxid := 0
	p := new(c2b.TransactionStatusResponse)

	id := startTrx(t, "post", strconv.Itoa(trxid), "5500")
	res = testTransactionQueryDryRun(t, strconv.Itoa(trxid), id)

	resp = doRequest(t, res.Bytes(), true)

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if err := xml.Unmarshal(b, p); err != nil {
		t.Errorf(err.Error())
		t.Fail()
	}
	if p.Body.Content.ReturnCode != "00" {
		t.Errorf("Expected: 00 Got: %s", p.Body.Content.ReturnCode)
		return
	}

}

//TestTrxStatusInvalidParams checks if empty IDS can be used to retrieve the status
func TestTrxStatusInvalidParams(t *testing.T) {
	defer c2bClient.Reset()
	var (
		resp *http.Response
		res  bytes.Buffer
	)
	trxid := 0
	p := new(c2b.ProcessCheckOutResponse)

	startTrx(t, "post", strconv.Itoa(trxid), "5500")
	res = testTransactionQueryDryRun(t, "", "")

	resp = doRequest(t, res.Bytes(), true)

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if err := xml.Unmarshal(b, p); err != nil {
		t.Errorf(err.Error())
		t.Fail()
	}

	if p.Body.Content.ReturnCode != "40" {
		t.Errorf("Expected: 40 Got: %s", p.Body.Content.ReturnCode)
		return
	}
}

//TestTrxStatusInvalidIDs checks if the the ids passed are non existent
func TestTrxStatusInvalidIDs(t *testing.T) {
	defer c2bClient.Reset()
	var (
		resp *http.Response
		res  bytes.Buffer
	)
	trxid := 0
	p := new(c2b.ProcessCheckOutResponse)

	startTrx(t, "post", strconv.Itoa(trxid), "5500")
	res = testTransactionQueryDryRun(t, "invalid", "invalid")

	resp = doRequest(t, res.Bytes(), true)

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if err := xml.Unmarshal(b, p); err != nil {
		t.Errorf(err.Error())
		t.Fail()
	}
	if p.Body.Content.ReturnCode != "12" {
		t.Errorf("Expected: 12 Got: %s", p.Body.Content.ReturnCode)
		return
	}
}
