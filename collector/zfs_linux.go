package collector

// +build linux

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/prometheus/common/log"
)

const zfsArcstatsProcpath = "spl/kstat/zfs/arcstats"

func zfsInitialize() error {
	return nil
}

func (p *zfsMetricProvider) PrepareUpdate() (err error) {

	err = p.prepareUpdateArcstats(zfsArcstatsProcpath)
	if err != nil {
		return
	}
	return nil
}

func (p *zfsMetricProvider) handleMiss(s zfsSysctl) (value zfsMetricValue, err error) {
	// all values are fetched in PrepareUpdate().
	return zfsErrorValue, fmt.Errorf("sysctl '%s' found")
}

func (p *zfsMetricProvider) prepareUpdateArcstats(zfsArcstatsProcpath string) (err error) {

	file, err := os.Open(procFilePath(zfsArcstatsProcpath))
	if err != nil {
		log.Debugf("Cannot open ZFS arcstats procfs file for reading. " +
			" Is the kernel module loaded?")
		return zfsNotAvailableError
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	parseLine := false
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())

		if !parseLine && parts[0] == "name" && parts[1] == "type" && parts[2] == "data" {
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
		log.Debugf("%s = %d", key, value)
		p.values[zfsSysctl(key)] = zfsMetricValue(value)
	}
	if !parseLine {
		return errors.New("did not parse a single arcstat metrics")
	}

	return scanner.Err()
}
