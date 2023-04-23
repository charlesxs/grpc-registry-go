package main

import (
	"context"
	"github.com/charlesxs/grpc-registry-go/config"
	hello "github.com/charlesxs/grpc-registry-go/examples/helloworld/stub"
	"github.com/charlesxs/grpc-registry-go/examples/helloworld/utils"
	"github.com/charlesxs/grpc-registry-go/gserver"
	"log"
)

func initConfig() (*config.ServerConfig, error) {
	cfg := &config.ServerConfig{}
	return cfg, utils.ReadConfig("server_config.json", cfg)
}

// 服务实现
type greeterServiceImpl struct {
	hello.UnimplementedGreeterServer
}

func (g *greeterServiceImpl) SayHello(ctx context.Context, in *hello.HelloRequest) (*hello.HelloReply, error) {
	log.Printf("received request: %s\n", in.GetName())
	return &hello.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	// 初始化配置
	cfg, err := initConfig()
	if err != nil {
		log.Fatalln(err)
	}

	// 初始化 gserver
	gs, err := gserver.New(cfg).Build()
	if err != nil {
		log.Fatalln(err)
	}

	// 注册服务
	hello.RegisterGreeterServer(gs.Server(), &greeterServiceImpl{})

	// gserver run
	log.Fatalln(gs.Run())
}
