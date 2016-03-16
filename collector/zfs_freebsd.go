package collector

import (
	"bufio"
	"io"
	"os/exec"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

/*
#cgo LDFLAGS:
#include <sys/param.h>
#include <sys/module.h>

int zfsModuleLoaded() {
	int modid = modfind("zfs");
	return modid < 0 ? 0 : -1;
}

*/
import "C"

func (c *zfsCollector) zfsAvailable() error {
	if C.zfsModuleLoaded() == 0 {
		return zfsNotAvailableError
	}
	return nil
}

const zfsArcstatsSysctl = "kstat.zfs.misc.arcstats"

func (c *zfsCollector) RunOnStdout(cmd *exec.Cmd, handler func(io.Reader) error) (err error) {

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	if err = cmd.Start(); err != nil {
		return
	}

	err = handler(stdout)
	if err != nil {
		return
	}

	if err = cmd.Wait(); err != nil {
		return
	}

	return err

}

func (c *zfsCollector) updateArcstats(ch chan<- prometheus.Metric) (err error) {

	cmd := exec.Command("sysctl", zfsArcstatsSysctl)

	err = c.RunOnStdout(cmd, func(stdout io.Reader) error {
		return c.parseArcstatsSysctlOutput(stdout, func(sysctl zfsSysctl, value zfsMetricValue) {
			ch <- c.ConstSysctlMetric(arc, sysctl, zfsMetricValue(value))
		})
	})
	return err

}

func (c *zfsCollector) parseArcstatsSysctlOutput(reader io.Reader, handler func(zfsSysctl, zfsMetricValue)) (err error) {

	// Decode values
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {

		fields := strings.Fields(scanner.Text())

		if len(fields) != 2 ||
			!strings.HasPrefix(fields[0], zfsArcstatsSysctl) ||
			!strings.HasSuffix(fields[0], ":") {

			log.Debugf("Skipping line of unknown format: %s", scanner.Text())
			continue

		}

		sysctl := zfsSysctl(strings.TrimSuffix(fields[0], ":"))
		value, err := strconv.Atoi(fields[1])
		if err != nil {
			return err
		}

		handler(sysctl, zfsMetricValue(value))
	}
	return scanner.Err()

}

func (c *zfsCollector) updatePoolStats(ch chan<- prometheus.Metric) (err error) {

	poolProperties := []string{"size", "free", "allocated", "capacity", "dedupratio", "fragmentation"}

	cmd := exec.Command("zpool", "get", "-pH", strings.Join(poolProperties, ","))

	err = c.RunOnStdout(cmd, func(stdout io.Reader) error {
		return c.parseZpoolOutput(stdout, func(pool, name string, value float64) {
			ch <- c.ConstZpoolMetric(pool, name, value)
		})
	})

	return err
}
