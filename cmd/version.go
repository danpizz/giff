package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var ShortVersion string
var Version string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s (%s)\n", ShortVersion, Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
