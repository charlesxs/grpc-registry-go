package gclient

import (
	"context"
	"fmt"
	"gitlab.corp.qunar.com/qgrpc-go/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GrpcClient 此对象应该全局单例，封装了grpc客户端以及客户端服务发现
type GrpcClient struct {
	// grpc client 可以同时支持多个应用的discovery和服务调用
	// 其中 map 的 key 便是server appCode，可以根据server appCode获取对应的 grpc.ClientConn
	conns map[string]*grpc.ClientConn

	options        []grpc.DialOption // grpc 连接参数
	defaultOptions []grpc.DialOption // 默认连接参数
	logger         *zap.Logger
	ctx            context.Context
}

func New() *GrpcClient {
	return &GrpcClient{
		conns: make(map[string]*grpc.ClientConn),
		defaultOptions: []grpc.DialOption{
			// 默认通讯不指定证书
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			// 默认的负载均衡策略是轮询, 若要更换策略，可以在qconfig grpc_client_config.json中的balance_policy中指定
			grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		},
		ctx: context.Background(),
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
	if config.ClientConfig == nil {
		return nil, fmt.Errorf("[%w]qconfig配置没有初始化, config_file=%s",
			ErrClientInit, config.ClientConfigFile)
	}

	// init logger
	if gc.logger == nil {
		if logger, err := zap.NewProduction(); err != nil {
			return nil, fmt.Errorf("[%w]初始化日志错误, err=%s", ErrClientInit, err)
		} else {
			gc.logger = logger
		}
	}

	// init balance policy
	if config.ClientConfig.BalancePolicy != "" {
		s := fmt.Sprintf(`{"loadBalancingPolicy":"%s"}`, config.ClientConfig.BalancePolicy)
		gc.defaultOptions = append(gc.defaultOptions, grpc.WithDefaultServiceConfig(s))
	}

	// merge options
	options := gc.defaultOptions[:]
	options = append(options, gc.options...)
	gc.options = options

	// init conns
	for _, serverCfg := range config.ClientConfig.ServersDiscovery {
		conn, err := gc.createConn(serverCfg)
		if err != nil {
			return nil, err
		}
		gc.conns[serverCfg.ServerAppCode] = conn
	}
	return gc, nil
}

func (gc *GrpcClient) GetConn(serverAppCode string) *grpc.ClientConn {
	if conn, ok := gc.conns[serverAppCode]; ok {
		// 查看生效的选项
		//rv := reflect.ValueOf(conn).Elem()
		//fmt.Println(rv.FieldByName("dopts"))

		return conn
	}
	return nil
}

func (gc *GrpcClient) createConn(serverCfg *config.RpcServerConfig) (*grpc.ClientConn, error) {
	factory, ok := Factories[serverCfg.Schema]
	if !ok {
		return nil, fmt.Errorf("不支持的schema, schema=%s", serverCfg.Schema)
	}

	err := factory.BuildOptions(gc.ctx, serverCfg)
	if err != nil {
		return nil, err
	}

	return factory.CreateConn(serverCfg, gc.options...)
}
