#!/bin/bash

if [[ "$(dirname $0)" != "." ]]; then
  cd $(dirname $0)
fi

[ -d pb ] && find ./pb -name "*.pb.go" | xargs rm -f

protoc **/*.proto --go_out=plugins=grpc:.
