// Copyright 2024 The Prometheus Authors
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

//go:build !noktls
// +build !noktls

package collector

import (
	"fmt"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

type ktlsCollector struct {
	fs     procfs.FS
	logger log.Logger
}

func init() {
	registerCollector("ktls", defaultDisabled, NewKTLSCollector)
}

// NewKTLSCollector returns a new Collector exposing kTLS stats.
func NewKTLSCollector(logger log.Logger) (Collector, error) {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}

	return &ktlsCollector{
		fs:     fs,
		logger: logger,
	}, nil
}

func (c *ktlsCollector) Update(ch chan<- prometheus.Metric) error {
	stat, err := c.fs.NewTLSStat()
	if err != nil {
		return err
	}

	ktlsCurrTxSwDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "ktls", "tls_curr_tx_sw"),
		"number of TX sessions currently installed where host handles cryptography",
		nil, nil,
	)
	ktlsCurrRxSwDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "ktls", "tls_curr_rx_sw"),
		"number of RX sessions currently installed where host handles cryptography",
		nil, nil,
	)
	ktlsCurrTxDeviceDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "ktls", "tls_curr_tx_device"),
		"number of TX sessions currently installed where NIC handles cryptography",
		nil, nil,
	)
	ktlsCurrRxDeviceDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "ktls", "tls_curr_rx_device"),
		"number of RX sessions currently installed where NIC handles cryptography",
		nil, nil,
	)
	ktlsTxDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "ktls", "tls_tx_sw_total"),
		"number of TX sessions opened with host cryptography",
		nil, nil,
	)
	ktlsRxDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "ktls", "tls_rx_sw_total"),
		"number of RX sessions opened with host cryptography",
		nil, nil,
	)
	ktlsTxDeviceDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "ktls", "tls_tx_device_total"),
		"number of TX sessions opened with NIC cryptograph",
		nil, nil,
	)
	ktlsRxDeviceDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "ktls", "tls_rx_device_total"),
		"number of RX sessions opened with NIC cryptograph",
		nil, nil,
	)
	ktlsDecryptErrorDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "ktls", "tls_decrypt_error_total"),
		"record decryption failed (e.g. due to incorrect authentication tag)",
		nil, nil,
	)
	ktlsRxDeviceResyncDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "ktls", "tls_rx_device_resync_total"),
		"number of RX resyncs sent to NICs handling cryptography",
		nil, nil,
	)
	ktlsDecryptRetryDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "ktls", "tls_decrypt_retry_total"),
		"number of RX records which had to be re-decrypted due to TLS_RX_EXPECT_NO_PAD mis-prediction",
		nil, nil,
	)
	ktlsRxNoPadViolationDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "ktls", "tls_no_pad_violation_total"),
		"number of data RX records which had to be re-decrypted due to TLS_RX_EXPECT_NO_PAD mis-prediction",
		nil, nil,
	)

	ch <- prometheus.MustNewConstMetric(ktlsCurrTxSwDesc, prometheus.GaugeValue, float64(stat.TLSCurrTxSw))
	ch <- prometheus.MustNewConstMetric(ktlsCurrRxSwDesc, prometheus.GaugeValue, float64(stat.TLSCurrTxSw))
	ch <- prometheus.MustNewConstMetric(ktlsCurrTxDeviceDesc, prometheus.GaugeValue, float64(stat.TLSCurrTxDevice))
	ch <- prometheus.MustNewConstMetric(ktlsCurrRxDeviceDesc, prometheus.GaugeValue, float64(stat.TLSCurrRxDevice))
	ch <- prometheus.MustNewConstMetric(ktlsTxDesc, prometheus.CounterValue, float64(stat.TLSTxSw))
	ch <- prometheus.MustNewConstMetric(ktlsRxDesc, prometheus.CounterValue, float64(stat.TLSRxSw))
	ch <- prometheus.MustNewConstMetric(ktlsTxDeviceDesc, prometheus.CounterValue, float64(stat.TLSTxDevice))
	ch <- prometheus.MustNewConstMetric(ktlsRxDeviceDesc, prometheus.CounterValue, float64(stat.TLSRxDevice))
	ch <- prometheus.MustNewConstMetric(ktlsDecryptErrorDesc, prometheus.CounterValue, float64(stat.TLSDecryptError))
	ch <- prometheus.MustNewConstMetric(ktlsRxDeviceResyncDesc, prometheus.CounterValue, float64(stat.TLSRxDeviceResync))
	ch <- prometheus.MustNewConstMetric(ktlsDecryptRetryDesc, prometheus.CounterValue, float64(stat.TLSDecryptRetry))
	ch <- prometheus.MustNewConstMetric(ktlsRxNoPadViolationDesc, prometheus.CounterValue, float64(stat.TLSRxNoPadViolation))

	return err
}
