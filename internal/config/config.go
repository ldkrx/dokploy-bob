package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type AvailableProviders int

const (
	Traefik AvailableProviders = iota
	Nginx
	Node
)

var providerNames = map[AvailableProviders]string{
	Traefik: "traefik",
	Nginx:   "nginx",
	Node:    "node",
}

func (prov AvailableProviders) String() string {
	return providerNames[prov]
}

type TraefikProvider struct {
	Target string `yaml:"target"`
}

type NginxProvider struct {
	Target string `yaml:"target"`
}

type Config struct {
	Providers struct {
		Traefik *TraefikProvider `yaml:"traefik,omitempty"`
		Nginx   *NginxProvider   `yaml:"nginx,omitempty"`
	} `yaml:"providers"`
	Services map[string]Service `yaml:"services"`
}

type Service struct {
	Domains   []string   `yaml:"domains"`
	Providers []string   `yaml:"providers"`
	Port      int        `yaml:"port"`
	PHP       *PHPConfig `yaml:"php,omitempty"`
}

func (s *Service) UnmarshalYAML(value *yaml.Node) error {
	type rawService struct {
		Domains   []string   `yaml:"domains"`
		Providers []string   `yaml:"providers"`
		Port      int        `yaml:"port"`
		PHP       *PHPConfig `yaml:"php,omitempty"`
	}
	var raw rawService
	if err := value.Decode(&raw); err != nil {
		return err
	}
	s.Domains = raw.Domains
	s.Port = raw.Port
	s.PHP = raw.PHP
	s.Providers = raw.Providers
	return nil
}

type PHPConfig struct {
	Version string `yaml:"version"`
	Root    string `yaml:"root"`
}

func (c *Config) Validate() error {
	if c.Providers.Nginx == nil && c.Providers.Traefik == nil {
		return fmt.Errorf("at least one target must be specified")
	}
	if len(c.Services) == 0 {
		return fmt.Errorf("at least one site must be specified")
	}
	for name, serv := range c.Services {
		if len(serv.Domains) == 0 {
			return fmt.Errorf("site %s must have at least one domain", name)
		}

		for _, prov := range serv.Providers {
			valid := false
			for _, v := range providerNames {
				if prov == v {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("site %s has an invalid provider: %s", name, prov)
			}
		}

		if serv.PHP != nil {
			if serv.PHP.Version == "" {
				return fmt.Errorf("site %s must specify a PHP version", name)
			}
			if serv.PHP.Root == "" {
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
