package cmd

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	gameController "ultimategaming.com/game/internal/controller/game"
	achievementgateway "ultimategaming.com/game/internal/gateway/achievement/http"
	metadatagateway "ultimategaming.com/game/internal/gateway/metadata/http"
	httphandler "ultimategaming.com/game/internal/handler/http"
	"ultimategaming.com/pkg/discovery/consul"
	discovery "ultimategaming.com/pkg/registry"
)

const serviceName = "game"

func main() {
	var port int
	flag.IntVar(&port, "port", 8083, "API Handler port")
	flag.Parse()
	log.Printf("Starting game service on port: %d", port)
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)
	metadataGateway := metadatagateway.New(registry)
	achievementGateway := achievementgateway.New(registry)
	ctrl := gameController.New(achievementGateway, metadataGateway)
	h := httphandler.New(ctrl)
	http.Handle("/game", http.HandlerFunc(h.GetGameDetails))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}