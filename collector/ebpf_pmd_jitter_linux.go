// Copyright 2025 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build linux && !noebpfpmdjitter

// Package collector provides the ebpf-pmd-jitter collector, which loads an
// in-tree eBPF program (built from collector/bpf/latency.c) to measure kernel
// stack packet latency and exposes PMD-style jitter and latency metrics.
// Requires: --collector.ebpf-pmd-jitter.object-path pointing at the compiled
// latency.o and optionally --collector.ebpf-pmd-jitter.interfaces to attach.
package collector

import (
	"bytes"
	"fmt"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/alecthomas/kingpin/v2"
	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	ebpfPmdJitterSubsystem  = "ebpf_pmd_jitter"
	latencyHistogramBuckets = 16
)

var (
	ebpfPmdJitterObjectPath = kingpin.Flag("collector.ebpf-pmd-jitter.object-path",
		"Path to compiled eBPF object (latency.o) built from collector/bpf/latency.c. If empty, collector reports no data.").Default("").String()
	ebpfPmdJitterInterfaces = kingpin.Flag("collector.ebpf-pmd-jitter.interfaces",
		"Comma-separated list of interfaces to attach XDP latency measurement to (e.g. lo,eth0).").Default("lo").String()
)

type interfaceStats struct {
	PacketsTotal     uint64
	BytesTotal       uint64
	LatencyNsTotal   uint64
	LatencyMinNs     uint64
	LatencyMaxNs     uint64
	XDPPackets       uint64
	TCIngressPackets uint64
	TCEgressPackets  uint64
	SoftirqTimeNs    uint64
}

var latencyHistogramBucketLabels = []string{
	"0-1us", "1-2us", "2-4us", "4-8us", "8-16us", "16-32us", "32-64us",
	"64-128us", "128-256us", "256-512us", "512-1024us", "1-2ms",
	"2-4ms", "4-8ms", "8-16ms", "16ms+",
}

const maxLoadErrorLabelLen = 200

func sanitizeForLabel(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) > maxLoadErrorLabelLen {
		s = s[:maxLoadErrorLabelLen-3] + "..."
	}
	return s
}

type ebpfPmdJitterCollector struct {
	objectPath    string
	interfaces    []string
	logger        *slog.Logger
	mu            sync.Mutex
	lastLoadError string
	// Set on first successful load
	coll           *ebpf.Collection
	xdpLinks       map[string]link.Link
	ifindexToName_ map[int]string
	// Descriptors
	packetsTotalDesc  *prometheus.Desc
	bytesTotalDesc    *prometheus.Desc
	latencyMinNsDesc  *prometheus.Desc
	latencyMaxNsDesc  *prometheus.Desc
	latencyAvgNsDesc  *prometheus.Desc
	jitterNsDesc      *prometheus.Desc
	histogramDesc     *prometheus.Desc
	globalPacketsDesc *prometheus.Desc
	globalLatencyDesc *prometheus.Desc
	collectorUpDesc   *prometheus.Desc
	objectPathSetDesc *prometheus.Desc
	loadErrorDesc     *prometheus.Desc
}

func init() {
	registerCollector("ebpf-pmd-jitter", defaultDisabled, NewEbpfPmdJitterCollector)
}

// NewEbpfPmdJitterCollector returns a collector that loads the eBPF latency
// program from the given object path and exposes PMD jitter and latency metrics.
func NewEbpfPmdJitterCollector(logger *slog.Logger) (Collector, error) {
	ifaces := strings.Split(*ebpfPmdJitterInterfaces, ",")
	for i := range ifaces {
		ifaces[i] = strings.TrimSpace(ifaces[i])
	}
	return &ebpfPmdJitterCollector{
		objectPath:     *ebpfPmdJitterObjectPath,
		interfaces:     ifaces,
		logger:         logger,
		xdpLinks:       make(map[string]link.Link),
		ifindexToName_: make(map[int]string),
		packetsTotalDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, ebpfPmdJitterSubsystem, "latency_packets_total"),
			"Total packets measured for latency (eBPF kernel stack latency).",
			[]string{"interface"}, nil,
		),
		bytesTotalDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, ebpfPmdJitterSubsystem, "latency_bytes_total"),
			"Total bytes processed (eBPF latency measurement).",
			[]string{"interface"}, nil,
		),
		latencyMinNsDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, ebpfPmdJitterSubsystem, "latency_min_ns"),
			"Minimum observed packet latency in nanoseconds (kernel stack).",
			[]string{"interface"}, nil,
		),
		latencyMaxNsDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, ebpfPmdJitterSubsystem, "latency_max_ns"),
			"Maximum observed packet latency in nanoseconds (kernel stack).",
			[]string{"interface"}, nil,
		),
		latencyAvgNsDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, ebpfPmdJitterSubsystem, "latency_avg_ns"),
			"Average packet latency in nanoseconds (kernel stack).",
			[]string{"interface"}, nil,
		),
		jitterNsDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, ebpfPmdJitterSubsystem, "pmd_jitter_ns"),
			"PMD-style latency jitter in nanoseconds (latency_max_ns - latency_min_ns).",
			[]string{"interface"}, nil,
		),
		histogramDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, ebpfPmdJitterSubsystem, "latency_histogram"),
			"Histogram of packet latencies (kernel stack) by bucket.",
			[]string{"bucket"}, nil,
		),
		globalPacketsDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, ebpfPmdJitterSubsystem, "global_packets_total"),
			"Total packets measured globally by eBPF latency program.",
			nil, nil,
		),
		globalLatencyDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, ebpfPmdJitterSubsystem, "global_latency_ns_total"),
			"Total latency in nanoseconds (global).",
			nil, nil,
		),
		collectorUpDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, ebpfPmdJitterSubsystem, "collector_up"),
			"Whether the eBPF PMD jitter collector loaded and attached successfully (1 = yes, 0 = no).",
			nil, nil,
		),
		objectPathSetDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, ebpfPmdJitterSubsystem, "object_path_configured"),
			"Whether --collector.ebpf-pmd-jitter.object-path was set to a non-empty value (1 = yes, 0 = no). Use with collector_up: 0+0 = path not set; 0+1 = path set but load/attach failed.",
			nil, nil,
		),
		loadErrorDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, ebpfPmdJitterSubsystem, "load_error"),
			"Set to 1 when the eBPF program failed to load or attach; the 'error' label contains the reason. Absent when collector_up is 1.",
			[]string{"error"}, nil,
		),
	}, nil
}

func (c *ebpfPmdJitterCollector) Update(ch chan<- prometheus.Metric) error {
	objectPathSet := 0.0
	if c.objectPath != "" {
		objectPathSet = 1.0
	}
	ch <- prometheus.MustNewConstMetric(c.objectPathSetDesc, prometheus.GaugeValue, objectPathSet)

	if c.objectPath == "" {
		ch <- prometheus.MustNewConstMetric(c.collectorUpDesc, prometheus.GaugeValue, 0)
		return ErrNoData
	}

	c.mu.Lock()
	if c.coll == nil {
		if err := c.loadAndAttach(); err != nil {
			c.mu.Unlock()
			c.lastLoadError = err.Error()
			c.logger.Info("ebpf-pmd-jitter: load/attach failed", "err", err)
			ch <- prometheus.MustNewConstMetric(c.collectorUpDesc, prometheus.GaugeValue, 0)
			ch <- prometheus.MustNewConstMetric(c.loadErrorDesc, prometheus.GaugeValue, 1, sanitizeForLabel(err.Error()))
			return ErrNoData
		}
	}
	coll := c.coll
	c.mu.Unlock()

	ch <- prometheus.MustNewConstMetric(c.collectorUpDesc, prometheus.GaugeValue, 1)

	// Global stats
	globalPackets := c.sumPerCPUArray(coll.Maps["global_packets"], 0)
	globalLatency := c.sumPerCPUArray(coll.Maps["global_latency_ns"], 0)
	ch <- prometheus.MustNewConstMetric(c.globalPacketsDesc, prometheus.GaugeValue, float64(globalPackets))
	ch <- prometheus.MustNewConstMetric(c.globalLatencyDesc, prometheus.GaugeValue, float64(globalLatency))

	// Histogram
	histMap := coll.Maps["latency_histogram"]
	if histMap != nil {
		for i := 0; i < latencyHistogramBuckets; i++ {
			count := c.sumPerCPUArray(histMap, uint32(i))
			ch <- prometheus.MustNewConstMetric(c.histogramDesc, prometheus.GaugeValue, float64(count), latencyHistogramBucketLabels[i])
		}
	}

	// Per-interface stats and jitter
	ifStatsMap := coll.Maps["interface_latency_stats"]
	if ifStatsMap != nil {
		c.collectInterfaceStats(ifStatsMap, ch)
	}

	return nil
}

func (c *ebpfPmdJitterCollector) loadAndAttach() error {
	objPath := filepath.Clean(c.objectPath)
	data, err := os.ReadFile(objPath)
	if err != nil {
		return fmt.Errorf("read object file: %w", err)
	}

	spec, err := ebpf.LoadCollectionSpecFromReader(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("load collection spec: %w", err)
	}

	// Map pinning must be on a BPF filesystem (e.g. /sys/fs/bpf), not /tmp.
	const bpfPinBase = "/sys/fs/bpf"
	pinDir := filepath.Join(bpfPinBase, fmt.Sprintf("node_exporter_%d", os.Getpid()))
	if err := os.MkdirAll(pinDir, 0700); err != nil {
		return fmt.Errorf("create bpf pin dir %s (is %s mounted?): %w", pinDir, bpfPinBase, err)
	}
	// Do not remove pinDir while the collection is in use; kernel holds map refs via pins.

	opts := ebpf.CollectionOptions{
		Maps: ebpf.MapOptions{PinPath: pinDir},
	}

	coll, err := ebpf.NewCollectionWithOptions(spec, opts)
	if err != nil {
		return fmt.Errorf("new collection: %w", err)
	}

	xdpProg := coll.Programs["xdp_latency_ingress"]
	if xdpProg == nil {
		coll.Close()
		return fmt.Errorf("program xdp_latency_ingress not found in object")
	}

	for _, ifname := range c.interfaces {
		if ifname == "" {
			continue
		}
		iface, err := net.InterfaceByName(ifname)
		if err != nil {
			c.logger.Debug("ebpf-pmd-jitter: interface not found", "interface", ifname, "err", err)
			continue
		}
		c.ifindexToName_[iface.Index] = ifname
		l, err := link.AttachXDP(link.XDPOptions{
			Program:   xdpProg,
			Interface: iface.Index,
			Flags:     link.XDPGenericMode,
		})
		if err != nil {
			c.logger.Debug("ebpf-pmd-jitter: attach XDP failed", "interface", ifname, "err", err)
			continue
		}
		c.xdpLinks[ifname] = l
	}

	c.coll = coll
	return nil
}

func (c *ebpfPmdJitterCollector) collectInterfaceStats(ifStatsMap *ebpf.Map, ch chan<- prometheus.Metric) {
	iter := ifStatsMap.Iterate()
	var key uint32
	for iter.Next(&key, nil) {
		var values []interfaceStats
		if err := ifStatsMap.Lookup(&key, &values); err != nil {
			continue
		}
		ifname := c.ifindexToName_[int(key)]
		if ifname == "" {
			ifname = c.resolveInterfaceName(int(key))
		}
		var total interfaceStats
		for _, v := range values {
			total.PacketsTotal += v.PacketsTotal
			total.BytesTotal += v.BytesTotal
			total.LatencyNsTotal += v.LatencyNsTotal
			total.XDPPackets += v.XDPPackets
			total.TCIngressPackets += v.TCIngressPackets
			total.TCEgressPackets += v.TCEgressPackets
			if v.LatencyMinNs > 0 && (total.LatencyMinNs == 0 || v.LatencyMinNs < total.LatencyMinNs) {
				total.LatencyMinNs = v.LatencyMinNs
			}
			if v.LatencyMaxNs > total.LatencyMaxNs {
				total.LatencyMaxNs = v.LatencyMaxNs
			}
		}

		ch <- prometheus.MustNewConstMetric(c.packetsTotalDesc, prometheus.GaugeValue, float64(total.PacketsTotal), ifname)
		ch <- prometheus.MustNewConstMetric(c.bytesTotalDesc, prometheus.GaugeValue, float64(total.BytesTotal), ifname)
		ch <- prometheus.MustNewConstMetric(c.latencyMinNsDesc, prometheus.GaugeValue, float64(total.LatencyMinNs), ifname)
		ch <- prometheus.MustNewConstMetric(c.latencyMaxNsDesc, prometheus.GaugeValue, float64(total.LatencyMaxNs), ifname)

		jitter := float64(0)
		if total.LatencyMaxNs >= total.LatencyMinNs {
			jitter = float64(total.LatencyMaxNs - total.LatencyMinNs)
		}
		ch <- prometheus.MustNewConstMetric(c.jitterNsDesc, prometheus.GaugeValue, jitter, ifname)

		if total.PacketsTotal > 0 {
			avg := float64(total.LatencyNsTotal) / float64(total.PacketsTotal)
			ch <- prometheus.MustNewConstMetric(c.latencyAvgNsDesc, prometheus.GaugeValue, avg, ifname)
		}
	}
}

func (c *ebpfPmdJitterCollector) resolveInterfaceName(ifindex int) string {
	if ifindex <= 0 {
		return fmt.Sprintf("if%d", ifindex)
	}
	if iface, err := net.InterfaceByIndex(ifindex); err == nil {
		return iface.Name
	}
	return fmt.Sprintf("if%d", ifindex)
}

func (c *ebpfPmdJitterCollector) sumPerCPUArray(m *ebpf.Map, key uint32) uint64 {
	if m == nil {
		return 0
	}
	var values []uint64
	if err := m.Lookup(&key, &values); err != nil {
		var single uint64
		if err := m.Lookup(&key, &single); err != nil {
			return 0
		}
		return single
	}
	var total uint64
	for _, v := range values {
		total += v
	}
	return total
}
