package yellowcard

import (
	"errors"
	"strings"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	e := &APIError{StatusCode: 404, Code: "NOT_FOUND", Message: "resource not found"}
	got := e.Error()
	if !strings.Contains(got, "404") {
		t.Errorf("Error() = %q, want it to contain status code", got)
	}
	if !strings.Contains(got, "NOT_FOUND") {
		t.Errorf("Error() = %q, want it to contain code", got)
	}
	if !strings.Contains(got, "resource not found") {
		t.Errorf("Error() = %q, want it to contain message", got)
	}
}

func TestIsNotFound_true(t *testing.T) {
	if !IsNotFound(&APIError{StatusCode: 404}) {
		t.Error("IsNotFound(404) = false, want true")
	}
}

func TestIsNotFound_false(t *testing.T) {
	if IsNotFound(&APIError{StatusCode: 400}) {
		t.Error("IsNotFound(400) = true, want false")
	}
}

func TestIsNotFound_nil(t *testing.T) {
	if IsNotFound(nil) {
		t.Error("IsNotFound(nil) = true, want false")
	}
}

func TestIsNotFound_nonAPIError(t *testing.T) {
	if IsNotFound(errors.New("some error")) {
		t.Error("IsNotFound(non-APIError) = true, want false")
	}
}

func TestIsUnauthorized_true(t *testing.T) {
	if !IsUnauthorized(&APIError{StatusCode: 401}) {
		t.Error("IsUnauthorized(401) = false, want true")
	}
}

func TestIsUnauthorized_false(t *testing.T) {
	if IsUnauthorized(&APIError{StatusCode: 403}) {
		t.Error("IsUnauthorized(403) = true, want false")
	}
}

func TestIsUnauthorized_nil(t *testing.T) {
	if IsUnauthorized(nil) {
		t.Error("IsUnauthorized(nil) = true, want false")
	}
}

func TestIsUnauthorized_nonAPIError(t *testing.T) {
	if IsUnauthorized(errors.New("some error")) {
		t.Error("IsUnauthorized(non-APIError) = true, want false")
	}
}
