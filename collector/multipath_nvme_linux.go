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

//go:build !nomultipath

package collector

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

// nvmeSubsystem represents an NVMe subsystem from /sys/class/nvme-subsystem/.
type nvmeSubsystem struct {
	Name        string // nvme-subsys0
	NQN         string // subsysnqn — unique identifier for the storage target
	Model       string
	Serial      string
	IOPolicy    string // numa, round-robin, queue-depth
	Controllers []nvmeController
}

// nvmeController represents a controller (path) within an NVMe subsystem.
type nvmeController struct {
	Name      string // nvme0
	State     string // live, connecting, resetting, deleting, dead
	Transport string // fc, tcp, rdma, pcie, loop
	Address   string // transport address (traddr=...,host_traddr=...)
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

// scanNVMeSubsystems reads /sys/class/nvme-subsystem/ to discover NVMe
// subsystems and their controllers (paths).
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

// emitNVMeSubsystemMetrics emits metrics for all NVMe subsystems.
func (c *multipathCollector) emitNVMeSubsystemMetrics(ch chan<- prometheus.Metric, subsystems []nvmeSubsystem) {
	for _, subsys := range subsystems {
		ch <- prometheus.MustNewConstMetric(c.nvmeSubsystemInfo, prometheus.GaugeValue, 1,
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
				ch <- prometheus.MustNewConstMetric(c.nvmePathState, prometheus.GaugeValue, val,
					subsys.Name, ctrl.Name, ctrl.Transport, s)
			}
		}

		ch <- prometheus.MustNewConstMetric(c.nvmeSubsystemPathsTotal, prometheus.GaugeValue, total, subsys.Name)
		ch <- prometheus.MustNewConstMetric(c.nvmeSubsystemPathsLive, prometheus.GaugeValue, live, subsys.Name)
	}
}
