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
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/ziscky/mock-pesa/common"

	"github.com/asaskevich/govalidator"
	"gopkg.in/mgo.v2/bson"
)

//validator type to validate parameters passed to the "SOAP" methods and respond accordingly
type validator func() *ProcessCheckoutResponse

//validate runs all the validators and renders response template if need be
func validate(rw http.ResponseWriter, validators ...validator) bool {
	tpl, _ := template.New("response").Parse(processCheckOutRespTPL)
	for _, v := range validators {
		resp := v()
		if resp != nil {
			tpl.Execute(rw, resp)
			return false
		}

	}
	return true
}

//validMerchID checks if the merchantID is actually registered on the "system"
func validMerchID(merchantID string, conf common.Config) validator {
	return func() *ProcessCheckoutResponse {
		if merchantID != conf.MerchantID {
			resp := new(ProcessCheckoutResponse)
			resp.ReturnCode = invalidMerchID
			resp.Description = "Invalid Merchant ID"
			resp.TransactionID = bson.NewObjectId().Hex()

			return resp
		}
		return nil
	}
}

//validAuthDetails checks if the password was encoded correctly
//Official mpesa docs are misguiding on this
//procedure -> (append merchantID,passkey,timestamp) -> (get sha256 sum encode to hex) -> (encode the result to standard base64
//as specified by  RFC 4648)
func validAuthDetails(password string, header CheckoutHeader) validator {
	return func() *ProcessCheckoutResponse {
		hash := sha256.New()
		hash.Write([]byte(fmt.Sprintf("%s%s%s", header.MerchantID, password, header.Timestamp)))
		password := base64.StdEncoding.EncodeToString([]byte(hex.EncodeToString(hash.Sum(nil))))
		if password == header.Password {
			return nil
		}
		resp := new(ProcessCheckoutResponse)
		resp.ReturnCode = authenticationFailure
		resp.Description = "Invalid Password"
		resp.TransactionID = bson.NewObjectId().Hex()

		return resp
	}
}

//validMSISDN the Official Mpesa docs specify that the phone number should begin with 254...
func validMSISDN(data string) validator {
	return func() *ProcessCheckoutResponse {
		flag := 0
		message := ""
		if len(data) < 12 {
			flag++
			message += "Invalid MSISDN Length,"
		}
		if len(data) > 3 {
			if data[:3] != "254" {
				flag++
				message += "MSISDN Format 254...,"
			}
		}
		if _, err := strconv.Atoi(data); err != nil {
			flag++
			message += "MSISDN not a number,"
		}
		if flag > 0 {
			resp := new(ProcessCheckoutResponse)
			resp.ReturnCode = incorrectMSISDN
			resp.Description = message
			resp.TransactionID = bson.NewObjectId().Hex()
			return resp
		}
		return nil
	}
}

//validCallBackURL the Mpesa *sys does not check for this* added as a convinience
func validCallBackURL(url string) validator {
	return func() *ProcessCheckoutResponse {
		if !govalidator.IsURL(url) {
			resp := new(ProcessCheckoutResponse)
			resp.ReturnCode = missingParameters
			resp.Description = "Invalid URL"
			resp.TransactionID = bson.NewObjectId().Hex()
			return resp
		}
		return nil
	}

}

//validCallBackMethod again the Mpesa sys *does not check for this* added as a convinience
//valid-> xml,post,get (case insensitive)
func validCallBackMethod(method string) validator {
	return func() *ProcessCheckoutResponse {
		if strings.ToUpper(method) == "POST" || strings.ToUpper(method) == "GET" || strings.ToUpper(method) == "XML" {
			return nil
		}
		resp := new(ProcessCheckoutResponse)
		resp.ReturnCode = missingParameters
		resp.Description = "Invalid Callback Method"
		resp.TransactionID = bson.NewObjectId().Hex()

		return resp

	}
}

//validMerchTrxID basically checks if the MERCHANT_TRANSACTION_ID is present
func validMerchTrxID(trx string) validator {
	return func() *ProcessCheckoutResponse {
		if len(trx) > 0 {
			return nil
		}
		resp := new(ProcessCheckoutResponse)
		resp.ReturnCode = missingParameters
		resp.Description = "Invalid Callback Method"
		resp.TransactionID = bson.NewObjectId().Hex()

		return resp
	}
}

//validPassedConfirmTrxID for the ConfirmTransactionRequest method either the MERCHANT_TRANSACTION_ID or TRX_ID can
//be specified, ensures atleast one is passed
func validPassedConfirmTrxID(sysTrx, merchTrx string) validator {
	return func() *ProcessCheckoutResponse {
		if len(sysTrx) == 0 && len(merchTrx) == 0 {
			resp := new(ProcessCheckoutResponse)
			resp.ReturnCode = missingParameters
			resp.Description = "Specify Either MERCHANT_TRANSACTION_ID/TRX_ID"
			resp.TransactionID = ""
			resp.ConfirmTrx = true
			return resp
		}
		return nil
	}
}

//validAmount checks if:
//amount is empty
//amount is not a float/double
//amount is greater than System allowed transaction amount(configurable)
//amount is less than System allowed minimum amount(configurable)
//amount is higher than the customer allowed daily limit
//confirm -> some amount checking is only done by the Safaricom API after confirmation
func validAmount(amount string, conf common.Config, confirm bool) validator {
	return func() *ProcessCheckoutResponse {
		resp := new(ProcessCheckoutResponse)
		if len(amount) == 0 {
			resp.ReturnCode = invalidAmount
			resp.Description = "Invalid Amount"
			resp.TransactionID = bson.NewObjectId().Hex()
			return resp
		}
		a, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			resp.ReturnCode = invalidAmount
			resp.Description = "Invalid Amount"
			resp.TransactionID = bson.NewObjectId().Hex()
			return resp
		}
		if confirm {
			if a > conf.MaxAmount {
				resp.ReturnCode = maxAmountReached
				resp.Description = "Amount above allowed maximum"
				resp.TransactionID = bson.NewObjectId().Hex()
				return resp
			}
			if a < conf.MinAmount {
				resp.ReturnCode = minAmountReached
				resp.Description = "Amount below allowed minimum"
				resp.TransactionID = bson.NewObjectId().Hex()
				return resp
			}
			if a > conf.MaxCustomerTransactionPerDay {
				resp.ReturnCode = maxDailyAmountReached
				resp.Description = "Customer max daily amount reached"
				resp.TransactionID = bson.NewObjectId().Hex()
				return resp
			}
		}
		return nil
	}
}
