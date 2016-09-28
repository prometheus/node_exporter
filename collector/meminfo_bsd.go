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

// +build freebsd darwin,amd64 dragonfly
// +build !nomeminfo

package collector

import (
	"errors"

	"github.com/prometheus/client_golang/prometheus"
)

/*
#include <stddef.h>
#include <sys/sysctl.h>

int _sysctl(const char* name) {
        int val;
        size_t size = sizeof(val);
        int res = sysctlbyname(name, &val, &size, NULL, 0);
        if (res == -1) {
                return -1;
        }
        if (size != sizeof(val)) {
                return -2;
        }
        return val;
}
*/
import "C"

const (
	memInfoSubsystem = "memory"
)

type meminfoCollector struct{}

func init() {
	Factories["meminfo"] = NewMeminfoCollector
}

// Takes a prometheus registry and returns a new Collector exposing
// Memory stats.
func NewMeminfoCollector() (Collector, error) {
	return &meminfoCollector{}, nil
}

func (c *meminfoCollector) Update(ch chan<- prometheus.Metric) (err error) {
	var pages map[string]C.int
	pages = make(map[string]C.int)

	size := C._sysctl(C.CString("vm.stats.vm.v_page_size"))
	if size == -1 {
		return errors.New("sysctl(vm.stats.vm.v_page_size) failed")
	}
	if size == -2 {
		return errors.New("sysctl(vm.stats.vm.v_page_size) failed, wrong buffer size")
	}

	pages["active"] = C._sysctl(C.CString("vm.stats.vm.v_active_count"))
	pages["inactive"] = C._sysctl(C.CString("vm.stats.vm.v_inactive_count"))
	pages["wire"] = C._sysctl(C.CString("vm.stats.vm.v_wire_count"))
	pages["cache"] = C._sysctl(C.CString("vm.stats.vm.v_cache_count"))
	pages["free"] = C._sysctl(C.CString("vm.stats.vm.v_free_count"))
	pages["swappgsin"] = C._sysctl(C.CString("vm.stats.vm.v_swappgsin"))
	pages["swappgsout"] = C._sysctl(C.CString("vm.stats.vm.v_swappgsout"))
	pages["total"] = C._sysctl(C.CString("vm.stats.vm.v_page_count"))

	for key := range pages {
		if pages[key] == -1 {
			return errors.New("sysctl() failed for " + key)
		}
		if pages[key] == -2 {
			return errors.New("sysctl() failed for " + key + ", wrong buffer size")
		}
	}

	for k, v := range pages {
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(Namespace, memInfoSubsystem, k),
				k+" from sysctl()",
				nil, nil,
			),
			// Convert metrics to kB (same as Linux meminfo).
			prometheus.UntypedValue, float64(v)*float64(size),
		)
	}
	return err
}
