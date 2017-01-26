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

var zfsNotAvailableError = errors.New("ZFS / ZFS statistics are not available")

type zfsSysctl string

// Metrics

type zfsMetric struct {
	subsystem string    // The Prometheus subsystem name.
	name      string    // The Prometheus name of the metric.
	sysctl    zfsSysctl // The sysctl of the ZFS metric.
}

// Collector

func init() {
	Factories["zfs"] = NewZFSCollector
}

type zfsCollector struct {
	zfsMetrics        []zfsMetric
	linuxProcpathBase string
	linuxPathMap      map[string]string
}

func NewZFSCollector() (Collector, error) {
	var z zfsCollector
	z.linuxProcpathBase = "spl/kstat/zfs"
	z.linuxPathMap = map[string]string{
		"zfsArc":       "arcstats",
		"zfsDmuTx":     "dmu_tx",
		"zfsFm":        "fm",
		"zfsFetch":     "zfetchstats",
		"zfsVdevCache": "vdev_cache_stats",
		"zfsXuio":      "xuio_stats",
		"zfsZil":       "zil",
	}
	return &z, nil
}

func (c *zfsCollector) Update(ch chan<- prometheus.Metric) (err error) {
	for subsystem := range c.linuxPathMap {
		err = c.updateZfsStats(subsystem, ch)
		switch {
		case err == zfsNotAvailableError:
			log.Debug(err)
			return nil
		case err != nil:
			return err
		}
	}

	// Pool stats
	return c.updatePoolStats(ch)
}

func (s zfsSysctl) metricName() string {
	parts := strings.Split(string(s), ".")
	return parts[len(parts)-1]
}

func (c *zfsCollector) constSysctlMetric(subsystem string, sysctl zfsSysctl, value int) prometheus.Metric {
	metricName := sysctl.metricName()

	return prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, string(subsystem), metricName),
			string(sysctl),
			nil,
			nil,
		),
		prometheus.UntypedValue,
		float64(value),
	)
}
