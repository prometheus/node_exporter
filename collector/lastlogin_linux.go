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

// +build !nolastlogin

package collector

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const lastLoginSubsystem = "last_login"

type lastLoginCollector struct {
	metric prometheus.Gauge
}

func init() {
	Factories["lastlogin"] = NewLastLoginCollector
}

// Takes a prometheus registry and returns a new Collector exposing
// load, seconds since last login and a list of tags as specified by config.
func NewLastLoginCollector() (Collector, error) {
	return &lastLoginCollector{
		metric: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: lastLoginSubsystem,
			Name:      "time",
			Help:      "The time of the last login.",
		}),
	}, nil
}

func (c *lastLoginCollector) Update(ch chan<- prometheus.Metric) (err error) {
	last, err := getLastLoginTime()
	if err != nil {
		return fmt.Errorf("couldn't get last seen: %s", err)
	}
	log.Debugf("Set node_last_login_time: %f", last)
	c.metric.Set(last)
	c.metric.Collect(ch)
	return err
}

func getLastLoginTime() (float64, error) {
	who := exec.Command("who", "/var/log/wtmp", "-l", "-u", "-s")

	output, err := who.StdoutPipe()
	if err != nil {
		return 0, err
	}

	err = who.Start()
	if err != nil {
		return 0, err
	}

	reader := bufio.NewReader(output)

	var last time.Time
	for {
		line, isPrefix, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if isPrefix {
			return 0, fmt.Errorf("line to long: %s(...)", line)
		}

		fields := strings.Fields(string(line))
		lastDate := fields[2]
		lastTime := fields[3]

		dateParts, err := splitToInts(lastDate, "-") // 2013-04-16
		if err != nil {
			return 0, fmt.Errorf("couldn't parse date in line '%s': %s", fields, err)
		}

		timeParts, err := splitToInts(lastTime, ":") // 11:33
		if err != nil {
			return 0, fmt.Errorf("couldn't parse time in line '%s': %s", fields, err)
		}

		last_t := time.Date(dateParts[0], time.Month(dateParts[1]), dateParts[2], timeParts[0], timeParts[1], 0, 0, time.UTC)
		last = last_t
	}
	err = who.Wait()
	if err != nil {
		return 0, err
	}

	return float64(last.Unix()), nil
}
