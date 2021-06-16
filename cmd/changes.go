package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	cfTypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"

	"github.com/danpizz/giff/pkg"
	"github.com/spf13/cobra"
)

func NewChangesCmd(cfClient pkg.CFAPI, apiClient pkg.API) *cobra.Command {
	changesCmd := &cobra.Command{
		Use:   "changes {stackname template-file [-p par1=val1 ... | -a par1=val1 ...] [--no-delete-changeset] | stack_arn} [--dump] [-v]",
		Short: "Show a human redable list of Cloudformation changes",
		Long:  "Create a temporary changeset and display an easy to read summary of the changes created by deploying a local template and some (optional) parameters",
		Run: func(cmd *cobra.Command, args []string) {
			if err := changes(cmd, cfClient, apiClient); err != nil {
				cmd.PrintErr(err)
				os.Exit(1)
			}
		},
		Args: func(cmd *cobra.Command, args []string) error {
			switch len(args) {
			case 1:
				if Parameters != "" || ParametersOverride != "" || NoDeleteChangeset {
					return fmt.Errorf("unaccepted flag")
				}
				ChangesetArn = args[0]
				return nil
			case 2:
				StackName = args[0]
				TemplateFileName = args[1]
				return nil
			}
			return fmt.Errorf("accepts 1 or 2 args, received %d", len(args))
		},
		Example: "giff change my-stack my-template.yaml -a Size=m4.tiny -v --no-delete-changeset\n" +
			"giff change arn:aws:cloudformation:us-east-1:123456789012:changeSet/SampleChangeSet-direct/1a2345b6-0000-00a0-a123-00abc0abc000 --dump",
	}
	changesCmd.Flags().StringVarP(&Parameters, "all-parameters", "a", "", "All the template parameters: \"par1=value1 para2=value2 ...\"")
	changesCmd.Flags().StringVarP(&ParametersOverride, "parameters-overrides", "p", "", "The input parameters for your stack template. If you don't specify a parameter, the stack's existing value is used. \"par1=value1 para2=value2 ...\"")
	changesCmd.Flags().BoolVar(&NoDeleteChangeset, "no-delete-changeset", false, "Don't remove the changeset, print its ARN")
	changesCmd.Flags().BoolVarP(&Dump, "dump", "d", false, "Print the raw changeset")
	return changesCmd
}

var TemplateFileName string
var StackName string

var Parameters string
var ParametersOverride string
var NoDeleteChangeset bool = false
var ChangesetArn string
var Dump bool = false

func init() {
	rootCmd.AddCommand(NewChangesCmd(nil, nil))
}

func changes(cmd *cobra.Command, cfClient pkg.CFAPI, apiClient pkg.API) (err error) {
	if cfClient == nil {
		cfClient, err = pkg.NewCFClient()
		if err != nil {
			return err
		}
	}
	if apiClient == nil {
		apiClient = pkg.APIClient{}
	}

	var parameters []cfTypes.Parameter
	var changesetArn string
	if ChangesetArn != "" {
		changesetArn = ChangesetArn
	} else {
		if ParametersOverride != "" || (Parameters == "" && ParametersOverride == "") {
			stackParameters, err := pkg.GetStackParameters(cfClient, aws.String(StackName))
			if err != nil {
				return err
			}
			p := pkg.ParameterListFromString(ParametersOverride)
			parameters, err = pkg.OverrideParameters(stackParameters, p)
			if err != nil {
				return err
			}
		} else {
			parameters = pkg.ParameterListFromString(Parameters)
		}

		PrintfV("Creating changeset...")
		templateBody, err := apiClient.ReadTemplateFile(TemplateFileName)
		if err != nil {
			return err
		}

		changesetArn, err = pkg.CreateChangeSet(cfClient, &StackName, &templateBody, parameters)
		if err != nil {
			PrintfV("\n")
			return err
		}
		PrintfV("ok\n")
		if NoDeleteChangeset {
			cmd.Printf("changeset arn: %s\n", changesetArn)
		}
	}

	describeChangesetOutput, err := pkg.WaitForChangeSet(cfClient, changesetArn, PrintfV)
	if err != nil {
		return err
	}

	extractedChanges, err := pkg.ExtractChanges(describeChangesetOutput)
	if err != nil {
		return err
	}

	printChanges(cmd, extractedChanges)

	if Dump {
		cmd.Println(PrettyJson(describeChangesetOutput))
	}

	if ChangesetArn != "" {
		return nil
	}

	if NoDeleteChangeset {
		return nil
	}
	PrintfV("Deleting changeset...")
	err = pkg.DeleteChangeset(cfClient, &changesetArn)
	if err != nil {
		return err
	}
	PrintfV("ok\n")

	return nil
}

func printChanges(cmd *cobra.Command, changes []pkg.GiffChange) {
	if len(changes) == 0 {
		cmd.PrintErrln("No changes")
	}
	for _, c := range changes {
		switch c.Action {
		case cfTypes.ChangeActionAdd:
			cmd.Printf("+     add: %s - %s\n", *c.LogicalResourceId, *c.ResourceType)
		case cfTypes.ChangeActionRemove:
			cmd.Printf("-  remove: %s - %s\n", *c.LogicalResourceId, *c.ResourceType)
		case cfTypes.ChangeActionModify:
			cmd.Printf("*  modify: %s (%s) - %s / replacement: %v\n", *c.LogicalResourceId, *c.PhysicalResourceId, *c.ResourceType, c.Replacement)
		case cfTypes.ChangeActionDynamic:
			cmd.Printf("* dynamic: %s (%s) - %s / replacement: %v\n", *c.LogicalResourceId, *c.PhysicalResourceId, *c.ResourceType, c.Replacement)
		case cfTypes.ChangeActionImport:
			cmd.Printf("+  import: %s (%s) - %s\n", *c.LogicalResourceId, *c.PhysicalResourceId, *c.ResourceType)
		default:
			cmd.Printf("%#v [unknown change type]\n", c)
		}
	}
}
