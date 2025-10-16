package config

type ProviderType int

const (
	Traefik ProviderType = iota
	Nginx
	Node
)

var providerTypeNames = map[ProviderType]string{
	Traefik: "traefik",
	Nginx:   "nginx",
	Node:    "node",
}

func (pt ProviderType) String() string {
	return providerTypeNames[pt]
}

type Provider struct {
	Target string `yaml:"target"`
}

type Config struct {
	Providers map[string]*Provider `yaml:"providers"`
	Services  map[string]*Service  `yaml:"services"`
}

type PHPConfig struct {
	Version string `yaml:"version"`
	Root    string `yaml:"root"`
}

type ProviderConfig interface {
	Validate() error
}

type ProviderInstance struct {
	Name   string
	Config ProviderConfig
}

type Service struct {
	Domains   []string           `yaml:"domains"`
	Providers []ProviderInstance `yaml:"providers"`
	Port      int                `yaml:"port"`
}
