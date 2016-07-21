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
         <TRX_ID>{{.TransactionID}}</TRX_ID>
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
	wsdl = `
    <wsdl:definitions targetNamespace="tns:ns"><wsdl:types><s:schema targetNamespace="tns:ns"><s:element name="CheckOutHeader"><s:complexType><s:sequence><s:element minOccurs="1" maxOccurs="1" name="MERCHANT_ID" type="s:string"/><s:element minOccurs="1" maxOccurs="1" name="PASSWORD" type="s:string"/><s:element minOccurs="1" maxOccurs="1" name="TIMESTAMP" type="s:string"/></s:sequence></s:complexType></s:element><s:element name="processCheckOutRequest"><s:complexType><s:sequence><s:element minOccurs="1" maxOccurs="1" name="MERCHANT_TRANSACTION_ID" type="s:string"/><s:element minOccurs="1" maxOccurs="1" name="REFERENCE_ID" type="s:string"/><s:element minOccurs="1" maxOccurs="1" name="AMOUNT" type="s:double"/><s:element minOccurs="1" maxOccurs="1" name="MSISDN" type="s:string"/><s:element minOccurs="0" maxOccurs="1" name="ENC_PARAMS" type="s:string"/><s:element minOccurs="1" maxOccurs="1" name="CALL_BACK_URL" type="s:string"/><s:element minOccurs="1" maxOccurs="1" name="CALL_BACK_METHOD" type="s:string"/><s:element minOccurs="0" maxOccurs="1" name="TIMESTAMP" type="s:string"/></s:sequence></s:complexType></s:element><s:element name="transactionStatusRequest"><s:complexType><s:sequence><s:element minOccurs="0" maxOccurs="1" name="TRX_ID" type="s:string"/><s:element minOccurs="0" maxOccurs="1" name="MERCHANT_TRANSACTION_ID" type="s:string"/></s:sequence></s:complexType></s:element><s:element name="processCheckOutResponse"><s:complexType><s:sequence><s:element minOccurs="1" name="RETURN_CODE" type="s:string"/><s:element minOccurs="1" name="DESCRIPTION" type="s:string"/><s:element minOccurs="1" name="TRX_ID" type="s:string"/><s:element minOccurs="1" name="ENC_PARAMS" type="s:string"/><s:element minOccurs="1" name="CUST_MSG" type="s:string"/></s:sequence></s:complexType></s:element><s:element name="transactionStatusResponse"><s:complexType><s:sequence><s:element minOccurs="1" name="MSISDN" type="s:string"/><s:element minOccurs="1" name="AMOUNT" type="s:string"/><s:element minOccurs="1" name="MPESA_TRX_DATE" type="s:string"/><s:element minOccurs="1" name="MPESA_TRX_ID" type="s:string"/><s:element minOccurs="1" name="TRX_STATUS" type="s:string"/><s:element minOccurs="1" name="RETURN_CODE" type="s:string"/><s:element minOccurs="1" name="DESCRIPTION" type="s:string"/><s:element minOccurs="1" name="MERCHANT_TRANSACTION_ID" type="s:string"/><s:element minOccurs="1" name="ENC_PARAMS" type="s:string"/><s:element minOccurs="1" name="TRX_ID" type="s:string"/></s:sequence></s:complexType></s:element><s:element name="transactionConfirmRequest"><s:complexType><s:sequence><s:element minOccurs="0" maxOccurs="1" name="TRX_ID" type="s:string"/><s:element minOccurs="0" maxOccurs="1" name="MERCHANT_TRANSACTION_ID" type="s:string"/></s:sequence></s:complexType></s:element><s:element name="transactionConfirmResponse"><s:complexType><s:sequence><s:element minOccurs="1" name="RETURN_CODE" type="s:string"/><s:element minOccurs="1" name="DESCRIPTION" type="s:string"/><s:element minOccurs="1" name="MERCHANT_TRANSACTION_ID" type="s:string"/><s:element minOccurs="1" name="TRX_ID" type="s:string"/></s:sequence></s:complexType></s:element><s:element name="ResultMsg"><s:complexType><s:sequence><s:element minOccurs="1" name="MSISDN" type="s:string"/><s:element minOccurs="1" name="AMOUNT" type="s:string"/><s:element minOccurs="1" name="MPESA_TRX_DATE" type="s:string"/><s:element minOccurs="1" name="MPESA_TRX_ID" type="s:string"/><s:element minOccurs="1" name="TRX_STATUS" type="s:string"/><s:element minOccurs="1" name="RETURN_CODE" type="s:string"/><s:element minOccurs="1" name="DESCRIPTION" type="s:string"/><s:element minOccurs="1" name="MERCHANT_TRANSACTION_ID" type="s:string"/><s:element minOccurs="1" name="ENC_PARAMS" type="s:string"/><s:element minOccurs="1" name="TRX_ID" type="s:string"/></s:sequence></s:complexType></s:element><s:element name="ResponseMsg" type="s:string"/></s:schema></wsdl:types><wsdl:message name="mpesaCheckOutRequest"><wsdl:part name="body" element="tns:processCheckOutRequest"/></wsdl:message><wsdl:message name="mpesaCheckOutHeader"><wsdl:part name="header" element="tns:CheckOutHeader"/></wsdl:message><wsdl:message name="mpesaCheckOutResponse"><wsdl:part name="parameters" element="tns:processCheckOutResponse"/></wsdl:message><wsdl:message name="mpesaTransactionRequest"><wsdl:part name="body" element="tns:transactionStatusRequest"/></wsdl:message><wsdl:message name="mpesaTransactionResponse"><wsdl:part name="parameters" element="tns:transactionStatusResponse"/></wsdl:message><wsdl:message name="mpesaConfirmRequest"><wsdl:part name="parameters" element="tns:transactionConfirmRequest"/></wsdl:message><wsdl:message name="mpesaConfirmResponse"><wsdl:part name="parameters" element="tns:transactionConfirmResponse"/></wsdl:message><wsdl:message name="ResultMessage"><wsdl:part name="ResultMsg" element="tns:ResultMsg">
        </wsdl:part></wsdl:message><wsdl:message name="ResponseMessage"><wsdl:part name="ResponseMsg" element="tns:ResponseMsg">
        </wsdl:part></wsdl:message><wsdl:portType name="LNMO_portType"><wsdl:operation name="processCheckOut"><wsdl:input message="tns:mpesaCheckOutRequest"/><wsdl:output message="tns:mpesaCheckOutResponse"/></wsdl:operation><wsdl:operation name="transactionStatusQuery"><wsdl:input message="tns:mpesaTransactionRequest"/><wsdl:output message="tns:mpesaTransactionResponse"/></wsdl:operation><wsdl:operation name="confirmTransaction"><wsdl:input message="tns:mpesaConfirmRequest"/><wsdl:output message="tns:mpesaConfirmResponse"/></wsdl:operation><wsdl:operation name="LNMOResult"><wsdl:input message="tns:ResultMessage">
            </wsdl:input><wsdl:output message="tns:ResponseMessage">
            </wsdl:output></wsdl:operation></wsdl:portType><wsdl:binding name="LNMO_binding" type="tns:LNMO_portType"><soap:binding transport="http://schemas.xmlsoap.org/soap/http"/><wsdl:operation name="processCheckOut"><soap:operation soapAction="" style="document"/><wsdl:input><soap:header message="tns:mpesaCheckOutHeader" part="header" use="literal"/><soap:body use="literal"/></wsdl:input><wsdl:output><soap:body use="literal"/></wsdl:output></wsdl:operation><wsdl:operation name="transactionStatusQuery"><soap:operation soapAction="" style="document"/><wsdl:input><soap:header message="tns:mpesaCheckOutHeader" part="header" use="literal"/><soap:body use="literal"/></wsdl:input><wsdl:output><soap:body use="literal"/></wsdl:output></wsdl:operation><wsdl:operation name="confirmTransaction"><soap:operation soapAction="" style="document"/><wsdl:input><soap:header message="tns:mpesaCheckOutHeader" part="header" use="literal"/><soap:body use="literal"/></wsdl:input><wsdl:output><soap:body use="literal"/></wsdl:output></wsdl:operation><wsdl:operation name="LNMOResult"><soap:operation soapAction="" style="document"/><wsdl:input><soap:body use="literal"/></wsdl:input><wsdl:output><soap:body use="literal"/></wsdl:output></wsdl:operation></wsdl:binding><wsdl:service name="lnmo_checkout_Service"><wsdl:port name="lnmo_checkout" binding="tns:LNMO_binding"><soap:address location="lnmo_checkout_server.php"/></wsdl:port></wsdl:service></wsdl:definitions>`
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
