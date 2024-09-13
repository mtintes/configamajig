package actions

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestFileReaders(t *testing.T) {
	AppFs = afero.NewMemMapFs()

	jsonData := "{\"key\":\"value\"}"
	yamlData := "key: value"
	genericData := "file contents"

	_ = afero.WriteFile(AppFs, "test.json", []byte(jsonData), 0644)
	_ = afero.WriteFile(AppFs, "test.yaml", []byte(yamlData), 0644)
	_ = afero.WriteFile(AppFs, "test.txt", []byte(genericData), 0644)

	jsonFile, err := SlurpJson("test.json")
	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{"key": "value"}, jsonFile)

	yamlFile, err := SlurpYaml("test.yaml")
	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{"key": "value"}, yamlFile)

	genericFile, err := SlurpGenericFile("test.txt")
	assert.Nil(t, err)
	assert.Equal(t, []byte("file contents"), genericFile)

}
