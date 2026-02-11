package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const defaultBaseURL = "https://api.honeycomb.io"

// DefaultTimeout is the default HTTP client timeout.
const DefaultTimeout = 30 * time.Second

type Client struct {
	APIKey  string
	BaseURL string
	HTTP    *http.Client
}

func NewClient(apiKey string, timeout time.Duration) *Client {
	return &Client{
		APIKey:  apiKey,
		BaseURL: defaultBaseURL,
		HTTP:    &http.Client{Timeout: timeout},
	}
}

func (c *Client) doBearer(req *http.Request, out any) error {
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	return c.doRequest(req, out)
}

func (c *Client) do(req *http.Request, out any) error {
	req.Header.Set("X-Honeycomb-Team", c.APIKey)
	return c.doRequest(req, out)
}

func (c *Client) doJSON(req *http.Request, out any) error {
	req.Header.Set("X-Honeycomb-Team", c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	return c.doRequest(req, out)
}

func (c *Client) doRequest(req *http.Request, out any) error {
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("API error (HTTP %d): (unreadable body)", resp.StatusCode)
		}
		return fmt.Errorf("API error (HTTP %d): %s", resp.StatusCode, string(body))
	}

	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return fmt.Errorf("decoding response: %w", err)
		}
	}
	return nil
}

// Generic CRUD helpers

// List retrieves a list of resources at the given path.
func List[T any](c *Client, ctx context.Context, path string) ([]T, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+path, nil)
	if err != nil {
		return nil, err
	}

	var result []T
	if err := c.do(req, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// Get retrieves a single resource at the given path.
func Get[T any](c *Client, ctx context.Context, path string) (*T, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+path, nil)
	if err != nil {
		return nil, err
	}

	var result T
	if err := c.do(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Create posts a new resource at the given path and returns the created resource.
func Create[T any](c *Client, ctx context.Context, path string, body *T) (*T, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("encoding request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	var result T
	if err := c.doJSON(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Update puts a resource at the given path and returns the updated resource.
func Update[T any](c *Client, ctx context.Context, path string, body *T) (*T, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("encoding request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.BaseURL+path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	var result T
	if err := c.doJSON(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete removes the resource at the given path.
func Delete(c *Client, ctx context.Context, path string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.BaseURL+path, nil)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}
