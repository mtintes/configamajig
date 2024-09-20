/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/mtintes/configamajig/actions"
	"github.com/spf13/cobra"
)

// replaceCmd represents the replace command
var replaceCmd = &cobra.Command{
	Use:   "replace",
	Short: "Replaces the variables in the input file using a config file",
	Long:  `This is for replacing the variables in the input file using a config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		configMapPath := cmd.Flag("config").Value.String()
		traceFileOutput := cmd.Flag("memoryTraceOut").Value.String()
		inputPath := cmd.Flag("input").Value.String()
		outputPath := cmd.Flag("output").Value.String()

		configurationMap, err := actions.ReadConfigurationMap(configMapPath)

		if err != nil {
			fmt.Println(err)
			return
		}
		actions.ReplaceCmd(configurationMap, inputPath, outputPath, traceFileOutput)
	},
}

func init() {
	rootCmd.AddCommand(replaceCmd)

	replaceCmd.Flags().StringP("config", "c", "", "config file defines the mapping of the variables (order, depth, etc.)")
	replaceCmd.Flags().StringP("input", "i", "", "input file to be replaced")
	replaceCmd.Flags().StringP("output", "o", "", "output file to be written")
	replaceCmd.Flags().StringP("memoryTraceOut", "t", "", "changes made during memory map setup")

	replaceCmd.MarkFlagRequired("config")
	replaceCmd.MarkFlagRequired("input")
	replaceCmd.MarkFlagRequired("output")

}
