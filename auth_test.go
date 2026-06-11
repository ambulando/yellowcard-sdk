package yellowcard

import (
	"testing"
	"time"
)

func TestSign(t *testing.T) {
	date := time.Now().UTC().Format(time.RFC3339)
	body := map[string]interface{}{
		"amount": 100,
	}
	got := sign("my-secret", "POST", "/v2/business/payments", date, body)
	if got == "" {
		t.Fatal("expected non-empty signature")
	}
	// deterministic: same inputs must produce the same signature
	got2 := sign("my-secret", "POST", "/v2/business/payments", date, body)
	if got != got2 {
		t.Fatalf("sign is not deterministic: %q != %q", got, got2)
	}
	// different secret must produce a different signature
	diff := sign("other-secret", "POST", "/v2/business/payments", date, body)
	if got == diff {
		t.Fatal("different secrets produced the same signature")
	}
}
