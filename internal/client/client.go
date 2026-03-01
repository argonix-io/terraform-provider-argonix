package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is the Argonix API client.
type Client struct {
	BaseURL        string
	APIKey         string
	OrganizationID string
	HTTPClient     *http.Client
}

// NewClient creates a new Argonix API client.
func NewClient(baseURL, apiKey, organizationID string) *Client {
	return &Client{
		BaseURL:        baseURL,
		APIKey:         apiKey,
		OrganizationID: organizationID,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// orgURL returns the base URL for org-scoped endpoints.
func (c *Client) orgURL() string {
	return fmt.Sprintf("%s/api/0.1/organizations/%s", c.BaseURL, c.OrganizationID)
}

// doRequest performs an HTTP request with auth headers.
func (c *Client) doRequest(ctx context.Context, method, url string, body interface{}) ([]byte, int, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("marshalling request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Api-Key %s", c.APIKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("reading response: %w", err)
	}

	return respBody, resp.StatusCode, nil
}

// Create sends a POST to the given endpoint and decodes the response into result.
func (c *Client) Create(ctx context.Context, endpoint string, payload interface{}, result interface{}) error {
	url := fmt.Sprintf("%s%s", c.orgURL(), endpoint)
	body, status, err := c.doRequest(ctx, http.MethodPost, url, payload)
	if err != nil {
		return err
	}
	if status < 200 || status >= 300 {
		return fmt.Errorf("API returned status %d: %s", status, string(body))
	}
	if result != nil {
		return json.Unmarshal(body, result)
	}
	return nil
}

// Read sends a GET to the given endpoint and decodes the response into result.
func (c *Client) Read(ctx context.Context, endpoint string, result interface{}) error {
	url := fmt.Sprintf("%s%s", c.orgURL(), endpoint)
	body, status, err := c.doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	if status == 404 {
		return &NotFoundError{Endpoint: endpoint}
	}
	if status < 200 || status >= 300 {
		return fmt.Errorf("API returned status %d: %s", status, string(body))
	}
	if result != nil {
		return json.Unmarshal(body, result)
	}
	return nil
}

// Update sends a PUT to the given endpoint and decodes the response into result.
func (c *Client) Update(ctx context.Context, endpoint string, payload interface{}, result interface{}) error {
	url := fmt.Sprintf("%s%s", c.orgURL(), endpoint)
	body, status, err := c.doRequest(ctx, http.MethodPut, url, payload)
	if err != nil {
		return err
	}
	if status < 200 || status >= 300 {
		return fmt.Errorf("API returned status %d: %s", status, string(body))
	}
	if result != nil {
		return json.Unmarshal(body, result)
	}
	return nil
}

// Patch sends a PATCH to the given endpoint and decodes the response into result.
func (c *Client) Patch(ctx context.Context, endpoint string, payload interface{}, result interface{}) error {
	url := fmt.Sprintf("%s%s", c.orgURL(), endpoint)
	body, status, err := c.doRequest(ctx, http.MethodPatch, url, payload)
	if err != nil {
		return err
	}
	if status < 200 || status >= 300 {
		return fmt.Errorf("API returned status %d: %s", status, string(body))
	}
	if result != nil {
		return json.Unmarshal(body, result)
	}
	return nil
}

// Delete sends a DELETE to the given endpoint.
func (c *Client) Delete(ctx context.Context, endpoint string) error {
	url := fmt.Sprintf("%s%s", c.orgURL(), endpoint)
	body, status, err := c.doRequest(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	if status < 200 || status >= 300 {
		return fmt.Errorf("API returned status %d: %s", status, string(body))
	}
	return nil
}

// List sends a GET to the given endpoint and decodes a paginated response.
func (c *Client) List(ctx context.Context, endpoint string, result interface{}) error {
	url := fmt.Sprintf("%s%s", c.orgURL(), endpoint)
	body, status, err := c.doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	if status < 200 || status >= 300 {
		return fmt.Errorf("API returned status %d: %s", status, string(body))
	}
	if result != nil {
		return json.Unmarshal(body, result)
	}
	return nil
}

// NotFoundError is returned when a resource is not found (404).
type NotFoundError struct {
	Endpoint string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("resource not found: %s", e.Endpoint)
}

// IsNotFound checks if an error is a NotFoundError.
func IsNotFound(err error) bool {
	_, ok := err.(*NotFoundError)
	return ok
}
