#!/bin/bash

cd .

find . -name "*.pb.go" | xargs rm -f

# 安装必要的工具
# brew install protobuf
go get -u google.golang.org/grpc@v1.68.0
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

GOBIN=$(go env GOPATH)/bin
export PATH=$PATH:$GOBIN

protoc ./*.proto \
    --go_out=paths=source_relative:. \
    --go-grpc_out=paths=source_relative:.
