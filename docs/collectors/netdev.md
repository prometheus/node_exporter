# netdev collector

The netdev collector exposes metrics about netdev.

## Configuration Flags

| Flag | Description | Default |
| --- | --- | --- |
| collector.netdev.address-info | Collect address-info for every device |  |
| collector.netdev.device-blacklist | DEPRECATED: Use collector.netdev.device-exclude |  |
| collector.netdev.device-exclude | Regexp of net devices to exclude (mutually exclusive to device-include). |  |
| collector.netdev.device-include | Regexp of net devices to include (mutually exclusive to device-exclude). |  |
| collector.netdev.device-whitelist | DEPRECATED: Use collector.netdev.device-include |  |
| collector.netdev.enable-detailed-metrics | Use (incompatible) metric names that provide more detailed stats on Linux |  |
| collector.netdev.label-ifalias | Add ifAlias label | false |
| collector.netdev.netlink | Use netlink to gather stats instead of /proc/net/dev. | true |

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_network_address_info | node network address by device | device, address, netmask, scope |
