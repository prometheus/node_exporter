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


int _get_ndevs();
int _get_stats(struct devinfo *info, Stats **stats);
