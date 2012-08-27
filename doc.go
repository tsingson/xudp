// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

/*
XUDP offers an API for extended UDP networking.

The point of this API is to offer persistent, reliable* two-way communication
without the overhead imposed by the TCP protocol. This makes it particularly
useful for high performance environments like multiplayer video games.

*) By reliable, we mean a different kind of reliability from TCP by
guaranteeing high throughput for time sensitive data.

*/
package xudp
