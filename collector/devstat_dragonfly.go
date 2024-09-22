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

//go:build !nodevstat
// +build !nodevstat

package collector

import (
	"errors"
	"fmt"
	"log/slog"

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

        return stats;
}
*/
import "C"

const (
	devstatSubsystem = "devstat"
)

type devstatCollector struct {
	bytesDesc     *prometheus.Desc
	transfersDesc *prometheus.Desc
	blocksDesc    *prometheus.Desc
	logger        *slog.Logger
}

func init() {
	registerCollector("devstat", defaultDisabled, NewDevstatCollector)
}

// NewDevstatCollector returns a new Collector exposing Device stats.
func NewDevstatCollector(logger *slog.Logger) (Collector, error) {
	return &devstatCollector{
		bytesDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, devstatSubsystem, "bytes_total"),
			"The total number of bytes transferred for reads and writes on the device.",
			[]string{"device"}, nil,
		),
		transfersDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, devstatSubsystem, "transfers_total"),
			"The total number of transactions completed.",
			[]string{"device"}, nil,
		),
		blocksDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, devstatSubsystem, "blocks_total"),
			"The total number of bytes given in terms of the devices blocksize.",
			[]string{"device"}, nil,
		),
		logger: logger,
	}, nil
}

func (c *devstatCollector) Update(ch chan<- prometheus.Metric) error {
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

		ch <- prometheus.MustNewConstMetric(c.bytesDesc, prometheus.CounterValue, float64(stats.bytes), device)
		ch <- prometheus.MustNewConstMetric(c.transfersDesc, prometheus.CounterValue, float64(stats.transfers), device)
		ch <- prometheus.MustNewConstMetric(c.blocksDesc, prometheus.CounterValue, float64(stats.blocks), device)
	}

	return nil
}
