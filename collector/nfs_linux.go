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
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

var (
	netLineRE  = regexp.MustCompile(`^net \d+ (\d+) (\d+) (\d+)$`)
	rpcLineRE  = regexp.MustCompile(`^rpc (\d+) (\d+) (\d+)$`)
	procLineRE = regexp.MustCompile(`^proc(\d+) \d+ (\d+( \d+)*)$`)

	nfsProcedures = map[string][]string{
		"2": []string{
			"null", "getattr", "setattr", "root", "lookup",
			"readlink", "read", "writecache", "write", "create",
			"remove", "rename", "link", "symlink", "mkdir",
			"rmdir", "readdir", "statfs",
		},
		"3": []string{
			"null", "getattr", "setattr", "lookup", "access",
			"readlink", "read", "write", "create", "mkdir",
			"symlink", "mknod", "remove", "rmdir", "rename",
			"link", "readdir", "readdirplus", "fsstat", "fsinfo",
			"pathconf", "commit",
		},
		"4": []string{
			"null", "read", "write", "commit", "open",
			"open_confirm", "open_noattr", "open_downgrade",
			"close", "setattr", "fsinfo", "renew", "setclientid",
			"setclientid_confirm", "lock", "lockt", "locku",
			"access", "getattr", "lookup", "lookup_root", "remove",
			"rename", "link", "symlink", "create", "pathconf",
			"statfs", "readlink", "readdir", "server_caps",
			"delegreturn", "getacl", "setacl", "fs_locations",
			"release_lockowner", "secinfo", "fsid_present",
			"exchange_id", "create_session", "destroy_session",
			"sequence", "get_lease_time", "reclaim_complete",
			"layoutget", "getdeviceinfo", "layoutcommit",
			"layoutreturn", "secinfo_no_name", "test_stateid",
			"free_stateid", "getdevicelist",
			"bind_conn_to_session", "destroy_clientid", "seek",
			"allocate", "deallocate", "layoutstats", "clone",
			"copy",
		},
	}

	nfsNetReadsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "nfs", "net_reads"),
		"Number of reads at the network layer.",
		[]string{"protocol"},
		nil,
	)
	nfsNetConnectionsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "nfs", "net_connections"),
		"Number of connections at the network layer.",
		[]string{"protocol"},
		nil,
	)

	nfsRpcOperationsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "nfs", "rpc_operations"),
		"Number of RPCs performed.",
		nil,
		nil,
	)
	nfsRpcRetransmissionsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "nfs", "rpc_retransmissions"),
		"Number of RPC transmissions performed.",
		nil,
		nil,
	)
	nfsRpcAuthenticationRefreshesDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "nfs", "rpc_authentication_refreshes"),
		"Number of RPC authentication refreshes performed.",
		nil,
		nil,
	)

	nfsProceduresDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "nfs", "procedures"),
		"Number of NFS procedures invoked.",
		[]string{"version", "procedure"},
		nil,
	)
)

type nfsCollector struct{}

func init() {
	Factories["nfs"] = NewNfsCollector
}

func NewNfsCollector() (Collector, error) {
	return &nfsCollector{}, nil
}

func (c *nfsCollector) Update(ch chan<- prometheus.Metric) (err error) {
	statsFile := procFilePath("net/rpc/nfs")
	content, err := ioutil.ReadFile(statsFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Debugf("Not collecting NFS statistics, as %s does not exist: %s", statsFile)
			return nil
		}
		return err
	}

	for _, line := range strings.Split(string(content), "\n") {
		if fields := netLineRE.FindStringSubmatch(line); fields != nil {
			value, _ := strconv.ParseFloat(fields[1], 64)
			ch <- prometheus.MustNewConstMetric(
				nfsNetReadsDesc, prometheus.CounterValue,
				value, "udp")

			value, _ = strconv.ParseFloat(fields[2], 64)
			ch <- prometheus.MustNewConstMetric(
				nfsNetReadsDesc, prometheus.CounterValue,
				value, "tcp")

			value, _ = strconv.ParseFloat(fields[3], 64)
			ch <- prometheus.MustNewConstMetric(
				nfsNetConnectionsDesc, prometheus.CounterValue,
				value, "tcp")
		} else if fields := rpcLineRE.FindStringSubmatch(line); fields != nil {
			value, _ := strconv.ParseFloat(fields[1], 64)
			ch <- prometheus.MustNewConstMetric(
				nfsRpcOperationsDesc,
				prometheus.CounterValue, value)

			value, _ = strconv.ParseFloat(fields[2], 64)
			ch <- prometheus.MustNewConstMetric(
				nfsRpcRetransmissionsDesc,
				prometheus.CounterValue, value)

			value, _ = strconv.ParseFloat(fields[3], 64)
			ch <- prometheus.MustNewConstMetric(
				nfsRpcAuthenticationRefreshesDesc,
				prometheus.CounterValue, value)
		} else if fields := procLineRE.FindStringSubmatch(line); fields != nil {
			version := fields[1]
			for procedure, count := range strings.Split(fields[2], " ") {
				value, _ := strconv.ParseFloat(count, 64)
				ch <- prometheus.MustNewConstMetric(
					nfsProceduresDesc,
					prometheus.CounterValue,
					value,
					version,
					nfsProcedures[version][procedure])
			}
		} else if line != "" {
			return errors.New("Failed to parse line: " + line)
		}
	}
	return nil
}
