// +build !noloadavg

package collector

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

func getLoad() ([]float64, error) {
	type loadavg struct {
		load  [3]uint32
		scale int
	}
	b, err := unix.SysctlRaw("vm.loadavg")
	if err != nil {
		return nil, err
	}
	load := *(*loadavg)(unsafe.Pointer((&b[0])))
	scale := float64(load.scale)
	return []float64{
		float64(load.load[0]) / scale,
		float64(load.load[1]) / scale,
		float64(load.load[2]) / scale,
	}, nil
}
