# diskstats

Exposes disk I/O statistics from `/proc/diskstats` and block device metadata from sysfs and udev.

Status: enabled by default

## Platforms

- Linux
- Darwin
- OpenBSD
- AIX

## Configuration

```
--collector.diskstats.device-include  Regexp of devices to include (mutually exclusive with device-exclude)
--collector.diskstats.device-exclude  Regexp of devices to exclude (default: ^(z?ram|loop|fd|(h|s|v|xv)d[a-z]|nvme\d+n\d+p)\d+$)
```

## Data Sources

| Source | Description |
|--------|-------------|
| `/proc/diskstats` | Disk I/O statistics |
| `/sys/block/<device>/` | Block device attributes |
| `/sys/block/<device>/queue/` | Block device queue stats |
| `/run/udev/data/b<major>:<minor>` | Udev device properties |

Kernel documentation: https://www.kernel.org/doc/Documentation/iostats.txt

## Metrics

### I/O Statistics

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `node_disk_reads_completed_total` | counter | `device` | Total number of reads completed successfully |
| `node_disk_reads_merged_total` | counter | `device` | Total number of reads merged |
| `node_disk_read_bytes_total` | counter | `device` | Total number of bytes read successfully |
| `node_disk_read_time_seconds_total` | counter | `device` | Total seconds spent by all reads |
| `node_disk_writes_completed_total` | counter | `device` | Total number of writes completed successfully |
| `node_disk_writes_merged_total` | counter | `device` | Total number of writes merged |
| `node_disk_written_bytes_total` | counter | `device` | Total number of bytes written successfully |
| `node_disk_write_time_seconds_total` | counter | `device` | Total seconds spent by all writes |
| `node_disk_io_now` | gauge | `device` | Number of I/Os currently in progress |
| `node_disk_io_time_seconds_total` | counter | `device` | Total seconds spent doing I/Os |
| `node_disk_io_time_weighted_seconds_total` | counter | `device` | Weighted seconds spent doing I/Os |

### Discard Statistics (Linux 4.18+)

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `node_disk_discards_completed_total` | counter | `device` | Total number of discards completed successfully |
| `node_disk_discards_merged_total` | counter | `device` | Total number of discards merged |
| `node_disk_discarded_sectors_total` | counter | `device` | Total number of sectors discarded successfully |
| `node_disk_discard_time_seconds_total` | counter | `device` | Total seconds spent by all discards |

### Flush Statistics (Linux 5.5+)

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `node_disk_flush_requests_total` | counter | `device` | Total number of flush requests completed successfully |
| `node_disk_flush_requests_time_seconds_total` | counter | `device` | Total seconds spent by all flush requests |

### Device Info

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `node_disk_info` | gauge | `device`, `major`, `minor`, `path`, `wwn`, `model`, `serial`, `revision`, `rotational` | Block device info, always 1 |
| `node_disk_filesystem_info` | gauge | `device`, `type`, `usage`, `uuid`, `version` | Filesystem info from udev, always 1 |
| `node_disk_device_mapper_info` | gauge | `device`, `name`, `uuid`, `vg_name`, `lv_name`, `lv_layer` | Device mapper info, always 1 |

### ATA Device Attributes

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `node_disk_ata_write_cache` | gauge | `device` | ATA disk has a write cache (1 if true) |
| `node_disk_ata_write_cache_enabled` | gauge | `device` | ATA disk write cache is enabled (1 if true) |
| `node_disk_ata_rotation_rate_rpm` | gauge | `device` | ATA disk rotation rate in RPM (0 for SSDs) |

## Notes

- Sector sizes in `/proc/diskstats` are always 512 bytes regardless of actual device sector size
- Time values in the kernel are in milliseconds; the collector converts to seconds
- Udev info metrics require readable `/run/udev/data/` directory
- Discard and flush metrics availability depends on kernel version
- The default exclude pattern filters out partition devices and RAM/loop devices
