# InfiniBand collector — AWS EFA support

This document describes the EFA (AWS Elastic Fabric Adapter) extension added
to the `infiniband` collector in the Baseten fork of `node_exporter`, and how
to test and validate it.

## Background — why the patch exists

EFA NICs register under `/sys/class/infiniband/<device>` like InfiniBand HCAs,
but they **do not** follow the IB spec for port byte/packet counters:

| Counter | Mellanox IB (`mlx5_*`) | AWS EFA (`rdmap*`) |
|---|---|---|
| TX bytes | `counters/port_xmit_data` (4-octet words per IB spec) | `hw_counters/tx_bytes` (raw bytes) |
| RX bytes | `counters/port_rcv_data` | `hw_counters/rx_bytes` |
| TX packets | `counters/port_xmit_packets` | `hw_counters/tx_pkts` |
| RX packets | `counters/port_rcv_packets` | `hw_counters/rx_pkts` |

Upstream `node_exporter` only reads the `counters/` path, so on EFA-equipped
hosts the IB-equivalent metrics report zero or are completely absent — the
collector reports the device exists (`state_id`, `rate_bytes_per_second`) but
shows no traffic. This silently breaks any dashboard that relies on
`node_infiniband_port_data_*_bytes_total` to monitor cross-host RDMA on AWS
p4d/p5/p6/g5.48xl/trn1 instances.

The patch teaches the collector to detect EFA NICs by PCI vendor ID
(`0x1d0f`, assigned by PCI-SIG to AWS) and route them through `hw_counters/`
for the affected counters. The metric names are reused, so existing IB
dashboards transparently start showing EFA throughput without query changes.

## Metric semantics

### Shared IB/EFA metric names

Both InfiniBand and EFA devices populate the following metrics; values come
from whichever sysfs path matches the device family:

- `node_infiniband_port_data_transmitted_bytes_total`
- `node_infiniband_port_data_received_bytes_total`
- `node_infiniband_port_packets_transmitted_total`
- `node_infiniband_port_packets_received_total`
- `node_infiniband_state_id` (4 = ACTIVE)
- `node_infiniband_physical_state_id` (5 = LinkUp)
- `node_infiniband_rate_bytes_per_second` (negotiated link rate)

The `device` label distinguishes them: `mlx5_*` / `ibp*` for Mellanox,
`rdmap*` for EFA.

### EFA-only diagnostic metrics

EFA exposes counters that have no IB-spec equivalent. These are emitted under
the `efa_` prefix to keep IB semantics clean and let dashboards opt in:

| Metric | Source file | Use |
|---|---|---|
| `node_infiniband_efa_rx_drops_total` | `hw_counters/rx_drops` | Inbound packet loss at the NIC |
| `node_infiniband_efa_retrans_packets_total` | `hw_counters/retrans_pkts` | Fabric-level retransmissions |
| `node_infiniband_efa_retrans_bytes_total` | `hw_counters/retrans_bytes` | Retransmitted bytes |
| `node_infiniband_efa_retrans_timeout_events_total` | `hw_counters/retrans_timeout_events` | Retransmit timeouts |
| `node_infiniband_efa_unresponsive_remote_events_total` | `hw_counters/unresponsive_remote_events` | Remote peer not responding |
| `node_infiniband_efa_impaired_remote_conn_events_total` | `hw_counters/impaired_remote_conn_events` | Degraded peer connection |
| `node_infiniband_efa_rdma_read_bytes_total` | `hw_counters/rdma_read_bytes` | RDMA Read bytes |
| `node_infiniband_efa_rdma_write_bytes_total` | `hw_counters/rdma_write_bytes` | RDMA Write bytes |

Sustained non-zero `retrans_*` or `rx_drops` typically indicates EFA fabric
congestion or unhealthy peers — a strong signal for NCCL/training perf
regressions.

## Behavior summary

- A device with `device/vendor` reading `0x1d0f` follows the EFA code path.
  All other devices (including Mellanox `0x15b3` or devices missing
  `device/vendor` entirely) follow the original IB code path with no change
  in behavior.
- IB error counters that have no EFA equivalent (`link_downed_total`,
  `symbol_error_total`, `port_discards_*`, `port_transmit_wait_total`,
  multicast/legacy counters, etc.) are **skipped** for EFA devices — they
  would be zero/empty anyway.
- The EFA detection is performed once per device, not per port, so the
  per-scrape overhead is one extra `stat`+`read` of a small file per IB
  device.
- A missing `hw_counters/` file silently emits no series for that counter,
  matching existing `pushCounter` behavior.

## Mapping to upstream procfs

The new logic lives in `collector/infiniband_linux.go`. The upstream
`prometheus/procfs` library is **not** modified; it still returns the IB-spec
view of each device. The collector now branches on EFA detection and reads
`hw_counters/` directly via `os.ReadFile` when needed. This isolates the
patch to a single file and avoids forking a transitive dependency.

## Running the unit tests

The patch ships with 22 test cases in
`collector/infiniband_linux_test.go` covering helper functions and
end-to-end behavior through a mock sysfs tree.

### Quick check — only the EFA tests

```bash
cd /workspace/node_exporter
go test ./collector/ -run 'EFA|TestUpdate_' -v
```

Expected: 22 PASS in ~20ms.

### All collector tests (sanity before push)

```bash
go test ./collector/ -v -race
```

### Single test case

```bash
go test -v ./collector/ -run '^TestUpdate_EFAReadsHWCounters$'
```

