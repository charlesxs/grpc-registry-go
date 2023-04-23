package gserver

import (
	"fmt"
	"gitlab.corp.qunar.com/qgrpc-go/config"
	"gitlab.corp.qunar.com/qgrpc-go/healthcheck"
	"gitlab.corp.qunar.com/qgrpc-go/registry"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"time"
)

// GrpcServer 对象应该全局单例，封装了grpc 注册中心以及服务注册
type GrpcServer struct {
	server *grpc.Server

	registry    registry.IRegistry // 注册中心
	healthcheck bool               // 指定是否遵循healthcheck协议, 默认true即遵循 (检测是否有healthcheck.html，有的话才将自己注册到注册中心)
	checker     *healthcheck.Checker

	options        []grpc.ServerOption
	defaultOptions []grpc.ServerOption
	localAddrs     []string

	logger *zap.Logger
}

func New() *GrpcServer {
	return &GrpcServer{
		defaultOptions: []grpc.ServerOption{
			grpc.ConnectionTimeout(60 * time.Second),
		},
		healthcheck: true,
	}
}

// WithDisableHealthcheck 指定healthcheck选项为true, 遵循qunar healthcheck协议
func (gs *GrpcServer) WithDisableHealthcheck() *GrpcServer {
	gs.healthcheck = false
	return gs
}

func (gs *GrpcServer) WithLogger(logger *zap.Logger) *GrpcServer {
	gs.logger = logger
	return gs
}

func (gs *GrpcServer) WithServerOptions(options ...grpc.ServerOption) *GrpcServer {
	gs.options = options
	return gs
}

func (gs *GrpcServer) Build() (*GrpcServer, error) {
	if config.ServerConfig == nil {
		return nil, fmt.Errorf("[%w]初始化grpc gserver 错误, qconfig配置没有初始化, config_file=%s",
			ErrServerInit, config.ServerConfigFile)
	}

	// init logger
	if gs.logger == nil {
		if logger, err := zap.NewProduction(); err != nil {
			return nil, fmt.Errorf("[%w]初始化日志错误, err=%s", ErrServerInit, err)
		} else {
			gs.logger = logger
		}
	}

	// get local addr
	addrs, err := LocalAddrs()
	if err != nil {
		return nil, fmt.Errorf("[%w] 获取本机ip地址错误, err=%s", ErrServerInit, err)
	}
	gs.localAddrs = addrs

	// merge options
	options := gs.defaultOptions[:]
	options = append(options, gs.options...)
	gs.options = options

	r, err := gs.buildRegistry()
	if err != nil {
		return nil, err
	}

	gs.registry = r
	gs.server = grpc.NewServer(gs.options...)

	// 初始化 checker
	hcInterval := 3 * time.Second // 默认检测间隔是3秒
	if config.ServerConfig.HealthcheckInterval > 0 {
		hcInterval = time.Duration(config.ServerConfig.HealthcheckInterval) * time.Second
	}
	gs.checker = healthcheck.NewChecker(hcInterval, gs.register, gs.unRegister, gs.logger)

	return gs, nil
}

func (gs *GrpcServer) Run() error {
	if err := gs.registerWithHC(); err != nil {
		return err
	}

	// 启动grpc gserver
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.ServerConfig.Port))
	if err != nil {
		return err
	}
	return gs.server.Serve(listener)
}

func (gs *GrpcServer) Server() *grpc.Server {
	return gs.server
}

func (gs *GrpcServer) registerWithHC() error {
	// 指定了要检测hc时，使用checker 后台检测
	if gs.healthcheck {
		go gs.checker.CheckForever()
		return nil
	}

	// 不需要检测hc时，直接注册
	if err := gs.register(); err != nil {
		return err
	}
	return nil
}

func (gs *GrpcServer) register() error {
	// 注册本机ip
	for _, addr := range gs.localAddrs {
		if err := gs.registry.Register(addr, config.ServerConfig.Port, nil); err != nil {
			return err
		}
	}
	return nil
}

func (gs *GrpcServer) unRegister() error {
	port := config.ServerConfig.Port
	for _, addr := range gs.localAddrs {
		if err := gs.registry.Unregister(addr, port); err != nil {
			return err
		}
	}
	return nil
}

func (gs *GrpcServer) buildRegistry() (registry.IRegistry, error) {
	schema := registry.Schema(config.ServerConfig.Schema)
	factory, ok := registry.Factories[schema]
	if !ok {
		return nil, fmt.Errorf("[%w]不支持对应的schema", ErrServerInit)
	}

	err := factory.BuildOptions(gs.logger)
	if err != nil {
		return nil, err
	}
	return factory.CreateRegistry()
}
