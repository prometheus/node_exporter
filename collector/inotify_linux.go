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
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"golang.org/x/sys/unix"
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
	success := 0.0
	err := c.tryAddWatch()
	if err == nil {
		success = 1
	} else {
		log.Debugf("inotify: not successful: %s", err)
	}
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, inotifySubsystem, "watch_success"),
			"inotify_add_watch() working as desired",
			nil, nil,
		),
		prometheus.GaugeValue, success,
	)
	return nil
}

// tryAddWatch attempts to register an inotify watcher for the root filesystem.
func (c *inotifyCollector) tryAddWatch() error {
	fd, err := unix.InotifyInit()
	if fd < 0 {
		// If this fails, this usually means fs.inotify.max_user_instances is exhausted.
		return fmt.Errorf("inotify_init() did not return a valid fd: %s", err)
	}
	defer unix.Close(fd)

	_, err = unix.InotifyAddWatch(fd, "/", unix.IN_CREATE)
	if err != nil {
		// If this fails, this usually means fs.inotify.max_user_watches is exhausted.
		return fmt.Errorf("inotify_add_watch() failed: %s", err)
	}
	return nil
}
