package yellowcard

import (
	"context"
	"fmt"
)

// PaymentsService handles payment creation and retrieval.
type PaymentsService struct{ client *httpClient }

// PaymentStatus enumerates known payment lifecycle states.
type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusProcessing PaymentStatus = "processing"
	PaymentStatusCompleted  PaymentStatus = "completed"
	PaymentStatusFailed     PaymentStatus = "failed"
	PaymentStatusExpired    PaymentStatus = "expired"
)

// Payment represents a YellowCard payment.
type Payment struct {
	ID          string        `json:"id"`
	SequenceID  string        `json:"sequenceId"`
	Status      PaymentStatus `json:"status"`
	Amount      float64       `json:"amount"`
	Currency    string        `json:"currency"`
	NetworkID   string        `json:"networkId"`
	ChannelID   string        `json:"channelId"`
	Reason      string        `json:"reason"`
	Destination Destination   `json:"destination"`
	CreatedAt   string        `json:"createdAt"`
	UpdatedAt   string        `json:"updatedAt"`
	ExpiresAt   string        `json:"expiresAt"`
}

// Destination holds recipient details for a payment.
type Destination struct {
	AccountName   string `json:"accountName"`
	AccountNumber string `json:"accountNumber"`
	NetworkID     string `json:"networkId"`
	Country       string `json:"country"`
}

// CreatePaymentRequest is the request body for creating a payment.
type CreatePaymentRequest struct {
	// SequenceID is a caller-assigned idempotency key.
	SequenceID  string      `json:"sequenceId"`
	Amount      float64     `json:"amount"`
	Currency    string      `json:"currency"`
	ChannelID   string      `json:"channelId"`
	Destination Destination `json:"destination"`
	Reason      string      `json:"reason,omitempty"`
}

// Create initiates a new payment.
func (s *PaymentsService) Create(ctx context.Context, req CreatePaymentRequest) (*Payment, error) {
	var p Payment
	if err := s.client.post(ctx, "/v2/business/payments", req, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// Get retrieves a payment by its ID.
func (s *PaymentsService) Get(ctx context.Context, id string) (*Payment, error) {
	var p Payment
	if err := s.client.get(ctx, fmt.Sprintf("/v2/business/payments/%s", id), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// GetBySequenceID retrieves a payment by its caller-assigned sequence ID.
func (s *PaymentsService) GetBySequenceID(ctx context.Context, sequenceID string) (*Payment, error) {
	var p Payment
	if err := s.client.get(ctx, fmt.Sprintf("/v2/business/payments/sequence/%s", sequenceID), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// Cancel cancels a pending payment.
func (s *PaymentsService) Cancel(ctx context.Context, id string) error {
	return s.client.delete(ctx, fmt.Sprintf("/v2/business/payments/%s", id))
}
