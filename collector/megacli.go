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

// +build !nomegacli

package collector

import (
	"bufio"
	"flag"
	"io"
	"os/exec"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	defaultMegaCli   = "megacli"
	adapterHeaderSep = "================"
)

var (
	megacliCommand = flag.String("collector.megacli.command", defaultMegaCli, "Command to run megacli.")
)

type megaCliCollector struct {
	cli string

	driveTemperature *prometheus.GaugeVec
	driveCounters    *prometheus.CounterVec
	drivePresence    *prometheus.GaugeVec
}

func init() {
	Factories["megacli"] = NewMegaCliCollector
}

// Takes a prometheus registry and returns a new Collector exposing
// RAID status through megacli.
func NewMegaCliCollector() (Collector, error) {
	return &megaCliCollector{
		cli: *megacliCommand,
		driveTemperature: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "megacli_drive_temperature_celsius",
			Help:      "megacli: drive temperature",
		}, []string{"enclosure", "slot"}),
		driveCounters: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      "megacli_drive_count",
			Help:      "megacli: drive error and event counters",
		}, []string{"enclosure", "slot", "type"}),
		drivePresence: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "megacli_adapter_disk_presence",
			Help:      "megacli: disk presence per adapter",
		}, []string{"type"}),
	}, nil
}

func (c *megaCliCollector) Update(ch chan<- prometheus.Metric) (err error) {
	err = c.updateAdapter()
	if err != nil {
		return err
	}
	err = c.updateDisks()
	c.driveTemperature.Collect(ch)
	c.driveCounters.Collect(ch)
	c.drivePresence.Collect(ch)
	return err
}

func parseMegaCliDisks(r io.Reader) (map[int]map[int]map[string]string, error) {
	var (
		stats   = map[int]map[int]map[string]string{}
		scanner = bufio.NewScanner(r)
		curEnc  = -1
		curSlot = -1
	)

	for scanner.Scan() {
		var err error
		text := strings.TrimSpace(scanner.Text())
		parts := strings.SplitN(text, ":", 2)
		if len(parts) != 2 { // Adapter #X
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		switch {
		case key == "Enclosure Device ID":
			curEnc, err = strconv.Atoi(value)
			if err != nil {
				return nil, err
			}
		case key == "Slot Number":
			curSlot, err = strconv.Atoi(value)
			if err != nil {
				return nil, err
			}
		case curSlot != -1 && curEnc != -1:
			if _, ok := stats[curEnc]; !ok {
				stats[curEnc] = map[int]map[string]string{}
			}
			if _, ok := stats[curEnc][curSlot]; !ok {
				stats[curEnc][curSlot] = map[string]string{}
			}
			stats[curEnc][curSlot][key] = value
		}
	}

	return stats, nil
}

func parseMegaCliAdapter(r io.Reader) (map[string]map[string]string, error) {
	var (
		raidStats = map[string]map[string]string{}
		scanner   = bufio.NewScanner(r)
		header    = ""
		last      = ""
	)

	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text == adapterHeaderSep {
			header = last
			raidStats[header] = map[string]string{}
			continue
		}
		last = text
		if header == "" { // skip Adapter #X and separator
			continue
		}
		parts := strings.SplitN(text, ":", 2)
		if len(parts) != 2 { // these section never include anything we are interested in
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		raidStats[header][key] = value

	}

	return raidStats, nil
}

func (c *megaCliCollector) updateAdapter() error {
	cmd := exec.Command(c.cli, "-AdpAllInfo", "-aALL")
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	stats, err := parseMegaCliAdapter(pipe)
	if err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}

	for k, v := range stats["Device Present"] {
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}
		c.drivePresence.WithLabelValues(k).Set(value)
	}
	return nil
}

func (c *megaCliCollector) updateDisks() error {
	var counters = []string{"Media Error Count", "Other Error Count", "Predictive Failure Count"}

	cmd := exec.Command(c.cli, "-PDList", "-aALL")
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	stats, err := parseMegaCliDisks(pipe)
	if err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}

	for enc, encStats := range stats {
		for slot, slotStats := range encStats {
			encStr := strconv.Itoa(enc)
			slotStr := strconv.Itoa(slot)

			tStr := slotStats["Drive Temperature"]
			if strings.Index(tStr, "C") > 0 {
				tStr = tStr[:strings.Index(tStr, "C")]
				t, err := strconv.ParseFloat(tStr, 64)
				if err != nil {
					return err
				}
				c.driveTemperature.WithLabelValues(encStr, slotStr).Set(t)
			}

			for _, i := range counters {
				counter, err := strconv.ParseFloat(slotStats[i], 64)
				if err != nil {
					return err
				}

				c.driveCounters.WithLabelValues(encStr, slotStr, i).Set(counter)
			}
		}
	}
	return nil
}
