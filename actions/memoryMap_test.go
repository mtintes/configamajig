package actions

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestMappings(t *testing.T) {

	masterMemoryMap := make(map[string]interface{})
	flatFile := make(map[string]interface{})
	flatFile["key"] = "value"
	flatFile["key2"] = make(map[string]interface{})
	flatFile["key2"].(map[string]interface{})["key3"] = "value3"

	mappings := []Mapping{
		{
			InPath: "key",
			ToPath: "newKey",
		},
		{
			InPath: "key2",
			ToPath: "newKey2",
		},
	}

	traces := []Trace{}
	masterMemoryMap, flatFile, _ = applyMappings(flatFile, mappings, masterMemoryMap, traces)

	assert.Equal(t, false, SafePropertyCheck(masterMemoryMap, "key"))
	assert.Equal(t, false, SafePropertyCheck(masterMemoryMap, "key2"))
	assert.Equal(t, true, SafePropertyCheck(masterMemoryMap, "newKey"))
	assert.Equal(t, true, SafePropertyCheck(masterMemoryMap, "newKey2"))
	assert.Equal(t, 0, len(flatFile))

	assert.Equal(t, "value", masterMemoryMap["newKey"])
	assert.Equal(t, "value3", masterMemoryMap["newKey2"].(map[string]interface{})["key3"])

}

func TestFileRestore(t *testing.T) {

	storedFiles := []StoredMemoryMap{
		{
			File: map[string]interface{}{
				"key": "value",
			},
			FileName: "testFile",
		},
	}

	file := returnStoredFile(storedFiles, "testFile")

	assert.Equal(t, map[string]interface{}{"key": "value"}, file)
}

func TestReadMemoryMap(t *testing.T) {

	AppFs = afero.NewMemMapFs()

	jsonData := "{\"key\":\"valuejson\", \"key2\": {\"key3\": \"value3\", \"key5\": \"value5\"}}"
	_ = afero.WriteFile(AppFs, "test.json", []byte(jsonData), 0644)

	yamlData := ` 
key: valueyaml
key2: 
  key3: valueyaml3
  key4: value4`

	_ = afero.WriteFile(AppFs, "test.yaml", []byte(yamlData), 0644)

	expectedMemoryMap := map[string]interface{}{
		"key": "valuejson",
		"key2": map[string]interface{}{
			"key3": "valueyaml3",
			"key4": "value4",
			"key5": "value5",
		},
		"newKey": "valueyaml",
	}

	ConfigurationMap := ConfigurationMap{
		Configs: []Config{
			{
				Path: "test.json",
				Mappings: []Mapping{
					{
						InPath: "key2",
						ToPath: "newKey",
					},
				},
				ApplyFile: "before",
			},
			{
				Path: "test.yaml",
				Mappings: []Mapping{
					{
						InPath: "key",
						ToPath: "newKey",
					},
				},
				ApplyFile: "after",
			},
		},
	}

	for i := 0; i < 1000; i++ {

		memoryMap, _, err := ReadMemoryMap(&ConfigurationMap)

		// fmt.Println(memoryMap["newKey"])

		assert.Nil(t, err)

		assert.Equal(t, "valuejson", memoryMap["key"])
		assert.Equal(t, "valueyaml3", memoryMap["key2"].(map[string]interface{})["key3"])
		assert.Equal(t, "value5", memoryMap["key2"].(map[string]interface{})["key5"])
		assert.Equal(t, "value4", memoryMap["key2"].(map[string]interface{})["key4"])
		assert.Equal(t, "valueyaml", memoryMap["newKey"])

		assert.Equal(t, expectedMemoryMap, memoryMap)

		i++
	}

}
