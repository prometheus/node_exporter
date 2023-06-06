// Copyright 2021 The Prometheus Authors
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
	"github.com/alecthomas/kingpin/v2"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	azureProvider = "azure"
)

var (
	provider = kingpin.Flag(
		"collector.ims.provider",
		"Name of the cloud provider to select instance meta service (azure)",
	).Default("disabled").String()
)

type disabledInstanceMetadataCollector struct {
}

func init() {
	registerCollector("ims", defaultEnabled, NewImsCollector)
}

// NewImsCollector returna a new Collector exposing instance metadata services information.
func NewImsCollector(logger log.Logger) (Collector, error) {
	level.Info(logger).Log("msg", "Parsed flag --collector.ims.provider", "flag", *provider)

	description := "A metric with a constant '1' value labelled by"
	buildFqdn := prometheus.BuildFQName(namespace, "ims", "info")

	switch *provider {
	case azureProvider:
		return NewAzureInstanceMetadataCollector(description, buildFqdn, logger), nil
	default:
		return &disabledInstanceMetadataCollector{}, nil
	}
}

// Default instance metadata collector, do nothing if collector is not configured.
func (c *disabledInstanceMetadataCollector) Update(ch chan<- prometheus.Metric) error {
	return nil
}
