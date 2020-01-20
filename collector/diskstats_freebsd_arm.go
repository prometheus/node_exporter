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
// +build arm

package collector

const (
	//size of devstat returned for this arch
	sizeOfDevstat = 240

	//byte offsets for stats for this arch
	busyTimeFracStart = 176
	busyTimeFracStop  = 184
	busyTimeSecStart  = 172
	busyTimeSecStop   = 176

	bytesReadStart  = 68
	bytesReadStop   = 76
	bytesWriteStart = 76
	bytesWriteStop  = 84

	deviceNameStart = 36
	deviceNameStop  = 52

	durationReadFracStart  = 140
	durationReadFracStop   = 148
	durationReadSecStart   = 136
	durationReadSecStop    = 140
	durationWriteFracStart = 152
	durationWriteFracStop  = 160
	durationWriteSecStart  = 148
	durationWriteSecStop   = 152

	operationsReadStart  = 100
	operationsReadStop   = 108
	operationsWriteStart = 108
	operationsWriteStop  = 116

	unitNumberStart = 56
	unitNumberStop  = 60
)
