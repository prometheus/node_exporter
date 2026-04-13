# runit collector

The runit collector exposes metrics about runit.

## Configuration Flags

| Flag | Description | Default |
| --- | --- | --- |
| collector.runit.servicedir | Path to runit service directory. | /etc/service |

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_service_desired_state | Desired state of runit service. | n/a |
| node_service_normal_state | Normal state of runit service. | n/a |
| node_service_state | State of runit service. | n/a |
| node_service_state_last_change_timestamp_seconds | Unix timestamp of the last runit service state change. | n/a |
