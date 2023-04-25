package main

import (
	"context"
	"github.com/charlesxs/grpc-registry-go/config"
	hello "github.com/charlesxs/grpc-registry-go/examples/helloworld/stub"
	"github.com/charlesxs/grpc-registry-go/examples/helloworld/utils"
	"github.com/charlesxs/grpc-registry-go/gclient"
	"log"
	"time"
)

func initConfig() (*config.ClientConfig, error) {
	cfg := &config.ClientConfig{}
	return cfg, utils.ReadConfig("client_config.yaml", cfg)
}

func main() {
	// init config
	cfg, err := initConfig()
	if err != nil {
		log.Fatalln(err)
	}

	// init gclient
	c, err := gclient.New(cfg).Build()
	if err != nil {
		log.Fatalln(err)
	}

	// invoke remote method
	greeterClient := hello.NewGreeterClient(c.GetConn("testApp"))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := greeterClient.SayHello(ctx, &hello.HelloRequest{Name: "charles"})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Greeting: ", r)
}
