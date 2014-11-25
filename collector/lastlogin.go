// +build !nolastlogin

package collector

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

const lastLoginSubsystem = "last_login"

type lastLoginCollector struct {
	config Config
	metric prometheus.Gauge
}

func init() {
	Factories["lastlogin"] = NewLastLoginCollector
}

// Takes a config struct and prometheus registry and returns a new Collector exposing
// load, seconds since last login and a list of tags as specified by config.
func NewLastLoginCollector(config Config) (Collector, error) {
	return &lastLoginCollector{
		config: config,
		metric: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: lastLoginSubsystem,
			Name:      "time",
			Help:      "The time of the last login.",
		}),
	}, nil
}

func (c *lastLoginCollector) Update(ch chan<- prometheus.Metric) (err error) {
	last, err := getLastLoginTime()
	if err != nil {
		return fmt.Errorf("Couldn't get last seen: %s", err)
	}
	glog.V(1).Infof("Set node_last_login_time: %f", last)
	c.metric.Set(last)
	c.metric.Collect(ch)
	return err
}

func getLastLoginTime() (float64, error) {
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

	return float64(last.Unix()), nil
}
