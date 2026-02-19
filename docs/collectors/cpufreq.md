# cpufreq

Exposes CPU frequency scaling statistics from sysfs.

Status: enabled by default

## Platforms

- Linux
- Solaris

## Data Sources

| Source | Description |
|--------|-------------|
| `/sys/devices/system/cpu/cpu*/cpufreq/` | Per-CPU frequency scaling data |

Kernel documentation:
- https://www.kernel.org/doc/Documentation/cpu-freq/user-guide.txt
- https://www.kernel.org/doc/Documentation/cpu-freq/governors.txt

## Metrics

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `node_cpu_frequency_hertz` | gauge | `cpu` | Current CPU thread frequency in hertz |
| `node_cpu_frequency_min_hertz` | gauge | `cpu` | Minimum CPU thread frequency in hertz |
| `node_cpu_frequency_max_hertz` | gauge | `cpu` | Maximum CPU thread frequency in hertz |
| `node_cpu_scaling_frequency_hertz` | gauge | `cpu` | Current scaled CPU thread frequency in hertz |
| `node_cpu_scaling_frequency_min_hertz` | gauge | `cpu` | Minimum scaled CPU thread frequency in hertz |
| `node_cpu_scaling_frequency_max_hertz` | gauge | `cpu` | Maximum scaled CPU thread frequency in hertz |
| `node_cpu_scaling_governor` | gauge | `cpu`, `governor` | Current CPU frequency governor (1 if active, 0 otherwise) |

## Labels

| Label | Description |
|-------|-------------|
| `cpu` | CPU name from sysfs (e.g., `cpu0`) |
| `governor` | Frequency governor name (e.g., `performance`, `powersave`, `ondemand`) |

## Notes

- Sysfs values are in kHz; the collector converts to Hz
- Metrics without `scaling` in the name reflect hardware limits from cpuinfo files; `scaling_*` metrics reflect current governor policy limits
- `node_cpu_scaling_governor` emits one metric per available governor per CPU, with value 1 for the active governor
- When this collector is enabled, the `cpu` collector does not expose `node_cpu_frequency_hertz` to avoid duplication
