package generator

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

func getDir(filePath string) string {
	lastSlash := strings.LastIndex(filePath, "/")
	if lastSlash == -1 {
		return "."
	}
	return filePath[:lastSlash]
}

func Process(targetPath string, data *[]byte) error {
	err := os.MkdirAll(getDir(targetPath), 0755)
	if err != nil {
		return err
	}

	return os.WriteFile(targetPath, *data, 0644)
}
