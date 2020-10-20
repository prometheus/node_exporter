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

// +build !nvme

package collector

/*
#include <linux/nvme_ioctl.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	nvmeSysPath            = "/sys/class/nvme"
	nvmeCollectorSubsystem = "nvme"
)

var (
	nvmeInterestedDevices = kingpin.Flag("collector.nvme.devices", "Comma-separated list of devices, e.g. /dev/nvme0,/dev/nvme1").String()
	nvmeLabels            = []string{"device"}
)

type nvmeCollector struct {
	devices  []*nvmeDevice
	logger   log.Logger
	smartLog *nvmeSmartLogDescriptor
}

type nvmeDevice struct {
	path string
	file *os.File
}

type nvmeSmartLogDescriptor struct {
	spareError                         *typedDesc
	temperatureError                   *typedDesc
	reliabilityError                   *typedDesc
	backupError                        *typedDesc
	compositeTemperature               *typedDesc
	availableSpare                     *typedDesc
	availableSpareThreshold            *typedDesc
	percentageUsed                     *typedDesc
	dataRead                           *typedDesc
	dataWritten                        *typedDesc
	hostReadCommands                   *typedDesc
	hostWriteCommands                  *typedDesc
	controllerBusyTime                 *typedDesc
	powerCycles                        *typedDesc
	powerOnHours                       *typedDesc
	unsafeShutdowns                    *typedDesc
	mediaAndDataIntegrityErrors        *typedDesc
	numberOfErrorInformationLogEntries *typedDesc
}

type nvmeSmartLog struct {
	spareError                         bool
	temperatureError                   bool
	reliabilityError                   bool
	backupError                        bool
	compositeTemperature               uint16
	availableSpare                     uint8
	availableSpareThreshold            uint8
	percentageUsed                     uint8
	dataRead                           uint64
	dataWritten                        uint64
	hostReadCommands                   uint64
	hostWriteCommands                  uint64
	controllerBusyTime                 uint64
	powerCycles                        uint64
	powerOnHours                       uint64
	unsafeShutdowns                    uint64
	mediaAndDataIntegrityErrors        uint64
	numberOfErrorInformationLogEntries uint64
}

func init() {
	registerCollector("nvme", defaultDisabled, NewNVMeCollector)
}

// NewNVMeCollector returns a new Collector exposing SMART log of NVMe devices.
func NewNVMeCollector(logger log.Logger) (Collector, error) {
	interestedDevices, err := nvmeParseInterestedList(*nvmeInterestedDevices)
	if err != nil {
		interestedDevices, err = nvmeFindDevices()
		if err != nil {
			return nil, fmt.Errorf("--collector.nvme.devices is not provided and finding devices from %s failed: %v", nvmeSysPath, err)
		}
	}

	var devices []*nvmeDevice
	for _, devPath := range interestedDevices {
		device, err := nvmeNewDevice(devPath)
		if err != nil {
			return nil, fmt.Errorf("unable to %v", err)
		}

		devices = append(devices, device)
	}

	return &nvmeCollector{
		devices:  devices,
		smartLog: nvmeNewSmartLogDescriptor(),
		logger:   logger,
	}, nil
}

func (c *nvmeCollector) Update(ch chan<- prometheus.Metric) error {
	for _, device := range c.devices {
		smartLog, err := device.smartLog()
		if err != nil {
			return err
		}
		c.smartLog.update(ch, smartLog, device.path)
	}
	return nil
}

func (d *nvmeSmartLogDescriptor) update(ch chan<- prometheus.Metric, smartLog *nvmeSmartLog, label string) {
	updateBool(ch, d.spareError, smartLog.spareError, label)
	updateBool(ch, d.temperatureError, smartLog.temperatureError, label)
	updateBool(ch, d.reliabilityError, smartLog.reliabilityError, label)
	updateBool(ch, d.backupError, smartLog.backupError, label)
	update(ch, d.compositeTemperature, float64(smartLog.compositeTemperature), label)
	update(ch, d.availableSpare, float64(smartLog.availableSpare), label)
	update(ch, d.availableSpareThreshold, float64(smartLog.availableSpareThreshold), label)
	update(ch, d.percentageUsed, float64(smartLog.percentageUsed), label)
	update(ch, d.dataRead, float64(smartLog.dataRead), label)
	update(ch, d.dataWritten, float64(smartLog.dataWritten), label)
	update(ch, d.hostReadCommands, float64(smartLog.hostReadCommands), label)
	update(ch, d.hostWriteCommands, float64(smartLog.hostWriteCommands), label)
	update(ch, d.controllerBusyTime, float64(smartLog.controllerBusyTime), label)
	update(ch, d.powerCycles, float64(smartLog.powerCycles), label)
	update(ch, d.powerOnHours, float64(smartLog.powerOnHours), label)
	update(ch, d.unsafeShutdowns, float64(smartLog.unsafeShutdowns), label)
	update(ch, d.mediaAndDataIntegrityErrors, float64(smartLog.mediaAndDataIntegrityErrors), label)
	update(ch, d.numberOfErrorInformationLogEntries, float64(smartLog.numberOfErrorInformationLogEntries), label)
}

func update(ch chan<- prometheus.Metric, desc *typedDesc, val float64, label string) {
	ch <- desc.mustNewConstMetric(val, label)
}

func updateBool(ch chan<- prometheus.Metric, desc *typedDesc, val bool, label string) {
	var v float64
	if val {
		v = 1
	} else {
		v = 0
	}
	update(ch, desc, v, label)
}

func nvmeNewSmartLogDescriptor() *nvmeSmartLogDescriptor {
	d := &nvmeSmartLogDescriptor{
		spareError: nvmeNewDescriptor(
			"spare_error",
			"Indicates that available spare space has fallen below the threshold.",
		),
		temperatureError: nvmeNewDescriptor(
			"temperature_error",
			"Indicates that a temperature is above an over temperature threshold or below an under temperature threshold.",
		),
		reliabilityError: nvmeNewDescriptor(
			"reliability_error",
			"Indicates that the NVM subsystem reliability has been degraded due to significant media related errors or any internal error that degrades NVM subsystemreliability.",
		),
		backupError: nvmeNewDescriptor(
			"backup_error",
			"Indicates that the volatile memory backup device has failed.",
		),
		compositeTemperature: nvmeNewDescriptor(
			"composite_temperature",
			"Current composite temperature of the controller in degrees Kelvin.",
		),
		availableSpare: nvmeNewDescriptor(
			"available_spare",
			"Normalized percentage (0 to 100) of the remaining spare capacity available.",
		),
		availableSpareThreshold: nvmeNewDescriptor(
			"available_spare_threshold",
			"Threshold of Available Spare.",
		),
		percentageUsed: nvmeNewDescriptor(
			"percentage_used",
			"Vendor specific estimate of the percentage of NVM subsystem life (0-255).",
		),
		dataRead: nvmeNewDescriptor(
			"data_read",
			"Kilobytes the host has read from the controller",
		),
		dataWritten: nvmeNewDescriptor(
			"data_written",
			"Kilobytes the host has written to the controller",
		),
		hostReadCommands: nvmeNewDescriptor(
			"host_read_commands",
			"Number of read commands completed by the controller.",
		),
		hostWriteCommands: nvmeNewDescriptor(
			"host_write_commands",
			"Number of write commands completed by the controller.",
		),
		controllerBusyTime: nvmeNewDescriptor(
			"controller_busy_time",
			"Minutes the controller is busy with I/O commands.",
		),
		powerCycles: nvmeNewDescriptor(
			"power_cycles",
			"Number of power cycles.",
		),
		powerOnHours: nvmeNewDescriptor(
			"power_on_hours",
			"Hours of power-on.",
		),
		unsafeShutdowns: nvmeNewDescriptor(
			"unsafe_shutdowns",
			"Number of unsafe shutdowns.",
		),
		mediaAndDataIntegrityErrors: nvmeNewDescriptor(
			"media_and_data_integrity_errors",
			"Number of occurrences where the controller detected an unrecovered data integrity error.",
		),
		numberOfErrorInformationLogEntries: nvmeNewDescriptor(
			"number_of_error_information_log_entries",
			"Number of Error Information log entries over the life of the controller.",
		),
	}
	return d
}

func nvmeNewDescriptor(name, help string) *typedDesc {
	return &typedDesc{
		desc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nvmeCollectorSubsystem, name),
			help, nvmeLabels, nil,
		),
		valueType: prometheus.GaugeValue,
	}
}

func nvmeParseInterestedList(list string) ([]string, error) {
	if len(list) == 0 {
		return nil, errors.New("list of interested devices is empty")
	}
	return strings.Split(list, ","), nil
}

func nvmeNewDevice(path string) (*nvmeDevice, error) {
	file, err := os.OpenFile(path, os.O_RDWR, 0755)
	if err != nil {
		return nil, err
	}
	return &nvmeDevice{path, file}, nil
}

func (d *nvmeDevice) smartLog() (*nvmeSmartLog, error) {
	var data [512]byte

	cmd := C.struct_nvme_admin_cmd{
		opcode:   0x02,
		nsid:     0xffffffff,
		addr:     C.ulonglong(uintptr(unsafe.Pointer(&data))),
		data_len: 512,
		cdw10:    0x007F0002,
	}

	err := nvmeIoctlAdminCmd(d.file.Fd(), uintptr(unsafe.Pointer(&cmd)), unsafe.Sizeof(cmd))
	if err != nil {
		return nil, err
	}

	return nvmeParseSmartLog(data), nil
}

// nvmeParseSmartLog parses `data` based on [nvme specification](https://nvmexpress.org/wp-content/uploads/NVM_Express_1_2_Gold_20141209.pdf) (Page 89, Figure 79).
func nvmeParseSmartLog(data [512]byte) *nvmeSmartLog {
	return &nvmeSmartLog{
		spareError:                         (data[0] & 0x1) != 0,
		temperatureError:                   (data[0] & 0x2) != 0,
		reliabilityError:                   (data[0] & 0x4) != 0,
		backupError:                        (data[0] & 0x8) != 0,
		compositeTemperature:               *(*uint16)(unsafe.Pointer(&data[1])),
		availableSpare:                     data[3],
		availableSpareThreshold:            data[4],
		percentageUsed:                     data[5],
		dataRead:                           *(*uint64)(unsafe.Pointer(&data[32])) * 500, // bytes = unit * 1000 * 512
		dataWritten:                        *(*uint64)(unsafe.Pointer(&data[48])) * 500, // bytes = unit * 1000 * 512
		hostReadCommands:                   *(*uint64)(unsafe.Pointer(&data[64])),
		hostWriteCommands:                  *(*uint64)(unsafe.Pointer(&data[80])),
		controllerBusyTime:                 *(*uint64)(unsafe.Pointer(&data[96])),
		powerCycles:                        *(*uint64)(unsafe.Pointer(&data[112])),
		powerOnHours:                       *(*uint64)(unsafe.Pointer(&data[128])),
		unsafeShutdowns:                    *(*uint64)(unsafe.Pointer(&data[144])),
		mediaAndDataIntegrityErrors:        *(*uint64)(unsafe.Pointer(&data[160])),
		numberOfErrorInformationLogEntries: *(*uint64)(unsafe.Pointer(&data[176])),
	}
}

func nvmeFindDevices() ([]string, error) {
	devices, err := ioutil.ReadDir(nvmeSysPath)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, dev := range devices {
		files = append(files, filepath.Join("/dev", dev.Name()))
	}
	return files, nil
}

func nvmeIoctlAdminCmd(fd, arg, size uintptr) error {
	op := uintptr(0xc0004e41) | (size << 16)
	_, _, ep := syscall.Syscall(syscall.SYS_IOCTL, fd, op, arg)
	if ep != 0 {
		return syscall.Errno(ep)
	}
	return nil
}
