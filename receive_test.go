package yellowcard

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var samplePayment = Payment{
	Id:       "pay-123",
	Status:   "pending",
	Amount:   100,
	Currency: "USD",
	Country:  "GH",
}

func TestPaymentsService_Create(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/business/receive" {
			t.Errorf("path = %s, want /business/receive", r.URL.Path)
		}
		var body ReceivePaymentRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.SequenceId != "seq-1" {
			t.Errorf("SequenceId = %q, want \"seq-1\"", body.SequenceId)
		}
		json.NewEncoder(w).Encode(samplePayment)
	}))
	defer srv.Close()

	svc := &PaymentsService{client: newTestHTTPClient(srv)}
	p, err := svc.Create(context.Background(), ReceivePaymentRequest{SequenceId: "seq-1", Amount: 100, Currency: "USD"})
	if err != nil {
		t.Fatal(err)
	}
	if p.Id != "pay-123" {
		t.Errorf("Id = %q, want \"pay-123\"", p.Id)
	}
	if p.Status != "pending" {
		t.Errorf("Status = %q, want \"pending\"", p.Status)
	}
}

func TestPaymentsService_Create_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{"code": "INVALID", "message": "invalid request"})
	}))
	defer srv.Close()

	svc := &PaymentsService{client: newTestHTTPClient(srv)}
	_, err := svc.Create(context.Background(), ReceivePaymentRequest{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !isAPIError(err, http.StatusUnprocessableEntity) {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestPaymentsService_Get(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/business/receive/pay-123" {
			t.Errorf("path = %s, want /business/receive/pay-123", r.URL.Path)
		}
		json.NewEncoder(w).Encode(samplePayment)
	}))
	defer srv.Close()

	svc := &PaymentsService{client: newTestHTTPClient(srv)}
	p, err := svc.Get(context.Background(), "pay-123")
	if err != nil {
		t.Fatal(err)
	}
	if p.Id != "pay-123" {
		t.Errorf("Id = %q, want \"pay-123\"", p.Id)
	}
}

func TestPaymentsService_Get_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"code": "NOT_FOUND", "message": "not found"})
	}))
	defer srv.Close()

	svc := &PaymentsService{client: newTestHTTPClient(srv)}
	_, err := svc.Get(context.Background(), "missing")
	if !IsNotFound(err) {
		t.Errorf("expected IsNotFound, got: %v", err)
	}
}

func TestPaymentsService_GetBySequenceID(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/seq-abc") {
			t.Errorf("path = %s, want suffix /seq-abc", r.URL.Path)
		}
		json.NewEncoder(w).Encode(samplePayment)
	}))
	defer srv.Close()

	svc := &PaymentsService{client: newTestHTTPClient(srv)}
	p, err := svc.GetBySequenceID(context.Background(), "seq-abc")
	if err != nil {
		t.Fatal(err)
	}
	if p.Id != "pay-123" {
		t.Errorf("Id = %q, want \"pay-123\"", p.Id)
	}
}

func TestPaymentsService_List(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/business/receive" {
			t.Errorf("path = %s, want /business/receive", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("perPage") != "10" {
			t.Errorf("perPage = %q, want \"10\"", q.Get("perPage"))
		}
		if q.Get("orderBy") != "desc" {
			t.Errorf("orderBy = %q, want \"desc\"", q.Get("orderBy"))
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"collections": []Payment{samplePayment},
		})
	}))
	defer srv.Close()

	svc := &PaymentsService{client: newTestHTTPClient(srv)}
	payments, err := svc.List(context.Background(), ReceiveQuery{PerPage: 10, OrderBy: Desc})
	if err != nil {
		t.Fatal(err)
	}
	if len(payments) != 1 {
		t.Fatalf("len = %d, want 1", len(payments))
	}
	if payments[0].Id != "pay-123" {
		t.Errorf("Id = %q, want \"pay-123\"", payments[0].Id)
	}
}

func TestPaymentsService_List_NoQuery(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "" {
			t.Errorf("unexpected query string: %s", r.URL.RawQuery)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"collections": []Payment{}})
	}))
	defer srv.Close()

	svc := &PaymentsService{client: newTestHTTPClient(srv)}
	_, err := svc.List(context.Background(), ReceiveQuery{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestPaymentsService_Accept(t *testing.T) {
	testPaymentAction(t, "accept", func(svc *PaymentsService) (*Payment, error) {
		return svc.Accept(context.Background(), "pay-123")
	})
}

func TestPaymentsService_Deny(t *testing.T) {
	testPaymentAction(t, "deny", func(svc *PaymentsService) (*Payment, error) {
		return svc.Deny(context.Background(), "pay-123")
	})
}

func TestPaymentsService_Cancel(t *testing.T) {
	testPaymentAction(t, "cancel", func(svc *PaymentsService) (*Payment, error) {
		return svc.Cancel(context.Background(), "pay-123")
	})
}

func TestPaymentsService_Refund(t *testing.T) {
	testPaymentAction(t, "refund", func(svc *PaymentsService) (*Payment, error) {
		return svc.Refund(context.Background(), "pay-123")
	})
}

// testPaymentAction verifies a POST action on /business/receive/{id}/{action}.
func testPaymentAction(t *testing.T, action string, fn func(*PaymentsService) (*Payment, error)) {
	t.Helper()
	wantPath := "/business/receive/pay-123/" + action
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != wantPath {
			t.Errorf("path = %s, want %s", r.URL.Path, wantPath)
		}
		json.NewEncoder(w).Encode(samplePayment)
	}))
	defer srv.Close()

	svc := &PaymentsService{client: newTestHTTPClient(srv)}
	p, err := fn(svc)
	if err != nil {
		t.Fatalf("%s: %v", action, err)
	}
	if p.Id != "pay-123" {
		t.Errorf("%s: Id = %q, want \"pay-123\"", action, p.Id)
	}
}

func isAPIError(err error, status int) bool {
	e, ok := err.(*APIError)
	return ok && e.StatusCode == status
}
