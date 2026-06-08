package yellowcard

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type httpClient struct {
	apiKey    string
	secretKey string
	baseURL   string
	http      *http.Client
}

func newHTTPClient(apiKey, secretKey string, cfg *config) *httpClient {
	return &httpClient{
		apiKey:    apiKey,
		secretKey: secretKey,
		baseURL:   cfg.baseURL,
		http:      cfg.httpClient,
	}
}

// do executes an authenticated API request and decodes the JSON response into out.
// Pass nil body for GET/DELETE requests. Pass a non-nil value to serialize as JSON.
func (c *httpClient) do(ctx context.Context, method, path string, body, out any) error {
	var (
		rawBody string
		reqBody io.Reader
	)
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("yellowcard: marshal request: %w", err)
		}
		rawBody = string(b)
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return fmt.Errorf("yellowcard: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	for k, v := range authHeaders(c.apiKey, c.secretKey, method, path, rawBody) {
		req.Header.Set(k, v)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("yellowcard: http: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("yellowcard: read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		apiErr := &APIError{StatusCode: resp.StatusCode}
		_ = json.Unmarshal(respBody, apiErr)
		return apiErr
	}

	if out != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, out); err != nil {
			return fmt.Errorf("yellowcard: decode response: %w", err)
		}
	}
	return nil
}

func (c *httpClient) get(ctx context.Context, path string, out any) error {
	return c.do(ctx, http.MethodGet, path, nil, out)
}

func (c *httpClient) post(ctx context.Context, path string, body, out any) error {
	return c.do(ctx, http.MethodPost, path, body, out)
}

func (c *httpClient) put(ctx context.Context, path string, body, out any) error {
	return c.do(ctx, http.MethodPut, path, body, out)
}

func (c *httpClient) delete(ctx context.Context, path string) error {
	return c.do(ctx, http.MethodDelete, path, nil, nil)
}
