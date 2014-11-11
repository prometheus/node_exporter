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
	procNetDev      = "/proc/net/dev"
	netDevSubsystem = "network"
)

var (
	netDevMetrics = map[string]*prometheus.GaugeVec{}
)

type netDevCollector struct {
	config Config
}

func init() {
	Factories["netdev"] = NewNetDevCollector
}

// Takes a config struct and prometheus registry and returns a new Collector exposing
// network device stats.
func NewNetDevCollector(config Config) (Collector, error) {
	c := netDevCollector{
		config: config,
	}
	return &c, nil
}

func (c *netDevCollector) Update(ch chan<- prometheus.Metric) (err error) {
	netDev, err := getNetDevStats()
	if err != nil {
		return fmt.Errorf("Couldn't get netstats: %s", err)
	}
	for direction, devStats := range netDev {
		for dev, stats := range devStats {
			for t, value := range stats {
				key := direction + "_" + t
				if _, ok := netDevMetrics[key]; !ok {
					netDevMetrics[key] = prometheus.NewGaugeVec(
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
				netDevMetrics[key].WithLabelValues(dev).Set(v)
			}
		}
	}
	for _, m := range netDevMetrics {
		m.Collect(ch)
	}
	return err
}

func getNetDevStats() (map[string]map[string]map[string]string, error) {
	file, err := os.Open(procNetDev)
	if err != nil {
		return nil, err
	}
	return parseNetDevStats(file)
}

func parseNetDevStats(r io.ReadCloser) (map[string]map[string]map[string]string, error) {
	defer r.Close()
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
