// Copyright 2020 The Prometheus Authors
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

//go:build !nointerrupts
// +build !nointerrupts

package collector

import (
	"fmt"
	"strconv"
	"unsafe"

	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sys/unix"
)

const (
	KERN_INTRCNT        = 63
	KERN_INTRCNT_NUM    = 1
	KERN_INTRCNT_CNT    = 2
	KERN_INTRCNT_NAME   = 3
	KERN_INTRCNT_VECTOR = 4
)

func nintr() _C_int {
	mib := [3]_C_int{unix.CTL_KERN, KERN_INTRCNT, KERN_INTRCNT_NUM}
	buf, err := sysctl(mib[:])
	if err != nil {
		return 0
	}
	return *(*_C_int)(unsafe.Pointer(&buf[0]))
}

func intr(idx _C_int) (itr interrupt, err error) {
	mib := [4]_C_int{unix.CTL_KERN, KERN_INTRCNT, KERN_INTRCNT_NAME, idx}
	buf, err := sysctl(mib[:])
	if err != nil {
		return
	}
	dev := *(*[128]byte)(unsafe.Pointer(&buf[0]))
	itr.device = string(dev[:])

	mib[2] = KERN_INTRCNT_VECTOR
	buf, err = sysctl(mib[:])
	if err != nil {
		return
	}
	itr.vector = *(*int)(unsafe.Pointer(&buf[0]))

	mib[2] = KERN_INTRCNT_CNT
	buf, err = sysctl(mib[:])
	if err != nil {
		return
	}
	count := *(*uint64)(unsafe.Pointer(&buf[0]))
	itr.values = []float64{float64(count)}
	return
}

var interruptLabelNames = []string{"cpu", "type", "devices"}

func (c *interruptsCollector) Update(ch chan<- prometheus.Metric) error {
	interrupts, err := getInterrupts()
	if err != nil {
		return fmt.Errorf("couldn't get interrupts: %s", err)
	}
	for dev, interrupt := range interrupts {
		for cpuNo, value := range interrupt.values {
			ch <- c.desc.mustNewConstMetric(
				value,
				strconv.Itoa(cpuNo),
				fmt.Sprintf("%d", interrupt.vector),
				dev,
			)
		}
	}
	return nil
}

type interrupt struct {
	vector int
	device string
	values []float64
}

func getInterrupts() (map[string]interrupt, error) {
	var interrupts = map[string]interrupt{}
	n := nintr()

	for i := _C_int(0); i < n; i++ {
		itr, err := intr(i)
		if err != nil {
			return nil, err
		}
		interrupts[itr.device] = itr
	}

	return interrupts, nil
}
