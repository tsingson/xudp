## Order

This sample program demonstrates how to use the `xudp.Connection`, along
with the Reliability plugin to do buffered, in-order reception of data.

This caches outgoing packets. It tracks which packets have been lost along
the way and resends them.

The receiving end returns packets in sequential order, as defined by
the packet sequence number.

