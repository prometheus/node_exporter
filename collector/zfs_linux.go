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

//go:build !nozfs
// +build !nozfs

package collector

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

// constants from https://github.com/zfsonlinux/zfs/blob/master/lib/libspl/include/sys/kstat.h
// kept as strings for comparison thus avoiding conversion to int
const (
	// kstatDataChar   = "0"
	// kstatDataInt32  = "1"
	// kstatDataUint32 = "2"
	kstatDataInt64  = "3"
	kstatDataUint64 = "4"
	// kstatDataLong   = "5"
	// kstatDataUlong  = "6"
	// kstatDataString = "7"
)

var (
	errZFSNotAvailable = errors.New("ZFS / ZFS statistics are not available")

	zfsPoolStatesName = [...]string{"online", "degraded", "faulted", "offline", "removed", "unavail", "suspended"}
)

type zfsCollector struct {
	linuxProcpathBase    string
	linuxZpoolIoPath     string
	linuxZpoolObjsetPath string
	linuxZpoolStatePath  string
	linuxPathMap         map[string]string
	logger               *slog.Logger
}

// NewZFSCollector returns a new Collector exposing ZFS statistics.
func NewZFSCollector(logger *slog.Logger) (Collector, error) {
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
			c.logger.Debug(err.Error())
			return ErrNoData
		}
	}

	for subsystem := range c.linuxPathMap {
		if err := c.updateZfsStats(subsystem, ch); err != nil {
			if err == errZFSNotAvailable {
				c.logger.Debug(err.Error())
				// ZFS /proc files are added as new features to ZFS arrive, it is ok to continue
				continue
			}
			return err
		}
	}

	// Pool stats
	return c.updatePoolStats(ch)
}

func (c *zfsCollector) openProcFile(path string) (*os.File, error) {
	file, err := os.Open(procFilePath(path))
	if err != nil {
		// file not found error can occur if:
		// 1. zfs module is not loaded
		// 2. zfs version does not have the feature with metrics -- ok to ignore
		c.logger.Debug("Cannot open file for reading", "path", procFilePath(path))
		return nil, errZFSNotAvailable
	}
	return file, nil
}

func (c *zfsCollector) updateZfsStats(subsystem string, ch chan<- prometheus.Metric) error {
	file, err := c.openProcFile(filepath.Join(c.linuxProcpathBase, c.linuxPathMap[subsystem]))
	if err != nil {
		return err
	}
	defer file.Close()

	return c.parseProcfsFile(file, c.linuxPathMap[subsystem], func(s zfsSysctl, v interface{}) {
		var valueAsFloat64 float64
		switch value := v.(type) {
		case int64:
			valueAsFloat64 = float64(value)
		case uint64:
			valueAsFloat64 = float64(value)
		}
		ch <- c.constSysctlMetric(subsystem, s, valueAsFloat64)
	})
}

func (c *zfsCollector) updatePoolStats(ch chan<- prometheus.Metric) error {
	zpoolPaths, err := filepath.Glob(procFilePath(filepath.Join(c.linuxProcpathBase, c.linuxZpoolIoPath)))
	if err != nil {
		return err
	}

	for _, zpoolPath := range zpoolPaths {
		file, err := os.Open(zpoolPath)
		if err != nil {
			// this file should exist, but there is a race where an exporting pool can remove the files -- ok to ignore
			c.logger.Debug("Cannot open file for reading", "path", zpoolPath)
			return errZFSNotAvailable
		}

		err = c.parsePoolProcfsFile(file, zpoolPath, func(poolName string, s zfsSysctl, v uint64) {
			ch <- c.constPoolMetric(poolName, s, v)
		})
		file.Close()
		if err != nil {
			return err
		}
	}

	zpoolObjsetPaths, err := filepath.Glob(procFilePath(filepath.Join(c.linuxProcpathBase, c.linuxZpoolObjsetPath)))
	if err != nil {
		return err
	}

	for _, zpoolPath := range zpoolObjsetPaths {
		file, err := os.Open(zpoolPath)
		if err != nil {
			// This file should exist, but there is a race where an exporting pool can remove the files. Ok to ignore.
			c.logger.Debug("Cannot open file for reading", "path", zpoolPath)
			return errZFSNotAvailable
		}

		err = c.parsePoolObjsetFile(file, zpoolPath, func(poolName string, datasetName string, s zfsSysctl, v uint64) {
			ch <- c.constPoolObjsetMetric(poolName, datasetName, s, v)
		})
		file.Close()
		if err != nil {
			return err
		}
	}

	zpoolStatePaths, err := filepath.Glob(procFilePath(filepath.Join(c.linuxProcpathBase, c.linuxZpoolStatePath)))
	if err != nil {
		return err
	}

	if zpoolStatePaths == nil {
		c.logger.Debug("No pool state files found")
		return nil
	}

	for _, zpoolPath := range zpoolStatePaths {
		file, err := os.Open(zpoolPath)
		if err != nil {
			// This file should exist, but there is a race where an exporting pool can remove the files. Ok to ignore.
			c.logger.Debug("Cannot open file for reading", "path", zpoolPath)
			return errZFSNotAvailable
		}

		err = c.parsePoolStateFile(file, zpoolPath, func(poolName string, stateName string, isActive uint64) {
			ch <- c.constPoolStateMetric(poolName, stateName, isActive)
		})

		file.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *zfsCollector) parseProcfsFile(reader io.Reader, fmtExt string, handler func(zfsSysctl, interface{})) error {
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

		// kstat data type (column 2) should be KSTAT_DATA_UINT64, otherwise ignore
		// TODO: when other KSTAT_DATA_* types arrive, much of this will need to be restructured
		key := fmt.Sprintf("kstat.zfs.misc.%s.%s", fmtExt, parts[0])
		switch parts[1] {
		case kstatDataUint64:
			value, err := strconv.ParseUint(parts[2], 10, 64)
			if err != nil {
				return fmt.Errorf("could not parse expected unsigned integer value for %q: %w", key, err)
			}
			handler(zfsSysctl(key), value)
		case kstatDataInt64:
			value, err := strconv.ParseInt(parts[2], 10, 64)
			if err != nil {
				return fmt.Errorf("could not parse expected signed integer value for %q: %w", key, err)
			}
			handler(zfsSysctl(key), value)
		}
	}
	if !parseLine {
		return fmt.Errorf("did not parse a single %q metric", fmtExt)
	}

	return scanner.Err()
}

func (c *zfsCollector) parsePoolProcfsFile(reader io.Reader, zpoolPath string, handler func(string, zfsSysctl, uint64)) error {
	scanner := bufio.NewScanner(reader)

	parseLine := false
	var fields []string
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())

		if !parseLine && len(line) >= 12 && line[0] == "nread" {
			//Start parsing from here.
			parseLine = true
			fields = make([]string, len(line))
			copy(fields, line)
			continue
		}
		if !parseLine {
			continue
		}

		zpoolPathElements := strings.Split(zpoolPath, "/")
		pathLen := len(zpoolPathElements)
		if pathLen < 2 {
			return fmt.Errorf("zpool path did not return at least two elements")
		}
		zpoolName := zpoolPathElements[pathLen-2]
		zpoolFile := zpoolPathElements[pathLen-1]

		for i, field := range fields {
			key := fmt.Sprintf("kstat.zfs.misc.%s.%s", zpoolFile, field)

			value, err := strconv.ParseUint(line[i], 10, 64)
			if err != nil {
				return fmt.Errorf("could not parse expected integer value for %q: %w", key, err)
			}
			handler(zpoolName, zfsSysctl(key), value)
		}
	}

	return scanner.Err()
}

func (c *zfsCollector) parsePoolObjsetFile(reader io.Reader, zpoolPath string, handler func(string, string, zfsSysctl, uint64)) error {
	scanner := bufio.NewScanner(reader)

	parseLine := false
	var zpoolName, datasetName string
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)

		if !parseLine && len(parts) == 3 && parts[0] == "name" && parts[1] == "type" && parts[2] == "data" {
			parseLine = true
			continue
		}

		if !parseLine || len(parts) < 3 {
			continue
		}
		if parts[0] == "dataset_name" {
			zpoolPathElements := strings.Split(zpoolPath, "/")
			pathLen := len(zpoolPathElements)
			zpoolName = zpoolPathElements[pathLen-2]
			datasetName = line[strings.Index(line, parts[2]):]
			continue
		}

		if parts[1] == kstatDataUint64 {
			key := fmt.Sprintf("kstat.zfs.misc.objset.%s", parts[0])
			value, err := strconv.ParseUint(parts[2], 10, 64)
			if err != nil {
				return fmt.Errorf("could not parse expected integer value for %q", key)
			}
			handler(zpoolName, datasetName, zfsSysctl(key), value)
		}
	}
	if !parseLine {
		return fmt.Errorf("did not parse a single %s %s metric", zpoolName, datasetName)
	}

	return scanner.Err()
}

func (c *zfsCollector) parsePoolStateFile(reader io.Reader, zpoolPath string, handler func(string, string, uint64)) error {
	scanner := bufio.NewScanner(reader)
	scanner.Scan()

	actualStateName, err := scanner.Text(), scanner.Err()
	if err != nil {
		return err
	}

	actualStateName = strings.ToLower(actualStateName)

	zpoolPathElements := strings.Split(zpoolPath, "/")
	pathLen := len(zpoolPathElements)
	if pathLen < 2 {
		return fmt.Errorf("zpool path did not return at least two elements")
	}

	zpoolName := zpoolPathElements[pathLen-2]

	for _, stateName := range zfsPoolStatesName {
		isActive := uint64(0)

		if actualStateName == stateName {
			isActive = 1
		}

		handler(zpoolName, stateName, isActive)
	}

	return nil
}

func (c *zfsCollector) constSysctlMetric(subsystem string, sysctl zfsSysctl, value float64) prometheus.Metric {
	metricName := sysctl.metricName()

	return prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, metricName),
			string(sysctl),
			nil,
			nil,
		),
		prometheus.UntypedValue,
		value,
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

type zfsSysctl string

func (s zfsSysctl) metricName() string {
	parts := strings.Split(string(s), ".")
	return strings.Replace(parts[len(parts)-1], "-", "_", -1)
}
