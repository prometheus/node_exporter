# rapl collector

The rapl collector exposes metrics about rapl.

## Configuration Flags

| Flag | Description | Default |
| --- | --- | --- |
| collector.rapl.enable-zone-label | Enables service unit metric unit_start_time_seconds |  |

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_rapl_joules_total | Current RAPL value in joules | index, path, rapl_zone |
