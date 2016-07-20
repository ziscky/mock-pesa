package main

import "fmt"

var (
	c2bClient      *TestClient
	callbackClient *CallBack
	port           = "7000"
	callbackPort   = "7001"
	merchantID     = "12345"
	passkey        = "54321"
	address        = fmt.Sprintf("http://127.0.0.1:%s/", port)
	callback       = fmt.Sprintf("http://127.0.0.1:%s/", callbackPort)
	invalidRequest = `
    <SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ns1="tns:ns">
    <SOAP-ENV:Body>
    <ns1:invalidRequest>
    <RETURN_CODE></RETURN_CODE>
    <DESCRIPTION></DESCRIPTION>
    <TRX_ID></TRX_ID>
    <MERCHANT_TRANSACTION_ID></MERCHANT_TRANSACTION_ID>
    <ENC_PARAMS></ENC_PARAMS>
    <CUST_MSG></CUST_MSG>
    </ns1:invalidRequest>
    </SOAP-ENV:Body>
    </SOAP-ENV:Envelope>
    `
	processCheckoutRequest = `
	<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:tns="tns:ns">
	<soapenv:Header>
    <tns:CheckOutHeader>
	<MERCHANT_ID>{{.MerchantID}}</MERCHANT_ID>
	<PASSWORD>{{.Password}}</PASSWORD>
	<TIMESTAMP>{{.Timestamp}}</TIMESTAMP>
	</tns:CheckOutHeader>
	</soapenv:Header>
	<soapenv:Body>
	<tns:processCheckOutRequest>
	<MERCHANT_TRANSACTION_ID>{{.TransactionID}}</MERCHANT_TRANSACTION_ID>
	<REFERENCE_ID>{{.TransactionID}}</REFERENCE_ID>
	<AMOUNT>{{.Amount}}</AMOUNT>
	<MSISDN>{{.MSISDN}}</MSISDN>
	<!--Optional:-->
	<ENC_PARAMS></ENC_PARAMS>
	<CALL_BACK_URL>{{.Callback}}</CALL_BACK_URL>
	<CALL_BACK_METHOD>{{.Method}}</CALL_BACK_METHOD>
	<TIMESTAMP>{{.Timestamp}}</TIMESTAMP>
	</tns:processCheckOutRequest>
	</soapenv:Body>
	</soapenv:Envelope>
	`
	confirmTransactionRequest = `
	<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:tns="tns:ns">
	<soapenv:Header>
    <tns:CheckOutHeader>
    <MERCHANT_ID>{{.MerchantID}}</MERCHANT_ID>
	<PASSWORD>{{.Password}}</PASSWORD>
	<TIMESTAMP>{{.Timestamp}}</TIMESTAMP>
    </tns:CheckOutHeader>
    </soapenv:Header>
    <soapenv:Body>
    <tns:transactionConfirmRequest>
	<TRX_ID>{{.TrxID}}</TRX_ID>
    <MERCHANT_TRANSACTION_ID>{{.TransactionID}}</MERCHANT_TRANSACTION_ID>
    </tns:transactionConfirmRequest>
    </soapenv:Body>
	</soapenv:Envelope>
	`
	transactionStatusRequest = `
	<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:tns="tns:ns">
	<soapenv:Header>
    <tns:CheckOutHeader>
    <MERCHANT_ID>{{.MerchantID}}</MERCHANT_ID>
	<PASSWORD>{{.Password}}</PASSWORD>
	<TIMESTAMP>{{.Timestamp}}</TIMESTAMP>
    </tns:CheckOutHeader>
    </soapenv:Header>
    <soapenv:Body>
    <tns:transactionStatusRequest>
	<TRX_ID>{{.TrxID}}</TRX_ID>
    <MERCHANT_TRANSACTION_ID>{{.TransactionID}}</MERCHANT_TRANSACTION_ID>
    </tns:transactionStatusRequest>
    </soapenv:Body>
	</soapenv:Envelope>
	`
)
