package healthcheck

import (
	"os"
)

// IHealthChecker 实现健康检测的接口
type IHealthChecker interface {
	// IsHealth 可用于做业务检查，其中 true 表示健康, false表示不健康
	IsHealth() bool
}

type fileHealthChecker struct {
	fpath string
}

func NewFileHealth(fpath string) *fileHealthChecker {
	return &fileHealthChecker{
		fpath: fpath,
	}
}

func (fh *fileHealthChecker) IsHealth() bool {
	_, err := os.Stat(fh.fpath)
	if err != nil {
		return false
	}
	return true
}
