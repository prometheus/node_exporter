// Copyright 2021 The Prometheus Authors
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

//go:build !nonvme
// +build !nonvme

package collector

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/sysfs"
)

type nvmeCollector struct {
	fs                             sysfs.FS
	logger                         log.Logger
	namespaceInfo                  *prometheus.Desc
	namespaceCapacityBytes         *prometheus.Desc
	namespaceSizeBytes             *prometheus.Desc
	namespaceUsedBytes             *prometheus.Desc
	namespaceLogicalBlockSizeBytes *prometheus.Desc
	info                           *prometheus.Desc
}

func init() {
	registerCollector("nvme", defaultEnabled, NewNVMeCollector)
}

// NewNVMeCollector returns a new Collector exposing NVMe stats.
func NewNVMeCollector(logger log.Logger) (Collector, error) {
	fs, err := sysfs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}

	info := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "nvme", "info"),
		"Non-numeric data from /sys/class/nvme/<device>, value is always 1.",
		[]string{"device", "firmware_revision", "model", "serial", "state", "cntlid"},
		nil,
	)
	namespaceInfo := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "nvme", "namespace_info"),
		"Information about NVMe namespaces. Value is always 1",
		[]string{"device", "nsid", "ana_state"}, nil,
	)

	namespaceCapacityBytes := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "nvme", "namespace_capacity_bytes"),
		"Capacity of the NVMe namespace in bytes. Computed as namespace_size * namespace_logical_block_size",
		[]string{"device", "nsid"}, nil,
	)

	namespaceSizeBytes := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "nvme", "namespace_size_bytes"),
		"Size of the NVMe namespace in bytes. Available in /sys/class/nvme/<device>/<namespace>/size",
		[]string{"device", "nsid"}, nil,
	)

	namespaceUsedBytes := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "nvme", "namespace_used_bytes"),
		"Used space of the NVMe namespace in bytes. Available in /sys/class/nvme/<device>/<namespace>/nuse",
		[]string{"device", "nsid"}, nil,
	)

	namespaceLogicalBlockSizeBytes := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "nvme", "namespace_logical_block_size_bytes"),
		"Logical block size of the NVMe namespace in bytes. Usually 4Kb. Available in /sys/class/nvme/<device>/<namespace>/queue/logical_block_size",
		[]string{"device", "nsid"}, nil,
	)

	return &nvmeCollector{
		fs:                             fs,
		logger:                         logger,
		namespaceInfo:                  namespaceInfo,
		namespaceCapacityBytes:         namespaceCapacityBytes,
		namespaceSizeBytes:             namespaceSizeBytes,
		namespaceUsedBytes:             namespaceUsedBytes,
		namespaceLogicalBlockSizeBytes: namespaceLogicalBlockSizeBytes,
		info:                           info,
	}, nil
}

func (c *nvmeCollector) Update(ch chan<- prometheus.Metric) error {
	devices, err := c.fs.NVMeClass()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			level.Debug(c.logger).Log("msg", "nvme statistics not found, skipping")
			return ErrNoData
		}
		return fmt.Errorf("error obtaining NVMe class info: %w", err)
	}

	for _, device := range devices {
		infoValue := 1.0

		devicePath := filepath.Join(*sysPath, "class/nvme", device.Name)
		cntlid, err := readUintFromFile(filepath.Join(devicePath, "cntlid"))
		if err != nil {
			level.Debug(c.logger).Log("msg", "failed to read cntlid", "device", device.Name, "err", err)
		}
		ch <- prometheus.MustNewConstMetric(c.info, prometheus.GaugeValue, infoValue, device.Name, device.FirmwareRevision, device.Model, device.Serial, device.State, strconv.FormatUint(cntlid, 10))
		// Find namespace directories.
		namespacePaths, err := filepath.Glob(filepath.Join(devicePath, "nvme[0-9]*c[0-9]*n[0-9]*"))
		if err != nil {
			level.Error(c.logger).Log("msg", "failed to list NVMe namespaces", "device", device.Name, "err", err)
			continue
		}
		re := regexp.MustCompile(`nvme[0-9]+c[0-9]+n([0-9]+)`)

		for _, namespacePath := range namespacePaths {

			// Read namespace data.
			match := re.FindStringSubmatch(filepath.Base(namespacePath))
			if len(match) == 0 {
				continue
			}
			nsid := match[1]
			nuse, err := readUintFromFile(filepath.Join(namespacePath, "nuse"))
			if err != nil {
				level.Debug(c.logger).Log("msg", "failed to read nuse", "device", device.Name, "namespace", match[0], "err", err)
			}
			nsze, err := readUintFromFile(filepath.Join(namespacePath, "size"))
			if err != nil {
				level.Debug(c.logger).Log("msg", "failed to read size", "device", device.Name, "namespace", match[0], "err", err)
			}
			lbaSize, err := readUintFromFile(filepath.Join(namespacePath, "queue", "logical_block_size"))
			if err != nil {
				level.Debug(c.logger).Log("msg", "failed to read queue/logical_block_size", "device", device.Name, "namespace", match[0], "err", err)
			}
			ncap := nsze * lbaSize
			anaState := "unknown"
			anaStateSysfs, err := os.ReadFile(filepath.Join(namespacePath, "ana_state"))
			if err == nil {
				anaState = strings.TrimSpace(string(anaStateSysfs))
			} else {
				level.Debug(c.logger).Log("msg", "failed to read ana_state", "device", device.Name, "namespace", match[0], "err", err)
			}

			ch <- prometheus.MustNewConstMetric(
				c.namespaceInfo,
				prometheus.GaugeValue,
				1.0,
				device.Name,
				nsid,
				anaState,
			)

			ch <- prometheus.MustNewConstMetric(
				c.namespaceCapacityBytes,
				prometheus.GaugeValue,
				float64(ncap),
				device.Name,
				nsid,
			)

			ch <- prometheus.MustNewConstMetric(
				c.namespaceSizeBytes,
				prometheus.GaugeValue,
				float64(nsze),
				device.Name,
				nsid,
			)

			ch <- prometheus.MustNewConstMetric(
				c.namespaceUsedBytes,
				prometheus.GaugeValue,
				float64(nuse*lbaSize),
				device.Name,
				nsid,
			)

			ch <- prometheus.MustNewConstMetric(
				c.namespaceLogicalBlockSizeBytes,
				prometheus.GaugeValue,
				float64(lbaSize),
				device.Name,
				nsid,
			)
		}
	}

	return nil
}