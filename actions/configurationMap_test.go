package actions

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestReadVersion1(t *testing.T) {

	AppFs = afero.NewMemMapFs()

	jsonData := `
	{
		"version": "1",
		"configs": [
			{
				"path": "testPath",
				"mappings": [
					{	
						"inPath": "key",
						"toPath": "newKey"
					}
						
				],
				"applyFile": "testFile"			
			},
			{
				"path": "testPath2",
				"mappings": [
					{
						"inPath": "key2",
						"toPath": "newKey2"
					},
					{
						"inPath": "key3",
						"toPath": "newKey3"
					}
				],
				"applyFile": "testFile2"
			}
		]
	}`

	_ = afero.WriteFile(AppFs, "test.json", []byte(jsonData), 0644)

	config, err := readVersion1("test.json")

	if err != nil {
		t.Errorf("Error was not nil")
	}

	if config.Version != "1" {
		t.Errorf("Version was not 1")
	}

	for i := 0; i < 100; i++ {
		assert.Equal(t, "testPath", config.Configs[0].Path)
		assert.Equal(t, "testPath2", config.Configs[1].Path)
	}

}
