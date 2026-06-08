package yellowcard

import "context"

// AccountsService handles business account and balance operations.
type AccountsService struct{ client *httpClient }

// Account represents a YellowCard business account.
type Account struct {
	ID       string  `json:"id"`
	Label    string  `json:"label"`
	Currency string  `json:"currency"`
	Balance  float64 `json:"balance"`
	Status   string  `json:"status"`
}

// List returns all accounts for the authenticated business.
func (s *AccountsService) List(ctx context.Context) ([]Account, error) {
	var resp struct {
		Accounts []Account `json:"accounts"`
	}
	if err := s.client.get(ctx, "/v2/business/accounts", &resp); err != nil {
		return nil, err
	}
	return resp.Accounts, nil
}
