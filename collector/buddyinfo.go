// Copyright 2017 The Prometheus Authors
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

// +build !nobuddyinfo
// +build !netbsd

package collector

import (
	"fmt"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/procfs"
)

const (
	buddyInfoSubsystem = "buddyinfo"
)

type buddyinfoCollector struct {
	desc *prometheus.Desc
}

func init() {
	registerCollector("buddyinfo", defaultDisabled, NewBuddyinfoCollector)
}

// NewBuddyinfoCollector returns a new Collector exposing buddyinfo stats.
func NewBuddyinfoCollector() (Collector, error) {
	desc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, buddyInfoSubsystem, "blocks"),
		"Count of free blocks according to size.",
		[]string{"node", "zone", "size"}, nil,
	)
	return &buddyinfoCollector{desc}, nil
}

// Update calls (*buddyinfoCollector).getBuddyInfo to get the platform specific
// buddyinfo metrics.
func (c *buddyinfoCollector) Update(ch chan<- prometheus.Metric) error {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return fmt.Errorf("failed to open procfs: %v", err)
	}

	buddyInfo, err := fs.NewBuddyInfo()
	if err != nil {
		return fmt.Errorf("couldn't get buddyinfo: %s", err)
	}

	log.Debugf("Set node_buddy: %#v", buddyInfo)
	for _, entry := range buddyInfo {
		for size, value := range entry.Sizes {
			ch <- prometheus.MustNewConstMetric(
				c.desc,
				prometheus.GaugeValue, value,
				entry.Node, entry.Zone, strconv.Itoa(size),
			)
		}
	}
	return nil
}
