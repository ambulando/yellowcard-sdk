package yellowcard

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

// RatesService handles exchange rate lookups.
type RatesService struct{ client *httpClient }

// Rate represents a currency exchange rate quote.
type Rate struct {
	Buy       float64   `json:"buy"`
	Sell      float64   `json:"sell"`
	Locale    string    `json:"locale"`
	RateId    string    `json:"rateId"`
	Code      string    `json:"code"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// List returns current exchange rates, optionally filtered.
func (s *RatesService) List(ctx context.Context, currency string) ([]Rate, error) {
	if currency == "" {
		return nil, fmt.Errorf("currency is required")
	}
	q := url.Values{}
	q.Set("currency", currency)
	path := "/business/rates"
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
