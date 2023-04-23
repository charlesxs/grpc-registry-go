package config

var ServerConfig *serverConfig

const ServerConfigFile = "grpc_server_config.json"

type serverConfig struct {
	AppCode             string              `mapstructure:"app_code"`             // 服务appCode
	Port                int                 `mapstructure:"port"`                 // 服务监听的port
	Schema              string              `mapstructure:"schema"`               // 指定registry 类型
	HealthcheckInterval int                 `mapstructure:"healthcheck_interval"` // 指定健康检测的间隔, 单位second
	EtcdRegistryConfig  *etcdRegistryConfig `mapstructure:"etcd_registry_config"` // 指定etc registry的配置
}

type etcdRegistryConfig struct {
	Endpoints            []string `mapstructure:"endpoints"`               // 指定 etcd 服务端地址
	DialTimeout          int64    `mapstructure:"dial_timeout"`            // etcd 链接超时时间, 单位 second, 默认5s
	DialKeepAliveTime    int64    `mapstructure:"dial_keep_alive_time"`    // etcd keepalive 时间, 单位second, 默认5s
	DialKeepAliveTimeout int64    `mapstructure:"dial_keep_alive_timeout"` // etcd keepalive 超时时间, 单位 second, 默认5s
	Username             string   `mapstructure:"username"`                // etcd 认证 username
	Password             string   `mapstructure:"password"`                // etcd 认证 password
	Interval             int64    `mapstructure:"interval"`                // 注册中心配置, 指定多久检测并重新续约
	LeaseTTL             int64    `mapstructure:"lease_ttl"`               // 注册中心配置, 租约时长
}
