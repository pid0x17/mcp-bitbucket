package bitbucket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HTTPClient struct {
	baseURL    string
	httpClient *http.Client
	username   string
	password   string
	token      string
}

// ---------------------------------------------------------
// Functional Client Options
// ---------------------------------------------------------

// Option represents a configuration function for HTTPClient.
type Option func(*HTTPClient)

func WithBaseURL(url string) Option {
	return func(c *HTTPClient) {
		c.baseURL = url
	}
}

func WithAuth(username, password string) Option {
	return func(c *HTTPClient) {
		c.username = username
		c.password = password
	}
}

func WithToken(token string) Option {
	return func(c *HTTPClient) {
		c.token = token
	}
}

func WithHTTPClient(client *http.Client) Option {
	return func(c *HTTPClient) {
		c.httpClient = client
	}
}

func NewHTTPClient(opts ...Option) *HTTPClient {
	// sane defaults
	c := &HTTPClient{
		baseURL: "https://api.bitbucket.org/2.0",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// ---------------------------------------------------------
// DTO
// ---------------------------------------------------------

type repoDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPrivate   bool   `json:"is_private"`
}

func (c *HTTPClient) GetRepository(ctx context.Context, workspace string, repoSlug string) (Repository, error) {
	url := fmt.Sprintf("%s/repositories/%s/%s", c.baseURL, workspace, repoSlug)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Repository{}, fmt.Errorf("failed to create request: %w", err)
	}

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	} else if c.username != "" && c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Repository{}, fmt.Errorf("http request failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Repository{}, fmt.Errorf("bitbucket API returned status: %d", resp.StatusCode)
	}

	var dto repoDTO
	if err := json.NewDecoder(resp.Body).Decode(&dto); err != nil {
		return Repository{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return Repository{
		Name:        dto.Name,
		Description: dto.Description,
		IsPrivate:   dto.IsPrivate,
	}, nil
}
