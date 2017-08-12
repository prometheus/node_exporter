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
	"fmt"

	"github.com/beevik/ntp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	ntpServer          = kingpin.Flag("collector.ntp.server", "NTP server to use for ntp collector.").Default("").String()
	ntpProtocolVersion = kingpin.Flag("collector.ntp.protocol-version", "NTP protocol version").Default("4").Int()
)

type ntpCollector struct {
	drift, stratum typedDesc
}

func init() {
	Factories["ntp"] = NewNtpCollector
}

// NewNtpCollector returns a new Collector exposing the offset between ntp and
// the current system time.
func NewNtpCollector() (Collector, error) {
	warnDeprecated("ntp")
	if *ntpServer == "" {
		return nil, fmt.Errorf("no NTP server specified, see -collector.ntp.server")
	}
	if *ntpProtocolVersion < 2 || *ntpProtocolVersion > 4 {
		return nil, fmt.Errorf("invalid NTP protocol version %d; must be 2, 3, or 4", *ntpProtocolVersion)
	}

	return &ntpCollector{
		drift: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "ntp", "drift_seconds"),
			"Time between system time and ntp time.",
			nil, nil,
		), prometheus.GaugeValue},
		stratum: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "ntp", "stratum"),
			"NTP server stratum.",
			nil, nil,
		), prometheus.GaugeValue},
	}, nil
}

func (c *ntpCollector) Update(ch chan<- prometheus.Metric) error {
	resp, err := ntp.Query(*ntpServer, *ntpProtocolVersion)
	if err != nil {
		return fmt.Errorf("couldn't get NTP drift: %s", err)
	}
	driftSeconds := resp.ClockOffset.Seconds()
	log.Debugf("Set ntp_drift_seconds: %f", driftSeconds)
	ch <- c.drift.mustNewConstMetric(driftSeconds)

	stratum := float64(resp.Stratum)
	log.Debugf("Set ntp_stratum: %f", stratum)
	ch <- c.stratum.mustNewConstMetric(stratum)
	return nil
}
