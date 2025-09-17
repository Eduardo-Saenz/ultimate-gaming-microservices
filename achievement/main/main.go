package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	achievementController "ultimategaming.com/achievement/internal/controller/achievement"
	httpHandler "ultimategaming.com/achievement/internal/handler/http"
	"ultimategaming.com/achievement/internal/repository/memory"
	"ultimategaming.com/pkg/discovery/consul"
	discovery "ultimategaming.com/pkg/registry"
)

const serviceName = "achievement"

func main() {
	// ===== Config =====
	var (
		port      int
		host      string
		consulAddr string
	)
	flag.IntVar(&port, "port", 8081, "HTTP port for the service")
	flag.StringVar(&host, "host", "127.0.0.1", "Service host/address to publish in Consul")
	flag.StringVar(&consulAddr, "consul", "127.0.0.1:8500", "Consul address (host:port)")
	flag.Parse()

	log.Printf("[%s] starting on %s:%d (consul=%s)", serviceName, host, port, consulAddr)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ===== Registry (wrapper del profe) =====
	registry, err := consul.NewRegistry(consulAddr)
	if err != nil {
		log.Fatalf("registry init error: %v", err)
	}

	instanceID := discovery.GenerateInstanceID(serviceName)
	serviceAddr := fmt.Sprintf("%s:%d", host, port)

	// Registro inicial
	if err := registry.Register(ctx, instanceID, serviceName, serviceAddr); err != nil {
		log.Fatalf("consul register error: %v", err)
	}
	log.Printf("[%s] registered in Consul as %s at %s", serviceName, instanceID, serviceAddr)

	// Heartbeat goroutine
	heartbeatStop := make(chan struct{})
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
					log.Printf("heartbeat error: %v", err)
				}
			case <-heartbeatStop:
				return
			}
		}
	}()

	// ===== Dependencias de la app =====
	repo := memory.New()
	ctrl := achievementController.New(repo)
	h := httpHandler.New(ctrl)

	// ===== Router local =====
	mux := http.NewServeMux()
	h.Register(mux)

	// (Opcional) Healthz para probar en navegador (el check real es TTL)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// ===== Servidor HTTP =====
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	// Arrancar servidor
	go func() {
		log.Printf("[%s] listening on %s", serviceName, srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server error: %v", err)
		}
	}()

	// ===== Graceful shutdown =====
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Printf("[%s] shutting down...", serviceName)

	// Detener heartbeat y deregistrar
	close(heartbeatStop)
	if err := registry.Deregister(ctx, instanceID, serviceName); err != nil {
		log.Printf("deregister error: %v", err)
	} else {
		log.Printf("[%s] deregistered from Consul", serviceName)
	}

	// Apagar HTTP con timeout
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
	log.Printf("[%s] stopped", serviceName)
}
