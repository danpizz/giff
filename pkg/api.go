package pkg

import (
	"errors"
	"io/ioutil"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	cf "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cfTypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/dchest/uniuri"
)

type API interface {
	ReadTemplateFile(templateFileName string) (body string, err error)
}
type APIClient struct {
}

func (APIClient) ReadTemplateFile(templateFileName string) (body string, err error) {
	templateFileBytes, err := ioutil.ReadFile(templateFileName)
	if err != nil {
		return
	}
	body = string(templateFileBytes)
	return
}

func GetStackParameters(api CFAPI, stackName *string) ([]cfTypes.Parameter, error) {
	describeStacksOutput, err := api.DescribeStacks(&cf.DescribeStacksInput{
		StackName: stackName,
	})
	if err != nil {
		return nil, err
	}
	if describeStacksOutput.Stacks == nil || len(describeStacksOutput.Stacks) != 1 {
		return nil, errors.New("cannot read stack parameters")
	}
	return describeStacksOutput.Stacks[0].Parameters, nil
}

const changesetBaseName = "giff"

func CreateChangeSet(api CFAPI, stackName *string, templateBody *string, parameterList []cfTypes.Parameter, tagList []cfTypes.Tag) (changeSetId string, err error) {

	capabilities := []cfTypes.Capability{cfTypes.CapabilityCapabilityNamedIam}

	createChangeSetInput := cf.CreateChangeSetInput{
		StackName:     stackName,
		ChangeSetName: aws.String(changesetBaseName + "-" + uniuri.New()),
		ChangeSetType: "UPDATE",
		TemplateBody:  templateBody,
		Capabilities:  capabilities,
	}
	if parameterList != nil {
		createChangeSetInput.Parameters = parameterList
	}
	if tagList != nil {
		createChangeSetInput.Tags = tagList
	}

	changesetOutput, err := api.CreateChangeSet(&createChangeSetInput)
	if err != nil {
		return "", err
	}
	return *changesetOutput.Id, nil
}

// no waiters in the aws-sdk-go-v2 for cloudformation yet
// https://github.com/aws/aws-sdk-go-v2/issues/1111
func WaitForChangeSet(api CFAPI, changeSetArn string, print func(string, ...interface{})) (out *cf.DescribeChangeSetOutput, err error) {
	var try int = 1
	const maxRetries = 20
	print("Reading changeset...")
	for {
		if try > maxRetries {
			print("\n")
			return &cf.DescribeChangeSetOutput{}, errors.New("max retries while waiting for changset")
		}
		out, err = api.DescribeChangeSet(&cf.DescribeChangeSetInput{
			ChangeSetName: aws.String(changeSetArn),
		})
		if err != nil {
			// delete changeset?
			return
		}
		if out.Status != cfTypes.ChangeSetStatusCreateComplete && out.Status != cfTypes.ChangeSetStatusFailed {
			time.Sleep(2 * time.Second)
			try += 1
			print(".")
			continue
		}
		print("ok\n")
		return
	}
}

func DeleteChangeset(api CFAPI, changeSetArn *string) error {
	_, err := api.DeleteChangeSet(&cf.DeleteChangeSetInput{
		ChangeSetName: changeSetArn,
	})
	return err
}
