// +build linux

package perf

import (
	"fmt"
	"strconv"
	"strings"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	// PERF_TYPE_TRACEPOINT is a kernel tracepoint.
	PERF_TYPE_TRACEPOINT = 2
)

// AvailableEvents returns a mapping of available subsystems and their
// corresponding list of available events.
func AvailableEvents() (map[string][]string, error) {
	events := map[string][]string{}
	// BUG(hodgesds): this should ideally check mounts for debugfs
	rawEvents, err := fileToStrings(TracingDir + "/available_events")
	// Events are colon delimited by type so parse the type and add sub
	// events appropriately.
	if err != nil {
		return events, err
	}
	for _, rawEvent := range rawEvents {
		splits := strings.Split(rawEvent, ":")
		if len(splits) <= 1 {
			continue
		}
		eventTypeEvents, found := events[splits[0]]
		if found {
			events[splits[0]] = append(eventTypeEvents, splits[1])
			continue
		}
		events[splits[0]] = []string{splits[1]}
	}
	return events, err
}

// AvailableSubsystems returns a slice of available subsystems.
func AvailableSubsystems() ([]string, error) {
	subsystems := []string{}
	// BUG(hodgesds): this should ideally check mounts for debugfs
	rawEvents, err := fileToStrings(TracingDir + "/available_events")
	// Events are colon delimited by type so parse the type and add sub
	// events appropriately.
	if err != nil {
		return subsystems, err
	}
	for _, rawEvent := range rawEvents {
		splits := strings.Split(rawEvent, ":")
		if len(splits) <= 1 {
			continue
		}
		subsystems = append(subsystems, splits[0])
	}
	return subsystems, nil
}

// AvailableTracers returns the list of available tracers.
func AvailableTracers() ([]string, error) {
	return fileToStrings(TracingDir + "/available_tracers")
}

// CurrentTracer returns the current tracer.
func CurrentTracer() (string, error) {
	res, err := fileToStrings(TracingDir + "/current_tracer")
	return res[0], err
}

// GetTracepointConfig is used to get the configuration for a trace event.
func GetTracepointConfig(subsystem, event string) (uint64, error) {
	res, err := fileToStrings(
		TracingDir + fmt.Sprintf("/events/%s/%s/id", subsystem, event))
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(res[0], 10, 64)
}

// ProfileTracepoint is used to profile a kernel tracepoint event for a
// specific PID. Events can be listed with `perf list` for Tracepoint Events or
// in the /sys/kernel/debug/tracing/events directory with the kind being the
// directory and the event being the subdirectory.
func ProfileTracepoint(subsystem, event string, pid, cpu int, opts ...int) (BPFProfiler, error) {
	config, err := GetTracepointConfig(subsystem, event)
	if err != nil {
		return nil, err
	}
	eventAttr := &unix.PerfEventAttr{
		Type:        PERF_TYPE_TRACEPOINT,
		Config:      config,
		Size:        uint32(unsafe.Sizeof(unix.PerfEventAttr{})),
		Bits:        unix.PerfBitDisabled | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
		Sample_type: PERF_SAMPLE_IDENTIFIER,
	}
	var eventOps int
	if len(opts) > 0 {
		eventOps = opts[0]
	}
	fd, err := unix.PerfEventOpen(
		eventAttr,
		pid,
		cpu,
		-1,
		eventOps,
	)
	if err != nil {
		return nil, err
	}

	return &profiler{
		fd: fd,
	}, nil
}

// TracepointEventAttr is used to return an PerfEventAttr for a trace event.
func TracepointEventAttr(subsystem, event string) (*unix.PerfEventAttr, error) {
	config, err := GetTracepointConfig(subsystem, event)
	if err != nil {
		return nil, err
	}
	return &unix.PerfEventAttr{
		Type:        PERF_TYPE_TRACEPOINT,
		Config:      config,
		Size:        uint32(unsafe.Sizeof(unix.PerfEventAttr{})),
		Bits:        unix.PerfBitDisabled | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
		Sample_type: PERF_SAMPLE_IDENTIFIER,
	}, nil
}
