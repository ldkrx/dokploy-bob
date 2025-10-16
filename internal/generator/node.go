package generator

import (
	"encoding/json"
	"ldriko/dokploy-bob/internal/config"
	"ldriko/dokploy-bob/internal/exporter"
)

type NodeConfig struct {
	Target   string         `json:"-"`
	Services []NodeServices `json:"apps"`
}

type NodeEnv struct {
	Port int `json:"PORT"`
}

type NodeServices struct {
	Name        string   `json:"name"`
	Script      string   `json:"script"`
	Args        []string `json:"args,omitempty"`
	Interpreter string   `json:"interpreter"`
	PostUpdate  []string `json:"post_update,omitempty"`
	Env         *NodeEnv `json:"env"`
}

func NewNodeConfig() *NodeConfig {
	return &NodeConfig{
		Services: []NodeServices{},
	}
}

func (nc *NodeConfig) SetTarget(target string) {
	nc.Target = target
}

func (nc *NodeConfig) AddService(name string, svc *config.Service, pi config.ProviderInstance) error {
	if npc, ok := pi.Config.(*config.NodeProviderConfig); ok {
		service := NodeServices{
			Name:        name,
			Script:      npc.Script,
			Args:        npc.Args,
			Interpreter: npc.Interpreter,
			PostUpdate:  npc.PostUpdate,
			Env: &NodeEnv{
				Port: svc.Port,
			},
		}

		nc.Services = append(nc.Services, service)
	}

	return nil
}

func (nc *NodeConfig) Export(path string) error {
	data, err := json.MarshalIndent(nc, "", "    ")
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
