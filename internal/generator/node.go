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

type NodeServices struct {
	Name        string                 `json:"name"`
	CWD         string                 `json:"cwd,omitempty"`
	Script      string                 `json:"script"`
	Args        []string               `json:"args,omitempty"`
	Interpreter string                 `json:"interpreter,omitempty"`
	PostUpdate  []string               `json:"post_update,omitempty"`
	Env         map[string]interface{} `json:"env"`
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
		env := map[string]interface{}{
			"PORT": svc.Port,
		}
		for k, v := range npc.Env {
			env[k] = v
		}
		script := npc.Script
		interpreter := npc.Interpreter
		if npc.UseNvmrc && npc.CWD != "" {
			script = "source ~/.nvm/nvm.sh && nvm use && " + npc.Script
			interpreter = "/bin/bash"
		}
		service := NodeServices{
			Name:        name,
			Script:      script,
			Args:        npc.Args,
			Interpreter: interpreter,
			CWD:         npc.CWD,
			PostUpdate:  npc.PostUpdate,
			Env:         env,
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
