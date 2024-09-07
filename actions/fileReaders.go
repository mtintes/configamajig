package actions

import (
	"encoding/json"
	"os"

	"gopkg.in/yaml.v3"
)

func slurpJson(input string) (interface{}, error) {

	file, err := os.ReadFile(input)

	if err != nil {
		return nil, err
	}

	var data interface{}
	err = json.Unmarshal(file, &data)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func slurpYaml(input string) (interface{}, error) {

	file, err := os.ReadFile(input)

	if err != nil {
		return nil, err
	}

	var data interface{}
	err = yaml.Unmarshal(file, &data)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func slurpGenericFile(inputPath string) ([]byte, error) {
	return os.ReadFile(inputPath)
}

func WriteFile(outputPath string, file string) error {
	return os.WriteFile(outputPath, []byte(file), 0644)
}