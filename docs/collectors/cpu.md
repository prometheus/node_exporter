# cpu collector

The cpu collector exposes metrics about cpu.

## Supported Platforms

- Linux

## Data Sources

- /proc/stat
- /sys/devices/system/cpu/

## Configuration Flags

| Flag | Description | Default |
| --- | --- | --- |
| collector.cpu.guest | Enables metric node_cpu_guest_seconds_total | true |
| collector.cpu.info | Enables metric cpu_info |  |
| collector.cpu.info.bugs-include | Filter the `bugs` field in cpuInfo with a value that must be a regular expression |  |
| collector.cpu.info.flags-include | Filter the `flags` field in cpuInfo with a value that must be a regular expression |  |

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_bug_info | The `bugs` field of CPU information from /proc/cpuinfo taken from the first core. | bug |
| node_context_switches_total | Number of context switches. | cpu |
| node_core_throttles_total | Number of times this CPU core has been throttled. | package, core |
| node_cpu_seconds_total | Seconds the CPUs spent in each mode. | cpu, mode |
| node_cpu_vulnerabilities_info | Details of each CPU vulnerability reported by sysfs. The value of the series is an int encoded state of the vulnerability. The same state is stored as a string in the label | codename, state, mitigation |
| node_flag_info | The `flags` field of CPU information from /proc/cpuinfo taken from the first core. | flag |
| node_flags | CPU flags. | cpu, flag |
| node_frequency_hertz | CPU frequency in hertz from /proc/cpuinfo. | package, core, cpu |
| node_guest_seconds_total | Seconds the CPUs spent in guests (VMs) for each mode. | cpu, mode |
| node_info | CPU information from /proc/cpuinfo. | package, core, cpu, vendor, family, model, model_name, microcode, stepping, cachesize |
| node_isolated | Whether each core is isolated, information from /sys/devices/system/cpu/isolated. | cpu |
| node_online | CPUs that are online and being scheduled. | cpu |
| node_package_throttles_total | Number of times this CPU package has been throttled. | package |
| node_physical_seconds_total | Seconds the physical CPUs spent in each mode. | cpu, mode |
| node_runqueue | Length of the run queue. | cpu |
| node_temperature_celsius | CPU temperature | cpu |
