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

// +build !nontp

package collector

import (
	"flag"
	"fmt"

	"github.com/beevik/ntp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

var (
	ntpServer          = flag.String("collector.ntp.server", "", "NTP server to use for ntp collector.")
	ntpProtocolVersion = flag.Int("collector.ntp.protocol-version", 4, "NTP protocol version")
)

type ntpCollector struct {
	drift   prometheus.Gauge
	stratum prometheus.Gauge
}

func init() {
	Factories["ntp"] = NewNtpCollector
}

// Takes a prometheus registry and returns a new Collector exposing
// the offset between ntp and the current system time.
func NewNtpCollector() (Collector, error) {
	if *ntpServer == "" {
		return nil, fmt.Errorf("no NTP server specified, see -collector.ntp.server")
	}
	if *ntpProtocolVersion < 2 || *ntpProtocolVersion > 4 {
		return nil, fmt.Errorf("invalid NTP protocol version %d; must be 2, 3, or 4", *ntpProtocolVersion)
	}

	return &ntpCollector{
		drift: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "ntp_drift_seconds",
			Help:      "Time between system time and ntp time.",
		}),
		stratum: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "ntp_stratum",
			Help:      "NTP server stratum.",
		}),
	}, nil
}

func (c *ntpCollector) Update(ch chan<- prometheus.Metric) (err error) {
	resp, err := ntp.Query(*ntpServer, *ntpProtocolVersion)
	if err != nil {
		return fmt.Errorf("couldn't get NTP drift: %s", err)
	}
	driftSeconds := resp.ClockOffset.Seconds()
	log.Debugf("Set ntp_drift_seconds: %f", driftSeconds)
	c.drift.Set(driftSeconds)
	c.drift.Collect(ch)

	stratum := float64(resp.Stratum)
	log.Debugf("Set ntp_stratum: %f", stratum)
	c.stratum.Set(stratum)
	c.stratum.Collect(ch)
	return nil
}
