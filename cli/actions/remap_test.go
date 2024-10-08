package actions

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestRemapCmd(t *testing.T) {

	AppFs = afero.NewMemMapFs()

	config := &ConfigurationMap{
		Configs: []Config{
			{
				Props: map[string]interface{}{
					"key": "value",
				},
				Mappings: []Mapping{},
			},
		},
	}

	RemapCmd(config, "outputPath.yaml", "")

	out, err := afero.ReadFile(AppFs, "outputPath.yaml")

	assert.Nil(t, err)
	assert.NotNil(t, out)
}
