package healthcheck

import (
	"os"
)

// IHealthChecker 实现健康检测的接口
type IHealthChecker interface {
	// IsHealth 可用于做业务检查，其中 true 表示健康, false表示不健康
	IsHealth() bool
}

// 基于文件的健康检测，当某个文件存在时便认为服务Ok,将自己注册到注册中心，当文件不存在时便将自己从注册中心债除掉
// 测试请看 checker_test.go
type fileHealthChecker struct {
	filePath string
}

func NewFileHealth(filePath string) IHealthChecker {
	return &fileHealthChecker{
		filePath: filePath,
	}
}

func (fh *fileHealthChecker) IsHealth() bool {
	_, err := os.Stat(fh.filePath)
	if err != nil {
		return false
	}
	return true
}
