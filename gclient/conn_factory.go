package gclient

import (
	"context"
	"github.com/charlesxs/grpc-registry-go/config"
	"google.golang.org/grpc"
)

// Factories 定义了已经注册的工厂(IConnFactory), 其中map的key是schema，目前实现了etcd的discovery方式
var Factories map[string]IConnFactory

func init() {
	Factories = map[string]IConnFactory{
		"etcd": newEtcdConnFactory(),
	}
}

// IConnFactory 接口定义了用于创建grpc.ClientConn的工厂, 解析、构建配置并创建grpc.ClientConn
type IConnFactory interface {
	BuildOptions(ctx context.Context, serverCfg *config.RpcServerConfig) error

	CreateConn(serverCfg *config.RpcServerConfig, options ...grpc.DialOption) (*grpc.ClientConn, error)
}
