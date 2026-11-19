// Copyright The Prometheus Authors
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

//go:build !nokexec

package collector

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/prometheus/client_golang/prometheus"
)

const kexecCollectorSubsystem = "kexec"

var (
	kexecLoadedDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, kexecCollectorSubsystem, "loaded"),
		"Whether a kexec kernel is loaded and ready for a kexec reboot, 1 is true, 0 is false.",
		nil, nil,
	)
	kexecCrashLoadedDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, kexecCollectorSubsystem, "crash_loaded"),
		"Whether a crash (kdump) kernel is loaded, 1 is true, 0 is false.",
		nil, nil,
	)
	kexecCrashSizeDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, kexecCollectorSubsystem, "crash_size_bytes"),
		"Size in bytes of the memory reserved for the crash (kdump) kernel, 0 means no reservation.",
		nil, nil,
	)

	kexecFiles = []kexecFile{
		{"kexec_loaded", kexecLoadedDesc},
		{"kexec_crash_loaded", kexecCrashLoadedDesc},
		{"kexec_crash_size", kexecCrashSizeDesc},
	}
)

type kexecFile struct {
	name string
	desc *prometheus.Desc
}

type kexecCollector struct {
	logger *slog.Logger
}

func init() {
	registerCollector(kexecCollectorSubsystem, defaultDisabled, NewKexecCollector)
}

// NewKexecCollector returns a collector exposing kexec and crash kernel (kdump) state from /sys/kernel
func NewKexecCollector(logger *slog.Logger) (Collector, error) {
	return &kexecCollector{logger: logger}, nil
}

func (c *kexecCollector) Update(ch chan<- prometheus.Metric) error {
	emitted := 0
	for _, f := range kexecFiles {
		path := sysFilePath(filepath.Join("kernel", f.name))
		val, err := readUintFromFile(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				c.logger.Debug("kexec file not available, skipping", "file", f.name, "err", err)
				continue
			}
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		ch <- prometheus.MustNewConstMetric(f.desc, prometheus.GaugeValue, float64(val))
		emitted++
	}

	if emitted == 0 {
		return ErrNoData
	}
	return nil
}
