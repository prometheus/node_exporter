// Copyright 2018 The Prometheus Authors
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

//go:build !nonfsd
// +build !nonfsd

package collector

import (
	"errors"
	"fmt"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/nfs"
)

// A nfsdCollector is a Collector which gathers metrics from /proc/net/rpc/nfsd.
// See: https://www.svennd.be/nfsd-stats-explained-procnetrpcnfsd/
type nfsdCollector struct {
	fs           nfs.FS
	requestsDesc *prometheus.Desc
	logger       log.Logger
}

func init() {
	registerCollector("nfsd", defaultEnabled, NewNFSdCollector)
}

const (
	nfsdSubsystem = "nfsd"
)

// NewNFSdCollector returns a new Collector exposing /proc/net/rpc/nfsd statistics.
func NewNFSdCollector(logger log.Logger) (Collector, error) {
	fs, err := nfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}

	return &nfsdCollector{
		fs: fs,
		requestsDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nfsdSubsystem, "requests_total"),
			"Total number NFSd Requests by method and protocol.",
			[]string{"proto", "method"}, nil,
		),
		logger: logger,
	}, nil
}

// Update implements Collector.
func (c *nfsdCollector) Update(ch chan<- prometheus.Metric) error {
	stats, err := c.fs.ServerRPCStats()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			level.Debug(c.logger).Log("msg", "Not collecting NFSd metrics", "err", err)
			return ErrNoData
		}
		return fmt.Errorf("failed to retrieve nfsd stats: %w", err)
	}

	c.updateNFSdReplyCacheStats(ch, &stats.ReplyCache)
	c.updateNFSdFileHandlesStats(ch, &stats.FileHandles)
	c.updateNFSdInputOutputStats(ch, &stats.InputOutput)
	c.updateNFSdThreadsStats(ch, &stats.Threads)
	c.updateNFSdReadAheadCacheStats(ch, &stats.ReadAheadCache)
	c.updateNFSdNetworkStats(ch, &stats.Network)
	c.updateNFSdServerRPCStats(ch, &stats.ServerRPC)
	c.updateNFSdRequestsv2Stats(ch, &stats.V2Stats)
	c.updateNFSdRequestsv3Stats(ch, &stats.V3Stats)
	c.updateNFSdRequestsv4Stats(ch, &stats.V4Ops)

	return nil
}

// updateNFSdReplyCacheStats collects statistics for the reply cache.
func (c *nfsdCollector) updateNFSdReplyCacheStats(ch chan<- prometheus.Metric, s *nfs.ReplyCache) {
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nfsdSubsystem, "reply_cache_hits_total"),
			"Total number of NFSd Reply Cache hits (client lost server response).",
			nil,
			nil,
		),
		prometheus.CounterValue,
		float64(s.Hits))
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nfsdSubsystem, "reply_cache_misses_total"),
			"Total number of NFSd Reply Cache an operation that requires caching (idempotent).",
			nil,
			nil,
		),
		prometheus.CounterValue,
		float64(s.Misses))
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nfsdSubsystem, "reply_cache_nocache_total"),
			"Total number of NFSd Reply Cache non-idempotent operations (rename/delete/…).",
			nil,
			nil,
		),
		prometheus.CounterValue,
		float64(s.NoCache))
}

// updateNFSdFileHandlesStats collects statistics for the file handles.
func (c *nfsdCollector) updateNFSdFileHandlesStats(ch chan<- prometheus.Metric, s *nfs.FileHandles) {
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nfsdSubsystem, "file_handles_stale_total"),
			"Total number of NFSd stale file handles",
			nil,
			nil,
		),
		prometheus.CounterValue,
		float64(s.Stale))
	// NOTE: Other FileHandles entries are unused in the kernel.
}

// updateNFSdInputOutputStats collects statistics for the bytes in/out.
func (c *nfsdCollector) updateNFSdInputOutputStats(ch chan<- prometheus.Metric, s *nfs.InputOutput) {
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nfsdSubsystem, "disk_bytes_read_total"),
			"Total NFSd bytes read.",
			nil,
			nil,
		),
		prometheus.CounterValue,
		float64(s.Read))
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nfsdSubsystem, "disk_bytes_written_total"),
			"Total NFSd bytes written.",
			nil,
			nil,
		),
		prometheus.CounterValue,
		float64(s.Write))
}

// updateNFSdThreadsStats collects statistics for kernel server threads.
func (c *nfsdCollector) updateNFSdThreadsStats(ch chan<- prometheus.Metric, s *nfs.Threads) {
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nfsdSubsystem, "server_threads"),
			"Total number of NFSd kernel threads that are running.",
			nil,
			nil,
		),
		prometheus.GaugeValue,
		float64(s.Threads))
}

// updateNFSdReadAheadCacheStats collects statistics for the read ahead cache.
func (c *nfsdCollector) updateNFSdReadAheadCacheStats(ch chan<- prometheus.Metric, s *nfs.ReadAheadCache) {
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nfsdSubsystem, "read_ahead_cache_size_blocks"),
			"How large the read ahead cache is in blocks.",
			nil,
			nil,
		),
		prometheus.GaugeValue,
		float64(s.CacheSize))
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nfsdSubsystem, "read_ahead_cache_not_found_total"),
			"Total number of NFSd read ahead cache not found.",
			nil,
			nil,
		),
		prometheus.CounterValue,
		float64(s.NotFound))
}

// updateNFSdNetworkStats collects statistics for network packets/connections.
func (c *nfsdCollector) updateNFSdNetworkStats(ch chan<- prometheus.Metric, s *nfs.Network) {
	packetDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nfsdSubsystem, "packets_total"),
		"Total NFSd network packets (sent+received) by protocol type.",
		[]string{"proto"},
		nil,
	)
	ch <- prometheus.MustNewConstMetric(
		packetDesc,
		prometheus.CounterValue,
		float64(s.UDPCount), "udp")
	ch <- prometheus.MustNewConstMetric(
		packetDesc,
		prometheus.CounterValue,
		float64(s.TCPCount), "tcp")
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nfsdSubsystem, "connections_total"),
			"Total number of NFSd TCP connections.",
			nil,
			nil,
		),
		prometheus.CounterValue,
		float64(s.TCPConnect))
}

// updateNFSdServerRPCStats collects statistics for kernel server RPCs.
func (c *nfsdCollector) updateNFSdServerRPCStats(ch chan<- prometheus.Metric, s *nfs.ServerRPC) {
	badRPCDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nfsdSubsystem, "rpc_errors_total"),
		"Total number of NFSd RPC errors by error type.",
		[]string{"error"},
		nil,
	)
	ch <- prometheus.MustNewConstMetric(
		badRPCDesc,
		prometheus.CounterValue,
		float64(s.BadFmt), "fmt")
	ch <- prometheus.MustNewConstMetric(
		badRPCDesc,
		prometheus.CounterValue,
		float64(s.BadAuth), "auth")
	ch <- prometheus.MustNewConstMetric(
		badRPCDesc,
		prometheus.CounterValue,
		float64(s.BadcInt), "cInt")
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nfsdSubsystem, "server_rpcs_total"),
			"Total number of NFSd RPCs.",
			nil,
			nil,
		),
		prometheus.CounterValue,
		float64(s.RPCCount))
}

// updateNFSdRequestsv2Stats collects statistics for NFSv2 requests.
func (c *nfsdCollector) updateNFSdRequestsv2Stats(ch chan<- prometheus.Metric, s *nfs.V2Stats) {
	const proto = "2"
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.GetAttr), proto, "GetAttr")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.SetAttr), proto, "SetAttr")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.Root), proto, "Root")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.Lookup), proto, "Lookup")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.ReadLink), proto, "ReadLink")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.Read), proto, "Read")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.WrCache), proto, "WrCache")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.Write), proto, "Write")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.Create), proto, "Create")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.Remove), proto, "Remove")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.Rename), proto, "Rename")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.Link), proto, "Link")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.SymLink), proto, "SymLink")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.MkDir), proto, "MkDir")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.RmDir), proto, "RmDir")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.ReadDir), proto, "ReadDir")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.FsStat), proto, "FsStat")
}

// updateNFSdRequestsv3Stats collects statistics for NFSv3 requests.
func (c *nfsdCollector) updateNFSdRequestsv3Stats(ch chan<- prometheus.Metric, s *nfs.V3Stats) {
	const proto = "3"
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.GetAttr), proto, "GetAttr")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.SetAttr), proto, "SetAttr")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.Lookup), proto, "Lookup")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.Access), proto, "Access")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.ReadLink), proto, "ReadLink")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.Read), proto, "Read")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.Write), proto, "Write")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.Create), proto, "Create")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.MkDir), proto, "MkDir")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.SymLink), proto, "SymLink")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.MkNod), proto, "MkNod")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.Remove), proto, "Remove")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.RmDir), proto, "RmDir")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.Rename), proto, "Rename")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.Link), proto, "Link")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.ReadDir), proto, "ReadDir")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.ReadDirPlus), proto, "ReadDirPlus")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.FsStat), proto, "FsStat")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.FsInfo), proto, "FsInfo")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.PathConf), proto, "PathConf")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.Commit), proto, "Commit")
}

// updateNFSdRequestsv4Stats collects statistics for NFSv4 requests.
func (c *nfsdCollector) updateNFSdRequestsv4Stats(ch chan<- prometheus.Metric, s *nfs.V4Ops) {
	const proto = "4"
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.Access), proto, "Access")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.Close), proto, "Close")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.Commit), proto, "Commit")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.Create), proto, "Create")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.DelegPurge), proto, "DelegPurge")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.DelegReturn), proto, "DelegReturn")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.GetAttr), proto, "GetAttr")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.GetFH), proto, "GetFH")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.Link), proto, "Link")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.Lock), proto, "Lock")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.Lockt), proto, "Lockt")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.Locku), proto, "Locku")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.Lookup), proto, "Lookup")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.LookupRoot), proto, "LookupRoot")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.Nverify), proto, "Nverify")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.Open), proto, "Open")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.OpenAttr), proto, "OpenAttr")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.OpenConfirm), proto, "OpenConfirm")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.OpenDgrd), proto, "OpenDgrd")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.PutFH), proto, "PutFH")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.Read), proto, "Read")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.ReadDir), proto, "ReadDir")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.ReadLink), proto, "ReadLink")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.Remove), proto, "Remove")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.Rename), proto, "Rename")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.Renew), proto, "Renew")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.RestoreFH), proto, "RestoreFH")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.SaveFH), proto, "SaveFH")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.SecInfo), proto, "SecInfo")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.SetAttr), proto, "SetAttr")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.Verify), proto, "Verify")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.Write), proto, "Write")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue,
		float64(s.RelLockOwner), proto, "RelLockOwner")
}
