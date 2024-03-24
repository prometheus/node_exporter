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
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	azureComputeURL = "http://169.254.169.254/metadata/instance/compute?api-version=2021-01-01&format=json"
)

type AzureInstanceMetadataCollector struct {
	InfoDesc   *prometheus.Desc
	labels     []string
	mapping    map[string]string
	logger     log.Logger
	computeAPI string
}

func NewAzureInstanceMetadataCollector(description string, fqdn string, logger log.Logger) *AzureInstanceMetadataCollector {
	mapping := map[string]string{
		"vm_sku":   "vmSize",
		"vm_zone":  "zone",
		"provider": "provider",
		"location": "location",
	}

	labels := make([]string, 0, len(mapping))

	for k := range mapping {
		labels = append(labels, k)
	}

	sort.Strings(labels)

	infoDesc := prometheus.NewDesc(
		fqdn,
		fmt.Sprintf("%s: %s .", description, strings.Join(labels, ", ")),
		labels, nil)

	return &AzureInstanceMetadataCollector{
		InfoDesc:   infoDesc,
		labels:     labels,
		mapping:    mapping,
		logger:     logger,
		computeAPI: azureComputeURL,
	}
}

func (c *AzureInstanceMetadataCollector) GetMetaDataServiceInfo() (map[string]interface{}, error) {
	// Disable proxy authentication from host system
	transport := http.DefaultTransport
	transport.(*http.Transport).Proxy = nil

	client := http.Client{
		Transport: transport,
	}

	req, err := http.NewRequest("GET", c.computeAPI, nil)
	if err != nil {
		return nil, fmt.Errorf("could not build http client to fetch auzure instance meta data service for compute")
	}

	req.Header = http.Header{
		"Metadata":     {"true"},
		"Content-Type": {"application/json"},
	}

	res, err := client.Do(req)
	if err != nil {
		level.Debug(c.logger).Log(
			"msg", "Error querying azure instance meta data service",
			"compute api", c.computeAPI,
			"err", err)
		return nil, fmt.Errorf("could not fetch azure compute instance metadata service")
	}

	defer res.Body.Close()
	var metadata map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&metadata); err != nil {
		level.Debug(c.logger).Log(
			"msg", "Error could not decode response body from azure instance meta data service",
			"compute api", c.computeAPI,
			"err", err)
		return nil, fmt.Errorf("could not decode azure compute instance metadata service")
	}

	return metadata, nil
}

func (c *AzureInstanceMetadataCollector) Update(ch chan<- prometheus.Metric) error {
	stats, err := c.GetMetaDataServiceInfo()
	if err != nil {
		return err
	}

	values := make([]string, 0, len(c.mapping))
	for i := range c.labels {
		statsKey := c.mapping[c.labels[i]]
		values = append(values, fmt.Sprintf("%v", stats[statsKey]))
	}
	ch <- prometheus.MustNewConstMetric(c.InfoDesc, prometheus.GaugeValue, 1.0, values...)

	return nil
}
