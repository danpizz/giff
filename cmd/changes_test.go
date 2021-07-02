package cmd

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	cf "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cfTypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/stretchr/testify/assert"
)

func TestChanges_no_args(t *testing.T) {
	cmd := NewChangesCmd(nil, nil)
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetErr(b)
	cmd.Execute()
	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, string(out), "Error: accepts 1 or 2 args, received 0")
}

type MockAPI struct {
}

func (MockAPI) ReadTemplateFile(templateFileName string) (body string, err error) {
	return "<template>", nil
}

type MockCFClientNoChanges struct {
}

func (client MockCFClientNoChanges) CreateChangeSet(params *cf.CreateChangeSetInput) (*cf.CreateChangeSetOutput, error) {
	return &cf.CreateChangeSetOutput{
		Id: aws.String("changesetID"),
	}, nil
}
func (client MockCFClientNoChanges) DescribeStacks(params *cf.DescribeStacksInput) (*cf.DescribeStacksOutput, error) {
	return &cf.DescribeStacksOutput{
		Stacks: []cfTypes.Stack{{Parameters: []cfTypes.Parameter{}}},
	}, nil
}
func (client MockCFClientNoChanges) GetTemplate(params *cf.GetTemplateInput) (*cf.GetTemplateOutput, error) {
	return &cf.GetTemplateOutput{}, nil
}
func (client MockCFClientNoChanges) DescribeChangeSet(params *cf.DescribeChangeSetInput) (*cf.DescribeChangeSetOutput, error) {
	return &cf.DescribeChangeSetOutput{
		Status: cfTypes.ChangeSetStatusCreateComplete,
	}, nil
}
func (client MockCFClientNoChanges) DeleteChangeSet(params *cf.DeleteChangeSetInput) (*cf.DeleteChangeSetOutput, error) {
	return &cf.DeleteChangeSetOutput{}, nil
}

func TestChanges_no_changes(t *testing.T) {
	cmd := NewChangesCmd(MockCFClientNoChanges{}, MockAPI{})
	cmd.SetArgs([]string{"stack", "template", "-p", "p1=v1"})
	b := bytes.NewBufferString("")
	cmd.SetOutput(b)
	cmd.Execute()
	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, string(out), "No changes")
}

type MockCFClientChanges struct {
}

var MockAction cfTypes.ChangeAction

func (client MockCFClientChanges) CreateChangeSet(params *cf.CreateChangeSetInput) (*cf.CreateChangeSetOutput, error) {
	return &cf.CreateChangeSetOutput{
		Id: aws.String(""),
	}, nil
}
func (client MockCFClientChanges) DescribeStacks(params *cf.DescribeStacksInput) (*cf.DescribeStacksOutput, error) {
	return &cf.DescribeStacksOutput{
		Stacks: []cfTypes.Stack{{Parameters: []cfTypes.Parameter{}}},
	}, nil
}
func (client MockCFClientChanges) GetTemplate(params *cf.GetTemplateInput) (*cf.GetTemplateOutput, error) {
	return &cf.GetTemplateOutput{}, nil
}
func (client MockCFClientChanges) DescribeChangeSet(params *cf.DescribeChangeSetInput) (*cf.DescribeChangeSetOutput, error) {
	return &cf.DescribeChangeSetOutput{
		Status: cfTypes.ChangeSetStatusCreateComplete,
		Changes: []cfTypes.Change{
			{
				ResourceChange: &cfTypes.ResourceChange{
					Action:             MockAction,
					LogicalResourceId:  aws.String("LogRId"),
					PhysicalResourceId: aws.String("PhyRId"),
					Replacement:        cfTypes.ReplacementTrue,
					ResourceType:       aws.String("RT"),
					Scope:              []cfTypes.ResourceAttribute{},
				},
				Type: cfTypes.ChangeTypeResource,
			},
		},
	}, nil
}
func (client MockCFClientChanges) DeleteChangeSet(params *cf.DeleteChangeSetInput) (*cf.DeleteChangeSetOutput, error) {
	return &cf.DeleteChangeSetOutput{}, nil
}

func TestChanges_all(t *testing.T) {
	cmd := NewChangesCmd(MockCFClientChanges{}, MockAPI{})
	cmd.SetArgs([]string{"stack", "template", "-p", "p1=v1"})
	b := bytes.NewBufferString("")
	cmd.SetOutput(b)

	t.Run("add", func(t *testing.T) {
		MockAction = cfTypes.ChangeActionAdd
		cmd.Execute()
		out, err := ioutil.ReadAll(b)
		if err != nil {
			t.Fatal(err)
		}
		assert.Contains(t, string(out), "+     add: LogRId - RT")
	})
	t.Run("remove", func(t *testing.T) {
		MockAction = cfTypes.ChangeActionRemove
		cmd.Execute()
		out, err := ioutil.ReadAll(b)
		if err != nil {
			t.Fatal(err)
		}
		assert.Contains(t, string(out), "-  remove: LogRId - RT")
	})
	t.Run("modify", func(t *testing.T) {
		MockAction = cfTypes.ChangeActionModify
		cmd.Execute()
		out, err := ioutil.ReadAll(b)
		if err != nil {
			t.Fatal(err)
		}
		assert.Contains(t, string(out), "*  modify: LogRId (PhyRId) - RT / replacement: True")
	})
	t.Run("dynamic", func(t *testing.T) {
		MockAction = cfTypes.ChangeActionDynamic
		cmd.Execute()
		out, err := ioutil.ReadAll(b)
		if err != nil {
			t.Fatal(err)
		}
		assert.Contains(t, string(out), "* dynamic: LogRId (PhyRId) - RT / replacement: True")
	})
	t.Run("import", func(t *testing.T) {
		MockAction = cfTypes.ChangeActionImport
		cmd.Execute()
		out, err := ioutil.ReadAll(b)
		if err != nil {
			t.Fatal(err)
		}
		assert.Contains(t, string(out), "+  import: LogRId (PhyRId) - RT")
	})
	t.Run("empty", func(t *testing.T) {
		MockAction = ""
		cmd.Execute()
		out, err := ioutil.ReadAll(b)
		if err != nil {
			t.Fatal(err)
		}
		assert.Contains(t, string(out), "[unknown change type]")
	})

}
