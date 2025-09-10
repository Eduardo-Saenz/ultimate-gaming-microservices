package achievementgateway

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	achievementModel "ultimategaming.com/achievement/pkg/model"
	"ultimategaming.com/game/internal/gateway"
	discovery "ultimategaming.com/pkg/registry"
)

type Gateway struct {
	registry discovery.Registry
}

func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry: registry}
}

// ListByRecord obtiene la lista de logros para un juego (recordID, recordType=game).
func (g *Gateway) ListByRecord(
	ctx context.Context,
	recordID achievementModel.RecordID,
	recordType achievementModel.RecordType,
) ([]achievementModel.Achievement, error) {

	addrs, err := g.registry.ServiceAddress(ctx, "achievement")
	if err != nil {
		return nil, err
	}
	// Endpoint típico estilo REST: GET /achievements?recordId=...&recordType=game
	url := "http://" + addrs[rand.Intn(len(addrs))] + "/achievements"

	log.Printf("Calling achievement service, request: GET %s", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	q := req.URL.Query()
	q.Add("recordId", string(recordID))
	q.Add("recordType", string(recordType))
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch {
	case resp.StatusCode == http.StatusNotFound:
		return nil, gateway.ErrNotFound
	case resp.StatusCode/100 != 2:
		return nil, fmt.Errorf("non-2xx response from achievements: %v", resp.Status)
	}

	var list []achievementModel.Achievement
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, err
	}
	return list, nil
}