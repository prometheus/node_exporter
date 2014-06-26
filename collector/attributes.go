// +build !noattributes

package collector

import (
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	attributes *prometheus.GaugeVec
)

type attributesCollector struct {
	config Config
}

func init() {
	Factories["attributes"] = NewAttributesCollector
}

// Takes a config struct and prometheus registry and returns a new Collector exposing
// labels from the config.
func NewAttributesCollector(config Config) (Collector, error) {
	c := attributesCollector{
		config: config,
	}
	labelNames := []string{}
	for l := range c.config.Attributes {
		labelNames = append(labelNames, l)
	}
	gv := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "attributes",
			Help:      "The node_exporter attributes.",
		},
		labelNames,
	)
	collector, err := prometheus.RegisterOrGet(gv)
	if err != nil {
		return nil, err
	}
	attributes = collector.(*prometheus.GaugeVec)
	return &c, nil
}

func (c *attributesCollector) Update() (updates int, err error) {
	glog.V(1).Info("Set node_attributes{%v}: 1", c.config.Attributes)
	attributes.Reset()
	attributes.With(c.config.Attributes).Set(1)
	return updates, err
}
