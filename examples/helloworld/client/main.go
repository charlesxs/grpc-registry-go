package main

import (
	"context"
	"github.com/charlesxs/grpc-registry-go/config"
	hello "github.com/charlesxs/grpc-registry-go/examples/helloworld/stub"
	"github.com/charlesxs/grpc-registry-go/gclient"
	"gitlab.corp.qunar.com/tcdev/qconfig-go/common/param"
	"gitlab.corp.qunar.com/tcdev/qconfig-go/qconfig"
	"log"
	"time"
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
		config.ClientConfigFile: &config.ClientConfig,
	})
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	c, err := gclient.New().Build()
	if err != nil {
		log.Fatalln(err)
	}

	greeterClient := hello.NewGreeterClient(c.GetConn(testAppCode))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := greeterClient.SayHello(ctx, &hello.HelloRequest{Name: "xs.xiao"})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Greeting: ", r)
}
