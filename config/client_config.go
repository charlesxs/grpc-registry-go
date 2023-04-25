package config

type ClientConfig struct {
	ServersDiscovery []*RpcServerConfig `mapstructure:"servers_discovery" json:"servers_discovery" yaml:"servers_discovery"` // 封装了多应用服务发现, gclient针对每个应用创建一个 grpc.ClientConn
	BalancePolicy    string             `mapstructure:"balance_policy" json:"balance_policy" yaml:"balance_policy"`          // 负载均衡策略, 默认为 round_robin
}

type RpcServerConfig struct {
	ServerApp  string      `mapstructure:"server_app" json:"server_app" yaml:"server_app"`    // 指定要连接的server端的app name
	Schema     string      `mapstructure:"schema" json:"schema" yaml:"schema"`                // 指定registry 类型协议
	EtcdConfig *EtcdConfig `mapstructure:"etcd_config" json:"etcd_config" yaml:"etcd_config"` // etcd 相关配置
}

type EtcdConfig struct {
	Endpoints            []string `mapstructure:"endpoints" json:"endpoints" yaml:"endpoints"`                                           // 指定 etcd 服务端地址
	DialTimeout          int64    `mapstructure:"dial_timeout" json:"dial_timeout" yaml:"dial_timeout" `                                 // etcd 链接超时时间, 单位 second, 默认5s
	DialKeepAliveTime    int64    `mapstructure:"dial_keep_alive_time" json:"dial_keep_alive_time" yaml:"dial_keep_alive_time"`          // etcd keepalive 时间, 单位second, 默认5s
	DialKeepAliveTimeout int64    `mapstructure:"dial_keep_alive_timeout" json:"dial_keep_alive_timeout" yaml:"dial_keep_alive_timeout"` // etcd keepalive 超时时间, 单位 second, 默认5s
	Username             string   `mapstructure:"username" json:"username" yaml:"username"`                                              // etcd 认证 username
	Password             string   `mapstructure:"password" json:"password" yaml:"password"`                                              // etcd 认证 password
}
