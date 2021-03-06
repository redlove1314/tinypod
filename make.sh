#!/bin/sh

echo "+-------------------------------------+"
echo "|               >> pod <<             |"
echo "+-------------------------------------+"

rm -rf bin

export GOROOT=$(go env GOROOT)
if [ -z "$(go env GOPATH)" ]; then
    export GOPATH=${HOME}/go:$PWD
else
    export GOPATH=$(go env GOPATH):$PWD
fi

echo "step 1/3: create build output directory."
if [ ! -e "./bin" ];then
    mkdir ./bin
fi

echo "step 2/3: install libs..."
go get github.com/urfave/cli

echo "step 3/3: build..."
go build -i -o bin/pod src/pod.go

echo "build success!"
