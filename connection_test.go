// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package xudp

import (
	"testing"
)

func TestServer(t *testing.T) {
	conn := NewConnection(ModeServer, MTU, ProtocolId)
	defer conn.Close()

}
