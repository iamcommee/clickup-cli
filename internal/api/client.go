package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	// BaseURL is the ClickUp API base URL
	BaseURL = "https://api.clickup.com/api/v2"
)

// Client is the ClickUp API client
type Client struct {
	httpClient *http.Client
	baseURL    string
	apiToken   string
	debug      bool
}

// NewClient creates a new API client
func NewClient(apiToken string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL:  BaseURL,
		apiToken: apiToken,
	}
}

// SetDebug enables or disables debug mode
func (c *Client) SetDebug(debug bool) {
	c.debug = debug
}

// doRequest performs an HTTP request
func (c *Client) doRequest(method, path string, query url.Values) ([]byte, error) {
	fullURL := c.baseURL + path
	if query != nil && len(query) > 0 {
		fullURL += "?" + query.Encode()
	}

	if c.debug {
		fmt.Printf("[DEBUG] %s %s\n", method, fullURL)
	}

	req, err := http.NewRequest(method, fullURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if c.debug {
		fmt.Printf("[DEBUG] Response status: %d\n", resp.StatusCode)
		if len(body) < 1000 {
			fmt.Printf("[DEBUG] Response body: %s\n", string(body))
		}
	}

	if resp.StatusCode >= 400 {
		var apiErr APIError
		if err := json.Unmarshal(body, &apiErr); err != nil {
			apiErr = APIError{
				StatusCode: resp.StatusCode,
				Message:    string(body),
			}
		}
		apiErr.StatusCode = resp.StatusCode
		return nil, &apiErr
	}

	return body, nil
}

// Get performs a GET request
func (c *Client) Get(path string, query url.Values) ([]byte, error) {
	return c.doRequest(http.MethodGet, path, query)
}
