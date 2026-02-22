// SPDX-License-Identifier: GPL-2.0
/*
 * eBPF packet latency measurement for node_exporter.
 *
 * Measures time packets spend in the kernel network stack (XDP to TC).
 * Used by the node_exporter ebpf-pmd-jitter collector to expose
 * PMD-style latency jitter and kernel stack latency metrics.
 *
 * Build: make build-bpf (from node_exporter root), or:
 *   clang -O2 -g -target bpf -c collector/bpf/latency.c -o collector/bpf/latency.o \
 *     -I/usr/include
 */

#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/ipv6.h>
#include <linux/in.h>
#include <linux/tcp.h>
#include <linux/udp.h>
#include <linux/pkt_cls.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_endian.h>

#define MAX_TRACKED_PACKETS 65536
#define LATENCY_BUCKET_COUNT 16

#ifndef TC_ACT_OK
#define TC_ACT_OK 0
#endif

struct packet_timestamp {
    __u64 timestamp_ns;
    __u32 ifindex;
    __u32 len;
};

struct interface_stats {
    __u64 packets_total;
    __u64 bytes_total;
    __u64 latency_ns_total;
    __u64 latency_min_ns;
    __u64 latency_max_ns;
    __u64 xdp_packets;
    __u64 tc_ingress_packets;
    __u64 tc_egress_packets;
    __u64 softirq_time_ns;
};

struct {
    __uint(type, BPF_MAP_TYPE_LRU_HASH);
    __uint(max_entries, MAX_TRACKED_PACKETS);
    __type(key, __u32);
    __type(value, struct packet_timestamp);
    __uint(pinning, LIBBPF_PIN_BY_NAME);
} packet_timestamps SEC(".maps");

struct {
    __uint(type, BPF_MAP_TYPE_PERCPU_HASH);
    __uint(max_entries, 256);
    __type(key, __u32);
    __type(value, struct interface_stats);
    __uint(pinning, LIBBPF_PIN_BY_NAME);
} interface_latency_stats SEC(".maps");

struct {
    __uint(type, BPF_MAP_TYPE_PERCPU_ARRAY);
    __uint(max_entries, LATENCY_BUCKET_COUNT);
    __type(key, __u32);
    __type(value, __u64);
    __uint(pinning, LIBBPF_PIN_BY_NAME);
} latency_histogram SEC(".maps");

struct {
    __uint(type, BPF_MAP_TYPE_PERCPU_ARRAY);
    __uint(max_entries, 1);
    __type(key, __u32);
    __type(value, __u64);
    __uint(pinning, LIBBPF_PIN_BY_NAME);
} global_packets SEC(".maps");

struct {
    __uint(type, BPF_MAP_TYPE_PERCPU_ARRAY);
    __uint(max_entries, 1);
    __type(key, __u32);
    __type(value, __u64);
    __uint(pinning, LIBBPF_PIN_BY_NAME);
} global_latency_ns SEC(".maps");

static __always_inline __u32 calculate_packet_hash(void *data, void *data_end)
{
    struct ethhdr *eth = data;
    __u32 hash = 0;

    if ((void *)(eth + 1) > data_end)
        return 0;

    hash = eth->h_source[0] ^ eth->h_source[5];
    hash ^= eth->h_dest[0] ^ eth->h_dest[5];

    if (eth->h_proto == bpf_htons(ETH_P_IP)) {
        struct iphdr *ip = (void *)(eth + 1);
        if ((void *)(ip + 1) > data_end)
            return hash;
        hash ^= ip->saddr ^ ip->daddr ^ ip->protocol ^ ip->id;
        if (ip->protocol == IPPROTO_TCP || ip->protocol == IPPROTO_UDP) {
            __u16 *ports = (void *)ip + (ip->ihl * 4);
            if ((void *)(ports + 2) <= data_end) {
                hash ^= ports[0] ^ ports[1];
            }
        }
    } else if (eth->h_proto == bpf_htons(ETH_P_IPV6)) {
        struct ipv6hdr *ip6 = (void *)(eth + 1);
        if ((void *)(ip6 + 1) > data_end)
            return hash;
        hash ^= ip6->saddr.s6_addr32[0] ^ ip6->saddr.s6_addr32[3];
        hash ^= ip6->daddr.s6_addr32[0] ^ ip6->daddr.s6_addr32[3];
        hash ^= ip6->nexthdr;
    }
    return hash;
}

static __always_inline __u32 get_latency_bucket(__u64 latency_ns)
{
    __u64 us = latency_ns / 1000;
    if (us < 1) return 0;
    if (us < 2) return 1;
    if (us < 4) return 2;
    if (us < 8) return 3;
    if (us < 16) return 4;
    if (us < 32) return 5;
    if (us < 64) return 6;
    if (us < 128) return 7;
    if (us < 256) return 8;
    if (us < 512) return 9;
    if (us < 1024) return 10;
    if (us < 2048) return 11;
    if (us < 4096) return 12;
    if (us < 8192) return 13;
    if (us < 16384) return 14;
    return 15;
}

static __always_inline void update_histogram(__u64 latency_ns)
{
    __u32 bucket = get_latency_bucket(latency_ns);
    __u64 *count = bpf_map_lookup_elem(&latency_histogram, &bucket);
    if (count)
        __sync_fetch_and_add(count, 1);
}

static __always_inline void update_interface_stats(__u32 ifindex, __u32 pkt_len,
                                                  __u64 latency_ns, int hook_type)
{
    struct interface_stats *stats = bpf_map_lookup_elem(&interface_latency_stats, &ifindex);
    if (!stats) {
        struct interface_stats new_stats = {};
        new_stats.latency_min_ns = latency_ns > 0 ? latency_ns : ~0ULL;
        new_stats.latency_max_ns = latency_ns;
        bpf_map_update_elem(&interface_latency_stats, &ifindex, &new_stats, BPF_ANY);
        stats = bpf_map_lookup_elem(&interface_latency_stats, &ifindex);
        if (!stats)
            return;
    }
    __sync_fetch_and_add(&stats->packets_total, 1);
    __sync_fetch_and_add(&stats->bytes_total, pkt_len);
    if (latency_ns > 0) {
        __sync_fetch_and_add(&stats->latency_ns_total, latency_ns);
        if (latency_ns < stats->latency_min_ns || stats->latency_min_ns == 0)
            stats->latency_min_ns = latency_ns;
        if (latency_ns > stats->latency_max_ns)
            stats->latency_max_ns = latency_ns;
    }
    if (hook_type == 0)
        __sync_fetch_and_add(&stats->xdp_packets, 1);
    else if (hook_type == 1)
        __sync_fetch_and_add(&stats->tc_ingress_packets, 1);
    else if (hook_type == 2)
        __sync_fetch_and_add(&stats->tc_egress_packets, 1);
}

static __always_inline void update_global_stats(__u64 latency_ns)
{
    __u32 key = 0;
    __u64 *packets = bpf_map_lookup_elem(&global_packets, &key);
    if (packets)
        __sync_fetch_and_add(packets, 1);
    if (latency_ns > 0) {
        __u64 *latency = bpf_map_lookup_elem(&global_latency_ns, &key);
        if (latency)
            __sync_fetch_and_add(latency, latency_ns);
    }
}

SEC("xdp")
int xdp_latency_ingress(struct xdp_md *ctx)
{
    void *data = (void *)(long)ctx->data;
    void *data_end = (void *)(long)ctx->data_end;
    __u32 pkt_len = (__u32)(data_end - data);
    __u64 now = bpf_ktime_get_ns();
    __u32 hash = calculate_packet_hash(data, data_end);
    if (hash == 0)
        goto out;
    struct packet_timestamp ts = {
        .timestamp_ns = now,
        .ifindex = ctx->ingress_ifindex,
        .len = pkt_len,
    };
    bpf_map_update_elem(&packet_timestamps, &hash, &ts, BPF_ANY);
    update_interface_stats(ctx->ingress_ifindex, pkt_len, 0, 0);
out:
    return XDP_PASS;
}

SEC("tc")
int tc_latency_ingress(struct __sk_buff *skb)
{
    void *data = (void *)(long)skb->data;
    void *data_end = (void *)(long)skb->data_end;
    __u64 now = bpf_ktime_get_ns();
    __u64 latency_ns = 0;
    __u32 hash = calculate_packet_hash(data, data_end);
    if (hash != 0) {
        struct packet_timestamp *ts = bpf_map_lookup_elem(&packet_timestamps, &hash);
        if (ts && ts->timestamp_ns > 0) {
            latency_ns = now - ts->timestamp_ns;
            update_histogram(latency_ns);
            update_global_stats(latency_ns);
        }
    }
    update_interface_stats(skb->ifindex, skb->len, latency_ns, 1);
    return TC_ACT_OK;
}

SEC("tc")
int tc_latency_egress(struct __sk_buff *skb)
{
    void *data = (void *)(long)skb->data;
    void *data_end = (void *)(long)skb->data_end;
    __u64 now = bpf_ktime_get_ns();
    __u64 latency_ns = 0;
    __u32 hash = calculate_packet_hash(data, data_end);
    if (hash != 0) {
        struct packet_timestamp *ts = bpf_map_lookup_elem(&packet_timestamps, &hash);
        if (ts && ts->timestamp_ns > 0) {
            latency_ns = now - ts->timestamp_ns;
            update_histogram(latency_ns);
            update_global_stats(latency_ns);
            bpf_map_delete_elem(&packet_timestamps, &hash);
        }
    }
    update_interface_stats(skb->ifindex, skb->len, latency_ns, 2);
    return TC_ACT_OK;
}

char _license[] SEC("license") = "GPL";
