package actions

import (
	"errors"
	"fmt"
)

func ReadMemoryMap(configurationMap *ConfigurationMap) (interface{}, error) {
	masterMemoryMap := make(map[string]interface{})
	storedFiles := make([]StoredMemoryMap, 0)
	mapLengthCount := 0
	for _, config := range configurationMap.Configs {
		mapLengthCount += len(config.Mappings) + 1
	}

	index := 0
	for _, config := range configurationMap.Configs {
		filePath := config.Path
		fileType := findFileType(filePath)

		var err error
		var file interface{}

		file = returnStoredFile(storedFiles, filePath)

		if file != nil {
			if fileType == "json" {
				file, err = slurpJson(filePath)
			} else if fileType == "yaml" {
				file, err = slurpYaml(filePath)
			} else {
				return nil, errors.New("file type not supported")
			}
		}

		if err != nil {
			return nil, err
		}

		fmt.Println(file)

		if config.ApplyFile == "before" {
			for key, value := range file.(map[string]interface{}) {
				masterMemoryMap[key] = value
			}
			file, masterMemoryMap = applyMappings(file, config.Mappings, masterMemoryMap)
		} else if config.ApplyFile == "later" {
			file, masterMemoryMap = applyMappings(file, config.Mappings, masterMemoryMap)
			storedFiles = append(storedFiles, StoredMemoryMap{file, filePath})
			//store the file and put in the map if asked later
		} else {
			file, masterMemoryMap = applyMappings(file, config.Mappings, masterMemoryMap)
			for key, value := range file.(map[string]interface{}) {
				masterMemoryMap[key] = value
			}
		}

		fmt.Println(file)
		index++

	}

	return masterMemoryMap, nil
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

func applyMappings(file interface{}, mappings []Mapping, masterMemoryMap map[string]interface{}) (interface{}, map[string]interface{}) {
	newFile := file

	for _, mapping := range mappings {
		if SafePropertyCheck(file, mapping.InPath) {
			masterMemoryMap[mapping.ToPath] = file.(map[string]interface{})[mapping.InPath]
			delete(newFile.(map[string]interface{}), mapping.InPath)
		}
	}

	return newFile, masterMemoryMap
}

func returnStoredFile(storedFiles []StoredMemoryMap, fileName string) interface{} {
	for _, file := range storedFiles {
		if file.FileName == fileName {
			return file.File
		}
	}
	return nil
}

type StoredMemoryMap struct {
	File     interface{}
	FileName string
}
