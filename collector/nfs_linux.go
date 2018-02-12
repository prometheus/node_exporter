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

	"github.com/prometheus/client_golang/prometheus"
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
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Null), proto, "null")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.GetAttr), proto, "getattr")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.SetAttr), proto, "setattr")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Root), proto, "root")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Lookup), proto, "lookup")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.ReadLink), proto, "readlink")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Read), proto, "read")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.WrCache), proto, "writecache")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Write), proto, "write")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Create), proto, "create")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Remove), proto, "remove")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Rename), proto, "rename")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Link), proto, "link")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.SymLink), proto, "symlink")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.MkDir), proto, "mkdir")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.RmDir), proto, "rmdir")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.ReadDir), proto, "readdir")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.FsStat), proto, "statfs")
}

// updateNFSRequestsv3Stats collects statistics for NFSv3 requests.
func (c *nfsCollector) updateNFSRequestsv3Stats(ch chan<- prometheus.Metric, s *nfs.V3Stats) {
	const proto = "3"
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Null), proto, "null")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.GetAttr), proto, "getattr")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.SetAttr), proto, "setattr")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Lookup), proto, "lookup")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Access), proto, "access")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.ReadLink), proto, "readlink")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Read), proto, "read")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Write), proto, "write")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Create), proto, "create")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.MkDir), proto, "mkdir")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.SymLink), proto, "symlink")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.MkNod), proto, "mknod")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Remove), proto, "remove")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.RmDir), proto, "rmdir")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Rename), proto, "rename")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Link), proto, "link")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.ReadDir), proto, "readdir")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.ReadDirPlus), proto, "readdirplus")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.FsStat), proto, "fsstat")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.FsInfo), proto, "fsinfo")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.PathConf), proto, "pathconf")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Commit), proto, "commit")
}

// updateNFSRequestsv4Stats collects statistics for NFSv4 requests.
func (c *nfsCollector) updateNFSRequestsv4Stats(ch chan<- prometheus.Metric, s *nfs.ClientV4Stats) {
	const proto = "4"
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Null), proto, "null")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Read), proto, "read")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Write), proto, "write")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Commit), proto, "commit")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Open), proto, "open")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.OpenConfirm), proto, "open_confirm")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.OpenNoattr), proto, "open_noattr")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.OpenDowngrade), proto, "open_downgrade")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Close), proto, "close")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Setattr), proto, "setattr")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.FsInfo), proto, "fsinfo")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Renew), proto, "renew")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.SetClientId), proto, "setclientid")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.SetClientIdConfirm), proto, "setclientid_confirm")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Lock), proto, "lock")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Lockt), proto, "lockt")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Locku), proto, "locku")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Access), proto, "access")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Getattr), proto, "getattr")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Lookup), proto, "lookup")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.LookupRoot), proto, "lookup_root")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Remove), proto, "remove")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Rename), proto, "rename")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Link), proto, "link")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Symlink), proto, "symlink")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Create), proto, "create")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Pathconf), proto, "pathconf")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.StatFs), proto, "statfs")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.ReadLink), proto, "readlink")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.ReadDir), proto, "readdir")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.ServerCaps), proto, "server_caps")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.DelegReturn), proto, "delegreturn")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.GetAcl), proto, "getacl")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.SetAcl), proto, "setacl")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.FsLocations), proto, "fs_locations")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.ReleaseLockowner), proto, "release_lockowner")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Secinfo), proto, "secinfo")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.FsidPresent), proto, "fsid_present")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.ExchangeId), proto, "exchange_id")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.CreateSession), proto, "create_session")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.DestroySession), proto, "destroy_session")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.Sequence), proto, "sequence")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.GetLeaseTime), proto, "get_lease_time")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.ReclaimComplete), proto, "reclaim_complete")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.LayoutGet), proto, "layoutget")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.GetDeviceInfo), proto, "getdeviceinfo")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.LayoutCommit), proto, "layoutcommit")
	ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
		float64(s.LayoutReturn), proto, "layoutreturn")
	// TODO: Enable after testing feature parity.
	// ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
	// 	float64(s.SecinfoNoName), proto, "secinfo_no_name")
	// ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
	// 	float64(s.TestStateId), proto, "test_stateid")
	// ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
	// 	float64(s.FreeStateId), proto, "free_stateid")
	// ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
	// 	float64(s.GetDeviceList), proto, "getdevicelist")
	// ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
	// 	float64(s.BindConnToSession), proto, "bind_conn_to_session")
	// ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
	// 	float64(s.DestroyClientId), proto, "destroy_clientid")
	// ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
	// 	float64(s.Seek), proto, "seek")
	// ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
	// 	float64(s.Allocate), proto, "allocate")
	// ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
	// 	float64(s.DeAllocate), proto, "deallocate")
	// ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
	// 	float64(s.LayoutStats), proto, "layoutstats")
	// ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue,
	// 	float64(s.Clone), proto, "clone")
}
