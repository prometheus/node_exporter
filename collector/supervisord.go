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

// +build !nosupervisord

package collector

import (
	"flag"

	"github.com/kolo/xmlrpc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

var (
	supervisordURL = flag.String("collector.supervisord.url", "http://localhost:9001/RPC2", "XML RPC endpoint")
)

type supervisordCollector struct {
	client         *xmlrpc.Client
	upDesc         *prometheus.Desc
	stateDesc      *prometheus.Desc
	exitStatusDesc *prometheus.Desc
	uptimeDesc     *prometheus.Desc
}

func init() {
	Factories["supervisord"] = NewSupervisordCollector
}

func NewSupervisordCollector() (Collector, error) {
	client, err := xmlrpc.NewClient(*supervisordURL, nil)
	if err != nil {
		return nil, err
	}

	var (
		subsystem  = "supervisord"
		labelNames = []string{"name", "group"}
	)
	return &supervisordCollector{
		client: client,
		upDesc: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "up"),
			"Process Up",
			labelNames,
			nil,
		),
		stateDesc: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "state"),
			"Process State",
			labelNames,
			nil,
		),
		exitStatusDesc: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "exit_status"),
			"Process Exit Status",
			labelNames,
			nil,
		),
		uptimeDesc: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "uptime"),
			"Process Uptime",
			labelNames,
			nil,
		),
	}, nil
}

func (c *supervisordCollector) isRunning(state int) bool {
	// http://supervisord.org/subprocess.html#process-states
	const (
		STOPPED  = 0
		STARTING = 10
		RUNNING  = 20
		BACKOFF  = 30
		STOPPING = 40
		EXITED   = 100
		FATAL    = 200
		UNKNOWN  = 1000
	)
	switch state {
	case STARTING, RUNNING, STOPPING:
		return true
	}
	return false
}

func (c *supervisordCollector) Update(ch chan<- prometheus.Metric) error {
	var infos []struct {
		Name          string `xmlrpc:"name"`
		Group         string `xmlrpc:"group"`
		Start         int    `xmlrpc:"start"`
		Stop          int    `xmlrpc:"stop"`
		Now           int    `xmlrpc:"now"`
		State         int    `xmlrpc:"state"`
		StateName     string `xmlrpc:"statename"`
		SpawnErr      string `xmlrpc:"spanerr"`
		ExitStatus    int    `xmlrpc:"exitstatus"`
		StdoutLogfile string `xmlrcp:"stdout_logfile"`
		StderrLogfile string `xmlrcp:"stderr_logfile"`
		PID           int    `xmlrpc:"pid"`
	}
	if err := c.client.Call("supervisor.getAllProcessInfo", nil, &infos); err != nil {
		return err
	}
	for _, info := range infos {
		lables := []string{info.Name, info.Group}

		ch <- prometheus.MustNewConstMetric(c.stateDesc, prometheus.GaugeValue, float64(info.State), lables...)
		ch <- prometheus.MustNewConstMetric(c.exitStatusDesc, prometheus.GaugeValue, float64(info.ExitStatus), lables...)

		if c.isRunning(info.State) {
			ch <- prometheus.MustNewConstMetric(c.upDesc, prometheus.GaugeValue, 1, lables...)
			ch <- prometheus.MustNewConstMetric(c.uptimeDesc, prometheus.CounterValue, float64(info.Now-info.Start), lables...)
		} else {
			ch <- prometheus.MustNewConstMetric(c.upDesc, prometheus.GaugeValue, 0, lables...)
			ch <- prometheus.MustNewConstMetric(c.uptimeDesc, prometheus.CounterValue, 0, lables...)
		}
		log.Debugf("%s:%s is %s on pid %d", info.Group, info.Name, info.StateName, info.PID)
	}

	return nil
}
