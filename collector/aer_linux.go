// Copyright 2024 The Prometheus Authors
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
//
//go:build !nonetclass && linux
// +build !nonetclass,linux

package collector

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strconv"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/sysfs"
)

var (
	aerIgnoredDevices   = kingpin.Flag("collector.aer.ignored-devices", "Regexp of aer devices to ignore for aer collector.").Default("^$").String()
	aerCorrectableRxErr = prometheus.NewDesc(prometheus.BuildFQName(namespace, "aer", "correctable_rx_err"),
		"Count of correctable receiver errors",
		[]string{"interface"}, nil,
	)
	aerCorrectableBadTLP = prometheus.NewDesc(prometheus.BuildFQName(namespace, "aer", "correctable_bad_tlp"),
		"Count of correctable bad TLPs",
		[]string{"interface"}, nil,
	)
	aerCorrectableBadDLLP = prometheus.NewDesc(prometheus.BuildFQName(namespace, "aer", "correctable_bad_dllp"),
		"Count of correctable bad DLLPs",
		[]string{"interface"}, nil,
	)
	aerCorrectableRollover = prometheus.NewDesc(prometheus.BuildFQName(namespace, "aer", "correctable_rollover"),
		"Count of correctable rollovers",
		[]string{"interface"}, nil,
	)
	aerCorrectableTimeout = prometheus.NewDesc(prometheus.BuildFQName(namespace, "aer", "correctable_timeout"),
		"Count of correctable replay timer timeouts",
		[]string{"interface"}, nil,
	)
	aerCorrectableNonFatalErr = prometheus.NewDesc(prometheus.BuildFQName(namespace, "aer", "correctable_non_fatal_err"),
		"Count of correctable advisory non-fatal errors",
		[]string{"interface"}, nil,
	)
	aerCorrectableCorrIntErr = prometheus.NewDesc(prometheus.BuildFQName(namespace, "aer", "correctable_corr_int_err"),
		"Count of correctable corrected internal errors",
		[]string{"interface"}, nil,
	)
	aerCorrectableHeaderOF = prometheus.NewDesc(prometheus.BuildFQName(namespace, "aer", "correctable_header_of"),
		"Count of correctable header log Overflows",
		[]string{"interface"}, nil,
	)
	aerUncorrectableUndefined = prometheus.NewDesc(prometheus.BuildFQName(namespace, "aer", "uncorrectable_undefined"),
		"Count of uncorrectable undefined errors",
		[]string{"interface", "fatal"}, nil,
	)
	aerUncorrectableDLP = prometheus.NewDesc(prometheus.BuildFQName(namespace, "aer", "uncorrectable_dlp"),
		"Count of uncorrectable data link protocol errors",
		[]string{"interface", "fatal"}, nil,
	)
	aerUncorrectableSDES = prometheus.NewDesc(prometheus.BuildFQName(namespace, "aer", "uncorrectable_sdes"),
		"Count of uncorrectable surprise down errors",
		[]string{"interface", "fatal"}, nil,
	)
	aerUncorrectableTLP = prometheus.NewDesc(prometheus.BuildFQName(namespace, "aer", "uncorrectable_tlp"),
		"Count of uncorrectable poisoned TLPs",
		[]string{"interface", "fatal"}, nil,
	)
	aerUncorrectableFCP = prometheus.NewDesc(prometheus.BuildFQName(namespace, "aer", "uncorrectable_fcp"),
		"Count of uncorrectable flow control protocol errors",
		[]string{"interface", "fatal"}, nil,
	)
	aerUncorrectableCmpltTO = prometheus.NewDesc(prometheus.BuildFQName(namespace, "aer", "uncorrectable_cmplt_to"),
		"Count of uncorrectable completion timeouts",
		[]string{"interface", "fatal"}, nil,
	)
	aerUncorrectableCmpltAbrt = prometheus.NewDesc(prometheus.BuildFQName(namespace, "aer", "uncorrectable_cmplt_abrt"),
		"Count of uncorrectable completer aborts",
		[]string{"interface", "fatal"}, nil,
	)
	aerUncorrectableUnxCmplt = prometheus.NewDesc(prometheus.BuildFQName(namespace, "aer", "uncorrectable_unx_cmplt"),
		"Count of uncorrectable unexpected completion errors",
		[]string{"interface", "fatal"}, nil,
	)
	aerUncorrectableRxOF = prometheus.NewDesc(prometheus.BuildFQName(namespace, "aer", "uncorrectable_rx_of"),
		"Count of uncorrectable receiver overflows",
		[]string{"interface", "fatal"}, nil,
	)
	aerUncorrectableMalfTLP = prometheus.NewDesc(prometheus.BuildFQName(namespace, "aer", "uncorrectable_malf_tlp"),
		"Count of uncorrectable malformed TLPs",
		[]string{"interface", "fatal"}, nil,
	)
	aerUncorrectableECRC = prometheus.NewDesc(prometheus.BuildFQName(namespace, "aer", "uncorrectable_ecrc"),
		"Count of uncorrectable ECRCs",
		[]string{"interface", "fatal"}, nil,
	)
	aerUncorrectableUnsupReq = prometheus.NewDesc(prometheus.BuildFQName(namespace, "aer", "uncorrectable_unsup_req"),
		"Count of uncorrectable unsupported requests",
		[]string{"interface", "fatal"}, nil,
	)
	aerUncorrectableACSViol = prometheus.NewDesc(prometheus.BuildFQName(namespace, "aer", "uncorrectable_acs_viol"),
		"Count of uncorrectable ACS violations",
		[]string{"interface", "fatal"}, nil,
	)
	aerUncorrectableUncorrIntErr = prometheus.NewDesc(prometheus.BuildFQName(namespace, "aer", "uncorrectable_uncorr_int_err"),
		"Count of uncorrectable uncorrectable internal errors",
		[]string{"interface", "fatal"}, nil,
	)
	aerUncorrectableBlockedTLP = prometheus.NewDesc(prometheus.BuildFQName(namespace, "aer", "uncorrectable_blocked_tlp"),
		"Count of uncorrectable MC blocked TLPs",
		[]string{"interface", "fatal"}, nil,
	)
	aerUncorrectableAtomicOpBlocked = prometheus.NewDesc(prometheus.BuildFQName(namespace, "aer", "uncorrectable_atomic_op_blocked"),
		"Count of uncorrectable AtomicOp egress blocked errors",
		[]string{"interface", "fatal"}, nil,
	)
	aerUncorrectableTLPBlockedErr = prometheus.NewDesc(prometheus.BuildFQName(namespace, "aer", "uncorrectable_tlp_blocked_err"),
		"Count of uncorrectable TLP prefix blocked errors",
		[]string{"interface", "fatal"}, nil,
	)
	aerUncorrectablePoisonTLPBlocked = prometheus.NewDesc(prometheus.BuildFQName(namespace, "aer", "uncorrectable_poison_tlp_blocked"),
		"Count of uncorrectable poison TLP prefix blocked errors",
		[]string{"interface", "fatal"}, nil,
	)
)

type aerCollector struct {
	fs                    sysfs.FS
	ignoredDevicesPattern *regexp.Regexp
	logger                *slog.Logger
}

func init() {
	registerCollector("aer", defaultDisabled, NewAerCollector)
}

// NewAerCollector returns a new Collector exposing aer stats.
func NewAerCollector(logger *slog.Logger) (Collector, error) {
	return makeAerCollector(logger)
}

func makeAerCollector(logger *slog.Logger) (*aerCollector, error) {
	fs, err := sysfs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}

	if *aerIgnoredDevices != "" {
		logger.Info("Parsed flag --collector.aer.ignored-devices", "flag", *aerIgnoredDevices)
	}
	pattern := regexp.MustCompile(*aerIgnoredDevices)
	return &aerCollector{
		fs:                    fs,
		ignoredDevicesPattern: pattern,
		logger:                logger,
	}, nil

}

func (c *aerCollector) Update(ch chan<- prometheus.Metric) error {
	counters, err := c.fs.AerCounters()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) || errors.Is(err, os.ErrPermission) {
			c.logger.Debug("Could not read netclass file", "err", err)
			return ErrNoData
		}
		return fmt.Errorf("could not get net class info: %w", err)
	}

	for deviceName, deviceCounters := range counters {
		if c.ignoredDevicesPattern.MatchString(deviceName) {
			continue
		}

		c.updateCorrectableCntrs(ch, deviceName, deviceCounters.Correctable)
		c.updateUncorrectableCntrs(ch, deviceName, deviceCounters.Fatal, true)
		c.updateUncorrectableCntrs(ch, deviceName, deviceCounters.NonFatal, false)

	}
	return nil
}

func (c *aerCollector) updateCorrectableCntrs(ch chan<- prometheus.Metric, deviceName string, counters sysfs.CorrectableAerCounters) {
	ch <- prometheus.MustNewConstMetric(aerCorrectableRxErr, prometheus.CounterValue,
		float64(counters.RxErr), deviceName)
	ch <- prometheus.MustNewConstMetric(aerCorrectableBadTLP, prometheus.CounterValue,
		float64(counters.BadTLP), deviceName)
	ch <- prometheus.MustNewConstMetric(aerCorrectableBadDLLP, prometheus.CounterValue,
		float64(counters.BadDLLP), deviceName)
	ch <- prometheus.MustNewConstMetric(aerCorrectableRollover, prometheus.CounterValue,
		float64(counters.Rollover), deviceName)
	ch <- prometheus.MustNewConstMetric(aerCorrectableTimeout, prometheus.CounterValue,
		float64(counters.Timeout), deviceName)
	ch <- prometheus.MustNewConstMetric(aerCorrectableNonFatalErr, prometheus.CounterValue,
		float64(counters.NonFatalErr), deviceName)
	ch <- prometheus.MustNewConstMetric(aerCorrectableCorrIntErr, prometheus.CounterValue,
		float64(counters.CorrIntErr), deviceName)
	ch <- prometheus.MustNewConstMetric(aerCorrectableHeaderOF, prometheus.CounterValue,
		float64(counters.HeaderOF), deviceName)
}

func (c *aerCollector) updateUncorrectableCntrs(ch chan<- prometheus.Metric, deviceName string, counters sysfs.UncorrectableAerCounters, fatal bool) {
	ch <- prometheus.MustNewConstMetric(aerUncorrectableUndefined, prometheus.CounterValue,
		float64(counters.Undefined), deviceName, strconv.FormatBool(fatal))
	ch <- prometheus.MustNewConstMetric(aerUncorrectableDLP, prometheus.CounterValue,
		float64(counters.DLP), deviceName, strconv.FormatBool(fatal))
	ch <- prometheus.MustNewConstMetric(aerUncorrectableSDES, prometheus.CounterValue,
		float64(counters.SDES), deviceName, strconv.FormatBool(fatal))
	ch <- prometheus.MustNewConstMetric(aerUncorrectableTLP, prometheus.CounterValue,
		float64(counters.TLP), deviceName, strconv.FormatBool(fatal))
	ch <- prometheus.MustNewConstMetric(aerUncorrectableFCP, prometheus.CounterValue,
		float64(counters.FCP), deviceName, strconv.FormatBool(fatal))
	ch <- prometheus.MustNewConstMetric(aerUncorrectableCmpltTO, prometheus.CounterValue,
		float64(counters.CmpltTO), deviceName, strconv.FormatBool(fatal))
	ch <- prometheus.MustNewConstMetric(aerUncorrectableCmpltAbrt, prometheus.CounterValue,
		float64(counters.CmpltAbrt), deviceName, strconv.FormatBool(fatal))
	ch <- prometheus.MustNewConstMetric(aerUncorrectableUnxCmplt, prometheus.CounterValue,
		float64(counters.UnxCmplt), deviceName, strconv.FormatBool(fatal))
	ch <- prometheus.MustNewConstMetric(aerUncorrectableRxOF, prometheus.CounterValue,
		float64(counters.RxOF), deviceName, strconv.FormatBool(fatal))
	ch <- prometheus.MustNewConstMetric(aerUncorrectableMalfTLP, prometheus.CounterValue,
		float64(counters.MalfTLP), deviceName, strconv.FormatBool(fatal))
	ch <- prometheus.MustNewConstMetric(aerUncorrectableECRC, prometheus.CounterValue,
		float64(counters.ECRC), deviceName, strconv.FormatBool(fatal))
	ch <- prometheus.MustNewConstMetric(aerUncorrectableUnsupReq, prometheus.CounterValue,
		float64(counters.UnsupReq), deviceName, strconv.FormatBool(fatal))
	ch <- prometheus.MustNewConstMetric(aerUncorrectableACSViol, prometheus.CounterValue,
		float64(counters.ACSViol), deviceName, strconv.FormatBool(fatal))
	ch <- prometheus.MustNewConstMetric(aerUncorrectableUncorrIntErr, prometheus.CounterValue,
		float64(counters.UncorrIntErr), deviceName, strconv.FormatBool(fatal))
	ch <- prometheus.MustNewConstMetric(aerUncorrectableBlockedTLP, prometheus.CounterValue,
		float64(counters.BlockedTLP), deviceName, strconv.FormatBool(fatal))
	ch <- prometheus.MustNewConstMetric(aerUncorrectableAtomicOpBlocked, prometheus.CounterValue,
		float64(counters.AtomicOpBlocked), deviceName, strconv.FormatBool(fatal))
	ch <- prometheus.MustNewConstMetric(aerUncorrectableTLPBlockedErr, prometheus.CounterValue,
		float64(counters.TLPBlockedErr), deviceName, strconv.FormatBool(fatal))
	ch <- prometheus.MustNewConstMetric(aerUncorrectablePoisonTLPBlocked, prometheus.CounterValue,
		float64(counters.PoisonTLPBlocked), deviceName, strconv.FormatBool(fatal))
}
