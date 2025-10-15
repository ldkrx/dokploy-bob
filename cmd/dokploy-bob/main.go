package main

import (
	"ldriko/dokploy-bob/internal/config"
	"ldriko/dokploy-bob/internal/generator"
	"log"
	"os"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		log.Fatal("Usage: dokploy-bob <config-file>")
	}

	configFile := args[1]
	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	cfg, err := config.Parse(&data)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	traefikCfg := generator.NewTraefikConfig()
	nginxCfg := generator.NewNginxConfig()

	for name, srv := range cfg.Services {
		for _, prov := range srv.Providers {
			switch prov {
			case config.Traefik.String():
				if err := traefikCfg.AddService(name, srv); err != nil {
					log.Fatalf("Error adding traefik service %s: %v", name, err)
				}
			case config.Nginx.String():
				if err := nginxCfg.AddService(name, srv); err != nil {
					log.Fatalf("Error adding nginx service %s: %v", name, err)
				}
			}
		}
	}

	err = traefikCfg.Export(cfg.Providers.Traefik.Target)
	if err != nil {
		log.Fatalf("Error exporting traefik config: %v", err)
	}

	err = nginxCfg.Export(cfg.Providers.Nginx.Target)
	if err != nil {
		log.Fatalf("Error exporting nginx config: %v", err)
	}
}
