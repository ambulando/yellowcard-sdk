package yellowcard

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAccountsService_List(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/business/account" {
			t.Errorf("path = %s, want /business/account", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"accounts": []Account{
				{Available: 500.00, Currency: "USD", CurrencyType: "fiat"},
				{Available: 1200.50, Currency: "GHS", CurrencyType: "fiat"},
			},
		})
	}))
	defer srv.Close()

	svc := &AccountsService{client: newTestHTTPClient(srv)}
	accounts, err := svc.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(accounts) != 2 {
		t.Fatalf("len = %d, want 2", len(accounts))
	}
	if accounts[0].Available != 500.00 {
		t.Errorf("Available = %f, want 500.00", accounts[0].Available)
	}
	if accounts[0].Currency != "USD" {
		t.Errorf("Currency = %q, want \"USD\"", accounts[0].Currency)
	}
	if accounts[0].CurrencyType != "fiat" {
		t.Errorf("CurrencyType = %q, want \"fiat\"", accounts[0].CurrencyType)
	}
	if accounts[1].Currency != "GHS" {
		t.Errorf("Currency = %q, want \"GHS\"", accounts[1].Currency)
	}
}

func TestAccountsService_List_Empty(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{"accounts": []Account{}})
	}))
	defer srv.Close()

	svc := &AccountsService{client: newTestHTTPClient(srv)}
	accounts, err := svc.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(accounts) != 0 {
		t.Errorf("len = %d, want 0", len(accounts))
	}
}

func TestAccountsService_List_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"code": "UNAUTHORIZED", "message": "invalid key"})
	}))
	defer srv.Close()

	svc := &AccountsService{client: newTestHTTPClient(srv)}
	_, err := svc.List(context.Background())
	if !IsUnauthorized(err) {
		t.Errorf("expected IsUnauthorized, got: %v", err)
	}
}
