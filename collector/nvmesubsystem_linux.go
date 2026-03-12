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
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type nvmeSubsystemCollector struct {
	logger         *slog.Logger
	scanSubsystems func() ([]nvmeSubsystem, error)

	subsystemInfo       *prometheus.Desc
	subsystemPathsTotal *prometheus.Desc
	subsystemPathsLive  *prometheus.Desc
	pathState           *prometheus.Desc
}

type nvmeSubsystem struct {
	Name        string
	NQN         string
	Model       string
	Serial      string
	IOPolicy    string
	Controllers []nvmeController
}

type nvmeController struct {
	Name      string
	State     string
	Transport string
	Address   string
}

var (
	nvmeControllerRE = regexp.MustCompile(`^nvme\d+$`)

	nvmeControllerStates = []string{
		"live", "connecting", "resetting", "dead", "unknown",
	}
)

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

func init() {
	registerCollector("nvmesubsystem", defaultDisabled, NewNVMeSubsystemCollector)
}

// NewNVMeSubsystemCollector returns a new Collector exposing NVMe-oF subsystem
// path health from /sys/class/nvme-subsystem/.
func NewNVMeSubsystemCollector(logger *slog.Logger) (Collector, error) {
	const subsystem = "nvmesubsystem"

	c := &nvmeSubsystemCollector{
		logger: logger,
		subsystemInfo: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "info"),
			"Non-numeric information about an NVMe subsystem.",
			[]string{"subsystem", "nqn", "model", "serial", "iopolicy"}, nil,
		),
		subsystemPathsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "paths_total"),
			"Total number of controller paths for an NVMe subsystem.",
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
	}

	c.scanSubsystems = func() ([]nvmeSubsystem, error) {
		return scanNVMeSubsystems(*sysPath)
	}

	return c, nil
}

func (c *nvmeSubsystemCollector) Update(ch chan<- prometheus.Metric) error {
	subsystems, err := c.scanSubsystems()
	if err != nil {
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

		ch <- prometheus.MustNewConstMetric(c.subsystemPathsTotal, prometheus.GaugeValue, total, subsys.Name)
		ch <- prometheus.MustNewConstMetric(c.subsystemPathsLive, prometheus.GaugeValue, live, subsys.Name)
	}

	return nil
}

func scanNVMeSubsystems(sysfsBase string) ([]nvmeSubsystem, error) {
	subsysBase := filepath.Join(sysfsBase, "class", "nvme-subsystem")

	entries, err := os.ReadDir(subsysBase)
	if err != nil {
		return nil, err
	}

	var subsystems []nvmeSubsystem
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), "nvme-subsys") {
			continue
		}
		subsysPath := filepath.Join(subsysBase, entry.Name())
		subsys, err := parseNVMeSubsystem(entry.Name(), subsysPath)
		if err != nil {
			continue
		}
		subsystems = append(subsystems, *subsys)
	}

	return subsystems, nil
}

func parseNVMeSubsystem(name, path string) (*nvmeSubsystem, error) {
	subsys := &nvmeSubsystem{Name: name}

	subsys.NQN = readSysfsString(filepath.Join(path, "subsysnqn"))
	subsys.Model = readSysfsString(filepath.Join(path, "model"))
	subsys.Serial = readSysfsString(filepath.Join(path, "serial"))
	subsys.IOPolicy = readSysfsString(filepath.Join(path, "iopolicy"))

	entries, err := os.ReadDir(path)
	if err != nil {
		return subsys, nil
	}

	for _, entry := range entries {
		if !nvmeControllerRE.MatchString(entry.Name()) {
			continue
		}
		ctrlPath := filepath.Join(path, entry.Name())
		ctrl := nvmeController{
			Name:      entry.Name(),
			State:     readSysfsString(filepath.Join(ctrlPath, "state")),
			Transport: readSysfsString(filepath.Join(ctrlPath, "transport")),
			Address:   readSysfsString(filepath.Join(ctrlPath, "address")),
		}
		subsys.Controllers = append(subsys.Controllers, ctrl)
	}

	return subsys, nil
}

func readSysfsString(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}
