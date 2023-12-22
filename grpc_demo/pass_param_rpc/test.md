```shell
cd /grpc_test/pass_param_rpc/server/main.go  && go run main.go

cd /grpc_test/pass_param_rpc/client/main.go  && go run main.go
````

---

````
server:
Start Server 8090 !!! 
2023/11/16 16:55:57 --> :authority[localhost:8090] <--
2023/11/16 16:55:57 --> content-type[application/grpc] <--
2023/11/16 16:55:57 --> user-agent[grpc-go/1.59.0] <--
2023/11/16 16:55:57 --> token[123456] <--
2023/11/16 16:55:57 Get Client will Msg


client:
Get server return message:"Hello will"
````