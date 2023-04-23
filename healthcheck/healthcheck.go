package healthcheck

import (
	"os"
	"path/filepath"
)

var HealthcheckFile string

func init() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	HealthcheckFile = filepath.Join(dir, "healthcheck.html")
}

func exists(path string) bool {
	_, err := os.Stat(path) //os.Stat 获取文件信息
	if err != nil {
		return false
	}
	return true
}

func IsHealth() bool {
	return exists(HealthcheckFile)
}
