// +build !nodevstat

#include <devstat.h>
#include <fcntl.h>
#include <libgeom.h>
#include <limits.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <devstat_freebsd.h>

int _get_ndevs() {
	struct statinfo current;
	struct devinfo info = {};
	current.dinfo = &info;

	devstat_checkversion(NULL);

	if (devstat_getdevs(NULL, &current) == -1)
		return -1;

	return current.dinfo->numdevs;
}

Stats _get_stats(int i) {
	struct statinfo current;
	struct devinfo info = {};
	current.dinfo = &info;

	devstat_getdevs(NULL, &current);

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
