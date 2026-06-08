package yellowcard

import "fmt"

// APIError represents an error response from the YellowCard API.
type APIError struct {
	StatusCode int
	Code       string `json:"code"`
	Message    string `json:"message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("yellowcard: %d %s — %s", e.StatusCode, e.Code, e.Message)
}

// IsNotFound reports whether the error is a 404 from the API.
func IsNotFound(err error) bool {
	var e *APIError
	return asAPIError(err, &e) && e.StatusCode == 404
}

// IsUnauthorized reports whether the error is a 401 from the API.
func IsUnauthorized(err error) bool {
	var e *APIError
	return asAPIError(err, &e) && e.StatusCode == 401
}

func asAPIError(err error, target **APIError) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*APIError); ok {
		*target = e
		return true
	}
	return false
}
