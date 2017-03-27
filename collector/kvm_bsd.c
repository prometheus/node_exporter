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

// +build !nomeminfo
// +build freebsd dragonfly

#include <fcntl.h>
#include <kvm.h>
#include <limits.h>
#include <paths.h>
#include <stdlib.h>

int _kvm_swap_used_pages(uint64_t *out) {
	const int total_only = 1; // from kvm_getswapinfo(3)

	kvm_t *kd;
	struct kvm_swap current;

	kd = kvm_open(NULL, _PATH_DEVNULL, NULL, O_RDONLY, NULL);
	if (kd == NULL) {
		return -1;
	}

	if (kvm_getswapinfo(kd, &current, total_only, 0) == -1) {
		goto error1;
	}

	if (kvm_close(kd) != 0) {
		return -1;
	}
	kd = NULL;

	*out = current.ksw_used;
	return 0;

error1:
	if (kd != NULL) {
		kvm_close(kd);
	}

	return -1;
}
