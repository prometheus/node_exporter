// +build !nonetDev

package collector

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	procNetDev = "/proc/net/dev"
)

var (
	netStatsMetrics = map[string]prometheus.Gauge{}
)

type netDevCollector struct {
	registry prometheus.Registry
	config   Config
}

func init() {
	Factories["netdev"] = NewNetDevCollector
}

// Takes a config struct and prometheus registry and returns a new Collector exposing
// network device stats.
func NewNetDevCollector(config Config, registry prometheus.Registry) (Collector, error) {
	c := netDevCollector{
		config:   config,
		registry: registry,
	}
	return &c, nil
}

func (c *netDevCollector) Update() (updates int, err error) {
	netStats, err := getNetStats()
	if err != nil {
		return updates, fmt.Errorf("Couldn't get netstats: %s", err)
	}
	for direction, devStats := range netStats {
		for dev, stats := range devStats {
			for t, value := range stats {
				key := direction + "_" + t
				if _, ok := netStatsMetrics[key]; !ok {
					netStatsMetrics[key] = prometheus.NewGauge()
					c.registry.Register(
						"node_network_"+key,
						t+" "+direction+" from /proc/net/dev",
						prometheus.NilLabels,
						netStatsMetrics[key],
					)
				}
				updates++
				v, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return updates, fmt.Errorf("Invalid value %s in netstats: %s", value, err)
				}
				netStatsMetrics[key].Set(map[string]string{"device": dev}, v)
			}
		}
	}
	return updates, err
}

func getNetStats() (map[string]map[string]map[string]string, error) {
	file, err := os.Open(procNetDev)
	if err != nil {
		return nil, err
	}
	return parseNetStats(file)
}

func parseNetStats(r io.ReadCloser) (map[string]map[string]map[string]string, error) {
	defer r.Close()
	netStats := map[string]map[string]map[string]string{}
	netStats["transmit"] = map[string]map[string]string{}
	netStats["receive"] = map[string]map[string]string{}

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
		parts := strings.Fields(string(scanner.Text()))
		if len(parts) != 2*len(header)+1 {
			return nil, fmt.Errorf("Invalid line in %s: %s",
				procNetDev, scanner.Text())
		}

		dev := parts[0][:len(parts[0])-1]
		receive, err := parseNetDevLine(parts[1:len(header)+1], header)
		if err != nil {
			return nil, err
		}

		transmit, err := parseNetDevLine(parts[len(header)+1:], header)
		if err != nil {
			return nil, err
		}
		netStats["transmit"][dev] = transmit
		netStats["receive"][dev] = receive
	}
	return netStats, nil
}

func parseNetDevLine(parts []string, header []string) (map[string]string, error) {
	devStats := map[string]string{}
	for i, v := range parts {
		devStats[header[i]] = v
	}
	return devStats, nil
}
