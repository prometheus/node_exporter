# Collector Documentation

Per-collector metric documentation. Each file documents one collector.

## Available Documentation

- [cpu](cpu.md) - CPU time statistics and metadata
- [cpufreq](cpufreq.md) - CPU frequency scaling statistics
- [diskstats](diskstats.md) - Disk I/O statistics
- [filesystem](filesystem.md) - Filesystem space and inode statistics
- [hwmon](hwmon.md) - Hardware monitoring sensors
- [meminfo](meminfo.md) - Memory statistics
- [netdev](netdev.md) - Network interface statistics
- [netstat](netstat.md) - Network protocol statistics
- [stat](stat.md) - Kernel/system statistics

## Structure

See [_TEMPLATE.md](_TEMPLATE.md) for the documentation template.

## Naming

Files are named `<collector_name>.md` matching the collector registration name (e.g., `cpu.md`, `filesystem.md`).

## Contributing

When adding or modifying a collector:
1. Update or create the corresponding documentation file
2. Ensure all metrics are listed with correct types and labels
3. Document any configuration flags
