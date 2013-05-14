package exporter

import (
	"bufio"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	procLoad = "/proc/loadavg"
)

type nativeCollector struct {
	loadAvg    prometheus.Gauge
	attributes prometheus.Gauge
	lastSeen   prometheus.Gauge
	name       string
	config     config
}

// Takes a config struct and prometheus registry and returns a new Collector exposing
// load, seconds since last login and a list of tags as specified by config.
func NewNativeCollector(config config, registry prometheus.Registry) (collector nativeCollector, err error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nativeCollector{}, fmt.Errorf("Couldn't get hostname: %s", err)
	}

	collector = nativeCollector{
		name:       "native_collector",
		config:     config,
		loadAvg:    prometheus.NewGauge(),
		attributes: prometheus.NewGauge(),
		lastSeen:   prometheus.NewGauge(),
	}

	registry.Register(
		"node_load",
		"node_exporter: system load.",
		map[string]string{"hostname": hostname},
		collector.loadAvg,
	)

	registry.Register(
		"node_last_login_seconds",
		"node_exporter: seconds since last login.",
		map[string]string{"hostname": hostname},
		collector.lastSeen,
	)

	registry.Register(
		"node_attributes",
		"node_exporter: system attributes.",
		map[string]string{"hostname": hostname},
		collector.attributes,
	)

	return collector, nil
}

func (c *nativeCollector) Name() string { return c.name }

func (c *nativeCollector) Update() (updates int, err error) {
	last, err := getSecondsSinceLastLogin()
	if err != nil {
		return updates, fmt.Errorf("Couldn't get last seen: %s", err)
	} else {
		updates++
		debug(c.Name(), "Set node_last_login_seconds: %f", last)
		c.lastSeen.Set(nil, last)
	}

	load, err := getLoad()
	if err != nil {
		return updates, fmt.Errorf("Couldn't get load: %s", err)
	} else {
		updates++
		debug(c.Name(), "Set node_load: %f", load)
		c.loadAvg.Set(nil, load)
	}
	debug(c.Name(), "Set node_attributes{%v}: 1", c.config.Attributes)
	c.attributes.Set(c.config.Attributes, 1)
	return updates, err
}

func getLoad() (float64, error) {
	data, err := ioutil.ReadFile(procLoad)
	if err != nil {
		return 0, err
	}
	parts := strings.Fields(string(data))
	load, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, fmt.Errorf("Could not parse load '%s': %s", parts[0], err)
	}
	return load, nil
}

func getSecondsSinceLastLogin() (float64, error) {
	who := exec.Command("who", "/var/log/wtmp", "-l", "-u", "-s")

	output, err := who.StdoutPipe()
	if err != nil {
		return 0, err
	}

	err = who.Start()
	if err != nil {
		return 0, err
	}

	reader := bufio.NewReader(output)

	var last time.Time
	for {
		line, isPrefix, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if isPrefix {
			return 0, fmt.Errorf("line to long: %s(...)", line)
		}

		fields := strings.Fields(string(line))
		lastDate := fields[2]
		lastTime := fields[3]

		dateParts, err := splitToInts(lastDate, "-") // 2013-04-16
		if err != nil {
			return 0, fmt.Errorf("Couldn't parse date in line '%s': %s", fields, err)
		}

		timeParts, err := splitToInts(lastTime, ":") // 11:33
		if err != nil {
			return 0, fmt.Errorf("Couldn't parse time in line '%s': %s", fields, err)
		}

		last_t := time.Date(dateParts[0], time.Month(dateParts[1]), dateParts[2], timeParts[0], timeParts[1], 0, 0, time.UTC)
		last = last_t
	}
	err = who.Wait()
	if err != nil {
		return 0, err
	}

	return float64(time.Now().Sub(last).Seconds()), nil
}
