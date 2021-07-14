---
title: "Grpc Notes"
description: 
date: 2021-07-14T14:02:57+08:00
image: 
math: 
license: 
hidden: false
draft: true
categories:
    - program
tag:
    - gRPC
---

## protobuf

> Protocol buffers are a language-neutral, platform-neutral extensible mechanism for serializing structured data.

Protocol Buffers（简称 protobuf）是一种不依赖语言以及平台的可扩展的数据序列化机制，类似于XML与Json，不过protobuf是一个二进制协议，理论上性能更好。

在gRPC中，protobuf除了用于作为数据传输的格式外，还作为接口描述语言（Interface description language，缩写*IDL*）用于生成server/client的接口代码。

### 使用

#### 需要的工具

- [protobuf complier](https://github.com/protocolbuffers/protobuf/releases)（简称protoc）：用于编译protobuf文件
- protobuf plugin：一部分语言如C++、Java、C#等，protoc直接可以编译，一部分语言如Go等则需要额外插件的支持
  - Go
    - [protoc-gen-go](https://github.com/golang/protobuf/tree/master/protoc-gen-go)：官方插件
    - [gogoprotobuf](https://github.com/gogo/protobuf): 第三方插件

本文后续没有特别说明，都是使用官方的`protoc-gen-go`插件



#### 编译

示例proto文件

```protobuf
// hello.proto
syntax = "proto3";
package hello;
option go_package = "pb/hello";

service HelloService {
  rpc SayHello(HelloReq) returns (HelloResp) {};
}


message HelloReq {
  string Name = 1;
  optional string other = 2;
}

message HelloResp {
  string Result = 1;
  optional string other = 2;
}
```



##### 只编译数据序列化的代码

```shell
protoc --go_out=. hello.proto
```

生成的pb.go只包含了message的读写方法，不包含gRPC的代码生成（即proto中的service）部分



##### 包含gRPC代码生成

```shell
protoc --go_out=plugins=grpc:. example.proto
```

在指定了grpc的插件后，生成的pb.go文件包还含了server/client相关得方法





