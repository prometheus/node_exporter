// Copyright 2015 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build !nodevstat

package collector

import (
	"errors"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

/*
#cgo LDFLAGS: -ldevstat -lkvm
#include <devstat.h>
#include <fcntl.h>
#include <libgeom.h>
#include <limits.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

typedef struct {
	uint64_t	read;
	uint64_t	write;
	uint64_t	free;
} Bytes;

typedef struct {
	uint64_t	other;
	uint64_t	read;
	uint64_t	write;
	uint64_t	free;
} Transfers;

typedef struct {
	double		other;
	double		read;
	double		write;
	double		free;
} Duration;

typedef struct {
	char		device[DEVSTAT_NAME_LEN];
	int		unit;
	Bytes		bytes;
	Transfers	transfers;
	Duration	duration;
	long		busyTime;
	uint64_t	blocks;
} Stats;

int _get_ndevs() {
	struct statinfo current;
	int num_devices;

	current.dinfo = (struct devinfo *)calloc(1, sizeof(struct devinfo));
	if (current.dinfo == NULL)
		return -2;

	devstat_checkversion(NULL);

	if (devstat_getdevs(NULL, &current) == -1)
		return -1;

	return current.dinfo->numdevs;
}

Stats _get_stats(int i) {
	struct statinfo current;
	int num_devices;

	current.dinfo = (struct devinfo *)calloc(1, sizeof(struct devinfo));
	devstat_getdevs(NULL, &current);

	num_devices = current.dinfo->numdevs;
	Stats stats;
	uint64_t bytes_read, bytes_write, bytes_free;
	uint64_t transfers_other, transfers_read, transfers_write, transfers_free;
	long double duration_other, duration_read, duration_write, duration_free;
	long double busy_time;
	uint64_t blocks;

	strcpy(stats.device, current.dinfo->devices[i].device_name);
	stats.unit = current.dinfo->devices[i].unit_number;
	devstat_compute_statistics(&current.dinfo->devices[i],
		NULL,
		1.0,
		DSM_TOTAL_BYTES_READ, &bytes_read,
		DSM_TOTAL_BYTES_WRITE, &bytes_write,
		DSM_TOTAL_BYTES_FREE, &bytes_free,
		DSM_TOTAL_TRANSFERS_OTHER, &transfers_other,
		DSM_TOTAL_TRANSFERS_READ, &transfers_read,
		DSM_TOTAL_TRANSFERS_WRITE, &transfers_write,
		DSM_TOTAL_TRANSFERS_FREE, &transfers_free,
		DSM_TOTAL_DURATION_OTHER, &duration_other,
		DSM_TOTAL_DURATION_READ, &duration_read,
		DSM_TOTAL_DURATION_WRITE, &duration_write,
		DSM_TOTAL_DURATION_FREE, &duration_free,
		DSM_TOTAL_BUSY_TIME, &busy_time,
		DSM_TOTAL_BLOCKS, &blocks,
		DSM_NONE);

	stats.bytes.read = bytes_read;
	stats.bytes.write = bytes_write;
	stats.bytes.free = bytes_free;
	stats.transfers.other = transfers_other;
	stats.transfers.read = transfers_read;
	stats.transfers.write = transfers_write;
	stats.transfers.free = transfers_free;
	stats.duration.other = duration_other;
	stats.duration.read = duration_read;
	stats.duration.write = duration_write;
	stats.duration.free = duration_free;
	stats.busyTime = busy_time;
	stats.blocks = blocks;

	return stats;
}
*/
import "C"

const (
	devstatSubsystem = "devstat"
)

type devstatCollector struct {
	bytes       *prometheus.CounterVec
	bytes_total *prometheus.CounterVec
	transfers   *prometheus.CounterVec
	duration    *prometheus.CounterVec
	busyTime    *prometheus.CounterVec
	blocks      *prometheus.CounterVec
}

func init() {
	Factories["devstat"] = NewDevstatCollector
}

// Takes a prometheus registry and returns a new Collector exposing
// Device stats.
func NewDevstatCollector() (Collector, error) {
	return &devstatCollector{
		bytes: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Subsystem: devstatSubsystem,
				Name:      "bytes_total",
				Help:      "The total number of bytes in transactions.",
			},
			[]string{"device", "type"},
		),
		transfers: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Subsystem: devstatSubsystem,
				Name:      "transfers_total",
				Help:      "The total number of transactions.",
			},
			[]string{"device", "type"},
		),
		duration: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Subsystem: devstatSubsystem,
				Name:      "duration_seconds_total",
				Help:      "The total duration of transactions in seconds.",
			},
			[]string{"device", "type"},
		),
		busyTime: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Subsystem: devstatSubsystem,
				Name:      "busy_time_seconds_total",
				Help:      "Total time the device had one or more transactions outstanding in seconds.",
			},
			[]string{"device"},
		),
		blocks: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Subsystem: devstatSubsystem,
				Name:      "blocks_transferred_total",
				Help:      "The total number of blocks transferred.",
			},
			[]string{"device"},
		),
	}, nil
}

func (c *devstatCollector) Update(ch chan<- prometheus.Metric) (err error) {
	count := C._get_ndevs()
	if count == -1 {
		return errors.New("devstat_getdevs() failed")
	}
	if count == -2 {
		return errors.New("calloc() failed")
	}

	for i := C.int(0); i < count; i++ {
		stats := C._get_stats(i)
		device := fmt.Sprintf("%s%d", C.GoString(&stats.device[0]), stats.unit)
		// Free metrics are disabled for now, please see PR #88 for more details.
		c.bytes.With(prometheus.Labels{"device": device, "type": "read"}).Set(float64(stats.bytes.read))
		c.bytes.With(prometheus.Labels{"device": device, "type": "write"}).Set(float64(stats.bytes.write))
		//c.bytes.With(prometheus.Labels{"device": device, "type": "free"}).Set(float64(stats.bytes.free))
		c.transfers.With(prometheus.Labels{"device": device, "type": "other"}).Set(float64(stats.transfers.other))
		c.transfers.With(prometheus.Labels{"device": device, "type": "read"}).Set(float64(stats.transfers.read))
		c.transfers.With(prometheus.Labels{"device": device, "type": "write"}).Set(float64(stats.transfers.write))
		//c.transfers.With(prometheus.Labels{"device": device, "type": "free"}).Set(float64(stats.transfers.free))
		c.duration.With(prometheus.Labels{"device": device, "type": "other"}).Set(float64(stats.duration.other))
		c.duration.With(prometheus.Labels{"device": device, "type": "read"}).Set(float64(stats.duration.read))
		c.duration.With(prometheus.Labels{"device": device, "type": "write"}).Set(float64(stats.duration.write))
		//c.duration.With(prometheus.Labels{"device": device, "type": "free"}).Set(float64(stats.duration.free))
		c.busyTime.With(prometheus.Labels{"device": device}).Set(float64(stats.busyTime))
		c.blocks.With(prometheus.Labels{"device": device}).Set(float64(stats.blocks))
	}

	c.bytes.Collect(ch)
	c.transfers.Collect(ch)
	c.duration.Collect(ch)
	c.busyTime.Collect(ch)
	c.blocks.Collect(ch)

	return err
}
