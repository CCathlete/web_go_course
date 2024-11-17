package views

import (
	"fmt"
	"os"

	"github.com/go-yaml/yaml"
)

func ParseYaml(yamlPath string) (any, error) {
	file, err := os.Open(yamlPath)
	if err != nil {
		return nil, fmt.Errorf("error while opening the file at %s: %w", yamlPath, err)
	}
	defer file.Close()

	var data interface{}
	err = yaml.NewDecoder(file).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("error while reading the file at %s: %w", yamlPath, err)
	}

	return data, nil
}
