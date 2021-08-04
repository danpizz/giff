![ci workflow](https://github.com/danpizz/giff/actions/workflows/default.yml/badge.svg)
# The `giff` CloudFormation differ

`giff` is an AWS CloudFormation stack diff tool.

* **diff**: shows the differences between a CloudFormation stack and a local template
* **changes**: creates a temporary changeset and displays an easy to read summary of the changes created by a local template and some (optional) parameters.

`giff` was inspired by the `cliff` tool you can find here: https://github.com/meetup/cliff

```
+     add: SampleRole2 - AWS::IAM::Role
*  modify: SampleRole (sample-giff-stack-sample-role) - AWS::IAM::Role / replacement: False / scope: Tags
```

## Template diffing

```
+  SampleRole2:
+    Type: AWS::IAM::Role
+    Properties:
+      RoleName: !Sub ${AWS::StackName}-sample-role-2
+      AssumeRolePolicyDocument:
+        Version: 2012-10-17
```

**Giff** downloads the template that was used to create the specified stack, then runs the `diff` command over that template and a local one showing the results.
By default **giff** uses `diff -u` but you can provide any command with the `-d` flag.


```
giff diff my-stack my-template.yaml -d colordiff
```

## Showing changes with temporary changesets

```
giff changes sample-1-stack testdata/sample-2.yaml -p OtherPolicyArn=newArn
+     add: SampleRole2 - AWS::IAM::Role
*  modify: SampleRole (sample-giff-stack-sample-role) - AWS::IAM::Role / replacement: False / scope: Tags
```

If used with two arguments, the stack name and the template file, `giff changes` shows the changes caused by deploying the specified template file over the named stack. 
It will create a temporary changeset, show a easy to read list of changes, and then delete the changeset.

#### Flags

`--parameters-overrides` a partial list of parameters `Param1=Value1 Param2=Value2`

`--all-parameters` the complete listof parameters. This flag will cause an error if there are some missing parameters.

`--tags` tags to associate to the stack

`--no-delete-changeset` don't delete the temporary changeset and print its ARN

`--dump` print the full raw changeset in JSON format

## Showing changes of existing changesets

With one single argument, a changeset ARN, giff will show a the list of changes caused by the changeset.

```
giff change arn:aws:cloudformation:us-east-1:123456789012:changeSet/SampleChangeSet-direct/1a2345b6-0000-00a0-a123-00abc0abc000
+     add: SampleRole2 - AWS::IAM::Role
*  modify: SampleRole (sample-giff-stack-sample-role) - AWS::IAM::Role / replacement: False / scope: Tags
```
