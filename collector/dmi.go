// Copyright 2021 The Prometheus Authors
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

//go:build linux && !nodmi
// +build linux,!nodmi

package collector

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/sysfs"
)

type dmiCollector struct {
	infoDesc *prometheus.Desc
	values   []string
}

func init() {
	registerCollector("dmi", defaultEnabled, NewDMICollector)
}

// NewDMICollector returns a new Collector exposing DMI information.
func NewDMICollector(logger log.Logger) (Collector, error) {
	fs, err := sysfs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}

	dmi, err := fs.DMIClass()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			level.Debug(logger).Log("msg", "Platform does not support Desktop Management Interface (DMI) information", "err", err)
			dmi = &sysfs.DMIClass{}
		} else {
			return nil, fmt.Errorf("failed to read Desktop Management Interface (DMI) information: %w", err)
		}
	}

	var labels, values []string
	for label, value := range map[string]*string{
		"bios_date":         dmi.BiosDate,
		"bios_release":      dmi.BiosRelease,
		"bios_vendor":       dmi.BiosVendor,
		"bios_version":      dmi.BiosVersion,
		"board_asset_tag":   dmi.BoardAssetTag,
		"board_name":        dmi.BoardName,
		"board_serial":      dmi.BoardSerial,
		"board_vendor":      dmi.BoardVendor,
		"board_version":     dmi.BoardVersion,
		"chassis_asset_tag": dmi.ChassisAssetTag,
		"chassis_serial":    dmi.ChassisSerial,
		"chassis_vendor":    dmi.ChassisVendor,
		"chassis_version":   dmi.ChassisVersion,
		"product_family":    dmi.ProductFamily,
		"product_name":      dmi.ProductName,
		"product_serial":    dmi.ProductSerial,
		"product_sku":       dmi.ProductSKU,
		"product_uuid":      dmi.ProductUUID,
		"product_version":   dmi.ProductVersion,
		"system_vendor":     dmi.SystemVendor,
	} {
		if value != nil {
			labels = append(labels, label)
			values = append(values, strings.ToValidUTF8(*value, "�"))
		}
	}

	// Construct DMI metric only once since it will not change until the next reboot.
	return &dmiCollector{
		infoDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "dmi", "info"),
			"A metric with a constant '1' value labeled by bios_date, bios_release, bios_vendor, bios_version, "+
				"board_asset_tag, board_name, board_serial, board_vendor, board_version, chassis_asset_tag, "+
				"chassis_serial, chassis_vendor, chassis_version, product_family, product_name, product_serial, "+
				"product_sku, product_uuid, product_version, system_vendor if provided by DMI.",
			labels, nil,
		),
		values: values,
	}, nil
}

func (c *dmiCollector) Update(ch chan<- prometheus.Metric) error {
	if len(c.values) == 0 {
		return ErrNoData
	}
	ch <- prometheus.MustNewConstMetric(c.infoDesc, prometheus.GaugeValue, 1.0, c.values...)
	return nil
}
