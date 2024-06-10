// Copyright 2021 The Prometheus Authors
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

//go:build !nofibrechannel
// +build !nofibrechannel

package collector

import (
	"fmt"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/sysfs"
)

const maxUint64 = ^uint64(0)

type fibrechannelCollector struct {
	fs          sysfs.FS
	metricDescs map[string]*prometheus.Desc
	logger      log.Logger
	subsystem   string
}

func init() {
	registerCollector("fibrechannel", defaultEnabled, NewFibreChannelCollector)
}

// NewFibreChannelCollector returns a new Collector exposing FibreChannel stats.
func NewFibreChannelCollector(logger log.Logger) (Collector, error) {
	var i fibrechannelCollector
	var err error

	i.fs, err = sysfs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}
	i.logger = logger

	// Detailed description for all metrics.
	descriptions := map[string]string{
		"dumped_frames_total":            "Number of dumped frames",
		"loss_of_signal_total":           "Number of times signal has been lost",
		"loss_of_sync_total":             "Number of failures on either bit or transmission word boundaries",
		"rx_frames_total":                "Number of frames received",
		"error_frames_total":             "Number of errors in frames",
		"invalid_tx_words_total":         "Number of invalid words transmitted by host port",
		"seconds_since_last_reset_total": "Number of seconds since last host port reset",
		"tx_words_total":                 "Number of words transmitted by host port",
		"invalid_crc_total":              "Invalid Cyclic Redundancy Check count",
		"nos_total":                      "Number Not_Operational Primitive Sequence received by host port",
		"fcp_packet_aborts_total":        "Number of aborted packets",
		"rx_words_total":                 "Number of words received by host port",
		"tx_frames_total":                "Number of frames transmitted by host port",
		"link_failure_total":             "Number of times the host port link has failed",
	}

	i.metricDescs = make(map[string]*prometheus.Desc)
	i.subsystem = "fibrechannel"

	for metricName, description := range descriptions {
		i.metricDescs[metricName] = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, i.subsystem, metricName),
			description,
			[]string{"fc_host"},
			nil,
		)
	}

	return &i, nil
}

func (c *fibrechannelCollector) pushMetric(ch chan<- prometheus.Metric, name string, value uint64, host string, valueType prometheus.ValueType) {
	ch <- prometheus.MustNewConstMetric(c.metricDescs[name], valueType, float64(value), host)
}

func (c *fibrechannelCollector) pushCounter(ch chan<- prometheus.Metric, name string, value *uint64, host string) {
	// Don't push counters that aren't implemented (a counter equal to maxUint64 is unimplemented by the HBA firmware)
	if value != nil && *value != maxUint64 {
		c.pushMetric(ch, name, *value, host, prometheus.CounterValue)
	}
}

func (c *fibrechannelCollector) Update(ch chan<- prometheus.Metric) error {
	hosts, err := c.fs.FibreChannelClass()
	if err != nil {
		if os.IsNotExist(err) {
			level.Debug(c.logger).Log("msg", "fibrechannel statistics not found, skipping")
			return ErrNoData
		}
		return fmt.Errorf("error obtaining FibreChannel class info: %s", err)
	}

	for _, host := range hosts {
		infoDesc := prometheus.NewDesc(
			prometheus.BuildFQName(namespace, c.subsystem, "info"),
			"Non-numeric data from /sys/class/fc_host/<host>, value is always 1.",
			[]string{"fc_host", "speed", "port_state", "port_type", "port_id", "port_name", "fabric_name", "symbolic_name", "supported_classes", "supported_speeds", "dev_loss_tmo"},
			nil,
		)
		infoValue := 1.0

		// First push the Host values
		ch <- prometheus.MustNewConstMetric(infoDesc, prometheus.GaugeValue, infoValue,
			nilToEmpty(host.Name), nilToEmpty(host.Speed), nilToEmpty(host.PortState), nilToEmpty(host.PortType),
			nilToEmpty(host.PortID), nilToEmpty(host.PortName), nilToEmpty(host.FabricName), nilToEmpty(host.SymbolicName),
			nilToEmpty(host.SupportedClasses), nilToEmpty(host.SupportedSpeeds), nilToEmpty(host.DevLossTMO))

		// Then the counters
		c.pushCounter(ch, "dumped_frames_total", host.Counters.DumpedFrames, nilToEmpty(host.Name))
		c.pushCounter(ch, "error_frames_total", host.Counters.ErrorFrames, nilToEmpty(host.Name))
		c.pushCounter(ch, "invalid_crc_total", host.Counters.InvalidCRCCount, nilToEmpty(host.Name))
		c.pushCounter(ch, "rx_frames_total", host.Counters.RXFrames, nilToEmpty(host.Name))
		c.pushCounter(ch, "rx_words_total", host.Counters.RXWords, nilToEmpty(host.Name))
		c.pushCounter(ch, "tx_frames_total", host.Counters.TXFrames, nilToEmpty(host.Name))
		c.pushCounter(ch, "tx_words_total", host.Counters.TXWords, nilToEmpty(host.Name))
		c.pushCounter(ch, "seconds_since_last_reset_total", host.Counters.SecondsSinceLastReset, nilToEmpty(host.Name))
		c.pushCounter(ch, "invalid_tx_words_total", host.Counters.InvalidTXWordCount, nilToEmpty(host.Name))
		c.pushCounter(ch, "link_failure_total", host.Counters.LinkFailureCount, nilToEmpty(host.Name))
		c.pushCounter(ch, "loss_of_sync_total", host.Counters.LossOfSyncCount, nilToEmpty(host.Name))
		c.pushCounter(ch, "loss_of_signal_total", host.Counters.LossOfSignalCount, nilToEmpty(host.Name))
		c.pushCounter(ch, "nos_total", host.Counters.NosCount, nilToEmpty(host.Name))
		c.pushCounter(ch, "fcp_packet_aborts_total", host.Counters.FCPPacketAborts, nilToEmpty(host.Name))
	}

	return nil
}

func nilToEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
