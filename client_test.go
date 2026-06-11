package yellowcard

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// newTestHTTPClient creates an httpClient pointed at the given test server.
func newTestHTTPClient(srv *httptest.Server) *httpClient {
	return newHTTPClient("test-api-key", "test-secret", &config{
		baseURL:    srv.URL,
		httpClient: http.DefaultClient,
	})
}

func TestClientAuthHeadersGET(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "YcHmacV1 test-api-key:") {
			t.Errorf("Authorization = %q, want prefix \"YcHmacV1 test-api-key:\"", auth)
		}
		ts := r.Header.Get("X-YC-Timestamp")
		if ts == "" {
			t.Error("X-YC-Timestamp header missing")
		}
		if _, err := time.Parse(time.RFC3339, ts); err != nil {
			t.Errorf("X-YC-Timestamp not RFC3339: %q", ts)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	var out struct{}
	_ = newTestHTTPClient(srv).do(context.Background(), http.MethodGet, "/test", nil, &out)
}

func TestClientAuthHeadersPOST(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "YcHmacV1 test-api-key:") {
			t.Errorf("Authorization = %q", auth)
		}
		if r.Header.Get("X-YC-Timestamp") == "" {
			t.Error("X-YC-Timestamp header missing")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	var out struct{}
	_ = newTestHTTPClient(srv).do(context.Background(), http.MethodPost, "/test", map[string]string{"k": "v"}, &out)
}

func TestClientDecodeResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"foo": "bar"})
	}))
	defer srv.Close()

	var out struct {
		Foo string `json:"foo"`
	}
	if err := newTestHTTPClient(srv).do(context.Background(), http.MethodGet, "/test", nil, &out); err != nil {
		t.Fatal(err)
	}
	if out.Foo != "bar" {
		t.Errorf("Foo = %q, want \"bar\"", out.Foo)
	}
}

func TestClientAPIError4xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"code":    "INVALID_REQUEST",
			"message": "bad input",
		})
	}))
	defer srv.Close()

	err := newTestHTTPClient(srv).do(context.Background(), http.MethodGet, "/test", nil, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("err type = %T, want *APIError", err)
	}
	if apiErr.StatusCode != http.StatusBadRequest {
		t.Errorf("StatusCode = %d, want 400", apiErr.StatusCode)
	}
	if apiErr.Code != "INVALID_REQUEST" {
		t.Errorf("Code = %q, want \"INVALID_REQUEST\"", apiErr.Code)
	}
	if apiErr.Message != "bad input" {
		t.Errorf("Message = %q, want \"bad input\"", apiErr.Message)
	}
}

func TestClientAPIError5xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	err := newTestHTTPClient(srv).do(context.Background(), http.MethodGet, "/test", nil, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("err type = %T, want *APIError", err)
	}
	if apiErr.StatusCode != http.StatusInternalServerError {
		t.Errorf("StatusCode = %d, want 500", apiErr.StatusCode)
	}
}

func TestClientNoErrorOn2xx(t *testing.T) {
	for _, status := range []int{200, 201, 204} {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(status)
		}))
		err := newTestHTTPClient(srv).do(context.Background(), http.MethodGet, "/test", nil, nil)
		srv.Close()
		if err != nil {
			t.Errorf("status %d: unexpected error: %v", status, err)
		}
	}
}
