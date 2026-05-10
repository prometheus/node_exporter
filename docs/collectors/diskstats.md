# diskstats collector

The diskstats collector exposes metrics about diskstats.

## Configuration Flags

| Flag | Description | Default |
| --- | --- | --- |
| collector.diskstats.device-exclude | Regexp of diskstats devices to exclude (mutually exclusive to device-include). |  |
| collector.diskstats.device-include | Regexp of diskstats devices to include (mutually exclusive to device-exclude). |  |
| collector.diskstats.ignored-devices | DEPRECATED: Use collector.diskstats.device-exclude |  |

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_ata_rotation_rate_rpm | ATA disk rotation rate in RPMs (0 for SSDs). | device |
| node_ata_write_cache | ATA disk has a write cache. | device |
| node_ata_write_cache_enabled | ATA disk has its write cache enabled. | device |
| node_block_size_bytes | Size of the block device in bytes. | n/a |
| node_device_mapper_info | Info about disk device mapper. | device, name, uuid, vg_name, lv_name, lv_layer |
| node_discard_time_seconds_total | This is the total number of seconds spent by all discards. | n/a |
| node_discarded_sectors_total | The total number of sectors discarded successfully. | n/a |
| node_discards_completed_total | The total number of discards completed successfully. | n/a |
| node_discards_merged_total | The total number of discards merged. | n/a |
| node_disk_io_time_seconds_total | Total seconds spent doing I/Os. | n/a |
| node_disk_read_bytes_total | The total number of bytes read successfully. | n/a |
| node_disk_read_time_seconds_total | The total number of seconds spent by all reads. | n/a |
| node_disk_reads_completed_total | The total number of reads completed successfully. | n/a |
| node_disk_write_time_seconds_total | This is the total number of seconds spent by all writes. | n/a |
| node_disk_writes_completed_total | The total number of writes completed successfully. | n/a |
| node_disk_written_bytes_total | The total number of bytes written successfully. | n/a |
| node_filesystem_info | Info about disk filesystem. | device, type, usage, uuid, version |
| node_flush_requests_time_seconds_total | This is the total number of seconds spent by all flush requests. | n/a |
| node_flush_requests_total | The total number of flush requests completed successfully | n/a |
| node_info | Info of /sys/block/<block_device>. | device, major, minor, path, wwn, model, serial, revision, rotational |
| node_io_now | The number of I/Os currently in progress. | n/a |
| node_io_time_weighted_seconds_total | The weighted # of seconds spent doing I/Os. | n/a |
| node_queue_depth | Number of requests in the queue. | n/a |
| node_read_errors_total | The total number of read errors. | n/a |
| node_read_retries_total | The total number of read retries. | n/a |
| node_read_sectors_total | The total number of sectors read successfully. | n/a |
| node_read_time_seconds_total | The total time spent servicing read requests. | n/a |
| node_reads_merged_total | The total number of reads merged. | n/a |
| node_transfers_to_disk_total | The total number of transfers from disk. | n/a |
| node_transfers_total | The total number of transfers to/from disk. | n/a |
| node_write_errors_total | The total number of write errors. | n/a |
| node_write_retries_total | The total number of write retries. | n/a |
| node_write_time_seconds_total | The total time spent servicing write requests. | n/a |
| node_writes_merged_total | The number of writes merged. | n/a |
| node_written_sectors_total | The total number of sectors written successfully. | n/a |
