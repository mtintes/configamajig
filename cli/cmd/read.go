/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(readCmd)

}
