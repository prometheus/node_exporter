# devstat collector

The devstat collector exposes metrics about devstat.

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_devstat_blocks_total | The total number of bytes given in terms of the devices blocksize. | device |
| node_devstat_blocks_transferred_total | The total number of blocks transferred. | device |
| node_devstat_busy_time_seconds_total | Total time the device had one or more transactions outstanding in seconds. | device |
| node_devstat_bytes_total | The total number of bytes transferred for reads and writes on the device. | device |
| node_devstat_duration_seconds_total | The total duration of transactions in seconds. | device, type |
| node_devstat_transfers_total | The total number of transactions completed. | device |
