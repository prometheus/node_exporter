# zoneinfo collector

The zoneinfo collector exposes metrics about zoneinfo.

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_zoneinfo_high_pages | Zone watermark pages_high | node, zone |
| node_zoneinfo_low_pages | Zone watermark pages_low | node, zone |
| node_zoneinfo_managed_pages | Present pages managed by the buddy system | node, zone |
| node_zoneinfo_min_pages | Zone watermark pages_min | node, zone |
| node_zoneinfo_nr_active_anon_pages | Number of anonymous pages recently more used | node, zone |
| node_zoneinfo_nr_active_file_pages | Number of active pages with file-backing | node, zone |
| node_zoneinfo_nr_anon_pages | Number of anonymous pages currently used by the system | node, zone |
| node_zoneinfo_nr_anon_transparent_hugepages | Number of anonymous transparent huge pages currently used by the system | node, zone |
| node_zoneinfo_nr_dirtied_total | Page dirtyings since bootup | node, zone |
| node_zoneinfo_nr_dirty_pages | Number of dirty pages | node, zone |
| node_zoneinfo_nr_file_pages | Number of file pages | node, zone |
| node_zoneinfo_nr_free_pages | Total number of free pages in the zone | node, zone |
| node_zoneinfo_nr_inactive_anon_pages | Number of anonymous pages recently less used | node, zone |
| node_zoneinfo_nr_inactive_file_pages | Number of inactive pages with file-backing | node, zone |
| node_zoneinfo_nr_isolated_anon_pages | Temporary isolated pages from anon lru | node, zone |
| node_zoneinfo_nr_isolated_file_pages | Temporary isolated pages from file lru | node, zone |
| node_zoneinfo_nr_kernel_stacks | Number of kernel stacks | node, zone |
| node_zoneinfo_nr_mapped_pages | Number of mapped pages | node, zone |
| node_zoneinfo_nr_mlock_stack_pages | mlock()ed pages found and moved off LRU | node, zone |
| node_zoneinfo_nr_shmem_pages | Number of shmem pages (included tmpfs/GEM pages) | node, zone |
| node_zoneinfo_nr_slab_reclaimable_pages | Number of reclaimable slab pages | node, zone |
| node_zoneinfo_nr_slab_unreclaimable_pages | Number of unreclaimable slab pages | node, zone |
| node_zoneinfo_nr_unevictable_pages | Number of unevictable pages | node, zone |
| node_zoneinfo_nr_writeback_pages | Number of writeback pages | node, zone |
| node_zoneinfo_nr_written_total | Page writings since bootup | node, zone |
| node_zoneinfo_numa_foreign_total | Was intended here, hit elsewhere | node, zone |
| node_zoneinfo_numa_hit_total | Allocated in intended node | node, zone |
| node_zoneinfo_numa_interleave_total | Interleaver preferred this zone | node, zone |
| node_zoneinfo_numa_local_total | Allocation from local node | node, zone |
| node_zoneinfo_numa_miss_total | Allocated in non intended node | node, zone |
| node_zoneinfo_numa_other_total | Allocation from other node | node, zone |
| node_zoneinfo_present_pages | Physical pages existing within the zone | node, zone |
| node_zoneinfo_scanned_pages | Pages scanned since last reclaim | node, zone |
| node_zoneinfo_spanned_pages | Total pages spanned by the zone, including holes | node, zone |
