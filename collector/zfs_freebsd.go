package collector

// +build freebsd

import (
	"fmt"
	"unsafe"
)

/*
#cgo LDFLAGS:
#include <fcntl.h>
#include <stdlib.h>
#include <sys/param.h>
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

*/
import "C"

func zfsInitialize() error {
	return nil
}

func (c *zfsMetricProvider) PrepareUpdate() error {
	return nil
}

func (p *zfsMetricProvider) handleFetchedMetricCacheMiss(c zfsSysctl) (zfsMetricValue, error) {

	sysctlCString := C.CString(string(c))
	defer C.free(unsafe.Pointer(sysctlCString))

	value := int(C.zfsIntegerSysctl(sysctlCString))

	if value == -1 {
		return zfsErrorValue, fmt.Errorf("Could not retrieve sysctl '%s'", c)
	}

	return zfsMetricValue(value), nil

}
