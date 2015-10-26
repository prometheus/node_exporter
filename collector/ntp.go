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
	"time"

	"github.com/beevik/ntp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/log"
)

var (
	ntpServer          = flag.String("collector.ntp.server", "", "NTP server to use for ntp collector.")
	ntpProtocolVersion = flag.Int("collector.ntp.protocol-version", 4, "NTP protocol version")
)

type ntpCollector struct {
	drift prometheus.Gauge
}

func init() {
	ntp.Version = byte(*ntpProtocolVersion)
	Factories["ntp"] = NewNtpCollector
}

// Takes a prometheus registry and returns a new Collector exposing
// the offset between ntp and the current system time.
func NewNtpCollector() (Collector, error) {
	if *ntpServer == "" {
		return nil, fmt.Errorf("no NTP server specifies, see --ntpServer")
	}

	return &ntpCollector{
		drift: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "ntp_drift_seconds",
			Help:      "Time between system time and ntp time.",
		}),
	}, nil
}

func (c *ntpCollector) Update(ch chan<- prometheus.Metric) (err error) {
	t, err := ntp.Time(*ntpServer)
	if err != nil {
		return fmt.Errorf("couldn't get ntp drift: %s", err)
	}
	drift := t.Sub(time.Now())
	log.Debugf("Set ntp_drift_seconds: %f", drift.Seconds())
	c.drift.Set(drift.Seconds())
	c.drift.Collect(ch)
	return err
}
