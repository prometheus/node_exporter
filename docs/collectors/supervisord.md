# supervisord collector

The supervisord collector exposes metrics about supervisord.

## Configuration Flags

| Flag | Description | Default |
| --- | --- | --- |
| collector.supervisord.url | XML RPC endpoint. | http://localhost:9001/RPC2 |

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_supervisord_exit_status | Process Exit Status | n/a |
| node_supervisord_start_time_seconds | Process start time | n/a |
| node_supervisord_state | Process State | n/a |
| node_supervisord_up | Process Up | n/a |
