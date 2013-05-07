package exporter

import (
	"bufio"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"io"
	"net"
	"regexp"
	"strconv"
	"strings"
)

const (
	muninAddress = "127.0.0.1:4949"
	muninProto   = "tcp"
)

var muninBanner = regexp.MustCompile(`# munin node at (.*)`)

type muninCollector struct {
	name           string
	hostname       string
	graphs         []string
	gaugePerMetric map[string]prometheus.Gauge
	config         config
	registry       prometheus.Registry
	connection     net.Conn
}

// Takes a config struct and prometheus registry and returns a new Collector scraping munin.
func NewMuninCollector(config config, registry prometheus.Registry) (c muninCollector, err error) {
	c = muninCollector{
		name:           "munin_collector",
		config:         config,
		registry:       registry,
		gaugePerMetric: make(map[string]prometheus.Gauge),
	}

	return c, err
}

func (c *muninCollector) Name() string { return c.name }

func (c *muninCollector) connect() (err error) {
	c.connection, err = net.Dial(muninProto, muninAddress)
	if err != nil {
		return err
	}
	debug(c.Name(), "Connected.")

	reader := bufio.NewReader(c.connection)
	head, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	matches := muninBanner.FindStringSubmatch(head)
	if len(matches) != 2 { // expect: # munin node at <hostname>
		return fmt.Errorf("Unexpected line: %s", head)
	}
	c.hostname = matches[1]
	debug(c.Name(), "Found hostname: %s", c.hostname)
	return err
}

func (c *muninCollector) muninCommand(cmd string) (reader *bufio.Reader, err error) {
	if c.connection == nil {
		err := c.connect()
		if err != nil {
			return reader, fmt.Errorf("Couldn't connect to munin: %s", err)
		}
	}
	reader = bufio.NewReader(c.connection)

	fmt.Fprintf(c.connection, cmd+"\n")

	_, err = reader.Peek(1)
	switch err {
	case io.EOF:
		debug(c.Name(), "not connected anymore, closing connection and reconnect.")
		c.connection.Close()
		err = c.connect()
		if err != nil {
			return reader, fmt.Errorf("Couldn't connect to %s: %s", muninAddress)
		}
		return c.muninCommand(cmd)
	case nil: //no error
		break
	default:
		return reader, fmt.Errorf("Unexpected error: %s", err)
	}

	return reader, err
}

func (c *muninCollector) muninList() (items []string, err error) {
	munin, err := c.muninCommand("list")
	if err != nil {
		return items, fmt.Errorf("Couldn't get list: %s", err)
	}

	response, err := munin.ReadString('\n') // we are only interested in the first line
	if err != nil {
		return items, fmt.Errorf("Couldn't read response: %s", err)
	}

	if response[0] == '#' { // # not expected here
		return items, fmt.Errorf("Error getting items: %s", response)
	}
	items = strings.Fields(strings.TrimRight(response, "\n"))
	return items, err
}

func (c *muninCollector) muninConfig(name string) (config map[string]map[string]string, graphConfig map[string]string, err error) {
	graphConfig = make(map[string]string)
	config = make(map[string]map[string]string)

	resp, err := c.muninCommand("config " + name)
	if err != nil {
		return config, graphConfig, fmt.Errorf("Couldn't get config for %s: %s", name, err)
	}

	for {
		line, err := resp.ReadString('\n')
		if err == io.EOF {
			debug(c.Name(), "EOF, retrying")
			return c.muninConfig(name)
		}
		if err != nil {
			return nil, nil, err
		}
		if line == ".\n" { // munin end marker
			break
		}
		if line[0] == '#' { // here it's just a comment, so ignore it
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			return nil, nil, fmt.Errorf("Line unexpected: %s", line)
		}
		key, value := parts[0], strings.TrimRight(strings.Join(parts[1:], " "), "\n")

		key_parts := strings.Split(key, ".")
		if len(key_parts) > 1 { // it's a metric config (metric.label etc)
			if _, ok := config[key_parts[0]]; !ok {
				config[key_parts[0]] = make(map[string]string)
			}
			config[key_parts[0]][key_parts[1]] = value
		} else {
			graphConfig[key_parts[0]] = value
		}
	}
	return config, graphConfig, err
}

func (c *muninCollector) registerMetrics() (err error) {
	items, err := c.muninList()
	if err != nil {
		return fmt.Errorf("Couldn't get graph list: %s", err)
	}

	for _, name := range items {
		c.graphs = append(c.graphs, name)
		configs, graphConfig, err := c.muninConfig(name)
		if err != nil {
			return fmt.Errorf("Couldn't get config for graph %s: %s", name, err)
		}

		for metric, config := range configs {
			metricName := name + "-" + metric
			desc := graphConfig["graph_title"] + ": " + config["label"]
			if config["info"] != "" {
				desc = desc + ", " + config["info"]
			}
			gauge := prometheus.NewGauge()
			debug(c.Name(), "Register %s: %s", metricName, desc)
			c.gaugePerMetric[metricName] = gauge
			c.registry.Register(metricName, desc, prometheus.NilLabels, gauge)
		}
	}
	return err
}

func (c *muninCollector) Update() (updates int, err error) {
	err = c.registerMetrics()
	if err != nil {
		return updates, fmt.Errorf("Couldn't register metrics: %s", err)
	}

	for _, graph := range c.graphs {
		munin, err := c.muninCommand("fetch " + graph)
		if err != nil {
			return updates, err
		}

		for {
			line, err := munin.ReadString('\n')
			line = strings.TrimRight(line, "\n")
			if err == io.EOF {
				debug(c.Name(), "unexpected EOF, retrying")
				return c.Update()
			}
			if err != nil {
				return updates, err
			}
			if len(line) == 1 && line[0] == '.' {
				break // end of list
			}

			parts := strings.Fields(line)
			if len(parts) != 2 {
				debug(c.Name(), "unexpected line: %s", line)
				continue
			}
			key, value_s := strings.Split(parts[0], ".")[0], parts[1]
			value, err := strconv.ParseFloat(value_s, 64)
			if err != nil {
				debug(c.Name(), "Couldn't parse value in line %s, malformed?", line)
				continue
			}
			labels := map[string]string{
				"hostname": c.hostname,
			}
			name := graph + "-" + key
			debug(c.Name(), "Set %s{%s}: %f\n", name, labels, value)
			c.gaugePerMetric[name].Set(labels, value)
			updates++
		}
	}
	return updates, err
}
