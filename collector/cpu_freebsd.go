// Copyright 2015 The Prometheus Authors
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

// +build !nocpu

package collector

import (
	"errors"
	"os"
	"strconv"
	"unsafe"

	"github.com/prometheus/client_golang/prometheus"
)

/*
#cgo LDFLAGS: -lkvm
#include <fcntl.h>
#include <kvm.h>
#include <stdlib.h>
#include <sys/param.h>
#include <sys/pcpu.h>
#include <sys/resource.h>
#include <sys/sysctl.h>
#include <sys/time.h>

long _clockrate() {
	struct clockinfo clockrate;
	size_t size = sizeof(clockrate);
	int res = sysctlbyname("kern.clockrate", &clockrate, &size, NULL, 0);
	if (res == -1) {
		return -1;
	}
	if (size != sizeof(clockrate)) {
		return -2;
	}
	return clockrate.stathz > 0 ? clockrate.stathz : clockrate.hz;
}

*/
import "C"

type statCollector struct {
	cpu *prometheus.CounterVec
}

func init() {
	Factories["cpu"] = NewStatCollector
}

// Takes a prometheus registry and returns a new Collector exposing
// CPU stats.
func NewStatCollector() (Collector, error) {
	return &statCollector{
		cpu: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Name:      "cpu_seconds_total",
				Help:      "Seconds the CPU spent in each mode.",
			},
			[]string{"cpu", "mode"},
		),
	}, nil
}

// Expose CPU stats using KVM.
func (c *statCollector) Update(ch chan<- prometheus.Metric) (err error) {
	if os.Geteuid() != 0 && os.Getegid() != 2 {
		return errors.New("caller should be either root user or kmem group to access /dev/mem")
	}

	var errbuf *C.char
	kd := C.kvm_open(nil, nil, nil, C.O_RDONLY, errbuf)
	if errbuf != nil {
		return errors.New("failed to call kvm_open()")
	}
	defer C.kvm_close(kd)

	// The cp_time variable is an array of CPUSTATES long integers -- in
	// the same format as the kern.cp_time sysctl.  According to the
	// comments in sys/kern/kern_clock.c, the frequency of this timer will
	// be stathz (or hz, if stathz is zero).
	clockrate, err := getClockRate()
	if err != nil {
		return err
	}

	ncpus := C.kvm_getncpus(kd)
	for i := 0; i < int(ncpus); i++ {
		pcpu := C.kvm_getpcpu(kd, C.int(i))
		cp_time := ((*C.struct_pcpu)(unsafe.Pointer(pcpu))).pc_cp_time
		c.cpu.With(prometheus.Labels{"cpu": strconv.Itoa(i), "mode": "user"}).Set(float64(cp_time[C.CP_USER]) / clockrate)
		c.cpu.With(prometheus.Labels{"cpu": strconv.Itoa(i), "mode": "nice"}).Set(float64(cp_time[C.CP_NICE]) / clockrate)
		c.cpu.With(prometheus.Labels{"cpu": strconv.Itoa(i), "mode": "system"}).Set(float64(cp_time[C.CP_SYS]) / clockrate)
		c.cpu.With(prometheus.Labels{"cpu": strconv.Itoa(i), "mode": "interrupt"}).Set(float64(cp_time[C.CP_INTR]) / clockrate)
		c.cpu.With(prometheus.Labels{"cpu": strconv.Itoa(i), "mode": "idle"}).Set(float64(cp_time[C.CP_IDLE]) / clockrate)
	}
	c.cpu.Collect(ch)
	return err
}

func getClockRate() (float64, error) {
	clockrate := C._clockrate()
	if clockrate == -1 {
		return 0, errors.New("sysctl(kern.clockrate) failed")
	} else if clockrate == -2 {
		return 0, errors.New("sysctl(kern.clockrate) failed, wrong buffer size")
	} else if clockrate <= 0 {
		return 0, errors.New("sysctl(kern.clockrate) bad clocktime")
	}
	return float64(clockrate), nil
}
