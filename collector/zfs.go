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

// +build linux
// +build !nozfs

package collector

import (
	"errors"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

var errZFSNotAvailable = errors.New("ZFS / ZFS statistics are not available")

type zfsSysctl string

func init() {
	registerCollector("zfs", defaultEnabled, NewZFSCollector)
}

type zfsCollector struct {
	linuxProcpathBase string
	linuxZpoolIoPath  string
	linuxPathMap      map[string]string
}

// NewZFSCollector returns a new Collector exposing ZFS statistics.
func NewZFSCollector() (Collector, error) {
	return &zfsCollector{
		linuxProcpathBase: "spl/kstat/zfs",
		linuxZpoolIoPath:  "/*/io",
		linuxPathMap: map[string]string{
			"zfs_abd":         "abdstats",
			"zfs_arc":         "arcstats",
			"zfs_dbuf":        "dbuf_stats",
			"zfs_dmu_tx":      "dmu_tx",
			"zfs_dnode":       "dnodestats",
			"zfs_fm":          "fm",
			"zfs_vdev_cache":  "vdev_cache_stats", // vdev_cache is deprecated
			"zfs_vdev_mirror": "vdev_mirror_stats",
			"zfs_xuio":        "xuio_stats", // no known consumers of the XUIO interface on Linux exist
			"zfs_zfetch":      "zfetchstats",
			"zfs_zil":         "zil",
		},
	}, nil
}

func (c *zfsCollector) Update(ch chan<- prometheus.Metric) error {
	for subsystem := range c.linuxPathMap {
		if err := c.updateZfsStats(subsystem, ch); err != nil {
			if err == errZFSNotAvailable {
				log.Debug(err)
				// ZFS /proc files are added as new features to ZFS arrive, it is ok to continue
				continue
			}
			return err
		}
	}

	// Pool stats
	return c.updatePoolStats(ch)
}

func (s zfsSysctl) metricName() string {
	parts := strings.Split(string(s), ".")
	return strings.Replace(parts[len(parts)-1], "-", "_", -1)
}

func (c *zfsCollector) constSysctlMetric(subsystem string, sysctl zfsSysctl, value uint64) prometheus.Metric {
	metricName := sysctl.metricName()

	return prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, metricName),
			string(sysctl),
			nil,
			nil,
		),
		prometheus.UntypedValue,
		float64(value),
	)
}

func (c *zfsCollector) constPoolMetric(poolName string, sysctl zfsSysctl, value uint64) prometheus.Metric {
	metricName := sysctl.metricName()

	return prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "zfs_zpool", metricName),
			string(sysctl),
			[]string{"zpool"},
			nil,
		),
		prometheus.UntypedValue,
		float64(value),
		poolName,
	)
}
