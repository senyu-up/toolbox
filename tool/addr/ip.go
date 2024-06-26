package addr

import (
	"errors"
	"net"
)

var ErrNetWorkNotConnected = errors.New("connected to the network")

// ExternalIP 获取ip
func ExternalIP() (string, error) {
	var (
		ifaces []net.Interface
		iface  net.Interface
		addrs  []net.Addr
		addr   net.Addr
		err    error
	)

	ifaces, err = net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface = range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err = iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr = range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip.String(), nil
		}
	}
	return "", ErrNetWorkNotConnected
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}
