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

	for name, service := range cfg.Sites {
		err := traefikCfg.AddService(name, service)
		if err != nil {
			log.Fatalf("Error adding traefik service %s: %v", name, err)
		}

		err = nginxCfg.AddService(name, service)
		if err != nil {
			log.Fatalf("Error adding nginx service %s: %v", name, err)
		}
	}

	traefikCfg.Export(cfg.Targets.Traefik)
	nginxCfg.Export(cfg.Targets.Nginx)
}
