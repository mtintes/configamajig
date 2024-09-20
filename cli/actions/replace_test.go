package actions

import (
	"testing"

	"github.com/spf13/afero"
)

func TestReplaceCmd(t *testing.T) {

	AppFs = afero.NewMemMapFs()

	jsonData := "{\"key\":\"{{if .key2 }}key{{else}}no_key{{end}}\"}"
	yamlData := "key2: value"

	_ = afero.WriteFile(AppFs, "test.json", []byte(jsonData), 0644)
	_ = afero.WriteFile(AppFs, "test.yaml", []byte(yamlData), 0644)

	ConfigurationMap := &ConfigurationMap{
		Version: "1.0.0",
		Configs: []Config{
			{
				Path:      "test.json",
				Mappings:  []Mapping{},
				ApplyFile: "after",
			},
			{
				Path:      "test.yaml",
				Mappings:  []Mapping{},
				ApplyFile: "after",
			},
		},
	}

	inputFile := "I am {{ .key }}"
	outputFile := "I am key"

	_ = afero.WriteFile(AppFs, "test.txt", []byte(inputFile), 0644)

	ReplaceCmd(ConfigurationMap, "test.txt", "output.txt", "")

	file, _ := afero.ReadFile(AppFs, "output.txt")

	if string(file) != outputFile {
		t.Errorf("Expected %s, got %s", outputFile, string(file))
	}

}
