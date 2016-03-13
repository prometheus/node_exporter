package collector

import (
	"fmt"
	"unsafe"
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

func (c *zfsMetricProvider) PrepareUpdate() error {
	if C.zfsModuleLoaded() == 0 {
		return zfsNotAvailableError
	}
	return nil
}

func (p *zfsMetricProvider) handleMiss(s zfsSysctl) (zfsMetricValue, error) {

	sysctlCString := C.CString(string(s))
	defer C.free(unsafe.Pointer(sysctlCString))

	value := int(C.zfsIntegerSysctl(sysctlCString))

	if value == -1 {
		return zfsErrorValue, fmt.Errorf("Could not retrieve sysctl '%s'", s)
	}

	return zfsMetricValue(value), nil

}
