package cmd

import (
	"encoding/json"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var Out io.Writer
var Err io.Writer

var rootCmd = &cobra.Command{
	Use:   "giff",
	Short: "Giff is a CloudFormation differ",
}

func Execute() {
	if Out != nil {
		rootCmd.SetOut(Out)
	} else {
		rootCmd.SetOut(os.Stdout)
	}
	if Err != nil {
		rootCmd.SetErr(Err)
	} else {
		rootCmd.SetErr(os.Stderr)
	}
	if err := rootCmd.Execute(); err != nil {
		rootCmd.PrintErrln(err)
		os.Exit(1)
	}
}

var verbose bool

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}

func PrintfV(format string, a ...interface{}) {
	if verbose {
		rootCmd.Printf(format, a...)
	}
}

func PrettyJson(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
