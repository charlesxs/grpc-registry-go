package registry

import "errors"

var ErrRegistry = errors.New("注册中心错误")
var ErrRegistryOption = errors.New("注册中心配置错误")
var ErrConfigNotFound = errors.New("未找到配置")
