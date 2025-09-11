package cmd

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	achievementController "ultimategaming.com/achievement/internal/controller/achievement"
	httpHandler "ultimategaming.com/achievement/internal/handler/http"
	"ultimategaming.com/achievement/internal/repository/memory"
	"ultimategaming.com/pkg/discovery/consul"
	discovery "ultimategaming.com/pkg/registry"
)

const serviceName = "achievement"

func main() {
	var port int
	flag.IntVar(&port, "port", 8082, "API handler port")
	flag.Parse()
	log.Printf("Starting achievement service on port %d", port)
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
	repo := memory.New()
	ctrl := achievementController.New(repo)
	h := httpHandler.New(ctrl)
	http.HandleFunc("/achievement", h.GetAchievement)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}