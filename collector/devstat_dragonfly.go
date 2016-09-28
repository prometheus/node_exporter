// Copyright 2016 The Prometheus Authors
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
#cgo LDFLAGS: -ldevstat
#include <devstat.h>
#include <stdlib.h>
#include <string.h>

typedef struct {
	char		device[DEVSTAT_NAME_LEN];
	int		unit;
	uint64_t	bytes;
	uint64_t	transfers;
	uint64_t	blocks;
	double		kb_per_transfer;
	double		transfers_per_second;
	double		mb_per_second;
	double		blocks_per_second;
	double		ms_per_transaction;
} Stats;

int _get_ndevs() {
	struct statinfo current;
	int num_devices;

	current.dinfo = (struct devinfo *)calloc(1, sizeof(struct devinfo));
	if (current.dinfo == NULL)
		return -2;

	checkversion();

	if (getdevs(&current) == -1)
		return -1;

	return current.dinfo->numdevs;
}

Stats _get_stats(int i) {
	struct statinfo current;
	int num_devices;

	current.dinfo = (struct devinfo *)calloc(1, sizeof(struct devinfo));
	getdevs(&current);

	num_devices = current.dinfo->numdevs;
	Stats stats;

	uint64_t total_bytes, total_transfers, total_blocks;
	long double kb_per_transfer, transfers_per_second, mb_per_second, blocks_per_second, ms_per_transaction;

	strcpy(stats.device, current.dinfo->devices[i].device_name);
	stats.unit = current.dinfo->devices[i].unit_number;
	compute_stats(&current.dinfo->devices[i],
		NULL,
		1.0,
		&total_bytes,
		&total_transfers,
		&total_blocks,
		&kb_per_transfer,
		&transfers_per_second,
		&mb_per_second,
		&blocks_per_second,
		&ms_per_transaction);

	stats.bytes = total_bytes;
	stats.transfers = total_transfers;
	stats.blocks = total_blocks;
	stats.kb_per_transfer =	kb_per_transfer;
	stats.transfers_per_second = transfers_per_second;
	stats.mb_per_second = mb_per_second;
	stats.blocks_per_second = blocks_per_second;
	stats.ms_per_transaction = ms_per_transaction;

        return stats;
}
*/
import "C"

const (
	devstatSubsystem = "devstat"
)

type devstatCollector struct {
	bytes_total             *prometheus.CounterVec
	transfers_total         *prometheus.CounterVec
	blocks_total            *prometheus.CounterVec
	bytes_per_transfer      *prometheus.GaugeVec
	transfers_per_second    *prometheus.GaugeVec
	bytes_per_second        *prometheus.GaugeVec
	blocks_per_second       *prometheus.GaugeVec
	seconds_per_transaction *prometheus.GaugeVec
}

func init() {
	Factories["devstat"] = NewDevstatCollector
}

// Takes a prometheus registry and returns a new Collector exposing
// Device stats.
func NewDevstatCollector() (Collector, error) {
	return &devstatCollector{
		bytes_total: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Subsystem: devstatSubsystem,
				Name:      "bytes_total",
				Help:      "The total number of bytes transferred for reads and writes on the device.",
			},
			[]string{"device"},
		),
		transfers_total: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Subsystem: devstatSubsystem,
				Name:      "transfers_total",
				Help:      "The total number of transactions completed.",
			},
			[]string{"device"},
		),
		blocks_total: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Subsystem: devstatSubsystem,
				Name:      "blocks_total",
				Help:      "The total number of bytes given in terms of the devices blocksize.",
			},
			[]string{"device"},
		),
		bytes_per_transfer: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Subsystem: devstatSubsystem,
				Name:      "bytes_per_transfer",
				Help:      "The average number of bytes per transfer.",
			},
			[]string{"device"},
		),
		transfers_per_second: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Subsystem: devstatSubsystem,
				Name:      "transfers_per_second",
				Help:      "The average number of transfers per second.",
			},
			[]string{"device"},
		),
		bytes_per_second: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Subsystem: devstatSubsystem,
				Name:      "bytes_per_second",
				Help:      "The average bytes per second.",
			},
			[]string{"device"},
		),
		blocks_per_second: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Subsystem: devstatSubsystem,
				Name:      "blocks_per_second",
				Help:      "The average blocks per second.",
			},
			[]string{"device"},
		),
		seconds_per_transaction: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Subsystem: devstatSubsystem,
				Name:      "seconds_per_transaction",
				Help:      "The average number of seconds per transaction.",
			},
			[]string{"device"},
		),
	}, nil
}

func (c *devstatCollector) Update(ch chan<- prometheus.Metric) (err error) {
	count := C._get_ndevs()
	if count == -1 {
		return errors.New("getdevs() failed")
	}
	if count == -2 {
		return errors.New("calloc() failed")
	}

	for i := C.int(0); i < count; i++ {
		stats := C._get_stats(i)
		device := fmt.Sprintf("%s%d", C.GoString(&stats.device[0]), stats.unit)

		c.bytes_total.With(prometheus.Labels{"device": device}).Set(float64(stats.bytes))
		c.transfers_total.With(prometheus.Labels{"device": device}).Set(float64(stats.transfers))
		c.blocks_total.With(prometheus.Labels{"device": device}).Set(float64(stats.blocks))
		c.bytes_per_transfer.With(prometheus.Labels{"device": device}).Set(float64(stats.kb_per_transfer) * 1000)
		c.bytes_per_second.With(prometheus.Labels{"device": device}).Set(float64(stats.mb_per_second) * 1000000)
		c.transfers_per_second.With(prometheus.Labels{"device": device}).Set(float64(stats.transfers_per_second))
		c.blocks_per_second.With(prometheus.Labels{"device": device}).Set(float64(stats.blocks_per_second))
		c.seconds_per_transaction.With(prometheus.Labels{"device": device}).Set(float64(stats.ms_per_transaction) / 1000)
	}

	c.bytes_total.Collect(ch)
	c.transfers_total.Collect(ch)
	c.blocks_total.Collect(ch)
	c.bytes_per_transfer.Collect(ch)
	c.bytes_per_second.Collect(ch)
	c.transfers_per_second.Collect(ch)
	c.blocks_per_second.Collect(ch)
	c.seconds_per_transaction.Collect(ch)

	return err
}
