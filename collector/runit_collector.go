// +build runit

package collector

import (
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/soundcloud/go-runit/runit"
)

type runitCollector struct {
	config       Config
	state        prometheus.Gauge
	stateDesired prometheus.Gauge
	stateNormal  prometheus.Gauge
}

func init() {
	Factories["runit"] = NewRunitCollector
}

func NewRunitCollector(config Config, registry prometheus.Registry) (Collector, error) {
	c := runitCollector{
		config:       config,
		state:        prometheus.NewGauge(),
		stateDesired: prometheus.NewGauge(),
		stateNormal:  prometheus.NewGauge(),
	}

	registry.Register(
		"node_service_state",
		"node_exporter: state of runit service.",
		prometheus.NilLabels,
		c.state,
	)

	registry.Register(
		"node_service_desired_state",
		"node_exporter: desired state of runit service.",
		prometheus.NilLabels,
		c.stateDesired,
	)

	registry.Register(
		"node_service_normal_state",
		"node_exporter: normal state of runit service.",
		prometheus.NilLabels,
		c.stateNormal,
	)

	return &c, nil
}

func (c *runitCollector) Update() (updates int, err error) {
	services, err := runit.GetServices("/etc/service")
	if err != nil {
		return 0, err
	}

	for _, service := range services {
		status, err := service.Status()
		if err != nil {
			glog.V(1).Infof("Couldn't get status for %s: %s, skipping...", service.Name, err)
			continue
		}

		glog.V(1).Infof("%s is %d on pid %d for %d seconds", service.Name, status.State, status.Pid, status.Duration)
		labels := map[string]string{
			"service": service.Name,
		}

		c.state.Set(labels, float64(status.State))
		c.stateDesired.Set(labels, float64(status.Want))
		if status.NormallyUp {
			c.stateNormal.Set(labels, 1)
		} else {
			c.stateNormal.Set(labels, 1)
		}
		updates += 3
	}

	return updates, err
}
