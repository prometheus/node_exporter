# nfs collector

The nfs collector exposes metrics about nfs.

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_nfs_connections_total | Total number of NFSd TCP connections. | n/a |
| node_nfs_packets_total | Total NFSd network packets (sent+received) by protocol type. | protocol |
| node_nfs_requests_total | Number of NFS procedures invoked. | proto, method |
| node_nfs_rpc_authentication_refreshes_total | Number of RPC authentication refreshes performed. | n/a |
| node_nfs_rpc_retransmissions_total | Number of RPC transmissions performed. | n/a |
| node_nfs_rpcs_total | Total number of RPCs performed. | n/a |
