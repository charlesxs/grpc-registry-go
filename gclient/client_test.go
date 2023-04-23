package gclient

import (
	"context"
	"fmt"
	"github.com/charlesxs/grpc-registry-go/config"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	cfg := &config.ClientConfig{
		ServersDiscovery: []*config.RpcServerConfig{
			{
				ServerApp: "testApp",
				Schema:    "etcd",
				EtcdConfig: &config.EtcdConfig{
					Endpoints: []string{"etcd.server.addr:2379"},
				},
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, _ := New(cfg).
		WithContext(ctx).
		Build()

	fmt.Println(client.GetConn("testApp"))
}
