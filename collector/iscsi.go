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

// +build !nofilefd

package collector

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type iscsiStatsCollector struct {
	targetName    string
	initiatorName string
	state         string
	ifacename     string
	sessionName   string
}

const (
	iscsiSessionNameRegex = "/sys/class/iscsi_session/session*"
)

func init() {
	registerCollector("iscsistats", defaultEnabled, NewIscsiStatsCollector)
}

func NewIscsiStatsCollector() (Collector, error) {
	return &iscsiStatsCollector{}, nil
}

func (c *iscsiStatsCollector) Update(ch chan<- prometheus.Metric) error {
	finalIscsiStats, err := getIscsiStats()
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

func getIscsiStats() (map[string]iscsiStatsCollector, error) {
	iscsiSessionDirNames, err := filepath.Glob(iscsiSessionNameRegex)
	if err != nil {
		return nil, err
	}
	var iscsiStats = make(map[string]iscsiStatsCollector)

	for _, iscsiSessionDir := range iscsiSessionDirNames {
		thisSession := filepath.Base(iscsiSessionDir)
		iscsiStat := iscsiStatsCollector{}
		iscsiStat.sessionName = thisSession

		targetnameFilepath := fmt.Sprintf("%v/targetname", iscsiSessionDir)
		targetName, err := ioutil.ReadFile(targetnameFilepath)
		if err != nil {
			return nil, err
		}
		iscsiStat.targetName = strings.TrimSuffix(string(targetName), "\n")

		initiatornameFilepath := fmt.Sprintf("%v/initiatorname", iscsiSessionDir)
		initiatorname, err := ioutil.ReadFile(initiatornameFilepath)
		if err != nil {
			return nil, err
		}
		iscsiStat.initiatorName = strings.TrimSuffix(string(initiatorname), "\n")

		stateFilepath := fmt.Sprintf("%v/state", iscsiSessionDir)
		state, err := ioutil.ReadFile(stateFilepath)
		if err != nil {
			return nil, err
		}
		iscsiStat.state = strings.TrimSuffix(string(state), "\n")

		ifacenameFilepath := fmt.Sprintf("%v/ifacename", iscsiSessionDir)
		ifacename, err := ioutil.ReadFile(ifacenameFilepath)
		if err != nil {
			return nil, err
		}
		iscsiStat.ifacename = strings.TrimSuffix(string(ifacename), "\n")

		iscsiStats[thisSession] = iscsiStat
	}

	return iscsiStats, nil
}
