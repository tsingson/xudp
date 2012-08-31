// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package ident

import (
	"net"
	"testing"
)

func TestPeerIdEquality(t *testing.T) {
	tests := []struct {
		IPA   string
		PortA int
		IPB   string
		PortB int
		Equal bool
	}{
		{"1.2.3.4", 12345, "1.2.3.4", 12345, true},
		{"::1", 12345, "::1", 12345, true},
		{"1.2.3.5", 12345, "1.2.3.4", 12345, false},
		{"1.2.3.4", 12346, "1.2.3.4", 12345, false},
		{"::1", 12346, "::1", 12345, false},
		{"", 0, "", 0, true},
	}

	for _, test := range tests {
		a := NewPeerId(net.ParseIP(test.IPA), test.PortA)
		b := NewPeerId(net.ParseIP(test.IPB), test.PortB)
		eq := compareHash(a, b)

		if eq != test.Equal {
			t.Fatalf("Want %v, have %v\n%s:%d\n%s:%d",
				test.Equal, eq, test.IPA, test.PortA, test.IPB, test.PortB)
		}
	}
}
