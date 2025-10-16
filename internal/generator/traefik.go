package generator

import (
	"fmt"
	"ldriko/dokploy-bob/internal/config"
	"ldriko/dokploy-bob/internal/exporter"
	"strings"
)

type TraefikConfig struct {
	Target string            `yaml:"-"`
	HTTP   TraefikHTTPConfig `yaml:"http"`
}

type TraefikHTTPConfig struct {
	Routers  map[string]TraefikRouter  `yaml:"routers"`
	Services map[string]TraefikService `yaml:"services"`
}

type TraefikRouter struct {
	Rule        string      `yaml:"rule"`
	EntryPoints []string    `yaml:"entryPoints"`
	Service     string      `yaml:"service"`
	TLS         *TraefikTLS `yaml:"tls,omitempty"`
}

type TraefikTLS struct {
	CertResolver string `yaml:"certResolver"`
}

type TraefikService struct {
	LoadBalancer TraefikLoadBalancer `yaml:"loadBalancer"`
}

type TraefikLoadBalancer struct {
	Servers []TraefikServer `yaml:"servers"`
}

type TraefikServer struct {
	URL string `yaml:"url"`
}

func NewTraefikConfig() *TraefikConfig {
	return &TraefikConfig{
		HTTP: TraefikHTTPConfig{
			Routers:  make(map[string]TraefikRouter),
			Services: make(map[string]TraefikService),
		},
	}
}

func (tc *TraefikConfig) SetTarget(target string) {
	tc.Target = target
}

func (tc *TraefikConfig) GetTarget() string {
	return tc.Target
}

func (tc *TraefikConfig) AddService(name string, svc *config.Service) error {
	var rules []string
	for _, domain := range svc.Domains {
		rules = append(rules, "Host(`"+domain+"`)")
	}

	router := TraefikRouter{
		Rule:        strings.Join(rules, " || "),
		EntryPoints: []string{"web", "websecure"},
		Service:     name,
		TLS:         &TraefikTLS{CertResolver: "letsencrypt"},
	}

	service := TraefikService{
		LoadBalancer: TraefikLoadBalancer{
			Servers: []TraefikServer{
				{URL: "http://172.17.0.1:" + fmt.Sprint(svc.Port)},
			},
		},
	}

	tc.HTTP.Routers[name] = router
	tc.HTTP.Services[name] = service

	return nil
}

func (tc *TraefikConfig) ToYAML() ([]byte, error) {
	data, err := exporter.MarshalToYAML(tc)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (tc *TraefikConfig) Export(path string) error {
	traefikYaml, err := tc.ToYAML()
	if err != nil {
		return err
	}

	err = exporter.Process(path, traefikYaml)
	if err != nil {
		return err
	}

	return nil
}
