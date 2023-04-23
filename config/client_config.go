package config

var ClientConfig *clientConfig

const ClientConfigFile = "grpc_client_config.json"

type clientConfig struct {
	ServersDiscovery []*RpcServerConfig `mapstructure:"servers_discovery"`
	BalancePolicy    string             `mapstructure:"balance_policy"`
}

type RpcServerConfig struct {
	ServerAppCode string      `mapstructure:"server_app_code"` // 指定要连接的server端的appCode
	Schema        string      `mapstructure:"schema"`          // 指定registry 类型
	EtcdConfig    *etcdConfig `mapstructure:"etcd_config"`
}

type etcdConfig struct {
	Endpoints            []string `mapstructure:"endpoints"`               // 指定 etcd 服务端地址
	DialTimeout          int64    `mapstructure:"dial_timeout"`            // etcd 链接超时时间, 单位 second, 默认5s
	DialKeepAliveTime    int64    `mapstructure:"dial_keep_alive_time"`    // etcd keepalive 时间, 单位second, 默认5s
	DialKeepAliveTimeout int64    `mapstructure:"dial_keep_alive_timeout"` // etcd keepalive 超时时间, 单位 second, 默认5s
	Username             string   `mapstructure:"username"`                // etcd 认证 username
	Password             string   `mapstructure:"password"`                // etcd 认证 password
}
