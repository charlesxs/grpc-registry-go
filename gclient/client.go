package gclient

import (
	"context"
	"fmt"
	"github.com/charlesxs/grpc-registry-go/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GrpcClient 此对象应该全局单例，封装了grpc客户端以及客户端服务发现
type GrpcClient struct {
	// grpc client 可以同时支持多个应用的discovery和服务调用
	// 其中 map 的 key 便是server appName，可以根据server appName获取对应的 grpc.ClientConn
	conns  map[string]*grpc.ClientConn
	config *config.ClientConfig

	options        []grpc.DialOption // grpc 连接参数
	defaultOptions []grpc.DialOption // 默认连接参数
	logger         *zap.Logger
	ctx            context.Context
}

func New(cfg *config.ClientConfig) *GrpcClient {
	return &GrpcClient{
		conns: make(map[string]*grpc.ClientConn),
		defaultOptions: []grpc.DialOption{
			// 默认通讯不指定证书
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			// 默认的负载均衡策略是轮询, 若要更换策略，可以在qconfig grpc_client_config.json中的balance_policy中指定
			grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		},
		ctx:    context.Background(),
		config: cfg,
	}
}

func (gc *GrpcClient) WithContext(ctx context.Context) *GrpcClient {
	gc.ctx = ctx
	return gc
}

func (gc *GrpcClient) WithLogger(logger *zap.Logger) *GrpcClient {
	gc.logger = logger
	return gc
}

func (gc *GrpcClient) WithDialOptions(options ...grpc.DialOption) *GrpcClient {
	gc.options = options
	return gc
}

func (gc *GrpcClient) Build() (*GrpcClient, error) {
	if gc.config == nil {
		return nil, fmt.Errorf("[%w] client config uninitialized", ErrConfig)
	}

	// init logger
	if gc.logger == nil {
		if logger, err := zap.NewProduction(); err != nil {
			return nil, fmt.Errorf("[%w] init logger error, err=%s", ErrClientInit, err)
		} else {
			gc.logger = logger
		}
	}

	// init balance policy
	if gc.config.BalancePolicy != "" {
		s := fmt.Sprintf(`{"loadBalancingPolicy":"%s"}`, gc.config.BalancePolicy)
		gc.defaultOptions = append(gc.defaultOptions, grpc.WithDefaultServiceConfig(s))
	}

	// merge options
	options := gc.defaultOptions[:]
	options = append(options, gc.options...)
	gc.options = options

	// init conns
	for _, serverCfg := range gc.config.ServersDiscovery {
		conn, err := gc.createConn(serverCfg)
		if err != nil {
			return nil, err
		}
		gc.conns[serverCfg.ServerApp] = conn
	}
	return gc, nil
}

func (gc *GrpcClient) GetConn(serverAppName string) *grpc.ClientConn {
	if conn, ok := gc.conns[serverAppName]; ok {
		return conn
	}
	return nil
}

func (gc *GrpcClient) createConn(serverCfg *config.RpcServerConfig) (*grpc.ClientConn, error) {
	factory, ok := Factories[serverCfg.Schema]
	if !ok {
		return nil, fmt.Errorf("unsupport schema, schema=%s", serverCfg.Schema)
	}

	err := factory.BuildOptions(gc.ctx, serverCfg)
	if err != nil {
		return nil, err
	}

	return factory.CreateConn(serverCfg, gc.options...)
}
