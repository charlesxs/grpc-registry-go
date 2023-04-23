## grpc-registry-go 封装了带有注册中心功能的grpc库

具体例子请参考: [grpc-registry-go/examples/hellworld](https://github.com/charlesxs/grpc-registry-go/tree/master/examples/helloworld)

### Overview

grpc-registry-go 基于grpc-go之上封装了注册中心的功能，支持基于健康检查的自动上下线。
grpc-registry-go 是以应用为维度的服务注册和服务发现，当前实现了etcd 方式的服务注册和服务发现，后续可实现其他类型的注册中心。


### QuickStart

-----

**Server端**

- 配置server端 config

> 最小化配置, 更多配置请查看 [grpc-registry-go/server/server_config.go](https://github.com/charlesxs/grpc-registry-go/blob/master/config/server_config.go)

```json

{
  "app_name": "server_app_name",
  "port": 8888,
  "schema": "etcd",
  "etcd_registry_config": {
    "endpoints": ["etcd.server.addr"]
  }
}

```

- protobuf声明一个rpc 服务，并暴露服务端 stub

> example 定义[grpc-registry-go/examples/hellworld/stub/hello.proto](https://github.com/charlesxs/grpc-registry-go/blob/master/examples/helloworld/stub/hello.proto)

```protobuf
syntax = "proto3";

package hello;
option go_package = "demo/hello";

service Greeter {
    rpc SayHello(HelloRequest) returns (HelloReply) {}
}

message HelloRequest {
    string name = 1;
}

message HelloReply {
    string message = 1;
}
```

- 编译 protobuf

```go
$ protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative helloworld/stub/hello.proto
```

- 业务代码

```go
type greeterServiceImpl struct {
	hello.UnimplementedGreeterServer
}

func (g *greeterServiceImpl) SayHello(ctx context.Context, in *hello.HelloRequest) (*hello.HelloReply, error) {
	log.Printf("received request: %s\n", in.GetName())
	return &hello.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	// 初始化qconfig
	....
	
	// 初始化 gserver
	gs, err := gserver.New().WithDisableHealthcheck().Build()
	if err != nil {
		log.Fatalln(err)
	}

	// 注册服务
	hello.RegisterGreeterServer(gs.Server(), &greeterServiceImpl{})

	// gserver run
	log.Fatalln(gs.Run())
}

```

**Client端**

- 配置 client 端config

> 最小化配置，更多配置请查看 [grpc-registry-go/config/client_config.go](https://github.com/charlesxs/grpc-registry-go/blob/master/config/client_config.go)

```json
{
  "servers_discovery": [
    {
      "server_app": "serverAppName",
      "schema": "etcd",
      "etcd_config": {
        "endpoints": ["etcd.server.addr"]
      }
    }
  ]
}

```

- 业务代码

```go
func main() {
	// 初始化config
    cfg := initConfig()
	....
	
	// 初始化 gclient
	c, err := gclient.New(cfg).Build()
	if err != nil {
		log.Fatalln(err)
	}

	greeterClient := hello.NewGreeterClient(c.GetConn("serverApp"))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := greeterClient.SayHello(ctx, &hello.HelloRequest{Name: "charles"})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Greeting: ", r)
}

```

### 指南

----

### gserver

特性：

```go
1. WithDisableHealthcheck() 
禁用healthcheck协议，默认支持healthcheck协议，只有应用的 healthcheck.html存在才会向注册中心注册自己

2. WithServerOptions(options ...grpc.ServerOption)
指定自定义的grpc.ServerOption， 默认的grpc.ServerOption有 grpc.ConnectionTimeout(60 * time.Second), 60s连接超时

3. WithLogger(logger *zap.Logger)
指定logger

```

 ### gclient

特性:

```go
1. WithDialOptions(options ...grpc.DialOption)
指定自定义的grpc.DialOption, 默认的DialOption 有
// 默认通讯不指定证书
grpc.WithTransportCredentials(insecure.NewCredentials())  
// 默认的负载均衡策略是轮询, 若要更换策略，可以在qconfig grpc_client_config.json中的balance_policy中指定
grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`)
注意: 通常不要自己指定resolvers grpc.WithResolvers(), resolver和对应的注册中心成对出现, gclient中已经默认实现了

2. WithLogger(logger *zap.Logger)
指定logger

3. WithContext(ctx context.Context)
指定 context, 此context会被用于连接超时使用
```

### 其他

```go
1. 如果要实现其他类型的registry, 只需要实现 IRegistry 和 IRegistryFactory 两个接口即可

2. 如果要实现其他类型的服务发现, 只需要实现 IConnFactory 接口即可
```

