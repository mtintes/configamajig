/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// autocompleteCmd represents the autocomplete command
var autocompleteCmd = &cobra.Command{
	Use:   "autocomplete",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		shell, _ := cmd.Flags().GetString("shell")

		switch shell {
		case "bash":
			rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			rootCmd.GenZshCompletion(os.Stdout)
		default:
			fmt.Println("Unsupported shell type")
		}

	},
}

func init() {
	rootCmd.AddCommand(autocompleteCmd)

	autocompleteCmd.Flags().StringP("shell", "s", "", "shell type (bash, zsh)")

	autocompleteCmd.MarkFlagRequired("shell")

}
