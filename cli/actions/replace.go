package actions

import (
	"fmt"

	"github.com/nqd/flat"
)

func ReplaceCmd(configMap string, inputPath string, outputPath string, traceFileOutput string) {
	configurationMap, err := ReadConfigurationMap(configMap)

	if err != nil {
		fmt.Println(err)
		return
	}

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
	inputFile, err := SlurpGenericFile(inputPath)

	if err != nil {
		fmt.Println(err)
		return
	}

	file, err := RunTemplate(inputFile, memoryMap)

	if err != nil {
		fmt.Println(err)
		return
	}

	err = WriteFile(outputPath, file)

	if err != nil {
		fmt.Println(err)
		return
	}

}
