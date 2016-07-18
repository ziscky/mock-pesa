package main

import (
	"fmt"
	"net/http"
	"testing"

	"bytes"

	"github.com/ziscky/mock-pesa/c2b"
	"github.com/ziscky/mock-pesa/common"
)

var (
	c2bClient      *TestClient
	port           = "7000"
	address        = fmt.Sprintf("http://localhost:%s/", port)
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
	processCheckoutRequest = ``
)

type TestClient struct {
	c2bClient *c2b.C2B
}

func (tc *TestClient) Reset() {
	tc.c2bClient.Stop()
	tc.c2bClient.Start()
}

func init() {
	c2bClient = &TestClient{
		c2b.NewAPI(port, common.Config{
			MaxAmount:                    70000,
			MinAmount:                    10,
			MerchantID:                   "12345",
			CallBackDelay:                0,
			SAGPasskey:                   "54321",
			MaxCustomerTransactionPerDay: 150000,
			EnabledAPIS:                  []string{"c2b"},
		},
		)}
	c2bClient.c2bClient.Start()
}

func TestInvalidHeader(t *testing.T) {
	defer c2bClient.Reset()
	var (
		req  *http.Request
		resp *http.Response
		err  error
	)
	req, err = http.NewRequest("POST", address, bytes.NewBufferString(invalidRequest))
	if err != nil {
		t.Error(err)
	}
	req.Close = true
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Error(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected: %d Got: %d", http.StatusBadRequest, resp.StatusCode)
	}

}

func TestInvalidOperation(t *testing.T) {
	defer c2bClient.Reset()
	var (
		req  *http.Request
		resp *http.Response
		err  error
	)

	req, err = http.NewRequest("POST", address, bytes.NewBufferString(invalidRequest))
	if err != nil {
		t.Error(err)
	}
	req.Close = true
	req.Header.Add("Content-Type", "text/xml")
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Error(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected: %d Got: %d", http.StatusNotFound, resp.StatusCode)
	}
}

func TestProcessCheckOut(t *testing.T)                      {}
func TestProcessCheckOutInvalidMSISDN(t *testing.T)         {}
func TestProcessCheckOutInvalidAuth(t *testing.T)           {}
func TestProcessCheckOutInvalidCallBackURL(t *testing.T)    {}
func TestProcessCheckOutInvalidCallBackMethod(t *testing.T) {}
func TestProcessCheckOutInvalidMerchTrxID(t *testing.T)     {}
func TestProcessCheckOutInvalidAmount(t *testing.T)         {}

func TestConfirmTrx(t *testing.T)              {}
func TestConfirmTrxWithCode(t *testing.T)      {}
func TestConfirmTrxInvalidAuth(t *testing.T)   {}
func TestConfirmTrxInvalidParams(t *testing.T) {}
func TestConfirmTrxInvalidIDs(t *testing.T)    {}
func TestConfirmTrxInvalidAmount(t *testing.T) {}

func TestTrxStatus(t *testing.T)              {}
func TestTrxStatusInvalidAuth(t *testing.T)   {}
func TestTrxStatusInvalidParams(t *testing.T) {}
func TestTrxStatusInvalidIDs(t *testing.T)    {}
