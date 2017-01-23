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
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const (
	zfsProcpathBase      = "spl/kstat/zfs/"
	zfsArcstatsExt       = "arcstats"
	zfsFetchstatsExt     = "zfetchstats"
	zfsVdevCacheStatsExt = "vdev_cache_stats"
	zfsXuioStatsExt      = "xuio_stats"
	zfsZilExt            = "zil"
)

func (c *zfsCollector) openProcFile(path string) (file *os.File, err error) {
	file, err = os.Open(procFilePath(path))
	if err != nil {
		log.Debugf("Cannot open %q for reading. Is the kernel module loaded?", procFilePath(path))
		err = zfsNotAvailableError
	}
	return
}

func (c *zfsCollector) updateArcstats(ch chan<- prometheus.Metric) (err error) {
	file, err := c.openProcFile(filepath.Join(zfsProcpathBase, zfsArcstatsExt))
	if err != nil {
		return err
	}
	defer file.Close()

	return c.parseProcfsFile(file, zfsArcstatsExt, func(s zfsSysctl, v zfsMetricValue) {
		ch <- c.constSysctlMetric(arc, s, v)
	})
}

func (c *zfsCollector) updateZfetchstats(ch chan<- prometheus.Metric) (err error) {
	file, err := c.openProcFile(filepath.Join(zfsProcpathBase, zfsFetchstatsExt))
	if err != nil {
		return err
	}
	defer file.Close()

	return c.parseProcfsFile(file, zfsFetchstatsExt, func(s zfsSysctl, v zfsMetricValue) {
		ch <- c.constSysctlMetric(zfetch, s, v)
	})
}

func (c *zfsCollector) updateZil(ch chan<- prometheus.Metric) (err error) {
	file, err := c.openProcFile(filepath.Join(zfsProcpathBase, zfsZilExt))
	if err != nil {
		return err
	}
	defer file.Close()

	return c.parseProcfsFile(file, zfsZilExt, func(s zfsSysctl, v zfsMetricValue) {
		ch <- c.constSysctlMetric(zil, s, v)
	})
}

func (c *zfsCollector) updateVdevCacheStats(ch chan<- prometheus.Metric) (err error) {
	file, err := c.openProcFile(filepath.Join(zfsProcpathBase, zfsVdevCacheStatsExt))
	if err != nil {
		return err
	}
	defer file.Close()

	return c.parseProcfsFile(file, zfsVdevCacheStatsExt, func(s zfsSysctl, v zfsMetricValue) {
		ch <- c.constSysctlMetric(vdevCache, s, v)
	})
}

func (c *zfsCollector) updateXuioStats(ch chan<- prometheus.Metric) (err error) {
	file, err := c.openProcFile(filepath.Join(zfsProcpathBase, zfsXuioStatsExt))
	if err != nil {
		return err
	}
	defer file.Close()

	return c.parseProcfsFile(file, zfsXuioStatsExt, func(s zfsSysctl, v zfsMetricValue) {
		ch <- c.constSysctlMetric(xuio, s, v)
	})
}

func (c *zfsCollector) parseProcfsFile(reader io.Reader, fmt_ext string, handler func(zfsSysctl, zfsMetricValue)) (err error) {
	scanner := bufio.NewScanner(reader)

	parseLine := false
	for scanner.Scan() {

		parts := strings.Fields(scanner.Text())

		if !parseLine && len(parts) == 3 && parts[0] == "name" && parts[1] == "type" && parts[2] == "data" {
			// Start parsing from here.
			parseLine = true
			continue
		}

		if !parseLine || len(parts) < 3 {
			continue
		}

		key := fmt.Sprintf("kstat.zfs.misc.%s.%s", fmt_ext, parts[0])

		value, err := strconv.Atoi(parts[2])
		if err != nil {
			return fmt.Errorf("could not parse expected integer value for %q", key)
		}
		handler(zfsSysctl(key), zfsMetricValue(value))

	}
	if !parseLine {
		return fmt.Errorf("did not parse a single %q metric", fmt_ext)
	}

	return scanner.Err()
}

func (c *zfsCollector) updatePoolStats(ch chan<- prometheus.Metric) (err error) {
	return nil
}
