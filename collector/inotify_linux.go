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

// +build !noinotify

package collector

import (
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	inotifySubsystem = "inotify"
)

type inotifyCollector struct {
	initSuccess     float64
	addWatchSuccess float64
}

func init() {
	registerCollector(inotifySubsystem, defaultEnabled, NewInotifyCollector)
}

// NewInotifyCollector returns a new Collector for inotify test results.
func NewInotifyCollector() (Collector, error) {
	return &inotifyCollector{}, nil
}

func (c *inotifyCollector) Update(ch chan<- prometheus.Metric) error {
	c.tryAddWatch()
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, inotifySubsystem, "init_success"),
			"inotify_init() working as desired",
			nil, nil,
		),
		prometheus.GaugeValue, c.initSuccess,
	)
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, inotifySubsystem, "add_watch_success"),
			"inotify_add_watch() working as desired",
			nil, nil,
		),
		prometheus.GaugeValue, c.addWatchSuccess,
	)
	return nil
}

// tryAddWatch attempts to register an inotify watcher for the root filesystem.
// The result is written to the inotifyCollector fields initSuccess and addWatchSuccess.
func (c *inotifyCollector) tryAddWatch() {
	// Start by considering inotify broken unless the opposite is proven.
	c.initSuccess = 0
	c.addWatchSuccess = 0
	fd, _ := syscall.InotifyInit()
	if fd < 0 {
		// If this fails, this usually means fs.inotify.max_user_instances is exhausted.
		return
	}
	defer syscall.Close(fd)
	c.initSuccess = 1

	_, err := syscall.InotifyAddWatch(fd, "/", syscall.IN_CREATE)
	if err != nil {
		// If this fails, this usually means fs.inotify.max_user_watches is exhausted.
		return
	}
	c.addWatchSuccess = 1
}
