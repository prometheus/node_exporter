# ntp collector

The ntp collector exposes metrics about ntp.

## Configuration Flags

| Flag | Description | Default |
| --- | --- | --- |
| collector.ntp.ip-ttl | IP TTL to use while sending NTP query | 1 |
| collector.ntp.local-offset-tolerance | Offset between local clock and local ntpd time to tolerate | 1ms |
| collector.ntp.max-distance | Max accumulated distance to the root | 3.46608s |
| collector.ntp.protocol-version | NTP protocol version | 4 |
| collector.ntp.server | NTP server to use for ntp collector | 127.0.0.1 |
| collector.ntp.server-is-local | Certify that collector.ntp.server address is not a public ntp server | false |
| collector.ntp.server-port | UDP port number to connect to on NTP server | 123 |

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_ntp_leap | NTPD leap second indicator, 2 bits. | n/a |
| node_ntp_offset_seconds | ClockOffset between NTP and local clock. | n/a |
| node_ntp_reference_timestamp_seconds | NTPD ReferenceTime, UNIX timestamp. | n/a |
| node_ntp_root_delay_seconds | NTPD RootDelay. | n/a |
| node_ntp_root_dispersion_seconds | NTPD RootDispersion. | n/a |
| node_ntp_rtt_seconds | RTT to NTPD. | n/a |
| node_ntp_sanity | NTPD sanity according to RFC5905 heuristics and configured limits. | n/a |
| node_ntp_stratum | NTPD stratum. | n/a |
