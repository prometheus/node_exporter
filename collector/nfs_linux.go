// Copyright 2016 The Prometheus Authors
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

//go:build !nonfs
// +build !nonfs

package collector

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/nfs"
)

const (
	nfsSubsystem = "nfs"
)

type nfsCollector struct {
	fs                                nfs.FS
	nfsNetReadsDesc                   *prometheus.Desc
	nfsNetConnectionsDesc             *prometheus.Desc
	nfsRPCOperationsDesc              *prometheus.Desc
	nfsRPCRetransmissionsDesc         *prometheus.Desc
	nfsRPCAuthenticationRefreshesDesc *prometheus.Desc
	nfsProceduresDesc                 *prometheus.Desc
	logger                            log.Logger
}

func init() {
	registerCollector("nfs", defaultEnabled, NewNfsCollector)
}

// NewNfsCollector returns a new Collector exposing NFS statistics.
func NewNfsCollector(logger log.Logger) (Collector, error) {
	fs, err := nfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}

	return &nfsCollector{
		fs: fs,
		nfsNetReadsDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nfsSubsystem, "packets_total"),
			"Total NFSd network packets (sent+received) by protocol type.",
			[]string{"protocol"},
			nil,
		),
		nfsNetConnectionsDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nfsSubsystem, "connections_total"),
			"Total number of NFSd TCP connections.",
			nil,
			nil,
		),
		nfsRPCOperationsDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nfsSubsystem, "rpcs_total"),
			"Total number of RPCs performed.",
			nil,
			nil,
		),
		nfsRPCRetransmissionsDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nfsSubsystem, "rpc_retransmissions_total"),
			"Number of RPC transmissions performed.",
			nil,
			nil,
		),
		nfsRPCAuthenticationRefreshesDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nfsSubsystem, "rpc_authentication_refreshes_total"),
			"Number of RPC authentication refreshes performed.",
			nil,
			nil,
		),
		nfsProceduresDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nfsSubsystem, "requests_total"),
			"Number of NFS procedures invoked.",
			[]string{"proto", "method"},
			nil,
		),
		logger: logger,
	}, nil
}

func (c *nfsCollector) Update(ch chan<- prometheus.Metric) error {
	stats, err := c.fs.ClientRPCStats()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			level.Debug(c.logger).Log("msg", "Not collecting NFS metrics", "err", err)
			return ErrNoData
		}
		return fmt.Errorf("failed to retrieve nfs stats: %w", err)
	}

	c.updateNFSNetworkStats(ch, &stats.Network)
	c.updateNFSClientRPCStats(ch, &stats.ClientRPC)
	c.updateNFSRequestsv2Stats(ch, &stats.V2Stats)
	c.updateNFSRequestsv3Stats(ch, &stats.V3Stats)
	c.updateNFSRequestsv4Stats(ch, &stats.ClientV4Stats)

	return nil
}

// updateNFSNetworkStats collects statistics for network packets/connections.
func (c *nfsCollector) updateNFSNetworkStats(ch chan<- prometheus.Metric, s *nfs.Network) {
	ch <- prometheus.MustNewConstMetric(c.nfsNetReadsDesc, prometheus.CounterValue,
		float64(s.UDPCount), "udp")
	ch <- prometheus.MustNewConstMetric(c.nfsNetReadsDesc, prometheus.CounterValue,
		float64(s.TCPCount), "tcp")
	ch <- prometheus.MustNewConstMetric(c.nfsNetConnectionsDesc, prometheus.CounterValue,
		float64(s.TCPConnect))
}

// updateNFSClientRPCStats collects statistics for kernel server RPCs.
func (c *nfsCollector) updateNFSClientRPCStats(ch chan<- prometheus.Metric, s *nfs.ClientRPC) {
	ch <- prometheus.MustNewConstMetric(c.nfsRPCOperationsDesc, prometheus.CounterValue,
		float64(s.RPCCount))
	ch <- prometheus.MustNewConstMetric(c.nfsRPCRetransmissionsDesc, prometheus.CounterValue,
		float64(s.Retransmissions))
	ch <- prometheus.MustNewConstMetric(c.nfsRPCAuthenticationRefreshesDesc, prometheus.CounterValue,
		float64(s.AuthRefreshes))
}

// updateNFSRequestsv2Stats collects statistics for NFSv2 requests.
func (c *nfsCollector) updateNFSRequestsv2Stats(ch chan<- prometheus.Metric, s *nfs.V2Stats) {
	const proto = "2"

	v := reflect.ValueOf(s).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
			float64(field.Uint()), proto, v.Type().Field(i).Name)
	}
}

// updateNFSRequestsv3Stats collects statistics for NFSv3 requests.
func (c *nfsCollector) updateNFSRequestsv3Stats(ch chan<- prometheus.Metric, s *nfs.V3Stats) {
	const proto = "3"

	v := reflect.ValueOf(s).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
			float64(field.Uint()), proto, v.Type().Field(i).Name)
	}
}

// updateNFSRequestsv4Stats collects statistics for NFSv4 requests.
func (c *nfsCollector) updateNFSRequestsv4Stats(ch chan<- prometheus.Metric, s *nfs.ClientV4Stats) {
	const proto = "4"

	v := reflect.ValueOf(s).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
			float64(field.Uint()), proto, v.Type().Field(i).Name)
	}
}
