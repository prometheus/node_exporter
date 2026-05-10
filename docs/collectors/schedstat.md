# schedstat collector

The schedstat collector exposes metrics about schedstat.

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_schedstat_running_seconds_total | Number of seconds CPU spent running a process. | cpu |
| node_schedstat_timeslices_total | Number of timeslices executed by CPU. | cpu |
| node_schedstat_waiting_seconds_total | Number of seconds spent by processing waiting for this CPU. | cpu |
