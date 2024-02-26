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
	"os"
	"strings"
	"testing"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

type testXfrmCollector struct {
	xc Collector
}

func (c testXfrmCollector) Collect(ch chan<- prometheus.Metric) {
	c.xc.Update(ch)
}

func (c testXfrmCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func TestXfrmStats(t *testing.T) {
	testcase := `# HELP node_xfrm_acquire_error_packets_total State hasn’t been fully acquired before use
	# TYPE node_xfrm_acquire_error_packets_total counter
	node_xfrm_acquire_error_packets_total 24532
	# HELP node_xfrm_fwd_hdr_error_packets_total Forward routing of a packet is not allowed
	# TYPE node_xfrm_fwd_hdr_error_packets_total counter
	node_xfrm_fwd_hdr_error_packets_total 6654
	# HELP node_xfrm_in_buffer_error_packets_total No buffer is left
	# TYPE node_xfrm_in_buffer_error_packets_total counter
	node_xfrm_in_buffer_error_packets_total 2
	# HELP node_xfrm_in_error_packets_total All errors not matched by other
	# TYPE node_xfrm_in_error_packets_total counter
	node_xfrm_in_error_packets_total 1
	# HELP node_xfrm_in_hdr_error_packets_total Header error
	# TYPE node_xfrm_in_hdr_error_packets_total counter
	node_xfrm_in_hdr_error_packets_total 4
	# HELP node_xfrm_in_no_pols_packets_total No policy is found for states e.g. Inbound SAs are correct but no SP is found
	# TYPE node_xfrm_in_no_pols_packets_total counter
	node_xfrm_in_no_pols_packets_total 65432
	# HELP node_xfrm_in_no_states_packets_total No state is found i.e. Either inbound SPI, address, or IPsec protocol at SA is wrong
	# TYPE node_xfrm_in_no_states_packets_total counter
	node_xfrm_in_no_states_packets_total 3
	# HELP node_xfrm_in_pol_block_packets_total Policy discards
	# TYPE node_xfrm_in_pol_block_packets_total counter
	node_xfrm_in_pol_block_packets_total 100
	# HELP node_xfrm_in_pol_error_packets_total Policy error
	# TYPE node_xfrm_in_pol_error_packets_total counter
	node_xfrm_in_pol_error_packets_total 10000
	# HELP node_xfrm_in_state_expired_packets_total State is expired
	# TYPE node_xfrm_in_state_expired_packets_total counter
	node_xfrm_in_state_expired_packets_total 7
	# HELP node_xfrm_in_state_invalid_packets_total State is invalid
	# TYPE node_xfrm_in_state_invalid_packets_total counter
	node_xfrm_in_state_invalid_packets_total 55555
	# HELP node_xfrm_in_state_mismatch_packets_total State has mismatch option e.g. UDP encapsulation type is mismatch
	# TYPE node_xfrm_in_state_mismatch_packets_total counter
	node_xfrm_in_state_mismatch_packets_total 23451
	# HELP node_xfrm_in_state_mode_error_packets_total Transformation mode specific error
	# TYPE node_xfrm_in_state_mode_error_packets_total counter
	node_xfrm_in_state_mode_error_packets_total 100
	# HELP node_xfrm_in_state_proto_error_packets_total Transformation protocol specific error e.g. SA key is wrong
	# TYPE node_xfrm_in_state_proto_error_packets_total counter
	node_xfrm_in_state_proto_error_packets_total 40
	# HELP node_xfrm_in_state_seq_error_packets_total Sequence error i.e. Sequence number is out of window
	# TYPE node_xfrm_in_state_seq_error_packets_total counter
	node_xfrm_in_state_seq_error_packets_total 6000
	# HELP node_xfrm_in_tmpl_mismatch_packets_total No matching template for states e.g. Inbound SAs are correct but SP rule is wrong
	# TYPE node_xfrm_in_tmpl_mismatch_packets_total counter
	node_xfrm_in_tmpl_mismatch_packets_total 51
	# HELP node_xfrm_out_bundle_check_error_packets_total Bundle check error
	# TYPE node_xfrm_out_bundle_check_error_packets_total counter
	node_xfrm_out_bundle_check_error_packets_total 555
	# HELP node_xfrm_out_bundle_gen_error_packets_total Bundle generation error
	# TYPE node_xfrm_out_bundle_gen_error_packets_total counter
	node_xfrm_out_bundle_gen_error_packets_total 43321
	# HELP node_xfrm_out_error_packets_total All errors which is not matched others
	# TYPE node_xfrm_out_error_packets_total counter
	node_xfrm_out_error_packets_total 1e+06
	# HELP node_xfrm_out_no_states_packets_total No state is found
	# TYPE node_xfrm_out_no_states_packets_total counter
	node_xfrm_out_no_states_packets_total 869
	# HELP node_xfrm_out_pol_block_packets_total Policy discards
	# TYPE node_xfrm_out_pol_block_packets_total counter
	node_xfrm_out_pol_block_packets_total 43456
	# HELP node_xfrm_out_pol_dead_packets_total Policy is dead
	# TYPE node_xfrm_out_pol_dead_packets_total counter
	node_xfrm_out_pol_dead_packets_total 7656
	# HELP node_xfrm_out_pol_error_packets_total Policy error
	# TYPE node_xfrm_out_pol_error_packets_total counter
	node_xfrm_out_pol_error_packets_total 1454
	# HELP node_xfrm_out_state_expired_packets_total State is expired
	# TYPE node_xfrm_out_state_expired_packets_total counter
	node_xfrm_out_state_expired_packets_total 565
	# HELP node_xfrm_out_state_invalid_packets_total State is invalid, perhaps expired
	# TYPE node_xfrm_out_state_invalid_packets_total counter
	node_xfrm_out_state_invalid_packets_total 28765
	# HELP node_xfrm_out_state_mode_error_packets_total Transformation mode specific error
	# TYPE node_xfrm_out_state_mode_error_packets_total counter
	node_xfrm_out_state_mode_error_packets_total 8
	# HELP node_xfrm_out_state_proto_error_packets_total Transformation protocol specific error
	# TYPE node_xfrm_out_state_proto_error_packets_total counter
	node_xfrm_out_state_proto_error_packets_total 4542
	# HELP node_xfrm_out_state_seq_error_packets_total Sequence error i.e. Sequence number overflow
	# TYPE node_xfrm_out_state_seq_error_packets_total counter
	node_xfrm_out_state_seq_error_packets_total 543
	`
	*procPath = "fixtures/proc"

	logger := log.NewLogfmtLogger(os.Stderr)
	c, err := NewXfrmCollector(logger)
	if err != nil {
		t.Fatal(err)
	}
	reg := prometheus.NewRegistry()
	reg.MustRegister(&testXfrmCollector{xc: c})

	sink := make(chan prometheus.Metric)
	go func() {
		err = c.Update(sink)
		if err != nil {
			panic(fmt.Errorf("failed to update collector: %s", err))
		}
		close(sink)
	}()

	err = testutil.GatherAndCompare(reg, strings.NewReader(testcase))
	if err != nil {
		t.Fatal(err)
	}
}
