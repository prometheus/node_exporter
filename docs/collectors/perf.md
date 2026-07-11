# perf collector

The perf collector exposes metrics about perf.

## Configuration Flags

| Flag | Description | Default |
| --- | --- | --- |
| collector.perf.cache-profilers | perf cache profilers that should be collected |  |
| collector.perf.cpus | List of CPUs from which perf metrics should be collected |  |
| collector.perf.disable-cache-profilers | disable perf cache profilers | false |
| collector.perf.disable-hardware-profilers | disable perf hardware profilers | false |
| collector.perf.disable-software-profilers | disable perf software profilers | false |
| collector.perf.hardware-profilers | perf hardware profilers that should be collected |  |
| collector.perf.software-profilers | perf software profilers that should be collected |  |
| collector.perf.tracepoint | perf tracepoint that should be collected |  |

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_perf_branch_instructions_total | Number of CPU branch instructions | cpu |
| node_perf_branch_misses_total | Number of CPU branch misses | cpu |
| node_perf_cache_bpu_read_hits_total | Number BPU read hits | cpu |
| node_perf_cache_bpu_read_misses_total | Number BPU read misses | cpu |
| node_perf_cache_l1_instr_read_misses_total | Number instruction L1 instruction read misses | cpu |
| node_perf_cache_l1d_read_hits_total | Number L1 data cache read hits | cpu |
| node_perf_cache_l1d_read_misses_total | Number L1 data cache read misses | cpu |
| node_perf_cache_l1d_write_hits_total | Number L1 data cache write hits | cpu |
| node_perf_cache_ll_read_hits_total | Number last level read hits | cpu |
| node_perf_cache_ll_read_misses_total | Number last level read misses | cpu |
| node_perf_cache_ll_write_hits_total | Number last level write hits | cpu |
| node_perf_cache_ll_write_misses_total | Number last level write misses | cpu |
| node_perf_cache_misses_total | Number of cache misses | cpu |
| node_perf_cache_refs_total | Number of cache references (non frequency scaled) | cpu |
| node_perf_cache_tlb_data_read_hits_total | Number of data TLB read hits | cpu |
| node_perf_cache_tlb_data_read_misses_total | Number of data TLB read misses | cpu |
| node_perf_cache_tlb_data_write_hits_total | Number of data TLB write hits | cpu |
| node_perf_cache_tlb_data_write_misses_total | Number of data TLB write misses | cpu |
| node_perf_cache_tlb_instr_read_hits_total | Number instruction TLB read hits | cpu |
| node_perf_cache_tlb_instr_read_misses_total | Number instruction TLB read misses | cpu |
| node_perf_context_switches_total | Number of context switches | cpu |
| node_perf_cpu_migrations_total | Number of CPU process migrations | cpu |
| node_perf_cpucycles_total | Number of CPU cycles (frequency scaled) | cpu |
| node_perf_instructions_total | Number of CPU instructions | cpu |
| node_perf_major_faults_total | Number of major page faults | cpu |
| node_perf_minor_faults_total | Number of minor page faults | cpu |
| node_perf_page_faults_total | Number of page faults | cpu |
| node_perf_ref_cpucycles_total | Number of CPU cycles | cpu |
| node_perf_stalled_cycles_backend_total | Number of stalled backend CPU cycles | cpu |
| node_perf_stalled_cycles_frontend_total | Number of stalled frontend CPU cycles | cpu |
