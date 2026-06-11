package yellowcard

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRatesService_List(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/business/rates" {
			t.Errorf("path = %s, want /business/rates", r.URL.Path)
		}
		if r.URL.Query().Get("currency") != "USD" {
			t.Errorf("currency = %q, want \"USD\"", r.URL.Query().Get("currency"))
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"rates": []Rate{
				{Code: "GHS", Buy: 12.5, Sell: 12.0},
				{Code: "NGN", Buy: 1600.0, Sell: 1580.0},
			},
		})
	}))
	defer srv.Close()

	svc := &RatesService{client: newTestHTTPClient(srv)}
	rates, err := svc.List(context.Background(), "USD")
	if err != nil {
		t.Fatal(err)
	}
	if len(rates) != 2 {
		t.Fatalf("len = %d, want 2", len(rates))
	}
	if rates[0].Code != "GHS" {
		t.Errorf("rates[0].Code = %q, want \"GHS\"", rates[0].Code)
	}
	if rates[0].Buy != 12.5 {
		t.Errorf("rates[0].Buy = %f, want 12.5", rates[0].Buy)
	}
	if rates[0].Sell != 12.0 {
		t.Errorf("rates[0].Sell = %f, want 12.0", rates[0].Sell)
	}
}

func TestRatesService_List_RequiresCurrency(t *testing.T) {
	svc := &RatesService{client: newTestHTTPClient(httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})))}
	_, err := svc.List(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty currency, got nil")
	}
}

func TestRatesService_List_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"code": "UNAUTHORIZED", "message": "invalid key"})
	}))
	defer srv.Close()

	svc := &RatesService{client: newTestHTTPClient(srv)}
	_, err := svc.List(context.Background(), "USD")
	if !IsUnauthorized(err) {
		t.Errorf("expected IsUnauthorized, got: %v", err)
	}
}
