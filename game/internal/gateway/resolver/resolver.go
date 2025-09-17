package discovery

import (
	"fmt"
	"net"

	"github.com/hashicorp/consul/api"
)

type Resolver interface {
	// Resolve devuelve "host:port" de una instancia saludable del servicio dado.
	Resolve(serviceName string) (string, error)
}

type ConsulResolver struct {
	client *api.Client
}

func NewConsulResolver(addr string) (*ConsulResolver, error) {
	cfg := api.DefaultConfig()
	if addr != "" {
		cfg.Address = addr
	}
	c, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &ConsulResolver{client: c}, nil
}

func (r *ConsulResolver) Resolve(serviceName string) (string, error) {
	// Filtramos por checks "passing"
	entries, _, err := r.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return "", err
	}
	if len(entries) == 0 {
		return "", fmt.Errorf("no healthy instances for service %q", serviceName)
	}

	// Estrategia simple: devolvemos la primera (podrías hacer round-robin)
	inst := entries[0].Service

	host := inst.Address
	if host == "" {
		host = entries[0].Node.Address
	}

	if ip := net.ParseIP(host); ip == nil && host == "" {
		return "", fmt.Errorf("invalid host for service %q", serviceName)
	}

	return fmt.Sprintf("%s:%d", host, inst.Port), nil
}