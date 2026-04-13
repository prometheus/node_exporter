# bcachefs collector

The bcachefs collector exposes metrics about bcachefs.

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_bcachefs_btree_cache_size_bytes | Btree cache memory usage in bytes. | uuid |
| node_bcachefs_btree_write_average_size_bytes | Average btree write size by type. | uuid, type |
| node_bcachefs_btree_writes_total | Number of btree writes by type. | uuid, type |
| node_bcachefs_compression_compressed_bytes | Compressed size by algorithm. | uuid, algorithm |
| node_bcachefs_compression_uncompressed_bytes | Uncompressed size by algorithm. | uuid, algorithm |
| node_bcachefs_device_bucket_size_bytes | Bucket size in bytes. | uuid, device |
| node_bcachefs_device_buckets | Total number of buckets. | uuid, device |
| node_bcachefs_device_durability | Device durability setting. | uuid, device |
| node_bcachefs_device_info | Device information. | uuid, device, label, state |
| node_bcachefs_device_io_done_bytes_total | IO bytes by operation type and data type. | uuid, device, operation, data_type |
| node_bcachefs_device_io_errors_total | IO errors by error type. | uuid, device, type |
| node_bcachefs_errors_total | Error count by error type. | uuid, error_type |
| node_bcachefs_info | Filesystem information. | uuid |
