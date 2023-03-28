package metering

import (
	"encoding/json"
	"fmt"
)

const (
	invoiceAPIName = "Invoice"
)

type InvoiceClient struct {
	BaseClient
}

func NewInvoiceClient(apiKey string, opts ...ClientOption) *InvoiceClient {
	bc := NewBaseClient(apiKey, opts...)
	ic := &InvoiceClient{BaseClient: *bc}
	ic.logf("Instantiating amberflo.io Invoice Client")
	return ic
}

type UserInvoice struct {
	AppliedCommitments        []AppliedCommitment `json:"appliedCommitments,omitempty"`
	InvoiceURI                string              `json:"invoiceUri,omitempty"`
	InvoiceKey                InvoiceKey          `json:"invoiceKey,omitempty"`
	PlanBillingPeriod         BillingPeriod       `json:"planBillingPeriod,omitempty"`
	InvoiceStartTimeInSeconds int                 `json:"invoiceStartTimeInSeconds,omitempty"`
	InvoiceEndTimeInSeconds   int                 `json:"invoiceEndTimeInSeconds,omitempty"`
	PlanName                  string              `json:"planName,omitempty"`
	InvoicePriceStatus        InvoicePriceStatus  `json:"invoicePriceStatus,omitempty"`
	PaymentStatus             PaymentStatus       `json:"paymentStatus,omitempty"`
	TotalBill                 TotalBill           `json:"totalBill,omitempty"`
}

type TotalBill struct {
	TotalPrice float64 `json:"totalPrice,omitempty"`
}

type InvoicePriceStatus string

const (
	InvoicePriceStatusOpen        InvoicePriceStatus = "OPEN"
	InvoicePriceStatusGracePeriod InvoicePriceStatus = "GRACE_PERIOD"
	InvoicePriceStatusPriceLocked InvoicePriceStatus = "PRICE_LOCKED"
)

type InvoiceKey struct {
	AccountID     string `json:"accountId,omitempty"`
	CustomerID    string `json:"customerId,omitempty"`
	Day           int    `json:"day,omitempty"`
	Month         int    `json:"month,omitempty"`
	Year          int    `json:"year,omitempty"`
	ProductID     string `json:"productId,omitempty"`
	ProductPlanID string `json:"productPlanId,omitempty"`
}

type AppliedCommitment struct {
	CommitmentId      string  `json:"commitmentId,omitempty"`
	FeeAmount         float64 `json:"feeAmount,omitempty"`
	FeeInBaseCurrency float64 `json:"feeInBaseCurrency,omitempty"`
	FeeInCredits      float64 `json:"feeInCredits,omitempty"`
	ObligationId      string  `json:"obligationId,omitempty"`
	ObligationName    string  `json:"obligationName,omitempty"`
}

func (ic *InvoiceClient) GetLatestInvoice(customerID string) (*UserInvoice, error) {
	signature := fmt.Sprintf("GetLatestInvoice(%v): ", customerID)

	url := fmt.Sprintf("%s/payments/billing/customer-product-invoice?customerId=%s&latest=true", Endpoint, customerID)

	body, err := ic.AmberfloHttpClient.sendHttpRequest(
		invoiceAPIName, url, "GET", nil)
	if err != nil {
		ic.logf("%s API error: %s", signature, err)
		return nil, fmt.Errorf("API error: %s", err)
	}

	var invoice UserInvoice
	err = json.Unmarshal(body, &invoice)
	if err != nil {
		return nil, fmt.Errorf("%s Error reading JSON body: %s", signature, err)
	}

	return &invoice, nil
}

func (ic *InvoiceClient) ListAllInvoices(customerID string) ([]UserInvoice, error) {
	signature := fmt.Sprintf("GetAllInvoices(%v): ", customerID)

	url := fmt.Sprintf("%s/payments/billing/customer-product-invoice/all?customerId=%s", Endpoint, customerID)

	body, err := ic.AmberfloHttpClient.sendHttpRequest(
		invoiceAPIName, url, "GET", nil)
	if err != nil {
		ic.logf("%s API error: %s", signature, err)
		return nil, fmt.Errorf("API error: %s", err)
	}

	var invoices []UserInvoice
	err = json.Unmarshal(body, &invoices)
	if err != nil {
		return nil, fmt.Errorf("%s Error reading JSON body: %s", signature, err)
	}

	return invoices, nil
}
