package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

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

func (c *Config) Validate() error {
	if c.Targets.Traefik == "" && c.Targets.Nginx == "" {
		return fmt.Errorf("at least one target must be specified")
	}
	if len(c.Sites) == 0 {
		return fmt.Errorf("at least one site must be specified")
	}
	for name, site := range c.Sites {
		if len(site.Domains) == 0 {
			return fmt.Errorf("site %s must have at least one domain", name)
		}
		if site.Provider != "php" && site.Provider != "node" {
			return fmt.Errorf("site %s has invalid provider %s", name, site.Provider)
		}
		if site.PHP != nil {
			if site.PHP.Version == "" {
				return fmt.Errorf("site %s must specify a PHP version", name)
			}
			if site.PHP.Root == "" {
				return fmt.Errorf("site %s must specify a PHP root", name)
			}
		}
	}
	return nil
}

func Parse(data *[]byte) (*Config, error) {
	var config Config
	err := yaml.Unmarshal(*data, &config)
	if err != nil {
		return nil, err
	}
	err = config.Validate()
	if err != nil {
		return nil, err
	}
	return &config, nil
}
