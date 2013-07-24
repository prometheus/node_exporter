package exporter

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/node_exporter/exporter/ganglia"
	"io"
	"net"
	"strings"
	"time"
)

const (
	gangliaAddress = "127.0.0.1:8649"
	gangliaProto   = "tcp"
	gangliaTimeout = 30 * time.Second
)

type gmondCollector struct {
	name     string
	Metrics  map[string]prometheus.Gauge
	config   config
	registry prometheus.Registry
}

// Takes a config struct and prometheus registry and returns a new Collector scraping ganglia.
func NewGmondCollector(config config, registry prometheus.Registry) (collector gmondCollector, err error) {
	collector = gmondCollector{
		name:     "gmond_collector",
		config:   config,
		Metrics:  make(map[string]prometheus.Gauge),
		registry: registry,
	}

	return collector, nil
}

func (c *gmondCollector) Name() string { return c.name }

func (c *gmondCollector) setMetric(name string, labels map[string]string, metric ganglia.Metric) {
	if _, ok := c.Metrics[name]; !ok {
		var desc string
		var title string
		for _, element := range metric.ExtraData.ExtraElements {
			switch element.Name {
			case "DESC":
				desc = element.Val
			case "TITLE":
				title = element.Val
			}
			if title != "" && desc != "" {
				break
			}
		}
		debug(c.Name(), "Register %s: %s", name, desc)
		gauge := prometheus.NewGauge()
		c.Metrics[name] = gauge
		c.registry.Register(name, desc, prometheus.NilLabels, gauge) // one gauge per metric!
	}
	debug(c.Name(), "Set %s{%s}: %f", name, labels, metric.Value)
	c.Metrics[name].Set(labels, metric.Value)
}

func (c *gmondCollector) Update() (updates int, err error) {
	conn, err := net.Dial(gangliaProto, gangliaAddress)
	debug(c.Name(), "gmondCollector Update")
	if err != nil {
		return updates, fmt.Errorf("Can't connect to gmond: %s", err)
	}
	conn.SetDeadline(time.Now().Add(gangliaTimeout))

	ganglia := ganglia.Ganglia{}
	decoder := xml.NewDecoder(bufio.NewReader(conn))
	decoder.CharsetReader = toUtf8

	err = decoder.Decode(&ganglia)
	if err != nil {
		return updates, fmt.Errorf("Couldn't parse xml: %s", err)
	}

	for _, cluster := range ganglia.Clusters {
		for _, host := range cluster.Hosts {

			for _, metric := range host.Metrics {
				name := strings.Replace(strings.ToLower(metric.Name), ".", "_", -1)

				var labels = map[string]string{
					"cluster": cluster.Name,
				}
				c.setMetric(name, labels, metric)
				updates++
			}
		}
	}
	return updates, err
}

func toUtf8(charset string, input io.Reader) (io.Reader, error) {
	return input, nil //FIXME
}
