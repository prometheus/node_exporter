// Copyright 2023 The Prometheus Authors
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

//go:build !noiio
// +build !noiio

package collector

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sys/unix"
)

var (
	iioInvalidMetricChars = regexp.MustCompile("[^a-z0-9:_]")
	iioFilenameFormat     = regexp.MustCompile(`^(?P<type>[^0-9]+)(?P<id>[0-9]*)?_(?P<property>.+)$`)
	iioLabelDesc          = []string{"chip", "sensor"}
	iioChipNameLabelDesc  = []string{"chip", "chip_name"}
	iioSensorTypes        = []string{
		"in_temp", "in_pressure",
	}
)

func init() {
	registerCollector("iio", defaultDisabled, NewIIOCollector)
}

type iioCollector struct {
	logger log.Logger
}

// NewIIOCollector returns a new Collector exposing /sys/bus/iio/devices.
func NewIIOCollector(logger log.Logger) (Collector, error) {
	return &iioCollector{logger}, nil
}

func iioAddValueFile(data map[string]map[string]string, sensor string, prop string, file string) {
	raw, err := iioSysReadFile(file)
	if err != nil {
		return
	}
	value := strings.Trim(string(raw), "\n")

	if _, ok := data[sensor]; !ok {
		data[sensor] = make(map[string]string)
	}

	data[sensor][prop] = value
}

// sysReadFile is a simplified os.ReadFile that invokes syscall.Read directly.
func iioSysReadFile(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// On some machines, iio drivers are broken and return EAGAIN.  This causes
	// Go's os.ReadFile implementation to poll forever.
	//
	// Since we either want to read data or bail immediately, do the simplest
	// possible read using system call directly.
	b := make([]byte, 128)
	n, err := unix.Read(int(f.Fd()), b)
	if err != nil {
		return nil, err
	}

	return b[:n], nil
}

// iioExplodeSensorFilename splits a sensor name into <type><num>_<property>.
func iioExplodeSensorFilename(filename string) (ok bool, sensorType string, sensorNum int, sensorProperty string) {
	matches := iioFilenameFormat.FindStringSubmatch(filename)
	if len(matches) == 0 {
		return false, sensorType, sensorNum, sensorProperty
	}
	for i, match := range iioFilenameFormat.SubexpNames() {
		if i >= len(matches) {
			return true, sensorType, sensorNum, sensorProperty
		}
		if match == "type" {
			sensorType = matches[i]
		}
		if match == "property" {
			sensorProperty = matches[i]
		}
		if match == "id" && len(matches[i]) > 0 {
			if num, err := strconv.Atoi(matches[i]); err == nil {
				sensorNum = num
			} else {
				return false, sensorType, sensorNum, sensorProperty
			}
		}
	}
	return true, sensorType, sensorNum, sensorProperty
}

func iioCollectSensorData(dir string, data map[string]map[string]string) error {
	sensorFiles, dirError := os.ReadDir(dir)
	if dirError != nil {
		return dirError
	}
	for _, file := range sensorFiles {
		filename := file.Name()
		ok, sensorType, sensorNum, sensorProperty := iioExplodeSensorFilename(filename)
		if !ok {
			continue
		}

		for _, t := range iioSensorTypes {
			if t == sensorType {
				iioAddValueFile(data, sensorType+strconv.Itoa(sensorNum), sensorProperty, filepath.Join(dir, file.Name()))
				break
			}
		}
	}
	return nil
}

func (c *iioCollector) updateIIO(ch chan<- prometheus.Metric, dir string) error {
	iioName, err := c.iioName(dir)
	if err != nil {
		return err
	}

	data := make(map[string]map[string]string)
	err = iioCollectSensorData(dir, data)
	if err != nil {
		return err
	}

	iioChipName, err := c.iioHumanReadableChipName(dir)
	if err == nil {
		// sensor chip metadata
		desc := prometheus.NewDesc(
			"node_iio_chip_names",
			"Annotation metric for human-readable chip names",
			iioChipNameLabelDesc,
			nil,
		)

		ch <- prometheus.MustNewConstMetric(
			desc,
			prometheus.GaugeValue,
			1.0,
			iioName,
			iioChipName,
		)
	}

	// Format all sensors.
	for sensor, sensorData := range data {

		_, sensorType, _, _ := explodeSensorFilename(sensor)

		labels := []string{iioName, sensor}
		if labelText, ok := sensorData["label"]; ok {
			label := strings.ToValidUTF8(labelText, "�")
			desc := prometheus.NewDesc("node_iio_sensor_label", "Label for given chip and sensor",
				[]string{"chip", "sensor", "label"}, nil)
			ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, 1.0, iioName, sensor, label)
		}

		prefix := "node_iio_" + sensorType

		for element, value := range sensorData {

			if element == "label" {
				continue
			}

			name := prefix
			if element == "input" {
				// input is actually the value
				if _, ok := sensorData[""]; ok {
					name = name + "_input"
				}
			} else if element != "" {
				name = name + "_" + cleanMetricName(element)
			}
			parsedValue, err := strconv.ParseFloat(value, 64)
			if err != nil {
				continue
			}

			if sensorType == "in_temp" && element == "input" {
				desc := prometheus.NewDesc(name+"_celsius", "Industrial I/O sensor for temperature ("+element+")", iioLabelDesc, nil)
				ch <- prometheus.MustNewConstMetric(
					desc, prometheus.GaugeValue, parsedValue*0.001, labels...)
				continue
			}

			if sensorType == "in_pressure" && element == "input" {
				desc := prometheus.NewDesc(name+"_pascal", "Industrial I/O sensor for pressure ("+element+")", iioLabelDesc, nil)
				ch <- prometheus.MustNewConstMetric(
					desc, prometheus.GaugeValue, parsedValue*1000.0, labels...)
				continue
			}

			// fallback, just dump the metric as is

			desc := prometheus.NewDesc(name, "Industrial I/O sensor "+sensorType+" element "+element, iioLabelDesc, nil)
			ch <- prometheus.MustNewConstMetric(
				desc, prometheus.GaugeValue, parsedValue, labels...)
		}
	}
	return nil
}

func (c *iioCollector) iioName(dir string) (string, error) {
	// generate a name for a sensor path

	// sensor numbering depends on the order of linux module loading and
	// is thus unstable.
	// However the path of the device has to be stable:
	// - /sys/devices/<bus>/<device>
	// Some hardware monitors have a "name" file that exports a human
	// readable name that can be used.

	// human readable names would be bat0 or coretemp, while a path string
	// could be platform_applesmc.768

	// preference 1: construct a name based on device name, always unique

	devicePath, devErr := filepath.EvalSymlinks(filepath.Join(dir, "device"))
	if devErr == nil {
		devPathPrefix, devName := filepath.Split(devicePath)
		_, devType := filepath.Split(strings.TrimRight(devPathPrefix, "/"))

		cleanDevName := cleanMetricName(devName)
		cleanDevType := cleanMetricName(devType)

		if cleanDevType != "" && cleanDevName != "" {
			return cleanDevType + "_" + cleanDevName, nil
		}

		if cleanDevName != "" {
			return cleanDevName, nil
		}
	}

	// preference 2: is there a name file
	sysnameRaw, nameErr := os.ReadFile(filepath.Join(dir, "name"))
	if nameErr == nil && string(sysnameRaw) != "" {
		cleanName := cleanMetricName(string(sysnameRaw))
		if cleanName != "" {
			return cleanName, nil
		}
	}

	// it looks bad, name and device don't provide enough information
	// return a iio:device[0-9]* name

	realDir, err := filepath.EvalSymlinks(dir)
	if err != nil {
		return "", err
	}

	// take the last path element, this will be iio:deviceX
	_, name := filepath.Split(realDir)
	cleanName := cleanMetricName(name)
	if cleanName != "" {
		return cleanName, nil
	}
	return "", errors.New("Could not derive a monitoring name for " + dir)
}

// iioHumanReadableChipName is similar to the methods in iioName, but with
// different precedences -- we can allow duplicates here.
func (c *iioCollector) iioHumanReadableChipName(dir string) (string, error) {
	sysnameRaw, nameErr := os.ReadFile(filepath.Join(dir, "name"))
	if nameErr != nil {
		return "", nameErr
	}

	if string(sysnameRaw) != "" {
		cleanName := cleanMetricName(string(sysnameRaw))
		if cleanName != "" {
			return cleanName, nil
		}
	}

	return "", errors.New("Could not derive a human-readable chip type for " + dir)
}

func (c *iioCollector) Update(ch chan<- prometheus.Metric) error {
	// Step 1: scan /sys/bus/iio/devices, resolve all symlinks and call
	//         updatesIIO for each folder

	iioPathName := filepath.Join(sysFilePath("bus"), "iio", "devices")

	iioDeviceFiles, err := os.ReadDir(iioPathName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			level.Debug(c.logger).Log("msg", "iio collector metrics are not available for this system")
			return ErrNoData
		}

		return err
	}
	for _, iioDeviceDir := range iioDeviceFiles {
		iioDevicePathName := filepath.Join(iioPathName, iioDeviceDir.Name())
		fileInfo, _ := os.Lstat(iioDevicePathName)

		if fileInfo.Mode()&os.ModeSymlink > 0 {
			fileInfo, err = os.Stat(iioDevicePathName)
			if err != nil {
				continue
			}
		}

		if !fileInfo.IsDir() {
			continue
		}

		if lastErr := c.updateIIO(ch, iioDevicePathName); lastErr != nil {
			err = lastErr
		}
	}
	return err
}
