package yellowcard

import (
	"testing"
	"time"
)

func TestSign(t *testing.T) {
	ts := time.UnixMilli(1700000000000)
	got := sign("my-secret", "POST", "/v2/business/payments", `{"amount":100}`, ts)
	if got == "" {
		t.Fatal("expected non-empty signature")
	}
	// deterministic: same inputs must produce the same signature
	got2 := sign("my-secret", "POST", "/v2/business/payments", `{"amount":100}`, ts)
	if got != got2 {
		t.Fatalf("sign is not deterministic: %q != %q", got, got2)
	}
	// different secret must produce a different signature
	diff := sign("other-secret", "POST", "/v2/business/payments", `{"amount":100}`, ts)
	if got == diff {
		t.Fatal("different secrets produced the same signature")
	}
}
