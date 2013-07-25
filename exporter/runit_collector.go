package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/soundcloud/go-runit/runit"
)

var ()

type runitCollector struct {
	name          string
	config        config
	serviceStatus prometheus.Gauge
}

func NewRunitCollector(config config, registry prometheus.Registry) (runitCollector, error) {
	c := runitCollector{
		name:          "runit_collector",
		config:        config,
		serviceStatus: prometheus.NewGauge(),
	}

	registry.Register(
		"node_service_status",
		"node_exporter: status of runit service.",
		prometheus.NilLabels,
		c.serviceStatus,
	)

	return c, nil
}

func (c *runitCollector) Name() string { return c.name }

func (c *runitCollector) Update() (updates int, err error) {
	services, err := runit.GetServices()
	if err != nil {
		return 0, err
	}

	for _, service := range services {
		status, err := service.Status()
		if err != nil {
			return 0, err
		}
		debug(c.Name(), "%s is %d on pid %d for %d seconds", service.Name, status.State, status.Pid, status.Duration)
		labels := map[string]string{
			"name":  service.Name,
			"state": runit.StateToString[status.State],
			"want":  runit.StateToString[status.Want],
		}

		if status.NormallyUp {
			labels["normally_up"] = "yes"
		} else {
			labels["normally_up"] = "no"
		}

		c.serviceStatus.Set(labels, float64(status.Duration))
		updates++
	}
	return updates, err
}
