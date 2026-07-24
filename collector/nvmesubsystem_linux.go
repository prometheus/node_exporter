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
	"sort"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/sysfs"
)

// nvmeControllerStateMap maps each kernel-reported controller state to
// its metric label value. States not in this map are reported as "unknown".
var nvmeControllerStateMap = map[string]string{
	"live":             "live",
	"connecting":       "connecting",
	"resetting":        "resetting",
	"dead":             "dead",
	"deleting":         "deleting",
	"deleting (no IO)": "deleting",
	"new":              "new",
}

// nvmeMetricStates is the sorted, unique set of metric-level states
// derived from nvmeControllerStateMap values plus "unknown".
var nvmeMetricStates = func() []string {
	seen := map[string]struct{}{"unknown": {}}
	for _, v := range nvmeControllerStateMap {
		seen[v] = struct{}{}
	}
	states := make([]string, 0, len(seen))
	for s := range seen {
		states = append(states, s)
	}
	sort.Strings(states)
	return states
}()

func normalizeControllerState(raw string) string {
	if s, ok := nvmeControllerStateMap[raw]; ok {
		return s
	}
	return "unknown"
}

func init() {
	registerCollector("nvmesubsystem", defaultDisabled, NewNVMeSubsystemCollector)
}

var (
	nvmesubsystemInfo = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "nvmesubsystem", "info"),
		"Non-numeric information about an NVMe subsystem.",
		[]string{"subsystem", "nqn", "model", "serial", "iopolicy"}, nil,
	)
	nvmesubsystemNamespaceInfo = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "nvmesubsystem", "namespace_info"),
		"Maps an NVMe namespace block device to its subsystem.",
		[]string{"subsystem", "device"}, nil,
	)
	nvmesubsystemPaths = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "nvmesubsystem", "paths"),
		"Number of controller paths for an NVMe subsystem.",
		[]string{"subsystem"}, nil,
	)
	nvmesubsystemPathsLive = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "nvmesubsystem", "paths_live"),
		"Number of controller paths in live state for an NVMe subsystem.",
		[]string{"subsystem"}, nil,
	)
	nvmesubsystemPathState = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "nvmesubsystem", "path_state"),
		"Current NVMe controller path state (1 for the current state, 0 for all others).",
		[]string{"subsystem", "controller", "transport", "state"}, nil,
	)
)

type nvmeSubsystemCollector struct {
	fs     sysfs.FS
	logger *slog.Logger
}

// NewNVMeSubsystemCollector returns a new Collector exposing NVMe-oF subsystem
// path health from /sys/class/nvme-subsystem/.
func NewNVMeSubsystemCollector(logger *slog.Logger) (Collector, error) {
	fs, err := sysfs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}

	return &nvmeSubsystemCollector{
		fs:     fs,
		logger: logger,
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
		ch <- prometheus.MustNewConstMetric(nvmesubsystemInfo, prometheus.GaugeValue, 1,
			subsys.Name, subsys.NQN, subsys.Model, subsys.Serial, subsys.IOPolicy)

		for _, ns := range subsys.Namespaces {
			ch <- prometheus.MustNewConstMetric(nvmesubsystemNamespaceInfo, prometheus.GaugeValue, 1,
				subsys.Name, ns)
		}

		total := float64(len(subsys.Controllers))
		var live float64
		for _, ctrl := range subsys.Controllers {
			state := normalizeControllerState(ctrl.State)
			if state == "live" {
				live++
			}

			for _, s := range nvmeMetricStates {
				val := 0.0
				if s == state {
					val = 1.0
				}
				ch <- prometheus.MustNewConstMetric(nvmesubsystemPathState, prometheus.GaugeValue, val,
					subsys.Name, ctrl.Name, ctrl.Transport, s)
			}
		}

		ch <- prometheus.MustNewConstMetric(nvmesubsystemPaths, prometheus.GaugeValue, total, subsys.Name)
		ch <- prometheus.MustNewConstMetric(nvmesubsystemPathsLive, prometheus.GaugeValue, live, subsys.Name)
	}

	return nil
}
