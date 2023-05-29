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

//go:build linux && !nozfs
// +build linux,!nozfs

package collector

import (
	"errors"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

var errZFSNotAvailable = errors.New("ZFS / ZFS statistics are not available")

type zfsSysctl string

func init() {
	registerCollector("zfs", defaultEnabled, NewZFSCollector)
}

type zfsCollector struct {
	linuxProcpathBase    string
	linuxZpoolIoPath     string
	linuxZpoolObjsetPath string
	linuxZpoolStatePath  string
	linuxPathMap         map[string]string
	logger               log.Logger
}

// NewZFSCollector returns a new Collector exposing ZFS statistics.
func NewZFSCollector(logger log.Logger) (Collector, error) {
	return &zfsCollector{
		linuxProcpathBase:    "spl/kstat/zfs",
		linuxZpoolIoPath:     "/*/io",
		linuxZpoolObjsetPath: "/*/objset-*",
		linuxZpoolStatePath:  "/*/state",
		linuxPathMap: map[string]string{
			"zfs_abd":         "abdstats",
			"zfs_arc":         "arcstats",
			"zfs_dbuf":        "dbufstats",
			"zfs_dmu_tx":      "dmu_tx",
			"zfs_dnode":       "dnodestats",
			"zfs_fm":          "fm",
			"zfs_vdev_cache":  "vdev_cache_stats", // vdev_cache is deprecated
			"zfs_vdev_mirror": "vdev_mirror_stats",
			"zfs_xuio":        "xuio_stats", // no known consumers of the XUIO interface on Linux exist
			"zfs_zfetch":      "zfetchstats",
			"zfs_zil":         "zil",
		},
		logger: logger,
	}, nil
}

func (c *zfsCollector) Update(ch chan<- prometheus.Metric) error {

	if _, err := c.openProcFile(c.linuxProcpathBase); err != nil {
		if err == errZFSNotAvailable {
			level.Debug(c.logger).Log("err", err)
			return ErrNoData
		}
	}

	for subsystem := range c.linuxPathMap {
		if err := c.updateZfsStats(subsystem, ch); err != nil {
			if err == errZFSNotAvailable {
				level.Debug(c.logger).Log("err", err)
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

func (c *zfsCollector) constPoolObjsetMetric(poolName string, datasetName string, sysctl zfsSysctl, value uint64) prometheus.Metric {
	metricName := sysctl.metricName()

	return prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "zfs_zpool_dataset", metricName),
			string(sysctl),
			[]string{"zpool", "dataset"},
			nil,
		),
		prometheus.UntypedValue,
		float64(value),
		poolName,
		datasetName,
	)
}

func (c *zfsCollector) constPoolStateMetric(poolName string, stateName string, isActive uint64) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "zfs_zpool", "state"),
			"kstat.zfs.misc.state",
			[]string{"zpool", "state"},
			nil,
		),
		prometheus.GaugeValue,
		float64(isActive),
		poolName,
		stateName,
	)
}
