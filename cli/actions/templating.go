package actions

import (
	"strings"
	"text/template"
)

func RunTemplate(inputFile []byte, memoryMap map[string]interface{}) (string, error) {

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
