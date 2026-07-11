# ipvs collector

The ipvs collector exposes metrics about ipvs.

## Configuration Flags

| Flag | Description | Default |
| --- | --- | --- |
| collector.ipvs.backend-labels | Comma separated list for IPVS backend stats labels. |  |

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_ipvs_backend_connections_active | The current active connections by local and remote address. | n/a |
| node_ipvs_backend_connections_inactive | The current inactive connections by local and remote address. | n/a |
| node_ipvs_backend_weight | The current backend weight by local and remote address. | n/a |
| node_ipvs_connections_total | The total number of connections made. | n/a |
| node_ipvs_incoming_bytes_total | The total amount of incoming data. | n/a |
| node_ipvs_incoming_packets_total | The total number of incoming packets. | n/a |
| node_ipvs_outgoing_bytes_total | The total amount of outgoing data. | n/a |
| node_ipvs_outgoing_packets_total | The total number of outgoing packets. | n/a |
