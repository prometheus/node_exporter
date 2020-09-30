// Copyright 2020 The Prometheus Authors
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

// +build !nozoneinfo

package collector

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-kit/kit/log/level"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	zoneinfoSubsystem      = "zoneinfo"
	badLengthErrorMessage  = "cannot parse %s. Skipping"
	convertionErrorMessage = "cannot parse %s: %v. Skipping"
)

var (
	zoneinfoPath = kingpin.Flag("collector.zoneinfo.path", "path to zoneinfo").Default("/proc/zoneinfo").String()
)

type zoneinfoCollector struct {
	metricDescs map[string]*prometheus.Desc
	logger      log.Logger
}

type Zone struct {
	Node              string
	Zone              string
	PerNodeStats      map[string]uint64
	Pages             map[string]uint64
	PageSets          []PageSet
	NodeUnreclaimable uint64
	StartPfm          uint64
}

type PageSet struct {
	CPU             uint64
	Count           uint64
	High            uint64
	Batch           uint64
	VMStatThreshold uint64
}

func init() {
	registerCollector("zoneinfo", defaultDisabled, NewZoneinfoCollector)
}

// NewZoneinfoCollector returns a new Collector exposing zone info metrics.
func NewZoneinfoCollector(logger log.Logger) (Collector, error) {
	return &zoneinfoCollector{
		metricDescs: map[string]*prometheus.Desc{},
		logger:      logger,
	}, nil
}

func (c *zoneinfoCollector) Update(ch chan<- prometheus.Metric) error {
	zones, err := getZoneInfo(c.logger)
	if err != nil {
		return fmt.Errorf("couldn't get zoneinfo: %v", err)
	}
	for _, zone := range zones {
		for metric, value := range zone.PerNodeStats {
			desc, ok := c.metricDescs[metric]
			if !ok {
				desc = prometheus.NewDesc(
					prometheus.BuildFQName(namespace, zoneinfoSubsystem, metric),
					fmt.Sprintf("Zone information field %s.", metric),
					[]string{"node", "zone"}, nil)
				c.metricDescs[metric] = desc
			}
			ch <- prometheus.MustNewConstMetric(desc, prometheus.UntypedValue, float64(value), zone.Node, zone.Zone)
		}
		for metric, value := range zone.Pages {
			metricName := "page_" + metric
			desc, ok := c.metricDescs[metric]
			if !ok {
				desc = prometheus.NewDesc(
					prometheus.BuildFQName(namespace, zoneinfoSubsystem, metricName),
					fmt.Sprintf("Zone information field %s.", metricName),
					[]string{"node", "zone"}, nil)
				c.metricDescs[metricName] = desc
			}
			ch <- prometheus.MustNewConstMetric(desc, prometheus.UntypedValue, float64(value), zone.Node, zone.Zone)
		}
		for _, pageSet := range zone.PageSets {
			desc, ok := c.metricDescs["PageSet"]
			if !ok {
				desc = prometheus.NewDesc(
					prometheus.BuildFQName(namespace, zoneinfoSubsystem, "pageset"),
					fmt.Sprintf("Zone information field %s.", "PageSet"),
					[]string{"node", "zone", "cpu", "high", "batch", "vm_stat_threshold"}, nil)
				c.metricDescs["PageSet"] = desc
			}
			ch <- prometheus.MustNewConstMetric(desc, prometheus.UntypedValue, float64(pageSet.Count),
				zone.Node, zone.Zone, strconv.Itoa(int(pageSet.CPU)), strconv.Itoa(int(pageSet.High)),
				strconv.Itoa(int(pageSet.Batch)), strconv.Itoa(int(pageSet.VMStatThreshold)))
		}
		metricName := "node_unreclaimable"
		desc, ok := c.metricDescs[metricName]
		if !ok {
			desc = prometheus.NewDesc(
				prometheus.BuildFQName(namespace, zoneinfoSubsystem, metricName),
				fmt.Sprintf("Zone information field %s.", metricName),
				[]string{"node", "zone"}, nil)
			c.metricDescs[metricName] = desc
		}
		ch <- prometheus.MustNewConstMetric(desc, prometheus.UntypedValue, float64(zone.NodeUnreclaimable), zone.Node, zone.Zone)

		metricName = "start_pfn"
		desc, ok = c.metricDescs[metricName]
		if !ok {
			desc = prometheus.NewDesc(
				prometheus.BuildFQName(namespace, zoneinfoSubsystem, metricName),
				fmt.Sprintf("Zone information field %s.", metricName),
				[]string{"node", "zone"}, nil)
			c.metricDescs[metricName] = desc
		}
		ch <- prometheus.MustNewConstMetric(desc, prometheus.UntypedValue, float64(zone.StartPfm), zone.Node, zone.Zone)
	}

	return nil
}

func getZoneInfo(logger log.Logger) ([]Zone, error) {
	file, err := os.Open(*zoneinfoPath)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	var (
		zones          []Zone
		zone           Zone
		lastPageSet    PageSet
		currentSection string
	)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "Node"):
			splittedLine := strings.SplitN(line, ",", 2)
			if len(splittedLine) != 2 {
				return nil, fmt.Errorf(badLengthErrorMessage, line)
			}
			nodeNumber := strings.Fields(splittedLine[0])[1]
			zoneName := strings.Fields(splittedLine[1])[1]
			zones = append(zones, zone)
			zone = Zone{
				Node:         nodeNumber,
				Zone:         zoneName,
				PerNodeStats: map[string]uint64{},
				Pages:        map[string]uint64{},
			}
			currentSection = "Node"
			continue

		case strings.HasPrefix(line, "  per-node stats"):
			currentSection = "per-node stats"
			continue

		case strings.HasPrefix(line, "  pages "):
			currentSection = "pages"
			// This case will look like: `pages free <value>` where pages is section name and free is key
			splittedLine := strings.Fields(line)
			if len(splittedLine) != 3 {
				level.Error(logger).Log("msg", fmt.Sprintf(badLengthErrorMessage, line))
				continue
			}
			value, err := strconv.Atoi(splittedLine[2])
			if err != nil {
				level.Error(logger).Log("msg", fmt.Sprintf(convertionErrorMessage, line, err))
				continue
			}
			zone.Pages[splittedLine[1]] = uint64(value)
			continue

		case strings.HasPrefix(line, "  pagesets"):
			currentSection = "pagesets"
			continue

		case strings.HasPrefix(line, "  node_unreclaimable:"):
			_, value, err := parseLine(line, 2)
			if err != nil {
				level.Error(logger).Log("msg", err)
				continue
			}
			zone.NodeUnreclaimable = uint64(value)
			continue

		case strings.HasPrefix(line, "  start_pfn:"):
			_, value, err := parseLine(line, 2)
			if err != nil {
				level.Error(logger).Log("msg", err)
				continue
			}
			zone.StartPfm = uint64(value)
			continue
		}

		switch currentSection {
		case "per-node stats":
			key, value, err := parseLine(line, 2)
			if err != nil {
				level.Error(logger).Log("msg", err)
				continue
			}
			zone.PerNodeStats[key] = uint64(value)

		case "pages":
			splittedLine := strings.Fields(line)
			// special case, value is list instead of uint and after that it goes back to per-node stats section
			if splittedLine[0] == "protection:" {
				for i, value := range splittedLine[1:] {
					intValue, err := strconv.Atoi(strings.Trim(value, "(),"))
					if err != nil {
						level.Error(logger).Log("msg", fmt.Sprintf(convertionErrorMessage, line, err))
					}
					zone.Pages["protection_"+strconv.Itoa(i)] = uint64(intValue)
				}
				currentSection = "per-node stats"
				continue
			}
			key, value, err := parseLine(line, 2)
			if err != nil {
				level.Error(logger).Log("msg", err)
				continue
			}
			zone.Pages[key] = uint64(value)
		case "pagesets":
			// special case, stat name might have spaces in it
			splittedLine := strings.Fields(line)
			if len(splittedLine) > 1 {
				if splittedLine[0] == "vm" && len(splittedLine) != 4 || splittedLine[0] != "vm" && len(splittedLine) != 2 {
					level.Error(logger).Log("msg", fmt.Sprintf(badLengthErrorMessage, line))
					continue
				}
			} else {
				level.Error(logger).Log("msg", fmt.Sprintf(badLengthErrorMessage, line))
				continue
			}
			value, err := strconv.Atoi(splittedLine[len(splittedLine)-1])
			if err != nil {
				level.Error(logger).Log("msg", fmt.Sprintf(convertionErrorMessage, line, err))
				continue
			}
			switch splittedLine[0] {
			case "cpu:":
				lastPageSet = PageSet{
					CPU: uint64(value),
				}
			case "count:":
				lastPageSet.Count = uint64(value)

			case "high:":
				lastPageSet.High = uint64(value)

			case "batch:":
				lastPageSet.Batch = uint64(value)

			case "vm":
				lastPageSet.VMStatThreshold = uint64(value)
				zone.PageSets = append(zone.PageSets, lastPageSet)
			}
		}
	}
	zones = append(zones, zone)
	return zones[1:], nil
}

func parseLine(line string, expectedLength int) (string, int, error) {
	splittedLine := strings.Fields(line)
	if len(splittedLine) != expectedLength {
		return "", 0, fmt.Errorf(badLengthErrorMessage, line)
	}
	value, err := strconv.Atoi(splittedLine[expectedLength-1])
	if err != nil {
		return "", 0, fmt.Errorf(convertionErrorMessage, line, err)
	}
	key := strings.Join(splittedLine[:expectedLength-1], "_")
	return key, value, nil
}
