// Copyright 2016 The Prometheus Authors
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

// +build !nodrbd

package collector

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

// Numerical metric provided by /proc/drbd.
type drbdNumericalMetric struct {
	desc       *prometheus.Desc
	valueType  prometheus.ValueType
	multiplier float64
}

func newDRBDNumericalMetric(name, desc string, valueType prometheus.ValueType, multiplier float64) drbdNumericalMetric {
	return drbdNumericalMetric{
		desc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "drbd", name),
			desc,
			[]string{"device"},
			nil,
		),
		valueType:  valueType,
		multiplier: multiplier,
	}
}

// String pair metric provided by /proc/drbd.
type drbdStringPairMetric struct {
	desc    *prometheus.Desc
	valueOK string
}

func (m *drbdStringPairMetric) isOkay(v string) float64 {
	if v == m.valueOK {
		return 1
	}

	return 0
}

func newDRBDStringPairMetric(name, desc, valueOK string) drbdStringPairMetric {
	return drbdStringPairMetric{
		desc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "drbd", name),
			desc,
			[]string{"device", "node"},
			nil,
		),
		valueOK: valueOK,
	}
}

type drbdCollector struct {
	numerical  map[string]drbdNumericalMetric
	stringPair map[string]drbdStringPairMetric
	connected  *prometheus.Desc
	logger     log.Logger
}

func init() {
	registerCollector("drbd", defaultDisabled, newDRBDCollector)
}

func newDRBDCollector(logger log.Logger) (Collector, error) {
	return &drbdCollector{
		numerical: map[string]drbdNumericalMetric{
			"ns": newDRBDNumericalMetric(
				"network_sent_bytes_total",
				"Total number of bytes sent via the network.",
				prometheus.CounterValue,
				1024,
			),
			"nr": newDRBDNumericalMetric(
				"network_received_bytes_total",
				"Total number of bytes received via the network.",
				prometheus.CounterValue,
				1,
			),
			"dw": newDRBDNumericalMetric(
				"disk_written_bytes_total",
				"Net data written on local hard disk; in bytes.",
				prometheus.CounterValue,
				1024,
			),
			"dr": newDRBDNumericalMetric(
				"disk_read_bytes_total",
				"Net data read from local hard disk; in bytes.",
				prometheus.CounterValue,
				1024,
			),
			"al": newDRBDNumericalMetric(
				"activitylog_writes_total",
				"Number of updates of the activity log area of the meta data.",
				prometheus.CounterValue,
				1,
			),
			"bm": newDRBDNumericalMetric(
				"bitmap_writes_total",
				"Number of updates of the bitmap area of the meta data.",
				prometheus.CounterValue,
				1,
			),
			"lo": newDRBDNumericalMetric(
				"local_pending",
				"Number of open requests to the local I/O sub-system.",
				prometheus.GaugeValue,
				1,
			),
			"pe": newDRBDNumericalMetric(
				"remote_pending",
				"Number of requests sent to the peer, but that have not yet been answered by the latter.",
				prometheus.GaugeValue,
				1,
			),
			"ua": newDRBDNumericalMetric(
				"remote_unacknowledged",
				"Number of requests received by the peer via the network connection, but that have not yet been answered.",
				prometheus.GaugeValue,
				1,
			),
			"ap": newDRBDNumericalMetric(
				"application_pending",
				"Number of block I/O requests forwarded to DRBD, but not yet answered by DRBD.",
				prometheus.GaugeValue,
				1,
			),
			"ep": newDRBDNumericalMetric(
				"epochs",
				"Number of Epochs currently on the fly.",
				prometheus.GaugeValue,
				1,
			),
			"oos": newDRBDNumericalMetric(
				"out_of_sync_bytes",
				"Amount of data known to be out of sync; in bytes.",
				prometheus.GaugeValue,
				1024,
			),
		},

		stringPair: map[string]drbdStringPairMetric{
			"ro": newDRBDStringPairMetric(
				"node_role_is_primary",
				"Whether the role of the node is in the primary state.",
				"Primary",
			),
			"ds": newDRBDStringPairMetric(
				"disk_state_is_up_to_date",
				"Whether the disk of the node is up to date.",
				"UpToDate",
			),
		},

		connected: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "drbd", "connected"),
			"Whether DRBD is connected to the peer.",
			[]string{"device"},
			nil,
		),
		logger: logger,
	}, nil
}

func (c *drbdCollector) Update(ch chan<- prometheus.Metric) error {
	statsFile := procFilePath("drbd")
	file, err := os.Open(statsFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			level.Debug(c.logger).Log("msg", "stats file does not exist, skipping", "file", statsFile, "err", err)
			return ErrNoData
		}

		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	device := "unknown"

	for scanner.Scan() {
		field := scanner.Text()

		kv := strings.Split(field, ":")
		if len(kv) != 2 {
			level.Debug(c.logger).Log("msg", "skipping invalid key:value pair", "field", field)
			continue
		}

		if id, err := strconv.ParseUint(kv[0], 10, 64); err == nil && kv[1] == "" {
			// New DRBD device encountered.
			device = fmt.Sprintf("drbd%d", id)
			continue
		}

		if m, ok := c.numerical[kv[0]]; ok {
			// Numerical value.
			v, err := strconv.ParseFloat(kv[1], 64)
			if err != nil {
				return err
			}

			ch <- prometheus.MustNewConstMetric(
				m.desc,
				m.valueType,
				v*m.multiplier,
				device,
			)

			continue
		}

		if m, ok := c.stringPair[kv[0]]; ok {
			// String pair value.
			values := strings.Split(kv[1], "/")
			ch <- prometheus.MustNewConstMetric(
				m.desc,
				prometheus.GaugeValue,
				m.isOkay(values[0]),
				device,
				"local",
			)

			ch <- prometheus.MustNewConstMetric(
				m.desc,
				prometheus.GaugeValue,
				m.isOkay(values[1]),
				device,
				"remote",
			)

			continue
		}

		if kv[0] == "cs" {
			// Connection state.
			var connected float64
			if kv[1] == "Connected" {
				connected = 1
			}

			ch <- prometheus.MustNewConstMetric(
				c.connected,
				prometheus.GaugeValue,
				connected,
				device,
			)

			continue
		}

		level.Debug(c.logger).Log("msg", "unhandled key-value pair", "key", kv[0], "value", kv[1])
	}

	return scanner.Err()
}
