package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/danpizz/giff/pkg"

	"github.com/spf13/cobra"
)

func NewDiffCmd(cfClient pkg.CFAPI, apiClient pkg.API) (diffCmd *cobra.Command) {
	diffCmd = &cobra.Command{
		Use:   "diff stackname template",
		Short: "Show the differences between a CloudFormation stack and a local template",
		Long:  "Dowload a stack template and run run the diff command over it and a local template",
		Run: func(cmd *cobra.Command, args []string) {
			if err := diff(cmd, args, cfClient, apiClient); err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
		},
		Args:    cobra.ExactArgs(2),
		Example: "giff diff my-stack my-template.yaml\n",
	}
	diffCmd.Flags().StringVarP(&diffCommand, "diff-command", "d", "diff", "Command on the PATH to use to create the diff")
	return diffCmd
}

var diffCommand string

func init() {
	rootCmd.AddCommand(NewDiffCmd(nil, nil))
}

func diff(cmd *cobra.Command, args []string, cfClient pkg.CFAPI, apiClient pkg.API) (err error) {
	stackName := args[0]
	templateFileName := args[1]
	if cfClient == nil {
		cfClient, err = pkg.NewCFClient()
		if err != nil {
			return err
		}
	}

	stackTemplateOut, err := cfClient.GetTemplate(&cloudformation.GetTemplateInput{
		StackName: &stackName,
	})
	if err != nil {
		return err
	}

	templateFileData, err := ioutil.ReadFile(templateFileName)
	if err != nil {
		return err
	}

	diffOut, err := pkg.Diff("giff", diffCommand, []byte(*stackTemplateOut.TemplateBody), templateFileData)
	if err != nil {
		return err
	}
	if len(diffOut) > 0 {
		cmd.Printf("%s\n", diffOut)
	}
	return nil
}
