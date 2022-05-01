package utils

import (
	"github.com/pkg/errors"
	"net"
)

// GetFreePort 获取本机未使用的端口号
func GetFreePort() (port int, err error) {
	var (
		addr     *net.TCPAddr
		listener net.Listener
	)
	if addr, err = net.ResolveTCPAddr("tcp", "localhost:0"); err != nil {
		return
	}
	if listener, err = net.ListenTCP("tcp", addr); err != nil {
		return
	}
	defer listener.Close()
	port = listener.Addr().(*net.TCPAddr).Port
	return
}

// GetLocalIp 取本机网卡 IP
func GetLocalIp() (ipv4 string, err error) {
	var (
		addrs   []net.Addr
		addr    net.Addr
		ipNet   *net.IPNet
		isIpNet bool
	)
	if addrs, err = net.InterfaceAddrs(); err != nil {
		return
	}
	// 取第一个非 IO 网卡 IP
	for _, addr = range addrs {
		if ipNet, isIpNet = addr.(*net.IPNet); isIpNet && ipNet.IP.IsLoopback() { // 跳过非 IP 和回路IP
			if ipNet.IP.To4() != nil { // 跳过 IPV6
				ipv4 = ipNet.IP.String()
				return
			}
		}
	}
	err = errors.New("no local IP address")
	return
}
