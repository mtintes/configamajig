package actions

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nqd/flat"
)

func ReadMemoryMap(configurationMap *ConfigurationMap) (map[string]interface{}, *[]Trace, error) {

	//traceOn := true
	var traces []Trace

	fmt.Println("---------------")

	var masterMemoryMap = make(map[string]interface{})
	storedFiles := make([]StoredMemoryMap, 0)
	mapLengthCount := 0
	for _, config := range configurationMap.Configs {
		mapLengthCount += len(config.Mappings) + 1
	}

	// fmt.Println(configurationMap.Configs[0])

	for _, config := range configurationMap.Configs {

		filePath := config.Path
		fileType := findFileType(filePath)

		var err error
		var file interface{}

		file = returnStoredFile(storedFiles, filePath)

		if file == nil {
			if fileType == "json" {
				file, err = SlurpJson(filePath)
			} else if fileType == "yaml" {
				file, err = SlurpYaml(filePath)
			} else {
				return nil, &traces, errors.New("file type not supported")
			}
		}

		if err != nil {
			return nil, &traces, err
		}

		//fmt.Println(file)

		flatFile, err := flat.Flatten(file.(map[string]interface{}), nil)

		if err != nil {
			return nil, &traces, err
		}
		fmt.Println("filePath:", filePath)
		if config.ApplyFile == "before" {
			masterMemoryMap = applyFlatFile(flatFile, masterMemoryMap)
			masterMemoryMap, flatFile, traces = applyMappings(flatFile, config.Mappings, masterMemoryMap, traces)
		} else if config.ApplyFile == "later" {
			for key, value := range flatFile {
				masterMemoryMap[key] = value
			}
		} else {
			masterMemoryMap, flatFile, traces = applyMappings(flatFile, config.Mappings, masterMemoryMap, traces)
			masterMemoryMap = applyFlatFile(flatFile, masterMemoryMap)
		}
		storedFiles = append(storedFiles, StoredMemoryMap{File: flatFile, FileName: filePath})
	}

	fmt.Println("Before:", masterMemoryMap["newKey"])
	unflatten, err := flat.Unflatten(masterMemoryMap, nil)

	if err != nil {
		return nil, &traces, err
	}

	fmt.Println("After", unflatten["newKey"])

	fmt.Println("---------------")
	return unflatten, &traces, nil
}

func findFileType(filePath string) string {
	return filePath[len(filePath)-4:]
}

func SafePropertyCheck(obj interface{}, key string) bool {
	if obj == nil {
		return false
	}

	if _, ok := obj.(map[string]interface{})[key]; ok {
		return true
	}

	return false
}

func returnStoredFile(storedFiles []StoredMemoryMap, fileName string) interface{} {
	for _, file := range storedFiles {
		if file.FileName == fileName {
			return file.File
		}
	}
	return nil
}

func applyMappings(flatFile map[string]interface{}, mappings []Mapping, masterMemoryMap map[string]interface{}, traces []Trace) (map[string]interface{}, map[string]interface{}, []Trace) {
	newFlatFile := flatFile

	for _, mapping := range mappings {
		for key, value := range flatFile {
			fmt.Println("Value:", value)

			//if key has partial match
			if strings.HasPrefix(key, mapping.InPath+".") || key == mapping.InPath {
				keySuffix := key[len(mapping.InPath):]
				for memKey := range masterMemoryMap {
					if strings.HasPrefix(memKey, mapping.ToPath) && keySuffix == "" {
						traces = append(traces, Trace{key: memKey, value: masterMemoryMap[memKey], changeType: "delete"})
						delete(masterMemoryMap, memKey)
					}
				}
				traces = append(traces, Trace{key: mapping.ToPath + keySuffix, value: value, changeType: "update"})
				masterMemoryMap[mapping.ToPath+keySuffix] = value
				delete(newFlatFile, key)
			}
		}
	}
	return masterMemoryMap, newFlatFile, traces
}

func applyFlatFile(flatFile map[string]interface{}, masterMemoryMap map[string]interface{}) map[string]interface{} {
	for key, value := range flatFile {
		for memKey := range masterMemoryMap {
			if strings.HasPrefix(memKey, key+".") || memKey == key {
				delete(masterMemoryMap, memKey)
			}
		}
		masterMemoryMap[key] = value
	}
	return masterMemoryMap
}

type StoredMemoryMap struct {
	File     map[string]interface{}
	FileName string
}

type Trace struct {
	key        string
	value      interface{}
	changeType string
}

func TracesToString(traces *[]Trace) string {
	var result string
	for _, trace := range *traces {
		result += trace.changeType
		result += " "
		result += trace.key
		result += " "
		result += fmt.Sprint(trace.value)
		result += "\n"
	}
	return result
}

func Unflatten(flat map[string]interface{}) (map[string]interface{}, error) {
	response := make(map[string]interface{})

	for key, value := range flat {
		current := response
		keys := strings.Split(key, ".")
		for i, k := range keys {
			if i == len(keys)-1 {
				current[k] = value
			} else {
				if _, ok := current[k]; !ok {
					current[k] = make(map[string]interface{})
				}
				current = current[k].(map[string]interface{})
			}
		}
	}

	return response, nil
}

func createTrace(masterMemoryMap map[string]interface{}, key string, value interface{}) Trace {
	if _, ok := masterMemoryMap[key]; ok {
		return Trace{key: key, value: value, changeType: "update"}
	} else {
		return Trace{key: key, value: value, changeType: "create"}
	}
}
