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

//go:build !noethtool

package collector

import (
	"fmt"
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

type testEthtoolCollector struct {
	dsc Collector
}

func (c testEthtoolCollector) Collect(ch chan<- prometheus.Metric) {
	c.dsc.Update(ch)
}

func (c testEthtoolCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func NewTestEthtoolCollector(logger *slog.Logger) (prometheus.Collector, error) {
	dsc, err := NewEthtoolTestCollector(logger)
	if err != nil {
		return testEthtoolCollector{}, err
	}
	return testEthtoolCollector{
		dsc: dsc,
	}, err
}

func NewEthtoolTestCollector(logger *slog.Logger) (Collector, error) {
	collector, err := makeEthtoolCollector(logger)
	if err != nil {
		return nil, err
	}
	collector.ethtool = &EthtoolFixture{
		fixturePath: "fixtures/ethtool/",
	}
	return collector, nil
}

func TestBuildEthtoolFQName(t *testing.T) {
	testcases := map[string]string{
		"port.rx_errors":               "node_ethtool_port_received_errors",
		"rx_errors":                    "node_ethtool_received_errors",
		"Queue[0] AllocFails":          "node_ethtool_queue_0_allocfails",
		"Tx LPI entry count":           "node_ethtool_transmitted_lpi_entry_count",
		"port.VF_admin_queue_requests": "node_ethtool_port_vf_admin_queue_requests",
		"[3]: tx_bytes":                "node_ethtool_3_transmitted_bytes",
		"     err":                     "node_ethtool_err",
	}

	for metric, expected := range testcases {
		got := buildEthtoolFQName(metric)
		if expected != got {
			t.Errorf("Expected '%s' but got '%s'", expected, got)
		}
	}
}

func TestEthToolCollector(t *testing.T) {
	testcase := `# HELP node_ethtool_align_errors Network interface align_errors
# TYPE node_ethtool_align_errors untyped
node_ethtool_align_errors{device="eth0"} 0
# HELP node_ethtool_info A metric with a constant '1' value labeled by bus_info, device, driver, expansion_rom_version, firmware_version, version.
# TYPE node_ethtool_info gauge
node_ethtool_info{bus_info="0000:00:1f.6",device="eth0",driver="e1000e",expansion_rom_version="",firmware_version="0.5-4",version="5.11.0-22-generic"} 1
# HELP node_ethtool_port_received_dropped Network interface port_rx_dropped
# TYPE node_ethtool_port_received_dropped untyped
node_ethtool_port_received_dropped{device="eth0"} 12028
# HELP node_ethtool_received_broadcast Network interface rx_broadcast
# TYPE node_ethtool_received_broadcast untyped
node_ethtool_received_broadcast{device="eth0"} 5792
# HELP node_ethtool_received_errors_total Number of received frames with errors
# TYPE node_ethtool_received_errors_total untyped
node_ethtool_received_errors_total{device="eth0"} 0
# HELP node_ethtool_received_missed Network interface rx_missed
# TYPE node_ethtool_received_missed untyped
node_ethtool_received_missed{device="eth0"} 401
# HELP node_ethtool_received_multicast Network interface rx_multicast
# TYPE node_ethtool_received_multicast untyped
node_ethtool_received_multicast{device="eth0"} 23973
# HELP node_ethtool_received_packets_total Network interface packets received
# TYPE node_ethtool_received_packets_total untyped
node_ethtool_received_packets_total{device="eth0"} 1.260062e+06
# HELP node_ethtool_received_unicast Network interface rx_unicast
# TYPE node_ethtool_received_unicast untyped
node_ethtool_received_unicast{device="eth0"} 1.230297e+06
# HELP node_ethtool_transmitted_aborted Network interface tx_aborted
# TYPE node_ethtool_transmitted_aborted untyped
node_ethtool_transmitted_aborted{device="eth0"} 0
# HELP node_ethtool_transmitted_errors_total Number of sent frames with errors
# TYPE node_ethtool_transmitted_errors_total untyped
node_ethtool_transmitted_errors_total{device="eth0"} 0
# HELP node_ethtool_transmitted_multi_collisions Network interface tx_multi_collisions
# TYPE node_ethtool_transmitted_multi_collisions untyped
node_ethtool_transmitted_multi_collisions{device="eth0"} 0
# HELP node_ethtool_transmitted_packets_total Network interface packets sent
# TYPE node_ethtool_transmitted_packets_total untyped
node_ethtool_transmitted_packets_total{device="eth0"} 961500
# HELP node_ethtool_transmitted_single_collisions Network interface tx_single_collisions
# TYPE node_ethtool_transmitted_single_collisions untyped
node_ethtool_transmitted_single_collisions{device="eth0"} 0
# HELP node_ethtool_transmitted_underrun Network interface tx_underrun
# TYPE node_ethtool_transmitted_underrun untyped
node_ethtool_transmitted_underrun{device="eth0"} 0
# HELP node_network_advertised_speed_bytes Combination of speeds and features offered by network device
# TYPE node_network_advertised_speed_bytes gauge
node_network_advertised_speed_bytes{device="eth0",duplex="full",mode="1000baseT"} 1.25e+08
node_network_advertised_speed_bytes{device="eth0",duplex="full",mode="100baseT"} 1.25e+07
node_network_advertised_speed_bytes{device="eth0",duplex="full",mode="10baseT"} 1.25e+06
node_network_advertised_speed_bytes{device="eth0",duplex="half",mode="100baseT"} 1.25e+07
node_network_advertised_speed_bytes{device="eth0",duplex="half",mode="10baseT"} 1.25e+06
# HELP node_network_asymmetricpause_advertised If this port device offers asymmetric pause capability
# TYPE node_network_asymmetricpause_advertised gauge
node_network_asymmetricpause_advertised{device="eth0"} 0
# HELP node_network_asymmetricpause_supported If this port device supports asymmetric pause frames
# TYPE node_network_asymmetricpause_supported gauge
node_network_asymmetricpause_supported{device="eth0"} 0
# HELP node_network_autonegotiate If this port is using autonegotiate
# TYPE node_network_autonegotiate gauge
node_network_autonegotiate{device="eth0"} 1
# HELP node_network_autonegotiate_advertised If this port device offers autonegotiate
# TYPE node_network_autonegotiate_advertised gauge
node_network_autonegotiate_advertised{device="eth0"} 1
# HELP node_network_autonegotiate_supported If this port device supports autonegotiate
# TYPE node_network_autonegotiate_supported gauge
node_network_autonegotiate_supported{device="eth0"} 1
# HELP node_network_pause_advertised If this port device offers pause capability
# TYPE node_network_pause_advertised gauge
node_network_pause_advertised{device="eth0"} 1
# HELP node_network_pause_supported If this port device supports pause frames
# TYPE node_network_pause_supported gauge
node_network_pause_supported{device="eth0"} 1
# HELP node_network_supported_port_info Type of ports or PHYs supported by network device
# TYPE node_network_supported_port_info gauge
node_network_supported_port_info{device="eth0",type="MII"} 1
node_network_supported_port_info{device="eth0",type="TP"} 1
# HELP node_network_supported_speed_bytes Combination of speeds and features supported by network device
# TYPE node_network_supported_speed_bytes gauge
node_network_supported_speed_bytes{device="eth0",duplex="full",mode="10000baseT"} 1.25e+09
node_network_supported_speed_bytes{device="eth0",duplex="full",mode="1000baseT"} 1.25e+08
node_network_supported_speed_bytes{device="eth0",duplex="full",mode="100baseT"} 1.25e+07
node_network_supported_speed_bytes{device="eth0",duplex="full",mode="10baseT"} 1.25e+06
node_network_supported_speed_bytes{device="eth0",duplex="half",mode="100baseT"} 1.25e+07
node_network_supported_speed_bytes{device="eth0",duplex="half",mode="10baseT"} 1.25e+06
`
	*sysPath = "fixtures/sys"

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	collector, err := NewEthtoolTestCollector(logger)
	if err != nil {
		t.Fatal(err)
	}
	c, err := NewTestEthtoolCollector(logger)
	if err != nil {
		t.Fatal(err)
	}
	reg := prometheus.NewRegistry()
	reg.MustRegister(c)

	sink := make(chan prometheus.Metric)
	go func() {
		err = collector.Update(sink)
		if err != nil {
			panic(fmt.Errorf("failed to update collector: %s", err))
		}
		close(sink)
	}()

	err = testutil.GatherAndCompare(reg, strings.NewReader(testcase))
	if err != nil {
		t.Fatal(err)
	}
}

func TestEthToolCollectorWithIfAlias(t *testing.T) {
	testcase := `# HELP node_ethtool_align_errors Network interface align_errors
# TYPE node_ethtool_align_errors untyped
node_ethtool_align_errors{device="eth0",ifalias="alias"} 0
# HELP node_ethtool_info A metric with a constant '1' value labeled by bus_info, device, ifalias, driver, expansion_rom_version, firmware_version, version.
# TYPE node_ethtool_info gauge
node_ethtool_info{bus_info="0000:00:1f.6",device="eth0",driver="e1000e",expansion_rom_version="",firmware_version="0.5-4",ifalias="alias",version="5.11.0-22-generic"} 1
# HELP node_ethtool_port_received_dropped Network interface port_rx_dropped
# TYPE node_ethtool_port_received_dropped untyped
node_ethtool_port_received_dropped{device="eth0",ifalias="alias"} 12028
# HELP node_ethtool_received_broadcast Network interface rx_broadcast
# TYPE node_ethtool_received_broadcast untyped
node_ethtool_received_broadcast{device="eth0",ifalias="alias"} 5792
# HELP node_ethtool_received_errors_total Number of received frames with errors
# TYPE node_ethtool_received_errors_total untyped
node_ethtool_received_errors_total{device="eth0",ifalias="alias"} 0
# HELP node_ethtool_received_missed Network interface rx_missed
# TYPE node_ethtool_received_missed untyped
node_ethtool_received_missed{device="eth0",ifalias="alias"} 401
# HELP node_ethtool_received_multicast Network interface rx_multicast
# TYPE node_ethtool_received_multicast untyped
node_ethtool_received_multicast{device="eth0",ifalias="alias"} 23973
# HELP node_ethtool_received_packets_total Network interface packets received
# TYPE node_ethtool_received_packets_total untyped
node_ethtool_received_packets_total{device="eth0",ifalias="alias"} 1.260062e+06
# HELP node_ethtool_received_unicast Network interface rx_unicast
# TYPE node_ethtool_received_unicast untyped
node_ethtool_received_unicast{device="eth0",ifalias="alias"} 1.230297e+06
# HELP node_ethtool_transmitted_aborted Network interface tx_aborted
# TYPE node_ethtool_transmitted_aborted untyped
node_ethtool_transmitted_aborted{device="eth0",ifalias="alias"} 0
# HELP node_ethtool_transmitted_errors_total Number of sent frames with errors
# TYPE node_ethtool_transmitted_errors_total untyped
node_ethtool_transmitted_errors_total{device="eth0",ifalias="alias"} 0
# HELP node_ethtool_transmitted_multi_collisions Network interface tx_multi_collisions
# TYPE node_ethtool_transmitted_multi_collisions untyped
node_ethtool_transmitted_multi_collisions{device="eth0",ifalias="alias"} 0
# HELP node_ethtool_transmitted_packets_total Network interface packets sent
# TYPE node_ethtool_transmitted_packets_total untyped
node_ethtool_transmitted_packets_total{device="eth0",ifalias="alias"} 961500
# HELP node_ethtool_transmitted_single_collisions Network interface tx_single_collisions
# TYPE node_ethtool_transmitted_single_collisions untyped
node_ethtool_transmitted_single_collisions{device="eth0",ifalias="alias"} 0
# HELP node_ethtool_transmitted_underrun Network interface tx_underrun
# TYPE node_ethtool_transmitted_underrun untyped
node_ethtool_transmitted_underrun{device="eth0",ifalias="alias"} 0
# HELP node_network_advertised_speed_bytes Combination of speeds and features offered by network device
# TYPE node_network_advertised_speed_bytes gauge
node_network_advertised_speed_bytes{device="eth0",duplex="full",ifalias="alias",mode="1000baseT"} 1.25e+08
node_network_advertised_speed_bytes{device="eth0",duplex="full",ifalias="alias",mode="100baseT"} 1.25e+07
node_network_advertised_speed_bytes{device="eth0",duplex="full",ifalias="alias",mode="10baseT"} 1.25e+06
node_network_advertised_speed_bytes{device="eth0",duplex="half",ifalias="alias",mode="100baseT"} 1.25e+07
node_network_advertised_speed_bytes{device="eth0",duplex="half",ifalias="alias",mode="10baseT"} 1.25e+06
# HELP node_network_asymmetricpause_advertised If this port device offers asymmetric pause capability
# TYPE node_network_asymmetricpause_advertised gauge
node_network_asymmetricpause_advertised{device="eth0",ifalias="alias"} 0
# HELP node_network_asymmetricpause_supported If this port device supports asymmetric pause frames
# TYPE node_network_asymmetricpause_supported gauge
node_network_asymmetricpause_supported{device="eth0",ifalias="alias"} 0
# HELP node_network_autonegotiate If this port is using autonegotiate
# TYPE node_network_autonegotiate gauge
node_network_autonegotiate{device="eth0",ifalias="alias"} 1
# HELP node_network_autonegotiate_advertised If this port device offers autonegotiate
# TYPE node_network_autonegotiate_advertised gauge
node_network_autonegotiate_advertised{device="eth0",ifalias="alias"} 1
# HELP node_network_autonegotiate_supported If this port device supports autonegotiate
# TYPE node_network_autonegotiate_supported gauge
node_network_autonegotiate_supported{device="eth0",ifalias="alias"} 1
# HELP node_network_pause_advertised If this port device offers pause capability
# TYPE node_network_pause_advertised gauge
node_network_pause_advertised{device="eth0",ifalias="alias"} 1
# HELP node_network_pause_supported If this port device supports pause frames
# TYPE node_network_pause_supported gauge
node_network_pause_supported{device="eth0",ifalias="alias"} 1
# HELP node_network_supported_port_info Type of ports or PHYs supported by network device
# TYPE node_network_supported_port_info gauge
node_network_supported_port_info{device="eth0",ifalias="alias",type="MII"} 1
node_network_supported_port_info{device="eth0",ifalias="alias",type="TP"} 1
# HELP node_network_supported_speed_bytes Combination of speeds and features supported by network device
# TYPE node_network_supported_speed_bytes gauge
node_network_supported_speed_bytes{device="eth0",duplex="full",ifalias="alias",mode="10000baseT"} 1.25e+09
node_network_supported_speed_bytes{device="eth0",duplex="full",ifalias="alias",mode="1000baseT"} 1.25e+08
node_network_supported_speed_bytes{device="eth0",duplex="full",ifalias="alias",mode="100baseT"} 1.25e+07
node_network_supported_speed_bytes{device="eth0",duplex="full",ifalias="alias",mode="10baseT"} 1.25e+06
node_network_supported_speed_bytes{device="eth0",duplex="half",ifalias="alias",mode="100baseT"} 1.25e+07
node_network_supported_speed_bytes{device="eth0",duplex="half",ifalias="alias",mode="10baseT"} 1.25e+06
`
	*ethtoolAddIfAliasLabel = true
	t.Cleanup(func() { *ethtoolAddIfAliasLabel = false })
	*sysPath = "fixtures/sys"

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	collector, err := NewEthtoolTestCollector(logger)
	if err != nil {
		t.Fatal(err)
	}
	c, err := NewTestEthtoolCollector(logger)
	if err != nil {
		t.Fatal(err)
	}
	reg := prometheus.NewRegistry()
	reg.MustRegister(c)

	sink := make(chan prometheus.Metric)
	go func() {
		err = collector.Update(sink)
		if err != nil {
			panic(fmt.Errorf("failed to update collector: %s", err))
		}
		close(sink)
	}()

	err = testutil.GatherAndCompare(reg, strings.NewReader(testcase))
	if err != nil {
		t.Fatal(err)
	}
}
