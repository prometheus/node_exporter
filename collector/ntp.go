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
	"net"
	"time"

	"github.com/beevik/ntp"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	hour24       = 24 * time.Hour // `time` does not export `Day` as Day != 24h because of DST
	ntpSubsystem = "ntp"
)

var (
	ntpServer          = kingpin.Flag("collector.ntp.server", "NTP server to use for ntp collector").Default("127.0.0.1").String()
	ntpProtocolVersion = kingpin.Flag("collector.ntp.protocol-version", "NTP protocol version").Default("4").Int()
	ntpServerIsLocal   = kingpin.Flag("collector.ntp.server-is-local", "Certify that collector.ntp.server address is the same local host as this collector.").Default("false").Bool()
	ntpIPTTL           = kingpin.Flag("collector.ntp.ip-ttl", "IP TTL to use while sending NTP query").Default("1").Int()
	// 3.46608s ~ 1.5s + PHI * (1 << maxPoll), where 1.5s is MAXDIST from ntp.org, it is 1.0 in RFC5905
	// max-distance option is used as-is without phi*(1<<poll)
	ntpMaxDistance     = kingpin.Flag("collector.ntp.max-distance", "Max accumulated distance to the root").Default("3.46608s").Duration()
	ntpOffsetTolerance = kingpin.Flag("collector.ntp.local-offset-tolerance", "Offset between local clock and local ntpd time to tolerate").Default("1ms").Duration()

	leapMidnight time.Time
)

type ntpCollector struct {
	stratum, leap, rtt, offset, reftime, rootDelay, rootDispersion, sanity typedDesc
}

func init() {
	registerCollector("ntp", defaultDisabled, NewNtpCollector)
}

// NewNtpCollector returns a new Collector exposing sanity of local NTP server.
// Default definition of "local" is:
// - collector.ntp.server address is a loopback address (or collector.ntp.server-is-mine flag is turned on)
// - the server is reachable with outgoin IP_TTL = 1
func NewNtpCollector() (Collector, error) {
	ipaddr := net.ParseIP(*ntpServer)
	if !*ntpServerIsLocal && (ipaddr == nil || !ipaddr.IsLoopback()) {
		return nil, fmt.Errorf("only IP address of local NTP server is valid for --collector.ntp.server")
	}

	if *ntpProtocolVersion < 2 || *ntpProtocolVersion > 4 {
		return nil, fmt.Errorf("invalid NTP protocol version %d; must be 2, 3, or 4", *ntpProtocolVersion)
	}

	if *ntpOffsetTolerance < 0 {
		return nil, fmt.Errorf("Offset tolerance must be non-negative")
	}

	return &ntpCollector{
		stratum: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, ntpSubsystem, "stratum"),
			"NTPD stratum.",
			nil, nil,
		), prometheus.GaugeValue},
		leap: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, ntpSubsystem, "leap"),
			"NTPD leap second indicator, 2 bits.",
			nil, nil,
		), prometheus.GaugeValue},
		rtt: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, ntpSubsystem, "rtt_seconds"),
			"RTT to NTPD.",
			nil, nil,
		), prometheus.GaugeValue},
		offset: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, ntpSubsystem, "offset_seconds"),
			"ClockOffset between NTP and local clock.",
			nil, nil,
		), prometheus.GaugeValue},
		reftime: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, ntpSubsystem, "reference_timestamp_seconds"),
			"NTPD ReferenceTime, UNIX timestamp.",
			nil, nil,
		), prometheus.GaugeValue},
		rootDelay: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, ntpSubsystem, "root_delay_seconds"),
			"NTPD RootDelay.",
			nil, nil,
		), prometheus.GaugeValue},
		rootDispersion: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, ntpSubsystem, "root_dispersion_seconds"),
			"NTPD RootDispersion.",
			nil, nil,
		), prometheus.GaugeValue},
		sanity: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, ntpSubsystem, "sanity"),
			"NTPD sanity according to RFC5905 heuristics and configured limits.",
			nil, nil,
		), prometheus.GaugeValue},
	}, nil
}

func (c *ntpCollector) Update(ch chan<- prometheus.Metric) error {
	resp, err := ntp.QueryWithOptions(*ntpServer, ntp.QueryOptions{
		Version: *ntpProtocolVersion,
		TTL:     *ntpIPTTL,
		Timeout: time.Second, // default `ntpdate` timeout
	})
	if err != nil {
		return fmt.Errorf("couldn't get SNTP reply: %s", err)
	}

	ch <- c.stratum.mustNewConstMetric(float64(resp.Stratum))
	ch <- c.leap.mustNewConstMetric(float64(resp.Leap))
	ch <- c.rtt.mustNewConstMetric(resp.RTT.Seconds())
	ch <- c.offset.mustNewConstMetric(resp.ClockOffset.Seconds())
	if resp.ReferenceTime.Unix() > 0 {
		// Go Zero is   0001-01-01 00:00:00 UTC
		// NTP Zero is  1900-01-01 00:00:00 UTC
		// UNIX Zero is 1970-01-01 00:00:00 UTC
		// so let's keep ALL ancient `reftime` values as zero
		ch <- c.reftime.mustNewConstMetric(float64(resp.ReferenceTime.UnixNano()) / 1e9)
	} else {
		ch <- c.reftime.mustNewConstMetric(0)
	}
	ch <- c.rootDelay.mustNewConstMetric(resp.RootDelay.Seconds())
	ch <- c.rootDispersion.mustNewConstMetric(resp.RootDispersion.Seconds())

	// Here is SNTP packet sanity check that is exposed to move burden of
	// configuration from node_exporter user to the developer.

	maxerr := *ntpOffsetTolerance
	if resp.Leap == ntp.LeapAddSecond || resp.Leap == ntp.LeapDelSecond {
		// state of leapMidnight is cached as leap flag is dropped right after midnight
		leapMidnight = resp.Time.Truncate(hour24).Add(hour24)
	}
	if leapMidnight.Add(-hour24).Before(resp.Time) && resp.Time.Before(leapMidnight.Add(hour24)) {
		// tolerate leap smearing
		maxerr += time.Second
	}

	if resp.Validate() == nil && resp.RootDistance <= *ntpMaxDistance && resp.MinError <= maxerr {
		ch <- c.sanity.mustNewConstMetric(1)
	} else {
		ch <- c.sanity.mustNewConstMetric(0)
	}

	return nil
}
