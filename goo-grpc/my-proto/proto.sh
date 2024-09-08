#!/bin/bash

cd .

find . -name "*.pb.go" | xargs rm -f

protoc ./*.proto \
    --go_out=plugins=grpc,paths=source_relative:.
