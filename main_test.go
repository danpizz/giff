package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/danpizz/giff/cmd"
	"github.com/stretchr/testify/assert"
)

// To run these tests you must have access to a AWS account and have the
// testdata/sample-1.yaml template deployed with the name "sample-giff-stack"

func TestCLI_diff_sample_2(t *testing.T) {
	os.Args = []string{"./giff", "diff", "sample-giff-stack", "testdata/sample-2.yaml"}
	b := bytes.NewBufferString("")
	cmd.Out = b
	main()
	out, _ := ioutil.ReadAll(b)
	assert.Contains(t, string(out), "+  # Adding this\n+  SampleRole2:")
}

func TestCLI_changes_sample_2(t *testing.T) {
	os.Args = []string{"./giff", "changes", "sample-giff-stack", "testdata/sample-2.yaml", "-t", "tag=tagdata"}
	b := bytes.NewBufferString("")
	cmd.Out = b
	main()
	out, _ := ioutil.ReadAll(b)
	assert.Exactly(t,
		"+     add: SampleRole2 - AWS::IAM::Role\n"+
			"*  modify: SampleRole (sample-giff-stack-sample-role) - AWS::IAM::Role / replacement: False / scope: Tags\n",
		string(out))
}

func TestCLI_changes_sample_3(t *testing.T) {
	os.Args = []string{"./giff", "changes", "sample-giff-stack", "testdata/sample-3.yaml", "-p", "MyTag=hello"}
	b := bytes.NewBufferString("")
	cmd.Out = b
	main()
	out, _ := ioutil.ReadAll(b)
	assert.Exactly(t,
		"+     add: SampleRole2 - AWS::IAM::Role\n"+
			"-  remove: SampleRole - AWS::IAM::Role\n",
		string(out))
}

func TestCLI_changes_param_and_tag(t *testing.T) {
	os.Args = []string{"./giff", "changes", "sample-giff-stack-2", "testdata/sample-volume.yaml", "-p", "Size=2", "-t", "MyTag=hello"}
	b := bytes.NewBufferString("")
	cmd.Out = b
	main()
	out, _ := ioutil.ReadAll(b)
	assert.Exactly(t,
		"*  modify: Volume (vol-049ee452fc2a8cd03) - AWS::EC2::Volume / replacement: False / scope: Properties Tags\n",
		string(out))
}
