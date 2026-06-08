// Package yellowcard provides a Go client for the YellowCard API,
// a crypto exchange and payment platform serving African markets.
//
// Usage:
//
//	client := yellowcard.New("api-key", "secret-key")
//	rates, err := client.Rates.List(ctx, yellowcard.RatesParams{})
package yellowcard

const (
	DefaultBaseURL = "https://api.yellowcard.io"
	SandboxBaseURL = "https://sandbox.yellowcard.io"
)

// Client is the entry point for all YellowCard API operations.
type Client struct {
	http     *httpClient
	Payments *PaymentsService
	Rates    *RatesService
	Networks *NetworksService
	Accounts *AccountsService
}

// New creates a new YellowCard API client.
func New(apiKey, secretKey string, opts ...Option) *Client {
	cfg := defaultConfig()
	for _, o := range opts {
		o(cfg)
	}
	hc := newHTTPClient(apiKey, secretKey, cfg)
	c := &Client{http: hc}
	c.Payments = &PaymentsService{client: hc}
	c.Rates = &RatesService{client: hc}
	c.Networks = &NetworksService{client: hc}
	c.Accounts = &AccountsService{client: hc}
	return c
}
