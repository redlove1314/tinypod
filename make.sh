#!/bin/sh

echo "+-------------------------------------+"
echo "|              >> httpd <<           |"
echo "+-------------------------------------+"

rm -rf bin

export GOROOT=$(go env GOROOT)
if [ -z "$(go env GOPATH)" ]; then
    export GOPATH=${HOME}/go:$PWD
else
    export GOPATH=$(go env GOPATH):$PWD
fi

echo "step 1/2: create build output directory."
if [ ! -e "./bin" ];then
    mkdir ./bin
fi

echo "step 2/2: build..."
go build -i -o bin/spider src/spider.go
echo "build success!"