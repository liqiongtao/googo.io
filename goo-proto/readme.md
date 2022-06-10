# 命令

- `protoc --proto_path=. --go_out=plugins=grpc,paths=source_relative:. goo-proto/**/*.proto`
- `option go_package` 声明 是为了让生成的其他 go 包（依赖方）可以正确 import 到本包（被依赖方）
- `--go_out=paths=source_relative:.` 参数 是为了让加了 `option go_package` 声明的 proto 文件可以将 go 代码编译到与其同目录

# 使用说明

- 导入包: `import "goo-proto/v1/message.proto";`
- 使用`message`: `goo.proto.v1.Empty`
- `goland` -> 语言和框架 -> `Protocol Buffers` -> `import paths`: `$GOPATH/src/github.com/liqiongtao/googo.io`
- 命令行编译
```
protoc --proto_path=. \
    --proto_path=$GOPATH/src/github.com/liqiongtao/googo.io \
    --go_out=plugins=grpc,paths=source_relative:. \
    ./**/*.proto
```