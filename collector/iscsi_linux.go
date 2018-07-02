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

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/sysfs"
)

type iscsiStatsCollector struct {
	targetName    string
	initiatorName string
	state         string
	ifacename     string
	sessionName   string
}

func init() {
	registerCollector("iscsistats", defaultEnabled, NewIscsiStatsCollector)
}

func NewIscsiStatsCollector() (Collector, error) {
	return &iscsiStatsCollector{}, nil
}

func (c *iscsiStatsCollector) Update(ch chan<- prometheus.Metric) error {
	finalIscsiStats, err := sysfs.getIscsiStats()
	if err != nil {
		return fmt.Errorf("couldn't get iscsi stats: %v", err)
	}

	for _, info := range finalIscsiStats {

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "", "iscsi_session_info"),
				fmt.Sprintf("ISCSI session information on the node"),
				nil,
				prometheus.Labels{
					"target_name":    info.targetName,
					"initiator_name": info.initiatorName,
					"session_state":  info.state,
					"interface_name": info.ifacename,
					"session_name":   info.sessionName,
				},
			),
			prometheus.GaugeValue, 1,
		)
	}
	return nil
}
