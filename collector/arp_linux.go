// +build !noarp

package collector

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

func getArpEntries() (map[string]uint32, error) {
	file, err := os.Open(procFilePath("net/arp"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parseArpEntries(file), nil
}

func parseArpEntries(data *os.File) map[string]uint32 {
	scanner := bufio.NewScanner(data)
	entries := make(map[string]uint32)

	for scanner.Scan() {
		columns := strings.Split(string(scanner.Text()), " ")

		if columns[0] != "IP" {
			deviceIndex := len(columns) - 1
			entries[columns[deviceIndex]]++
		}
	}

	return entries
}

func (c *arpCollector) Update(ch chan<- prometheus.Metric) error {
	entries, err := getArpEntries()
	if err != nil {
		return fmt.Errorf("could not get ARP entries: %s", err)
	}

	for device, entryCount := range entries {
		ch <- prometheus.MustNewConstMetric(
			c.count, prometheus.GaugeValue, float64(entryCount), device)
	}

	return nil
}
