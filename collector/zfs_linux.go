package collector

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const (
	zfsArcstatsProcpath = "spl/kstat/zfs/arcstats"
)

func (c *zfsCollector) zfsAvailable() (err error) {
	file, err := c.openArcstatsFile()
	if err != nil {
		file.Close()
	}
	return err
}

func (c *zfsCollector) openArcstatsFile() (file *os.File, err error) {
	file, err = os.Open(procFilePath(zfsArcstatsProcpath))
	if err != nil {
		log.Debugf("Cannot open '%s' for reading. Is the kernel module loaded?", procFilePath(zfsArcstatsProcpath))
		err = zfsNotAvailableError
	}
	return
}

func (c *zfsCollector) updateArcstats(ch chan<- prometheus.Metric) (err error) {

	file, err := c.openArcstatsFile()
	if err != nil {
		return err
	}
	defer file.Close()

	return c.parseArcstatsProcfsFile(file, func(s zfsSysctl, v zfsMetricValue) {
		ch <- c.ConstSysctlMetric(arc, s, v)
	})

}

func (c *zfsCollector) parseArcstatsProcfsFile(reader io.Reader, handler func(zfsSysctl, zfsMetricValue)) (err error) {

	scanner := bufio.NewScanner(reader)

	parseLine := false
	for scanner.Scan() {

		parts := strings.Fields(scanner.Text())

		if !parseLine && len(parts) == 3 && parts[0] == "name" && parts[1] == "type" && parts[2] == "data" {
			// Start parsing from here.
			parseLine = true
			continue
		}

		if !parseLine || len(parts) < 3 {
			continue
		}

		key := fmt.Sprintf("kstat.zfs.misc.arcstats.%s", parts[0])

		value, err := strconv.Atoi(parts[2])
		if err != nil {
			return fmt.Errorf("could not parse expected integer value for '%s'", key)
		}
		handler(zfsSysctl(key), zfsMetricValue(value))

	}
	if !parseLine {
		return errors.New("did not parse a single arcstat metrics")
	}

	return scanner.Err()
}

func (c *zfsCollector) updatePoolStats(ch chan<- prometheus.Metric) (err error) {
	return nil
}
