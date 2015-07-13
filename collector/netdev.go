// +build !nonetdev

package collector

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/log"
)

const (
	procNetDev      = "/proc/net/dev"
	netDevSubsystem = "network"
)

var (
	fieldSep             = regexp.MustCompile("[ :] *")
	netdevIgnoredDevices = flag.String("collector.netdev.ignored-devices", "^$", "Regexp of net devices to ignore for netdev collector.")
)

type netDevCollector struct {
	ignoredDevicesPattern *regexp.Regexp

	metrics map[string]*prometheus.GaugeVec
}

func init() {
	Factories["netdev"] = NewNetDevCollector
}

// Takes a prometheus registry and returns a new Collector exposing
// network device stats.
func NewNetDevCollector() (Collector, error) {
	return &netDevCollector{
		ignoredDevicesPattern: regexp.MustCompile(*netdevIgnoredDevices),
		metrics:               map[string]*prometheus.GaugeVec{},
	}, nil
}

func (c *netDevCollector) Update(ch chan<- prometheus.Metric) (err error) {
	netDev, err := getNetDevStats(c.ignoredDevicesPattern)
	if err != nil {
		return fmt.Errorf("Couldn't get netstats: %s", err)
	}
	for direction, devStats := range netDev {
		for dev, stats := range devStats {
			for t, value := range stats {
				key := direction + "_" + t
				if _, ok := c.metrics[key]; !ok {
					c.metrics[key] = prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Namespace: Namespace,
							Subsystem: netDevSubsystem,
							Name:      key,
							Help:      fmt.Sprintf("%s %s from /proc/net/dev.", t, direction),
						},
						[]string{"device"},
					)
				}
				v, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return fmt.Errorf("Invalid value %s in netstats: %s", value, err)
				}
				c.metrics[key].WithLabelValues(dev).Set(v)
			}
		}
	}
	for _, m := range c.metrics {
		m.Collect(ch)
	}
	return err
}

func getNetDevStats(ignore *regexp.Regexp) (map[string]map[string]map[string]string, error) {
	file, err := os.Open(procNetDev)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parseNetDevStats(file, ignore)
}

func parseNetDevStats(r io.Reader, ignore *regexp.Regexp) (map[string]map[string]map[string]string, error) {
	netDev := map[string]map[string]map[string]string{}
	netDev["transmit"] = map[string]map[string]string{}
	netDev["receive"] = map[string]map[string]string{}

	scanner := bufio.NewScanner(r)
	scanner.Scan() // skip first header
	scanner.Scan()
	parts := strings.Split(string(scanner.Text()), "|")
	if len(parts) != 3 { // interface + receive + transmit
		return nil, fmt.Errorf("Invalid header line in %s: %s",
			procNetDev, scanner.Text())
	}
	header := strings.Fields(parts[1])
	for scanner.Scan() {
		line := strings.TrimLeft(string(scanner.Text()), " ")
		parts := fieldSep.Split(line, -1)
		if len(parts) != 2*len(header)+1 {
			return nil, fmt.Errorf("Invalid line in %s: %s",
				procNetDev, scanner.Text())
		}

		dev := parts[0][:len(parts[0])]
		receive, err := parseNetDevLine(parts[1:len(header)+1], header)
		if err != nil {
			return nil, err
		}

		transmit, err := parseNetDevLine(parts[len(header)+1:], header)
		if err != nil {
			return nil, err
		}

		if ignore.MatchString(dev) {
			log.Debugf("Ignoring device: %s", dev)
			continue
		}
		netDev["transmit"][dev] = transmit
		netDev["receive"][dev] = receive
	}
	return netDev, nil
}

func parseNetDevLine(parts []string, header []string) (map[string]string, error) {
	devStats := map[string]string{}
	for i, v := range parts {
		devStats[header[i]] = v
	}
	return devStats, nil
}
