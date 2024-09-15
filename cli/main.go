/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import "github.com/mtintes/configamajig/cli/cmd"

var version string
var commit string
var date string

func main() {

	cmd.SetVersionInfo(version, commit, date)
	cmd.Execute()
}
