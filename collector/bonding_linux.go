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

// +build !nobonding

package collector

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type bondingCollector struct {
	slaves, active typedDesc
}

func init() {
	registerCollector("bonding", defaultEnabled, NewBondingCollector)
}

// NewBondingCollector returns a newly allocated bondingCollector.
// It exposes the number of configured and active slave of linux bonding interfaces.
func NewBondingCollector() (Collector, error) {
	return &bondingCollector{
		slaves: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "bonding", "slaves"),
			"Number of configured slaves per bonding interface.",
			[]string{"master"}, nil,
		), prometheus.GaugeValue},
		active: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "bonding", "active"),
			"Number of active slaves per bonding interface.",
			[]string{"master"}, nil,
		), prometheus.GaugeValue},
	}, nil
}

// Update reads and exposes bonding states, implements Collector interface. Caution: This works only on linux.
func (c *bondingCollector) Update(ch chan<- prometheus.Metric) error {
	statusfile := sysFilePath("class/net")
	bondingStats, err := readBondingStats(statusfile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Debugf("Not collecting bonding, file does not exist: %s", statusfile)
			return nil
		}
		return err
	}
	for master, status := range bondingStats {
		ch <- c.slaves.mustNewConstMetric(float64(status[0]), master)
		ch <- c.active.mustNewConstMetric(float64(status[1]), master)
	}
	return nil
}

func readBondingStats(root string) (status map[string][2]int, err error) {
	status = map[string][2]int{}
	masters, err := ioutil.ReadFile(path.Join(root, "bonding_masters"))
	if err != nil {
		return nil, err
	}
	for _, master := range strings.Fields(string(masters)) {
		slaves, err := ioutil.ReadFile(path.Join(root, master, "bonding", "slaves"))
		if err != nil {
			return nil, err
		}
		sstat := [2]int{0, 0}
		for _, slave := range strings.Fields(string(slaves)) {
			state, err := ioutil.ReadFile(path.Join(root, master, fmt.Sprintf("lower_%s", slave), "operstate"))
			if os.IsNotExist(err) {
				// some older? kernels use slave_ prefix
				state, err = ioutil.ReadFile(path.Join(root, master, fmt.Sprintf("slave_%s", slave), "operstate"))
			}
			if err != nil {
				return nil, err
			}
			sstat[0]++
			if strings.TrimSpace(string(state)) == "up" {
				sstat[1]++
			}
		}
		status[master] = sstat
	}
	return status, err
}
