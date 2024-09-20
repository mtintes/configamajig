package actions

import (
	"fmt"

	"github.com/nqd/flat"
)

func ReplaceCmd(configurationMap *ConfigurationMap, inputPath string, outputPath string, traceFileOutput string) {

	memoryMap, traces, err := ReadMemoryMap(configurationMap)

	if err != nil {
		fmt.Println(err)
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

	if err == nil {
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

}
