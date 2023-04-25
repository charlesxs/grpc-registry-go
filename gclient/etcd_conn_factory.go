package gclient

import (
	"context"
	"fmt"
	"github.com/charlesxs/grpc-registry-go/config"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"time"
)

// EtcdConnOptions 用于etcd注册中心服务发现配置
type EtcdConnOptions struct {
	ctx        context.Context
	EtcdConfig *clientv3.Config // EtcdConfig 指定etcd 配置
}

type etcdConnFactory struct {
	Options map[string]*EtcdConnOptions
}

func newEtcdConnFactory() *etcdConnFactory {
	return &etcdConnFactory{
		Options: make(map[string]*EtcdConnOptions),
	}
}

func (factory *etcdConnFactory) CreateConn(serverCfg *config.RpcServerConfig, options ...grpc.DialOption) (*grpc.ClientConn, error) {
	cfg, ok := factory.Options[serverCfg.ServerApp]
	if !ok {
		return nil, fmt.Errorf("[%w] not found server app config, appName=%s", ErrCreateConn, serverCfg.ServerApp)
	}

	if err := factory.checkConfig(serverCfg); err != nil {
		return nil, err
	}

	etcdClient, err := clientv3.New(*cfg.EtcdConfig)
	if err != nil {
		return nil, fmt.Errorf("[%w] etcd client init error, err=%s", ErrCreateConn, err)
	}

	// etcdClient connection test to avoid the resolver build stuck
	ctx, cancel := context.WithTimeout(cfg.ctx, cfg.EtcdConfig.DialTimeout)
	defer cancel()
	_, err = etcdClient.Get(ctx, serverCfg.ServerApp)
	if err != nil {
		return nil, fmt.Errorf("[%w] etcd connection error, err=%s", ErrCreateConn, err)
	}

	builder, err := resolver.NewBuilder(etcdClient)
	if err != nil {
		return nil, fmt.Errorf("[%w] build etcd resolver builder error, err=%s", ErrCreateConn, err)
	}

	// 加入resolvers
	options = append(options, grpc.WithResolvers(builder))

	return grpc.DialContext(
		cfg.ctx,
		fmt.Sprintf("%s:///%s", serverCfg.Schema, serverCfg.ServerApp),
		options...,
	)
}

func (factory *etcdConnFactory) BuildOptions(ctx context.Context, serverCfg *config.RpcServerConfig) error {
	c := serverCfg.EtcdConfig
	if c == nil {
		return fmt.Errorf("[%w] not found etcd config", ErrConfigNotFound)
	}

	etcdConnConfig := &EtcdConnOptions{
		ctx: ctx,
		EtcdConfig: &clientv3.Config{
			Endpoints: c.Endpoints,
			Username:  c.Username,
			Password:  c.Password,
		},
	}
	if c.DialTimeout > 0 {
		etcdConnConfig.EtcdConfig.DialTimeout = time.Duration(c.DialTimeout) * time.Second
	}
	if c.DialKeepAliveTime > 0 {
		etcdConnConfig.EtcdConfig.DialKeepAliveTime = time.Duration(c.DialKeepAliveTime) * time.Second
	}
	if c.DialKeepAliveTimeout > 0 {
		etcdConnConfig.EtcdConfig.DialKeepAliveTimeout = time.Duration(c.DialKeepAliveTimeout) * time.Second
	}
	factory.Options[serverCfg.ServerApp] = etcdConnConfig
	return nil
}

func (factory *etcdConnFactory) checkConfig(serverCfg *config.RpcServerConfig) error {
	if serverCfg.ServerApp == "" {
		return fmt.Errorf("[%w] did not specify app name", ErrConfig)
	}

	c := factory.Options[serverCfg.ServerApp]

	// etcd 配置校验以及设置默认值
	defaultTimeout := 5 * time.Second
	if c.EtcdConfig == nil || len(c.EtcdConfig.Endpoints) <= 0 {
		return fmt.Errorf("[%w] did not specify etcd endpoints", ErrConfig)
	}

	if c.EtcdConfig.DialTimeout == 0 {
		c.EtcdConfig.DialTimeout = defaultTimeout
	}
	if c.EtcdConfig.DialKeepAliveTime == 0 {
		c.EtcdConfig.DialKeepAliveTime = defaultTimeout
	}
	if c.EtcdConfig.DialKeepAliveTimeout == 0 {
		c.EtcdConfig.DialKeepAliveTimeout = defaultTimeout
	}
	return nil
}
