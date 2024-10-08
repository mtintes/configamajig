/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/mtintes/configamajig/actions"
	"github.com/spf13/cobra"
)

// remapCmd represents the remap command
var remapCmd = &cobra.Command{
	Use:   "remap",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		configMapPath := cmd.Flag("config").Value.String()
		traceFileOutput := cmd.Flag("memoryTraceOut").Value.String()
		outputPath := cmd.Flag("output").Value.String()

		configurationMap, err := actions.ReadConfigurationMap(configMapPath)

		if err != nil {
			fmt.Println(err)
			return
		}

		actions.RemapCmd(configurationMap, outputPath, traceFileOutput)

	},
}

func init() {
	rootCmd.AddCommand(remapCmd)

	remapCmd.Flags().StringP("config", "c", "", "config file defines the mapping of the variables (order, depth, etc.)")
	remapCmd.Flags().StringP("output", "o", "", "output file to be written")
	remapCmd.Flags().StringP("memoryTraceOut", "t", "", "changes made during memory map setup")

	remapCmd.MarkFlagRequired("config")
	remapCmd.MarkFlagRequired("output")
}
