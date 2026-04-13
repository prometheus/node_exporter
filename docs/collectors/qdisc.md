# qdisc collector

The qdisc collector exposes metrics about qdisc.

## Configuration Flags

| Flag | Description | Default |
| --- | --- | --- |
| collector.qdisc.device-exclude | Regexp of qdisc devices to exclude (mutually exclusive to device-include). |  |
| collector.qdisc.device-include | Regexp of qdisc devices to include (mutually exclusive to device-exclude). |  |
| collector.qdisc.fixtures | test fixtures to use for qdisc collector end-to-end testing |  |
| collector.qdisk.device-exclude | DEPRECATED: Use collector.qdisc.device-exclude |  |
| collector.qdisk.device-include | DEPRECATED: Use collector.qdisc.device-include |  |

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_qdisc_backlog | Number of bytes currently in queue to be sent. | device, kind |
| node_qdisc_bytes_total | Number of bytes sent. | device, kind |
| node_qdisc_current_queue_length | Number of packets currently in queue to be sent. | device, kind |
| node_qdisc_drops_total | Number of packets dropped. | device, kind |
| node_qdisc_overlimits_total | Number of overlimit packets. | device, kind |
| node_qdisc_packets_total | Number of packets sent. | device, kind |
| node_qdisc_requeues_total | Number of packets dequeued, not transmitted, and requeued. | device, kind |
