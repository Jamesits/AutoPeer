package main

import (
	"github.com/asaskevich/govalidator"
	"net"
	"strings"
)

// (not in use currently)
// Get the actual router ID from user input
// a: can be IP address or a net interface name
func getRouterId(a string) net.IP {
	var routerId net.IP

	if govalidator.IsIPv4(a) {
		routerId = net.ParseIP(a)
	} else {
		// try to find a interface name
		ifaces, _ := net.Interfaces()
		for _, iface := range ifaces {
			found := false

			if strings.EqualFold(iface.Name, a) {
				addrs, _ := iface.Addrs()
				for _, addr := range addrs {
					ip := addr.(*net.IPNet).IP
					if isIPv6(ip) || ip.IsLoopback() {
						// is not a IPv4 or is loopback
						continue
					} else {
						routerId = ip
						found = true
					}
				}
			}

			if found {
				break
			}
		}

		if routerId == nil {
			// TODO: use the smallest IP in all interfaces
			routerId = net.ParseIP("127.0.0.1")
		}
	}

	return routerId
}
