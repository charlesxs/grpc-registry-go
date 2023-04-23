package gserver

import (
	"github.com/charlesxs/grpc-registry-go/config"
	"log"
	"testing"
)

func TestServer(t *testing.T) {
	cfg := &config.ServerConfig{
		AppName: "testApp",
		Port:    8081,
		Schema:  "etcd",
		EtcdRegistryConfig: &config.EtcdRegistryConfig{
			Endpoints: []string{"etcd.server.addr:2379"},
		},
	}

	server, err := New(cfg).Build()
	if err != nil {
		log.Fatalln(err)
	}

	log.Fatalln(server.Run())
}
