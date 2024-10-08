package actions

import (
	"encoding/json"
	"fmt"

	"github.com/nqd/flat"
	"gopkg.in/yaml.v3"
)

func RemapCmd(configurationMap *ConfigurationMap, outputPath string, traceFileOutput string) {

	memoryMap, traces, err := ReadMemoryMap(configurationMap)

	if err != nil {
		fmt.Println(err)
		return
	}

	if traceFileOutput != "" {
		traceOutputType := findFileType(traceFileOutput)
		flatMemoryMap, err := flat.Flatten(memoryMap, nil)

		if err != nil {
			fmt.Println(err)
			return
		}
		table, err := traceToTable(traces, configurationMap.Configs, traceOutputType, flatMemoryMap)

		if err != nil {
			fmt.Println(err)
			return
		}

		WriteFile(traceFileOutput, table)
	}

	outputType := findFileType(outputPath)

	var mappedFileString []byte

	if outputType == "json" {
		mappedFileString, err = json.Marshal(memoryMap)

		if err != nil {
			fmt.Println(err)
			return
		}

	} else if outputType == "yaml" {
		mappedFileString, err = yaml.Marshal(memoryMap)

		if err != nil {
			fmt.Println(err)
			return
		}
	}

	err = WriteFile(outputPath, string(mappedFileString))

	if err != nil {
		fmt.Println(err)
		return
	}

}
