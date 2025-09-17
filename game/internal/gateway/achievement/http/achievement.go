package achievement

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	discovery "ultimategaming.com/game/internal/gateway/resolver"
)

type Client interface {
	ListByGame(gameID string) ([]AchievementDTO, error)
}

type AchievementDTO struct {
	ID          string `json:"id"`
	GameID      string `json:"gameId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Points      int    `json:"points"`
	Secret      bool   `json:"secret"`
}

type HTTPClient struct {
	resolver           discovery.Resolver
	serviceName        string // normalmente "achievement"
	httpClient         *http.Client
	basePathAchievements string // "/achievements"
}

func NewHTTPClient(res discovery.Resolver, serviceName string) *HTTPClient {
	if serviceName == "" {
		serviceName = "achievement"
	}
	return &HTTPClient{
		resolver:    res,
		serviceName: serviceName,
		httpClient: &http.Client{
			Timeout: 3 * time.Second,
		},
		basePathAchievements: "/achievements",
	}
}

func (c *HTTPClient) ListByGame(gameID string) ([]AchievementDTO, error) {
	addr, err := c.resolver.Resolve(c.serviceName)
	if err != nil {
		return nil, fmt.Errorf("resolver error: %w", err)
	}

	u := url.URL{
		Scheme: "http",
		Host:   addr,
		Path:   c.basePathAchievements,
	}
	q := u.Query()
	q.Set("gameId", gameID)
	u.RawQuery = q.Encode()

	req, _ := http.NewRequest(http.MethodGet, u.String(), nil)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("achievements not found for game %s", gameID)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status from achievement: %s", resp.Status)
	}

	var out []AchievementDTO
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}
	return out, nil
}
