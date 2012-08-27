// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package xudp

import (
	"fmt"
	"net"
)

// findLocalAddr finds a valid and suitable local address.
// It returns the address and whether or not it is an IPv6 address.
func findLocalAddr() (net.IP, bool) {
	var ip net.IP
	var isIPv6 bool

	iface, err := findLocalInterface()

	if err != nil {
		return nil, false
	}

	addrlist, _ := iface.Addrs()

	for _, addr := range addrlist {
		ipnet, ok := addr.(*net.IPNet)

		if !ok {
			continue
		}

		isIPv6, ip = false, ipnet.IP.To4()

		if len(ip) == 0 {
			isIPv6, ip = true, ipnet.IP.To16()
		}

		break
	}

	return ip, isIPv6
}

// findLocalInterface finds a network interface suitable for data transport.
func findLocalInterface() (*net.Interface, error) {
	list, err := net.Interfaces()

	if err != nil {
		return nil, fmt.Errorf("findLocalInterface: %v", err)
	}

	for i := range list {
		iface := list[i]

		if iface.Flags&net.FlagUp != 0 &&
			iface.Flags&net.FlagBroadcast != 0 &&
			iface.Flags&net.FlagMulticast != 0 {

			addr, err := iface.Addrs()

			if err != nil {
				continue
			}

			if len(addr) > 0 {
				return &iface, nil
			}
		}
	}

	return nil, fmt.Errorf("findLocalInterface: No valid interface could be found.")
}
