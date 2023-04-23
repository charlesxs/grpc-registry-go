package registry

// Factories 指定所有支持的注册中心的工厂
var Factories map[Schema]IRegistryFactory

func init() {
	Factories = map[Schema]IRegistryFactory{
		SchemaEtcd: &etcdRegistryFactory{},
	}
}
