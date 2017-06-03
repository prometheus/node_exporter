// Package iostat presents I/O statistics.
package iostat

import "time"

// DriveStats represents I/O statistics of a drive.
type DriveStats struct {
	Name      string // drive name
	Size      int64  // total drive size in bytes
	BlockSize int64  // block size in bytes

	BytesRead      int64
	BytesWritten   int64
	NumRead        int64
	NumWrite       int64
	TotalReadTime  time.Duration
	TotalWriteTime time.Duration
	ReadLatency    time.Duration
	WriteLatency   time.Duration
}
