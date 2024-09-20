package actions

import (
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

func RunTemplate(inputFile []byte, memoryMap map[string]interface{}) (string, error) {

	tmpl := template.Must(template.New("replace").Funcs(sprig.FuncMap()).Parse(string(inputFile)))

	b := new(strings.Builder)
	err := tmpl.Execute(b, memoryMap)
	if err != nil {
		return "", err
	}

	return b.String(), nil

}
