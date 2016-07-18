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

var (
	//processCheckOutRespTPL is the response for the ConfirmTransactionRequest and ProcessCheckoutRequest methods
	processCheckOutRespTPL = `
    <SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ns1="tns:ns">
   <SOAP-ENV:Body>
      <ns1:processCheckOutResponse>
         <RETURN_CODE>{{.ReturnCode}}</RETURN_CODE>
         <DESCRIPTION>{{.Description}}</DESCRIPTION>
         <TRX_ID>{{.TransactionID}}}</TRX_ID>
         {{if .ConfirmTrx}}
         <MERCHANT_TRANSACTION_ID>{{.MerchantTransID}}</MERCHANT_TRANSACTION_ID>
         {{else}}
         <ENC_PARAMS>{{.EncodedParams}}</ENC_PARAMS>
         <CUST_MSG>{{.Message}}</CUST_MSG>
         {{end}}
      </ns1:processCheckOutResponse>
   </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

	//callBackRespTPL is the template for the user callback request and TransactionStatusRequest method
	callBackRespTPL = `
    <SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ns1="tns:ns">
   <SOAP-ENV:Body>
      <ns1:ResultMsg>
      <MSISDN ns1:type="xsd:string">{{.MSISDN}}</MSISDN>
      <AMOUNT ns1:type="xsd:string">{{.Amount}}</AMOUNT>
      <M-PESA_TRX_DATE ns1:type="xsd:string">{{.MpesaTrxDate}}</M-PESA_TRX_DATE>
      <M-PESA_TRX_ID ns1:type="xsd:string">{{.MpesaTrxID}}</M-PESA_TRX_ID>
      <TRX_STATUS ns1:type="xsd:string">{{.TrxStatus}}</TRX_STATUS>
      <RETURN_CODE ns1:type="xsd:string">{{.ReturnCode}}</RETURN_CODE>
      <DESCRIPTION ns1:type="xsd:string">{{.Description}}</DESCRIPTION> 
      <MERCHANT_TRANSACTION_ID ns1:type="xsd:string">{{.MerchantTrxID}}</MERCHANT_TRANSACTION_ID>       
      <ENC_PARAMS ns1:type="xsd:string">{{.EncodedParams}}</ENC_PARAMS>    
      <TRX_ID ns1:type="xsd:string">{{.TrxID}}</TRX_ID>
      </ns1:ResultMsg>
   </SOAP-ENV:Body>
   </SOAP-ENV:Envelope>
    `
)

//ProcessCheckoutResponse data for the above template
type ProcessCheckoutResponse struct {
	ReturnCode      string
	Description     string
	TransactionID   string
	MerchantTransID string
	EncodedParams   string
	Message         string
	ConfirmTrx      bool
}

//Transaction data for the above template
type Transaction struct {
	MSISDN         string
	Amount         string
	MpesaTrxDate   string
	MpesaTrxID     string
	TrxStatus      string
	ReturnCode     string
	Description    string
	MerchantTrxID  string
	EncodedParams  string
	TrxID          string
	CallBackURL    string
	CallBackMethod string
}
