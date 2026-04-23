# netstat collector

The netstat collector exposes metrics about netstat.

## Configuration Flags

| Flag | Description | Default |
| --- | --- | --- |
| collector.netstat.fields | Regexp of fields to return for netstat collector. | ^(.*_(InErrors|InErrs)|Ip_Forwarding|Ip(6|Ext)_(InOctets|OutOctets)|Icmp6?_(InMsgs|OutMsgs)|TcpExt_(Listen.*|Syncookies.*|TCPSynRetrans|TCPTimeouts|TCPOFOQueue|TCPRcvQDrop)|Tcp_(ActiveOpens|InSegs|OutSegs|OutRsts|PassiveOpens|RetransSegs|CurrEstab)|Udp6?_(InDatagrams|OutDatagrams|NoPorts|RcvbufErrors|SndbufErrors))$ |

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_netstat_tcp_receive_packets_total | TCP packets received | n/a |
| node_netstat_tcp_transmit_packets_total | TCP packets sent | n/a |
