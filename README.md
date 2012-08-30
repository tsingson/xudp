## xudp

XUDP offers an API for extended UDP networking.

The point of this API is to offer persistent, reliable two-way communication
without the overhead imposed by the TCP protocol. This makes it particularly
useful for environments like multiplayer video games, where low latency and is
fast transfer of time sensitive data is of paramount importance.

Features include:

* IPv4 and IPv6 support.
* NAT punch-through: Reliable identification of peers behind the same
  public IP/NAT setup.
* Highly redundant reception acknowledgement by piggybacking multiple
  ACKs on regular data packets.
* Expose event handlers for cases where indivual packets are lost or ACK'ed.
  This allows the host application to implement resending of lost packets.
  Our library does not do this, because this behaviour is the reason why
  TCP is slow. We therefor leave it to the host application to determine
  what to do when packets are lost.


### Usage

    go get github.com/jteeuwen/xudp


### License

Unless otherwise stated, all of the work in this project is subject to a
1-clause BSD license. Its contents can be found in the enclosed LICENSE file.

