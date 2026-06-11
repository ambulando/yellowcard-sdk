package yellowcard

import "context"

// AccountsService handles business account and balance operations.
type AccountsService struct{ client *httpClient }

// Account represents a YellowCard business account.
type Account struct {
	Available    float64 `json:"available"`
	Currency     string  `json:"currency"`
	CurrencyType string  `json:"currencyType"`
}

// List returns all accounts for the authenticated business.
func (s *AccountsService) List(ctx context.Context) ([]Account, error) {
	var resp struct {
		Accounts []Account `json:"accounts"`
	}
	if err := s.client.get(ctx, "/business/account", &resp); err != nil {
		return nil, err
	}
	return resp.Accounts, nil
}
