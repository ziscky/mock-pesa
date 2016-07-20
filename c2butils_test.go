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
