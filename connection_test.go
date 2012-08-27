// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package xudp

import (
	"testing"
)

const ProtocolId = 'X'<<24 | 'U'<<16 | 'D'<<8 | 'P'

func TestServer(t *testing.T) {
	conn := NewConnection(ModeServer, ProtocolId)
	defer conn.Close()

}
