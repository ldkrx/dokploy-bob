package config

import "gopkg.in/yaml.v3"

func (s *Service) UnmarshalYAML(value *yaml.Node) error {
	type rawService struct {
		Domains   []string      `yaml:"domains"`
		Providers []interface{} `yaml:"providers"`
		Port      int           `yaml:"port"`
	}
	var raw rawService
	if err := value.Decode(&raw); err != nil {
		return err
	}
	s.Domains = raw.Domains
	s.Port = raw.Port
	for _, item := range raw.Providers {
		if str, ok := item.(string); ok {
			s.Providers = append(s.Providers, ProviderInstance{Name: str})
		} else if m, ok := item.(map[string]interface{}); ok {
			for name, config := range m {
				pi := ProviderInstance{Name: name}
				if ctor, ok := providerConfigFactories[name]; ok {
					cfg := ctor()
					configBytes, _ := yaml.Marshal(config)
					yaml.Unmarshal(configBytes, cfg)
					pi.Config = cfg
				}
				// For traefik, Config remains nil
				s.Providers = append(s.Providers, pi)
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
