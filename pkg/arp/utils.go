package arp

import (
	"net"
	"sync"
)

//  http://play.golang.org/p/m8TNTtygK0
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func hosts(cidr string) ([]net.IP, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	if x := ip.To4(); x == nil {
		return nil, nil
	}

	var ips []net.IP
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		newIP := make([]byte, net.IPv4len)
		copy(newIP, ip)
		ips = append(ips, newIP)
	}
	// remove network address and broadcast address
	return ips[1 : len(ips)-1], nil
}

type discoveryTable struct {
	sync.Mutex
	discovered map[string]struct{}
}
