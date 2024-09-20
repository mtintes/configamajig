package actions

import "fmt"

func GenerateConfigCmd() {
	version1File := `{
	"version": "1.0.0",
	"configs": [
		{
			"path": "", # Path to the file
			"mappings": [
				"InPath":"",	# Path to the key in the input file that you want to redirect
				"OutPath":"",	# Path that InPath will be redirected to
			],
			"applyFile": "after" # before/after/later - when to apply the file in relation to the mappings.
		}
	]
}`

	fmt.Println(version1File)

}
