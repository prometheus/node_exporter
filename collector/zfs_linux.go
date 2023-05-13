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
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-kit/log/level"
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

var zfsPoolStatesName = []string{"online", "degraded", "faulted", "offline", "removed", "unavail", "suspended"}

func (c *zfsCollector) openProcFile(path string) (*os.File, error) {
	file, err := os.Open(procFilePath(path))
	if err != nil {
		// file not found error can occur if:
		// 1. zfs module is not loaded
		// 2. zfs version does not have the feature with metrics -- ok to ignore
		level.Debug(c.logger).Log("msg", "Cannot open file for reading", "path", procFilePath(path))
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

	return c.parseProcfsFile(file, c.linuxPathMap[subsystem], func(s zfsSysctl, v uint64) {
		ch <- c.constSysctlMetric(subsystem, s, v)
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
			level.Debug(c.logger).Log("msg", "Cannot open file for reading", "path", zpoolPath)
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
			level.Debug(c.logger).Log("msg", "Cannot open file for reading", "path", zpoolPath)
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
		level.Debug(c.logger).Log("msg", "No pool state files found")
		return nil
	}

	for _, zpoolPath := range zpoolStatePaths {
		file, err := os.Open(zpoolPath)
		if err != nil {
			// This file should exist, but there is a race where an exporting pool can remove the files. Ok to ignore.
			level.Debug(c.logger).Log("msg", "Cannot open file for reading", "path", zpoolPath)
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

func (c *zfsCollector) parseProcfsFile(reader io.Reader, fmtExt string, handler func(zfsSysctl, uint64)) error {
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
		if parts[1] == kstatDataUint64 || parts[1] == kstatDataInt64 {
			key := fmt.Sprintf("kstat.zfs.misc.%s.%s", fmtExt, parts[0])
			value, err := strconv.ParseUint(parts[2], 10, 64)
			if err != nil {
				return fmt.Errorf("could not parse expected integer value for %q", key)
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
		parts := strings.Fields(scanner.Text())

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
			datasetName = parts[2]
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
