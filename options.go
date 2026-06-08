package yellowcard

import "net/http"

type config struct {
	baseURL    string
	httpClient *http.Client
}

func defaultConfig() *config {
	return &config{
		baseURL:    DefaultBaseURL,
		httpClient: http.DefaultClient,
	}
}

// Option configures the client.
type Option func(*config)

// WithBaseURL overrides the API base URL.
func WithBaseURL(url string) Option {
	return func(c *config) { c.baseURL = url }
}

// WithSandbox points the client at the YellowCard sandbox environment.
func WithSandbox() Option {
	return WithBaseURL(SandboxBaseURL)
}

// WithHTTPClient replaces the default HTTP client (e.g. to set timeouts or add tracing).
func WithHTTPClient(hc *http.Client) Option {
	return func(c *config) { c.httpClient = hc }
}
