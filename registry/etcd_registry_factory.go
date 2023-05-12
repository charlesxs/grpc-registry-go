package registry

import (
	"fmt"
	"github.com/charlesxs/grpc-registry-go/config"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"time"
)

// EtcdRegistryOptions etcd注册中心配置选项
type EtcdRegistryOptions struct {
	AppName    string           // AppName 指定应用名称，也会作为 etcd 的注册路径
	Interval   time.Duration    // Interval 指定多久检测并重新续约
	LeaseTTL   int64            // LeaseTTL 租约时长
	Logger     *zap.Logger      // Logger 日志
	EtcdConfig *clientv3.Config // EtcdConfig 指定etcd 配置
}

type etcdRegistryFactory struct {
	Options *EtcdRegistryOptions
}

// BuildOptions 从qconfig中构建registry配置
func (factory *etcdRegistryFactory) BuildOptions(cfg *config.ServerConfig) error {
	etcdConfig := cfg.EtcdRegistryConfig
	if etcdConfig == nil {
		return fmt.Errorf("[%w] did not specify etcd registry config", ErrConfigNotFound)
	}

	o := &EtcdRegistryOptions{
		AppName: cfg.AppName,
		Logger:  cfg.Logger,
		EtcdConfig: &clientv3.Config{
			Endpoints: etcdConfig.Endpoints,
			Username:  etcdConfig.Username,
			Password:  etcdConfig.Password,
		},
	}

	if etcdConfig.LeaseTTL > 0 {
		o.LeaseTTL = etcdConfig.LeaseTTL
	}
	if etcdConfig.Interval > 0 {
		o.Interval = time.Duration(etcdConfig.Interval) * time.Second
	}
	if etcdConfig.DialTimeout > 0 {
		o.EtcdConfig.DialTimeout = time.Duration(etcdConfig.DialTimeout) * time.Second
	}
	if etcdConfig.DialKeepAliveTime > 0 {
		o.EtcdConfig.DialKeepAliveTime = time.Duration(etcdConfig.DialKeepAliveTime) * time.Second
	}
	if etcdConfig.DialKeepAliveTimeout > 0 {
		o.EtcdConfig.DialKeepAliveTimeout = time.Duration(etcdConfig.DialKeepAliveTimeout) * time.Second
	}

	factory.Options = o
	return nil
}

// CreateRegistry 创建etcd registry
func (factory *etcdRegistryFactory) CreateRegistry() (IRegistry, error) {
	if err := factory.checkOptions(); err != nil {
		return nil, err
	}

	client, err := clientv3.New(*factory.Options.EtcdConfig)
	if err != nil {
		return nil, err
	}

	return newEtcdRegistry(client, factory.Options)
}

func (factory *etcdRegistryFactory) checkOptions() error {
	o := factory.Options
	// 检查app name
	if o.AppName == "" {
		return fmt.Errorf("[%w] app name must be not empty", ErrRegistryOption)
	}

	// 设置默认检测间隔
	if o.Interval == 0 {
		o.Interval = 10 * time.Second
	}
	if o.LeaseTTL == 0 {
		o.LeaseTTL = 12
	}

	if o.Interval > time.Duration(o.LeaseTTL)*time.Second {
		return fmt.Errorf("[%w] Interval can not largeer than LeaseTTL", ErrRegistryOption)
	}

	if o.Logger == nil {
		logger, _ := zap.NewProduction()
		o.Logger = logger
	}

	// etcd 配置校验以及设置默认值
	defaultTimeout := 5 * time.Second
	if o.EtcdConfig == nil || len(o.EtcdConfig.Endpoints) <= 0 {
		return fmt.Errorf("[%w] did not specify etcd endpoints", ErrRegistryOption)
	}

	o.EtcdConfig.Logger = o.Logger

	if o.EtcdConfig.DialTimeout == 0 {
		o.EtcdConfig.DialTimeout = defaultTimeout
	}
	if o.EtcdConfig.DialKeepAliveTime == 0 {
		o.EtcdConfig.DialKeepAliveTime = defaultTimeout
	}
	if o.EtcdConfig.DialKeepAliveTimeout == 0 {
		o.EtcdConfig.DialKeepAliveTimeout = defaultTimeout
	}

	return nil
}
