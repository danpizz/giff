package pkg

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"

	cf "github.com/aws/aws-sdk-go-v2/service/cloudformation"
)

type CFAPI interface {
	CreateChangeSet(params *cf.CreateChangeSetInput) (*cf.CreateChangeSetOutput, error)
	DescribeChangeSet(params *cf.DescribeChangeSetInput) (*cf.DescribeChangeSetOutput, error)
	DescribeStacks(params *cf.DescribeStacksInput) (*cf.DescribeStacksOutput, error)
	DeleteChangeSet(params *cf.DeleteChangeSetInput) (*cf.DeleteChangeSetOutput, error)
	GetTemplate(params *cf.GetTemplateInput) (*cf.GetTemplateOutput, error)
}
type CFClient struct {
	*cf.Client
}

func (client CFClient) CreateChangeSet(params *cf.CreateChangeSetInput) (*cf.CreateChangeSetOutput, error) {
	return client.Client.CreateChangeSet(context.TODO(), params)
}
func (client CFClient) DescribeStacks(params *cf.DescribeStacksInput) (*cf.DescribeStacksOutput, error) {
	return client.Client.DescribeStacks(context.TODO(), params)
}
func (client CFClient) GetTemplate(params *cf.GetTemplateInput) (*cf.GetTemplateOutput, error) {
	return client.Client.GetTemplate(context.TODO(), params)
}
func (client CFClient) DescribeChangeSet(params *cf.DescribeChangeSetInput) (*cf.DescribeChangeSetOutput, error) {
	return client.Client.DescribeChangeSet(context.TODO(), params)
}
func (client CFClient) DeleteChangeSet(params *cf.DeleteChangeSetInput) (*cf.DeleteChangeSetOutput, error) {
	return client.Client.DeleteChangeSet(context.TODO(), params)
}

func NewCFClient() (*CFClient, error) {
	awsCfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	cfClient := cf.NewFromConfig(awsCfg)
	return &CFClient{
		Client: cfClient,
	}, nil
}
