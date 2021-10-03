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

//go:build openbsd && !amd64 && !nointerrupts
// +build openbsd,!amd64,!nointerrupts

package collector

import (
	"fmt"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

/*
#include <sys/param.h>
#include <sys/types.h>
#include <sys/sysctl.h>

#include <errno.h>
#include <stdio.h>

struct intr
{
	int		vector;
	char		device[128];
	u_int64_t	count;
};

int sysctl_nintr(void)
{
	int nintr, mib[4];
	size_t siz;

	mib[0] = CTL_KERN;
	mib[1] = KERN_INTRCNT;
	mib[2] = KERN_INTRCNT_NUM;
	siz = sizeof(nintr);
	if (sysctl(mib, 3, &nintr, &siz, NULL, 0) < 0) {
		return -1;
	}

	return nintr;
}

int
sysctl_intr(struct intr *intr, int idx)
{
	int mib[4];
	size_t siz;
	u_quad_t cnt;

	mib[0] = CTL_KERN;
	mib[1] = KERN_INTRCNT;
	mib[2] = KERN_INTRCNT_NAME;
	mib[3] = idx;
	siz = sizeof intr->device;
	if (sysctl(mib, 4, intr->device, &siz, NULL, 0) < 0) {
		return -1;
	}

	mib[0] = CTL_KERN;
	mib[1] = KERN_INTRCNT;
	mib[2] = KERN_INTRCNT_VECTOR;
	mib[3] = idx;
	siz = sizeof intr->vector;
	if (sysctl(mib, 4, &intr->vector, &siz, NULL, 0) < 0) {
		return -1;
	}

	mib[0] = CTL_KERN;
	mib[1] = KERN_INTRCNT;
	mib[2] = KERN_INTRCNT_CNT;
	mib[3] = idx;
	siz = sizeof(cnt);
	if (sysctl(mib, 4, &cnt, &siz, NULL, 0) < 0) {
		return -1;
	}

	intr->count = cnt;

	return 1;
}
*/
import "C"

var (
	interruptLabelNames = []string{"cpu", "type", "devices"}
)

func (c *interruptsCollector) Update(ch chan<- prometheus.Metric) error {
	interrupts, err := getInterrupts()
	if err != nil {
		return fmt.Errorf("couldn't get interrupts: %w", err)
	}
	for dev, interrupt := range interrupts {
		for cpuNo, value := range interrupt.values {
			ch <- c.desc.mustNewConstMetric(
				value,
				strconv.Itoa(cpuNo),
				strconv.Itoa(interrupt.vector),
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
	var (
		cintr      C.struct_intr
		interrupts = map[string]interrupt{}
	)

	nintr := C.sysctl_nintr()

	for i := C.int(0); i < nintr; i++ {
		_, err := C.sysctl_intr(&cintr, i)
		if err != nil {
			return nil, err
		}

		dev := C.GoString(&cintr.device[0])

		interrupts[dev] = interrupt{
			vector: int(cintr.vector),
			device: dev,
			// XXX: openbsd appears to only handle interrupts on cpu 0.
			values: []float64{float64(cintr.count)},
		}
	}

	return interrupts, nil
}
