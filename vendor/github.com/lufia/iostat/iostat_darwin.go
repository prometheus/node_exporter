// +build darwin

package iostat

// #cgo LDFLAGS: -framework CoreFoundation -framework IOKit
// #include <stdint.h>
// #include "iostat_darwin.h"
import "C"
import (
	"time"
)

// ReadDriveStats returns statictics of each of the drives.
func ReadDriveStats() ([]*DriveStats, error) {
	var buf [C.NDRIVE]C.DriveStats
	n, err := C.readdrivestat(&buf[0], C.int(len(buf)))
	if err != nil {
		return nil, err
	}
	stats := make([]*DriveStats, n)
	for i := 0; i < int(n); i++ {
		stats[i] = &DriveStats{
			Name:           C.GoString(&buf[i].name[0]),
			Size:           int64(buf[i].size),
			BlockSize:      int64(buf[i].blocksize),
			BytesRead:      int64(buf[i].read),
			BytesWritten:   int64(buf[i].written),
			NumRead:        int64(buf[i].nread),
			NumWrite:       int64(buf[i].nwrite),
			TotalReadTime:  time.Duration(buf[i].readtime),
			TotalWriteTime: time.Duration(buf[i].writetime),
			ReadLatency:    time.Duration(buf[i].readlat),
			WriteLatency:   time.Duration(buf[i].writelat),
		}
	}
	return stats, nil
}
