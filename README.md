# Microservices-Template
the Microservices Template for Golang

## 安装Protobuf编译器

下载地址：`https://github.com/protocolbuffers/protobuf/releases`

## 安装 Protobuf-Go 插件
安装 protoc-gen-go 插件, 用于生成 *.pb.go 文件：
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```
安装 protoc-gen-go-grpc 插件：用于生成 *_grpc.pb.go 文件：
```bash
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## 编写proto文件

```proto
syntax = "proto3";
package user;

message User {
    string name = 1;
    int32 age = 2;
}
```

## 生成 *.pb.go 文件
- --go_out=. --go_opt=paths=source_relative 用以在 .proto 文件同目录下生成 goods.pb.go
- --go-grpc_opt=paths=source_relative 用以在 .proto 文件同目录下生成 goods_grpc.pb.go

```bash
# paths 可选：source_relative，import，$PREFIX
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative user.proto
```

## 提供RESTful API服务

### 生成样板代码
```bash
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2
protoc -I=proto \
   --go_out=entity --go_opt=paths=source_relative \
   --go-grpc_out=entity --go-grpc_opt=paths=source_relative \
   --grpc-gateway_out=entity --grpc-gateway_opt=paths=source_relative \
   proto/entity.proto
```

# 参考链接：
1. https://protobuf.dev/reference/go/go-generated/
