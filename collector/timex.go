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

// +build linux
// +build !notimex

package collector

// #include <sys/timex.h>
import "C"

import (
	"fmt"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	// The system clock is not synchronized to a reliable server.
	timeError = C.TIME_ERROR
	// The timex.Status time resolution bit, 0 = microsecond, 1 = nanoseconds.
	staNano = C.STA_NANO
	// 1 second in
	nanoSeconds  = 1000000000
	microSeconds = 1000000
)

type timexCollector struct {
	offset,
	freq,
	maxerror,
	esterror,
	status,
	constant,
	tick,
	ppsfreq,
	jitter,
	shift,
	stabil,
	jitcnt,
	calcnt,
	errcnt,
	stbcnt,
	tai,
	syncStatus typedDesc
}

func init() {
	registerCollector("timex", defaultEnabled, NewTimexCollector)
}

// NewTimexCollector returns a new Collector exposing adjtime(3) stats.
func NewTimexCollector() (Collector, error) {
	const subsystem = "timex"

	return &timexCollector{
		offset: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "offset_seconds"),
			"Time offset in between local system and reference clock.",
			nil, nil,
		), prometheus.GaugeValue},
		freq: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "frequency_adjustment"),
			"Local clock frequency adjustment.",
			nil, nil,
		), prometheus.GaugeValue},
		maxerror: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "maxerror_seconds"),
			"Maximum error in seconds.",
			nil, nil,
		), prometheus.GaugeValue},
		esterror: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "estimated_error_seconds"),
			"Estimated error in seconds.",
			nil, nil,
		), prometheus.GaugeValue},
		status: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "status"),
			"Value of the status array bits.",
			nil, nil,
		), prometheus.GaugeValue},
		constant: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "loop_time_constant"),
			"Phase-locked loop time constant.",
			nil, nil,
		), prometheus.GaugeValue},
		tick: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "tick_seconds"),
			"Seconds between clock ticks.",
			nil, nil,
		), prometheus.GaugeValue},
		ppsfreq: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "pps_frequency"),
			"Pulse per second frequency.",
			nil, nil,
		), prometheus.GaugeValue},
		jitter: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "pps_jitter_seconds"),
			"Pulse per second jitter.",
			nil, nil,
		), prometheus.GaugeValue},
		shift: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "pps_shift_seconds"),
			"Pulse per second interval duration.",
			nil, nil,
		), prometheus.GaugeValue},
		stabil: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "pps_stability"),
			"Pulse per second stability.",
			nil, nil,
		), prometheus.CounterValue},
		jitcnt: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "pps_jitter_count"),
			"Pulse per second count of jitter limit exceeded events.",
			nil, nil,
		), prometheus.CounterValue},
		calcnt: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "pps_calibration_count"),
			"Pulse per second count of calibration intervals.",
			nil, nil,
		), prometheus.CounterValue},
		errcnt: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "pps_error_count"),
			"Pulse per second count of calibration errors.",
			nil, nil,
		), prometheus.CounterValue},
		stbcnt: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "pps_stability_exceeded_count"),
			"Pulse per second count of stability limit exceeded events.",
			nil, nil,
		), prometheus.GaugeValue},
		tai: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "tai_offset"),
			"International Atomic Time (TAI) offset.",
			nil, nil,
		), prometheus.GaugeValue},
		syncStatus: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "sync_status"),
			"Is clock synchronized to a reliable server (1 = yes, 0 = no).",
			nil, nil,
		), prometheus.GaugeValue},
	}, nil
}

func (c *timexCollector) Update(ch chan<- prometheus.Metric) error {
	var syncStatus float64
	var divisor float64
	var timex = new(syscall.Timex)

	status, err := syscall.Adjtimex(timex)
	if err != nil {
		return fmt.Errorf("failed to retrieve adjtimex stats: %v", err)
	}

	if status == timeError {
		syncStatus = 0
	} else {
		syncStatus = 1
	}
	if (timex.Status & staNano) != 0 {
		divisor = nanoSeconds
	} else {
		divisor = microSeconds
	}
	ch <- c.syncStatus.mustNewConstMetric(syncStatus)
	ch <- c.offset.mustNewConstMetric(float64(timex.Offset) / divisor)
	ch <- c.freq.mustNewConstMetric(float64(timex.Freq))
	ch <- c.maxerror.mustNewConstMetric(float64(timex.Maxerror) / microSeconds)
	ch <- c.esterror.mustNewConstMetric(float64(timex.Esterror) / microSeconds)
	ch <- c.status.mustNewConstMetric(float64(timex.Status))
	ch <- c.constant.mustNewConstMetric(float64(timex.Constant))
	ch <- c.tick.mustNewConstMetric(float64(timex.Tick) / microSeconds)
	ch <- c.ppsfreq.mustNewConstMetric(float64(timex.Ppsfreq))
	ch <- c.jitter.mustNewConstMetric(float64(timex.Jitter) / divisor)
	ch <- c.shift.mustNewConstMetric(float64(timex.Shift))
	ch <- c.stabil.mustNewConstMetric(float64(timex.Stabil))
	ch <- c.jitcnt.mustNewConstMetric(float64(timex.Jitcnt))
	ch <- c.calcnt.mustNewConstMetric(float64(timex.Calcnt))
	ch <- c.errcnt.mustNewConstMetric(float64(timex.Errcnt))
	ch <- c.stbcnt.mustNewConstMetric(float64(timex.Stbcnt))
	ch <- c.tai.mustNewConstMetric(float64(timex.Tai))

	return nil
}
