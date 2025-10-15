package exporter

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func MarshalToYAML(v interface{}) ([]byte, error) {
	var sb strings.Builder
	encoder := yaml.NewEncoder(&sb)
	encoder.SetIndent(2)
	err := encoder.Encode(v)
	if err != nil {
		return nil, err
	}
	encoder.Close()
	return []byte(sb.String()), nil
}

func GetDir(filePath string) string {
	lastSlash := strings.LastIndex(filePath, "/")
	if lastSlash == -1 {
		return "."
	}
	return filePath[:lastSlash]
}

func EnsureDir(path string) error {
	err := os.MkdirAll(GetDir(path), 0755)
	if err != nil {
		return err
	}

	return nil
}

func Process(targetPath string, data []byte) error {
	EnsureDir(targetPath)
	return os.WriteFile(targetPath, data, 0644)
}
