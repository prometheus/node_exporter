# stat

Exposes kernel and system statistics from `/proc/stat`.

Status: enabled by default

## Platforms

- Linux

## Configuration

```
--collector.stat.softirq  Export softirq calls per vector (default: false)
```

## Data Sources

| Source | Description |
|--------|-------------|
| `/proc/stat` | Kernel/system statistics |

Documentation:
- https://docs.kernel.org/filesystems/proc.html
- `proc(5)` manpage

## Metrics

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `node_intr_total` | counter | | Total number of interrupts serviced |
| `node_context_switches_total` | counter | | Total number of context switches |
| `node_forks_total` | counter | | Total number of forks |
| `node_boot_time_seconds` | gauge | | Node boot time in Unix timestamp |
| `node_procs_running` | gauge | | Number of processes in runnable state |
| `node_procs_blocked` | gauge | | Number of processes blocked waiting for I/O |

### Softirq Metrics (--collector.stat.softirq)

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `node_softirqs_total` | counter | `vector` | Number of softirq calls per vector |

Softirq vectors:

| Vector | Description |
|--------|-------------|
| `hi` | High-priority tasklets |
| `timer` | Timer interrupts |
| `net_tx` | Network transmit |
| `net_rx` | Network receive |
| `block` | Block device |
| `block_iopoll` | Block I/O polling |
| `tasklet` | Tasklet processing |
| `sched` | Scheduler |
| `hrtimer` | High-resolution timer |
| `rcu` | Read-copy-update |

## Notes

- `node_boot_time_seconds` is a Unix timestamp; use `time() - node_boot_time_seconds` for uptime
- `node_intr_total` is the sum of all interrupt counts; for per-interrupt details, use the `interrupts` collector
- `node_procs_running` and `node_procs_blocked` are instantaneous values
- Softirq metrics are disabled by default to reduce cardinality; for detailed softirq stats, see the `softirqs` collector
