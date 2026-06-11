package yellowcard

import (
	"context"
	"time"
)

// NetworksService handles supported blockchain networks and payment channels.
type NetworksService struct{ client *httpClient }

// Network represents a supported blockchain network or payment channel.
type Network struct {
	Id                       string    `json:"id"`
	Code                     string    `json:"code"`
	UpdatedAt                time.Time `json:"updatedAt"`
	Status                   string    `json:"status"`
	CreatedAt                time.Time `json:"createdAt"`
	AccountNumberType        string    `json:"accountNumberType"`
	Country                  string    `json:"country"`
	Name                     string    `json:"name"`
	ChannelIds               []string  `json:"channelIds"`
	CountryAccountNumberType string    `json:"countryAccountNumberType"`
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
// https://sandbox.api.yellowcard.io/business/networks
func (s *NetworksService) List(ctx context.Context, country string) ([]Network, error) {
	var resp []Network
	path := "/business/networks"
	if country != "" {
		path += "?country=" + country
	}
	if err := s.client.get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Channels returns all supported payment channels, optionally filtered by country code.
func (s *NetworksService) Channels(ctx context.Context, country string) ([]Channel, error) {
	path := "/business/channels"
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
