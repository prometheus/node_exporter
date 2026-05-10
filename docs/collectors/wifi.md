# wifi collector

The wifi collector exposes metrics about wifi.

## Configuration Flags

| Flag | Description | Default |
| --- | --- | --- |
| collector.wifi.fixtures | test fixtures to use for wifi collector metrics |  |

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_wifi_interface_frequency_hertz | The current frequency a WiFi interface is operating at, in hertz. | device |
| node_wifi_station_beacon_loss_total | The total number of times a station has detected a beacon loss. | n/a |
| node_wifi_station_connected_seconds_total | The total number of seconds a station has been connected to an access point. | n/a |
| node_wifi_station_inactive_seconds | The number of seconds since any wireless activity has occurred on a station. | n/a |
| node_wifi_station_info | Labeled WiFi interface station information as provided by the operating system. | device, bssid, ssid, mode |
| node_wifi_station_receive_bits_per_second | The current WiFi receive bitrate of a station, in bits per second. | n/a |
| node_wifi_station_receive_bytes_total | The total number of bytes received by a WiFi station. | n/a |
| node_wifi_station_received_packets_total | The total number of packets received by a station. | n/a |
| node_wifi_station_signal_dbm | The current WiFi signal strength, in decibel-milliwatts (dBm). | n/a |
| node_wifi_station_transmit_bits_per_second | The current WiFi transmit bitrate of a station, in bits per second. | n/a |
| node_wifi_station_transmit_bytes_total | The total number of bytes transmitted by a WiFi station. | n/a |
| node_wifi_station_transmit_failed_total | The total number of times a station has failed to send a packet. | n/a |
| node_wifi_station_transmit_retries_total | The total number of times a station has had to retry while sending a packet. | n/a |
| node_wifi_station_transmitted_packets_total | The total number of packets transmitted by a station. | n/a |
