# Correlating Exported Metrics: Kernel Stack Cost per Packet and AF_PACKET

**Diagram:** For a visual overview of the new features and how they interact with the Linux system, kernel, eBPF program, interfaces, DPDK, and hugepages, open [node-exporter-new-features.excalidraw.json](node-exporter-new-features.excalidraw.json) in [Excalidraw](https://excalidraw.com) (or the VS Code Excalidraw extension).

This document ties together the node_exporter metrics that help you **detect and quantify** high cost per packet in the Linux kernel network stack, and to **correlate** that cost with AF_PACKET (and similar) usage. It explains how the exported metrics relate to each other and how to conclude that AF_PACKET (or the kernel packet path) is increasing cost per packet.

---

## 1. Why “cost per packet” and AF_PACKET matter

### 1.1 Cost per packet in the kernel stack

Every packet that goes through the kernel network stack consumes CPU and time:

- **Interrupt / NAPI** → **netif_receive_skb** → protocol handling → **delivery to sockets** (e.g. AF_PACKET, AF_INET) → **kfree_skb** / **consume_skb**.

The **time from driver (or XDP) entry until the packet leaves the stack** is the “kernel stack latency” for that packet. The **CPU time spent per packet** (e.g. in softirq, or in kernel system time) is the “cost per packet.” When that cost or latency goes up, throughput drops and latency increases.

### 1.2 How AF_PACKET increases cost per packet

**AF_PACKET** (and similar mechanisms that receive a copy of every packet in the kernel) increase cost per packet because:

1. **Extra work per packet**: For each packet, the kernel may fan out copies to one or more AF_PACKET sockets (`packet_rcv()` and related paths), then run `consume_skb()` (or equivalent) for each copy.
2. **More CPU in softirq**: That work usually runs in **softirq** (e.g. **NET_RX**). So when many AF_PACKET sockets are open (e.g. `tcpdump`, capture tools), you see:
   - **Higher softirq time** on the CPUs handling that traffic.
   - **Longer time per packet in the stack** (higher “kernel stack latency”).
   - **Higher jitter** (variance in that latency).
3. **Impact on reserved/isolated CPUs**: If the same CPUs are used for both AF_PACKET and for DPDK or other latency-sensitive workloads, those CPUs can hit 90–100% softirq and you get throughput drops (e.g. 30–70%), connection issues, and higher latency.

So: **AF_PACKET (and similar kernel copy-to-userspace paths) are one major way the kernel stack’s cost per packet is increased.** The metrics below let you observe that increase and correlate it with softirq, CPU isolation, and throughput.

---

## 2. What the eBPF “kernel stack latency” actually measures

The **ebpf-pmd-jitter** collector uses an in-tree eBPF program (`collector/bpf/latency.c`) that measures **time spent in the kernel network stack** for each packet:

| Stage | Hook | Meaning |
|--------|------|--------|
| **Start** | **XDP** (`xdp_latency_ingress`) | Packet timestamp when it first enters the stack (driver / XDP). |
| **End** | **TC ingress** or **TC egress** | Packet timestamp when it reaches TC (after protocol and socket handling). |

So:

- **Latency** = `(time at TC) - (time at XDP)` = time the packet spent in the kernel path between XDP and TC (NAPI, netif_receive_skb, protocol, **AF_PACKET/socket delivery**, etc.).
- **Jitter** = variance of that latency (e.g. max − min over a window).

When AF_PACKET (or any extra per-packet work) increases, this **kernel stack latency** and **jitter** tend to increase. So these metrics are a **direct indicator of “cost per packet” in the stack**, and help you conclude that the kernel path (including AF_PACKET) is one way that cost is increased.

---

## 3. Metric catalog: what is exported and how it correlates

All metric names below are as exported by node_exporter when the corresponding collector is enabled. **Coverage at `localhost:9192/metrics`:** metrics appear only if the collector that produces them is enabled and (where applicable) the kernel or filesystem provides the data (e.g. conntrack metrics only if the `nf_conntrack` module is loaded). Default-enabled collectors (cpu, meminfo, netdev, netstat, sockstat, conntrack when loaded, pressure, etc.) are always present on a typical run; **nodeconfig**, **ebpf-pmd-jitter**, **zoneinfo**, **meminfo_numa**, and **pcidevice** are disabled by default—enable them as in §6 to see those metrics at `:9192/metrics`. Everything in this document works as expected when the listed collectors are enabled and the eBPF object/BPF fs are set up for ebpf-pmd-jitter.

### 3.1 Kernel stack cost and jitter (eBPF)

| Metric | Description | Correlation to “cost per packet” / AF_PACKET |
|--------|-------------|----------------------------------------------|
| `node_ebpf_pmd_jitter_latency_min_ns` | Min time (ns) packet in stack (XDP→TC) per interface. | Baseline; rises when stack does more work per packet (e.g. AF_PACKET). |
| `node_ebpf_pmd_jitter_latency_max_ns` | Max time (ns) packet in stack per interface. | Spikes when some packets pay extra cost (e.g. copy to many sockets). |
| `node_ebpf_pmd_jitter_latency_avg_ns` | Average time (ns) in stack per interface. | Direct measure of average cost per packet in the kernel path. |
| `node_ebpf_pmd_jitter_pmd_jitter_ns` | Jitter (max − min) per interface. | High jitter often goes with AF_PACKET/capture (variable extra work per packet). |
| `node_ebpf_pmd_jitter_latency_histogram` | Distribution of latency by bucket (0–1µs … 16ms+). | Shift to higher buckets = more packets with high stack cost. |
| `node_ebpf_pmd_jitter_global_packets_total` | Total packets measured globally. | Volume; combine with latency to reason about total cost. |
| `node_ebpf_pmd_jitter_global_latency_ns_total` | Total latency (ns) globally. | Total “cost” in time; divide by packets for average cost per packet. |
| `node_ebpf_pmd_jitter_collector_up` | 1 if eBPF loaded/attached. | Must be 1 for above metrics to be meaningful. |
| `node_ebpf_pmd_jitter_object_path_configured` | 1 if object-path set. | With `collector_up=0`, distinguishes “not configured” vs “load failed”. |
| `node_ebpf_pmd_jitter_load_error` | 1 on load/attach failure; `error` label has reason. | Debug why kernel stack metrics are missing. |

**Collector:** `ebpf-pmd-jitter` (optional). Requires eBPF object and `/sys/fs/bpf` mounted.

### 3.2 Node configuration and isolation (runbook context)

| Metric | Description | Correlation to “cost per packet” / AF_PACKET |
|--------|-------------|----------------------------------------------|
| `node_nodeconfig_cores_dedicated` | 1 if any CPU in `/sys/.../cpu/isolated`. | Indicates reserved/isolated CPUs; high softirq here + AF_PACKET is a common failure mode. |
| `node_nodeconfig_pcie_nic_min_link_width` | Min PCIe link width among NICs. | Runbook context (e.g. slot correctness). |
| `node_nodeconfig_pcie_slot_ok` | 1 if min NIC width ≥ 16. | Runbook context. |
| `node_nodeconfig_memory_banks_full` | 1 if DMI says all memory slots populated. | Runbook context. |

**Collector:** `nodeconfig` (optional).

### 3.3 CPU and softirq (where the cost shows up)

| Metric | Description | Correlation to “cost per packet” / AF_PACKET |
|--------|-------------|----------------------------------------------|
| `node_cpu_seconds_total{cpu="...", mode="softirq"}` | CPU time in softirq per CPU. | **Primary signal**: AF_PACKET work runs in softirq; high rate = high cost. |
| `node_cpu_seconds_total{mode="system"}` | Kernel (system) time. | Kernel stack and AF_PACKET also show up here. |
| `node_cpu_isolated{cpu="..."}` | 1 if that CPU is isolated. | Isolated CPUs often used for DPDK/capture; combine with softirq to find overload. |

**Collector:** `cpu` (default).

### 3.4 Softirq breakdown (NET_RX and others)

| Metric | Description | Correlation to “cost per packet” / AF_PACKET |
|--------|-------------|----------------------------------------------|
| `node_softirqs_functions_total{cpu="...", type="NET_RX"}` | NET_RX softirq count per CPU. | **NET_RX** dominates when copying to many AF_PACKET sockets; compare before/during capture. |

**Collector:** `softirqs` (optional).

### 3.5 Network throughput and drops

| Metric | Description | Correlation to “cost per packet” / AF_PACKET |
|--------|-------------|----------------------------------------------|
| `node_network_receive_bytes_total`, `node_network_transmit_bytes_total` | Bytes per interface. | Throughput; drop in rate can correlate with high stack cost / AF_PACKET. |
| `node_network_receive_drop_total`, `node_network_receive_errs_total` | Drops/errors per interface. | Can increase when CPUs are overloaded by softirq/AF_PACKET. |

**Collector:** `netdev` (default).

### 3.6 Hugepages (runbook / DPDK context)

Hugepages are **not** exposed by the new nodeconfig or ebpf-pmd-jitter collectors. They come from existing **meminfo** (and optionally **meminfo_numa**, **zoneinfo**) collectors, which read `/proc/meminfo` (and NUMA/zone stats).

| Metric | Description | Correlation to “cost per packet” / AF_PACKET |
|--------|-------------|----------------------------------------------|
| `node_memory_HugePages_Total` | Total number of huge pages configured. | Runbook: DPDK and high-throughput packet paths often use hugepages; non-zero indicates hugepage pool is configured. |
| `node_memory_HugePages_Free` | Free huge pages. | Runbook: low free count can indicate DPDK/other apps consuming the pool. |
| `node_memory_HugePages_Rsvd`, `node_memory_HugePages_Surp` | Reserved and surplus huge pages. | Runbook context for capacity. |
| `node_memory_Hugepagesize_bytes` | Size of one huge page (e.g. 2 MiB). | Runbook context. |
| `node_memory_AnonHugePages_bytes` | Anonymous memory backed by transparent huge pages. | General memory use; not specific to DPDK. |
| `node_memory_numa_HugePages_*` (with label `node`) | Per-NUMA-node huge page counts (meminfo_numa). | Runbook: when correlating with isolated CPUs, NUMA locality and hugepage usage matter for DPDK/PMD. |
| `node_zoneinfo_nr_anon_transparent_hugepages` | Transparent huge pages (zoneinfo). | General THP usage. |

**How hugepages factor in:** They do **not** measure kernel stack cost or AF_PACKET directly. They provide **runbook and environment context** on nodes where you are correlating AF_PACKET vs DPDK:

- On **DPDK nodes**, PMD and packet buffers often use **hugepages**. If you see high softirq on isolated CPUs (possible AF_PACKET/capture) and you expect DPDK to be using hugepages, checking `node_memory_HugePages_Total` / `_Free` confirms hugepage configuration and usage; together with `node_nodeconfig_cores_dedicated` and `node_cpu_isolated`, you have a fuller picture (isolated CPUs, PCIe slot, memory banks, hugepages) for the same runbook.
- So: use hugepage metrics as **context** alongside nodeconfig (cores_dedicated, pcie_slot_ok, memory_banks_full) when reasoning about DPDK vs kernel path (AF_PACKET) and cost per packet.

**Collectors:** `meminfo` (default), `meminfo_numa` (optional), `zoneinfo` (default).

### 3.7 Conntrack (nf_conntrack) — connection tracking table

When the kernel connection tracking table is full, **new flows cannot be established** and packets are dropped. This is a common cause of “random” packet loss and connection failures under load.

| Metric | Description | Correlation |
|--------|-------------|-------------|
| `node_nf_conntrack_entries` | Current number of allocated conntrack flow entries. | Nearing limit → table filling; correlate with drops. |
| `node_nf_conntrack_entries_limit` | Maximum size of the conntrack table (`nf_conntrack_max`). | When `entries` ≈ `entries_limit`, new entries fail. |
| `node_nf_conntrack_stat_found` | Number of successful conntrack lookups. | Normal operation. |
| `node_nf_conntrack_stat_invalid` | Packets that could not be tracked. | Can increase with bad packets or table pressure. |
| `node_nf_conntrack_stat_ignore` | Packets already connected (existing entry). | Normal. |
| `node_nf_conntrack_stat_insert` | New entries inserted. | Insert rate vs limit drives fullness. |
| `node_nf_conntrack_stat_insert_failed` | Insert attempts that failed (e.g. duplicate). | Correlate with table pressure. |
| `node_nf_conntrack_stat_drop` | **Packets dropped due to conntrack failure** (allocation or helper). | **Primary drop signal** when conntrack is the cause. |
| `node_nf_conntrack_stat_early_drop` | **Entries dropped to make room** (table was full). | **Direct signal** that the table hit max; packets were dropped. |
| `node_nf_conntrack_stat_search_restart` | Lookups restarted due to hashtable resize. | High rate can indicate churn. |

**Collector:** `conntrack` (default when conntrack module loaded). Source: `/proc/sys/net/netfilter/nf_conntrack_count`, `nf_conntrack_max`, `/proc/net/nf_conntrack` stats.

### 3.8 Netstat — TCP/UDP and buffer-related counters

From `/proc/net/netstat` and `/proc/net/snmp`. The default netstat collector exposes a subset of fields (see `--collector.netstat.fields`). Key metrics for TCP buffers, listen queue, and drops:

| Metric | Description | Correlation |
|--------|-------------|-------------|
| `node_netstat_Tcp_CurrEstab` | Currently established TCP connections. | Connection count; high with many flows. |
| `node_netstat_Tcp_ActiveOpens`, `node_netstat_Tcp_PassiveOpens` | Open attempts. | Connection churn. |
| `node_netstat_Tcp_InSegs`, `node_netstat_Tcp_OutSegs` | Segments in/out. | Traffic volume. |
| `node_netstat_Tcp_RetransSegs` | Retransmitted segments. | Loss or latency; can rise with buffer pressure. |
| `node_netstat_TcpExt_TCPRcvQDrop` | **Segments dropped because receive queue was full**. | **TCP receive buffer full** → application not reading fast enough or buffer too small. |
| `node_netstat_TcpExt_ListenOverflows` | **Times the listen queue overflowed** (SYN queue full). | **Listen backlog full** → new connections dropped. |
| `node_netstat_TcpExt_ListenDrops` | **Times a SYN was dropped** (e.g. backlog full). | **Direct listen-queue drop** signal. |
| `node_netstat_TcpExt_TCPTimeouts` | TCP timeouts. | Can increase with loss or congestion. |
| `node_netstat_TcpExt_TCPOFOQueue` | Out-of-order queue length. | Reordering / burst. |
| `node_netstat_Udp_RcvbufErrors` | **UDP receive buffer errors** (drops). | **UDP socket receive buffer full**. |
| `node_netstat_Udp_SndbufErrors` | **UDP send buffer errors** (drops). | **UDP socket send buffer full**. |
| `node_netstat_Udp6_RcvbufErrors`, `node_netstat_Udp6_SndbufErrors` | Same for IPv6. | Same interpretation. |

**Collector:** `netstat` (default). Source: `/proc/net/netstat`, `/proc/net/snmp`, `/proc/net/snmp6`.

### 3.9 Sockstat — socket memory and usage

Socket layer memory and in-use socket counts. Relevant for buffer and connection scaling.

| Metric | Description | Correlation |
|--------|-------------|-------------|
| `node_sockstat_sockets_used` | Number of IPv4 sockets in use. | Total socket usage. |
| `node_sockstat_TCP_inuse` | TCP sockets in use. | TCP connection count (similar to CurrEstab but from sock layer). |
| `node_sockstat_TCP_orphan` | Orphaned TCP (no user ref). | Can grow under load. |
| `node_sockstat_TCP_tw` | TIME_WAIT sockets. | High with many short-lived connections. |
| `node_sockstat_TCP_alloc` | TCP sockets allocated. | Allocation count. |
| `node_sockstat_TCP_mem` | **TCP socket memory (pages)**. | **TCP buffer memory**; high = many/big buffers. |
| `node_sockstat_TCP_mem_bytes` | **TCP socket memory in bytes**. | **Direct TCP buffer memory** (mem × page size). |
| `node_sockstat_UDP_inuse`, `node_sockstat_UDP_mem`, `node_sockstat_UDP_mem_bytes` | Same for UDP. | UDP buffer and usage. |

**Collector:** `sockstat` (default). Source: `/proc/net/sockstat`, `/proc/net/sockstat6`.

### 3.10 NUMA — locality and cross-node access

When the kernel or userspace (e.g. OVS, DPDK) allocates memory on a “wrong” NUMA node (e.g. NIC on node 0, process on node 1), **latency and cost per packet increase** due to remote memory access. These metrics help detect NUMA-unfriendly placement.

| Metric | Description | Correlation |
|--------|-------------|-------------|
| `node_zoneinfo_numa_hit_total{node="...", zone="..."}` | Allocations satisfied from the intended node. | Local allocations. |
| `node_zoneinfo_numa_miss_total{node="...", zone="..."}` | **Allocations satisfied from another node** (remote). | **High miss** → NUMA-unfriendly; higher latency. |
| `node_zoneinfo_numa_foreign_total` | “Intended here, hit elsewhere.” | Another view of remote allocation. |
| `node_zoneinfo_numa_local_total` | Allocations from local node. | Local. |
| `node_zoneinfo_numa_other_total` | Allocations from other node. | **High other** → cross-node access; bad for latency. |
| `node_memory_numa_*{node="..."}` | Per-NUMA memory stats (meminfo_numa). | Memory usage per node; pair with NIC/process placement. |
| `node_pcidevice_numa_node` | NUMA node of PCI device (e.g. NIC). | Which node the NIC is on; compare to process/OVS NUMA. |

**Collectors:** `zoneinfo` (optional), `meminfo_numa` (optional), `pcidevice` (optional). Sources: `/proc/zoneinfo`, sysfs NUMA meminfo/numastat, PCI sysfs.

### 3.11 Pressure — memory and I/O stall

PSI (Pressure Stall Information) indicates when workloads are stalled due to memory or I/O. Can correlate with buffer pressure and cost per packet.

| Metric | Description | Correlation |
|--------|-------------|-------------|
| `node_pressure_memory_stalled_seconds_total` | Time no process could make progress due to memory. | Memory pressure; can accompany buffer/socket pressure. |
| `node_pressure_memory_waiting_seconds_total` | Time processes waited for memory. | Same. |
| `node_pressure_io_stalled_seconds_total`, `node_pressure_io_waiting_seconds_total` | I/O pressure. | Disk/network stack contention. |
| `node_pressure_cpu_waiting_seconds_total` | CPU pressure. | CPU saturation. |

**Collector:** `pressure` (default). Source: `/proc/pressure/*`.

---

## 4. How the metrics correlate: decision flow

Use this flow to **correlate** metrics and **support the conclusion** that AF_PACKET (or the kernel packet path) is increasing cost per packet:

1. **Kernel stack cost (eBPF)**  
   - Check `node_ebpf_pmd_jitter_latency_avg_ns`, `node_ebpf_pmd_jitter_pmd_jitter_ns`, and `node_ebpf_pmd_jitter_latency_histogram`.  
   - **Rising** average latency or jitter, or histogram shifting to higher buckets → **higher cost per packet in the kernel stack**.

2. **Where the cost appears (CPU)**  
   - Check `rate(node_cpu_seconds_total{mode="softirq"}[5m])` and, if needed, `node_softirqs_functions_total{type="NET_RX"}`.  
   - **High softirq share** (e.g. > 0.9) and/or **high NET_RX** on the same CPUs that handle the traffic → cost is in the **kernel receive path** (where AF_PACKET runs).

3. **Context (isolation and throughput)**  
   - `node_nodeconfig_cores_dedicated == 1` and `node_cpu_isolated{cpu="X"} == 1`: reserved CPUs in use.  
   - If those CPUs are the ones with high softirq and high kernel stack latency → typical “reserved CPUs hammered by kernel path” pattern.  
   - `node_network_*` throughput down and/or drops up **together with** high softirq and high eBPF latency → **consistent with** high cost per packet (e.g. AF_PACKET) degrading performance.

4. **Conclusion**  
   - **High kernel stack latency/jitter (eBPF) + high softirq (and optionally NET_RX) + throughput degradation (and optionally drops)** → the **kernel stack** (and mechanisms like **AF_PACKET** that add per-packet work in that path) is **one way** cost per packet is increased.  
   - You can then correlate with known AF_PACKET use (tcpdump, capture tools) or with tooling (e.g. `/proc/net/packet`, tracing) to confirm AF_PACKET specifically.

5. **Optional: runbook context (nodeconfig, hugepages)**  
   - On nodes that may run **DPDK** or other high-throughput packet paths: use `node_nodeconfig_cores_dedicated`, `node_nodeconfig_pcie_slot_ok`, and **hugepage metrics** (`node_memory_HugePages_Total`, `node_memory_HugePages_Free`, etc.) as **context**—they do not measure cost per packet but help confirm isolation, NIC slot, and hugepage configuration when you are deciding between “kernel path (AF_PACKET) overload” vs “DPDK/PMD setup.”

6. **Other drop and latency causes (see §7)**  
   - **Conntrack full (Example A):** `node_nf_conntrack_entries` ≈ `node_nf_conntrack_entries_limit` and `node_nf_conntrack_stat_early_drop` or `_stat_drop` increasing → packet drops from full connection tracking table.  
   - **TCP/UDP buffer full (Example C):** `node_netstat_TcpExt_TCPRcvQDrop`, `ListenOverflows`, or `node_netstat_Udp_*bufErrors` increasing → drops from socket buffers; can be exacerbated by high cost per packet (AF_PACKET).  
   - **NUMA-unfriendly path (Example B):** high `node_zoneinfo_numa_miss_total` (or numa_other) + high eBPF latency → kernel/OVS path paying remote-memory cost per packet.

---

## 5. PromQL examples

### 5.1 Average kernel stack latency (when eBPF is used)

```promql
# Per-interface average (only when collector_up == 1)
node_ebpf_pmd_jitter_latency_avg_ns
  and on() node_ebpf_pmd_jitter_collector_up == 1
```

### 5.2 Jitter (max − min) per interface

```promql
node_ebpf_pmd_jitter_pmd_jitter_ns
  and on() node_ebpf_pmd_jitter_collector_up == 1
```

### 5.3 Softirq share of CPU (all CPUs)

```promql
rate(node_cpu_seconds_total{mode="softirq"}[5m])
  / ignoring(mode) group_left()
rate(node_cpu_seconds_total[5m])
```

### 5.4 High softirq on isolated CPUs (alert pattern)

```promql
(
  rate(node_cpu_seconds_total{mode="softirq"}[5m])
  / on(cpu) group_left()
  rate(node_cpu_seconds_total[5m])
) > 0.9
and on(cpu) node_cpu_isolated == 1
```

### 5.5 Correlation: high kernel stack latency and high softirq

```promql
# Example: nodes where both kernel stack avg latency and softirq share are high
(
  node_ebpf_pmd_jitter_latency_avg_ns > 1000
  and on() node_ebpf_pmd_jitter_collector_up == 1
)
and on(instance) (
  rate(node_cpu_seconds_total{mode="softirq"}[5m])
  / ignoring(mode) group_left()
  rate(node_cpu_seconds_total[5m])
) > 0.7
```

(Adjust thresholds and grouping to match your labels and SLOs.)

### 5.6 Conntrack table near full (Example A)

```promql
# Usage ratio: entries / limit (alert when > 0.9)
node_nf_conntrack_entries / node_nf_conntrack_entries_limit > 0.9
```

```promql
# Any conntrack drop indicates table pressure or allocation failure
increase(node_nf_conntrack_stat_drop[5m]) > 0
or
increase(node_nf_conntrack_stat_early_drop[5m]) > 0
```

### 5.7 TCP receive buffer drops (Example C)

```promql
# TCP receive queue drops (segments dropped due to full rcvbuf)
increase(node_netstat_TcpExt_TCPRcvQDrop[5m]) > 0
```

```promql
# Listen queue overflows (SYN backlog full)
increase(node_netstat_TcpExt_ListenOverflows[5m]) > 0
```

### 5.8 NUMA miss ratio (Example B)

```promql
# Per-node NUMA miss rate (high = allocations served from wrong node)
rate(node_zoneinfo_numa_miss_total[5m])
  / (rate(node_zoneinfo_numa_hit_total[5m]) + rate(node_zoneinfo_numa_miss_total[5m]))
```

---

## 6. Collectors to enable

| Collector | Default | Enable with | Purpose |
|-----------|---------|-------------|---------|
| conntrack | on | (default, if loaded) | **Example A**: `node_nf_conntrack_entries`, `_entries_limit`, `_stat_drop`, `_stat_early_drop`. |
| cpu | on | (default) | Softirq and system time, `node_cpu_isolated`. |
| meminfo | on | (default) | **Hugepages**: `node_memory_HugePages_*`, etc. Runbook/DPDK context. |
| meminfo_numa | off | `--collector.meminfo_numa` | Per-NUMA memory; **Example B** (NUMA/OVS). |
| netdev | on | (default) | Throughput and drops: `node_network_*`. |
| netstat | on | (default) | **Example C**: `node_netstat_TcpExt_TCPRcvQDrop`, `ListenOverflows`, `Udp_*bufErrors`. |
| nodeconfig | off | `--collector.nodeconfig` | `node_nodeconfig_cores_dedicated`, PCIe/memory context. |
| pcidevice | off | `--collector.pcidevice` | **Example B**: `node_pcidevice_numa_node` (NIC NUMA). |
| pressure | on | (default) | Memory/IO stall; can correlate with buffer pressure. |
| softirqs | off | `--collector.softirqs` | NET_RX and other softirq vectors. |
| sockstat | on | (default) | **Example C**: `node_sockstat_TCP_mem_bytes`, socket usage. |
| zoneinfo | off | `--collector.zoneinfo` | **Example B**: `node_zoneinfo_numa_*` (NUMA hit/miss). |
| ebpf-pmd-jitter | off | `--collector.ebpf-pmd-jitter --collector.ebpf-pmd-jitter.object-path=...` | Kernel stack latency and jitter; requires eBPF object and `/sys/fs/bpf`. |

---

## 7. Detailed examples: conntrack drops, NUMA latency, TCP buffers

The following three examples show how to use the exported metrics at `localhost:9192/metrics` (or your node_exporter endpoint) to diagnose specific failure modes. Each example ties back to **cost per packet** and, where relevant, **AF_PACKET** (or kernel path) impact.

---

### Example A: Full nf_conntrack table leads to packet drops

**What happens:** The kernel connection tracking table has a fixed maximum (`nf_conntrack_max`). Every new flow (e.g. new TCP connection, UDP flow, or NAT session) allocates an entry. When the table is full:

1. New entries cannot be allocated.
2. The kernel **drops packets** that would create new flows.
3. You see `node_nf_conntrack_stat_drop` and/or `node_nf_conntrack_stat_early_drop` increase; established flows may keep working while new connections fail.

**Exported metrics that matter:**

| Metric | How to use it |
|--------|----------------|
| `node_nf_conntrack_entries` | Current number of conntrack entries. If it stays near the limit, the table is full. |
| `node_nf_conntrack_entries_limit` | Maximum table size. Compare to `entries`; ratio close to 1.0 = full. |
| `node_nf_conntrack_stat_early_drop` | **Increments when the kernel drops existing entries to make room** (table was full). **Primary counter** for “table full” drops. |
| `node_nf_conntrack_stat_drop` | **Increments when packets are dropped due to conntrack failure** (allocation failed or protocol helper). Includes both “no room” and other conntrack failures. |
| `node_nf_conntrack_stat_insert` | New entries inserted. High rate with high `entries` → table fills quickly. |
| `node_network_receive_drop_total`, `node_network_receive_errs_total` | Interface-level drops. Can rise when conntrack drops packets before they are “received” by a socket; correlate with conntrack drop spikes. |

**Decision flow:**

1. Check **usage ratio**: `node_nf_conntrack_entries / node_nf_conntrack_entries_limit`. If consistently > 0.85–0.95, the table is under heavy pressure.
2. Check **drop counters**: `increase(node_nf_conntrack_stat_early_drop[5m]) > 0` or `increase(node_nf_conntrack_stat_drop[5m]) > 0`. Any increase confirms conntrack-related drops.
3. Correlate in time with **symptoms**: new connection failures, timeouts, or interface drops (`node_network_receive_drop_total`).
4. **Conclusion:** Full conntrack table → packet drops for new flows. Remediation: increase `nf_conntrack_max` (and often `nf_conntrack_buckets`), or reduce connection churn / shorten timeouts.

**PromQL (alerts):**

```promql
# Alert when table usage > 90%
(node_nf_conntrack_entries / node_nf_conntrack_entries_limit) > 0.9

# Alert when any early_drop (table was full)
increase(node_nf_conntrack_stat_early_drop[5m]) > 0
```

**Relation to AF_PACKET / cost per packet:** Conntrack is in the same kernel path as AF_PACKET. When the table is full, **every packet that would create a new flow pays the cost of a failed conntrack lookup/insert and then drop**—so CPU is still spent (cost per packet) but the packet is dropped. High conntrack churn can also increase softirq and latency (more lookups, resizes); correlate with `node_cpu_seconds_total{mode="softirq"}` and eBPF latency if available.

---

### Example B: High packet latency caused by OVS (or kernel path) not NUMA-aware optimized

**What happens:** On NUMA systems, memory access is faster when the CPU and the memory (or device) are on the same NUMA node. If Open vSwitch (OVS) or the kernel datapath allocates buffers or runs work on a different node than the NIC (or the other end of the path), you get **remote memory access** → higher latency and higher cost per packet. OVS not bound to the NIC’s NUMA node is a common cause.

**Exported metrics that matter:**

| Metric | How to use it |
|--------|----------------|
| `node_zoneinfo_numa_hit_total{node="N", zone="Z"}` | Allocations satisfied from the intended node (local). |
| `node_zoneinfo_numa_miss_total{node="N", zone="Z"}` | **Allocations satisfied from another node (remote).** High rate here (or high miss ratio) = many allocations served from the “wrong” node. |
| `node_zoneinfo_numa_other_total` | Allocations from a non-local node. High = cross-node access. |
| `node_zoneinfo_numa_local_total` | Local allocations. Compare to `numa_other_total`; if “other” is a large fraction, workload is not NUMA-local. |
| `node_memory_numa_*{node="..."}` (meminfo_numa) | Per-node memory usage. Helps see which node the process/OVS is using. |
| `node_pcidevice_numa_node` | **NUMA node of each PCI device (e.g. NIC).** Compare to the NUMA node of the process handling the traffic; mismatch suggests non-local path. |
| `node_ebpf_pmd_jitter_latency_avg_ns`, `node_ebpf_pmd_jitter_pmd_jitter_ns` | **Kernel stack latency and jitter.** When NUMA is wrong, latency and jitter often go up (remote access adds delay and variance). |
| `node_cpu_seconds_total{mode="softirq"}`, `node_cpu_isolated` | Which CPUs are busy with softirq; if those CPUs are on a different node than the NIC, you have a NUMA mismatch. |

**Decision flow:**

1. Identify **NIC NUMA node**: use `node_pcidevice_numa_node` for the relevant NIC(s).
2. Check **NUMA allocation quality**: compute miss ratio per node, e.g. `rate(node_zoneinfo_numa_miss_total[5m]) / (rate(node_zoneinfo_numa_hit_total[5m]) + rate(node_zoneinfo_numa_miss_total[5m]))`. High ratio on the node(s) handling packet I/O → allocations often served from remote node.
3. Correlate with **latency**: if eBPF is available, `node_ebpf_pmd_jitter_latency_avg_ns` and `node_ebpf_pmd_jitter_pmd_jitter_ns` rising together with high NUMA miss → kernel path (including OVS if in kernel) is paying remote-memory cost per packet.
4. **Conclusion:** High NUMA miss (or high “other”) + high kernel stack latency/jitter → packet path is not NUMA-optimized; consider binding OVS/process to the NIC’s NUMA node and ensuring hugepage/memory allocation is local.

**PromQL (NUMA miss ratio per node):**

```promql
# Per-node NUMA miss ratio (by zone)
sum by (node) (rate(node_zoneinfo_numa_miss_total[5m]))
  / (sum by (node) (rate(node_zoneinfo_numa_hit_total[5m]))
      + sum by (node) (rate(node_zoneinfo_numa_miss_total[5m])))
```

**Relation to AF_PACKET / cost per packet:** Both “OVS not NUMA-aware” and “AF_PACKET on wrong node” increase **cost per packet** (extra cycles for remote access or extra copy). The eBPF latency metrics reflect the combined effect: higher latency and jitter when the path is not NUMA-local.

---

### Example C: TCP buffer monitoring and cost

**What happens:** TCP and UDP socket buffers are finite. If the application does not read (or send) fast enough, the kernel **drops segments** or **refuses new connections** (listen queue). This shows up as TCP receive-queue drops, listen overflows, or UDP buffer errors. Buffer pressure can also correlate with high CPU (e.g. softirq) or AF_PACKET: if the kernel is busy copying to many AF_PACKET sockets, it may not drain TCP buffers in time → more drops.

**Exported metrics that matter:**

| Metric | How to use it |
|--------|----------------|
| `node_netstat_TcpExt_TCPRcvQDrop` | **Segments dropped because the TCP receive queue was full.** Application not reading fast enough or rcvbuf too small. |
| `node_netstat_TcpExt_ListenOverflows` | **Times the listen queue overflowed** (SYN backlog full). New connections dropped. |
| `node_netstat_TcpExt_ListenDrops` | **Times a SYN was dropped** (e.g. listen backlog full). |
| `node_netstat_Udp_RcvbufErrors`, `node_netstat_Udp_SndbufErrors` | **UDP receive/send buffer errors** (drops). |
| `node_netstat_Udp6_RcvbufErrors`, `node_netstat_Udp6_SndbufErrors` | Same for IPv6. |
| `node_sockstat_TCP_mem`, `node_sockstat_TCP_mem_bytes` | **Total TCP socket memory (pages/bytes).** High = many or large TCP buffers; correlate with memory pressure. |
| `node_sockstat_TCP_inuse`, `node_sockstat_sockets_used` | Socket counts. High with high TCPRcvQDrop → many connections and buffer pressure. |
| `node_netstat_Tcp_CurrEstab`, `node_netstat_Tcp_RetransSegs` | Established connections and retransmits; context for load and loss. |
| `node_pressure_memory_stalled_seconds_total`, `node_pressure_memory_waiting_seconds_total` | Memory pressure; can accompany buffer pressure. |

**Decision flow:**

1. **TCP receive queue drops:** `increase(node_netstat_TcpExt_TCPRcvQDrop[5m]) > 0`. Any increase = kernel dropped segments due to full rcvbuf.
2. **Listen queue:** `increase(node_netstat_TcpExt_ListenOverflows[5m]) > 0` or `ListenDrops` increasing → new connections dropped (backlog full).
3. **UDP buffers:** `increase(node_netstat_Udp_RcvbufErrors[5m]) > 0` or `SndbufErrors` → UDP socket buffer full.
4. **Context:** Compare with `node_sockstat_TCP_mem_bytes` (total TCP buffer memory), `node_sockstat_TCP_inuse` (connection count), and memory pressure. If TCPRcvQDrop or ListenOverflows rise while softirq or eBPF latency is high, **cost per packet** (e.g. AF_PACKET or other kernel work) may be delaying buffer processing and contributing to drops.

**PromQL (alerts):**

```promql
# TCP receive queue drops
increase(node_netstat_TcpExt_TCPRcvQDrop[5m]) > 0

# Listen queue overflows
increase(node_netstat_TcpExt_ListenOverflows[5m]) > 0

# UDP buffer errors
increase(node_netstat_Udp_RcvbufErrors[5m]) > 0 or increase(node_netstat_Udp_SndbufErrors[5m]) > 0
```

**Relation to AF_PACKET / cost per packet:** When the kernel spends more time per packet (e.g. copying to AF_PACKET sockets), it can drain TCP/UDP socket buffers more slowly. So **high cost per packet** (observed via eBPF latency and softirq) can **contribute** to buffer full → TCPRcvQDrop, ListenOverflows, or RcvbufErrors. Use TCP/UDP buffer metrics as **outcomes** and eBPF + softirq as **cause** to conclude that the kernel path (including AF_PACKET) is increasing cost per packet and indirectly causing buffer drops.

---

## 8. Summary

- **Cost per packet** in the kernel = CPU time and wall-clock time spent per packet in the stack (interrupt/NAPI → protocol → socket delivery → free).
- **AF_PACKET** (and similar) **increase** that cost by adding per-packet work (e.g. `packet_rcv()`, copies, `consume_skb()`) in softirq.
- **Exported metrics** that matter:
  - **eBPF**: `node_ebpf_pmd_jitter_latency_*_ns`, `node_ebpf_pmd_jitter_pmd_jitter_ns`, `node_ebpf_pmd_jitter_latency_histogram`, global totals → **direct** kernel stack cost and jitter.
  - **CPU**: `node_cpu_seconds_total{mode="softirq"}`, `node_cpu_isolated` → **where** the cost shows up.
  - **Softirqs**: `node_softirqs_functions_total{type="NET_RX"}` → **receive path** (AF_PACKET-heavy).
  - **Nodeconfig**: `node_nodeconfig_cores_dedicated`, `node_nodeconfig_pcie_slot_ok`, etc. → **runbook context** (isolated CPUs, PCIe, memory banks).
  - **Meminfo**: `node_memory_HugePages_*`, `node_memory_Hugepagesize_bytes` → **runbook context** (hugepages / DPDK).
  - **Netdev**: `node_network_*` → **throughput and drops** (outcome of high cost).
  - **Conntrack** (Example A): `node_nf_conntrack_entries`, `node_nf_conntrack_entries_limit`, `node_nf_conntrack_stat_drop`, `node_nf_conntrack_stat_early_drop` → **full table → packet drops**.
  - **Netstat** (Example C): `node_netstat_TcpExt_TCPRcvQDrop`, `ListenOverflows`, `node_netstat_Udp_*bufErrors` → **TCP/UDP buffer drops**.
  - **Sockstat**: `node_sockstat_TCP_mem_bytes`, `node_sockstat_TCP_inuse` → **socket/buffer usage**.
  - **NUMA** (Example B): `node_zoneinfo_numa_*`, `node_memory_numa_*`, `node_pcidevice_numa_node` → **NUMA-unfriendly path → high latency**.
  - **Pressure**: `node_pressure_memory_*`, `node_pressure_io_*` → **stall**; can correlate with buffer pressure.

**Detailed examples:** **Example A** (full nf_conntrack → packet drops) uses conntrack metrics and drop counters. **Example B** (high packet latency from OVS not NUMA-optimized) uses NUMA hit/miss, pcidevice NUMA node, and eBPF latency. **Example C** (TCP buffer monitoring) uses netstat TCPRcvQDrop, ListenOverflows, UDP buffer errors, and sockstat; it ties buffer drops to cost per packet and AF_PACKET.

By correlating **rising kernel stack latency/jitter** with **high softirq** (and optionally NET_RX, conntrack drops, TCP buffer drops, or NUMA miss), you can **determine that the kernel stack—and mechanisms like AF_PACKET—are one way the cost per packet is increased**, and then narrow down to AF_PACKET via known capture activity or system tooling.

---

## 9. Functional test

A **Linux-only functional test** (`kernel_stack_af_packet_functional_test.go`) generates the scenarios above in a **dedicated network namespace** at **Gbps-scale data rates** where possible, stresses the kernel stack, and checks that the corresponding metrics and **.pcap** captures are produced. The test logs effective throughput (e.g. Gbps) for traffic scenarios and documents limitations of what the metrics can show in this environment (§9.1).

- **Scenario A (conntrack):** Lowers `nf_conntrack_max` in the netns, runs a TCP server that holds many connections, opens more than the limit from the host, then scrapes `node_nf_conntrack_entries`, `_entries_limit`, `_stat_drop`, `_stat_early_drop` and validates that table pressure and/or drops are reflected; traffic is captured to `scenario_a_conntrack.pcap`. **Limitation:** node_exporter may read `/proc` from the host, so conntrack metrics can reflect the host table, not the netns (§9.1).
- **Scenario C (listen overflow):** Server with listen backlog 1; many connections are opened quickly; asserts `node_netstat_TcpExt_ListenOverflows` / `ListenDrops` and captures to `scenario_c_listen.pcap`.
- **Scenario C (TCP rcvbuf):** Server with small `SO_RCVBUF` and slow read; client sends at high rate until rcvbuf fills; asserts `node_netstat_TcpExt_TCPRcvQDrop` and captures to `scenario_c_rcvq.pcap`.
- **Scenario B (traffic for NUMA):** High-throughput traffic over many connections to exercise the stack; on real NUMA hardware, `node_zoneinfo_numa_*` and `node_pcidevice_numa_node` would correlate with latency. **Limitation:** NUMA topology is not replicated in the netns (§9.1).
- **Traffic + pcap:** High data-rate traffic (target Gbps-scale) while capturing so netdev receive bytes, softirq, and drop metrics reflect per-packet cost; validates that the capture produces a valid `.pcap`. The test logs effective rate (Gbps) achieved.

**Requirements:** Linux, root (for `ip netns`, `tcpdump`), `node_exporter` binary (e.g. `make build`), and the stress server built from `cmd/kernel_stack_stress_server`.

**Run:** From the repo root, build the binary and run the test as root. If `go` is not in root’s PATH (e.g. it lives under `/usr/local/go/bin`), use the full path or preserve `PATH`:

```bash
make build
# Option A: preserve your PATH so root can find go
sudo env "PATH=$PATH" go test -v . -run TestKernelStackAFPacketScenarios -timeout 120s

# Option B: use full path to go
sudo $(which go) test -v . -run TestKernelStackAFPacketScenarios -timeout 120s
```

The test creates netns `kernel_stack_ftest`, a veth pair, starts node_exporter inside the netns on port 9192, runs each scenario, scrapes metrics from the host, and writes `.pcap` files under `/tmp/node_exporter_kernel_stack_pcaps_<timestamp>/`. It skips if not root or if `node_exporter` / `tcpdump` / `ip` are unavailable.

### 9.1 Limitations of the functional test (what the metrics can and cannot show)

The test runs in a **network namespace** with **veth** pairs. The following limitations determine what the exported metrics can and cannot demonstrate in this environment:

| Area | Limitation | What the metrics show in the test |
|------|------------|-----------------------------------|
| **Conntrack (Example A)** | node_exporter often shares the host’s `/proc` (or the test’s procfs is host-mounted). Conntrack counters and `nf_conntrack_max` are then from the **host**, not the netns. | The test can fill the netns conntrack table and cause drops in the netns; if the exporter reads host `/proc`, `node_nf_conntrack_entries_limit` will be the host limit (e.g. 262144), not the netns value (e.g. 50). The test logs this and only asserts conntrack pressure when the scraped limit looks like the netns value. |
| **NUMA (Example B)** | NUMA topology and PCI NUMA node are properties of real hardware. A netns has no separate NUMA topology; veth is not a physical NIC. | Traffic is generated at high rate to stress the stack. `node_zoneinfo_numa_*` and `node_pcidevice_numa_node` are **not** exercised by the test in a meaningful way; they are relevant on multi-NUMA machines with real NICs. The test documents that NUMA correlation requires real hardware. |
| **eBPF kernel stack latency** | The ebpf-pmd-jitter collector attaches to real devices (XDP/TC). In the test, interfaces are veth; the eBPF object may not be loaded or may not attach. | The test does **not** require eBPF to pass. Kernel stack latency and jitter metrics are meaningful when the collector is enabled and attached on real NICs under Gbps-scale load. |
| **Throughput (Gbps)** | veth is loopback-like; achievable rate is high but not representative of a physical 1/10/25 Gbps NIC. | The test aims for **Gbps-scale** data rates (many connections, large transfers) so that `node_network_receive_bytes_total`, softirq, and drop metrics reflect real stress. Achieved rate is logged (e.g. effective Gbps) so limitations of the environment are clear. |
| **AF_PACKET cost** | Per-packet cost (softirq, NET_RX) increases when tcpdump or another capture is active. Correlation of softirq with capture requires sustained high packet rate. | The test runs tcpdump in the netns during traffic and validates pcaps. It generates enough traffic so that netdev and (on real hardware) softirq metrics would show the effect of capture; in the test environment, the main outcome is valid metrics and pcaps. |
| **TCPRcvQDrop / ListenOverflows** | These are global kernel counters (netstat/sockstat). When node_exporter runs in the netns and reads netns `/proc`, the counters are netns-local. | The test reliably triggers ListenOverflows and TCPRcvQDrop in the netns and asserts the corresponding metrics, demonstrating that the **metrics correctly reflect** buffer and listen-queue drops under load. |
