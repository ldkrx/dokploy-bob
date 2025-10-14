package config

import "gopkg.in/yaml.v3"

type Config struct {
	Targets struct {
		Traefik string `yaml:"traefik"`
		Nginx   string `yaml:"nginx"`
	} `yaml:"targets"`
	Sites map[string]Service `yaml:"sites"`
}

type Service struct {
	Domains  []string   `yaml:"domains"`
	Provider string     `yaml:"provider"`
	Port     int        `yaml:"port"`
	PHP      *PHPConfig `yaml:"php,omitempty"`
}

type PHPConfig struct {
	Version string `yaml:"version"`
	Root    string `yaml:"root"`
}

func Parse(data *[]byte) (*Config, error) {
	var config Config
	err := yaml.Unmarshal(*data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
