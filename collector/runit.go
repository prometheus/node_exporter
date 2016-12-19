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

// +build !norunit

package collector

import (
	"flag"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/soundcloud/go-runit/runit"
)

var runitServiceDir = flag.String(
	"collector.runit.servicecdir",
	"/etc/service",
	"Path to runit service directory.")

type runitCollector struct {
	state, stateDesired, stateNormal, stateTimestamp *prometheus.GaugeVec
}

func init() {
	Factories["runit"] = NewRunitCollector
}

func NewRunitCollector() (Collector, error) {
	var (
		subsystem   = "service"
		constLabels = prometheus.Labels{"supervisor": "runit"}
		labelNames  = []string{"service"}
	)

	return &runitCollector{
		state: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   Namespace,
				Subsystem:   subsystem,
				Name:        "state",
				Help:        "State of runit service.",
				ConstLabels: constLabels,
			},
			labelNames,
		),
		stateDesired: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   Namespace,
				Subsystem:   subsystem,
				Name:        "desired_state",
				Help:        "Desired state of runit service.",
				ConstLabels: constLabels,
			},
			labelNames,
		),
		stateNormal: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   Namespace,
				Subsystem:   subsystem,
				Name:        "normal_state",
				Help:        "Normal state of runit service.",
				ConstLabels: constLabels,
			},
			labelNames,
		),
		stateTimestamp: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   Namespace,
				Subsystem:   subsystem,
				Name:        "state_last_change_timestamp_seconds",
				Help:        "Unix timestamp of the last runit service state change.",
				ConstLabels: constLabels,
			},
			labelNames,
		),
	}, nil
}

func (c *runitCollector) Update(ch chan<- prometheus.Metric) error {
	services, err := runit.GetServices(*runitServiceDir)
	if err != nil {
		return err
	}

	for _, service := range services {
		status, err := service.Status()
		if err != nil {
			log.Debugf("Couldn't get status for %s: %s, skipping...", service.Name, err)
			continue
		}

		log.Debugf("%s is %d on pid %d for %d seconds", service.Name, status.State, status.Pid, status.Duration)
		c.state.WithLabelValues(service.Name).Set(float64(status.State))
		c.stateDesired.WithLabelValues(service.Name).Set(float64(status.Want))
		c.stateTimestamp.WithLabelValues(service.Name).Set(float64(status.Timestamp.Unix()))
		if status.NormallyUp {
			c.stateNormal.WithLabelValues(service.Name).Set(1)
		} else {
			c.stateNormal.WithLabelValues(service.Name).Set(0)
		}
	}
	c.state.Collect(ch)
	c.stateDesired.Collect(ch)
	c.stateNormal.Collect(ch)
	c.stateTimestamp.Collect(ch)

	return nil
}
