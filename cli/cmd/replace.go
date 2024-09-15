/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
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
		configMap := cmd.Flag("config").Value.String()

		traceFileOutput := cmd.Flag("memoryTraceOut").Value.String()

		inputPath := cmd.Flag("input").Value.String()

		outputPath := cmd.Flag("output").Value.String()
		actions.ReplaceCmd(configMap, inputPath, outputPath, traceFileOutput)
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
