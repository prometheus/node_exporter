package collector

import (
	"fmt"
	"unsafe"

	"github.com/prometheus/client_golang/prometheus"
)

/*
#cgo LDFLAGS:
#include <fcntl.h>
#include <stdlib.h>
#include <sys/param.h>
#include <sys/module.h>
#include <sys/pcpu.h>
#include <sys/resource.h>
#include <sys/sysctl.h>
#include <sys/time.h>

int zfsIntegerSysctl(const char *name) {

	int value;
	size_t value_size = sizeof(value);
	if (sysctlbyname(name,  &value, &value_size, NULL, 0) != -1 ||
	    value_size != sizeof(value)) {
		return -1;
	}
	return value;

}

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

func (c *zfsCollector) updateArcstats(ch chan<- prometheus.Metric) (err error) {

	for _, metric := range c.zfsMetrics {

		sysctlCString := C.CString(string(metric.sysctl))
		defer C.free(unsafe.Pointer(sysctlCString))

		value := int(C.zfsIntegerSysctl(sysctlCString))

		if value == -1 {
			return fmt.Errorf("Could not retrieve value for metric '%v'", metric)
		}

		ch <- metric.ConstMetric(zfsMetricValue(value))
	}

	return err
}
