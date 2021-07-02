package pkg

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	cf "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cfTypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/stretchr/testify/assert"
)

func TestParameterListFromString(t *testing.T) {
	list := ParameterListFromString("p1=v1 p2=v2")
	assert.Exactly(t,
		[]cfTypes.Parameter{
			{
				ParameterKey:   aws.String("p1"),
				ParameterValue: aws.String("v1"),
			},
			{
				ParameterKey:   aws.String("p2"),
				ParameterValue: aws.String("v2"),
			},
		},
		list)
}
func TestOverrideParameters(t *testing.T) {
	stackParametersList :=
		[]cfTypes.Parameter{
			{
				ParameterKey:   aws.String("p1"),
				ParameterValue: aws.String("v1"),
			},
			{
				ParameterKey:   aws.String("p2"),
				ParameterValue: aws.String("v2"),
			},
		}
	inputParameters :=
		[]cfTypes.Parameter{
			{
				ParameterKey:   aws.String("p1"),
				ParameterValue: aws.String("n1"),
			},
			{
				ParameterKey:   aws.String("p3"),
				ParameterValue: aws.String("v3"),
			},
		}

	result, _ := OverrideParameters(stackParametersList, inputParameters)
	assert.Exactly(t,
		[]cfTypes.Parameter{
			{
				ParameterKey:     aws.String("p1"),
				ParameterValue:   aws.String("n1"),
				UsePreviousValue: aws.Bool(false),
			},
			{
				ParameterKey:     aws.String("p2"),
				ParameterValue:   nil,
				UsePreviousValue: aws.Bool(true),
			},
			{
				ParameterKey:     aws.String("p3"),
				ParameterValue:   aws.String("v3"),
				UsePreviousValue: aws.Bool(false),
			},
		},
		result)
}

func TestOverrideParameters_empty(t *testing.T) {
	stackParametersList :=
		[]cfTypes.Parameter{
			{
				ParameterKey:   aws.String("p1"),
				ParameterValue: aws.String("v1"),
			},
			{
				ParameterKey:   aws.String("p2"),
				ParameterValue: aws.String("p2"),
			},
		}
	inputParameters := []cfTypes.Parameter{}
	result, _ := OverrideParameters(stackParametersList, inputParameters)
	assert.Exactly(t,
		[]cfTypes.Parameter{
			{
				ParameterKey:     aws.String("p1"),
				ParameterValue:   nil,
				UsePreviousValue: aws.Bool(true),
			},
			{
				ParameterKey:     aws.String("p2"),
				ParameterValue:   nil,
				UsePreviousValue: aws.Bool(true),
			},
		},
		result)
}

func TestExtractChanges_Tag(t *testing.T) {

	content, _ := ioutil.ReadFile("testdata/020-tag.changeset.json")
	var changeSet cf.DescribeChangeSetOutput
	json.Unmarshal([]byte(content), &changeSet)
	changes, _ := ExtractChanges(&changeSet)

	LogicalResourceId := "MyEC2Instance"
	PhysicalResourceId := "i-1abc23d4"
	ResourceType := "AWS::EC2::Instance"
	assert.Exactly(
		t,
		changes,
		[]GiffChange{
			{
				Action:             cfTypes.ChangeActionModify,
				LogicalResourceId:  &LogicalResourceId,
				PhysicalResourceId: &PhysicalResourceId,
				Replacement:        cfTypes.ReplacementFalse,
				ResourceType:       &ResourceType,
			},
		},
		"error",
	)
}

func TestExtractChanges_Parameter(t *testing.T) {

	content, _ := ioutil.ReadFile("testdata/030-parameter.changeset.json")
	var changeSet cf.DescribeChangeSetOutput
	json.Unmarshal([]byte(content), &changeSet)
	changes, _ := ExtractChanges(&changeSet)

	LogicalResourceId := "MyEC2Instance"
	PhysicalResourceId := "i-1abc23d4"
	ResourceType := "AWS::EC2::Instance"
	assert.Exactly(
		t,
		changes,
		[]GiffChange{
			{
				Action:             cfTypes.ChangeActionModify,
				LogicalResourceId:  &LogicalResourceId,
				PhysicalResourceId: &PhysicalResourceId,
				Replacement:        cfTypes.ReplacementFalse,
				ResourceType:       &ResourceType,
			},
		},
		"error",
	)
}

func TestExtractChanges_Replacement(t *testing.T) {

	content, _ := ioutil.ReadFile("testdata/040-replacement.changeset.json")
	var changeSet cf.DescribeChangeSetOutput
	json.Unmarshal([]byte(content), &changeSet)
	changes, _ := ExtractChanges(&changeSet)

	LogicalResourceId := "MyEC2Instance"
	PhysicalResourceId := "i-7bef86f8"
	ResourceType := "AWS::EC2::Instance"
	assert.Exactly(
		t,
		changes,
		[]GiffChange{
			{
				Action:             cfTypes.ChangeActionModify,
				LogicalResourceId:  &LogicalResourceId,
				PhysicalResourceId: &PhysicalResourceId,
				Replacement:        cfTypes.ReplacementTrue,
				ResourceType:       &ResourceType,
			},
		},
		"error",
	)
}

func TestExtractChanges_AddAndRemove(t *testing.T) {

	content, _ := ioutil.ReadFile("testdata/050-add-and-remove.changeset.json")
	var changeSet cf.DescribeChangeSetOutput
	json.Unmarshal([]byte(content), &changeSet)
	changes, _ := ExtractChanges(&changeSet)

	assert.Exactly(
		t,
		changes,
		[]GiffChange{
			{
				Action:             cfTypes.ChangeActionAdd,
				LogicalResourceId:  aws.String("AutoScalingGroup"),
				PhysicalResourceId: nil,
				Replacement:        "",
				ResourceType:       aws.String("AWS::AutoScaling::AutoScalingGroup"),
			},
			{
				Action:             cfTypes.ChangeActionAdd,
				LogicalResourceId:  aws.String("LaunchConfig"),
				PhysicalResourceId: nil,
				Replacement:        "",
				ResourceType:       aws.String("AWS::AutoScaling::LaunchConfiguration"),
			},
			{
				Action:             cfTypes.ChangeActionRemove,
				LogicalResourceId:  aws.String("MyEC2Instance"),
				PhysicalResourceId: aws.String("i-1abc23d4"),
				Replacement:        "",
				ResourceType:       aws.String("AWS::EC2::Instance"),
			},
		},
		"error",
	)
}
