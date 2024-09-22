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
	"log/slog"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/sysfs"
)

type nvmeCollector struct {
	fs     sysfs.FS
	logger *slog.Logger
}

func init() {
	registerCollector("nvme", defaultEnabled, NewNVMeCollector)
}

// NewNVMeCollector returns a new Collector exposing NVMe stats.
func NewNVMeCollector(logger *slog.Logger) (Collector, error) {
	fs, err := sysfs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}

	return &nvmeCollector{
		fs:     fs,
		logger: logger,
	}, nil
}

func (c *nvmeCollector) Update(ch chan<- prometheus.Metric) error {
	devices, err := c.fs.NVMeClass()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			c.logger.Debug("nvme statistics not found, skipping")
			return ErrNoData
		}
		return fmt.Errorf("error obtaining NVMe class info: %w", err)
	}

	for _, device := range devices {
		infoDesc := prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "nvme", "info"),
			"Non-numeric data from /sys/class/nvme/<device>, value is always 1.",
			[]string{"device", "firmware_revision", "model", "serial", "state"},
			nil,
		)
		infoValue := 1.0
		ch <- prometheus.MustNewConstMetric(infoDesc, prometheus.GaugeValue, infoValue, device.Name, device.FirmwareRevision, device.Model, device.Serial, device.State)
	}

	return nil
}
