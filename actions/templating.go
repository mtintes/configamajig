package actions

import (
	"strings"
	"text/template"
)

func RunTemplate(inputPath string, memoryMap map[string]interface{}) (string, error) {

	inputFile, err := slurpGenericFile(inputPath)

	if err != nil {
		return "", err
	}

	tmpl, err := template.New("replace").Parse(string(inputFile))

	if err != nil {
		return "", err
	}
	b := new(strings.Builder)
	err = tmpl.Execute(b, memoryMap)
	if err != nil {
		return "", err
	}

	return b.String(), nil

}

// func TemplateMemoryMap(memoryMap map[string]interface{}) (map[string]interface{}, error) {

// 	counter := 0
// 	for counter < 10 {
// 		for key, value := range memoryMap {
// 			if value != nil {
// 				if str, ok := value.(string); ok {
// 					tmpl, err := template.New("replace").Parse(str)

// 					if err != nil {
// 						return nil, err
// 					}
// 					b := new(strings.Builder)
// 					err = tmpl.Execute(b, memoryMap)
// 					if err != nil {
// 						return nil, err
// 					}

// 					memoryMap[key] = b.String()
// 				}
// 			}
// 		}
// 		counter++
// 	}

// 	return memoryMap, nil

// }
