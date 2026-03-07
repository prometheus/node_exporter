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

//go:build !noedac

package collector

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	edacSubsystem = "edac"
)

var (
	edacMemControllerRE = regexp.MustCompile(`.*devices/system/edac/mc/mc([0-9]*)`)
	edacMemCsrowRE      = regexp.MustCompile(`.*devices/system/edac/mc/mc[0-9]*/csrow([0-9]*)`)
)

type edacCollector struct {
	logger *slog.Logger
}

func init() {
	registerCollector("edac", defaultEnabled, NewEdacCollector)
}

var (
	edacCeCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, edacSubsystem, "correctable_errors_total"),
		"Total correctable memory errors.",
		[]string{"controller"}, nil,
	)
	edacUeCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, edacSubsystem, "uncorrectable_errors_total"),
		"Total uncorrectable memory errors.",
		[]string{"controller"}, nil,
	)
	edacCsRowCECount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, edacSubsystem, "csrow_correctable_errors_total"),
		"Total correctable memory errors for this csrow.",
		[]string{"controller", "csrow"}, nil,
	)
	edacCsRowUECount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, edacSubsystem, "csrow_uncorrectable_errors_total"),
		"Total uncorrectable memory errors for this csrow.",
		[]string{"controller", "csrow"}, nil,
	)
)

// NewEdacCollector returns a new Collector exposing edac stats.
func NewEdacCollector(logger *slog.Logger) (Collector, error) {
	return &edacCollector{
		logger: logger,
	}, nil
}

func (c *edacCollector) Update(ch chan<- prometheus.Metric) error {
	memControllers, err := filepath.Glob(sysFilePath("devices/system/edac/mc/mc[0-9]*"))
	if err != nil {
		return err
	}
	for _, controller := range memControllers {
		controllerMatch := edacMemControllerRE.FindStringSubmatch(controller)
		if controllerMatch == nil {
			return fmt.Errorf("controller string didn't match regexp: %s", controller)
		}
		controllerNumber := controllerMatch[1]

		value, err := readUintFromFile(filepath.Join(controller, "ce_count"))
		if err != nil {
			return fmt.Errorf("couldn't get ce_count for controller %s: %w", controllerNumber, err)
		}
		ch <- prometheus.MustNewConstMetric(
			edacCeCount, prometheus.CounterValue, float64(value), controllerNumber)

		value, err = readUintFromFile(filepath.Join(controller, "ce_noinfo_count"))
		if err != nil {
			return fmt.Errorf("couldn't get ce_noinfo_count for controller %s: %w", controllerNumber, err)
		}
		ch <- prometheus.MustNewConstMetric(
			edacCsRowCECount, prometheus.CounterValue, float64(value), controllerNumber, "unknown")

		value, err = readUintFromFile(filepath.Join(controller, "ue_count"))
		if err != nil {
			return fmt.Errorf("couldn't get ue_count for controller %s: %w", controllerNumber, err)
		}
		ch <- prometheus.MustNewConstMetric(
			edacUeCount, prometheus.CounterValue, float64(value), controllerNumber)

		value, err = readUintFromFile(filepath.Join(controller, "ue_noinfo_count"))
		if err != nil {
			return fmt.Errorf("couldn't get ue_noinfo_count for controller %s: %w", controllerNumber, err)
		}
		ch <- prometheus.MustNewConstMetric(
			edacCsRowUECount, prometheus.CounterValue, float64(value), controllerNumber, "unknown")

		// For each controller, walk the csrow directories.
		csrows, err := filepath.Glob(controller + "/csrow[0-9]*")
		if err != nil {
			return err
		}
		for _, csrow := range csrows {
			csrowMatch := edacMemCsrowRE.FindStringSubmatch(csrow)
			if csrowMatch == nil {
				return fmt.Errorf("csrow string didn't match regexp: %s", csrow)
			}
			csrowNumber := csrowMatch[1]

			value, err = readUintFromFile(filepath.Join(csrow, "ce_count"))
			if err != nil {
				return fmt.Errorf("couldn't get ce_count for controller/csrow %s/%s: %w", controllerNumber, csrowNumber, err)
			}
			ch <- prometheus.MustNewConstMetric(
				edacCsRowCECount, prometheus.CounterValue, float64(value), controllerNumber, csrowNumber)

			value, err = readUintFromFile(filepath.Join(csrow, "ue_count"))
			if err != nil {
				return fmt.Errorf("couldn't get ue_count for controller/csrow %s/%s: %w", controllerNumber, csrowNumber, err)
			}
			ch <- prometheus.MustNewConstMetric(
				edacCsRowUECount, prometheus.CounterValue, float64(value), controllerNumber, csrowNumber)
		}
	}

	return err
}
