package gserver

import (
	"net"
)

// LocalAddrs 获取本机所有非回环地址
func LocalAddrs() ([]string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	var result []string
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ip, ok := address.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				result = append(result, ip.IP.String())
			}
		}
	}
	return result, nil
}
