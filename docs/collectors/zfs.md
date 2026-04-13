# zfs collector

The zfs collector exposes metrics about zfs.

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_zfs_abdstats_linear_count_total | ZFS ARC buffer data linear count | n/a |
| node_zfs_abdstats_linear_data_bytes | ZFS ARC buffer data linear data size | n/a |
| node_zfs_abdstats_scatter_chunk_waste_bytes | ZFS ARC buffer data scatter chunk waste | n/a |
| node_zfs_abdstats_scatter_count_total | ZFS ARC buffer data scatter count | n/a |
| node_zfs_abdstats_scatter_data_bytes | ZFS ARC buffer data scatter data size | n/a |
| node_zfs_abdstats_struct_bytes | ZFS ARC buffer data struct size | n/a |
| node_zfs_arcstats_anon_bytes | ZFS ARC anon size | n/a |
| node_zfs_arcstats_c_bytes | ZFS ARC target size | n/a |
| node_zfs_arcstats_c_max_bytes | ZFS ARC maximum size | n/a |
| node_zfs_arcstats_c_min_bytes | ZFS ARC minimum size | n/a |
| node_zfs_arcstats_data_bytes | ZFS ARC data size | n/a |
| node_zfs_arcstats_demand_data_hits_total | ZFS ARC demand data hits | n/a |
| node_zfs_arcstats_demand_data_misses_total | ZFS ARC demand data misses | n/a |
| node_zfs_arcstats_demand_metadata_hits_total | ZFS ARC demand metadata hits | n/a |
| node_zfs_arcstats_demand_metadata_misses_total | ZFS ARC demand metadata misses | n/a |
| node_zfs_arcstats_hdr_bytes | ZFS ARC header size | n/a |
| node_zfs_arcstats_hits_total | ZFS ARC hits | n/a |
| node_zfs_arcstats_mfu_bytes | ZFS ARC MFU size | n/a |
| node_zfs_arcstats_mfu_ghost_hits_total | ZFS ARC MFU ghost hits | n/a |
| node_zfs_arcstats_mfu_ghost_size | ZFS ARC MFU ghost size | n/a |
| node_zfs_arcstats_misses_total | ZFS ARC misses | n/a |
| node_zfs_arcstats_mru_bytes | ZFS ARC MRU size | n/a |
| node_zfs_arcstats_mru_ghost_bytes | ZFS ARC MRU ghost size | n/a |
| node_zfs_arcstats_mru_ghost_hits_total | ZFS ARC MRU ghost hits | n/a |
| node_zfs_arcstats_other_bytes | ZFS ARC other size | n/a |
| node_zfs_arcstats_p_bytes | ZFS ARC MRU target size | n/a |
| node_zfs_arcstats_size_bytes | ZFS ARC size | n/a |
| node_zfs_zfetchstats_hits_total | ZFS cache fetch hits | n/a |
| node_zfs_zfetchstats_misses_total | ZFS cache fetch misses | n/a |
| node_zfs_zpool_state | kstat.zfs.misc.state | zpool, state |
