package config

import "fmt"

type NginxProviderConfig struct {
	Type    string    `yaml:"type"`
	Root    string    `yaml:"root"`
	Include []string  `yaml:"include,omitempty"`
	PHP     PHPConfig `yaml:"php"`
}

func (npc *NginxProviderConfig) Validate() error {
	if npc.Type != "php" && npc.Type != "static" {
		return fmt.Errorf("nginx provider must have type php or static")
	}
	if npc.Type == "php" && npc.PHP.Version == "" {
		return fmt.Errorf("nginx provider must specify PHP version for type php")
	}
	if npc.Root == "" {
		return fmt.Errorf("nginx provider must specify root")
	}
	return nil
}

type NodeProviderConfig struct {
	Script      string            `yaml:"script"`
	Args        []string          `yaml:"args,omitempty"`
	Interpreter string            `yaml:"interpreter,omitempty"`
	CWD         string            `yaml:"cwd,omitempty"`
	PostUpdate  []string          `yaml:"post_update,omitempty"`
	Env         map[string]string `yaml:"env,omitempty"`
	UseNvmrc    bool              `yaml:"use_nvmrc,omitempty"`
}

func (npc *NodeProviderConfig) Validate() error {
	if npc.Script == "" {
		return fmt.Errorf("node provider must specify script")
	}
	return nil
}

var providerConfigFactories = map[string]func() ProviderConfig{
	"nginx": func() ProviderConfig { return &NginxProviderConfig{} },
	"node":  func() ProviderConfig { return &NodeProviderConfig{} },
}
