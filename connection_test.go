// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package xudp

import (
	"testing"
)

const (
	MTU        = 1400
	ProtocolId = 0xBADBEEF
	ServerPort = 12345
	ClientPort = 54321
)

func TestConnection(t *testing.T) {
	var err error

	server := New(MTU, ProtocolId)
	client := New(MTU, ProtocolId)

	if err = server.Open(ServerPort); err != nil {
		t.Fatal(err)
	}

	defer server.Close()

	if err = client.Open(ClientPort); err != nil {
		t.Fatal(err)
	}

	defer client.Close()
}
