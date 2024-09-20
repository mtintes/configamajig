/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/mtintes/configamajig/actions"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "generates a starter config file",
	Long:  `This command will generate a starter config file for you to use as a template.`,
	Run: func(cmd *cobra.Command, args []string) {
		actions.GenerateConfigCmd()
	},
}

func init() {
	generateCmd.AddCommand(configCmd)

}
