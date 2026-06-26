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
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	edacSubsystem = "edac"
)

var (
	edacMemControllerRE = regexp.MustCompile(`.*devices/system/edac/mc/mc([0-9]*)`)
	edacMemDimmRE       = regexp.MustCompile(`.*devices/system/edac/mc/mc[0-9]*/dimm([0-9]*)`)
)

type edacCollector struct {
	ceCount        *prometheus.Desc
	ueCount        *prometheus.Desc
	channelCECount *prometheus.Desc
	channelUECount *prometheus.Desc
	logger         *slog.Logger
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

		ceCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, edacSubsystem, "correctable_errors_total"),
			"Total correctable memory errors.",
			[]string{"controller"},
			nil,
		),

		ueCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, edacSubsystem, "uncorrectable_errors_total"),
			"Total uncorrectable memory errors.",
			[]string{"controller"},
			nil,
		),

		channelCECount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, edacSubsystem, "channel_correctable_errors_total"),
			"Total correctable memory errors for this channel.",
			[]string{"controller", "csrow", "channel", "dimm_label"},
			nil,
		),

		channelUECount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, edacSubsystem, "channel_uncorrectable_errors_total"),
			"Total uncorrectable memory errors for this channel.",
			[]string{"controller", "csrow", "channel", "dimm_label"},
			nil,
		),

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
		if err == nil {
			ch <- prometheus.MustNewConstMetric(
				c.ceCount,
				prometheus.CounterValue,
				float64(value),
				controllerNumber,
			)
		}

		value, err = readUintFromFile(filepath.Join(controller, "ue_count"))
		if err == nil {
			ch <- prometheus.MustNewConstMetric(
				c.ueCount,
				prometheus.CounterValue,
				float64(value),
				controllerNumber,
			)
		}
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

		csrows, err := filepath.Glob(controller + "/csrow[0-9]*")

		if err != nil {
			return err
		}

		for _, csrow := range csrows {
			base := filepath.Base(csrow)

			match := regexp.MustCompile(`csrow([0-9]+)`).FindStringSubmatch(base)
			if match == nil {
				continue
			}
			csrowNumber := match[1]

			channelFiles, err := filepath.Glob(csrow + "/ch*_ce_count")
			if err != nil {
				return err
			}
			ch <- prometheus.MustNewConstMetric(
				edacCsRowCECount, prometheus.CounterValue, float64(value), controllerNumber, csrowNumber)

			for _, chFile := range channelFiles {

				base := filepath.Base(chFile)

				match := regexp.MustCompile(`ch([0-9]+)_ce_count`).FindStringSubmatch(base)
				if match == nil {
					continue
				}

				channelNumber := match[1]
				label := "unknown"
				labelBytes, err := os.ReadFile(filepath.Join(csrow, "ch"+channelNumber+"_dimm_label"))
				if err == nil {
					label = strings.TrimSpace(string(labelBytes))
					// format label
					label = strings.ReplaceAll(label, "#", "")
					label = strings.ReplaceAll(label, "csrow", "_csrow")
					label = strings.ReplaceAll(label, "channel", "_channel")
				}
				value, err := readUintFromFile(chFile)
				if err == nil {
					ch <- prometheus.MustNewConstMetric(
						c.channelCECount,
						prometheus.CounterValue,
						float64(value),
						controllerNumber,
						csrowNumber,
						channelNumber,
						label,
					)
				}

				value, err = readUintFromFile(filepath.Join(csrow, "ch"+channelNumber+"_ue_count"))
				if err == nil {
					ch <- prometheus.MustNewConstMetric(
						c.channelUECount,
						prometheus.CounterValue,
						float64(value),
						controllerNumber,
						csrowNumber,
						channelNumber,
						label,
					)
				}
			}
			ch <- prometheus.MustNewConstMetric(
				edacCsRowUECount, prometheus.CounterValue, float64(value), controllerNumber, csrowNumber)
		}
	}

	return nil
}
