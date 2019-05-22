// +build linux

package perf

import (
	"encoding/binary"
	"fmt"
	"syscall"

	"go.uber.org/multierr"
	"golang.org/x/sys/unix"
)

// ErrNoLeader is returned when a leader of a GroupProfiler is not defined.
var ErrNoLeader = fmt.Errorf("No leader defined")

// GroupProfileValue is returned from a GroupProfiler.
type GroupProfileValue struct {
	Events      uint64
	TimeEnabled uint64
	TimeRunning uint64
	Values      []uint64
}

// GroupProfiler is used to setup a group profiler.
type GroupProfiler interface {
	Start() error
	Reset() error
	Stop() error
	Close() error
	Profile() (*GroupProfileValue, error)
}

// groupProfiler implements the GroupProfiler interface.
type groupProfiler struct {
	fds []int // leader is always element 0
}

// NewGroupProfiler returns a GroupProfiler.
func NewGroupProfiler(pid, cpu, opts int, eventAttrs ...unix.PerfEventAttr) (GroupProfiler, error) {
	fds := make([]int, len(eventAttrs))

	for i, eventAttr := range eventAttrs {
		// common configs
		eventAttr.Size = EventAttrSize
		eventAttr.Sample_type = PERF_SAMPLE_IDENTIFIER

		// Leader fd must be opened first
		if i == 0 {
			// leader specific configs
			eventAttr.Bits = unix.PerfBitDisabled | unix.PerfBitExcludeHv
			eventAttr.Read_format = unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED | unix.PERF_FORMAT_GROUP

			fd, err := unix.PerfEventOpen(
				&eventAttr,
				pid,
				cpu,
				-1,
				opts,
			)
			if err != nil {
				return nil, err
			}
			fds[i] = fd
			continue
		}

		// non leader configs
		eventAttr.Read_format = unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED | unix.PERF_FORMAT_GROUP
		eventAttr.Bits = unix.PerfBitExcludeHv

		fd, err := unix.PerfEventOpen(
			&eventAttr,
			pid,
			cpu,
			fds[0],
			opts,
		)
		if err != nil {
			// cleanup any old Fds
			for ii, fd2 := range fds {
				if ii == i {
					break
				}
				err = multierr.Append(err, unix.Close(fd2))
			}
			return nil, err
		}
		fds[i] = fd
	}

	return &groupProfiler{
		fds: fds,
	}, nil
}

// Start is used to start the GroupProfiler.
func (p *groupProfiler) Start() error {
	if len(p.fds) == 0 {
		return ErrNoLeader
	}
	return unix.IoctlSetInt(p.fds[0], unix.PERF_EVENT_IOC_ENABLE, 0)
}

// Reset is used to reset the GroupProfiler.
func (p *groupProfiler) Reset() error {
	if len(p.fds) == 0 {
		return ErrNoLeader
	}
	return unix.IoctlSetInt(p.fds[0], unix.PERF_EVENT_IOC_RESET, 0)
}

// Stop is used to stop the GroupProfiler.
func (p *groupProfiler) Stop() error {
	if len(p.fds) == 0 {
		return ErrNoLeader
	}
	return unix.IoctlSetInt(p.fds[0], unix.PERF_EVENT_IOC_DISABLE, 0)
}

// Close is used to close the GroupProfiler.
func (p *groupProfiler) Close() error {
	var err error
	for _, fd := range p.fds {
		err = multierr.Append(err, unix.Close(fd))
	}
	return err
}

// Profile is used to return the GroupProfileValue of the GroupProfiler.
func (p *groupProfiler) Profile() (*GroupProfileValue, error) {
	nEvents := len(p.fds)
	if nEvents == 0 {
		return nil, ErrNoLeader
	}

	// read format of the raw event looks like this:
	/*
		     struct read_format {
			 u64 nr;            // The number of events /
			 u64 time_enabled;  // if PERF_FORMAT_TOTAL_TIME_ENABLED
			 u64 time_running;  // if PERF_FORMAT_TOTAL_TIME_RUNNING
			 struct {
			     u64 value;     // The value of the event
			     u64 id;        // if PERF_FORMAT_ID
			 } values[nr];
		     };
	*/

	buf := make([]byte, 24+8*nEvents)
	_, err := syscall.Read(p.fds[0], buf)
	if err != nil {
		return nil, err
	}

	val := &GroupProfileValue{
		Events:      binary.LittleEndian.Uint64(buf[0:8]),
		TimeEnabled: binary.LittleEndian.Uint64(buf[8:16]),
		TimeRunning: binary.LittleEndian.Uint64(buf[16:24]),
		Values:      make([]uint64, len(p.fds)),
	}

	offset := 24
	for i := range p.fds {
		val.Values[i] = binary.LittleEndian.Uint64(buf[offset : offset+8])
		offset += 8
	}

	return val, nil
}
