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

package collector

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
	"github.com/prometheus/procfs/nfs"
)

// A nfsdCollector is a Collector which gathers metrics from /proc/net/rpc/nfsd.
// See: https://www.svennd.be/nfsd-stats-explained-procnetrpcnfsd/
type nfsdCollector struct {
	fs procfs.FS
}

func init() {
	registerCollector("nfsd", defaultEnabled, NewNFSdCollector)
}

const (
	nfsdSubsystem = "nfsd"
)

// NewNFSdCollector returns a new Collector exposing /proc/net/rpc/nfsd statistics.
func NewNFSdCollector() (Collector, error) {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %v", err)
	}

	return &nfsdCollector{
		fs: fs,
	}, nil
}

// Update implements Collector.
func (c *nfsdCollector) Update(ch chan<- prometheus.Metric) error {
	stats, err := c.fs.NFSdServerRPCStats()
	if err != nil {
		return fmt.Errorf("failed to retrieve nfsd stats: %v", err)
	}

	c.updateNFSdReplyCacheStats(ch, &stats.ReplyCache)

	return nil
}

// updateNFSdReplyCacheStats collects statistics for /proc/net/rpc/nfsd.
func (c *nfsdCollector) updateNFSdReplyCacheStats(ch chan<- prometheus.Metric, s *nfs.ReplyCache) {
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nfsdSubsystem, "reply_cache_hits_total"),
			"NFSd Reply Cache client did not receive a reply and decided to re-transmit its request and the reply was cached. (bad).",
			nil,
			nil,
		),
		prometheus.CounterValue,
		float64(s.Hits))
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nfsdSubsystem, "reply_cache_misses_total"),
			"NFSd Reply Cache an operation that requires caching (idempotent).",
			nil,
			nil,
		),
		prometheus.CounterValue,
		float64(s.Misses))
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nfsdSubsystem, "reply_cache_nocache_total"),
			"NFSd Reply Cache non-idempotent operations (rename/delete/…).",
			nil,
			nil,
		),
		prometheus.CounterValue,
		float64(s.NoCache))
}
