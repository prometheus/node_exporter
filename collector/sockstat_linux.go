// +build !nonetstat

package collector

import (
	"bufio"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	procSockStat      = "/proc/net/sockstat"
	sockStatSubsystem = "sockstat"
)

// Used for calculating the total memory bytes on TCP and UDP
var pageSize = os.Getpagesize()

type sockStatCollector struct {
	metrics map[string]prometheus.Gauge
}

func init() {
	Factories["sockstat"] = NewSockStatCollector
}

// NewSockStatCollector returns a new Collector exposing socket stats
func NewSockStatCollector() (Collector, error) {
	return &sockStatCollector{
		metrics: map[string]prometheus.Gauge{},
	}, nil
}

func (c *sockStatCollector) Update(ch chan<- prometheus.Metric) (err error) {
	sockStats, err := getSockStats(procSockStat)
	if err != nil {
		return fmt.Errorf("couldn't get sockstats: %s", err)
	}
	for protocol, protocolStats := range sockStats {
		for name, value := range protocolStats {
			key := protocol + "_" + name
			if _, ok := c.metrics[key]; !ok {
				c.metrics[key] = prometheus.NewGauge(
					prometheus.GaugeOpts{
						Namespace: Namespace,
						Subsystem: sockStatSubsystem,
						Name:      key,
						Help:      fmt.Sprintf("%s %s from /proc/net/sockstat.", protocol, name),
					},
				)
			}
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("invalid value %s in sockstats: %s", value, err)
			}
			c.metrics[key].Set(v)
		}
	}
	for _, m := range c.metrics {
		m.Collect(ch)
	}
	return err
}

func getSockStats(fileName string) (map[string]map[string]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parseSockStats(file, fileName)
}

func parseSockStats(r io.Reader, fileName string) (map[string]map[string]string, error) {
	var (
		sockStat = map[string]map[string]string{}
		scanner  = bufio.NewScanner(r)
	)

	for scanner.Scan() {
		line := strings.Split(string(scanner.Text()), " ")
		// Remove trailing :
		protocol := line[0][:len(line[0])-1]
		sockStat[protocol] = map[string]string{}

		for i := 1; i < len(line) && i+1 < len(line); i++ {
			sockStat[protocol][line[i]] = line[i+1]
			i++
		}
	}

	// The mem options are reported in pages
	// Multiply them by the pagesize to get bytes
	// Update TCP Mem
	pageCount, err := strconv.Atoi(sockStat["TCP"]["mem"])
	if err != nil {
		return nil, fmt.Errorf("invalid value %s in sockstats: %s", sockStat["TCP"]["mem"], err)
	}
	sockStat["TCP"]["mem"] = strconv.Itoa(pageCount * pageSize)

	// Update UDP Mem
	pageCount, err = strconv.Atoi(sockStat["UDP"]["mem"])
	if err != nil {
		return nil, fmt.Errorf("invalid value %s in sockstats: %s", sockStat["UDP"]["mem"], err)
	}
	sockStat["UDP"]["mem"] = strconv.Itoa(pageCount * pageSize)

	return sockStat, nil
}
