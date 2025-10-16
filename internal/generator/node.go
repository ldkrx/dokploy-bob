package generator

import (
	"encoding/json"
	"ldriko/dokploy-bob/internal/config"
	"ldriko/dokploy-bob/internal/exporter"
)

type NodeConfig struct {
	Target   string                  `json:"-"`
	Services map[string]NodeServices `json:"apps"`
}

type NodeEnv struct {
	Port int `json:"PORT"`
}

type NodeServices struct {
	Name        string   `json:"name"`
	Script      string   `json:"script"`
	Interpreter string   `json:"interpreter"`
	PostUpdate  []string `json:"post_update,omitempty"`
	Env         *NodeEnv `json:"env"`
}

func NewNodeConfig() *NodeConfig {
	return &NodeConfig{
		Services: make(map[string]NodeServices),
	}
}

func (nc *NodeConfig) SetTarget(target string) {
	nc.Target = target
}

func (nc *NodeConfig) AddService(name string, svc *config.Service) error {
	service := NodeServices{
		Name:        name,
		Script:      svc.Node.Script,
		Interpreter: svc.Node.Interpreter,
		PostUpdate:  svc.Node.PostUpdate,
		Env: &NodeEnv{
			Port: svc.Port,
		},
	}

	nc.Services[name] = service

	return nil
}

func (nc *NodeConfig) Export(path string) error {
	data, err := json.MarshalIndent(nc, "", "  ")
	if err != nil {
		return err
	}

	if err := exporter.Process(path, data); err != nil {
		return err
	}

	return nil
}

func (nc *NodeConfig) GetTarget() string {
	return nc.Target
}
