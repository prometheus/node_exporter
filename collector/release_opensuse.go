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

// +build !novmstat

package collector

import (
	"bufio"
	// 	"fmt"
	"os"
	// 	"regexp"
	// 	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	// 	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	releaseSubsystem = "release"
)

type releaseCollector struct{}
type release struct {
	ID        string
	VersionID string
}

func init() {
	registerCollector("release", defaultEnabled, newReleaseCollector)
}

func newReleaseCollector() (Collector, error) {
	return &releaseCollector{}, nil
}

func (c *releaseCollector) Update(ch chan<- prometheus.Metric) error {
	release, err := getRelease()
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, releaseSubsystem, "info"),
			"/etc/os-release information",
			[]string{"id", "versionid"}, nil),
		prometheus.GaugeValue,
		1,
		release.ID,
		release.VersionID,
	)
	return nil
}

func getRelease() (release, error) {
	release := release{
		ID:        "",
		VersionID: "",
	}
	file, err := os.Open(etcFilePath("os-release"))
	if err != nil {
		return release, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "=")
		if "ID" == parts[0] {
			release.ID = parts[1]
		} else if "VERSION_ID" == parts[0] {
			release.VersionID = parts[1]
		}
	}

	return release, scanner.Err()
}
