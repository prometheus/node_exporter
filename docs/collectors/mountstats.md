# mountstats collector

The mountstats collector exposes metrics about mountstats.

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_mountstats_nfs_age_seconds_total | The age of the NFS mount in seconds. | n/a |
| node_mountstats_nfs_direct_read_bytes_total | Number of bytes read using the read() syscall in O_DIRECT mode. | n/a |
| node_mountstats_nfs_direct_write_bytes_total | Number of bytes written using the write() syscall in O_DIRECT mode. | n/a |
| node_mountstats_nfs_event_attribute_invalidate_total | Number of times cached inode attributes are invalidated. | n/a |
| node_mountstats_nfs_event_data_invalidate_total | Number of times an inode cache is cleared. | n/a |
| node_mountstats_nfs_event_dnode_revalidate_total | Number of times cached dentry nodes are re-validated from the server. | n/a |
| node_mountstats_nfs_event_inode_revalidate_total | Number of times cached inode attributes are re-validated from the server. | n/a |
| node_mountstats_nfs_event_jukebox_delay_total | Number of times the NFS server indicated EJUKEBOX; retrieving data from offline storage. | n/a |
| node_mountstats_nfs_event_pnfs_read_total | Number of NFS v4.1+ pNFS reads. | n/a |
| node_mountstats_nfs_event_pnfs_write_total | Number of NFS v4.1+ pNFS writes. | n/a |
| node_mountstats_nfs_event_short_read_total | Number of times the NFS server gave less data than expected while reading. | n/a |
| node_mountstats_nfs_event_short_write_total | Number of times the NFS server wrote less data than expected while writing. | n/a |
| node_mountstats_nfs_event_silly_rename_total | Number of times a file was removed while still open by another process. | n/a |
| node_mountstats_nfs_event_truncation_total | Number of times files have been truncated. | n/a |
| node_mountstats_nfs_event_vfs_access_total | Number of times permissions have been checked. | n/a |
| node_mountstats_nfs_event_vfs_file_release_total | Number of times files have been closed and released. | n/a |
| node_mountstats_nfs_event_vfs_flush_total | Number of pending writes that have been forcefully flushed to the server. | n/a |
| node_mountstats_nfs_event_vfs_fsync_total | Number of times fsync() has been called on directories and files. | n/a |
| node_mountstats_nfs_event_vfs_getdents_total | Number of times directory entries have been read with getdents(). | n/a |
| node_mountstats_nfs_event_vfs_lock_total | Number of times locking has been attempted on a file. | n/a |
| node_mountstats_nfs_event_vfs_lookup_total | Number of times a directory lookup has occurred. | n/a |
| node_mountstats_nfs_event_vfs_open_total | Number of times cached inode attributes are invalidated. | n/a |
| node_mountstats_nfs_event_vfs_read_page_total | Number of pages read directly via mmap()'d files. | n/a |
| node_mountstats_nfs_event_vfs_read_pages_total | Number of times a group of pages have been read. | n/a |
| node_mountstats_nfs_event_vfs_setattr_total | Number of times directory entries have been read with getdents(). | n/a |
| node_mountstats_nfs_event_vfs_update_page_total | Number of updates (and potential writes) to pages. | n/a |
| node_mountstats_nfs_event_vfs_write_page_total | Number of pages written directly via mmap()'d files. | n/a |
| node_mountstats_nfs_event_vfs_write_pages_total | Number of times a group of pages have been written. | n/a |
| node_mountstats_nfs_event_write_extension_total | Number of times a file has been grown due to writes beyond its existing end. | n/a |
| node_mountstats_nfs_operations_major_timeouts_total | Number of times a request has had a major timeout for a given operation. | n/a |
| node_mountstats_nfs_operations_queue_time_seconds_total | Duration all requests spent queued for transmission for a given operation before they were sent, in seconds. | n/a |
| node_mountstats_nfs_operations_received_bytes_total | Number of bytes received for a given operation, including RPC headers and payload. | n/a |
| node_mountstats_nfs_operations_request_time_seconds_total | Duration all requests took from when a request was enqueued to when it was completely handled for a given operation, in seconds. | n/a |
| node_mountstats_nfs_operations_requests_total | Number of requests performed for a given operation. | n/a |
| node_mountstats_nfs_operations_response_time_seconds_total | Duration all requests took to get a reply back after a request for a given operation was transmitted, in seconds. | n/a |
| node_mountstats_nfs_operations_sent_bytes_total | Number of bytes sent for a given operation, including RPC headers and payload. | n/a |
| node_mountstats_nfs_operations_transmissions_total | Number of times an actual RPC request has been transmitted for a given operation. | n/a |
| node_mountstats_nfs_read_bytes_total | Number of bytes read using the read() syscall. | n/a |
| node_mountstats_nfs_read_pages_total | Number of pages read directly via mmap()'d files. | n/a |
| node_mountstats_nfs_total_read_bytes_total | Number of bytes read from the NFS server, in total. | n/a |
| node_mountstats_nfs_total_write_bytes_total | Number of bytes written to the NFS server, in total. | n/a |
| node_mountstats_nfs_transport_backlog_queue_total | Total number of items added to the RPC backlog queue. | n/a |
| node_mountstats_nfs_transport_bad_transaction_ids_total | Number of times the NFS server sent a response with a transaction ID unknown to this client. | n/a |
| node_mountstats_nfs_transport_bind_total | Number of times the client has had to establish a connection from scratch to the NFS server. | n/a |
| node_mountstats_nfs_transport_connect_total | Number of times the client has made a TCP connection to the NFS server. | n/a |
| node_mountstats_nfs_transport_idle_time_seconds | Duration since the NFS mount last saw any RPC traffic, in seconds. | n/a |
| node_mountstats_nfs_transport_maximum_rpc_slots | Maximum number of simultaneously active RPC requests ever used. | n/a |
| node_mountstats_nfs_transport_pending_queue_total | Total number of items added to the RPC transmission pending queue. | n/a |
| node_mountstats_nfs_transport_receives_total | Number of RPC responses for this mount received from the NFS server. | n/a |
| node_mountstats_nfs_transport_sending_queue_total | Total number of items added to the RPC transmission sending queue. | n/a |
| node_mountstats_nfs_transport_sends_total | Number of RPC requests for this mount sent to the NFS server. | n/a |
| node_mountstats_nfs_write_bytes_total | Number of bytes written using the write() syscall. | n/a |
| node_mountstats_nfs_write_pages_total | Number of pages written directly via mmap()'d files. | n/a |
