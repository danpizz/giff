package pkg

import (
	"encoding/json"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	cf "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cfTypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

func ParameterListFromString(parametersString string) []cfTypes.Parameter {
	var parameterList []cfTypes.Parameter
	for _, keyValue := range strings.Split(parametersString, " ") {
		p := strings.Split(keyValue, "=")
		if len(p) != 2 {
			continue
		}
		parameterList = append(parameterList, cfTypes.Parameter{
			ParameterKey:   &p[0],
			ParameterValue: &p[1],
		})
	}
	return parameterList
}

func OverrideParameters(stackParameters []cfTypes.Parameter, inputParameters []cfTypes.Parameter) ([]cfTypes.Parameter, error) {
	for _, v := range inputParameters {
		modified := false
		for i, sp := range stackParameters {
			if *sp.ParameterKey == *v.ParameterKey {
				stackParameters[i].ParameterValue = v.ParameterValue
				stackParameters[i].UsePreviousValue = aws.Bool(false)
				modified = true
			}
		}
		if !modified {
			// parameters added in the template but not in the stack
			stackParameters = append(stackParameters, cfTypes.Parameter{
				ParameterKey:     v.ParameterKey,
				ParameterValue:   v.ParameterValue,
				UsePreviousValue: aws.Bool(false),
			})
		}
	}

	for i, p := range stackParameters {
		if p.UsePreviousValue == nil {
			stackParameters[i].UsePreviousValue = aws.Bool(true)
			stackParameters[i].ParameterValue = nil
		}
	}
	return stackParameters, nil
}

type GiffChange struct {
	Action             cfTypes.ChangeAction
	LogicalResourceId  *string
	PhysicalResourceId *string
	Replacement        cfTypes.Replacement
	ResourceType       *string
}

func PrettyJson(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func ExtractChanges(describeChangeSetOutput *cf.DescribeChangeSetOutput) ([]GiffChange, error) {
	var changes []GiffChange
	if len(describeChangeSetOutput.Changes) <= 0 {
		return make([]GiffChange, 0), nil
	}
	for _, c := range describeChangeSetOutput.Changes {
		var change GiffChange
		if c.Type != cfTypes.ChangeTypeResource {
			continue
		}
		change.Action = c.ResourceChange.Action
		change.LogicalResourceId = c.ResourceChange.LogicalResourceId
		change.PhysicalResourceId = c.ResourceChange.PhysicalResourceId
		change.Replacement = c.ResourceChange.Replacement
		change.ResourceType = c.ResourceChange.ResourceType
		changes = append(changes, change)
	}
	return changes, nil
}
