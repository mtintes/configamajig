package actions

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nqd/flat"
)

func ReadMemoryMap(configurationMap *ConfigurationMap) (interface{}, error) {
	var masterMemoryMap = make(map[string]interface{})
	storedFiles := make([]StoredMemoryMap, 0)
	mapLengthCount := 0
	for _, config := range configurationMap.Configs {
		mapLengthCount += len(config.Mappings) + 1
	}

	for _, config := range configurationMap.Configs {
		filePath := config.Path
		fileType := findFileType(filePath)

		var err error
		var file interface{}

		file = returnStoredFile(storedFiles, filePath)

		if file == nil {
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

		flatFile, err := flat.Flatten(file.(map[string]interface{}), nil)

		if err != nil {
			return nil, err
		}

		if config.ApplyFile == "before" {
			for key, value := range flatFile {
				masterMemoryMap[key] = value
			}
			masterMemoryMap, flatFile = applyMappings(flatFile, config.Mappings, masterMemoryMap)
		} else if config.ApplyFile == "later" {
			for key, value := range flatFile {
				masterMemoryMap[key] = value
			}
		} else {
			masterMemoryMap, flatFile = applyMappings(flatFile, config.Mappings, masterMemoryMap)
			for key, value := range flatFile {
				masterMemoryMap[key] = value
			}
		}
		storedFiles = append(storedFiles, StoredMemoryMap{File: flatFile, FileName: filePath})
	}

	return flat.Unflatten((masterMemoryMap), nil)
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

func applyMappings(flatFile map[string]interface{}, mappings []Mapping, masterMemoryMap map[string]interface{}) (map[string]interface{}, map[string]interface{}) {
	for _, mapping := range mappings {
		for key, value := range flatFile {
			if strings.HasPrefix(key, mapping.InPath) {
				keySuffix := key[len(mapping.InPath):]
				masterMemoryMap[mapping.ToPath+keySuffix] = value
				delete(flatFile, key)
			}
		}
	}
	return masterMemoryMap, flatFile
}

type StoredMemoryMap struct {
	File     map[string]interface{}
	FileName string
}
