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

// +build !nosystemd

package collector

import (
	"fmt"
	"strings"

	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	"github.com/prometheus/client_golang/prometheus"
)

type systemdCollector struct {
	metric *prometheus.GaugeVec
}

type systemdService struct {
	Name  string
	State string
}

func init() {
	Factories["systemd"] = NewSystemdCollector
}

// NewSystemdCollector is used to create a new collector for Systemd services
func NewSystemdCollector() (Collector, error) {
	return &systemdCollector{
		metric: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Name:      "systemd_services",
				Help:      "service states.",
			},
			[]string{"service"},
		),
	}, nil
}

func (c *systemdCollector) Update(ch chan<- prometheus.Metric) error {
	infos, err := getServiceData()
	if err != nil {
		return err
	}

	for _, info := range infos {
		c.metric.WithLabelValues(info.Name).Set(toMetric(info.State))
	}

	c.metric.Collect(ch)
	return err
}

func toMetric(state string) float64 {
	if state == "active" {
		return 1.0
	}

	return 0.0
}

func getServiceData() ([]systemdService, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	node, err := introspect.Call(conn.Object("org.freedesktop.systemd1", "/org/freedesktop/systemd1/unit"))
	if err != nil {
		return nil, err
	}

	services := getServices(node)
	return getData(services, conn), nil
}

func getData(services []string, conn *dbus.Conn) []systemdService {
	var data []systemdService

	for _, s := range services {
		d := systemdService{}

		objectPath := dbus.ObjectPath(fmt.Sprintf("/org/freedesktop/systemd1/unit/%s", s))
		obj := conn.Object("org.freedesktop.systemd1", objectPath)
		state := obj.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.systemd1.Unit", "ActiveState")

		currentState := strings.Replace(state.Body[0].(dbus.Variant).String(), "\"", "", 2)
		d.State = currentState

		info := obj.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.systemd1.Unit", "Id")
		id := strings.Replace(info.Body[0].(dbus.Variant).String(), "\"", "", 2)
		d.Name = id

		data = append(data, d)
	}

	return data
}

func getServices(node *introspect.Node) []string {
	var services []string
	for _, n := range node.Children {
		if strings.HasSuffix(n.Name, "2eservice") {
			services = append(services, n.Name)
		}
	}

	return services
}
