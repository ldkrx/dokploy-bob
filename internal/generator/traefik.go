package generator

import (
	"fmt"
	"ldriko/dokploy-bob/internal/config"
	"strings"
)

type TraefikConfig struct {
	HTTP HTTPConfig `yaml:"http"`
}

type HTTPConfig struct {
	Routers  map[string]Router  `yaml:"routers"`
	Services map[string]Service `yaml:"services"`
}

type Router struct {
	Rule        string   `yaml:"rule"`
	EntryPoints []string `yaml:"entryPoints"`
	Service     string   `yaml:"service"`
	TLS         *TLS     `yaml:"tls,omitempty"`
}

type TLS struct {
	CertResolver string `yaml:"certResolver"`
}

type Service struct {
	LoadBalancer LoadBalancer `yaml:"loadBalancer"`
}

type LoadBalancer struct {
	Servers []Server `yaml:"servers"`
}

type Server struct {
	URL string `yaml:"url"`
}

func NewTraefik() *TraefikConfig {
	return &TraefikConfig{
		HTTP: HTTPConfig{
			Routers:  make(map[string]Router),
			Services: make(map[string]Service),
		},
	}
}

func (tc *TraefikConfig) AddService(name *string, svc *config.Service) error {
	var rules []string
	for _, domain := range svc.Domains {
		rules = append(rules, "Host(`"+domain+"`)")
	}

	router := Router{
		Rule:        strings.Join(rules, " || "),
		EntryPoints: []string{"web", "websecure"},
		Service:     *name,
		TLS:         &TLS{CertResolver: "letsencrypt"},
	}

	service := Service{
		LoadBalancer: LoadBalancer{
			Servers: []Server{
				{URL: "http://" + *name + ":" + fmt.Sprint(svc.Port)},
			},
		},
	}

	tc.HTTP.Routers[*name] = router
	tc.HTTP.Services[*name] = service

	return nil
}

func (tc *TraefikConfig) ToYAML() ([]byte, error) {
	data, err := MarshalToYAML(tc)
	if err != nil {
		return nil, err
	}

	return data, nil
}
