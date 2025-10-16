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

	for name, gc := range generator.Configs {
		gc.SetTarget(cfg.Providers[name].Target)
	}

	for name, srv := range cfg.Services {
		for _, prov := range srv.Providers {
			if gc, ok := generator.Configs[prov]; ok {
				if err := gc.AddService(name, srv); err != nil {
					log.Fatalf("Error adding %s service %s: %v", prov, name, err)
				}
			}
		}
	}

	for _, gc := range generator.Configs {
		if gc.GetTarget() != "" {
			if err := gc.Export(gc.GetTarget()); err != nil {
				log.Fatalf("Error exporting config: %v", err)
			}
		}
	}
}
