# nfsd collector

The nfsd collector exposes metrics about nfsd.

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_nfsd_connections_total | Total number of NFSd TCP connections. | n/a |
| node_nfsd_disk_bytes_read_total | Total NFSd bytes read. | n/a |
| node_nfsd_disk_bytes_written_total | Total NFSd bytes written. | n/a |
| node_nfsd_file_handles_stale_total | Total number of NFSd stale file handles | n/a |
| node_nfsd_packets_total | Total NFSd network packets (sent+received) by protocol type. | proto |
| node_nfsd_read_ahead_cache_not_found_total | Total number of NFSd read ahead cache not found. | n/a |
| node_nfsd_read_ahead_cache_size_blocks | How large the read ahead cache is in blocks. | n/a |
| node_nfsd_reply_cache_hits_total | Total number of NFSd Reply Cache hits (client lost server response). | n/a |
| node_nfsd_reply_cache_misses_total | Total number of NFSd Reply Cache an operation that requires caching (idempotent). | n/a |
| node_nfsd_reply_cache_nocache_total | Total number of NFSd Reply Cache non-idempotent operations (rename/delete/…). | n/a |
| node_nfsd_requests_total | Total number NFSd Requests by method and protocol. | proto, method |
| node_nfsd_rpc_errors_total | Total number of NFSd RPC errors by error type. | error |
| node_nfsd_server_rpcs_total | Total number of NFSd RPCs. | n/a |
| node_nfsd_server_threads | Total number of NFSd kernel threads that are running. | n/a |
