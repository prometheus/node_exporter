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

package collector

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/procfs"
	"github.com/prometheus/procfs/nfs"
)

const (
	nfsSubsystem = "nfs"
)

type nfsCollector struct {
	fs                                procfs.FS
	nfsNetReadsDesc                   *prometheus.Desc
	nfsNetConnectionsDesc             *prometheus.Desc
	nfsRPCOperationsDesc              *prometheus.Desc
	nfsRPCRetransmissionsDesc         *prometheus.Desc
	nfsRPCAuthenticationRefreshesDesc *prometheus.Desc
	nfsProceduresDesc                 *prometheus.Desc
}

func init() {
	registerCollector("nfs", defaultDisabled, NewNfsCollector)
}

// NewNfsCollector returns a new Collector exposing NFS statistics.
func NewNfsCollector() (Collector, error) {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %v", err)
	}

	return &nfsCollector{
		fs: fs,
		nfsNetReadsDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nfsSubsystem, "net_reads_total"),
			"Number of reads at the network layer.",
			[]string{"protocol"},
			nil,
		),
		nfsNetConnectionsDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nfsSubsystem, "net_connections_total"),
			"Number of connections at the network layer.",
			[]string{"protocol"},
			nil,
		),
		nfsRPCOperationsDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nfsSubsystem, "rpc_operations_total"),
			"Number of RPCs performed.",
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
			prometheus.BuildFQName(namespace, "nfs", "procedures_total"),
			"Number of NFS procedures invoked.",
			[]string{"version", "procedure"},
			nil,
		),
	}, nil
}

func (c *nfsCollector) Update(ch chan<- prometheus.Metric) error {
	stats, err := c.fs.NFSClientRPCStats()
	if err != nil {
		if os.IsNotExist(err) {
			log.Debugf("Not collecting NFS metrics: %s", err)
			return nil
		}
		return fmt.Errorf("failed to retrieve nfs stats: %v", err)
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
		float64(s.TCPConnect), "tcp")
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
		name := strings.ToLower(v.Type().Field(i).Name)

		switch name {
		case "wrcache":
			name = "writecache"
		case "fsstat":
			name = "statfs"
		}

		ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
			float64(field.Uint()), proto, name)
	}
}

// updateNFSRequestsv3Stats collects statistics for NFSv3 requests.
func (c *nfsCollector) updateNFSRequestsv3Stats(ch chan<- prometheus.Metric, s *nfs.V3Stats) {
	const proto = "3"

	v := reflect.ValueOf(s).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		name := strings.ToLower(v.Type().Field(i).Name)

		ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
			float64(field.Uint()), proto, name)
	}
}

// updateNFSRequestsv4Stats collects statistics for NFSv4 requests.
func (c *nfsCollector) updateNFSRequestsv4Stats(ch chan<- prometheus.Metric, s *nfs.ClientV4Stats) {
	const proto = "4"

	v := reflect.ValueOf(s).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		name := strings.ToLower(v.Type().Field(i).Name)

		switch name {
		case "openconfirm":
			name = "open_confirm"
		case "opendowngrade":
			name = "open_downgrade"
		case "opennoattr":
			name = "open_noattr"
		case "setclientidconfirm":
			name = "setclientid_confirm"
		case "lookuproot":
			name = "lookup_root"
		case "servercaps":
			name = "server_caps"
		case "fslocations":
			name = "fs_locations"
		case "releaselockowner":
			name = "release_lockowner"
		case "fsidpresent":
			name = "fsid_present"
		case "exchangeid":
			name = "exchange_id"
		case "createsession":
			name = "create_session"
		case "destroysession":
			name = "destroy_session"
		case "getleasetime":
			name = "get_lease_time"
		case "reclaimcomplete":
			name = "reclaim_complete"
		// TODO: Enable these metrics
		case "secinfononame", "teststateid", "freestateid", "getdevicelist", "bindconntosession", "destroyclientid", "seek", "allocate", "deallocate", "layoutstats", "clone":
			continue
		}

		ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
			float64(field.Uint()), proto, name)
	}
}
