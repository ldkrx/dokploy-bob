package generator

import (
	"fmt"
	"ldriko/dokploy-bob/internal/config"
)

type NginxConfig struct {
	Services map[string]NginxService
}

type NginxService struct {
	ServerName []string
	PHP        struct {
		Version string
		Root    string
	}
	AccessLog string
	ErrorLog  string
}

func NewNginxConfig() *NginxConfig {
	return &NginxConfig{
		Services: make(map[string]NginxService),
	}
}

func (nc *NginxConfig) AddService(name string, svc config.Service) error {
	service := NginxService{
		ServerName: svc.Domains,
		AccessLog:  fmt.Sprintf("/var/log/nginx/%s.access.log", name),
		ErrorLog:   fmt.Sprintf("/var/log/nginx/%s.error.log", name),
	}

	service.PHP.Version = svc.PHP.Version
	service.PHP.Root = svc.PHP.Root

	nc.Services[name] = service

	return nil
}

func (nc *NginxConfig) Export(path string) error {
	// todo: implement nginx config parsing

	return nil
}
