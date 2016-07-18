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

import "encoding/xml"

//Header required for every "SOAP" request for authentication
type Header struct {
	XMLName        xml.Name       `xml:"Header"`
	CheckoutHeader CheckoutHeader `xml:"CheckOutHeader"`
}

//CheckoutHeader contains the actual authentication information
type CheckoutHeader struct {
	XMLName    xml.Name `xml:"CheckOutHeader"`
	MerchantID string   `xml:"MERCHANT_ID"`
	Password   string   `xml:"PASSWORD"`
	Timestamp  string   `xml:"TIMESTAMP"`
}

//ProcessCheckoutRequest structure for the ProcessCheckoutRequest request
type ProcessCheckoutRequest struct {
	XMLName xml.Name            `xml:"Envelope"`
	Header  Header              `xml:"Header"`
	Body    ProcessCheckoutBody `xml:"Body"`
}

//ProcessCheckoutBody structure for the ProcessCheckoutRequest body
type ProcessCheckoutBody struct {
	XMLName         xml.Name        `xml:"Body"`
	ProcessCheckout ProcessCheckout `xml:"processCheckOutRequest"`
}

//ProcessCheckout contains actual parameters for the ProcessCheckoutRequest method
type ProcessCheckout struct {
	XMLName         xml.Name `xml:"processCheckOutRequest"`
	MerchantTransID string   `xml:"MERCHANT_TRANSACTION_ID"`
	ReferenceID     string   `xml:"REFERENCE_ID"`
	Amount          string   `xml:"AMOUNT"`
	MSISDN          string   `xml:"MSISDN"`
	EncodedParams   string   `xml:"ENC_PARAMS"`
	CallBackURL     string   `xml:"CALL_BACK_URL"`
	CallBackMethod  string   `xml:"CALL_BACK_METHOD"`
	Timestamp       string   `xml:"TIMESTAMP"`
}

//ConfirmTransactionRequest structure for the ConfirmTransactionRequest method
type ConfirmTransactionRequest struct {
	XMLName xml.Name               `xml:"Envelope"`
	Header  Header                 `xml:"Header"`
	Body    ConfirmTransactionBody `xml:"Body"`
}

//ConfirmTransactionBody contains the body for the ConfirmTransactionRequest method
type ConfirmTransactionBody struct {
	XMLName            xml.Name           `xml:"Body"`
	ConfirmTransaction ConfirmTransaction `xml:"transactionConfirmRequest"`
}

//ConfirmTransaction contains actual parameters for the ConfirmTransactionRequest method
type ConfirmTransaction struct {
	XMLName         xml.Name `xml:"transactionConfirmRequest"`
	TransID         string   `xml:"TRX_ID"`
	MerchantTransID string   `xml:"MERCHANT_TRANSACTION_ID"`
}

//TransactionStatusRequest structure for the TransactionStatusRequest method
type TransactionStatusRequest struct {
	XMLName xml.Name              `xml:"Envelope"`
	Header  Header                `xml:"Header"`
	Body    TransactionStatusBody `xml:"Body"`
}

//TransactionStatusBody contains the body for the TransactionStatusRequest method
type TransactionStatusBody struct {
	XMLName           xml.Name          `xml:"Body"`
	TransactionStatus TransactionStatus `xml:"transactionStatusRequest"`
}

//TransactionStatus contains the actual parameters for the TransactionStatusRequest method
type TransactionStatus struct {
	XMLName         xml.Name `xml:"transactionStatusRequest"`
	TransID         string   `xml:"TRX_ID"`
	MerchantTransID string   `xml:"MERCHANT_TRANSACTION_ID"`
}
