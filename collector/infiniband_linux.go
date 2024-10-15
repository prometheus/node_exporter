// Copyright 2017-2019 The Prometheus Authors
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

//go:build !noinfiniband
// +build !noinfiniband

package collector

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/sysfs"
)

type infinibandCollector struct {
	fs          sysfs.FS
	metricDescs map[string]*prometheus.Desc
	logger      *slog.Logger
	subsystem   string
}

func init() {
	registerCollector("infiniband", defaultEnabled, NewInfiniBandCollector)
}

// NewInfiniBandCollector returns a new Collector exposing InfiniBand stats.
func NewInfiniBandCollector(logger *slog.Logger) (Collector, error) {
	var i infinibandCollector
	var err error

	i.fs, err = sysfs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}
	i.logger = logger

	// Detailed description for all metrics.
	descriptions := map[string]string{
		"legacy_multicast_packets_received_total":    "Number of multicast packets received",
		"legacy_multicast_packets_transmitted_total": "Number of multicast packets transmitted",
		"legacy_data_received_bytes_total":           "Number of data octets received on all links",
		"legacy_packets_received_total":              "Number of data packets received on all links",
		"legacy_unicast_packets_received_total":      "Number of unicast packets received",
		"legacy_unicast_packets_transmitted_total":   "Number of unicast packets transmitted",
		"legacy_data_transmitted_bytes_total":        "Number of data octets transmitted on all links",
		"legacy_packets_transmitted_total":           "Number of data packets received on all links",
		"excessive_buffer_overrun_errors_total":      "Number of times that OverrunErrors consecutive flow control update periods occurred, each having at least one overrun error.",
		"link_downed_total":                          "Number of times the link failed to recover from an error state and went down",
		"link_error_recovery_total":                  "Number of times the link successfully recovered from an error state",
		"local_link_integrity_errors_total":          "Number of times that the count of local physical errors exceeded the threshold specified by LocalPhyErrors.",
		"multicast_packets_received_total":           "Number of multicast packets received (including errors)",
		"multicast_packets_transmitted_total":        "Number of multicast packets transmitted (including errors)",
		"physical_state_id":                          "Physical state of the InfiniBand port (0: no change, 1: sleep, 2: polling, 3: disable, 4: shift, 5: link up, 6: link error recover, 7: phytest)",
		"port_constraint_errors_received_total":      "Number of packets received on the switch physical port that are discarded",
		"port_constraint_errors_transmitted_total":   "Number of packets not transmitted from the switch physical port",
		"port_data_received_bytes_total":             "Number of data octets received on all links",
		"port_data_transmitted_bytes_total":          "Number of data octets transmitted on all links",
		"port_discards_received_total":               "Number of inbound packets discarded by the port because the port is down or congested",
		"port_discards_transmitted_total":            "Number of outbound packets discarded by the port because the port is down or congested",
		"port_errors_received_total":                 "Number of packets containing an error that were received on this port",
		"port_packets_received_total":                "Number of packets received on all VLs by this port (including errors)",
		"port_packets_transmitted_total":             "Number of packets transmitted on all VLs from this port (including errors)",
		"port_transmit_wait_total":                   "Number of ticks during which the port had data to transmit but no data was sent during the entire tick",
		"rate_bytes_per_second":                      "Maximum signal transfer rate",
		"state_id":                                   "State of the InfiniBand port (0: no change, 1: down, 2: init, 3: armed, 4: active, 5: act defer)",
		"unicast_packets_received_total":             "Number of unicast packets received (including errors)",
		"unicast_packets_transmitted_total":          "Number of unicast packets transmitted (including errors)",
		"port_receive_remote_physical_errors_total":  "Number of packets marked with the EBP (End of Bad Packet) delimiter received on the port.",
		"port_receive_switch_relay_errors_total":     "Number of packets that could not be forwarded by the switch.",
		"symbol_error_total":                         "Number of minor link errors detected on one or more physical lanes.",
		"vl15_dropped_total":                         "Number of incoming VL15 packets dropped due to resource limitations.",

		"active_ahs":                                 "Number of active ahs.",
		"active_cqs":                                 "Number of active_cqs.",
		"active_mrs":                                 "Number of active_mrs.",
		"active_mws":                                 "Number of active_mws.",
		"active_pds":                                 "Number of active_pds.",
		"active_qps":                                 "Number of active_qps.",
		"active_rc_qps":                              "Number of active_rc_qps.",
		"active_srqs":                                "Number of active_srqs.",
		"active_ud_qps":                              "Number of active_ud_qps.",
		"bad_resp_err":                               "Number of bad_resp_err.",
		"db_fifo_register":                           "Number of db_fifo_register.",
		"duplicate_request":                          "Number of duplicate_requests.",
		"implied_nak_seq_err":                        "Number of implied_nak_seq_err.",
		"lifespan":                                   "Lifespan.",
		"local_ack_timeout_err":                      "Number of local_ack_timeout_err.",
		"local_protection_err":                       "Number of local_protection_err.",
		"local_qp_op_err":                            "Number of local_qp_op_err.",
		"max_retry_exceeded":                         "Number of max_retry_exceeded.",
		"mem_mgmt_op_err":                            "Number of mem_mgmt_op_err.",
		"missing_resp":                               "Number of missing_resp.",
		"np_cnp_sent":                                "Number of np_cnp_sent.",
		"np_ecn_marked_roce_packets":                 "Number of np_ecn_marked_roce_packets.",
		"oos_drop_count":                             "Number of oos_drop_count.",
		"out_of_buffer":                              "Number of out_of_buffer.",
		"out_of_sequence":                            "Number of out_of_sequence.",
		"pacing_alerts":                              "Number of pacing_alerts.",
		"pacing_complete":                            "Number of pacing_complete.",
		"pacing_reschedule":                          "Number of pacing_reschedule.",
		"packet_seq_err":                             "Number of packet_seq_err.",
		"recoverable_errors":                         "Number of recoverable_errors.",
		"remote_access_err":                          "Number of remote_access_err.",
		"remote_invalid_req_err":                     "Number of remote_invalid_req_err.",
		"remote_op_err":                              "Number of remote_op_err.",
		"req_cqe_error":                              "Number of req_cqe_error.",
		"req_cqe_flush_error":                        "Number of req_cqe_flush_error.",
		"req_remote_access_errors":                   "Number of req_remote_access_errors.",
		"req_remote_invalid_request":                 "Number of req_remote_invalid_request.",
		"res_cmp_err":                                "Number of res_cmp_err.",
		"res_cq_load_err":                            "Number of res_cq_load_err.",
		"res_exceed_max":                             "Number of res_exceed_max.",
		"res_exceeds_wqe":                            "Number of res_exceeds_wqe.",
		"res_invalid_dup_rkey":                       "Number of res_invalid_dup_rkey.",
		"res_irrq_oflow":                             "Number of res_irrq_oflow.",
		"resize_cq_cnt":                              "Number of resize_cq_cnt.",
		"res_length_mismatch":                        "Number of res_length_mismatch.",
		"res_mem_err":                                "Number of res_mem_err.",
		"res_opcode_err":                             "Number of res_opcode_err.",
		"resp_cqe_error":                             "Number of resp_cqe_error.",
		"resp_cqe_flush_error":                       "Number of resp_cqe_flush_error.",
		"resp_local_length_error":                    "Number of resp_local_length_error.",
		"resp_remote_access_errors":                  "Number of resp_remote_access_errors.",
		"res_rem_inv_err":                            "Number of res_rem_inv_err.",
		"res_rx_domain_err":                          "Number of inbound res_rx_domain_err.",
		"res_rx_invalid_rkey":                        "Number of inbound res_rx_invalid_rkey.",
		"res_rx_no_perm":                             "Number of inbound res_rx_no_perm.",
		"res_rx_pci_err":                             "Number of inbound res_rx_pci_err.",
		"res_rx_range_err":                           "Number of inbound res_rx_range_err.",
		"res_srq_err":                                "Number of res_srq_err.",
		"res_srq_load_err":                           "Number of res_srq_load_err.",
		"res_tx_domain_err":                          "Number of outbound res_tx_domain_err.",
		"res_tx_invalid_rkey":                        "Number of outbound res_tx_invalid_rkey.",
		"res_tx_no_perm":                             "Number of outbound res_tx_no_perm.",
		"res_tx_pci_err":                             "Number of outbound res_tx_pci_err.",
		"res_tx_range_err":                           "Number of outbound res_tx_range_err.",
		"res_unaligned_atomic":                       "Number of res_unaligned_atomic.",
		"res_unsup_opcode":                           "Number of res_unsup_opcode.",
		"res_wqe_format_err":                         "Number of res_wqe_format_err.",
		"rnr_nak_retry_err":                          "Number of rnr_nak_retry_err.",
		"rnr_naks_rcvd":                              "Number of rnr_naks_rcvd.",
		"roce_adp_retrans_to":                        "Number of roce_adp_retrans_to.",
		"roce_adp_retrans":                           "Number of roce_adp_retrans.",
		"roce_slow_restart_cnps":                     "Number of roce_slow_restart_cnps.",
		"roce_slow_restart_trans":                    "Number of roce_slow_restart_trans.",
		"roce_slow_restart":                          "Number of roce_slow_restart.",
		"rp_cnp_handled":                             "Number of rp_cnp_handled.",
		"rp_cnp_ignored":                             "Number of rp_cnp_ignored.",
		"rx_atomic_requests":                         "Number of rx_atomic_requests.",
		"rx_atomic_req":                              "Number of inbound atomic_req packets.",
		"rx_bytes":                                   "Number of inbound data octets rx_bytes.",
		"rx_cnp_pkts":                                "Number of inbound cnp packets.",
		"rx_dct_connect":                             "Number of inbound dct_connect packets.",
		"rx_ecn_marked_pkts":                         "Number of inbound ecn marked packets.",
		"rx_good_bytes":                              "Number of inbound good data octets.",
		"rx_good_pkts":                               "Number of inbound packets rx_good_pkts.",
		"rx_icrc_encapsulated":                       "Number of inbound icrc_encapsulated.",
		"rx_out_of_buffer":                           "Number of inbound out_of_buffer.",
		"rx_pkts":                                    "Number of inbound packets.",
		"rx_read_requests":                           "Number of inbound read_requests.",
		"rx_read_req":                                "Number of inbound read_req.",
		"rx_read_resp":                               "Number of inbound read_resp.",
		"rx_roce_discards":                           "Number of inbound roce discards.",
		"rx_roce_errors":                             "Number of inbound roce errors.",
		"rx_roce_good_bytes":                         "Number of inbound roce good data octets",
		"rx_roce_good_pkts":                          "Number of inbound roce good packets.",
		"rx_roce_only_bytes":                         "Number of inbound roce only data octets .",
		"rx_roce_only_pkts":                          "Number of inbound roce only packets.",
		"rx_send_req":                                "Number of inbound send_req.",
		"rx_write_requests":                          "Number of inbound write_requests.",
		"rx_write_req":                               "Number of inbound write_req.",
		"seq_err_naks_rcvd":                          "Number of seq_err_naks_rcvd.",
		"to_retransmits":                             "Number of to_retransmits.",
		"tx_atomic_req":                              "Number of outbound atomic_req.",
		"tx_bytes":                                   "Number of outbound data octets.",
		"tx_cnp_pkts":                                "Number of outbound cnp packets.",
		"tx_pkts":                                    "Number of outbound packets.",
		"tx_read_req":                                "Number of outbound read_req.",
		"tx_read_resp":                               "Number of outbound read_resp.",
		"tx_roce_discards":                           "Number of outbound roce discards.",
		"tx_roce_errors":                             "Number of outbound roce errors.",
		"tx_roce_only_bytes":                         "Number of roce only outbound data octets",
		"tx_roce_only_pkts":                          "Number of outbound roce only packets.",
		"tx_send_req":                                "Number of outbound send_req.",
		"tx_write_req":                               "Number of outbound write_req.",
		"unrecoverable_err":                          "Number of unrecoverable_err.",
		"watermark_ahs":                              "Number of watermark_ahs.",
		"watermark_cqs":                              "Number of watermark_cqs.",
		"watermark_mrs":                              "Number of watermark_mrs.",
		"watermark_mws":                              "Number of watermark_mws.",
		"watermark_pds":                              "Number of watermark_pds.",
		"watermark_qps":                              "Number of watermark_qps.",
		"watermark_rc_qps":                           "Number of watermark_rc_qps.",
		"watermark_srqs":                             "Number of watermark_srqs.",
		"watermark_ud_qps":                           "Number of watermark_ud_qps.",
	}

	i.metricDescs = make(map[string]*prometheus.Desc)
	i.subsystem = "infiniband"

	for metricName, description := range descriptions {
		i.metricDescs[metricName] = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, i.subsystem, metricName),
			description,
			[]string{"device", "port"},
			nil,
		)
	}

	return &i, nil
}

func (c *infinibandCollector) pushMetric(ch chan<- prometheus.Metric, name string, value uint64, deviceName string, port string, valueType prometheus.ValueType) {
	ch <- prometheus.MustNewConstMetric(c.metricDescs[name], valueType, float64(value), deviceName, port)
}

func (c *infinibandCollector) pushCounter(ch chan<- prometheus.Metric, name string, value *uint64, deviceName string, port string) {
	if value != nil {
		c.pushMetric(ch, name, *value, deviceName, port, prometheus.CounterValue)
	}
}

func (c *infinibandCollector) Update(ch chan<- prometheus.Metric) error {
	devices, err := c.fs.InfiniBandClass()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			c.logger.Debug("infiniband statistics not found, skipping")
			return ErrNoData
		}
		return fmt.Errorf("error obtaining InfiniBand class info: %w", err)
	}

	for _, device := range devices {
		infoDesc := prometheus.NewDesc(
			prometheus.BuildFQName(namespace, c.subsystem, "info"),
			"Non-numeric data from /sys/class/infiniband/<device>, value is always 1.",
			[]string{"device", "board_id", "firmware_version", "hca_type"},
			nil,
		)
		infoValue := 1.0
		ch <- prometheus.MustNewConstMetric(infoDesc, prometheus.GaugeValue, infoValue, device.Name, device.BoardID, device.FirmwareVersion, device.HCAType)

		for _, port := range device.Ports {
			portStr := strconv.FormatUint(uint64(port.Port), 10)

			c.pushMetric(ch, "state_id", uint64(port.StateID), port.Name, portStr, prometheus.GaugeValue)
			c.pushMetric(ch, "physical_state_id", uint64(port.PhysStateID), port.Name, portStr, prometheus.GaugeValue)
			c.pushMetric(ch, "rate_bytes_per_second", port.Rate, port.Name, portStr, prometheus.GaugeValue)

			c.pushCounter(ch, "legacy_multicast_packets_received_total", port.Counters.LegacyPortMulticastRcvPackets, port.Name, portStr)
			c.pushCounter(ch, "legacy_multicast_packets_transmitted_total", port.Counters.LegacyPortMulticastXmitPackets, port.Name, portStr)
			c.pushCounter(ch, "legacy_data_received_bytes_total", port.Counters.LegacyPortRcvData64, port.Name, portStr)
			c.pushCounter(ch, "legacy_packets_received_total", port.Counters.LegacyPortRcvPackets64, port.Name, portStr)
			c.pushCounter(ch, "legacy_unicast_packets_received_total", port.Counters.LegacyPortUnicastRcvPackets, port.Name, portStr)
			c.pushCounter(ch, "legacy_unicast_packets_transmitted_total", port.Counters.LegacyPortUnicastXmitPackets, port.Name, portStr)
			c.pushCounter(ch, "legacy_data_transmitted_bytes_total", port.Counters.LegacyPortXmitData64, port.Name, portStr)
			c.pushCounter(ch, "legacy_packets_transmitted_total", port.Counters.LegacyPortXmitPackets64, port.Name, portStr)
			c.pushCounter(ch, "excessive_buffer_overrun_errors_total", port.Counters.ExcessiveBufferOverrunErrors, port.Name, portStr)
			c.pushCounter(ch, "link_downed_total", port.Counters.LinkDowned, port.Name, portStr)
			c.pushCounter(ch, "link_error_recovery_total", port.Counters.LinkErrorRecovery, port.Name, portStr)
			c.pushCounter(ch, "local_link_integrity_errors_total", port.Counters.LocalLinkIntegrityErrors, port.Name, portStr)
			c.pushCounter(ch, "multicast_packets_received_total", port.Counters.MulticastRcvPackets, port.Name, portStr)
			c.pushCounter(ch, "multicast_packets_transmitted_total", port.Counters.MulticastXmitPackets, port.Name, portStr)
			c.pushCounter(ch, "port_constraint_errors_received_total", port.Counters.PortRcvConstraintErrors, port.Name, portStr)
			c.pushCounter(ch, "port_constraint_errors_transmitted_total", port.Counters.PortXmitConstraintErrors, port.Name, portStr)
			c.pushCounter(ch, "port_data_received_bytes_total", port.Counters.PortRcvData, port.Name, portStr)
			c.pushCounter(ch, "port_data_transmitted_bytes_total", port.Counters.PortXmitData, port.Name, portStr)
			c.pushCounter(ch, "port_discards_received_total", port.Counters.PortRcvDiscards, port.Name, portStr)
			c.pushCounter(ch, "port_discards_transmitted_total", port.Counters.PortXmitDiscards, port.Name, portStr)
			c.pushCounter(ch, "port_errors_received_total", port.Counters.PortRcvErrors, port.Name, portStr)
			c.pushCounter(ch, "port_packets_received_total", port.Counters.PortRcvPackets, port.Name, portStr)
			c.pushCounter(ch, "port_packets_transmitted_total", port.Counters.PortXmitPackets, port.Name, portStr)
			c.pushCounter(ch, "port_transmit_wait_total", port.Counters.PortXmitWait, port.Name, portStr)
			c.pushCounter(ch, "unicast_packets_received_total", port.Counters.UnicastRcvPackets, port.Name, portStr)
			c.pushCounter(ch, "unicast_packets_transmitted_total", port.Counters.UnicastXmitPackets, port.Name, portStr)
			c.pushCounter(ch, "port_receive_remote_physical_errors_total", port.Counters.PortRcvRemotePhysicalErrors, port.Name, portStr)
			c.pushCounter(ch, "port_receive_switch_relay_errors_total", port.Counters.PortRcvSwitchRelayErrors, port.Name, portStr)
			c.pushCounter(ch, "symbol_error_total", port.Counters.SymbolError, port.Name, portStr)
			c.pushCounter(ch, "vl15_dropped_total", port.Counters.VL15Dropped, port.Name, portStr)
			c.pushCounter(ch, "active_ahs", port.HwCounters.ActiveAhs, port.Name, portStr)
			c.pushCounter(ch, "active_cqs", port.HwCounters.ActiveCqs, port.Name, portStr)
			c.pushCounter(ch, "active_mrs", port.HwCounters.ActiveMrs, port.Name, portStr)
			c.pushCounter(ch, "active_mws", port.HwCounters.ActiveMws, port.Name, portStr)
			c.pushCounter(ch, "active_pds", port.HwCounters.ActivePds, port.Name, portStr)
			c.pushCounter(ch, "active_qps", port.HwCounters.ActiveQps, port.Name, portStr)
			c.pushCounter(ch, "active_rc_qps", port.HwCounters.ActiveRcQps, port.Name, portStr)
			c.pushCounter(ch, "active_srqs", port.HwCounters.ActiveSrqs, port.Name, portStr)
			c.pushCounter(ch, "active_ud_qps", port.HwCounters.ActiveUdQps, port.Name, portStr)
			c.pushCounter(ch, "bad_resp_err", port.HwCounters.BadRespErr, port.Name, portStr)
			c.pushCounter(ch, "db_fifo_register", port.HwCounters.DbFifoRegister, port.Name, portStr)
			c.pushCounter(ch, "duplicate_request", port.HwCounters.DuplicateRequest, port.Name, portStr)
			c.pushCounter(ch, "implied_nak_seq_err", port.HwCounters.ImpliedNakSeqErr, port.Name, portStr)
			c.pushCounter(ch, "lifespan", port.HwCounters.Lifespan, port.Name, portStr)
			c.pushCounter(ch, "local_ack_timeout_err", port.HwCounters.LocalAckTimeoutErr, port.Name, portStr)
			c.pushCounter(ch, "local_protection_err", port.HwCounters.LocalProtectionErr, port.Name, portStr)
			c.pushCounter(ch, "local_qp_op_err", port.HwCounters.LocalQpOpErr, port.Name, portStr)
			c.pushCounter(ch, "max_retry_exceeded", port.HwCounters.MaxRetryExceeded, port.Name, portStr)
			c.pushCounter(ch, "mem_mgmt_op_err", port.HwCounters.MemMgmtOpErr, port.Name, portStr)
			c.pushCounter(ch, "missing_resp", port.HwCounters.MissingResp, port.Name, portStr)
			c.pushCounter(ch, "np_cnp_sent", port.HwCounters.NpCnpSent, port.Name, portStr)
			c.pushCounter(ch, "np_ecn_marked_roce_packets", port.HwCounters.NpEcnMarkedRocePackets, port.Name, portStr)
			c.pushCounter(ch, "oos_drop_count", port.HwCounters.OosDropCount, port.Name, portStr)
			c.pushCounter(ch, "out_of_buffer", port.HwCounters.OutOfBuffer, port.Name, portStr)
			c.pushCounter(ch, "out_of_sequence", port.HwCounters.OutOfSequence, port.Name, portStr)
			c.pushCounter(ch, "pacing_alerts", port.HwCounters.PacingAlerts, port.Name, portStr)
			c.pushCounter(ch, "pacing_complete", port.HwCounters.PacingComplete, port.Name, portStr)
			c.pushCounter(ch, "pacing_reschedule", port.HwCounters.PacingReschedule, port.Name, portStr)
			c.pushCounter(ch, "packet_seq_err", port.HwCounters.PacketSeqErr, port.Name, portStr)
			c.pushCounter(ch, "recoverable_errors", port.HwCounters.RecoverableErrors, port.Name, portStr)
			c.pushCounter(ch, "remote_access_err", port.HwCounters.RemoteAccessErr, port.Name, portStr)
			c.pushCounter(ch, "remote_invalid_req_err", port.HwCounters.RemoteInvalidReqErr, port.Name, portStr)
			c.pushCounter(ch, "remote_op_err", port.HwCounters.RemoteOpErr, port.Name, portStr)
			c.pushCounter(ch, "req_cqe_error", port.HwCounters.ReqCqeError, port.Name, portStr)
			c.pushCounter(ch, "req_cqe_flush_error", port.HwCounters.ReqCqeFlushError, port.Name, portStr)
			c.pushCounter(ch, "req_remote_access_errors", port.HwCounters.ReqRemoteAccessErrors, port.Name, portStr)
			c.pushCounter(ch, "req_remote_invalid_request", port.HwCounters.ReqRemoteInvalidRequest, port.Name, portStr)
			c.pushCounter(ch, "res_cmp_err", port.HwCounters.ResCmpErr, port.Name, portStr)
			c.pushCounter(ch, "res_cq_load_err", port.HwCounters.ResCqLoadErr, port.Name, portStr)
			c.pushCounter(ch, "res_exceed_max", port.HwCounters.ResExceedMax, port.Name, portStr)
			c.pushCounter(ch, "res_exceeds_wqe", port.HwCounters.ResExceedsWqe, port.Name, portStr)
			c.pushCounter(ch, "res_invalid_dup_rkey", port.HwCounters.ResInvalidDupRkey, port.Name, portStr)
			c.pushCounter(ch, "res_irrq_oflow", port.HwCounters.ResIrrqOflow, port.Name, portStr)
			c.pushCounter(ch, "resize_cq_cnt", port.HwCounters.ResizeCqCnt, port.Name, portStr)
			c.pushCounter(ch, "res_length_mismatch", port.HwCounters.ResLengthMismatch, port.Name, portStr)
			c.pushCounter(ch, "res_mem_err", port.HwCounters.ResMemErr, port.Name, portStr)
			c.pushCounter(ch, "res_opcode_err", port.HwCounters.ResOpcodeErr, port.Name, portStr)
			c.pushCounter(ch, "resp_cqe_error", port.HwCounters.RespCqeError, port.Name, portStr)
			c.pushCounter(ch, "resp_cqe_flush_error", port.HwCounters.RespCqeFlushError, port.Name, portStr)
			c.pushCounter(ch, "resp_local_length_error", port.HwCounters.RespLocalLengthError, port.Name, portStr)
			c.pushCounter(ch, "resp_remote_access_errors", port.HwCounters.RespRemoteAccessErrors, port.Name, portStr)
			c.pushCounter(ch, "res_rem_inv_err", port.HwCounters.ResRemInvErr, port.Name, portStr)
			c.pushCounter(ch, "res_rx_domain_err", port.HwCounters.ResRxDomainErr, port.Name, portStr)
			c.pushCounter(ch, "res_rx_invalid_rkey", port.HwCounters.ResRxInvalidRkey, port.Name, portStr)
			c.pushCounter(ch, "res_rx_no_perm", port.HwCounters.ResRxNoPerm, port.Name, portStr)
			c.pushCounter(ch, "res_rx_pci_err", port.HwCounters.ResRxPciErr, port.Name, portStr)
			c.pushCounter(ch, "res_rx_range_err", port.HwCounters.ResRxRangeErr, port.Name, portStr)
			c.pushCounter(ch, "res_srq_err", port.HwCounters.ResSrqErr, port.Name, portStr)
			c.pushCounter(ch, "res_srq_load_err", port.HwCounters.ResSrqLoadErr, port.Name, portStr)
			c.pushCounter(ch, "res_tx_domain_err", port.HwCounters.ResTxDomainErr, port.Name, portStr)
			c.pushCounter(ch, "res_tx_invalid_rkey", port.HwCounters.ResTxInvalidRkey, port.Name, portStr)
			c.pushCounter(ch, "res_tx_no_perm", port.HwCounters.ResTxNoPerm, port.Name, portStr)
			c.pushCounter(ch, "res_tx_pci_err", port.HwCounters.ResTxPciErr, port.Name, portStr)
			c.pushCounter(ch, "res_tx_range_err", port.HwCounters.ResTxRangeErr, port.Name, portStr)
			c.pushCounter(ch, "res_unaligned_atomic", port.HwCounters.ResUnalignedAtomic, port.Name, portStr)
			c.pushCounter(ch, "res_unsup_opcode", port.HwCounters.ResUnsupOpcode, port.Name, portStr)
			c.pushCounter(ch, "res_wqe_format_err", port.HwCounters.ResWqeFormatErr, port.Name, portStr)
			c.pushCounter(ch, "rnr_nak_retry_err", port.HwCounters.RnrNakRetryErr, port.Name, portStr)
			c.pushCounter(ch, "rnr_naks_rcvd", port.HwCounters.RnrNaksRcvd, port.Name, portStr)
			c.pushCounter(ch, "roce_adp_retrans_to", port.HwCounters.RoceAdpRetransTo, port.Name, portStr)
			c.pushCounter(ch, "roce_adp_retrans", port.HwCounters.RoceAdpRetrans, port.Name, portStr)
			c.pushCounter(ch, "roce_slow_restart_cnps", port.HwCounters.RoceSlowRestartCnps, port.Name, portStr)
			c.pushCounter(ch, "roce_slow_restart_trans", port.HwCounters.RoceSlowRestartTrans, port.Name, portStr)
			c.pushCounter(ch, "roce_slow_restart", port.HwCounters.RoceSlowRestart, port.Name, portStr)
			c.pushCounter(ch, "rp_cnp_handled", port.HwCounters.RpCnpHandled, port.Name, portStr)
			c.pushCounter(ch, "rp_cnp_ignored", port.HwCounters.RpCnpIgnored, port.Name, portStr)
			c.pushCounter(ch, "rx_atomic_requests", port.HwCounters.RxAtomicRequests, port.Name, portStr)
			c.pushCounter(ch, "rx_atomic_req", port.HwCounters.RxAtomicReq, port.Name, portStr)
			c.pushCounter(ch, "rx_bytes", port.HwCounters.RxBytes, port.Name, portStr)
			c.pushCounter(ch, "rx_cnp_pkts", port.HwCounters.RxCnpPkts, port.Name, portStr)
			c.pushCounter(ch, "rx_dct_connect", port.HwCounters.RxDctConnect, port.Name, portStr)
			c.pushCounter(ch, "rx_ecn_marked_pkts", port.HwCounters.RxEcnMarkedPkts, port.Name, portStr)
			c.pushCounter(ch, "rx_good_bytes", port.HwCounters.RxGoodBytes, port.Name, portStr)
			c.pushCounter(ch, "rx_good_pkts", port.HwCounters.RxGoodPkts, port.Name, portStr)
			c.pushCounter(ch, "rx_icrc_encapsulated", port.HwCounters.RxIcrcEncapsulated, port.Name, portStr)
			c.pushCounter(ch, "rx_out_of_buffer", port.HwCounters.RxOutOfBuffer, port.Name, portStr)
			c.pushCounter(ch, "rx_pkts", port.HwCounters.RxPkts, port.Name, portStr)
			c.pushCounter(ch, "rx_read_requests", port.HwCounters.RxReadRequests, port.Name, portStr)
			c.pushCounter(ch, "rx_read_req", port.HwCounters.RxReadReq, port.Name, portStr)
			c.pushCounter(ch, "rx_read_resp", port.HwCounters.RxReadResp, port.Name, portStr)
			c.pushCounter(ch, "rx_roce_discards", port.HwCounters.RxRoceDiscards, port.Name, portStr)
			c.pushCounter(ch, "rx_roce_errors", port.HwCounters.RxRoceErrors, port.Name, portStr)
			c.pushCounter(ch, "rx_roce_good_bytes", port.HwCounters.RxRoceGoodBytes, port.Name, portStr)
			c.pushCounter(ch, "rx_roce_good_pkts", port.HwCounters.RxRoceGoodPkts, port.Name, portStr)
			c.pushCounter(ch, "rx_roce_only_bytes", port.HwCounters.RxRoceOnlyBytes, port.Name, portStr)
			c.pushCounter(ch, "rx_roce_only_pkts", port.HwCounters.RxRoceOnlyPkts, port.Name, portStr)
			c.pushCounter(ch, "rx_send_req", port.HwCounters.RxSendReq, port.Name, portStr)
			c.pushCounter(ch, "rx_write_requests", port.HwCounters.RxWriteRequests, port.Name, portStr)
			c.pushCounter(ch, "rx_write_req", port.HwCounters.RxWriteReq, port.Name, portStr)
			c.pushCounter(ch, "seq_err_naks_rcvd", port.HwCounters.SeqErrNaksRcvd, port.Name, portStr)
			c.pushCounter(ch, "to_retransmits", port.HwCounters.ToRetransmits, port.Name, portStr)
			c.pushCounter(ch, "tx_atomic_req", port.HwCounters.TxAtomicReq, port.Name, portStr)
			c.pushCounter(ch, "tx_bytes", port.HwCounters.TxBytes, port.Name, portStr)
			c.pushCounter(ch, "tx_cnp_pkts", port.HwCounters.TxCnpPkts, port.Name, portStr)
			c.pushCounter(ch, "tx_pkts", port.HwCounters.TxPkts, port.Name, portStr)
			c.pushCounter(ch, "tx_read_req", port.HwCounters.TxReadReq, port.Name, portStr)
			c.pushCounter(ch, "tx_read_resp", port.HwCounters.TxReadResp, port.Name, portStr)
			c.pushCounter(ch, "tx_roce_discards", port.HwCounters.TxRoceDiscards, port.Name, portStr)
			c.pushCounter(ch, "tx_roce_errors", port.HwCounters.TxRoceErrors, port.Name, portStr)
			c.pushCounter(ch, "tx_roce_only_bytes", port.HwCounters.TxRoceOnlyBytes, port.Name, portStr)
			c.pushCounter(ch, "tx_roce_only_pkts", port.HwCounters.TxRoceOnlyPkts, port.Name, portStr)
			c.pushCounter(ch, "tx_send_req", port.HwCounters.TxSendReq, port.Name, portStr)
			c.pushCounter(ch, "tx_write_req", port.HwCounters.TxWriteReq, port.Name, portStr)
			c.pushCounter(ch, "unrecoverable_err", port.HwCounters.UnrecoverableErr, port.Name, portStr)
			c.pushCounter(ch, "watermark_ahs", port.HwCounters.WatermarkAhs, port.Name, portStr)
			c.pushCounter(ch, "watermark_cqs", port.HwCounters.WatermarkCqs, port.Name, portStr)
			c.pushCounter(ch, "watermark_mrs", port.HwCounters.WatermarkMrs, port.Name, portStr)
			c.pushCounter(ch, "watermark_mws", port.HwCounters.WatermarkMws, port.Name, portStr)
			c.pushCounter(ch, "watermark_pds", port.HwCounters.WatermarkPds, port.Name, portStr)
			c.pushCounter(ch, "watermark_qps", port.HwCounters.WatermarkQps, port.Name, portStr)
			c.pushCounter(ch, "watermark_rc_qps", port.HwCounters.WatermarkRcQps, port.Name, portStr)
			c.pushCounter(ch, "watermark_srqs", port.HwCounters.WatermarkSrqs, port.Name, portStr)
			c.pushCounter(ch, "watermark_ud_qps", port.HwCounters.WatermarkUdQps, port.Name, portStr)
		}
	}

	return nil
}
