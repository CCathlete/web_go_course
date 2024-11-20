package views

import (
	"fmt"
	"os"
	"webGo/templates"

	"github.com/go-yaml/yaml"
)

func ParseYaml(yamlPath string) (any, error) {
	file, err := templates.FS.Open("QA.yaml")
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

func ParseYamlNoEmbedding(yamlPath string) (any, error) {
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
