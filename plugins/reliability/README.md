## Reliability

Reliability implements the algorithms needed to make a reliable connection
reliable. This means it manages sent, received, pending ACKs and ACK'ed
packet queues. In addtion, it tracks bandwidth use and round trip timing.

It exposes event handlers which notify the host of any specific packets
which have been either ACK'ed or lost. The host application can then
chose to take whatever action is necessary. For any lost packets, it may
chose to resend the lost payload if necessary.

It achieves all this by adding three 32-bit integer fields to the packet
header. The first one is a numerical Sequence value, which identifies the
specific packet.

The second field is the ACK field. It holds the sequence number of a
packet we have previously received and are acknowledging to the other peer.

The third field is an ACK vector. Combined with the ACK field, this allows
us to piggyback up to 33 ACKS simultaneously in a single data packet.
Even when a number of packets are lost, this creates a highly redundant
packet acknowledgement mechanism.

To illustrate:

* We have received packet with sequence #100.
* We set the ACK field of our next outgoing packet to #100 to
  confirm its reception.
* Any older packets (#68 - #99) which have been received, are then
  encoded in the ACK vector. For packet #99, we set bit 1 of the vector.
  For packet #98, we set bit 2, etc. Any packets in this range which we
  have not received, keep their bit index at 0.
  This allows the other end to see at a glance which of the last 33
  packets should be marked as lost or not.

This plugin does not resend lost packets itself.
We simply offer a way for the host application to know about lost packets
and leave it up to them to decide what to do. The reason for this is
performance. The TCP protocol is pretty slow for some applications because
it detects lost packets, stops sending new ones and keeps resends the lost one
until it has a reception confirmation. This creates a pile up of packet data
which is very much undesirable in high performance, low-latency environments
like video games. And also the reason why many video games use UDP instead
of TCP.

The way this plugin works, should not be confused with guaranteed in-order
reception of packet data by either end of the connection. This is not what
we do. For the same reasons as stated in the previous paragraph, this is
left to the host application.


### Usage

    go get github.com/jteeuwen/xudp/plugins/reliability


### License

Unless otherwise stated, all of the work in this project is subject to a
1-clause BSD license. Its contents can be found in the enclosed LICENSE file.



