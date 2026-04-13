# tapestats collector

The tapestats collector exposes metrics about tapestats.

## Configuration Flags

| Flag | Description | Default |
| --- | --- | --- |
| collector.tapestats.ignored-devices | Regexp of devices to ignore for tapestats. | ^$ |

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_io_now | The number of I/Os currently outstanding to this device. | n/a |
| node_io_others_total | The number of I/Os issued to the tape drive other than read or write commands. The time taken to complete these commands uses the following calculation io_time_seconds_total-read_time_seconds_total-write_time_seconds_total | n/a |
| node_io_time_seconds_total | The amount of time spent waiting for all I/O to complete (including read and write). This includes tape movement commands such as seeking between file or set marks and implicit tape movement such as when rewind on close tape devices are used. | n/a |
| node_read_bytes_total | The number of bytes read from the tape drive. | n/a |
| node_read_time_seconds_total | The amount of time spent waiting for read requests to complete. | n/a |
| node_reads_completed_total | The number of read requests issued to the tape drive. | n/a |
| node_residual_total | The number of times during a read or write we found the residual amount to be non-zero. This should mean that a program is issuing a read larger thean the block size on tape. For write not all data made it to tape. | n/a |
| node_write_time_seconds_total | The amount of time spent waiting for write requests to complete. | n/a |
| node_writes_completed_total | The number of write requests issued to the tape drive. | n/a |
| node_written_bytes_total | The number of bytes written to the tape drive. | n/a |
