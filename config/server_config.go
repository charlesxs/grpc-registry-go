package config

import "go.uber.org/zap"

type ServerConfig struct {
	// 公共配置属性
	AppName             string              `mapstructure:"app_name" json:"app_name" yaml:"app_name"`                                     // 服务app name
	Port                int                 `mapstructure:"port" json:"port" yaml:"port"`                                                 // 服务监听的port
	Schema              string              `mapstructure:"schema" json:"schema" yaml:"schema"`                                           // 指定registry 类型
	HealthcheckInterval int                 `mapstructure:"healthcheck_interval" json:"healthcheck_interval" yaml:"healthcheck_interval"` // 指定健康检测的间隔, 单位second
	EtcdRegistryConfig  *EtcdRegistryConfig `mapstructure:"etcd_registry_config" json:"etcd_registry_config" yaml:"etcd_registry_config"` // 指定etc registry的配置

	// 私有化属性
	Logger *zap.Logger `mapstructure:"-" json:"-" yaml:"-"`
}

type EtcdRegistryConfig struct {
	Endpoints            []string `mapstructure:"endpoints" json:"endpoints" yaml:"endpoints"`                                           // 指定 etcd 服务端地址
	DialTimeout          int64    `mapstructure:"dial_timeout" json:"dial_timeout" yaml:"dial_timeout"`                                  // etcd 链接超时时间, 单位 second, 默认5s
	DialKeepAliveTime    int64    `mapstructure:"dial_keep_alive_time" json:"dial_keep_alive_time" yaml:"dial_keep_alive_time"`          // etcd keepalive 时间, 单位second, 默认5s
	DialKeepAliveTimeout int64    `mapstructure:"dial_keep_alive_timeout" json:"dial_keep_alive_timeout" yaml:"dial_keep_alive_timeout"` // etcd keepalive 超时时间, 单位 second, 默认5s
	Username             string   `mapstructure:"username" json:"username" yaml:"username"`                                              // etcd 认证 username
	Password             string   `mapstructure:"password" json:"password" yaml:"password"`                                              // etcd 认证 password
	Interval             int64    `mapstructure:"interval" json:"interval" yaml:"interval"`                                              // 注册中心配置, 指定多久检测并重新续约, 默认10s
	LeaseTTL             int64    `mapstructure:"lease_ttl" json:"lease_ttl" yaml:"lease_ttl"`                                           // 注册中心配置, 租约时长, 默认12s
}
