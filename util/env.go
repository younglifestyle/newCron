package util

import (
	"errors"
	"fmt"
	"net"
	"os"
)

var (
	// AppIP agent ip
	appIP = GetIP()
	// HostName machine hostname
	hostName = GetHostName("")
)

// GetHostName ...
func GetHostName(hostName string) string {
	if host := os.Getenv(hostName); host != "" {
		return host
	}
	name, err := os.Hostname()
	if err != nil {
		return fmt.Sprintf("error:%s", err.Error())
	}
	return name
}

// GetIP ...
func GetIP() string {
	ip, err := IP()
	if err != nil {
		return fmt.Sprintf("ip.error:%s", err.Error())
	}
	return ip
}

func IP() (string, error) {
	ip, err := externalIP()
	if err != nil {
		return "", err
	}
	return ip.String(), nil
}

func externalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addres, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addres {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network error")
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
		return nil
	}
	return ip
}


func ReturnHostName() string {
	return hostName
}

func ReturnAppIp() string {
	return appIP
}
