/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"math/rand/v2"

	"github.com/mtintes/configamajig/actions"
	"github.com/spf13/cobra"
)

// replaceCmd represents the replace command
var replaceCmd = &cobra.Command{
	Use:   "replace",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("replace called")
		configMap := cmd.Flag("config").Value.String()
		configurationMap, err := actions.ReadConfigurationMap(configMap)

		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(configurationMap)

		memoryMap, traces, err := actions.ReadMemoryMap(configurationMap)

		if err != nil {
			fmt.Println(err)
			return
		}

		random := rand.IntN(1000)
		actions.WriteFile(fmt.Sprintf("traces/trace_%d.json", random), actions.TracesToString(traces))

		// templatedMemoryMap, err := actions.TemplateMemoryMap(memoryMap.(map[string]interface{}))

		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }

		inputPath := cmd.Flag("input").Value.String()
		inputFile, err := actions.SlurpGenericFile(inputPath)

		if err != nil {
			fmt.Println(err)
			return
		}

		file, err := actions.RunTemplate(inputFile, memoryMap)

		if err != nil {
			fmt.Println(err)
			return
		}

		outputPath := cmd.Flag("output").Value.String()
		err = actions.WriteFile(outputPath, file)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(replaceCmd)

	replaceCmd.Flags().StringP("config", "c", "", "config file defines the mapping of the variables (order, depth, etc.)")
	replaceCmd.Flags().StringP("input", "i", "", "input file to be replaced")
	replaceCmd.Flags().StringP("output", "o", "", "output file to be written")

}
