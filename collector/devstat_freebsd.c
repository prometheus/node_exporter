// Copyright 2017 The Prometheus Authors
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

#include <devstat.h>
#include <fcntl.h>
#include <libgeom.h>
#include <limits.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <devstat_freebsd.h>


int _get_stats(struct devinfo *info, Stats **stats) {
	struct statinfo current;
	current.dinfo = info;

	if (devstat_getdevs(NULL, &current) == -1) {
		return -1;
	}

	Stats *p = (Stats*)calloc(current.dinfo->numdevs, sizeof(Stats));
	for (int i = 0; i < current.dinfo->numdevs; i++) {
		uint64_t bytes_read, bytes_write, bytes_free;
		uint64_t transfers_other, transfers_read, transfers_write, transfers_free;
		long double duration_other, duration_read, duration_write, duration_free;
		long double busy_time;
		uint64_t blocks;

		strcpy(p[i].device, current.dinfo->devices[i].device_name);
		p[i].unit = current.dinfo->devices[i].unit_number;
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

		p[i].bytes.read = bytes_read;
		p[i].bytes.write = bytes_write;
		p[i].bytes.free = bytes_free;
		p[i].transfers.other = transfers_other;
		p[i].transfers.read = transfers_read;
		p[i].transfers.write = transfers_write;
		p[i].transfers.free = transfers_free;
		p[i].duration.other = duration_other;
		p[i].duration.read = duration_read;
		p[i].duration.write = duration_write;
		p[i].duration.free = duration_free;
		p[i].busyTime = busy_time;
		p[i].blocks = blocks;
	}

	*stats = p;
	return current.dinfo->numdevs;
}
