package yellowcard

import (
	"context"
	"fmt"
	"net/url"
)

// RatesService handles exchange rate lookups.
type RatesService struct{ client *httpClient }

// Rate represents a currency exchange rate quote.
type Rate struct {
	From      string  `json:"from"`
	To        string  `json:"to"`
	Rate      float64 `json:"rate"`
	MinAmount float64 `json:"minAmount"`
	MaxAmount float64 `json:"maxAmount"`
	ExpiresAt string  `json:"expiresAt"`
}

// RatesParams filters the rates listing.
type RatesParams struct {
	// Filter by source currency (e.g. "USD")
	From string
	// Filter by destination currency (e.g. "GHS")
	To string
}

// List returns current exchange rates, optionally filtered.
func (s *RatesService) List(ctx context.Context, params RatesParams) ([]Rate, error) {
	q := url.Values{}
	if params.From != "" {
		q.Set("from", params.From)
	}
	if params.To != "" {
		q.Set("to", params.To)
	}
	path := "/v2/rates"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}
	var resp struct {
		Rates []Rate `json:"rates"`
	}
	if err := s.client.get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return resp.Rates, nil
}

// Get returns the rate for a specific currency pair.
func (s *RatesService) Get(ctx context.Context, from, to string) (*Rate, error) {
	var rate Rate
	if err := s.client.get(ctx, fmt.Sprintf("/v2/rates/%s/%s", from, to), &rate); err != nil {
		return nil, err
	}
	return &rate, nil
}
