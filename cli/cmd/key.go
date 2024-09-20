/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/mtintes/configamajig/actions"
	"github.com/spf13/cobra"
)

// keyCmd represents the key command
var keyCmd = &cobra.Command{
	Use:   "key",
	Short: "Find the value of a single key",
	Long:  `This is for finding the value of a single key`,
	Run: func(cmd *cobra.Command, args []string) {

		key := ""
		if args[0] == "" {
			fmt.Println("Please include the key to be read")
			return
		} else {
			fmt.Println("Key to be read: ", args[0])
			key = args[0]
		}

		configFilePath := cmd.Flag("config").Value.String()
		outputFilePath := cmd.Flag("output").Value.String()
		traceOutFilePath := cmd.Flag("memoryTraceOut").Value.String()

		configurationMap, err := actions.ReadConfigurationMap(configFilePath)

		if err != nil {
			fmt.Println(err)
			return
		}

		keyValue, err := actions.ReadKeyCmd(configurationMap, key, outputFilePath, traceOutFilePath)

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Key value: ", keyValue)
	},
}

func init() {
	readCmd.AddCommand(keyCmd)

	keyCmd.Flags().StringP("config", "c", "", "Config file defines the mapping of the variables (order, depth, etc.)")
	keyCmd.Flags().StringP("output", "o", "", "Optional: output file to be written. (yaml/json) Default is stdout")
	keyCmd.Flags().StringP("memoryTraceOut", "t", "", "Changes made during memory map setup")

	keyCmd.MarkFlagRequired("config")
}
