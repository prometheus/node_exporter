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

//go:build !noxfrm
// +build !noxfrm

package collector

import (
	"fmt"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

type xfrmCollector struct {
	fs     procfs.FS
	logger log.Logger
}

func init() {
	registerCollector("xfrm", defaultDisabled, NewXfrmCollector)
}

// NewXfrmCollector returns a new Collector exposing XFRM stats.
func NewXfrmCollector(logger log.Logger) (Collector, error) {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}

	return &xfrmCollector{
		fs:     fs,
		logger: logger,
	}, nil
}

var (
	xfrmInErrorDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "xfrm", "in_error_packets_total"),
		"All errors not matched by other",
		nil, nil,
	)
	xfrmInBufferErrorDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "xfrm", "in_buffer_error_packets_total"),
		"No buffer is left",
		nil, nil,
	)
	xfrmInHdrErrorDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "xfrm", "in_hdr_error_packets_total"),
		"Header error",
		nil, nil,
	)
	xfrmInNoStatesDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "xfrm", "in_no_states_packets_total"),
		"No state is found i.e. Either inbound SPI, address, or IPsec protocol at SA is wrong",
		nil, nil,
	)
	xfrmInStateProtoErrorDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "xfrm", "in_state_proto_error_packets_total"),
		"Transformation protocol specific error e.g. SA key is wrong",
		nil, nil,
	)
	xfrmInStateModeErrorDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "xfrm", "in_state_mode_error_packets_total"),
		"Transformation mode specific error",
		nil, nil,
	)
	xfrmInStateSeqErrorDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "xfrm", "in_state_seq_error_packets_total"),
		"Sequence error i.e. Sequence number is out of window",
		nil, nil,
	)
	xfrmInStateExpiredDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "xfrm", "in_state_expired_packets_total"),
		"State is expired",
		nil, nil,
	)
	xfrmInStateMismatchDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "xfrm", "in_state_mismatch_packets_total"),
		"State has mismatch option e.g. UDP encapsulation type is mismatch",
		nil, nil,
	)
	xfrmInStateInvalidDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "xfrm", "in_state_invalid_packets_total"),
		"State is invalid",
		nil, nil,
	)
	xfrmInTmplMismatchDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "xfrm", "in_tmpl_mismatch_packets_total"),
		"No matching template for states e.g. Inbound SAs are correct but SP rule is wrong",
		nil, nil,
	)
	xfrmInNoPolsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "xfrm", "in_no_pols_packets_total"),
		"No policy is found for states e.g. Inbound SAs are correct but no SP is found",
		nil, nil,
	)
	xfrmInPolBlockDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "xfrm", "in_pol_block_packets_total"),
		"Policy discards",
		nil, nil,
	)
	xfrmInPolErrorDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "xfrm", "in_pol_error_packets_total"),
		"Policy error",
		nil, nil,
	)
	xfrmOutErrorDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "xfrm", "out_error_packets_total"),
		"All errors which is not matched others",
		nil, nil,
	)
	xfrmOutBundleGenErrorDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "xfrm", "out_bundle_gen_error_packets_total"),
		"Bundle generation error",
		nil, nil,
	)
	xfrmOutBundleCheckErrorDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "xfrm", "out_bundle_check_error_packets_total"),
		"Bundle check error",
		nil, nil,
	)
	xfrmOutNoStatesDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "xfrm", "out_no_states_packets_total"),
		"No state is found",
		nil, nil,
	)
	xfrmOutStateProtoErrorDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "xfrm", "out_state_proto_error_packets_total"),
		"Transformation protocol specific error",
		nil, nil,
	)
	xfrmOutStateModeErrorDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "xfrm", "out_state_mode_error_packets_total"),
		"Transformation mode specific error",
		nil, nil,
	)
	xfrmOutStateSeqErrorDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "xfrm", "out_state_seq_error_packets_total"),
		"Sequence error i.e. Sequence number overflow",
		nil, nil,
	)
	xfrmOutStateExpiredDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "xfrm", "out_state_expired_packets_total"),
		"State is expired",
		nil, nil,
	)
	xfrmOutPolBlockDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "xfrm", "out_pol_block_packets_total"),
		"Policy discards",
		nil, nil,
	)
	xfrmOutPolDeadDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "xfrm", "out_pol_dead_packets_total"),
		"Policy is dead",
		nil, nil,
	)
	xfrmOutPolErrorDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "xfrm", "out_pol_error_packets_total"),
		"Policy error",
		nil, nil,
	)
	xfrmFwdHdrErrorDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "xfrm", "fwd_hdr_error_packets_total"),
		"Forward routing of a packet is not allowed",
		nil, nil,
	)
	xfrmOutStateInvalidDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "xfrm", "out_state_invalid_packets_total"),
		"State is invalid, perhaps expired",
		nil, nil,
	)
	xfrmAcquireErrorDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "xfrm", "acquire_error_packets_total"),
		"State hasn’t been fully acquired before use",
		nil, nil,
	)
)

func (c *xfrmCollector) Update(ch chan<- prometheus.Metric) error {
	stat, err := c.fs.NewXfrmStat()
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(xfrmInErrorDesc, prometheus.CounterValue, float64(stat.XfrmInError))
	ch <- prometheus.MustNewConstMetric(xfrmInBufferErrorDesc, prometheus.CounterValue, float64(stat.XfrmInBufferError))
	ch <- prometheus.MustNewConstMetric(xfrmInHdrErrorDesc, prometheus.CounterValue, float64(stat.XfrmInHdrError))
	ch <- prometheus.MustNewConstMetric(xfrmInNoStatesDesc, prometheus.CounterValue, float64(stat.XfrmInNoStates))
	ch <- prometheus.MustNewConstMetric(xfrmInStateProtoErrorDesc, prometheus.CounterValue, float64(stat.XfrmInStateProtoError))
	ch <- prometheus.MustNewConstMetric(xfrmInStateModeErrorDesc, prometheus.CounterValue, float64(stat.XfrmInStateModeError))
	ch <- prometheus.MustNewConstMetric(xfrmInStateSeqErrorDesc, prometheus.CounterValue, float64(stat.XfrmInStateSeqError))
	ch <- prometheus.MustNewConstMetric(xfrmInStateExpiredDesc, prometheus.CounterValue, float64(stat.XfrmInStateExpired))
	ch <- prometheus.MustNewConstMetric(xfrmInStateMismatchDesc, prometheus.CounterValue, float64(stat.XfrmInStateMismatch))
	ch <- prometheus.MustNewConstMetric(xfrmInStateInvalidDesc, prometheus.CounterValue, float64(stat.XfrmInStateInvalid))
	ch <- prometheus.MustNewConstMetric(xfrmInTmplMismatchDesc, prometheus.CounterValue, float64(stat.XfrmInTmplMismatch))
	ch <- prometheus.MustNewConstMetric(xfrmInNoPolsDesc, prometheus.CounterValue, float64(stat.XfrmInNoPols))
	ch <- prometheus.MustNewConstMetric(xfrmInPolBlockDesc, prometheus.CounterValue, float64(stat.XfrmInPolBlock))
	ch <- prometheus.MustNewConstMetric(xfrmInPolErrorDesc, prometheus.CounterValue, float64(stat.XfrmInPolError))
	ch <- prometheus.MustNewConstMetric(xfrmOutErrorDesc, prometheus.CounterValue, float64(stat.XfrmOutError))
	ch <- prometheus.MustNewConstMetric(xfrmOutBundleGenErrorDesc, prometheus.CounterValue, float64(stat.XfrmOutBundleGenError))
	ch <- prometheus.MustNewConstMetric(xfrmOutBundleCheckErrorDesc, prometheus.CounterValue, float64(stat.XfrmOutBundleCheckError))
	ch <- prometheus.MustNewConstMetric(xfrmOutNoStatesDesc, prometheus.CounterValue, float64(stat.XfrmOutNoStates))
	ch <- prometheus.MustNewConstMetric(xfrmOutStateProtoErrorDesc, prometheus.CounterValue, float64(stat.XfrmOutStateProtoError))
	ch <- prometheus.MustNewConstMetric(xfrmOutStateModeErrorDesc, prometheus.CounterValue, float64(stat.XfrmOutStateModeError))
	ch <- prometheus.MustNewConstMetric(xfrmOutStateSeqErrorDesc, prometheus.CounterValue, float64(stat.XfrmOutStateSeqError))
	ch <- prometheus.MustNewConstMetric(xfrmOutStateExpiredDesc, prometheus.CounterValue, float64(stat.XfrmOutStateExpired))
	ch <- prometheus.MustNewConstMetric(xfrmOutPolBlockDesc, prometheus.CounterValue, float64(stat.XfrmOutPolBlock))
	ch <- prometheus.MustNewConstMetric(xfrmOutPolDeadDesc, prometheus.CounterValue, float64(stat.XfrmOutPolDead))
	ch <- prometheus.MustNewConstMetric(xfrmOutPolErrorDesc, prometheus.CounterValue, float64(stat.XfrmOutPolError))
	ch <- prometheus.MustNewConstMetric(xfrmFwdHdrErrorDesc, prometheus.CounterValue, float64(stat.XfrmFwdHdrError))
	ch <- prometheus.MustNewConstMetric(xfrmOutStateInvalidDesc, prometheus.CounterValue, float64(stat.XfrmOutStateInvalid))
	ch <- prometheus.MustNewConstMetric(xfrmAcquireErrorDesc, prometheus.CounterValue, float64(stat.XfrmAcquireError))

	return err
}
