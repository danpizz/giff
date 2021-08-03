#!/usr/bin/env bash

short_desc=$(git describe --dirty)
long_desc=$(git describe --long --dirty)
ldflags="-X github.com/danpizz/giff/cmd.ShortVersion=${short_desc} -X github.com/danpizz/giff/cmd.Version=${long_desc}"

deploy-test-data() {
    aws cloudformation deploy \
        --template-file testdata/sample-1.yaml \
        --stack-name sample-giff-stack \
        --parameter-overrides OtherPolicyArn=arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore \
        --capabilities CAPABILITY_IAM CAPABILITY_NAMED_IAM \
        --tags Tag1="hello"
    aws cloudformation deploy \
        --template-file testdata/sample-volume.yaml \
        --stack-name sample-giff-stack-2 \
        --parameter-overrides Zone=eu-west-1a Size=1 \
        --capabilities CAPABILITY_IAM CAPABILITY_NAMED_IAM \
        --tags Tag1="hello"
}

clean() {
    rm giff giff_*
    go clean
    rm c.out coverage.html
}

build() {
    go build -ldflags="${ldflags}" -v
}

build_all() {
    GOOS=darwin GOARCH=amd64 go build -ldflags="${ldflags}" -o "giff_darwin_amd64_${long_desc}"
    GOOS=darwin GOARCH=arm64 go build -ldflags="${ldflags}" -o "giff_darwin_arm64_${long_desc}"
    GOOS=linux GOARCH=amd64 go build -ldflags="${ldflags}" -o "giff_linux_amd64_${long_desc}"
    GOOS=linux GOARCH=amd64 go build -ldflags="${ldflags}" -o "giff_linux_arm64_${long_desc}"
}

test() {
    go test -v ./cmd ./pkg
}

test_aws() {
    go test -v .
}

coverage() {
    go test ./cmd ./pkg -cover -coverprofile=c.out
    go tool cover -html=c.out -o coverage.html
}

"$@"
