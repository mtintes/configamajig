package actions

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nqd/flat"
)

func ReadMemoryMap(configurationMap *ConfigurationMap) (map[string]interface{}, *[]Trace, error) {

	var traces []Trace

	var masterMemoryMap = make(map[string]interface{})
	storedFiles := make([]StoredMemoryMap, 0)

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

		flatFile, err := flat.Flatten(file.(map[string]interface{}), nil)

		if err != nil {
			return nil, &traces, err
		}

		if config.ApplyFile == "before" {
			masterMemoryMap, traces = applyFlatFile(flatFile, masterMemoryMap, traces)
			masterMemoryMap, flatFile, traces = applyMappings(flatFile, config.Mappings, masterMemoryMap, traces)
		} else if config.ApplyFile == "later" {
			masterMemoryMap, flatFile, traces = applyMappings(flatFile, config.Mappings, masterMemoryMap, traces)
		} else {
			masterMemoryMap, flatFile, traces = applyMappings(flatFile, config.Mappings, masterMemoryMap, traces)
			masterMemoryMap, traces = applyFlatFile(flatFile, masterMemoryMap, traces)
		}
		storedFiles = append(storedFiles, StoredMemoryMap{File: flatFile, FileName: filePath})
	}

	masterMemoryMap, traces, err := templateMemoryMap(masterMemoryMap, traces)

	if err != nil {
		return nil, &traces, err
	}

	unflatten, err := flat.Unflatten(masterMemoryMap, nil)

	if err != nil {
		return nil, &traces, err
	}

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

			if strings.HasPrefix(key, mapping.InPath+".") || key == mapping.InPath {
				keySuffix := key[len(mapping.InPath):]
				for memKey := range masterMemoryMap {
					if strings.HasPrefix(memKey, mapping.ToPath) && keySuffix == "" {
						traces = append(traces, Trace{key: memKey, oldValue: masterMemoryMap[memKey], changeType: "delete"})
						delete(masterMemoryMap, memKey)
					}
				}
				traces = append(traces, createTrace(masterMemoryMap, mapping.ToPath+keySuffix, value))

				masterMemoryMap[mapping.ToPath+keySuffix] = value
				delete(newFlatFile, key)
			}
		}
	}
	return masterMemoryMap, newFlatFile, traces
}

func applyFlatFile(flatFile map[string]interface{}, masterMemoryMap map[string]interface{}, traces []Trace) (map[string]interface{}, []Trace) {
	for key, value := range flatFile {
		for memKey := range masterMemoryMap {
			if strings.HasPrefix(memKey, key+".") || memKey == key {
				traces = append(traces, Trace{key: memKey, value: masterMemoryMap[memKey], changeType: "delete"})
				delete(masterMemoryMap, memKey)
			}
		}
		traces = append(traces, createTrace(masterMemoryMap, key, value))

		masterMemoryMap[key] = value
	}
	return masterMemoryMap, traces
}

type StoredMemoryMap struct {
	File     map[string]interface{}
	FileName string
}

type Trace struct {
	key        string
	value      interface{}
	oldValue   interface{}
	changeType string
}

func TracesToString(traces *[]Trace) string {
	var result string
	for _, trace := range *traces {
		result += trace.changeType
		result += " "
		result += trace.key
		result += " "
		result += fmt.Sprint("Value: ", trace.value)
		if trace.oldValue != nil {
			result += " "
			result += fmt.Sprint("Old Value: ", trace.oldValue)
		}
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
		return Trace{key: key, value: value, oldValue: masterMemoryMap[key], changeType: "update"}
	} else {
		return Trace{key: key, value: value, changeType: "create"}
	}
}

func templateMemoryMap(flatFile map[string]interface{}, traces []Trace) (map[string]interface{}, []Trace, error) {

	for i := 0; i < 10; i++ {
		unflattenedFile, err := flat.Unflatten(flatFile, nil)

		if err != nil {
			return nil, traces, err
		}

		for key, value := range flatFile {
			if value, ok := value.(string); ok {
				oldValue := value
				flatFile[key], err = RunTemplate([]byte(value), unflattenedFile)
				if oldValue != flatFile[key] {
					traces = append(traces, Trace{key: key, value: flatFile[key], oldValue: oldValue, changeType: "template"})
				}
				if err != nil {
					return nil, traces, err
				}
			}
		}
	}
	return flatFile, traces, nil
}
