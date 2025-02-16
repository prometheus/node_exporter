// Copyright 2023 The Prometheus Authors
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

//go:build !nocpu
// +build !nocpu

package collector

import (
	"errors"
	"log/slog"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unsafe"

	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sys/unix"

	"howett.net/plist"
)

const (
	_IOC_OUT        = uint(0x40000000)
	_IOC_IN         = uint(0x80000000)
	_IOC_INOUT      = (_IOC_IN | _IOC_OUT)
	_IOCPARM_MASK   = uint(0x1fff)
	_IOCPARM_SHIFT  = uint(16)
	_IOCGROUP_SHIFT = uint(8)
)

type clockinfo struct {
	hz     int32 // clock frequency
	tick   int32 // micro-seconds per hz tick
	spare  int32
	stathz int32 // statistics clock frequency
	profhz int32 // profiling clock frequency
}

type cputime struct {
	user float64
	nice float64
	sys  float64
	intr float64
	idle float64
}

type plistref struct {
	pref_plist unsafe.Pointer
	pref_len   uint
}

type sysmonValues struct {
	CurValue    int    `plist:"cur-value"`
	Description string `plist:"description"`
	State       string `plist:"state"`
	Type        string `plist:"type"`
}

type sysmonProperty []sysmonValues

type sysmonProperties map[string]sysmonProperty

func _IOC(inout uint, group byte, num uint, len uintptr) uint {
	return ((inout) | ((uint(len) & _IOCPARM_MASK) << _IOCPARM_SHIFT) | (uint(group) << _IOCGROUP_SHIFT) | (num))
}

func _IOWR(group byte, num uint, len uintptr) uint {
	return _IOC(_IOC_INOUT, group, num, len)
}

func ioctl(fd int, nr uint, typ byte, size uintptr, retptr unsafe.Pointer) error {
	_, _, errno := unix.Syscall(
		unix.SYS_IOCTL,
		uintptr(fd),
		uintptr(_IOWR(typ, nr, size)),
		uintptr(retptr),
	)
	if errno != 0 {
		return errno
	}
	return nil
}

func readSysmonProperties() (sysmonProperties, error) {
	fd, err := unix.Open(rootfsFilePath("/dev/sysmon"), unix.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer unix.Close(fd)

	var retptr plistref

	if err = ioctl(fd, 0, 'E', unsafe.Sizeof(retptr), unsafe.Pointer(&retptr)); err != nil {
		return nil, err
	}
	defer unix.Syscall(unix.SYS_MUNMAP, uintptr(retptr.pref_plist), uintptr(retptr.pref_len), uintptr(0))
	bytes := unsafe.Slice((*byte)(unsafe.Pointer(retptr.pref_plist)), retptr.pref_len-1)

	var props sysmonProperties
	if _, err = plist.Unmarshal(bytes, &props); err != nil {
		return nil, err
	}
	return props, nil
}

func sortFilterSysmonProperties(props sysmonProperties, prefix string) []string {
	var keys []string
	for key := range props {
		if !strings.HasPrefix(key, prefix) {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func convertTemperatures(prop sysmonProperty, res map[int]float64) error {

	for _, val := range prop {
		if val.State == "invalid" || val.State == "unknown" || val.State == "" {
			continue
		}

		re := regexp.MustCompile("^cpu([0-9]+) temperature$")
		core := re.FindStringSubmatch(val.Description)[1]
		ncore, _ := strconv.Atoi(core)
		temperature := ((float64(uint64(val.CurValue))) / 1000000) - 273.15
		res[ncore] = temperature
	}
	return nil
}

func getCPUTemperatures() (map[int]float64, error) {

	res := make(map[int]float64)

	// Read all properties
	props, err := readSysmonProperties()
	if err != nil {
		return res, err
	}

	keys := sortFilterSysmonProperties(props, "coretemp")
	for idx, _ := range keys {
		convertTemperatures(props[keys[idx]], res)
	}

	return res, nil
}

func getCPUTimes() ([]cputime, error) {
	const states = 5

	clockb, err := unix.SysctlRaw("kern.clockrate")
	if err != nil {
		return nil, err
	}
	clock := *(*clockinfo)(unsafe.Pointer(&clockb[0]))

	var cpufreq float64
	if clock.stathz > 0 {
		cpufreq = float64(clock.stathz)
	} else {
		cpufreq = float64(clock.hz)
	}

	ncpusb, err := unix.SysctlRaw("hw.ncpu")
	if err != nil {
		return nil, err
	}
	ncpus := int(*(*uint32)(unsafe.Pointer(&ncpusb[0])))

	if ncpus < 1 {
		return nil, errors.New("Invalid cpu number")
	}

	var times []float64
	for ncpu := 0; ncpu < ncpus; ncpu++ {
		cpb, err := unix.SysctlRaw("kern.cp_time", ncpu)
		if err != nil {
			return nil, err
		}
		for len(cpb) >= int(unsafe.Sizeof(uint64(0))) {
			t := *(*uint64)(unsafe.Pointer(&cpb[0]))
			times = append(times, float64(t)/cpufreq)
			cpb = cpb[unsafe.Sizeof(uint64(0)):]
		}
	}

	cpus := make([]cputime, len(times)/states)
	for i := 0; i < len(times); i += states {
		cpu := &cpus[i/states]
		cpu.user = times[i]
		cpu.nice = times[i+1]
		cpu.sys = times[i+2]
		cpu.intr = times[i+3]
		cpu.idle = times[i+4]
	}
	return cpus, nil
}

type statCollector struct {
	cpu    typedDesc
	temp   typedDesc
	logger *slog.Logger
}

func init() {
	registerCollector("cpu", defaultEnabled, NewStatCollector)
}

// NewStatCollector returns a new Collector exposing CPU stats.
func NewStatCollector(logger *slog.Logger) (Collector, error) {
	return &statCollector{
		cpu: typedDesc{nodeCPUSecondsDesc, prometheus.CounterValue},
		temp: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "temperature_celsius"),
			"CPU temperature",
			[]string{"cpu"}, nil,
		), prometheus.GaugeValue},
		logger: logger,
	}, nil
}

// Expose CPU stats using sysctl.
func (c *statCollector) Update(ch chan<- prometheus.Metric) error {
	// We want time spent per-cpu per CPUSTATE.
	// CPUSTATES (number of CPUSTATES) is defined as 5U.
	// Order: CP_USER | CP_NICE | CP_SYS | CP_IDLE | CP_INTR
	// sysctl kern.cp_time.x provides CPUSTATES long integers:
	//  (space-separated list of the above variables, where
	//   x stands for the number of the CPU core)
	//
	// Each value is a counter incremented at frequency
	//   kern.clockrate.(stathz | hz)
	//
	// Look into sys/kern/kern_clock.c for details.

	cpuTimes, err := getCPUTimes()
	if err != nil {
		return err
	}

	cpuTemperatures, err := getCPUTemperatures()
	if err != nil {
		return err
	}

	for cpu, t := range cpuTimes {
		lcpu := strconv.Itoa(cpu)
		ch <- c.cpu.mustNewConstMetric(float64(t.user), lcpu, "user")
		ch <- c.cpu.mustNewConstMetric(float64(t.nice), lcpu, "nice")
		ch <- c.cpu.mustNewConstMetric(float64(t.sys), lcpu, "system")
		ch <- c.cpu.mustNewConstMetric(float64(t.intr), lcpu, "interrupt")
		ch <- c.cpu.mustNewConstMetric(float64(t.idle), lcpu, "idle")

		if temp, ok := cpuTemperatures[cpu]; ok {
			ch <- c.temp.mustNewConstMetric(temp, lcpu)
		} else {
			c.logger.Debug("no temperature information for CPU", "cpu", cpu)
			ch <- c.temp.mustNewConstMetric(math.NaN(), lcpu)
		}
	}
	return err
}
