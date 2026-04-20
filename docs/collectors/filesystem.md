# filesystem

Exposes filesystem statistics including space usage and inode counts.

Status: enabled by default

## Platforms

- Linux
- Darwin
- FreeBSD
- NetBSD
- OpenBSD
- Dragonfly
- AIX

## Configuration

```
--collector.filesystem.mount-points-exclude  Regexp of mount points to exclude (mutually exclusive to mount-points-include)
--collector.filesystem.mount-points-include  Regexp of mount points to include (mutually exclusive to mount-points-exclude)
--collector.filesystem.fs-types-exclude      Regexp of filesystem types to exclude (mutually exclusive to fs-types-include)
--collector.filesystem.fs-types-include      Regexp of filesystem types to include (mutually exclusive to fs-types-exclude)
```

Default exclusions vary by platform. On Linux, virtual filesystems like `tmpfs`, `devtmpfs`, `sysfs`, `proc` are excluded by default.

## Data Sources

| Source | Description |
|--------|-------------|
| `/proc/self/mounts` | Mount points (Linux) |
| `/proc/self/mountinfo` | Mount info with major/minor device numbers (Linux) |
| `statfs(2)` | Filesystem statistics syscall |

Documentation:
- https://docs.kernel.org/filesystems/proc.html
- `statfs(2)` manpage

## Metrics

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `node_filesystem_size_bytes` | gauge | `device`, `mountpoint`, `fstype`, `device_error` | Filesystem size in bytes |
| `node_filesystem_free_bytes` | gauge | `device`, `mountpoint`, `fstype`, `device_error` | Filesystem free space in bytes |
| `node_filesystem_avail_bytes` | gauge | `device`, `mountpoint`, `fstype`, `device_error` | Filesystem space available to non-root users in bytes |
| `node_filesystem_files` | gauge | `device`, `mountpoint`, `fstype`, `device_error` | Filesystem total file nodes (inodes) |
| `node_filesystem_files_free` | gauge | `device`, `mountpoint`, `fstype`, `device_error` | Filesystem free file nodes (inodes) |
| `node_filesystem_readonly` | gauge | `device`, `mountpoint`, `fstype`, `device_error` | Filesystem read-only status (1 = read-only) |
| `node_filesystem_device_error` | gauge | `device`, `mountpoint`, `fstype`, `device_error` | Error occurred getting statistics (1 = error) |
| `node_filesystem_mount_info` | gauge | `device`, `major`, `minor`, `mountpoint` | Filesystem mount information (always 1) |
| `node_filesystem_purgeable_bytes` | gauge | `device`, `mountpoint`, `fstype`, `device_error` | Purgeable space in bytes (Darwin only) |

## Labels

| Label | Description |
|-------|-------------|
| `device` | Block device path (e.g., `/dev/sda1`) |
| `mountpoint` | Mount path (e.g., `/`, `/home`) |
| `fstype` | Filesystem type (e.g., `ext4`, `xfs`, `btrfs`) |
| `device_error` | Error message if device stat failed, empty otherwise |
| `major` | Device major number (mount_info only) |
| `minor` | Device minor number (mount_info only) |

## Notes

- `free_bytes` includes reserved blocks; `avail_bytes` is what non-root users can use
- When `device_error` is set (value = 1), only `readonly` and `device_error` metrics are emitted
- Duplicate mounts (same device, mountpoint, fstype) are deduplicated
- Network filesystems may cause hangs if unreachable; consider excluding with `--collector.filesystem.fs-types-exclude`
- `purgeable_bytes` is Darwin-specific and includes space reclaimable by the OS
