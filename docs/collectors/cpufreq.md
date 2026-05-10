# cpufreq collector

The cpufreq collector exposes metrics about cpufreq.

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_frequency_hertz | Current CPU thread frequency in hertz. | cpu |
| node_frequency_max_hertz | Maximum CPU thread frequency in hertz. | cpu |
| node_frequency_min_hertz | Minimum CPU thread frequency in hertz. | cpu |
| node_scaling_frequency_hertz | Current scaled CPU thread frequency in hertz. | cpu |
| node_scaling_frequency_max_hertz | Maximum scaled CPU thread frequency in hertz. | cpu |
| node_scaling_frequency_min_hertz | Minimum scaled CPU thread frequency in hertz. | cpu |
| node_scaling_governor | Current enabled CPU frequency governor. | cpu, governor |
