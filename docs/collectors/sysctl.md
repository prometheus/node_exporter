# sysctl collector

The sysctl collector exposes metrics about sysctl.

## Configuration Flags

| Flag | Description | Default |
| --- | --- | --- |
| collector.sysctl.include | Select sysctl metrics to include |  |
| collector.sysctl.include-info | Select sysctl metrics to include as info metrics |  |

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_sysctl_info | sysctl info | name, value, index |
