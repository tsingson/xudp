// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

/*
XUDP offers an API for extended UDP networking.

The point of this API is to offer persistent, reliable two-way communication
without the overhead imposed by the TCP protocol. This makes it particularly
useful for environments like multiplayer video games, where low latency and is
fast transfer of time sensitive data is of paramount importance.

Features include:

	* IPv4 and IPv6 support.
	* Highly redundant reception acknowledgement by piggybacking multiple
	  ACKs on regular data packets.
	* Exposes event handlers for cases where indivual packets are lost or ACK'ed.
	  This allows the host application to implement resending of lost packets.
	  Our library does not do this for performance reasons. To be more
	  precise: when TCP detects a packet loss, it stops the sending of
	  everything else until the lost packet has been re-sent and ACK'ed by
	  the other end. For applications where time-sensitive data should
	  go through as fast as possible, this is very much not what we want.
	  We therefor leave it to the host application to determine
	  what to do when packets are lost. It can resend packets selectively
	  while not preventing the reception of remaining data.
	* This package explicitely does **not** guarantee in-order reception
	  of packets, for the same reason described in the previous point.
*/
package xudp
