package netex

import (
	"errors"
	"net"
)

type InterfaceAddrs func() ([]net.Addr, error)

func PrivateIPv4(interfaceAddrs InterfaceAddrs) (net.IP, error) {
	as, err := interfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, a := range as {
		ipnet, ok := a.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}

		ip := ipnet.IP.To4()
		if IsPrivateIPv4(ip) {
			return ip, nil
		}
	}
	return nil, errors.New("no private ip address")
}

func IsPrivateIPv4(ip net.IP) bool {
	// Allow private IP addresses (RFC1918) and link-local addresses (RFC3927)
	return ip != nil &&
		(ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) || ip[0] == 192 && ip[1] == 168 || ip[0] == 169 && ip[1] == 254)
}

func Lower16BitPrivateIP(interfaceAddrs InterfaceAddrs) (int, error) {
	ip, err := PrivateIPv4(interfaceAddrs)
	if err != nil {
		return 0, err
	}

	return int(ip[2])<<8 + int(ip[3]), nil
}
