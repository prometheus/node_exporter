// Copyright 2020 The Prometheus Authors
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

// +build freebsd
// +build arm64

package collector

const (
	//size of devstat returned for this arch
	sizeOfDevstat = 288

	//byte offsets for stats for this arch
	busyTimeFracStart      = 200
	busyTimeFracStop       = 208
	busyTimeSecStart       = 192
	busyTimeSecStop        = 200
	bytesReadStart         = 72
	bytesReadStop          = 80
	bytesWriteStart        = 80
	bytesWriteStop         = 88
	deviceNameStart        = 44
	deviceNameStop         = 60
	durationReadFracStart  = 152
	durationReadFracStop   = 160
	durationReadSecStart   = 144
	durationReadSecStop    = 152
	durationWriteFracStart = 168
	durationWriteFracStop  = 176
	durationWriteSecStart  = 160
	durationWriteSecStop   = 168
	operationsReadStart    = 104
	operationsReadStop     = 112
	operationsWriteStart   = 112
	operationsWriteStop    = 120
	unitNumberStart        = 60
	unitNumberStop         = 64
)
