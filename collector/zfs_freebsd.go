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

func (c *zfsCollector) PrepareUpdate() error {
	if C.zfsModuleLoaded() == 0 {
		return zfsNotAvailableError
	}
	return nil
}

const zfsArcstatsSysctl = "kstat.zfs.misc.arcstats"

func (c *zfsCollector) updateArcstats(ch chan<- prometheus.Metric) (err error) {

	cmd := exec.Command("sysctl", zfsArcstatsSysctl)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	if err = cmd.Start(); err != nil {
		return
	}

	err = c.parseArcstatsSysctlOutput(stdout, func(sysctl zfsSysctl, value zfsMetricValue) {
		ch <- c.ConstSysctlMetric(arc, sysctl, zfsMetricValue(value))
	})
	if err != nil {
		return
	}

	if err = cmd.Wait(); err != nil {
		return
	}

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
