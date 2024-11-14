# 命令

```
# brew install protobuf

# 安装必要的工具
go get -u google.golang.org/grpc@v1.68.0
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 设置环境变量
GOBIN=$(go env GOPATH)/bin
export PATH=$PATH:$GOBIN

# 生成 protobuf 代码
protoc \
    --proto_path=. \
    --go_out=paths=source_relative:. \
    --go-grpc_out=paths=source_relative:. \
    ./goo-proto/**/*.proto
```

# 使用说明

- 导入包: `import "goo-proto/v1/message.proto";`
- 使用`message`: `goo.proto.v1.Empty`
- `goland` -> 语言和框架 -> `Protocol Buffers` -> `import paths`: `$GOPATH/src/github.com/liqiongtao/googo.io`

