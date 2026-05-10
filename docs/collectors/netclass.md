# netclass collector

The netclass collector exposes metrics about netclass.

## Configuration Flags

| Flag | Description | Default |
| --- | --- | --- |
| collector.netclass.ignore-invalid-speed | Ignore devices where the speed is invalid. This will be the default behavior in 2.x. |  |
| collector.netclass.ignored-devices | Regexp of net devices to ignore for netclass collector. | ^$ |
| collector.netclass.netlink | Use netlink to gather stats instead of /proc/net/dev. | false |
| collector.netclass_rtnl.with-stats | Expose the statistics for each network device, replacing netdev collector. |  |

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_subsystem_altnames_info | Non-numeric data of <iface> altname, value is always 1. | altname, device |
| node_subsystem_info | Non-numeric data from /sys/class/net/<iface>, value is always 1. | device, address, broadcast, duplex, operstate, adminstate, ifalias |
| node_subsystem_up | Value is 1 if operstate is 'up', 0 otherwise. | device |
