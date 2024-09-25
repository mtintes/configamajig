package actions

import (
	"fmt"
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
	masterMemoryMap, flatFile, _ = applyMappings(flatFile, mappings, masterMemoryMap, traces, "test.json")

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

func TestOmegaSuperDuperBonkersMemoryMap(t *testing.T) {

	AppFs = afero.NewMemMapFs()

	jsonGlobalData := `{
		"GlobalAppConfigs": {	
			"connectionString": "value.server.com",
			"buildVersion": "development",
			"buildNumber": 1,
			"replacementInfo": "yes"
			}
		}`
	_ = afero.WriteFile(AppFs, "test.json", []byte(jsonGlobalData), 0644)

	jsonRegionalData := `{
		"RegionalAppConfigs": {
			"connectionString": "regionalValue.server.com",
			"buildVersion": "integration",
			"buildNumber": 2
		},
		"Apps": {
			"app1": {
				"specificConfig": {
					"key": "valueregional",
					"key3": "value3"
				}
			}		
		}
	}`
	_ = afero.WriteFile(AppFs, "test2.json", []byte(jsonRegionalData), 0644)

	jsonLocalData := `{
		"LocalAppConfigs": {
			"connectionString": "{{ .AppConfigs.replacementInfo }}.server.com",
			"buildVersion": "production",
			"buildNumber": 3
			},
		"Apps": {
			"app1": {
				"connectionString": "app1Value.server.com",
				"buildVersion": "productionApp1",
				"buildNumber": 4,
				"specificConfig": {
					"key": "value",
					"key2": "value2"
				}
			},
			"app2": {
				"buildVersion": "productionApp2",
				"buildNumber": 5
			}
		}
	}`
	_ = afero.WriteFile(AppFs, "test3.json", []byte(jsonLocalData), 0644)

	jsonDeploymentData := `{
		"DeploymentConfigs": {
			"cluster": "clusterValue",
			"region": "regionValue",
			"client": "clientValue"
		}
	}`
	_ = afero.WriteFile(AppFs, "test4.json", []byte(jsonDeploymentData), 0644)

	config := ConfigurationMap{
		Configs: []Config{
			{
				Path: "test.json",
				Mappings: []Mapping{
					{
						InPath: "GlobalAppConfigs",
						ToPath: "AppConfigs",
					},
				},
				ApplyFile: "after",
			},
			{
				Path: "test2.json",
				Mappings: []Mapping{
					{
						InPath: "RegionalAppConfigs",
						ToPath: "AppConfigs",
					},
					{
						InPath: "Apps",
						ToPath: "AppConfigs",
					},
				},
				ApplyFile: "after",
			},
			{
				Path: "test3.json",
				Mappings: []Mapping{
					{
						InPath: "LocalAppConfigs",
						ToPath: "AppConfigs",
					},
					{
						InPath: "Apps",
						ToPath: "AppConfigs",
					},
				},
				ApplyFile: "after",
			},
			{
				Path:      "test4.json",
				Mappings:  []Mapping{},
				ApplyFile: "after",
			},
		},
	}

	expectedMemoryMap := map[string]interface{}{
		"AppConfigs": map[string]interface{}{
			"connectionString": "yes.server.com",
			"buildVersion":     "production",
			"buildNumber":      3,
			"replacementInfo":  "yes",
			"app1": map[string]interface{}{
				"connectionString": "app1Value.server.com",
				"buildVersion":     "productionApp1",
				"buildNumber":      4,
				"specificConfig": map[string]interface{}{
					"key":  "value",
					"key2": "value2",
					"key3": "value3",
				},
			},
			"app2": map[string]interface{}{
				"buildVersion": "productionApp2",
				"buildNumber":  5,
			},
		},
		"DeploymentConfigs": map[string]interface{}{
			"cluster": "clusterValue",
			"region":  "regionValue",
			"client":  "clientValue",
		},
	}

	memoryMap, _, err := ReadMemoryMap(&config)

	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprint(expectedMemoryMap), fmt.Sprint(memoryMap)) // ints get turned into floats

}

func TestReadMemoryMap_DeepLinks(t *testing.T) {

	AppFs = afero.NewMemMapFs()

	jsonData1 := `{
		"key": {
			"key2": {
				"key3": {
					"key4": {
						"key5": "value"
					}
				}
			}
		}
	}`

	_ = afero.WriteFile(AppFs, "test.json", []byte(jsonData1), 0644)

	jsonData2 := `{
		"key": {
			"key2": {
				"key3": "newValue"
			}
		}
	}`

	_ = afero.WriteFile(AppFs, "test2.json", []byte(jsonData2), 0644)

	config := ConfigurationMap{
		Configs: []Config{
			{
				Path:      "test.json",
				Mappings:  []Mapping{},
				ApplyFile: "after",
			},
			{
				Path:      "test2.json",
				Mappings:  []Mapping{},
				ApplyFile: "after",
			},
		},
	}

	expectedMemoryMap := map[string]interface{}{
		"key": map[string]interface{}{
			"key2": map[string]interface{}{
				"key3": "newValue",
			},
		},
	}

	memoryMap, traces, err := ReadMemoryMap(&config)

	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprint(expectedMemoryMap), fmt.Sprint(memoryMap))
	assert.Equal(t, 3, len(*traces))

}

func TestReadMemoryMap_DeepLinks_Replaces(t *testing.T) {

	AppFs = afero.NewMemMapFs()

	jsonData := `{
		"key": "{{ .key2.key3 }}",
		"key2": {
			"key3": "{{ .key2.key4 }}",
			"key4": "key5"
		}
	}`

	_ = afero.WriteFile(AppFs, "test.json", []byte(jsonData), 0644)

	config := ConfigurationMap{
		Configs: []Config{
			{
				Path:      "test.json",
				Mappings:  []Mapping{},
				ApplyFile: "after",
			},
		},
	}

	expectedMemoryMap := map[string]interface{}{
		"key": "key5",
		"key2": map[string]interface{}{
			"key3": "key5",
			"key4": "key5",
		},
	}

	memoryMap, _, err := ReadMemoryMap(&config)

	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprint(expectedMemoryMap), fmt.Sprint(memoryMap))

}

func TestReadMemoryMap_10Replaces(t *testing.T) {

	AppFs = afero.NewMemMapFs()

	jsonData := `{
		"key": "{{ .key2 }}",
		"key2": "{{ .key3 }}",
		"key3": "{{ .key4 }}",
		"key4": "{{ .key5 }}",
		"key5": "{{ .key6 }}",
		"key6": "{{ .key7 }}",
		"key7": "{{ .key8 }}",
		"key8": "{{ .key9 }}",
		"key9": "{{ .key10 }}",
		"key10": "value"
	}`
	_ = afero.WriteFile(AppFs, "test.json", []byte(jsonData), 0644)

	config := ConfigurationMap{
		Configs: []Config{
			{
				Path:      "test.json",
				Mappings:  []Mapping{},
				ApplyFile: "after",
			},
		},
	}

	expectedMemoryMap := map[string]interface{}{
		"key":   "value",
		"key2":  "value",
		"key3":  "value",
		"key4":  "value",
		"key5":  "value",
		"key6":  "value",
		"key7":  "value",
		"key8":  "value",
		"key9":  "value",
		"key10": "value",
	}

	memoryMap, traces, err := ReadMemoryMap(&config)

	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprint(expectedMemoryMap), fmt.Sprint(memoryMap))
	//templating is very inefficient and needs to be run over and over
	assert.Equal(t, 35, len(*traces))
}

func TestReadMemoryMap_Traces(t *testing.T) {

	AppFs = afero.NewMemMapFs()

	jsonData := `{
		"key": "value"
	}`

	_ = afero.WriteFile(AppFs, "test.json", []byte(jsonData), 0644)

	config := ConfigurationMap{
		Configs: []Config{
			{
				Path:      "test.json",
				Mappings:  []Mapping{},
				ApplyFile: "after",
			},
		},
	}

	_, traces, err := ReadMemoryMap(&config)

	_ = traces
	assert.Nil(t, err)
	assert.Equal(t, 1, len(*traces))

}

func TestReadMemoryMap_PropToRoot(t *testing.T) {

	AppFs = afero.NewMemMapFs()

	jsonData := `{
		"key1": {
			"key2": {
				"key3": "value"}}
	}`

	_ = afero.WriteFile(AppFs, "test.json", []byte(jsonData), 0644)

	config := ConfigurationMap{
		Configs: []Config{
			{
				Path: "test.json",
				Mappings: []Mapping{
					{
						InPath: "key1.key2",
						ToPath: "",
					},
				},
				ApplyFile: "after",
			},
		},
	}

	expectedMemoryMap := map[string]interface{}{
		"key3": "value",
	}

	memoryMap, _, err := ReadMemoryMap(&config)

	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprint(expectedMemoryMap), fmt.Sprint(memoryMap))
}

func TestReadMemoryMap_Prop(t *testing.T) {

	AppFs = afero.NewMemMapFs()

	jsonData := `{
		"key1": {
			"key2": {
				"key3": "value"}}
	}`

	_ = afero.WriteFile(AppFs, "test.json", []byte(jsonData), 0644)

	config := ConfigurationMap{
		Configs: []Config{
			{
				Path: "test.json",
				Mappings: []Mapping{
					{
						InPath: "key1.key2",
						ToPath: "newKey",
					},
				},
				ApplyFile: "after",
			},
			{
				Props: map[string]interface{}{
					"key3": "value",
				},
				Mappings: []Mapping{
					{
						InPath: "key3",
						ToPath: "newKey2",
					},
				},
			},
		},
	}

	expectedMemoryMap := map[string]interface{}{
		"newKey": map[string]interface{}{
			"key3": "value",
		},
		"newKey2": "value",
	}

	memoryMap, _, err := ReadMemoryMap(&config)

	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprint(expectedMemoryMap), fmt.Sprint(memoryMap))

}

func TestMemoryMapWithArrays(t *testing.T) {

	AppFs = afero.NewMemMapFs()

	jsonData := `{
		"key": [
			{
				"key1": "value"
			},
			{
				"key2": "value2"
			}
		],
		"keyTemplate": "{{ .key | toRawJson}}"
	}`

	_ = afero.WriteFile(AppFs, "test.json", []byte(jsonData), 0644)

	config := ConfigurationMap{
		Configs: []Config{
			{
				Path:      "test.json",
				Mappings:  []Mapping{},
				ApplyFile: "after",
			},
		},
	}

	expectedMemoryMap := map[string]interface{}{
		"key": []interface{}{
			map[string]interface{}{
				"key1": "value",
			},
			map[string]interface{}{
				"key2": "value2",
			},
		},
		"keyTemplate": "[{\"key1\":\"value\"},{\"key2\":\"value2\"}]",
	}

	memoryMap, _, err := ReadMemoryMap(&config)

	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprint(expectedMemoryMap), fmt.Sprint(memoryMap))
}

func TestMemoryMapWithJSON(t *testing.T) {

	AppFs = afero.NewMemMapFs()

	jsonData := `{
		"key": {
			"key1": "value",
			"key2": "value2",
			"key3": {
				"key4": "value4"
			}
		},
		"keyTemplate": "{{ .key | toJson}}"
	}`

	_ = afero.WriteFile(AppFs, "test.json", []byte(jsonData), 0644)

	config := ConfigurationMap{
		Configs: []Config{
			{
				Path:      "test.json",
				Mappings:  []Mapping{},
				ApplyFile: "after",
			},
		},
	}

	expectedMemoryMap := map[string]interface{}{
		"key": map[string]interface{}{
			"key1": "value",
			"key2": "value2",
			"key3": map[string]interface{}{
				"key4": "value4",
			},
		},
		"keyTemplate": "{\"key1\":\"value\",\"key2\":\"value2\",\"key3\":{\"key4\":\"value4\"}}",
	}

	memoryMap, _, err := ReadMemoryMap(&config)

	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprint(expectedMemoryMap), fmt.Sprint(memoryMap))
}

func TestMemoryMapWithIntPropsIsOutlawed(t *testing.T) {

	AppFs = afero.NewMemMapFs()

	jsonData := `{
		"key": {
			"0": 1,
			"1": 2
		},
		"keyTemplate": "{{ .key | toJson}}"
	}`

	_ = afero.WriteFile(AppFs, "test.json", []byte(jsonData), 0644)

	config := ConfigurationMap{
		Configs: []Config{
			{
				Path:      "test.json",
				Mappings:  []Mapping{},
				ApplyFile: "after",
			},
		},
	}

	expectedMemoryMap := map[string]interface{}{
		"key": map[string]interface{}{
			"0": 1,
			"1": 2,
		},
		"keyTemplate": "{\"0\":1,\"1\":2}",
	}

	memoryMap, _, err := ReadMemoryMap(&config)

	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprint(expectedMemoryMap), fmt.Sprint(memoryMap))
}
