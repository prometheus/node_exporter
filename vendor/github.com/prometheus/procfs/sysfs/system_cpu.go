// Copyright 2018 The Prometheus Authors
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

// +build !windows

package sysfs

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/prometheus/procfs/internal/util"
)

// SystemCPUCpufreqStats contains stats from devices/system/cpu/cpu[0-9]*/cpufreq/...
type SystemCPUCpufreqStats struct {
	Name                     string
	CpuinfoCurrentFrequency  *uint64
	CpuinfoMinimumFrequency  *uint64
	CpuinfoMaximumFrequency  *uint64
	CpuinfoTransitionLatency *uint64
	ScalingCurrentFrequency  *uint64
	ScalingMinimumFrequency  *uint64
	ScalingMaximumFrequency  *uint64
	AvailableGovernors       string
	Driver                   string
	Govenor                  string
	RelatedCpus              string
	SetSpeed                 string
}

// TODO: Add topology support.

// TODO: Add thermal_throttle support.

// NewSystemCpufreq returns CPU frequency metrics for all CPUs.
func NewSystemCpufreq() ([]SystemCPUCpufreqStats, error) {
	fs, err := NewFS(DefaultMountPoint)
	if err != nil {
		return []SystemCPUCpufreqStats{}, err
	}

	return fs.NewSystemCpufreq()
}

// NewSystemCpufreq returns CPU frequency metrics for all CPUs.
func (fs FS) NewSystemCpufreq() ([]SystemCPUCpufreqStats, error) {
	var cpufreq = &SystemCPUCpufreqStats{}

	cpus, err := filepath.Glob(fs.Path("devices/system/cpu/cpu[0-9]*"))
	if err != nil {
		return []SystemCPUCpufreqStats{}, err
	}

	systemCpufreq := []SystemCPUCpufreqStats{}
	for _, cpu := range cpus {
		cpuName := filepath.Base(cpu)
		cpuNum := strings.TrimPrefix(cpuName, "cpu")

		cpuCpufreqPath := filepath.Join(cpu, "cpufreq")
		if _, err := os.Stat(cpuCpufreqPath); os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return []SystemCPUCpufreqStats{}, err
		}

		cpufreq, err = parseCpufreqCpuinfo(cpuCpufreqPath)
		if err != nil {
			return []SystemCPUCpufreqStats{}, err
		}
		cpufreq.Name = cpuNum
		systemCpufreq = append(systemCpufreq, *cpufreq)
	}

	return systemCpufreq, nil
}

func parseCpufreqCpuinfo(cpuPath string) (*SystemCPUCpufreqStats, error) {
	uintFiles := []string{
		"cpuinfo_cur_freq",
		"cpuinfo_max_freq",
		"cpuinfo_min_freq",
		"cpuinfo_transition_latency",
		"scaling_cur_freq",
		"scaling_max_freq",
		"scaling_min_freq",
	}
	uintOut := make([]*uint64, len(uintFiles))

	for i, f := range uintFiles {
		v, err := util.ReadUintFromFile(filepath.Join(cpuPath, f))
		if err != nil {
			if os.IsNotExist(err) || os.IsPermission(err) {
				continue
			}
			return &SystemCPUCpufreqStats{}, err
		}

		uintOut[i] = &v
	}

	stringFiles := []string{
		"scaling_available_governors",
		"scaling_driver",
		"scaling_governor",
		"related_cpus",
		"scaling_setspeed",
	}
	stringOut := make([]string, len(stringFiles))
	var err error

	for i, f := range stringFiles {
		stringOut[i], err = util.SysReadFile(filepath.Join(cpuPath, f))
		if err != nil {
			return &SystemCPUCpufreqStats{}, err
		}
	}

	return &SystemCPUCpufreqStats{
		CpuinfoCurrentFrequency:  uintOut[0],
		CpuinfoMaximumFrequency:  uintOut[1],
		CpuinfoMinimumFrequency:  uintOut[2],
		CpuinfoTransitionLatency: uintOut[3],
		ScalingCurrentFrequency:  uintOut[4],
		ScalingMaximumFrequency:  uintOut[5],
		ScalingMinimumFrequency:  uintOut[6],
		AvailableGovernors:       stringOut[0],
		Driver:                   stringOut[1],
		Govenor:                  stringOut[2],
		RelatedCpus:              stringOut[3],
		SetSpeed:                 stringOut[4],
	}, nil
}
