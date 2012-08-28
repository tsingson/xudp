// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package xudp

import (
	"testing"
)

const (
	MTU        = 1400
	ProtocolId = 0x12345678
	Port       = 12345
)

func TestServer(t *testing.T) {
	conn := New(Server, MTU, ProtocolId)
	err := conn.Open(Port)

	if err != nil {
		t.Fatal(err)
	}

	defer conn.Close()
}
