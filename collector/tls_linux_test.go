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
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

type testTLSCollector struct {
	c Collector
}

func (c testTLSCollector) Collect(ch chan<- prometheus.Metric) {
	c.c.Update(ch)
}

func (c testTLSCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func NewTestTLSCollector(t *testing.T, logger *slog.Logger) prometheus.Collector {
	t.Helper()

	c, err := NewTLSCollector(logger)
	if err != nil {
		t.Fatal(err)
	}
	return testTLSCollector{c: c}
}

func TestTLS(t *testing.T) {
	if _, err := kingpin.CommandLine.Parse([]string{}); err != nil {
		t.Fatal(err)
	}
	*procPath = "fixtures/proc"

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	c := NewTestTLSCollector(t, logger)

	expected := `
# HELP node_tls_curr_rx_device Number of RX sessions currently installed where NIC handles cryptography.
# TYPE node_tls_curr_rx_device gauge
node_tls_curr_rx_device 4
# HELP node_tls_curr_rx_sw Number of RX sessions currently installed where host handles cryptography.
# TYPE node_tls_curr_rx_sw gauge
node_tls_curr_rx_sw 2
# HELP node_tls_curr_tx_device Number of TX sessions currently installed where NIC handles cryptography.
# TYPE node_tls_curr_tx_device gauge
node_tls_curr_tx_device 3
# HELP node_tls_curr_tx_sw Number of TX sessions currently installed where host handles cryptography.
# TYPE node_tls_curr_tx_sw gauge
node_tls_curr_tx_sw 1
# HELP node_tls_decrypt_error_total Number of record decryption failures, e.g. due to incorrect authentication tags.
# TYPE node_tls_decrypt_error_total counter
node_tls_decrypt_error_total 9
# HELP node_tls_decrypt_retry_total Number of retried record decryption attempts.
# TYPE node_tls_decrypt_retry_total counter
node_tls_decrypt_retry_total 11
# HELP node_tls_rx_device_resync_total Number of RX resyncs sent to NICs handling cryptography.
# TYPE node_tls_rx_device_resync_total counter
node_tls_rx_device_resync_total 10
# HELP node_tls_rx_device_total Number of RX sessions opened with NIC cryptography.
# TYPE node_tls_rx_device_total counter
node_tls_rx_device_total 8
# HELP node_tls_rx_nopad_violation_total Number of RX no-pad violations.
# TYPE node_tls_rx_nopad_violation_total counter
node_tls_rx_nopad_violation_total 12
# HELP node_tls_rx_sw_total Number of RX sessions opened with host cryptography.
# TYPE node_tls_rx_sw_total counter
node_tls_rx_sw_total 6
# HELP node_tls_tx_device_total Number of TX sessions opened with NIC cryptography.
# TYPE node_tls_tx_device_total counter
node_tls_tx_device_total 7
# HELP node_tls_tx_sw_total Number of TX sessions opened with host cryptography.
# TYPE node_tls_tx_sw_total counter
node_tls_tx_sw_total 5
`

	if err := testutil.CollectAndCompare(c, strings.NewReader(expected)); err != nil {
		t.Fatal(err)
	}
}
