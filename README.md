## xudp

**Note**: This is work in progress and needs more testing.

XUDP offers an API for extended UDP networking.

The point of this API is to offer persistent, reliable two-way communication
without the overhead imposed by the TCP protocol. This makes it particularly
useful for environments like multiplayer video games, where low latency and
fast transfer of time sensitive data is of paramount importance.

What we support:

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
  We therefore leave it to the host application to determine
  what to do when packets are lost. It can resend packets selectively
  while not preventing the reception of remaining data.


What do we not support:

* This package explicitely does **not** guarantee in-order reception
  of packets, for the same reason described in the feature point on
  event handlers.
* Packet fragmentation, encryption or compression. These are all high
  level abstractions that are best left to the host application, because
  networking requirements are very different from one case to the next.

Much of the code in this package is ported from the guides published [here][1]

[1]: http://gafferongames.com/networking-for-game-programmers/udp-vs-tcp/


### Extensibility

The package is deliberately kept low level. It should serve as a foundation
for more appropriate networking code in your own program. You can create
your own higher level abstraction which embeds xudp.Connection and does any
additional handling you may require. For an example of this, refer to
the ident/connection.go file.


### Usage

    go get github.com/jteeuwen/xudp


### License

Unless otherwise stated, all of the work in this project is subject to a
1-clause BSD license. Its contents can be found in the enclosed LICENSE file.

