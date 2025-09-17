package metadata

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	discovery "ultimategaming.com/game/internal/gateway/resolver"
)

type Client interface {
	Get(gameID string) (*MetadataDTO, error)
}

type MetadataDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Genre       string `json:"genre,omitempty"`
	Developer   string `json:"developer,omitempty"`
	ReleaseYear int    `json:"releaseYear,omitempty"`
}

type HTTPClient struct {
	resolver      discovery.Resolver
	serviceName   string // normalmente "metadata"
	httpClient    *http.Client
	basePathMeta  string // "/metadata"
}

func NewHTTPClient(res discovery.Resolver, serviceName string) *HTTPClient {
	if serviceName == "" {
		serviceName = "metadata"
	}
	return &HTTPClient{
		resolver:    res,
		serviceName: serviceName,
		httpClient: &http.Client{
			Timeout: 3 * time.Second,
		},
		basePathMeta: "/metadata",
	}
}

func (c *HTTPClient) Get(gameID string) (*MetadataDTO, error) {
	addr, err := c.resolver.Resolve(c.serviceName)
	if err != nil {
		return nil, fmt.Errorf("resolver error: %w", err)
	}

	u := url.URL{
		Scheme: "http",
		Host:   addr,
		Path:   c.basePathMeta,
	}
	q := u.Query()
	q.Set("id", gameID)
	u.RawQuery = q.Encode()

	req, _ := http.NewRequest(http.MethodGet, u.String(), nil)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("metadata not found for game %s", gameID)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status from metadata: %s", resp.Status)
	}

	var out MetadataDTO
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}
	return &out, nil
}
