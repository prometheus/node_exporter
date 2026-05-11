# cpu

Exposes CPU time statistics from `/proc/stat` and CPU metadata from `/proc/cpuinfo` and sysfs.

Status: enabled by default

## Platforms

- Linux
- Darwin
- Dragonfly
- FreeBSD
- NetBSD
- OpenBSD
- Solaris
- AIX

## Configuration

```
--collector.cpu.guest              Enable node_cpu_guest_seconds_total metric (default: true)
--collector.cpu.info               Enable node_cpu_info metric (default: false)
--collector.cpu.info.flags-include Regex filter for CPU flags to include in node_cpu_flag_info
--collector.cpu.info.bugs-include  Regex filter for CPU bugs to include in node_cpu_bug_info
```

Setting `--collector.cpu.info.flags-include` or `--collector.cpu.info.bugs-include` implicitly enables `--collector.cpu.info`.

## Data Sources

| Source | Description |
|--------|-------------|
| `/proc/stat` | CPU time counters per core and mode |
| `/proc/cpuinfo` | CPU metadata (vendor, model, flags, bugs) |
| `/sys/devices/system/cpu/cpu*/topology/` | Physical package and core IDs |
| `/sys/devices/system/cpu/cpu*/thermal_throttle/` | Thermal throttling counters |
| `/sys/devices/system/cpu/cpu*/online` | CPU online status |
| `/sys/devices/system/cpu/isolated` | Isolated CPUs list |

## Metrics

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `node_cpu_seconds_total` | counter | `cpu`, `mode` | Seconds the CPUs spent in each mode |
| `node_cpu_guest_seconds_total` | counter | `cpu`, `mode` | Seconds the CPUs spent in guest (VM) mode |
| `node_cpu_info` | gauge | `package`, `core`, `cpu`, `vendor`, `family`, `model`, `model_name`, `microcode`, `stepping`, `cachesize` | CPU metadata, always 1 |
| `node_cpu_frequency_hertz` | gauge | `package`, `core`, `cpu` | CPU frequency from /proc/cpuinfo (only when cpufreq collector disabled) |
| `node_cpu_flag_info` | gauge | `flag` | CPU flag presence from first core, always 1 |
| `node_cpu_bug_info` | gauge | `bug` | CPU bug presence from first core, always 1 |
| `node_cpu_core_throttles_total` | counter | `package`, `core` | Thermal throttle events per core |
| `node_cpu_package_throttles_total` | counter | `package` | Thermal throttle events per package |
| `node_cpu_isolated` | gauge | `cpu` | CPU isolation status (1 if isolated) |
| `node_cpu_online` | gauge | `cpu` | CPU online status (1 if online) |

## Labels

| Label | Description |
|-------|-------------|
| `cpu` | Logical CPU number (0-indexed) |
| `mode` | CPU time mode: `user`, `nice`, `system`, `idle`, `iowait`, `irq`, `softirq`, `steal` |
| `package` | Physical CPU package ID |
| `core` | Physical core ID within package |

## Notes

- `node_cpu_guest_seconds_total` values are also included in `node_cpu_seconds_total` (user and nice modes)
- Counter values may jump backwards on CPU hotplug events; the collector handles this by resetting stats when idle jumps back more than 3 seconds
- `node_cpu_flag_info` and `node_cpu_bug_info` are only exposed from the first CPU core
- `node_cpu_frequency_hertz` is only exposed when the `cpufreq` collector is disabled to avoid duplicate metrics
- Linux-specific metrics: throttle counters, isolated, online status
