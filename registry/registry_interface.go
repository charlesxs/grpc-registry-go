package registry

import (
	"github.com/charlesxs/grpc-registry-go/config"
)

// IRegistry 接口
type IRegistry interface {
	// Register 以应用维度注册一个 endpoint
	Register(addr string, port int, metadata interface{}) error

	// Deregister 以应用维度取消一个endpoint的注册记录
	Deregister(addr string, port int) error
}

// IRegistryFactory registry工厂，用于创建不同的registry
type IRegistryFactory interface {
	// BuildOptions 构建 registry 配置选项
	BuildOptions(cfg *config.ServerConfig) error

	// CreateRegistry 创建 registry
	CreateRegistry() (IRegistry, error)
}
