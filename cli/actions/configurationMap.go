package actions

import (
	"encoding/json"
	"errors"

	"github.com/spf13/afero"
)

type ConfigurationMap struct {
	Version string   `json:"version" yaml:"version"`
	Configs []Config `json:"configs" yaml:"configs"`
}

type Config struct {
	Props     map[string]interface{} `json:"props" yaml:"props"`
	Path      string                 `json:"path" yaml:"path"`
	Mappings  []Mapping              `json:"mappings" yaml:"mappings"`
	ApplyFile string                 `json:"applyFile" yaml:"applyFile"`
}

type Mapping struct {
	InPath string `json:"inPath" yaml:"inPath"`
	ToPath string `json:"toPath" yaml:"toPath"`
}

func ReadConfigurationMap(input string) (*ConfigurationMap, error) {
	version, _ := findVersion(input)

	switch version {
	case "1.0.0":
		return readVersion1(input)
	default:
		return nil, errors.New("config file version not supported")
	}
}

func findVersion(input string) (string, error) {

	file, err := afero.ReadFile(AppFs, input)

	if err != nil {
		return "", err
	}

	var data map[string]interface{}
	err = json.Unmarshal(file, &data)

	if err != nil {
		return "", err
	}

	version := data["version"].(string)

	return version, nil

}

func readVersion1(input string) (*ConfigurationMap, error) {

	file, err := afero.ReadFile(AppFs, input)

	if err != nil {
		return nil, err
	}

	var data ConfigurationMap
	err = json.Unmarshal(file, &data)

	if err != nil {
		return nil, err
	}

	return &data, nil

}
