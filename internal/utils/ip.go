package utils

import (
	"errors"
	"fmt"
	"net"
)

var ErrInvalidIP = errors.New("untrusted IP network")

func GetOutboundIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	return conn.LocalAddr().String(), nil
}

func CheckIPBelongsNetwork(incomingIP string, trustedSubnet *net.IPNet) error {
	ip := net.ParseIP(incomingIP)
	if ip == nil || !trustedSubnet.Contains(ip) {
		return fmt.Errorf("%v - %w", ip, ErrInvalidIP)
	}
	return nil
}

func GetTrustedSubnet(cidr string) (*net.IPNet, error) {
	var (
		trustedSubnet *net.IPNet
		err           error
	)

	if cidr != "" {
		_, trustedSubnet, err = net.ParseCIDR(cidr)
		if err != nil {
			return nil, fmt.Errorf("can't parse trusted subnet, reason: %w", err)
		}
	}
	return trustedSubnet, nil
}
