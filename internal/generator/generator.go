package generator

import "ldriko/dokploy-bob/internal/config"

type GeneratorInterface interface {
	AddService(name string, svc *config.Service) error
	Export(path string) error
	SetTarget(target string)
	GetTarget() string
}

var Configs = map[string]GeneratorInterface{
	config.Traefik.String(): NewTraefikConfig(),
	config.Nginx.String():   NewNginxConfig(),
	config.Node.String():    NewNodeConfig(),
}
