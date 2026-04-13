# partition collector

The partition collector exposes metrics about partition.

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_partition_cpus_online | Number of online CPUs in the partition. | n/a |
| node_partition_cpus_pool | Number of physical CPUs in the pool. | n/a |
| node_partition_cpus_sys | Number of physical CPUs in the system. | n/a |
| node_partition_entitled_capacity | Entitled processor capacity of the partition in CPU units (e.g. 1.0 = one core). | n/a |
| node_partition_memory_max | Maximum memory of the partition in bytes. | n/a |
| node_partition_memory_online | Online memory of the partition in bytes. | n/a |
| node_partition_power_save_mode | Power save mode of the partition (1 for enabled, 0 for disabled). | n/a |
| node_partition_smt_threads | Number of SMT threads per core. | n/a |
