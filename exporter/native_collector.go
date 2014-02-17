// +build !nonative

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
	procLoad    = "/proc/loadavg"
	procMemInfo = "/proc/meminfo"
)

type nativeCollector struct {
	loadAvg    prometheus.Gauge
	attributes prometheus.Gauge
	lastSeen   prometheus.Gauge
	memInfo    prometheus.Gauge
	name       string
	config     config
}

func init() {
	collectorFactories = append(collectorFactories, NewNativeCollector)
}

// Takes a config struct and prometheus registry and returns a new Collector exposing
// load, seconds since last login and a list of tags as specified by config.
func NewNativeCollector(config config, registry prometheus.Registry) (Collector, error) {
	c := nativeCollector{
		name:       "native_collector",
		config:     config,
		loadAvg:    prometheus.NewGauge(),
		attributes: prometheus.NewGauge(),
		lastSeen:   prometheus.NewGauge(),
		memInfo:    prometheus.NewGauge(),
	}

	registry.Register(
		"node_load",
		"node_exporter: system load.",
		prometheus.NilLabels,
		c.loadAvg,
	)

	registry.Register(
		"node_last_login_seconds",
		"node_exporter: seconds since last login.",
		prometheus.NilLabels,
		c.lastSeen,
	)

	registry.Register(
		"node_attributes",
		"node_exporter: system attributes.",
		prometheus.NilLabels,
		c.attributes,
	)

	registry.Register(
		"node_mem",
		"node_exporter: memory details.",
		prometheus.NilLabels,
		c.memInfo,
	)

	return &c, nil
}

func (c *nativeCollector) Name() string { return c.name }

func (c *nativeCollector) Update() (updates int, err error) {
	last, err := getSecondsSinceLastLogin()
	if err != nil {
		return updates, fmt.Errorf("Couldn't get last seen: %s", err)
	}
	updates++
	debug(c.Name(), "Set node_last_login_seconds: %f", last)
	c.lastSeen.Set(nil, last)

	load, err := getLoad()
	if err != nil {
		return updates, fmt.Errorf("Couldn't get load: %s", err)
	}
	updates++
	debug(c.Name(), "Set node_load: %f", load)
	c.loadAvg.Set(nil, load)

	debug(c.Name(), "Set node_attributes{%v}: 1", c.config.Attributes)
	c.attributes.Set(c.config.Attributes, 1)

	memInfo, err := getMemInfo()
	if err != nil {
		return updates, fmt.Errorf("Couldn't get meminfo: %s", err)
	}
	debug(c.Name(), "Set node_mem: %#v", memInfo)
	for k, v := range memInfo {
		updates++
		fv, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return updates, fmt.Errorf("Invalid value in meminfo: %s", err)
		}
		c.memInfo.Set(map[string]string{"type": k}, fv)
	}
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

func getMemInfo() (map[string]string, error) {
	memInfo := map[string]string{}
	fh, err := os.Open(procMemInfo)
	if err != nil {
		return nil, err
	}
	defer fh.Close()
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(string(line))
		key := ""
		switch len(parts) {
		case 2: // no unit
			key = parts[0][:len(parts[0])-1] // remove trailing : from key
		case 3: // has unit
			key = fmt.Sprintf("%s_%s", parts[0][:len(parts[0])-1], parts[2])
		default:
			return nil, fmt.Errorf("Invalid line in %s: %s", procMemInfo, line)
		}
		memInfo[key] = parts[1]
	}
	return memInfo, nil

}
