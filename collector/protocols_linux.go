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

package collector

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	subsystem = "protocols"
)

var (
	netProtocolFilter = kingpin.Flag("collector.protocols.filter", "Regex of protocols to return for the collector.").Default("^tcp.*").String()
)

type protocolsCollector struct {
	fs     procfs.FS
	logger log.Logger
}

func init() {
	registerCollector(subsystem, defaultEnabled, NewProtocolsCollector)
}

// NewProtocolsCollector returns a Collector exposing net/protocols stats
func NewProtocolsCollector(logger log.Logger) (Collector, error) {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}
	return &protocolsCollector{fs, logger}, nil
}

// Update implements Collector and exposes /proc/net/protocols metrics
func (c *protocolsCollector) Update(ch chan<- prometheus.Metric) error {

	protocolStats, err := c.fs.NetProtocols()
	if err != nil {
		return fmt.Errorf("couldn't get protocols: %w", err)
	}

	// In the interest of reudcing cardinality we are only interested in the
	// first 8 fields as subsequent fields are not numerical or likely to change
	// over time.
	for _, p := range protocolStats {
		p.Name = strings.Replace(strings.ToLower(p.Name), "-", "", -1)
		re := regexp.MustCompile(*netProtocolFilter)
		if !re.MatchString(p.Name) {
			continue
		}

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(
					namespace,
					subsystem,
					fmt.Sprintf("%s_size", p.Name),
				),
				fmt.Sprintf("The size, in bytes, of %s", p.Name),
				nil,
				nil,
			),
			prometheus.GaugeValue,
			float64(p.Size),
		)

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(
					namespace,
					subsystem,
					fmt.Sprintf("%s_sockets", p.Name),
				),
				fmt.Sprintf("Number of sockets in use by %s", p.Name),
				nil,
				nil,
			),
			prometheus.GaugeValue,
			float64(p.Sockets),
		)

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(
					namespace,
					subsystem,
					fmt.Sprintf("%s_memory", p.Name),
				),
				fmt.Sprintf("Total number of 4KB pages allocated by %s", p.Name),
				nil,
				nil,
			),
			prometheus.GaugeValue,
			float64(p.Memory),
		)

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(
					namespace,
					subsystem,
					fmt.Sprintf("%s_pressure", p.Name),
				),
				fmt.Sprintf("Indicates whether %s is experiencing memory pressure; 1 = true, 0 = false, -1 = not implemented", p.Name),
				nil,
				nil,
			),
			prometheus.GaugeValue,
			float64(p.Pressure),
		)

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(
					namespace,
					subsystem,
					fmt.Sprintf("%s_maxhdr", p.Name),
				),
				fmt.Sprintf("Max header size for %s", p.Name),
				nil,
				nil,
			),
			prometheus.GaugeValue,
			float64(p.MaxHeader),
		)
		slabVal := float64(0.0)
		if p.Slab {
			slabVal = 1.0
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(
					namespace,
					subsystem,
					fmt.Sprintf("%s_slab", p.Name),
				),
				fmt.Sprintf("Bool indicating whether %s is allocated from SLAB", p.Name),
				nil,
				nil,
			),
			prometheus.GaugeValue,
			slabVal,
		)
	}
	return nil
}
