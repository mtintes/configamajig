package actions

import (
	"encoding/json"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

var AppFs = afero.NewOsFs()

func SlurpJson(input string) (interface{}, error) {

	file, err := afero.ReadFile(AppFs, input)

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

func SlurpYaml(input string) (interface{}, error) {

	file, err := afero.ReadFile(AppFs, input)

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

func SlurpGenericFile(inputPath string) ([]byte, error) {
	return afero.ReadFile(AppFs, inputPath)
}

func WriteFile(outputPath string, file string) error {
	return afero.WriteFile(AppFs, outputPath, []byte(file), 0644)
}
