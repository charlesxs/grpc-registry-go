package main

import (
	"context"
	"github.com/charlesxs/grpc-registry-go/config"
	hello "github.com/charlesxs/grpc-registry-go/examples/helloworld/stub"
	"github.com/charlesxs/grpc-registry-go/gserver"
	"gitlab.corp.qunar.com/tcdev/qconfig-go/common/param"
	"gitlab.corp.qunar.com/tcdev/qconfig-go/qconfig"
	"log"
)

const (
	testAppCode = "ops_watcher_gwhb"
	testToken   = "BvetjV3K+Mj6FJ6Qigy+yfc9AyTZfjA/0vlOEq1ZlvhF/csWuT77AN7PMtnT4H7IgiEdT0WlYDEhCTn922tY+HfuMSpzeOgMoTJbm0wpDpdKOgxN29AKf9vU69GjpLOEPTs7YHY1iC3DwuzEESCmrt7A0IW1/Eybxd4EstBFno4="
)

func init() {
	clientParam := param.QConfigClientParam{
		AppInfo: param.AppInfoParam{
			AppCode: testAppCode,
			Token:   testToken,
		},
	}

	err := qconfig.Init(clientParam, param.QConfigFiles{
		config.ServerConfigFile: &config.ServerConfig,
	})
	if err != nil {
		log.Fatalln(err)
	}
}

type greeterServiceImpl struct {
	hello.UnimplementedGreeterServer
}

func (g *greeterServiceImpl) SayHello(ctx context.Context, in *hello.HelloRequest) (*hello.HelloReply, error) {
	log.Printf("received request: %s\n", in.GetName())
	return &hello.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
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
