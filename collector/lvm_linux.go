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

// +build linux

package collector

import (
	// #include <lvm2app.h>
	// #cgo LDFLAGS: -llvm2app
	"C"
	"errors"
	"unsafe"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	lvmVolumeGroupSizeDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "lvm", "volume_group_size_bytes"),
		"Size of the LVM volume groups in bytes.",
		[]string{"volume_group"},
		nil)
	lvmLogicalVolumeSizeDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "lvm", "logical_volume_size_bytes"),
		"Size of the LVM logical volumes in bytes.",
		[]string{"volume_group", "logical_volume"},
		nil)
)

type lvmCollector struct{}

func init() {
	Factories["lvm"] = NewLvmCollector
}

// NewLvmCollector creates a new collector for LVM volume group and
// logical volume sizes.
func NewLvmCollector() (Collector, error) {
	return &lvmCollector{}, nil
}

func (c *lvmCollector) Update(ch chan<- prometheus.Metric) error {
	l := C.lvm_init(nil)
	if l == nil {
		return errors.New("Failed to initialize LVM")
	}
	defer C.lvm_quit(l)

	// Iterate over all volume groups.
	vgs := C.lvm_list_vg_names(l)
	for vgItem := vgs.n; vgItem != vgs; vgItem = vgItem.n {
		vgNameC := (*C.struct_lvm_str_list)(unsafe.Pointer(vgItem)).str
		vgName := C.GoString(vgNameC)
		mode := C.CString("r")
		vg := C.lvm_vg_open(l, vgNameC, mode, 0)
		C.free(unsafe.Pointer(mode))
		if vg == nil {
			return errors.New("Failed to open volume group " + vgName)
		}
		defer C.lvm_vg_close(vg)

		// Report volume group size.
		ch <- prometheus.MustNewConstMetric(
			lvmVolumeGroupSizeDesc,
			prometheus.GaugeValue,
			float64(C.lvm_vg_get_size(vg)),
			vgName)

		// Iterate over all logical volumes within the volume group.
		lvs := C.lvm_vg_list_lvs(vg)
		for lvItem := lvs.n; lvItem != lvs; lvItem = lvItem.n {
			lv := (*C.struct_lvm_lv_list)(unsafe.Pointer(lvItem)).lv
			lvNameC := C.lvm_lv_get_name(lv)
			lvName := C.GoString(lvNameC)

			// Report logical volume size.
			ch <- prometheus.MustNewConstMetric(
				lvmLogicalVolumeSizeDesc,
				prometheus.GaugeValue,
				float64(C.lvm_lv_get_size(lv)),
				vgName,
				lvName)
		}
	}
	return nil
}
