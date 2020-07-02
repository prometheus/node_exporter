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

// +build openbsd
// +build !norcctl

package collector

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var (
	specialSvcs = map[string]struct{}{
		"accounting":   struct{}{},
		"check_quotas": struct{}{},
		"ipsec":        struct{}{},
		"library_aslr": struct{}{},
		"multicast":    struct{}{},
		"pf":           struct{}{},
		"spamd_black":  struct{}{},
	}
)

const (
	rcctl        = "/usr/sbin/rcctl"
	statusOK     = "(ok)"
	statusFailed = "(failed)"

	serviceUpdateInterval = time.Second * 30
)

type rcctlCollector struct {
	sync.RWMutex
	state, stateDesired typedDesc
	logger              log.Logger

	services []string
}

func init() {
	registerCollector("rcctl", defaultDisabled, NewRcctlCollector)
}

// NewRunitCollector returns a new Collector exposing runit statistics.
func NewRcctlCollector(logger log.Logger) (Collector, error) {
	var (
		subsystem   = "service"
		constLabels = prometheus.Labels{"supervisor": "rcctl"}
		labelNames  = []string{"service"}
	)
	c := &rcctlCollector{
		state: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "state"),
			"State of rcctl service.",
			labelNames, constLabels,
		), prometheus.GaugeValue},
		stateDesired: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "desired_state"),
			"Desired state of rcctl service.",
			labelNames, constLabels,
		), prometheus.GaugeValue},
		logger: logger,
	}
	go c.getServices()
	return c, nil
}

func (c *rcctlCollector) getServices() {
	for {
		var svcs []string
		cmd := exec.Command(rcctl, "ls", "on")
		on, err := cmd.Output()
		if err != nil {
			level.Debug(c.logger).Log("msg", "rcctl list services", "error", err)
			time.Sleep(serviceUpdateInterval)
			continue
		}
		for _, s := range strings.Split(string(on), "\n") {
			if s == "" {
				continue
			}
			if _, ok := specialSvcs[s]; ok {
				continue
			}
			svcs = append(svcs, s)
		}
		c.Lock()
		c.services = svcs
		c.Unlock()
		time.Sleep(serviceUpdateInterval)
	}
}

func (c *rcctlCollector) Update(ch chan<- prometheus.Metric) error {
	c.RLock()
	args := append([]string{"check"}, c.services...)
	c.RUnlock()
	cmd := exec.Command(rcctl, args...)
	states, err := cmd.Output()
	if err != nil {
		level.Debug(c.logger).Log("msg", "rcctl check services", "error", err)
	}
	for _, s := range strings.Split(string(states), "\n") {
		if s == "" {
			continue
		}
		status := 0.0
		if strings.Contains(s, statusOK) {
			status = 1.0
			s = strings.TrimSuffix(s, statusOK)
		} else {
			s = strings.TrimSuffix(s, statusFailed)
		}
		level.Debug(c.logger).Log("msg", "Current status", "Service", s, "Status", status)
		ch <- c.state.mustNewConstMetric(status, s)
		ch <- c.stateDesired.mustNewConstMetric(1.0, s)
	}
	return nil
}
