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

//go:build !nonvmesubsystem

package collector

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/sysfs"
)

var nvmeControllerStates = []string{
	"live", "connecting", "resetting", "dead", "unknown",
}

func normalizeControllerState(raw string) string {
	switch raw {
	case "live", "connecting", "resetting", "dead":
		return raw
	case "deleting", "deleting (no IO)", "new":
		return raw
	default:
		return "unknown"
	}
}

type nvmeSubsystemCollector struct {
	fs     sysfs.FS
	logger *slog.Logger

	subsystemInfo      *prometheus.Desc
	subsystemPaths     *prometheus.Desc
	subsystemPathsLive *prometheus.Desc
	pathState          *prometheus.Desc
}

func init() {
	registerCollector("nvmesubsystem", defaultDisabled, NewNVMeSubsystemCollector)
}

// NewNVMeSubsystemCollector returns a new Collector exposing NVMe-oF subsystem
// path health from /sys/class/nvme-subsystem/.
func NewNVMeSubsystemCollector(logger *slog.Logger) (Collector, error) {
	const subsystem = "nvmesubsystem"

	fs, err := sysfs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}

	return &nvmeSubsystemCollector{
		fs:     fs,
		logger: logger,
		subsystemInfo: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "info"),
			"Non-numeric information about an NVMe subsystem.",
			[]string{"subsystem", "nqn", "model", "serial", "iopolicy"}, nil,
		),
		subsystemPaths: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "paths"),
			"Number of controller paths for an NVMe subsystem.",
			[]string{"subsystem"}, nil,
		),
		subsystemPathsLive: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "paths_live"),
			"Number of controller paths in live state for an NVMe subsystem.",
			[]string{"subsystem"}, nil,
		),
		pathState: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "path_state"),
			"Current NVMe controller path state (1 for the current state, 0 for all others).",
			[]string{"subsystem", "controller", "transport", "state"}, nil,
		),
	}, nil
}

func (c *nvmeSubsystemCollector) Update(ch chan<- prometheus.Metric) error {
	subsystems, err := c.fs.NVMeSubsystemClass()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) || errors.Is(err, os.ErrPermission) {
			c.logger.Debug("Could not read NVMe subsystem info", "err", err)
			return ErrNoData
		}
		return fmt.Errorf("failed to scan NVMe subsystems: %w", err)
	}

	for _, subsys := range subsystems {
		ch <- prometheus.MustNewConstMetric(c.subsystemInfo, prometheus.GaugeValue, 1,
			subsys.Name, subsys.NQN, subsys.Model, subsys.Serial, subsys.IOPolicy)

		total := float64(len(subsys.Controllers))
		var live float64
		for _, ctrl := range subsys.Controllers {
			state := normalizeControllerState(ctrl.State)
			if state == "live" {
				live++
			}

			for _, s := range nvmeControllerStates {
				val := 0.0
				if s == state {
					val = 1.0
				}
				ch <- prometheus.MustNewConstMetric(c.pathState, prometheus.GaugeValue, val,
					subsys.Name, ctrl.Name, ctrl.Transport, s)
			}
		}

		ch <- prometheus.MustNewConstMetric(c.subsystemPaths, prometheus.GaugeValue, total, subsys.Name)
		ch <- prometheus.MustNewConstMetric(c.subsystemPathsLive, prometheus.GaugeValue, live, subsys.Name)
	}

	return nil
}
