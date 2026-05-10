# softnet collector

The softnet collector exposes metrics about softnet.

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_softnet_backlog_len | Softnet backlog status | cpu |
| node_softnet_cpu_collision_total | Number of collision occur while obtaining device lock while transmitting | cpu |
| node_softnet_dropped_total | Number of dropped packets | cpu |
| node_softnet_flow_limit_count_total | Number of times flow limit has been reached | cpu |
| node_softnet_processed_total | Number of processed packets | cpu |
| node_softnet_received_rps_total | Number of times cpu woken up received_rps | cpu |
| node_softnet_times_squeezed_total | Number of times processing packets ran out of quota | cpu |
