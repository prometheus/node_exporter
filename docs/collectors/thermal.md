# thermal collector

The thermal collector exposes metrics about thermal.

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_cpu_available_cpu | Reflects how many, if any, CPUs have been taken offline. Represented as an integer number of CPUs (0 - Max CPUs). | n/a |
| node_cpu_scheduler_limit_ratio | Represents the percentage (0-100) of CPU time available. 100% at normal operation. The OS may limit this time for a percentage less than 100%. | n/a |
| node_cpu_speed_limit_ratio | Defines the speed & voltage limits placed on the CPU. Represented as a percentage (0-100) of maximum CPU speed. | n/a |
| node_cur_state | Current throttle state of the cooling device | name, type |
| node_max_state | Maximum throttle state of the cooling device | name, type |
| node_temp | Zone temperature in Celsius | zone, type |
| node_temperature_celsius | Temperature of the thermal sensor in Celsius. | sensor |
