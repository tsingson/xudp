// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package xudp

import (
	"fmt"
	"testing"
)

const (
	MTU        = 1400
	ProtocolId = 'X'<<24 | 'U'<<16 | 'D'<<8 | 'P'
)

func TestSocket(t *testing.T) {
	sock := newSocket(MTU, ProtocolId)

	err := sock.Open(12345)

	if err != nil {
		t.Fatal(err)
	}

	defer sock.Close()

	recv := sock.Poll()

	for {
		select {
		case packet := <-recv:
			if packet == nil {
				return
			}

			fmt.Printf("%v\n", packet)
		}
	}
}
