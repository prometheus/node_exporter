# netstat

Exposes network statistics from `/proc/net/netstat`, `/proc/net/snmp`, and `/proc/net/snmp6`.

Status: enabled by default

## Platforms

- Linux

## Configuration

```
--collector.netstat.fields  Regexp of fields to return (default: see below)
```

Default pattern:
```
^(.*_(InErrors|InErrs)|Ip_Forwarding|Ip(6|Ext)_(InOctets|OutOctets)|Icmp6?_(InMsgs|OutMsgs)|TcpExt_(Listen.*|Syncookies.*|TCPSynRetrans|TCPTimeouts|TCPOFOQueue|TCPRcvQDrop)|Tcp_(ActiveOpens|InSegs|OutSegs|OutRsts|PassiveOpens|RetransSegs|CurrEstab)|Udp6?_(InDatagrams|OutDatagrams|NoPorts|RcvbufErrors|SndbufErrors))$
```

### Examples

Expose all available metrics:
```
--collector.netstat.fields=".*"
```

TCP metrics only (basic and extended):
```
--collector.netstat.fields="^Tcp(Ext)?_.*"
```

Only error-related metrics:
```
--collector.netstat.fields=".*_(InErrors|InErrs|Drops|Timeouts|Retrans).*"
```

Minimal set (bytes in/out and established connections):
```
--collector.netstat.fields="^(IpExt_(InOctets|OutOctets)|Tcp_CurrEstab)$"
```

Add memory pressure metrics to the default set:
```
--collector.netstat.fields="^(.*_(InErrors|InErrs)|Ip_Forwarding|Ip(6|Ext)_(InOctets|OutOctets)|Icmp6?_(InMsgs|OutMsgs)|TcpExt_(Listen.*|Syncookies.*|TCPSynRetrans|TCPTimeouts|TCPOFOQueue|TCPRcvQDrop|TCPMemoryPressures.*)|Tcp_(ActiveOpens|InSegs|OutSegs|OutRsts|PassiveOpens|RetransSegs|CurrEstab)|Udp6?_(InDatagrams|OutDatagrams|NoPorts|RcvbufErrors|SndbufErrors))$"
```

## Data Sources

| Source | Description |
|--------|-------------|
| `/proc/net/netstat` | Extended TCP statistics (TcpExt, IpExt) |
| `/proc/net/snmp` | SNMP MIB statistics (Ip, Icmp, Tcp, Udp) |
| `/proc/net/snmp6` | IPv6 SNMP statistics (Ip6, Icmp6, Udp6) |

Documentation:
- https://docs.kernel.org/networking/snmp_counter.html
- https://docs.kernel.org/filesystems/proc.html (Table 1-9: Network info in /proc/net)
- `netstat(8)` manpage (`netstat -s` displays the same statistics)

## Metrics

Metrics are dynamically generated as `node_netstat_<Protocol>_<Field>` based on the fields regex filter.

### Default Exposed Metrics

#### IP

| Metric | Type | Description |
|--------|------|-------------|
| `node_netstat_Ip_Forwarding` | untyped | IP forwarding status |
| `node_netstat_IpExt_InOctets` | untyped | Total incoming bytes |
| `node_netstat_IpExt_OutOctets` | untyped | Total outgoing bytes |
| `node_netstat_Ip6_InOctets` | untyped | Total incoming IPv6 bytes |
| `node_netstat_Ip6_OutOctets` | untyped | Total outgoing IPv6 bytes |

#### ICMP

| Metric | Type | Description |
|--------|------|-------------|
| `node_netstat_Icmp_InMsgs` | untyped | Total incoming ICMP messages |
| `node_netstat_Icmp_OutMsgs` | untyped | Total outgoing ICMP messages |
| `node_netstat_Icmp6_InMsgs` | untyped | Total incoming ICMPv6 messages |
| `node_netstat_Icmp6_OutMsgs` | untyped | Total outgoing ICMPv6 messages |

#### TCP

| Metric | Type | Description |
|--------|------|-------------|
| `node_netstat_Tcp_ActiveOpens` | untyped | Active connection openings |
| `node_netstat_Tcp_PassiveOpens` | untyped | Passive connection openings |
| `node_netstat_Tcp_InSegs` | untyped | Incoming segments |
| `node_netstat_Tcp_OutSegs` | untyped | Outgoing segments |
| `node_netstat_Tcp_RetransSegs` | untyped | Retransmitted segments |
| `node_netstat_Tcp_OutRsts` | untyped | Outgoing resets |
| `node_netstat_Tcp_CurrEstab` | untyped | Currently established connections |

#### TCP Extended

| Metric | Type | Description |
|--------|------|-------------|
| `node_netstat_TcpExt_ListenOverflows` | untyped | Listen queue overflows |
| `node_netstat_TcpExt_ListenDrops` | untyped | Dropped incoming connections |
| `node_netstat_TcpExt_SyncookiesSent` | untyped | SYN cookies sent |
| `node_netstat_TcpExt_SyncookiesRecv` | untyped | SYN cookies received |
| `node_netstat_TcpExt_SyncookiesFailed` | untyped | SYN cookies failed |
| `node_netstat_TcpExt_TCPSynRetrans` | untyped | SYN retransmissions |
| `node_netstat_TcpExt_TCPTimeouts` | untyped | TCP timeouts |
| `node_netstat_TcpExt_TCPOFOQueue` | untyped | Out-of-order queue usage |
| `node_netstat_TcpExt_TCPRcvQDrop` | untyped | Receive queue drops |

#### UDP

| Metric | Type | Description |
|--------|------|-------------|
| `node_netstat_Udp_InDatagrams` | untyped | Incoming UDP datagrams |
| `node_netstat_Udp_OutDatagrams` | untyped | Outgoing UDP datagrams |
| `node_netstat_Udp_NoPorts` | untyped | Datagrams to unknown ports |
| `node_netstat_Udp_RcvbufErrors` | untyped | Receive buffer errors |
| `node_netstat_Udp_SndbufErrors` | untyped | Send buffer errors |
| `node_netstat_Udp6_InDatagrams` | untyped | Incoming UDPv6 datagrams |
| `node_netstat_Udp6_OutDatagrams` | untyped | Outgoing UDPv6 datagrams |

#### Error Metrics

Any field matching `.*_(InErrors|InErrs)` is included by default.

## Notes

- All metrics are exposed as `untyped` since the kernel doesn't indicate whether values are counters or gauges
- `/proc/net/snmp6` may not exist on systems with IPv6 disabled
- Customize `--collector.netstat.fields` to expose additional or fewer metrics
- Field names match the kernel's naming convention exactly
