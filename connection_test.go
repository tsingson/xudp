// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package xudp

import (
	"net"
	"testing"
	"time"
)

var Payload = []byte("Hello, world!")

func TestConn(t *testing.T) {
	ca := initConn(t, 12345)
	cb := initConn(t, 12346)

	defer ca.Close()
	defer cb.Close()

	ca.Send(&net.UDPAddr{Port: 12346}, Payload)

	go loop(t, ca)
	go loop(t, cb)

	<-time.After(time.Second / 2)
}

func loop(t *testing.T, c *Connection) {
	for {
		addr, payload, err := c.Recv()

		if err != nil {
			return
		}

		if len(payload) != len(Payload) {
			t.Errorf("Payload size mismatch: Want %d, have %d",
				len(Payload), len(payload))
			return
		}

		for i := range payload {
			if Payload[i] != payload[i] {
				t.Errorf("Payload mismatch at %d: Want %d, have %d",
					i, Payload[i] != payload[i])
				return
			}
		}

		err = c.Send(addr, payload)

		if err != nil {
			return
		}
	}
}

func initConn(t *testing.T, port int) *Connection {
	c := New(1400)
	err := c.Open(port)

	if err != nil {
		t.Fatal(err)
	}

	return c
}
