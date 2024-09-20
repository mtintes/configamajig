package actions

import (
	"errors"
	"fmt"

	"github.com/nqd/flat"
)

func ReadKeyCmd(configurationMap *ConfigurationMap, key string, outputFilePath string, traceOutFilePath string) (string, error) {

	memoryMap, traces, err := ReadMemoryMap(configurationMap)

	if err != nil {
		return "", err
	}

	if traceOutFilePath != "" {
		traceOutputType := findFileType(traceOutFilePath)
		flatMemoryMap, err := flat.Flatten(memoryMap, nil)

		if err != nil {
			fmt.Println(err)
			return "", err
		}
		table, err := traceToTableForSingleKey(traces, configurationMap.Configs, traceOutputType, key, flatMemoryMap[key])

		if err != nil {
			fmt.Println(err)
			return "", err
		}

		WriteFile(traceOutFilePath, table)
	}

	flatMemoryMap, err := flat.Flatten(memoryMap, nil)

	if err != nil {
		return "", err
	}

	if outputFilePath == "" {
		return flatMemoryMap[key].(string), nil
	} else {
		outputFileType := findFileType(outputFilePath)

		switch outputFileType {
		case "json":
			jsonString := "{\"" + key + "\":" + flatMemoryMap[key].(string) + "}"
			err = WriteFile(outputFilePath, jsonString)
		case "yaml":
			yamlString := key + ": " + flatMemoryMap[key].(string)
			err = WriteFile(outputFilePath, yamlString)
		default:
			return "", errors.New("output file type not supported. (json/yaml)")
		}
	}
	return flatMemoryMap[key].(string), err
}
