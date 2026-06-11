package yellowcard

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNetworksService_List(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/business/networks" {
			t.Errorf("path = %s, want /business/networks", r.URL.Path)
		}
		if r.URL.RawQuery != "" {
			t.Errorf("unexpected query: %s", r.URL.RawQuery)
		}
		json.NewEncoder(w).Encode([]Network{
			{Id: "net-1", Code: "MTN_GH", Country: "GH", Name: "MTN Ghana"},
		})
	}))
	defer srv.Close()

	svc := &NetworksService{client: newTestHTTPClient(srv)}
	nets, err := svc.List(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}
	if len(nets) != 1 {
		t.Fatalf("len = %d, want 1", len(nets))
	}
	if nets[0].Id != "net-1" {
		t.Errorf("Id = %q, want \"net-1\"", nets[0].Id)
	}
	if nets[0].Code != "MTN_GH" {
		t.Errorf("Code = %q, want \"MTN_GH\"", nets[0].Code)
	}
}

func TestNetworksService_List_WithCountry(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("country") != "GH" {
			t.Errorf("country = %q, want \"GH\"", r.URL.Query().Get("country"))
		}
		json.NewEncoder(w).Encode([]Network{{Id: "net-1", Country: "GH"}})
	}))
	defer srv.Close()

	svc := &NetworksService{client: newTestHTTPClient(srv)}
	_, err := svc.List(context.Background(), "GH")
	if err != nil {
		t.Fatal(err)
	}
}

func TestNetworksService_Channels(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/business/channels" {
			t.Errorf("path = %s, want /business/channels", r.URL.Path)
		}
		if r.URL.RawQuery != "" {
			t.Errorf("unexpected query: %s", r.URL.RawQuery)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"channels": []Channel{
				{ID: "ch-1", Name: "MTN Mobile Money", Type: "mobileMoney", Currency: "GHS", Country: "GH", Status: "active"},
			},
		})
	}))
	defer srv.Close()

	svc := &NetworksService{client: newTestHTTPClient(srv)}
	channels, err := svc.Channels(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}
	if len(channels) != 1 {
		t.Fatalf("len = %d, want 1", len(channels))
	}
	if channels[0].ID != "ch-1" {
		t.Errorf("ID = %q, want \"ch-1\"", channels[0].ID)
	}
	if channels[0].Type != "mobileMoney" {
		t.Errorf("Type = %q, want \"mobileMoney\"", channels[0].Type)
	}
}

func TestNetworksService_Channels_WithCountry(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("country") != "NG" {
			t.Errorf("country = %q, want \"NG\"", r.URL.Query().Get("country"))
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"channels": []Channel{}})
	}))
	defer srv.Close()

	svc := &NetworksService{client: newTestHTTPClient(srv)}
	_, err := svc.Channels(context.Background(), "NG")
	if err != nil {
		t.Fatal(err)
	}
}

func TestNetworksService_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"code": "UNAUTHORIZED", "message": "invalid key"})
	}))
	defer srv.Close()

	svc := &NetworksService{client: newTestHTTPClient(srv)}
	_, err := svc.Channels(context.Background(), "")
	if !IsUnauthorized(err) {
		t.Errorf("expected IsUnauthorized, got: %v", err)
	}
}
