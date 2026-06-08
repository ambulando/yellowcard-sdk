package yellowcard

import "context"

// NetworksService handles supported blockchain networks and payment channels.
type NetworksService struct{ client *httpClient }

// Network represents a supported blockchain network or payment channel.
type Network struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Status     string   `json:"status"`
	Countries  []string `json:"countries"`
	MinAmount  float64  `json:"minAmount"`
	MaxAmount  float64  `json:"maxAmount"`
	FeePercent float64  `json:"feePercent"`
	FeeFixed   float64  `json:"feeFixed"`
}

// Channel represents a payment channel (mobile money, bank transfer, etc.).
type Channel struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Country  string `json:"country"`
	Currency string `json:"currency"`
	Type     string `json:"type"` // e.g. "mobileMoney", "bankTransfer"
	Status   string `json:"status"`
}

// List returns all supported networks.
func (s *NetworksService) List(ctx context.Context) ([]Network, error) {
	var resp struct {
		Networks []Network `json:"networks"`
	}
	if err := s.client.get(ctx, "/v2/networks", &resp); err != nil {
		return nil, err
	}
	return resp.Networks, nil
}

// Channels returns all supported payment channels, optionally filtered by country code.
func (s *NetworksService) Channels(ctx context.Context, country string) ([]Channel, error) {
	path := "/v2/channels"
	if country != "" {
		path += "?country=" + country
	}
	var resp struct {
		Channels []Channel `json:"channels"`
	}
	if err := s.client.get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return resp.Channels, nil
}
