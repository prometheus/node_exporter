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

package collector

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mdlayher/wifi"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type wifiCollector struct {
	InterfaceFrequencyHertz *prometheus.Desc

	StationConnectedSecondsTotal *prometheus.Desc
	StationInactiveSeconds       *prometheus.Desc
	StationReceiveBitsPerSecond  *prometheus.Desc
	StationTransmitBitsPerSecond *prometheus.Desc
	StationSignalDBM             *prometheus.Desc
	StationTransmitRetriesTotal  *prometheus.Desc
	StationTransmitFailedTotal   *prometheus.Desc
	StationBeaconLossTotal       *prometheus.Desc
}

var (
	collectorWifi = flag.String("collector.wifi", "", "test fixtures to use for wifi collector metrics")
)

func init() {
	Factories["wifi"] = NewWifiCollector
}

var _ wifiStater = &wifi.Client{}

// wifiStater is an interface used to swap out a *wifi.Client for end to end tests.
type wifiStater interface {
	Close() error
	Interfaces() ([]*wifi.Interface, error)
	StationInfo(ifi *wifi.Interface) (*wifi.StationInfo, error)
}

func NewWifiCollector() (Collector, error) {
	const (
		subsystem = "wifi"
	)

	var (
		labels = []string{"device"}
	)

	return &wifiCollector{
		InterfaceFrequencyHertz: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "interface_frequency_hertz"),
			"The current frequency a WiFi interface is operating at, in hertz.",
			labels,
			nil,
		),

		StationConnectedSecondsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "station_connected_seconds_total"),
			"The total number of seconds a station has been connected to an access point.",
			labels,
			nil,
		),

		StationInactiveSeconds: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "station_inactive_seconds"),
			"The number of seconds since any wireless activity has occurred on a station.",
			labels,
			nil,
		),

		StationReceiveBitsPerSecond: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "station_receive_bits_per_second"),
			"The current WiFi receive bitrate of a station, in bits per second.",
			labels,
			nil,
		),

		StationTransmitBitsPerSecond: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "station_transmit_bits_per_second"),
			"The current WiFi transmit bitrate of a station, in bits per second.",
			labels,
			nil,
		),

		StationSignalDBM: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "station_signal_dbm"),
			"The current WiFi signal strength, in decibel-milliwatts (dBm).",
			labels,
			nil,
		),

		StationTransmitRetriesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "station_transmit_retries_total"),
			"The total number of times a station has had to retry while sending a packet.",
			labels,
			nil,
		),

		StationTransmitFailedTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "station_transmit_failed_total"),
			"The total number of times a station has failed to send a packet.",
			labels,
			nil,
		),

		StationBeaconLossTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "station_beacon_loss_total"),
			"The total number of times a station has detected a beacon loss.",
			labels,
			nil,
		),
	}, nil
}

func (c *wifiCollector) Update(ch chan<- prometheus.Metric) error {
	stat, err := newWifiStater(*collectorWifi)
	if err != nil {
		// Cannot access wifi metrics, report no error
		if os.IsNotExist(err) {
			log.Debug("wifi collector metrics are not available for this system")
			return nil
		}

		return fmt.Errorf("failed to access wifi data: %v", err)
	}
	defer stat.Close()

	ifis, err := stat.Interfaces()
	if err != nil {
		return fmt.Errorf("failed to retrieve wifi interfaces: %v", err)
	}

	for _, ifi := range ifis {
		// Only collect metrics on stations for now
		if ifi.Type != wifi.InterfaceTypeStation {
			continue
		}

		ch <- prometheus.MustNewConstMetric(
			c.InterfaceFrequencyHertz,
			prometheus.GaugeValue,
			mHzToHz(ifi.Frequency),
			ifi.Name,
		)

		info, err := stat.StationInfo(ifi)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}

			return fmt.Errorf("failed to retrieve station info for device %s: %v",
				ifi.Name, err)
		}

		c.updateStationStats(ch, ifi.Name, info)
	}

	return nil
}

func (c *wifiCollector) updateStationStats(ch chan<- prometheus.Metric, device string, info *wifi.StationInfo) {
	ch <- prometheus.MustNewConstMetric(
		c.StationConnectedSecondsTotal,
		prometheus.CounterValue,
		info.Connected.Seconds(),
		device,
	)

	ch <- prometheus.MustNewConstMetric(
		c.StationInactiveSeconds,
		prometheus.GaugeValue,
		info.Inactive.Seconds(),
		device,
	)

	ch <- prometheus.MustNewConstMetric(
		c.StationReceiveBitsPerSecond,
		prometheus.GaugeValue,
		float64(info.ReceiveBitrate),
		device,
	)

	ch <- prometheus.MustNewConstMetric(
		c.StationTransmitBitsPerSecond,
		prometheus.GaugeValue,
		float64(info.TransmitBitrate),
		device,
	)

	ch <- prometheus.MustNewConstMetric(
		c.StationSignalDBM,
		prometheus.GaugeValue,
		float64(info.Signal),
		device,
	)

	ch <- prometheus.MustNewConstMetric(
		c.StationTransmitRetriesTotal,
		prometheus.CounterValue,
		float64(info.TransmitRetries),
		device,
	)

	ch <- prometheus.MustNewConstMetric(
		c.StationTransmitFailedTotal,
		prometheus.CounterValue,
		float64(info.TransmitFailed),
		device,
	)

	ch <- prometheus.MustNewConstMetric(
		c.StationBeaconLossTotal,
		prometheus.CounterValue,
		float64(info.BeaconLoss),
		device,
	)
}

func mHzToHz(mHz int) float64 {
	return float64(mHz) * 1000 * 1000
}

// All code below this point is used to assist with end-to-end tests for
// the wifi collector, since wifi devices are not available in CI.

// newWifiStater determines if mocked test fixtures from files should be used for
// collecting wifi metrics, or if package wifi should be used.
func newWifiStater(fixtures string) (wifiStater, error) {
	if fixtures != "" {
		return &mockWifiStater{
			fixtures: fixtures,
		}, nil
	}

	return wifi.New()
}

var _ wifiStater = &mockWifiStater{}

type mockWifiStater struct {
	fixtures string
}

func (s *mockWifiStater) unmarshalJSONFile(filename string, v interface{}) error {
	b, err := ioutil.ReadFile(filepath.Join(s.fixtures, filename))
	if err != nil {
		return err
	}

	return json.Unmarshal(b, v)
}

func (s *mockWifiStater) Close() error { return nil }

func (s *mockWifiStater) Interfaces() ([]*wifi.Interface, error) {
	var ifis []*wifi.Interface
	if err := s.unmarshalJSONFile("interfaces.json", &ifis); err != nil {
		return nil, err
	}

	return ifis, nil
}

func (s *mockWifiStater) StationInfo(ifi *wifi.Interface) (*wifi.StationInfo, error) {
	p := filepath.Join(ifi.Name, "stationinfo.json")

	var info wifi.StationInfo
	if err := s.unmarshalJSONFile(p, &info); err != nil {
		return nil, err
	}

	return &info, nil
}
