# stat collector

The stat collector exposes metrics about stat.

## Configuration Flags

| Flag | Description | Default |
| --- | --- | --- |
| collector.stat.softirq | Export softirq calls per vector | false |

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_boot_time_seconds | Node boot time, in unixtime. | n/a |
| node_context_switches_total | Total number of context switches. | n/a |
| node_forks_total | Total number of forks. | n/a |
| node_intr_total | Total number of interrupts serviced. | n/a |
| node_procs_blocked | Number of processes blocked waiting for I/O to complete. | n/a |
| node_procs_running | Number of processes in runnable state. | n/a |
| node_softirqs_total | Number of softirq calls. | vector |
