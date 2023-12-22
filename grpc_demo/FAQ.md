### 学习顺序
- simple_rpc
- stream_rpc
- alts_rpc
- pkg_rpc
- pass_param_rpc
- interceptor_rpc 

### FAQ-1
执行 `protoc -I . test.proto --go_out=plugins=grpc:.` 生成文件会报错 
`protoc-gen-go: plugins are not supported; use 'protoc --go-grpc_out=...' to generate gRPC`

> 最新版本的protoc不再支持使用插件来生成代码。所以你需要更新你的命令为, 可以尝试使用protoc --go-grpc_out=...命令来生成gRPC代码

### FAQ-2
一定主要 `message` 中字段的顺序问题

```protobuf
message TestRequest {
  string name = 1;
  string sex = 1;
}
````
修改为
```protobuf
message TestRequest {
  string sex = 1;
  string name = 1;
}
````
如果在生成了文件后 pb 文件后修改了字段的顺序，那么生成的代码中字段顺序也会改变，这样就会导致服务端和客户端的调用失败
