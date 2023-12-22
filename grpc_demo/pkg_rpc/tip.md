这里可以重定向输出的 proto 文件位置, 通过设置 `go_package`

```shell
protoc --go_out=. --go-grpc_out=.  pkg.proto
````
```protobuf
syntax = "proto3";

import "google/protobuf/empty.proto";
import "google/protobuf/any.proto";
import "google/protobuf/api.proto";
import "google/protobuf/descriptor.proto";

option go_package = "/var/test/pb";

package pb;

message testRequest {
  string name = 1;
  string age = 2;
}

service  PkgServer {
  rpc Test  (testRequest) returns (google.protobuf.Empty);
}
````
````
└── var
    └── test
        └── pb
            ├── pkg.pb.go
            └── pkg_grpc.pb.go
````