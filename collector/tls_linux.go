// Copyright The Prometheus Authors
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

//go:build !notls

package collector

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

const tlsSubsystem = "tls"

type tlsStatDesc struct {
	typedDesc
	value func(stat *procfs.TLSStat) float64
}

type tlsCollector struct {
	descs  []tlsStatDesc
	fs     procfs.FS
	logger *slog.Logger
}

func init() {
	registerCollector("tls", defaultDisabled, NewTLSCollector)
}

// NewTLSCollector returns a new Collector exposing kernel TLS statistics
// from /proc/net/tls_stat, disabled by default.
func NewTLSCollector(logger *slog.Logger) (Collector, error) {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}

	return &tlsCollector{
		descs: []tlsStatDesc{
			{
				typedDesc: newTLSTypedDesc("curr_tx_sw", prometheus.GaugeValue,
					"Number of TX sessions currently installed where host handles cryptography."),
				value: func(stat *procfs.TLSStat) float64 { return float64(stat.TLSCurrTxSw) },
			},
			{
				typedDesc: newTLSTypedDesc("curr_rx_sw", prometheus.GaugeValue,
					"Number of RX sessions currently installed where host handles cryptography."),
				value: func(stat *procfs.TLSStat) float64 { return float64(stat.TLSCurrRxSw) },
			},
			{
				typedDesc: newTLSTypedDesc("curr_tx_device", prometheus.GaugeValue,
					"Number of TX sessions currently installed where NIC handles cryptography."),
				value: func(stat *procfs.TLSStat) float64 { return float64(stat.TLSCurrTxDevice) },
			},
			{
				typedDesc: newTLSTypedDesc("curr_rx_device", prometheus.GaugeValue,
					"Number of RX sessions currently installed where NIC handles cryptography."),
				value: func(stat *procfs.TLSStat) float64 { return float64(stat.TLSCurrRxDevice) },
			},
			{
				typedDesc: newTLSTypedDesc("tx_sw_total", prometheus.CounterValue,
					"Number of TX sessions opened with host cryptography."),
				value: func(stat *procfs.TLSStat) float64 { return float64(stat.TLSTxSw) },
			},
			{
				typedDesc: newTLSTypedDesc("rx_sw_total", prometheus.CounterValue,
					"Number of RX sessions opened with host cryptography."),
				value: func(stat *procfs.TLSStat) float64 { return float64(stat.TLSRxSw) },
			},
			{
				typedDesc: newTLSTypedDesc("tx_device_total", prometheus.CounterValue,
					"Number of TX sessions opened with NIC cryptography."),
				value: func(stat *procfs.TLSStat) float64 { return float64(stat.TLSTxDevice) },
			},
			{
				typedDesc: newTLSTypedDesc("rx_device_total", prometheus.CounterValue,
					"Number of RX sessions opened with NIC cryptography."),
				value: func(stat *procfs.TLSStat) float64 { return float64(stat.TLSRxDevice) },
			},
			{
				typedDesc: newTLSTypedDesc("decrypt_error_total", prometheus.CounterValue,
					"Number of record decryption failures, e.g. due to incorrect authentication tags."),
				value: func(stat *procfs.TLSStat) float64 { return float64(stat.TLSDecryptError) },
			},
			{
				typedDesc: newTLSTypedDesc("rx_device_resync_total", prometheus.CounterValue,
					"Number of RX resyncs sent to NICs handling cryptography."),
				value: func(stat *procfs.TLSStat) float64 { return float64(stat.TLSRxDeviceResync) },
			},
			{
				typedDesc: newTLSTypedDesc("decrypt_retry_total", prometheus.CounterValue,
					"Number of retried record decryption attempts."),
				value: func(stat *procfs.TLSStat) float64 { return float64(stat.TLSDecryptRetry) },
			},
			{
				typedDesc: newTLSTypedDesc("rx_nopad_violation_total", prometheus.CounterValue,
					"Number of RX no-pad violations."),
				value: func(stat *procfs.TLSStat) float64 { return float64(stat.TLSRxNoPadViolation) },
			},
		},
		fs:     fs,
		logger: logger,
	}, nil
}

func newTLSTypedDesc(name string, valueType prometheus.ValueType, help string) typedDesc {
	return typedDesc{
		desc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, tlsSubsystem, name),
			help,
			nil, nil,
		),
		valueType: valueType,
	}
}

func (c *tlsCollector) Update(ch chan<- prometheus.Metric) error {
	stats, err := c.fs.NewTLSStat()
	if err != nil {
		// The file only exists while the kernel TLS module is loaded.
		if errors.Is(err, os.ErrNotExist) {
			c.logger.Debug("TLS statistics unavailable", "err", err)
			return ErrNoData
		}
		return fmt.Errorf("failed to get TLS statistics: %w", err)
	}

	for _, statDesc := range c.descs {
		ch <- statDesc.mustNewConstMetric(statDesc.value(&stats))
	}

	return nil
}
