// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package ident

import (
	"crypto/sha256"
	"net"
)

// Size of a single Peer id in bytes.
const PeerIdSize = 32

// A PeerId is a SHA256 hash of the peer's internal IP address and port.
// This is included in every packet from the given peer.
type PeerId []byte

// NewPeerId generates an SHA256 hash for the given IP and port.
func NewPeerId(ip net.IP, port int) PeerId {
	data := make([]byte, len(ip)+4)

	data[0] = byte(port >> 24)
	data[1] = byte(port >> 16)
	data[2] = byte(port >> 8)
	data[3] = byte(port)
	copy(data[4:], ip)

	hm := sha256.New()
	hm.Write(data)
	return PeerId(hm.Sum(nil))
}

// compareHash compares two hashes and returns true if we consider them equal.
//
// A constant time comparison is used to prevent timing attacks from
// being performed. With a normal bytes.Equal(a, b) comparison, an attacker can 
// time how long this function takes to complete. The longer it takes
// to return, the more of the hash he knows will be correct. A constant time
// comparison always runs in the same amount of time, regardless of the
// hash contents; thus eliminating the timing attack vector.
func compareHash(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	var v byte

	for i := 0; i < len(a); i++ {
		v |= a[i] ^ b[i]
	}

	v = ^v
	v &= v >> 4
	v &= v >> 2
	v &= v >> 1
	return v == 1
}
