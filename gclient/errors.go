package gclient

import "errors"

var ErrClientInit = errors.New("grpcClient初始化错误")
var ErrConfigNotFound = errors.New("未找到配置")
var ErrConfig = errors.New("配置错误")
var ErrCreateConn = errors.New("创建链接失败")
