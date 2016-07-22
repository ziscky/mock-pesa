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
      {{if .ConfirmTrx}}
      <ns1:transactionConfirmResponse>
      {{else}}
      <ns1:processCheckOutResponse>
      {{end}}
         <RETURN_CODE>{{.ReturnCode}}</RETURN_CODE>
         <DESCRIPTION>{{.Description}}</DESCRIPTION>
         <TRX_ID>{{.TransactionID}}</TRX_ID>
         {{if .ConfirmTrx}}
         <MERCHANT_TRANSACTION_ID>{{.MerchantTransID}}</MERCHANT_TRANSACTION_ID>
         {{else}}
         <ENC_PARAMS>{{.EncodedParams}}</ENC_PARAMS>
         <CUST_MSG>{{.Message}}</CUST_MSG>
         {{end}}
      {{if .ConfirmTrx}}
      </ns1:transactionConfirmResponse>
      {{else}}
      </ns1:processCheckOutResponse>
      {{end}}
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
	transactionRespTPL = `
        <SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ns1="tns:ns">
   <SOAP-ENV:Body>
      <ns1:transactionStatusResponse>
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
      </ns1:transactionStatusResponse>
   </SOAP-ENV:Body>
   </SOAP-ENV:Envelope>
    `
	wsdl = `<?xml version="1.0" encoding="UTF-8"?>\r\n<wsdl:definitions xmlns:wsdl="http://schemas.xmlsoap.org/wsdl/" xmlns:s="http://www.w3.org/2001/XMLSchema" xmlns:soap="http://schemas.xmlsoap.org/wsdl/soap/" xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:tns="tns:ns" targetNamespace="tns:ns">\r\n   <wsdl:types>\r\n      <s:schema targetNamespace="tns:ns">\r\n\t\t <s:element name="CheckOutHeader">\r\n            <s:complexType>\r\n               <s:sequence>\r\n                  <s:element minOccurs="1" maxOccurs="1" name="MERCHANT_ID" type="s:string" />\r\n                  <s:element minOccurs="1" maxOccurs="1" name="PASSWORD" type="s:string" />\r\n                  <s:element minOccurs="1" maxOccurs="1" name="TIMESTAMP" type="s:string" />\r\n               </s:sequence>\r\n            </s:complexType>\r\n         </s:element>\r\n         <s:element name="processCheckOutRequest">\r\n            <s:complexType>\r\n               <s:sequence>\r\n                  <s:element minOccurs="1" maxOccurs="1" name="MERCHANT_TRANSACTION_ID" type="s:string" />\r\n                  <s:element minOccurs="1" maxOccurs="1" name="REFERENCE_ID" type="s:string" />\r\n                  <s:element minOccurs="1" maxOccurs="1" name="AMOUNT" type="s:double" />\r\n                  <s:element minOccurs="1" maxOccurs="1" name="MSISDN" type="s:string" />\r\n                  <s:element minOccurs="0" maxOccurs="1" name="ENC_PARAMS" type="s:string" />\r\n                  <s:element minOccurs="1" maxOccurs="1" name="CALL_BACK_URL" type="s:string" />\r\n                  <s:element minOccurs="1" maxOccurs="1" name="CALL_BACK_METHOD" type="s:string" />\r\n\t\t  <s:element minOccurs="0" maxOccurs="1" name="TIMESTAMP" type="s:string" />\r\n               </s:sequence>\r\n            </s:complexType>\r\n         </s:element>\r\n\t\t  <s:element name="transactionStatusRequest">\r\n            <s:complexType>\r\n               <s:sequence>\r\n                  <s:element minOccurs="0" maxOccurs="1" name="TRX_ID" type="s:string" />\r\n                  <s:element minOccurs="0" maxOccurs="1" name="MERCHANT_TRANSACTION_ID" type="s:string" />\r\n               </s:sequence>\r\n            </s:complexType>\r\n         </s:element>\r\n\t\r\n         <s:element name="processCheckOutResponse">\r\n            <s:complexType>\r\n               <s:sequence>\r\n                  <s:element minOccurs="1" name="RETURN_CODE" type="s:string" />\r\n                  <s:element minOccurs="1" name="DESCRIPTION" type="s:string" />\r\n                  <s:element minOccurs="1" name="TRX_ID" type="s:string" />\r\n                  <s:element minOccurs="1" name="ENC_PARAMS" type="s:string" />\r\n                  <s:element minOccurs="1" name="CUST_MSG" type="s:string" />\r\n               </s:sequence>\r\n            </s:complexType>\r\n         </s:element>\r\n        <s:element name="transactionStatusResponse">\r\n            <s:complexType>\r\n               <s:sequence>\r\n\t\t<s:element minOccurs="1" name="MSISDN" type="s:string" />\r\n\t\t<s:element minOccurs="1" name="AMOUNT" type="s:string" />\r\n\t\t<s:element minOccurs="1" name="MPESA_TRX_DATE" type="s:string" />\r\n\t\t<s:element minOccurs="1" name="MPESA_TRX_ID" type="s:string" />\r\n\t\t<s:element minOccurs="1" name="TRX_STATUS" type="s:string" />\r\n\t\t<s:element minOccurs="1" name="RETURN_CODE" type="s:string" />\r\n\t\t<s:element minOccurs="1" name="DESCRIPTION" type="s:string" />\r\n\t\t<s:element minOccurs="1" name="MERCHANT_TRANSACTION_ID" type="s:string" />\r\n\t\t<s:element minOccurs="1" name="ENC_PARAMS" type="s:string" />\r\n\t\t<s:element minOccurs="1" name="TRX_ID" type="s:string" />\r\n               </s:sequence>\r\n            </s:complexType>\r\n         </s:element>\r\n         <s:element name="transactionConfirmRequest">\r\n            <s:complexType>\r\n               <s:sequence>\r\n                  <s:element minOccurs="0" maxOccurs="1" name="TRX_ID" type="s:string" />\r\n                  <s:element minOccurs="0" maxOccurs="1" name="MERCHANT_TRANSACTION_ID" type="s:string" />\r\n               </s:sequence>\r\n            </s:complexType>\r\n         </s:element>\r\n\t <s:element name="transactionConfirmResponse">\r\n            <s:complexType>\r\n               <s:sequence>\r\n                  <s:element minOccurs="1" name="RETURN_CODE" type="s:string" />\r\n                  <s:element minOccurs="1" name="DESCRIPTION" type="s:string" />\r\n                  <s:element minOccurs="1" name="MERCHANT_TRANSACTION_ID" type="s:string" />\r\n\t\t  <s:element minOccurs="1" name="TRX_ID" type="s:string" />\r\n               </s:sequence>\r\n            </s:complexType>\r\n         </s:element>\r\n        <s:element name="ResultMsg">\r\n            <s:complexType>\r\n               <s:sequence>\r\n\t\t<s:element minOccurs="1" name="MSISDN" type="s:string" />\r\n\t\t<s:element minOccurs="1" name="AMOUNT" type="s:string" />\r\n\t\t<s:element minOccurs="1" name="MPESA_TRX_DATE" type="s:string" />\r\n\t\t<s:element minOccurs="1" name="MPESA_TRX_ID" type="s:string" />\r\n\t\t<s:element minOccurs="1" name="TRX_STATUS" type="s:string" />\r\n\t\t<s:element minOccurs="1" name="RETURN_CODE" type="s:string" />\r\n\t\t<s:element minOccurs="1" name="DESCRIPTION" type="s:string" />\r\n\t\t<s:element minOccurs="1" name="MERCHANT_TRANSACTION_ID" type="s:string" />\r\n\t\t<s:element minOccurs="1" name="ENC_PARAMS" type="s:string" />\r\n\t\t<s:element minOccurs="1" name="TRX_ID" type="s:string" />\r\n               </s:sequence>\r\n            </s:complexType>\r\n         </s:element>\r\n       <s:element name="ResponseMsg" type="s:string" />\r\n      </s:schema>\r\n   </wsdl:types>\r\n   <wsdl:message name="mpesaCheckOutRequest">\r\n      <wsdl:part name="body" element="tns:processCheckOutRequest" />\r\n   </wsdl:message>\r\n    <wsdl:message name="mpesaCheckOutHeader">\r\n\t  <wsdl:part name="header" element="tns:CheckOutHeader" />\r\n   </wsdl:message>\r\n   <wsdl:message name="mpesaCheckOutResponse">\r\n      <wsdl:part name="parameters" element="tns:processCheckOutResponse" />\r\n   </wsdl:message>\r\n   <wsdl:message name="mpesaTransactionRequest">\r\n      <wsdl:part name="body" element="tns:transactionStatusRequest" />\r\n   </wsdl:message>\r\n   <wsdl:message name="mpesaTransactionResponse">\r\n      <wsdl:part name="parameters" element="tns:transactionStatusResponse" />\r\n   </wsdl:message>\r\n   <wsdl:message name="mpesaConfirmRequest">\r\n      <wsdl:part name="parameters" element="tns:transactionConfirmRequest" />\r\n   </wsdl:message>\r\n   <wsdl:message name="mpesaConfirmResponse">\r\n      <wsdl:part name="parameters" element="tns:transactionConfirmResponse" />\r\n   </wsdl:message>\r\n   <wsdl:message name="ResultMessage">\r\n        <wsdl:part name="ResultMsg" element="tns:ResultMsg">\r\n        </wsdl:part>\r\n    </wsdl:message>\r\n    <wsdl:message name="ResponseMessage">\r\n        <wsdl:part name="ResponseMsg" element="tns:ResponseMsg">\r\n        </wsdl:part>\r\n    </wsdl:message>\r\n   <wsdl:portType name="LNMO_portType">\r\n      <wsdl:operation name="processCheckOut">\r\n         <wsdl:input message="tns:mpesaCheckOutRequest" />\r\n         <wsdl:output message="tns:mpesaCheckOutResponse" />\r\n      </wsdl:operation>\r\n\t  <wsdl:operation name="transactionStatusQuery">\r\n         <wsdl:input message="tns:mpesaTransactionRequest" />\r\n         <wsdl:output message="tns:mpesaTransactionResponse" />\r\n      </wsdl:operation>\r\n\t  <wsdl:operation name="confirmTransaction">\r\n         <wsdl:input message="tns:mpesaConfirmRequest" />\r\n         <wsdl:output message="tns:mpesaConfirmResponse" />\r\n      </wsdl:operation>\r\n      <wsdl:operation name="LNMOResult">\r\n            <wsdl:input message="tns:ResultMessage">\r\n            </wsdl:input>\r\n            <wsdl:output message="tns:ResponseMessage">\r\n            </wsdl:output>\r\n        </wsdl:operation>\r\n   </wsdl:portType>\r\n   <wsdl:binding name="LNMO_binding" type="tns:LNMO_portType">\r\n      <soap:binding transport="http://schemas.xmlsoap.org/soap/http" />\r\n      <wsdl:operation name="processCheckOut">\r\n         <soap:operation soapAction="" style="document" />\r\n         <wsdl:input>\r\n\t\t<soap:header message="tns:mpesaCheckOutHeader" part="header" use="literal"/>\r\n            <soap:body use="literal" />\r\n         </wsdl:input>\r\n         <wsdl:output>\r\n            <soap:body use="literal" />\r\n         </wsdl:output>\r\n      </wsdl:operation>\r\n\t  <wsdl:operation name="transactionStatusQuery">\r\n         <soap:operation soapAction="" style="document" />\r\n         <wsdl:input>\r\n\t\t\t<soap:header message="tns:mpesaCheckOutHeader" part="header" use="literal"/>\r\n            <soap:body use="literal" />\r\n         </wsdl:input>\r\n         <wsdl:output>\r\n            <soap:body use="literal" />\r\n         </wsdl:output>\r\n      </wsdl:operation>\r\n\t <wsdl:operation name="confirmTransaction">\r\n         <soap:operation soapAction="" style="document" />\r\n         <wsdl:input>\r\n\t\t\t<soap:header message="tns:mpesaCheckOutHeader" part="header" use="literal"/>\r\n            <soap:body use="literal" />\r\n         </wsdl:input>\r\n         <wsdl:output>\r\n            <soap:body use="literal" />\r\n         </wsdl:output>\r\n      </wsdl:operation>\r\n      <wsdl:operation name="LNMOResult">\r\n            <soap:operation soapAction="" style="document"/>\r\n            <wsdl:input>\r\n                <soap:body use="literal" />\r\n            </wsdl:input>\r\n            <wsdl:output>\r\n                <soap:body use="literal" />\r\n            </wsdl:output>\r\n        </wsdl:operation>\r\n   </wsdl:binding>\r\n   <wsdl:service name="lnmo_checkout_Service">\r\n      <wsdl:port name="lnmo_checkout" binding="tns:LNMO_binding">\r\n         <soap:address location="lnmo_checkout_server.php" />\r\n      </wsdl:port>\r\n   </wsdl:service>\r\n</wsdl:definitions>`
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
