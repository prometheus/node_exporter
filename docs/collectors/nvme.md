# nvme collector

The nvme collector exposes metrics about nvme.

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_nvme_info | Non-numeric data from /sys/class/nvme/<device>, value is always 1. | device, firmware_revision, model, serial, state, cntlid |
| node_nvme_namespace_capacity_bytes | Capacity of the NVMe namespace in bytes. Computed as namespace_size * namespace_logical_block_size | device, nsid |
| node_nvme_namespace_info | Information about NVMe namespaces. Value is always 1 | device, nsid, ana_state |
| node_nvme_namespace_logical_block_size_bytes | Logical block size of the NVMe namespace in bytes. Usually 4Kb. Available in /sys/class/nvme/<device>/<namespace>/queue/logical_block_size | device, nsid |
| node_nvme_namespace_size_bytes | Size of the NVMe namespace in bytes. Available in /sys/class/nvme/<device>/<namespace>/size | device, nsid |
| node_nvme_namespace_used_bytes | Used space of the NVMe namespace in bytes. Available in /sys/class/nvme/<device>/<namespace>/nuse | device, nsid |
