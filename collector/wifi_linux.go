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

//go:build !nowifi
// +build !nowifi

package collector

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/alecthomas/kingpin/v2"
	"github.com/mdlayher/wifi"
	"github.com/prometheus/client_golang/prometheus"
)

type wifiCollector struct {
	interfaceFrequencyHertz *prometheus.Desc
	stationInfo             *prometheus.Desc

	stationConnectedSecondsTotal *prometheus.Desc
	stationInactiveSeconds       *prometheus.Desc
	stationReceiveBitsPerSecond  *prometheus.Desc
	stationTransmitBitsPerSecond *prometheus.Desc
	stationReceiveBytesTotal     *prometheus.Desc
	stationTransmitBytesTotal    *prometheus.Desc
	stationSignalDBM             *prometheus.Desc
	stationTransmitRetriesTotal  *prometheus.Desc
	stationTransmitFailedTotal   *prometheus.Desc
	stationBeaconLossTotal       *prometheus.Desc

	logger *slog.Logger
}

var (
	collectorWifi = kingpin.Flag("collector.wifi.fixtures", "test fixtures to use for wifi collector metrics").Default("").String()
)

func init() {
	registerCollector("wifi", defaultDisabled, NewWifiCollector)
}

var _ wifiStater = &wifi.Client{}

// wifiStater is an interface used to swap out a *wifi.Client for end to end tests.
type wifiStater interface {
	BSS(ifi *wifi.Interface) (*wifi.BSS, error)
	Close() error
	Interfaces() ([]*wifi.Interface, error)
	StationInfo(ifi *wifi.Interface) ([]*wifi.StationInfo, error)
}

// NewWifiCollector returns a new Collector exposing Wifi statistics.
func NewWifiCollector(logger *slog.Logger) (Collector, error) {
	const (
		subsystem = "wifi"
	)

	var (
		labels = []string{"device", "mac_address"}
	)

	return &wifiCollector{
		interfaceFrequencyHertz: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "interface_frequency_hertz"),
			"The current frequency a WiFi interface is operating at, in hertz.",
			[]string{"device"},
			nil,
		),

		stationInfo: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "station_info"),
			"Labeled WiFi interface station information as provided by the operating system.",
			[]string{"device", "bssid", "ssid", "mode"},
			nil,
		),

		stationConnectedSecondsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "station_connected_seconds_total"),
			"The total number of seconds a station has been connected to an access point.",
			labels,
			nil,
		),

		stationInactiveSeconds: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "station_inactive_seconds"),
			"The number of seconds since any wireless activity has occurred on a station.",
			labels,
			nil,
		),

		stationReceiveBitsPerSecond: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "station_receive_bits_per_second"),
			"The current WiFi receive bitrate of a station, in bits per second.",
			labels,
			nil,
		),

		stationTransmitBitsPerSecond: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "station_transmit_bits_per_second"),
			"The current WiFi transmit bitrate of a station, in bits per second.",
			labels,
			nil,
		),

		stationReceiveBytesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "station_receive_bytes_total"),
			"The total number of bytes received by a WiFi station.",
			labels,
			nil,
		),

		stationTransmitBytesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "station_transmit_bytes_total"),
			"The total number of bytes transmitted by a WiFi station.",
			labels,
			nil,
		),

		stationSignalDBM: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "station_signal_dbm"),
			"The current WiFi signal strength, in decibel-milliwatts (dBm).",
			labels,
			nil,
		),

		stationTransmitRetriesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "station_transmit_retries_total"),
			"The total number of times a station has had to retry while sending a packet.",
			labels,
			nil,
		),

		stationTransmitFailedTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "station_transmit_failed_total"),
			"The total number of times a station has failed to send a packet.",
			labels,
			nil,
		),

		stationBeaconLossTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "station_beacon_loss_total"),
			"The total number of times a station has detected a beacon loss.",
			labels,
			nil,
		),
		logger: logger,
	}, nil
}

func (c *wifiCollector) Update(ch chan<- prometheus.Metric) error {
	stat, err := newWifiStater(*collectorWifi)
	if err != nil {
		// Cannot access wifi metrics, report no error.
		if errors.Is(err, os.ErrNotExist) {
			c.logger.Debug("wifi collector metrics are not available for this system")
			return ErrNoData
		}
		if errors.Is(err, os.ErrPermission) {
			c.logger.Debug("wifi collector got permission denied when accessing metrics")
			return ErrNoData
		}

		return fmt.Errorf("failed to access wifi data: %w", err)
	}
	defer stat.Close()

	ifis, err := stat.Interfaces()
	if err != nil {
		return fmt.Errorf("failed to retrieve wifi interfaces: %w", err)
	}

	for _, ifi := range ifis {
		// Some virtual devices have no "name" and should be skipped.
		if ifi.Name == "" {
			continue
		}

		c.logger.Debug("probing wifi device with type", "wifi", ifi.Name, "type", ifi.Type)

		ch <- prometheus.MustNewConstMetric(
			c.interfaceFrequencyHertz,
			prometheus.GaugeValue,
			mHzToHz(ifi.Frequency),
			ifi.Name,
		)

		// When a statistic is not available for a given interface, package wifi
		// returns a os.ErrNotExist error.  We leverage this to only export
		// metrics which are actually valid for given interface types.

		bss, err := stat.BSS(ifi)
		switch {
		case err == nil:
			c.updateBSSStats(ch, ifi.Name, bss)
		case errors.Is(err, os.ErrNotExist):
			c.logger.Debug("BSS information not found for wifi device", "name", ifi.Name)
		default:
			return fmt.Errorf("failed to retrieve BSS for device %s: %v",
				ifi.Name, err)
		}

		stations, err := stat.StationInfo(ifi)
		switch {
		case err == nil:
			for _, station := range stations {
				c.updateStationStats(ch, ifi.Name, station)
			}
		case errors.Is(err, os.ErrNotExist):
			c.logger.Debug("station information not found for wifi device", "name", ifi.Name)
		default:
			return fmt.Errorf("failed to retrieve station info for device %q: %v",
				ifi.Name, err)
		}
	}

	return nil
}

func (c *wifiCollector) updateBSSStats(ch chan<- prometheus.Metric, device string, bss *wifi.BSS) {
	// Synthetic metric which provides wifi station info, such as SSID, BSSID, etc.
	ch <- prometheus.MustNewConstMetric(
		c.stationInfo,
		prometheus.GaugeValue,
		1,
		device,
		bss.BSSID.String(),
		bss.SSID,
		bssStatusMode(bss.Status),
	)
}

func (c *wifiCollector) updateStationStats(ch chan<- prometheus.Metric, device string, info *wifi.StationInfo) {
	ch <- prometheus.MustNewConstMetric(
		c.stationConnectedSecondsTotal,
		prometheus.CounterValue,
		info.Connected.Seconds(),
		device,
		info.HardwareAddr.String(),
	)

	ch <- prometheus.MustNewConstMetric(
		c.stationInactiveSeconds,
		prometheus.GaugeValue,
		info.Inactive.Seconds(),
		device,
		info.HardwareAddr.String(),
	)

	ch <- prometheus.MustNewConstMetric(
		c.stationReceiveBitsPerSecond,
		prometheus.GaugeValue,
		float64(info.ReceiveBitrate),
		device,
		info.HardwareAddr.String(),
	)

	ch <- prometheus.MustNewConstMetric(
		c.stationTransmitBitsPerSecond,
		prometheus.GaugeValue,
		float64(info.TransmitBitrate),
		device,
		info.HardwareAddr.String(),
	)

	ch <- prometheus.MustNewConstMetric(
		c.stationReceiveBytesTotal,
		prometheus.CounterValue,
		float64(info.ReceivedBytes),
		device,
		info.HardwareAddr.String(),
	)

	ch <- prometheus.MustNewConstMetric(
		c.stationTransmitBytesTotal,
		prometheus.CounterValue,
		float64(info.TransmittedBytes),
		device,
		info.HardwareAddr.String(),
	)

	ch <- prometheus.MustNewConstMetric(
		c.stationSignalDBM,
		prometheus.GaugeValue,
		float64(info.Signal),
		device,
		info.HardwareAddr.String(),
	)

	ch <- prometheus.MustNewConstMetric(
		c.stationTransmitRetriesTotal,
		prometheus.CounterValue,
		float64(info.TransmitRetries),
		device,
		info.HardwareAddr.String(),
	)

	ch <- prometheus.MustNewConstMetric(
		c.stationTransmitFailedTotal,
		prometheus.CounterValue,
		float64(info.TransmitFailed),
		device,
		info.HardwareAddr.String(),
	)

	ch <- prometheus.MustNewConstMetric(
		c.stationBeaconLossTotal,
		prometheus.CounterValue,
		float64(info.BeaconLoss),
		device,
		info.HardwareAddr.String(),
	)
}

func mHzToHz(mHz int) float64 {
	return float64(mHz) * 1000 * 1000
}

func bssStatusMode(status wifi.BSSStatus) string {
	switch status {
	case wifi.BSSStatusAuthenticated, wifi.BSSStatusAssociated:
		return "client"
	case wifi.BSSStatusIBSSJoined:
		return "ad-hoc"
	default:
		return "unknown"
	}
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
	b, err := os.ReadFile(filepath.Join(s.fixtures, filename))
	if err != nil {
		return err
	}

	return json.Unmarshal(b, v)
}

func (s *mockWifiStater) Close() error { return nil }

func (s *mockWifiStater) BSS(ifi *wifi.Interface) (*wifi.BSS, error) {
	p := filepath.Join(ifi.Name, "bss.json")

	var bss wifi.BSS
	if err := s.unmarshalJSONFile(p, &bss); err != nil {
		return nil, err
	}

	return &bss, nil
}

func (s *mockWifiStater) Interfaces() ([]*wifi.Interface, error) {
	var ifis []*wifi.Interface
	if err := s.unmarshalJSONFile("interfaces.json", &ifis); err != nil {
		return nil, err
	}

	return ifis, nil
}

func (s *mockWifiStater) StationInfo(ifi *wifi.Interface) ([]*wifi.StationInfo, error) {
	p := filepath.Join(ifi.Name, "stationinfo.json")

	var stations []*wifi.StationInfo
	if err := s.unmarshalJSONFile(p, &stations); err != nil {
		return nil, err
	}

	return stations, nil
}
