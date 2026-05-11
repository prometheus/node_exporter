# meminfo

Exposes memory statistics from `/proc/meminfo`.

Status: enabled by default

## Platforms

- Linux
- Darwin
- OpenBSD
- NetBSD
- AIX

## Data Sources

| Source | Description |
|--------|-------------|
| `/proc/meminfo` | Memory statistics |

Kernel documentation: https://www.kernel.org/doc/Documentation/filesystems/proc.txt (search for "meminfo")

## Metrics

Metrics are dynamically generated from `/proc/meminfo` fields. Each field `FieldName` with value in kB becomes `node_memory_FieldName_bytes` (converted to bytes).

### Common Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `node_memory_MemTotal_bytes` | gauge | Total usable RAM |
| `node_memory_MemFree_bytes` | gauge | Free RAM |
| `node_memory_MemAvailable_bytes` | gauge | Available memory for starting new applications |
| `node_memory_Buffers_bytes` | gauge | Memory used by kernel buffers |
| `node_memory_Cached_bytes` | gauge | Memory used by page cache and slabs |
| `node_memory_SwapTotal_bytes` | gauge | Total swap space |
| `node_memory_SwapFree_bytes` | gauge | Free swap space |
| `node_memory_SwapCached_bytes` | gauge | Swap space cached in RAM |

### Active/Inactive Memory

| Metric | Type | Description |
|--------|------|-------------|
| `node_memory_Active_bytes` | gauge | Memory recently used |
| `node_memory_Inactive_bytes` | gauge | Memory not recently used |
| `node_memory_Active_anon_bytes` | gauge | Active anonymous memory |
| `node_memory_Inactive_anon_bytes` | gauge | Inactive anonymous memory |
| `node_memory_Active_file_bytes` | gauge | Active file-backed memory |
| `node_memory_Inactive_file_bytes` | gauge | Inactive file-backed memory |

### Slab Memory

| Metric | Type | Description |
|--------|------|-------------|
| `node_memory_Slab_bytes` | gauge | Kernel slab memory |
| `node_memory_SReclaimable_bytes` | gauge | Reclaimable slab memory |
| `node_memory_SUnreclaim_bytes` | gauge | Unreclaimable slab memory |

### Huge Pages

| Metric | Type | Description |
|--------|------|-------------|
| `node_memory_HugePages_Total` | gauge | Total huge pages (count, not bytes) |
| `node_memory_HugePages_Free` | gauge | Free huge pages (count) |
| `node_memory_HugePages_Rsvd` | gauge | Reserved huge pages (count) |
| `node_memory_HugePages_Surp` | gauge | Surplus huge pages (count) |
| `node_memory_Hugepagesize_bytes` | gauge | Size of each huge page |

### Virtual Memory

| Metric | Type | Description |
|--------|------|-------------|
| `node_memory_VmallocTotal_bytes` | gauge | Total vmalloc address space |
| `node_memory_VmallocUsed_bytes` | gauge | Used vmalloc address space |
| `node_memory_VmallocChunk_bytes` | gauge | Largest contiguous vmalloc block |

### Other

| Metric | Type | Description |
|--------|------|-------------|
| `node_memory_Dirty_bytes` | gauge | Memory waiting to be written to disk |
| `node_memory_Writeback_bytes` | gauge | Memory being written to disk |
| `node_memory_Mapped_bytes` | gauge | Files mapped into memory |
| `node_memory_Shmem_bytes` | gauge | Shared memory |
| `node_memory_KernelStack_bytes` | gauge | Kernel stack memory |
| `node_memory_PageTables_bytes` | gauge | Page table memory |
| `node_memory_CommitLimit_bytes` | gauge | Total memory available for allocation |
| `node_memory_Committed_AS_bytes` | gauge | Total memory allocated |

## Notes

- Available metrics vary by kernel version and configuration
- `MemAvailable` requires Linux 3.14+
- HugePages metrics are counts, not byte values
- All meminfo metrics are gauges
- Darwin, OpenBSD, NetBSD, and AIX have platform-specific implementations with different available metrics
