# filesystem collector

The filesystem collector exposes metrics about filesystem.

## Configuration Flags

| Flag | Description | Default |
| --- | --- | --- |
| collector.filesystem.fs-types-exclude | Regexp of filesystem types to exclude for filesystem collector. (mutually exclusive to fs-types-include) |  |
| collector.filesystem.fs-types-include | Regexp of filesystem types to exclude for filesystem collector. (mutually exclusive to fs-types-exclude) |  |
| collector.filesystem.ignored-fs-types | Regexp of filesystem types to ignore for filesystem collector. |  |
| collector.filesystem.ignored-mount-points | Regexp of mount points to ignore for filesystem collector. |  |
| collector.filesystem.mount-points-exclude | Regexp of mount points to exclude for filesystem collector. (mutually exclusive to mount-points-include) |  |
| collector.filesystem.mount-points-include | Regexp of mount points to include for filesystem collector. (mutually exclusive to mount-points-exclude) |  |
| collector.filesystem.mount-timeout | how long to wait for a mount to respond before marking it as stale | 5s |
| collector.filesystem.stat-workers | how many stat calls to process simultaneously | 4 |

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_filesystem_avail_bytes | Filesystem space available to non-root users in bytes. | n/a |
| node_filesystem_device_error | Whether an error occurred while getting statistics for the given device. | n/a |
| node_filesystem_files | Filesystem total file nodes. | n/a |
| node_filesystem_files_free | Filesystem total free file nodes. | n/a |
| node_filesystem_free_bytes | Filesystem free space in bytes. | n/a |
| node_filesystem_mount_info | Filesystem mount information. | device, major, minor, mountpoint |
| node_filesystem_purgeable_bytes | Filesystem space available including purgeable space (MacOS specific). | n/a |
| node_filesystem_readonly | Filesystem read-only status. | n/a |
| node_filesystem_size_bytes | Filesystem size in bytes. | n/a |
