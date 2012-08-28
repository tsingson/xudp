## xudp

XUDP offers an API for extended UDP networking.

The point of this API is to offer persistent, reliable* two-way communication
without the overhead imposed by the TCP protocol. This makes it particularly
useful for high performance environments like multiplayer video games.

*) By reliable, we mean a different kind of reliability from TCP by
guaranteeing high throughput for time sensitive data.

Features include:

* IPv4 and IPv6 support.
* NAT punch-through: Reliable identification of peers behind the same
  public IP/NAT setup.
* Packet fragmentation: Sending and receiving of data that spans
  multiple packets. The data is guaranteed to be received and reassembled
  in the correct order.
* Highly redundant reception acknowledgement by piggybacking multiple
  ACKs on regular data packets.
* A single connection type which functions as a server and client at
  the same time. It hides the internals required to make this work.
  To the host, we are simply working with a two-way connection on which
  we can read/write data at will. 


### Usage

    go get github.com/jteeuwen/xudp


### License

Unless otherwise stated, all of the work in this project is subject to a
1-clause BSD license. Its contents can be found in the enclosed LICENSE file.

