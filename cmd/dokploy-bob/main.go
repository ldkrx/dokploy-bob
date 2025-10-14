package main

import (
	"ldriko/dokploy-bob/config"
	"ldriko/dokploy-bob/generator"
	"log"
	"os"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		log.Fatal("Usage: go run main.go <config-file>")
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

	traefikCfg := generator.NewTraefik()

	for name, service := range cfg.Sites {
		err := traefikCfg.AddService(&name, &service)
		if err != nil {
			log.Fatalf("Error adding site %s: %v", name, err)
		}
	}

	traefikYaml, err := traefikCfg.ToYAML()
	if err != nil {
		log.Fatalf("Error converting Traefik config to YAML: %v", err)
	}

	err = generator.Process(cfg.Targets.Traefik, &traefikYaml)
	if err != nil {
		log.Fatalf("Error writing Traefik config to file: %v", err)
	}
}
