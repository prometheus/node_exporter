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

// +build !nomeminfo

package collector

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

/*
#include <sys/param.h>
#include <sys/types.h>
#include <sys/sysctl.h>

int
sysctl_uvmexp(struct uvmexp *uvmexp)
{
	static int uvmexp_mib[] = {CTL_VM, VM_UVMEXP};
	size_t sz = sizeof(struct uvmexp);

	if(sysctl(uvmexp_mib, 2, uvmexp, &sz, NULL, 0) < 0)
		return -1;

	return 0;
}

*/
import "C"

const (
	memInfoSubsystem = "memory"
)

type meminfoCollector struct {
	metrics map[string]prometheus.Gauge
}

func init() {
	Factories["meminfo"] = NewMeminfoCollector
}

// Takes a prometheus registry and returns a new Collector exposing
// Memory stats.
func NewMeminfoCollector() (Collector, error) {
	return &meminfoCollector{
		metrics: map[string]prometheus.Gauge{},
	}, nil
}

func (c *meminfoCollector) Update(ch chan<- prometheus.Metric) (err error) {
	var pages map[string]C.int
	var uvmexp C.struct_uvmexp

	if _, err := C.sysctl_uvmexp(&uvmexp); err != nil {
		return fmt.Errorf("sysctl CTL_VM VM_UVMEXP failed: %v", err)
	}

	size := uvmexp.pagesize

	pages = make(map[string]C.int)

	pages["active_bytes"] = uvmexp.active
	pages["inactive_bytes"] = uvmexp.inactive
	pages["wire_bytes"] = uvmexp.wired
	pages["cache_bytes"] = uvmexp.vnodepages
	pages["free_bytes"] = uvmexp.free
	pages["swappgsin_bytes"] = uvmexp.pgswapin
	pages["swappgsout_bytes"] = uvmexp.pgswapout
	pages["total_bytes"] = uvmexp.npages

	log.Debugf("Set node_mem: %#v", pages)
	for k, v := range pages {
		if _, ok := c.metrics[k]; !ok {
			c.metrics[k] = prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: Namespace,
				Subsystem: memInfoSubsystem,
				Name:      k,
				Help:      k + " from sysctl()",
			})
		}
		// Convert metrics to kB
		c.metrics[k].Set(float64(v) * float64(size))
		c.metrics[k].Collect(ch)
	}
	return err
}
