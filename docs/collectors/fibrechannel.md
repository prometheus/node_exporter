# fibrechannel collector

The fibrechannel collector exposes metrics about fibrechannel.

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_subsystem_info | Non-numeric data from /sys/class/fc_host/<host>, value is always 1. | fc_host, speed, port_state, port_type, port_id, port_name, fabric_name, symbolic_name, supported_classes, supported_speeds, dev_loss_tmo |
