// +build !nonetStat

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
	procNetStat       = "/proc/net/netstat"
	netStatsSubsystem = "netstat"
)

var (
	netStatsMetrics = map[string]prometheus.Gauge{}
)

type netStatCollector struct {
	config Config
}

func init() {
	Factories["netstat"] = NewNetStatCollector
}

// NewNetStatCollector takes a config struct and returns
// a new Collector exposing network stats.
func NewNetStatCollector(config Config) (Collector, error) {
	c := netStatCollector{
		config: config,
	}
	return &c, nil
}

func (c *netStatCollector) Update(ch chan<- prometheus.Metric) (err error) {
	netStats, err := getNetStats()
	if err != nil {
		return fmt.Errorf("couldn't get netstats: %s", err)
	}
	for protocol, protocolStats := range netStats {
		for name, value := range protocolStats {
			key := protocol + "_" + name
			if _, ok := netStatsMetrics[key]; !ok {
				netStatsMetrics[key] = prometheus.NewGauge(
					prometheus.GaugeOpts{
						Namespace: Namespace,
						Subsystem: netStatsSubsystem,
						Name:      key,
						Help:      fmt.Sprintf("%s %s from /proc/net/netstat.", protocol, name),
					},
				)
			}
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("invalid value %s in netstats: %s", value, err)
			}
			netStatsMetrics[key].Set(v)
		}
	}
	for _, m := range netStatsMetrics {
		m.Collect(ch)
	}
	return err
}

func getNetStats() (map[string]map[string]string, error) {
	file, err := os.Open(procNetStat)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parseNetStats(file)
}

func parseNetStats(r io.Reader) (map[string]map[string]string, error) {
	var (
		netStats = map[string]map[string]string{}
		scanner  = bufio.NewScanner(r)
	)

	for scanner.Scan() {
		nameParts := strings.Split(string(scanner.Text()), " ")
		scanner.Scan()
		valueParts := strings.Split(string(scanner.Text()), " ")
		// Remove trailing :.
		protocol := nameParts[0][:len(nameParts[0])-1]
		netStats[protocol] = map[string]string{}
		if len(nameParts) != len(valueParts) {
			return nil, fmt.Errorf("mismatch field count mismatch in %s: %s",
				procNetStat, protocol)
		}
		for i := 1; i < len(nameParts); i++ {
			netStats[protocol][nameParts[i]] = valueParts[i]
		}
	}

	return netStats, nil
}
