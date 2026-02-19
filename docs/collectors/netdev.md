# netdev

Exposes network interface statistics such as bytes transferred, packets, errors, and drops.

Status: enabled by default

## Platforms

- Linux
- Darwin
- FreeBSD
- OpenBSD
- Dragonfly
- AIX

## Configuration

```
--collector.netdev.device-include          Regexp of devices to include (mutually exclusive to device-exclude)
--collector.netdev.device-exclude          Regexp of devices to exclude (mutually exclusive to device-include)
--collector.netdev.address-info            Collect address info for every device (default: false)
--collector.netdev.enable-detailed-metrics Use detailed metric names on Linux (default: false)
```

### Examples

Exclude virtual and container interfaces:
```
--collector.netdev.device-exclude="^(veth|docker|br-|virbr|cni|flannel|cali).*"
```

Monitor only physical ethernet interfaces:
```
--collector.netdev.device-include="^(eth|ens|enp|eno)[0-9]+"
```

Exclude loopback only:
```
--collector.netdev.device-exclude="^lo$"
```

Include bonded interfaces and their members:
```
--collector.netdev.device-include="^(bond[0-9]+|eth[0-9]+)$"
```

Enable IP address information for all interfaces:
```
--collector.netdev.address-info
```

Use detailed error metrics (breaks compatibility with default metric names):
```
--collector.netdev.enable-detailed-metrics
```

## Data Sources

| Source | Description |
|--------|-------------|
| `/proc/net/dev` | Network device statistics (Linux) |
| `/sys/class/net/` | Network device info (Linux) |
| `getifaddrs(3)` | Interface addresses (all platforms) |

Documentation:
- https://docs.kernel.org/networking/statistics.html
- `netdevice(7)` manpage

## Metrics

All metrics have the `device` label and are counters with `_total` suffix.

### Standard Metrics (default)

| Metric | Type | Description |
|--------|------|-------------|
| `node_network_receive_bytes_total` | counter | Bytes received |
| `node_network_receive_packets_total` | counter | Packets received |
| `node_network_receive_errs_total` | counter | Receive errors |
| `node_network_receive_drop_total` | counter | Packets dropped on receive |
| `node_network_receive_fifo_total` | counter | FIFO buffer errors on receive |
| `node_network_receive_frame_total` | counter | Frame errors on receive |
| `node_network_receive_compressed_total` | counter | Compressed packets received |
| `node_network_receive_multicast_total` | counter | Multicast packets received |
| `node_network_transmit_bytes_total` | counter | Bytes transmitted |
| `node_network_transmit_packets_total` | counter | Packets transmitted |
| `node_network_transmit_errs_total` | counter | Transmit errors |
| `node_network_transmit_drop_total` | counter | Packets dropped on transmit |
| `node_network_transmit_fifo_total` | counter | FIFO buffer errors on transmit |
| `node_network_transmit_colls_total` | counter | Collisions detected |
| `node_network_transmit_carrier_total` | counter | Carrier errors on transmit |
| `node_network_transmit_compressed_total` | counter | Compressed packets transmitted |

### Detailed Metrics (--collector.netdev.enable-detailed-metrics)

When enabled, exposes more granular error counters instead of aggregated values:

| Metric | Type | Description |
|--------|------|-------------|
| `node_network_receive_errors_total` | counter | Total receive errors |
| `node_network_receive_dropped_total` | counter | Dropped packets (excludes missed) |
| `node_network_receive_missed_errors_total` | counter | Missed packets |
| `node_network_receive_fifo_errors_total` | counter | FIFO overrun errors |
| `node_network_receive_length_errors_total` | counter | Length errors |
| `node_network_receive_over_errors_total` | counter | Ring buffer overflow |
| `node_network_receive_crc_errors_total` | counter | CRC errors |
| `node_network_receive_frame_errors_total` | counter | Frame alignment errors |
| `node_network_transmit_errors_total` | counter | Total transmit errors |
| `node_network_transmit_dropped_total` | counter | Dropped packets |
| `node_network_transmit_fifo_errors_total` | counter | FIFO errors |
| `node_network_transmit_aborted_errors_total` | counter | Aborted transmissions |
| `node_network_transmit_carrier_errors_total` | counter | Carrier errors |
| `node_network_transmit_heartbeat_errors_total` | counter | Heartbeat errors |
| `node_network_transmit_window_errors_total` | counter | Window errors |
| `node_network_collisions_total` | counter | Collision count |

### Address Info (--collector.netdev.address-info)

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `node_network_address_info` | gauge | `device`, `address`, `netmask`, `scope` | Network address info (always 1) |

## Labels

| Label | Description |
|-------|-------------|
| `device` | Interface name (e.g., `eth0`, `ens192`, `lo`) |
| `address` | IP address (address_info only) |
| `netmask` | CIDR prefix length (address_info only) |
| `scope` | Address scope: `global`, `link-local`, `interface-local` (address_info only) |

## Notes

- Default metrics match `/proc/net/dev` column names for compatibility
- Detailed metrics provide per-error-type breakdown but change metric names
- Virtual interfaces (veth, docker, etc.) are included by default; use `--collector.netdev.device-exclude` to filter
- Loopback (`lo`) is included by default
