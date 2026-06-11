package yellowcard

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

// PaymentsService handles payment creation and retrieval.
type PaymentsService struct{ client *httpClient }

type Payment struct {
	Amount                int       `json:"amount"`
	ApiKey                string    `json:"apiKey"`
	ChannelId             string    `json:"channelId"`
	ConvertedAmount       int       `json:"convertedAmount"`
	Country               string    `json:"country"`
	CreatedAt             time.Time `json:"createdAt"`
	Currency              string    `json:"currency"`
	DirectSettlement      bool      `json:"directSettlement"`
	ExpiresAt             time.Time `json:"expiresAt"`
	FiatWallet            string    `json:"fiatWallet"`
	Id                    string    `json:"id"`
	PartnerFeeAmountLocal int       `json:"partnerFeeAmountLocal"`
	PartnerFeeAmountUSD   int       `json:"partnerFeeAmountUSD"`
	PartnerId             string    `json:"partnerId"`
	Rate                  float64   `json:"rate"`
	Recipient             Recipient `json:"recipient"`
	RequestSource         string    `json:"requestSource"`
	SequenceId            string    `json:"sequenceId"`
	ServiceFeeAmountLocal int       `json:"serviceFeeAmountLocal"`
	ServiceFeeAmountUSD   int       `json:"serviceFeeAmountUSD"`
	Source                Source    `json:"source"`
	Status                string    `json:"status"`
	UpdatedAt             time.Time `json:"updatedAt"`
}

type Recipient struct {
	Address  string `json:"address"`
	Country  string `json:"country"`
	Dob      string `json:"dob"`
	Email    string `json:"email"`
	IdNumber string `json:"idNumber"`
	IdType   string `json:"idType"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
}

type Source struct {
	AccountType   string `json:"accountType"`
	AccountNumber string `json:"accountNumber"`
	NetworkId     string `json:"networkId"`
}

type ChannelType string

type CustomerType string

type SortRangeBy string

type OrderBy string

const (
	Retail      CustomerType = "retail"
	Institution CustomerType = "institution"
	Bank        ChannelType  = "bank"
	Momo        CustomerType = "momo"
	CreatedAt   SortRangeBy  = "createdAt"
	UpdatedAt   SortRangeBy  = "updatedAt"
	Desc        OrderBy      = "desc"
	Asc         OrderBy      = "asc"
)

type ReceivePaymentRequest struct {
	Recipient    Recipient    `json:"recipient"`
	Source       Source       `json:"source"`
	ChannelId    string       `json:"channelId"`
	SequenceId   string       `json:"sequenceId"`
	Amount       int          `json:"amount"`
	Currency     string       `json:"currency"`
	Country      string       `json:"country"`
	Reason       string       `json:"reason"`
	ForceAccept  bool         `json:"forceAccept"`
	CustomerType CustomerType `json:"customerType"`
	ChannelType  ChannelType  `json:"channelType"`
}

// ReceiveQuery holds the query parameters for listing receive payments.
type ReceiveQuery struct {
	EndDate   string `url:"endDate,omitempty"`
	StartDate string `url:"startDate,omitempty"`
	StartAt   int32  `url:"startAt,omitempty"`
	PerPage   int32  `url:"perPage,omitempty"`
	// RangeBy filters by date field: "createdAt" (default) or "updatedAt"
	RangeBy SortRangeBy `url:"rangeBy,omitempty"`
	// SortBy sorts by date field: "createdAt" (default) or "updatedAt"
	SortBy SortRangeBy `url:"sortBy,omitempty"`
	// OrderBy sets sort direction: "asc" or "desc" (default)
	OrderBy OrderBy `url:"orderBy,omitempty"`
}

// Create initiates a new receive payment.
func (s *PaymentsService) Create(ctx context.Context, req ReceivePaymentRequest) (*Payment, error) {
	var p Payment
	if err := s.client.post(ctx, "/business/receive", req, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// Get retrieves a payment by its ID.
func (s *PaymentsService) Get(ctx context.Context, id string) (*Payment, error) {
	var p Payment
	if err := s.client.get(ctx, fmt.Sprintf("/business/receive/%s", id), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// Accept accepts a payment by its ID.
// https://sandbox.api.yellowcard.io/business/receive/{id}/accept
func (s *PaymentsService) Accept(ctx context.Context, id string) (*Payment, error) {
	var p Payment
	if err := s.client.post(ctx, fmt.Sprintf("/business/receive/%s/accept", id), nil, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// Deny denies a payment by its ID.
// https://sandbox.api.yellowcard.io/business/receive/{id}/deny
func (s *PaymentsService) Deny(ctx context.Context, id string) (*Payment, error) {
	var p Payment
	if err := s.client.post(ctx, fmt.Sprintf("/business/receive/%s/deny", id), nil, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// Cancel cancels a payment by its ID.
// https://sandbox.api.yellowcard.io/business/receive/{id}/cancel
func (s *PaymentsService) Cancel(ctx context.Context, id string) (*Payment, error) {
	var p Payment
	if err := s.client.post(ctx, fmt.Sprintf("/business/receive/%s/cancel", id), nil, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// Refund refunds a payment by its ID.
// https://sandbox.api.yellowcard.io/business/receive/{id}/refund
func (s *PaymentsService) Refund(ctx context.Context, id string) (*Payment, error) {
	var p Payment
	if err := s.client.post(ctx, fmt.Sprintf("/business/receive/%s/refund", id), nil, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// GetBySequenceID retrieves a payment by its caller-assigned sequence ID.
// https://sandbox.api.yellowcard.io/business/receive/sequence-id/{id}
func (s *PaymentsService) GetBySequenceID(ctx context.Context, sequenceID string) (*Payment, error) {
	var p Payment
	if err := s.client.get(ctx, fmt.Sprintf("/v2/business/receive/sequence-id/%s", sequenceID), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// List returns receive payments matching the query filters.
func (s *PaymentsService) List(ctx context.Context, query ReceiveQuery) ([]Payment, error) {
	q := url.Values{}
	if query.StartDate != "" {
		q.Set("startDate", query.StartDate)
	}
	if query.EndDate != "" {
		q.Set("endDate", query.EndDate)
	}
	if query.StartAt != 0 {
		q.Set("startAt", fmt.Sprintf("%d", query.StartAt))
	}
	if query.PerPage != 0 {
		q.Set("perPage", fmt.Sprintf("%d", query.PerPage))
	}
	if query.RangeBy != "" {
		q.Set("rangeBy", string(query.RangeBy))
	}
	if query.SortBy != "" {
		q.Set("sortBy", string(query.SortBy))
	}
	if query.OrderBy != "" {
		q.Set("orderBy", string(query.OrderBy))
	}
	path := "/business/receive"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}
	var resp struct {
		Collections []Payment `json:"collections"`
	}
	if err := s.client.get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return resp.Collections, nil
}
