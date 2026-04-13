# ethtool collector

The ethtool collector exposes metrics about ethtool.

## Configuration Flags

| Flag | Description | Default |
| --- | --- | --- |
| collector.ethtool.device-exclude | Regexp of ethtool devices to exclude (mutually exclusive to device-include). |  |
| collector.ethtool.device-include | Regexp of ethtool devices to include (mutually exclusive to device-exclude). |  |
| collector.ethtool.metrics-include | Regexp of ethtool stats to include. | .* |

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_ethtool_info | A metric with a constant '1' value labeled by bus_info, device, driver, expansion_rom_version, firmware_version, version. | bus_info, device, driver, expansion_rom_version, firmware_version, version |
| node_ethtool_received_bytes_total | Network interface bytes received | device |
| node_ethtool_received_dropped_total | Number of received frames dropped | device |
| node_ethtool_received_errors_total | Number of received frames with errors | device |
| node_ethtool_received_packets_total | Network interface packets received | device |
| node_ethtool_transmitted_bytes_total | Network interface bytes sent | device |
| node_ethtool_transmitted_errors_total | Number of sent frames with errors | device |
| node_ethtool_transmitted_packets_total | Network interface packets sent | device |
| node_network_advertised_speed_bytes | Combination of speeds and features offered by network device | device, duplex, mode |
| node_network_asymmetricpause_advertised | If this port device offers asymmetric pause capability | device |
| node_network_asymmetricpause_supported | If this port device supports asymmetric pause frames | device |
| node_network_autonegotiate | If this port is using autonegotiate | device |
| node_network_autonegotiate_advertised | If this port device offers autonegotiate | device |
| node_network_autonegotiate_supported | If this port device supports autonegotiate | device |
| node_network_pause_advertised | If this port device offers pause capability | device |
| node_network_pause_supported | If this port device supports pause frames | device |
| node_network_supported_port_info | Type of ports or PHYs supported by network device | device, type |
| node_network_supported_speed_bytes | Combination of speeds and features supported by network device | device, duplex, mode |
