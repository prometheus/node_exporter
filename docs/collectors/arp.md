# arp collector

The arp collector exposes metrics about arp.

## Configuration Flags

| Flag | Description | Default |
| --- | --- | --- |
| collector.arp.device-exclude | Regexp of arp devices to exclude (mutually exclusive to device-include). |  |
| collector.arp.device-include | Regexp of arp devices to include (mutually exclusive to device-exclude). |  |
| collector.arp.netlink | Use netlink to gather stats instead of /proc/net/arp. | true |

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_arp_entries | ARP entries by device | device |
