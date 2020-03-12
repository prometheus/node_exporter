// Copyright 2015 The Prometheus Authors
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

// +build !noipvs

package collector

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

type exarDxCollector struct {
	Collector
	fs                                                                                                                      procfs.FS
	dxTotalFinCmds, dxFailedCmds                                                                                            typedDesc
	up, dxTemperature, dxCounterOverTemperature, dxCounterPCIeError, dxCounterILKError, dxCounterIntegrityError             typedDesc
	dxCounterDeviceHang, dxCounterLinkDown, dxCounterLinkWidthDegradation, dxCounterLinkSpeedDegradation, dxCounterHotReset typedDesc
	logger                                                                                                                  log.Logger
}

func init() {
	registerCollector("exar", defaultDisabled, NewExarDxCollector)
}

// NewExarDxCollector sets up a new collector for ExarDx metrics. It accepts the
// "procfs" config parameter to override the default proc location (/proc).
func NewExarDxCollector(logger log.Logger) (Collector, error) {
	return newExarDxCollector(logger)
}

func newExarDxCollector(logger log.Logger) (*exarDxCollector, error) {
	var (
		dxCmdStatisticsLabelNames = []string{
			"device",
			"ring",
			"type",
			"Device0",
			"Status",
			"Link_Speed",
			"Link_Width",
		}
		dxCmdStatusLabelNames = []string{
			"Device0",
			"Status",
			"Link_Speed",
			"Link_Width",
		}
		c         exarDxCollector
		err       error
		subsystem = "exar"
	)

	c.logger = logger
	c.fs, err = procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}

	c.dxTotalFinCmds = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "totalFinCmds"),
		"The total number of totalFinCmds.",
		dxCmdStatisticsLabelNames, nil,
	), prometheus.CounterValue}
	c.dxFailedCmds = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "failedCmds"),
		"The total number of failedCmds.",
		dxCmdStatisticsLabelNames, nil,
	), prometheus.CounterValue}

	c.up = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "up"),
		"up.",
		nil, nil,
	), prometheus.GaugeValue}
	c.dxTemperature = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "Temperature"),
		"Temperature",
		dxCmdStatusLabelNames, nil,
	), prometheus.GaugeValue}

	c.dxCounterOverTemperature = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "Counter_Over_Temperature"),
		"Counter_Over_Temperature.",
		dxCmdStatusLabelNames, nil,
	), prometheus.CounterValue}

	c.dxCounterPCIeError = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "Counter_PCIe_Error"),
		"Counter_PCIe_Error.",
		dxCmdStatusLabelNames, nil,
	), prometheus.CounterValue}

	c.dxCounterILKError = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "Counter_ILK_Error"),
		"Counter_ILK_Error.",
		dxCmdStatusLabelNames, nil,
	), prometheus.CounterValue}

	c.dxCounterIntegrityError = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "Counter_Integrity_Error"),
		"Counter_Integrity_Error.",
		dxCmdStatusLabelNames, nil,
	), prometheus.CounterValue}

	c.dxCounterDeviceHang = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "Counter_Device_Hang"),
		"Counter_Device_Hang.",
		dxCmdStatusLabelNames, nil,
	), prometheus.CounterValue}

	c.dxCounterLinkDown = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "Counter_Link_Down"),
		"Counter_Link_Down.",
		dxCmdStatusLabelNames, nil,
	), prometheus.CounterValue}

	c.dxCounterLinkWidthDegradation = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "Counter_Link_Width_Degradation"),
		"Counter_Link_Width_Degradation.",
		dxCmdStatusLabelNames, nil,
	), prometheus.CounterValue}

	c.dxCounterLinkSpeedDegradation = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "Counter_Link_Speed_Degradation"),
		"Counter_Link_Speed_Degradation.",
		dxCmdStatusLabelNames, nil,
	), prometheus.CounterValue}

	c.dxCounterHotReset = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "Counter_Hot_Reset"),
		"Counter_Hot_Reset.",
		dxCmdStatusLabelNames, nil,
	), prometheus.CounterValue}

	return &c, nil
}

func (c *exarDxCollector) Update(ch chan<- prometheus.Metric) error {

	exarStatus, err := c.fs.ExarDxStatus()
	if err != nil {
		// Cannot access ipvs metrics, report no error.
		if os.IsNotExist(err) {
			level.Debug(c.logger).Log("msg", "Exar collector metrics are not available for this system")
			ch <- c.up.mustNewConstMetric(float64(0))
			return ErrNoData
		}
		ch <- c.up.mustNewConstMetric(float64(0))
		return fmt.Errorf("could not get Exar stats: %s", err)
	}
	StatuslabelValues := []string{
		GetMapString(exarStatus.Status, "Device0"),
		GetMapString(exarStatus.Status, "Status"),
		GetMapString(exarStatus.Status, "Link_Speed"),
		GetMapString(exarStatus.Status, "Link_Width"),
	}
	ch <- c.dxTemperature.mustNewConstMetric(float64(GetMapUint64(exarStatus.Status, "Temperature")), StatuslabelValues...)
	ch <- c.dxCounterOverTemperature.mustNewConstMetric(float64(GetMapUint64(exarStatus.Status, "Counter_Over_Temperature")), StatuslabelValues...)
	ch <- c.dxCounterPCIeError.mustNewConstMetric(float64(GetMapUint64(exarStatus.Status, "Counter_PCIe_Error")), StatuslabelValues...)
	ch <- c.dxCounterILKError.mustNewConstMetric(float64(GetMapUint64(exarStatus.Status, "Counter_ILK_Error")), StatuslabelValues...)
	ch <- c.dxCounterIntegrityError.mustNewConstMetric(float64(GetMapUint64(exarStatus.Status, "Counter_Integrity_Error")), StatuslabelValues...)
	ch <- c.dxCounterDeviceHang.mustNewConstMetric(float64(GetMapUint64(exarStatus.Status, "Counter_Device_Hang")), StatuslabelValues...)
	ch <- c.dxCounterLinkDown.mustNewConstMetric(float64(GetMapUint64(exarStatus.Status, "Counter_Link_Down")), StatuslabelValues...)
	ch <- c.dxCounterLinkWidthDegradation.mustNewConstMetric(float64(GetMapUint64(exarStatus.Status, "Counter_Link_Width_Degradation")), StatuslabelValues...)
	ch <- c.dxCounterLinkSpeedDegradation.mustNewConstMetric(float64(GetMapUint64(exarStatus.Status, "Counter_Link_Speed_Degradation")), StatuslabelValues...)
	ch <- c.dxCounterHotReset.mustNewConstMetric(float64(GetMapUint64(exarStatus.Status, "Counter_Hot_Reset")), StatuslabelValues...)

	exar, err := c.fs.ExarDx()
	if err != nil {
		// Cannot access ipvs metrics, report no error.
		if os.IsNotExist(err) {
			level.Debug(c.logger).Log("msg", "Exar collector metrics are not available for this system")
			ch <- c.up.mustNewConstMetric(float64(0))
			return ErrNoData
		}
		ch <- c.up.mustNewConstMetric(float64(0))
		return fmt.Errorf("could not get Exar stats: %s", err)
	}
	for _, v := range exar.Rows {
		labelValues := []string{
			v.Device,
			v.Ring,
			v.Type,
		}
		labelValues = append(labelValues, StatuslabelValues...)
		ch <- c.dxTotalFinCmds.mustNewConstMetric(float64(v.TotalFinCmds), labelValues...)
		ch <- c.dxFailedCmds.mustNewConstMetric(float64(v.FailedCmds), labelValues...)
	}

	ch <- c.up.mustNewConstMetric(float64(1))

	return nil
}

func GetMapString(m map[string]string, value string) string {
	if v, ok := m[value]; ok {
		return v
	} else {
		return ""
	}
}

func GetMapUint64(m map[string]string, value string) uint64 {
	if v, ok := m[value]; ok {
		r, err := strconv.ParseUint(v, 16, 32)
		if err != nil {
			return uint64(0)
		}
		return r
	} else {
		return uint64(0)
	}
}
