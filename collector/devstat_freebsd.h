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
