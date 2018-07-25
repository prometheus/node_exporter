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

package collector

import (
	"fmt"
	"sync"
)

// #cgo LDFLAGS: -lkvm
// #include "kvm_bsd.h"
import "C"

type kvm struct {
	mu     sync.Mutex
	hasErr bool
}

func (k *kvm) SwapUsedPages() (value uint64, err error) {
	k.mu.Lock()
	defer k.mu.Unlock()
	if C._kvm_swap_used_pages((*C.uint64_t)(&value)) == -1 {
		k.hasErr = true
		return 0, fmt.Errorf("couldn't get kvm stats")
	}

	return value, nil
}
