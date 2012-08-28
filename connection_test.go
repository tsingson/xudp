// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package xudp

import (
	"io"
	"log"
	"net"
	"testing"
)

const (
	MTU        = 1400
	ProtocolId = 0xBADBEEF
	ServerPort = 12345
	ClientPort = 54321
)

var serverAddr = &net.UDPAddr{
	IP:   net.ParseIP("[::1]"),
	Port: ServerPort,
}

func TestConnection(t *testing.T) {
	var err error

	server := NewConnection(MTU, ProtocolId)
	if err = server.Open(ServerPort); err != nil {
		t.Errorf("server.Open: %v", err)
		return
	}

	defer server.Close()

	client := NewConnection(MTU, ProtocolId)
	if err = client.Open(ClientPort); err != nil {
		t.Errorf("client.Open: %v", err)
		return
	}

	defer client.Close()

	go echo(t, server)
	go echo(t, client)

	client.Send(serverAddr, []byte("Hello, server."))
}

func echo(t *testing.T, c *Connection) {
	var addr *net.UDPAddr
	var data []byte
	var err error

	for {
		addr, data, err = c.Recv()

		if err != nil {
			if err == io.EOF {
				continue
			}

			t.Errorf("Recv: %v", err)
			return
		}

		log.Printf("%s: %v", addr, data)

		err = c.Send(addr, data)

		if err != nil {
			t.Errorf("Send: %v", err)
			return
		}
	}
}
