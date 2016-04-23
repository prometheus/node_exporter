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

// +build !nogops

package collector

import (
	"regexp"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

type gopsCollector struct {
	ignoredNetDevices *regexp.Regexp
	cpu               *prometheus.Desc
	mem               *prometheus.Desc
	swap              *prometheus.Desc
	swapIn            *prometheus.Desc
	swapOut           *prometheus.Desc
	bootTime          *prometheus.Desc
	netBytesSent      *prometheus.Desc
	netBytesRecv      *prometheus.Desc
	netPacketsSent    *prometheus.Desc
	netPacketsRecv    *prometheus.Desc
	netErrorsSent     *prometheus.Desc
	netErrorsRecv     *prometheus.Desc
}

func init() {
	Factories["gops"] = NewGopsCollector
}

func NewGopsCollector() (Collector, error) {
	ignoredNetDevices, err := regexp.Compile(*netdevIgnoredDevices)
	if err != nil {
		return nil, err
	}
	// FIXME: Add PlatformInformation but collected once at start?
	return &gopsCollector{
		ignoredNetDevices: ignoredNetDevices,
		cpu: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "cpu"),
			"Seconds the cpus spent in each mode.",
			[]string{"cpu", "mode"}, nil,
		),
		mem: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "memory"),
			"Memory in bytes by state.",
			[]string{"state"}, nil,
		),
		swap: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "memory", "swap"),
			"Swap in bytes by state.",
			[]string{"state"}, nil,
		),
		swapIn: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "memory", "swap_in"),
			"Total number of bytes swapped in from disk.",
			nil, nil,
		),
		swapOut: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "memory", "swap_out"),
			"Total number of bytes swapped out from disk.",
			nil, nil,
		),
		bootTime: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "boot", "time"),
			"Node boot time, in unixtime.",
			nil, nil,
		),
		netBytesSent: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "network", "transmit_bytes_total"),
			"Total number of bytes transmitted from device.",
			[]string{"device"},
			nil,
		),
		netBytesRecv: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "network", "receive_bytes_total"),
			"Total number of bytes received by device.",
			[]string{"device"},
			nil,
		),
		netPacketsSent: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "network", "transmit_packets_total"),
			"Total number of packets transmitted from device.",
			[]string{"device"},
			nil,
		),
		netPacketsRecv: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "network", "receive_packets_total"),
			"Total number of packets received by device.",
			[]string{"device"},
			nil,
		),
		netErrorsSent: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "network", "transmit_errors_total"),
			"Total number of errors while transmitting from device.",
			[]string{"device", "error"},
			nil,
		),
		netErrorsRecv: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "network", "receive_errors_total"),
			"Total number of errors while receiving by device.",
			[]string{"device", "error"},
			nil,
		),
	}, nil
}

func (c *gopsCollector) Update(ch chan<- prometheus.Metric) (err error) {
	for _, uf := range []func(ch chan<- prometheus.Metric) error{
		c.updateCPU,
		c.updateMemory,
		c.updateSwap,
		c.updateHost,
		c.updateNet,
	} {
		if err := uf(ch); err != nil {
			if strings.Contains(err.Error(), "not implemented") {
				continue // continue with next on not-implemented
			}
			return err
		}
	}
	return nil
}

func (c *gopsCollector) updateCPU(ch chan<- prometheus.Metric) error {
	times, err := cpu.Times(true)
	if err != nil {
		return err
	}
	for _, t := range times {
		for k, v := range map[string]float64{
			"user":      t.User,
			"system":    t.System,
			"idle":      t.Idle,
			"nice":      t.Nice,
			"iowait":    t.Iowait,
			"irq":       t.Irq,
			"softirq":   t.Softirq,
			"steal":     t.Steal,
			"guest":     t.Guest,
			"guestNice": t.GuestNice,
			"stolen":    t.Stolen,
		} {
			ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, v, t.CPU, k)
		}
	}
	return nil
}

// Memory
func (c *gopsCollector) updateMemory(ch chan<- prometheus.Metric) error {
	vmstat, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	for k, v := range map[string]uint64{
		"used":     vmstat.Used,
		"free":     vmstat.Free,
		"active":   vmstat.Active,
		"inactive": vmstat.Inactive,
		"wired":    vmstat.Wired,
		"buffers":  vmstat.Buffers,
		"cached":   vmstat.Cached,
	} {
		ch <- prometheus.MustNewConstMetric(c.mem, prometheus.GaugeValue, float64(v), k)
	}
	return nil
}

func (c *gopsCollector) updateSwap(ch chan<- prometheus.Metric) error {
	// - Swap
	sstat, err := mem.SwapMemory()
	if err != nil {
		return err
	}
	for k, v := range map[string]uint64{
		"used": sstat.Used,
		"free": sstat.Free,
	} {
		ch <- prometheus.MustNewConstMetric(c.swap, prometheus.GaugeValue, float64(v), k)
	}
	ch <- prometheus.MustNewConstMetric(c.swapIn, prometheus.CounterValue, float64(sstat.Sin))
	ch <- prometheus.MustNewConstMetric(c.swapOut, prometheus.CounterValue, float64(sstat.Sout))

	return nil
}

func (c *gopsCollector) updateHost(ch chan<- prometheus.Metric) error {
	// Host
	bt, err := host.BootTime()
	if err != nil {
		return err
	}
	ch <- prometheus.MustNewConstMetric(c.bootTime, prometheus.GaugeValue, float64(bt))

	return nil
}

func (c *gopsCollector) updateNet(ch chan<- prometheus.Metric) error {
	// Net
	netStats, err := net.IOCounters(true)
	if err != nil {
		return err
	}
	for _, stat := range netStats {
		if c.ignoredNetDevices.MatchString(stat.Name) {
			continue
		}
		for m, v := range map[*prometheus.Desc]uint64{
			c.netBytesSent: stat.BytesSent,
			c.netBytesRecv: stat.BytesRecv,

			c.netPacketsSent: stat.PacketsSent,
			c.netPacketsRecv: stat.PacketsRecv,
		} {
			ch <- prometheus.MustNewConstMetric(m, prometheus.CounterValue, float64(v), stat.Name)
		}
		ch <- prometheus.MustNewConstMetric(c.netErrorsSent, prometheus.CounterValue, float64(stat.Errout), stat.Name, "error")
		ch <- prometheus.MustNewConstMetric(c.netErrorsRecv, prometheus.CounterValue, float64(stat.Errin), stat.Name, "error")
		ch <- prometheus.MustNewConstMetric(c.netErrorsSent, prometheus.CounterValue, float64(stat.Dropout), stat.Name, "drop")
		ch <- prometheus.MustNewConstMetric(c.netErrorsRecv, prometheus.CounterValue, float64(stat.Dropin), stat.Name, "drop")
	}

	return nil
}
