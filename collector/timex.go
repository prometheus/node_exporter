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

//go:build linux && !notimex
// +build linux,!notimex

package collector

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sys/unix"
)

const (
	// The system clock is not synchronized to a reliable
	// server (TIME_ERROR).
	timeError = 5
	// The timex.Status time resolution bit (STA_NANO),
	// 0 = microsecond, 1 = nanoseconds.
	staNano = 0x2000

	// 1 second in
	nanoSeconds  = 1000000000
	microSeconds = 1000000

	// See NOTES in adjtimex(2).
	ppm16frac = 1000000.0 * 65536.0
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
	logger *slog.Logger
}

func init() {
	registerCollector("timex", defaultEnabled, NewTimexCollector)
}

// NewTimexCollector returns a new Collector exposing adjtime(3) stats.
func NewTimexCollector(logger *slog.Logger) (Collector, error) {
	const subsystem = "timex"

	return &timexCollector{
		offset: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "offset_seconds"),
			"Time offset in between local system and reference clock.",
			nil, nil,
		), prometheus.GaugeValue},
		freq: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "frequency_adjustment_ratio"),
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
			prometheus.BuildFQName(namespace, subsystem, "pps_frequency_hertz"),
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
			prometheus.BuildFQName(namespace, subsystem, "pps_stability_hertz"),
			"Pulse per second stability, average of recent frequency changes.",
			nil, nil,
		), prometheus.GaugeValue},
		jitcnt: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "pps_jitter_total"),
			"Pulse per second count of jitter limit exceeded events.",
			nil, nil,
		), prometheus.CounterValue},
		calcnt: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "pps_calibration_total"),
			"Pulse per second count of calibration intervals.",
			nil, nil,
		), prometheus.CounterValue},
		errcnt: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "pps_error_total"),
			"Pulse per second count of calibration errors.",
			nil, nil,
		), prometheus.CounterValue},
		stbcnt: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "pps_stability_exceeded_total"),
			"Pulse per second count of stability limit exceeded events.",
			nil, nil,
		), prometheus.CounterValue},
		tai: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "tai_offset_seconds"),
			"International Atomic Time (TAI) offset.",
			nil, nil,
		), prometheus.GaugeValue},
		syncStatus: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "sync_status"),
			"Is clock synchronized to a reliable server (1 = yes, 0 = no).",
			nil, nil,
		), prometheus.GaugeValue},
		logger: logger,
	}, nil
}

func (c *timexCollector) Update(ch chan<- prometheus.Metric) error {
	var syncStatus float64
	var divisor float64
	var timex = new(unix.Timex)

	status, err := unix.Adjtimex(timex)
	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			c.logger.Debug("Not collecting timex metrics", "err", err)
			return ErrNoData
		}
		return fmt.Errorf("failed to retrieve adjtimex stats: %w", err)
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
	ch <- c.freq.mustNewConstMetric(1 + float64(timex.Freq)/ppm16frac)
	ch <- c.maxerror.mustNewConstMetric(float64(timex.Maxerror) / microSeconds)
	ch <- c.esterror.mustNewConstMetric(float64(timex.Esterror) / microSeconds)
	ch <- c.status.mustNewConstMetric(float64(timex.Status))
	ch <- c.constant.mustNewConstMetric(float64(timex.Constant))
	ch <- c.tick.mustNewConstMetric(float64(timex.Tick) / microSeconds)
	ch <- c.ppsfreq.mustNewConstMetric(float64(timex.Ppsfreq) / ppm16frac)
	ch <- c.jitter.mustNewConstMetric(float64(timex.Jitter) / divisor)
	ch <- c.shift.mustNewConstMetric(float64(timex.Shift))
	ch <- c.stabil.mustNewConstMetric(float64(timex.Stabil) / ppm16frac)
	ch <- c.jitcnt.mustNewConstMetric(float64(timex.Jitcnt))
	ch <- c.calcnt.mustNewConstMetric(float64(timex.Calcnt))
	ch <- c.errcnt.mustNewConstMetric(float64(timex.Errcnt))
	ch <- c.stbcnt.mustNewConstMetric(float64(timex.Stbcnt))
	ch <- c.tai.mustNewConstMetric(float64(timex.Tai))

	return nil
}
