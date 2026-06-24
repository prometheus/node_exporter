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
	edacMemCsrowRE      = regexp.MustCompile(`.*devices/system/edac/mc/mc([0-9]*)/csrow([0-9]*)`)
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

	edacChannelCECount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, edacSubsystem, "channel_correctable_errors_total"),
		"Total correctable memory errors for this channel.",
		[]string{"controller", "csrow", "channel", "dimm_label"}, nil,
	)

	edacChannelUECount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, edacSubsystem, "channel_uncorrectable_errors_total"),
		"Total uncorrectable memory errors for this channel.",
		[]string{"controller", "csrow", "channel", "dimm_label"}, nil,
	)
)

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

		// Controller CE count
		value, err := readUintFromFile(filepath.Join(controller, "ce_count"))
		if err != nil {
			return fmt.Errorf("couldn't get ce_count for controller %s: %w", controllerNumber, err)
		}

		ch <- prometheus.MustNewConstMetric(
			edacCeCount,
			prometheus.CounterValue,
			float64(value),
			controllerNumber,
		)

		// Controller UE count
		value, err = readUintFromFile(filepath.Join(controller, "ue_count"))
		if err != nil {
			return fmt.Errorf("couldn't get ue_count for controller %s: %w", controllerNumber, err)
		}

		ch <- prometheus.MustNewConstMetric(
			edacUeCount,
			prometheus.CounterValue,
			float64(value),
			controllerNumber,
		)

		csrows, err := filepath.Glob(controller + "/csrow[0-9]*")
		if err != nil {
			return err
		}

		for _, csrow := range csrows {
			csrowMatch := edacMemCsrowRE.FindStringSubmatch(csrow)
			if csrowMatch == nil {
				return fmt.Errorf("csrow string didn't match regexp: %s", csrow)
			}

			csrowNumber := csrowMatch[2]

			channelFiles, err := filepath.Glob(csrow + "/ch[0-9]*_ce_count")
			if err != nil {
				return err
			}

			for _, chFile := range channelFiles {
				match := regexp.MustCompile(`ch([0-9]+)_ce_count`).FindStringSubmatch(filepath.Base(chFile))
				if match == nil {
					continue
				}

				channelNumber := match[1]

				label := fmt.Sprintf(
					"mc%s_csrow%s_channel%s",
					controllerNumber,
					csrowNumber,
					channelNumber,
				)

				value, err := readUintFromFile(chFile)
				if err == nil {
					ch <- prometheus.MustNewConstMetric(
						edacChannelCECount,
						prometheus.CounterValue,
						float64(value),
						controllerNumber,
						csrowNumber,
						channelNumber,
						label,
					)
				}

				ueFile := filepath.Join(
					csrow,
					fmt.Sprintf("ch%s_ue_count", channelNumber),
				)

				value, err = readUintFromFile(ueFile)
				if err == nil {
					ch <- prometheus.MustNewConstMetric(
						edacChannelUECount,
						prometheus.CounterValue,
						float64(value),
						controllerNumber,
						csrowNumber,
						channelNumber,
						label,
					)
				}
			}
		}
	}

	return nil
}
