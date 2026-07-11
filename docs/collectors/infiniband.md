# infiniband collector

The infiniband collector exposes metrics about infiniband.

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_subsystem_info | Non-numeric data from /sys/class/infiniband/<device>, value is always 1. | device, board_id, firmware_version, hca_type |
