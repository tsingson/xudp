// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package xudp

import (
	"testing"
)

func TestSocket(t *testing.T) {
	sock := NewSocket()

	err := sock.Open(12345)

	if err != nil {
		t.Fatal(err)
	}

	defer sock.Close()
}
