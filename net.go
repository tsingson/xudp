// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package xudp

import (
	"net"
)

// LocalIP returns the first available local IP address.
// This is the subnet address if the host is located in a subnet.
func LocalIP() net.IP {
	// Connect to a random machine somewhere. It's irrelevant
	// where to, as long as it's not the loopback address.

	addr, err := net.ResolveUDPAddr("udp", "192.168.0.0:0")

	if err != nil {
		return nil
	}

	conn, err := net.DialUDP("udp", nil, addr)

	if err != nil {
		return nil
	}

	addr = conn.LocalAddr().(*net.UDPAddr)
	conn.Close()

	return addr.IP
}
