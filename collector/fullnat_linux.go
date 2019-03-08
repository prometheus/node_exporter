// Copyright 2015 The Prometheus Authors
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

// +build !noipvs

package collector

import (
	"fmt"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/procfs"
)

type fnatCollector struct {
	Collector
	fs                                                                                                                        procfs.FS
	backendConnectionsActive, backendConnectionsInact, backendWeight                                                          typedDesc
	connections, incomingPackets, outgoingPackets, incomingBytes, outgoingBytes                                               typedDesc
	fullnatAddToaOk, fullnatAddToaFailLen, fullnatAddToaHeadFull, fullnatAddToaFailMem                                        typedDesc
	fullnatAddToaFailProto, fullnatConnReused, fullnatConnReusedClose, fullnatConnReusedTimewait                              typedDesc
	fullnatConnReusedFinwait, fullnatConnReusedClosewait, fullnatConnReusedLastack, fullnatConnReusedEstab                    typedDesc
	synproxyRsError, synproxyNullAck, synproxyBadAck, synproxyOkAck, synproxySynCnt                                           typedDesc
	synproxyAckstorm, synproxySynsendQlen, synproxyConnReused, synproxyConnReusedClose                                        typedDesc
	synproxyConnReusedTimewait, synproxyConnReusedFinwait, synproxyConnReusedClosewait, synproxyConnReusedLastack             typedDesc
	defenceIpFragDrop, defenceIpFragGather, defenceTcpDrop, defenceUdpDrop, fastXmitReject                                    typedDesc
	fastXmitPass, fastXmitSkbCopy, fastXmitNoMac, fastXmitSynproxySave, fastXmitDevLost                                       typedDesc
	rstInSynSent, rstOutSynSent, rstInEstablished, rstOutEstablished, groPass, lroReject, xmitUnexpectedMtu, connSchedUnreach typedDesc
	stat                                                                                                                      typedDesc
}

func init() {
	registerCollector("fnat", defaultDisabled, NewFNATCollector)
}

// NewFNATCollector sets up a new collector for FNAT metrics. It accepts the
// "procfs" config parameter to override the default proc location (/proc).
func NewFNATCollector() (Collector, error) {
	return newFNATCollector()
}

func newFNATCollector() (*fnatCollector, error) {
	var (
		fnatStatLabelNames = []string{
			"CPU",
		}
		fnatBackendLabelNames = []string{
			"local_address",
			"local_port",
			"remote_address",
			"remote_port",
			"proto",
		}
		c         fnatCollector
		err       error
		subsystem = "fnat"
	)

	c.fs, err = procfs.NewFS(*procPath)
	if err != nil {
		return nil, err
	}

	c.connections = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "connections_total"),
		"The total number of connections made.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.incomingPackets = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "incoming_packets_total"),
		"The total number of incoming packets.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.outgoingPackets = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "outgoing_packets_total"),
		"The total number of outgoing packets.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.incomingBytes = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "incoming_bytes_total"),
		"The total amount of incoming data.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.outgoingBytes = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "outgoing_bytes_total"),
		"The total amount of outgoing data.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.backendConnectionsActive = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "backend_connections_active"),
		"The current active connections by local and remote address.",
		fnatBackendLabelNames, nil,
	), prometheus.GaugeValue}
	c.backendConnectionsInact = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "backend_connections_inactive"),
		"The current inactive connections by local and remote address.",
		fnatBackendLabelNames, nil,
	), prometheus.GaugeValue}
	c.backendWeight = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "backend_weight"),
		"The current backend weight by local and remote address.",
		fnatBackendLabelNames, nil,
	), prometheus.GaugeValue}

	c.fullnatAddToaOk = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "fullnat_add_toa_ok_total"),
		"fullnat_add_toa_ok_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.fullnatAddToaFailLen = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "fullnat_add_toa_fail_len_total"),
		"fullnat_add_toa_fail_len_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.fullnatAddToaHeadFull = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "fullnat_add_toa_head_full_total"),
		"fullnat_add_toa_head_full_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.fullnatAddToaFailMem = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "fullnat_add_toa_fail_mem_total"),
		"fullnat_add_toa_fail_mem_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.fullnatAddToaFailProto = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "fullnat_add_toa_fail_proto_total"),
		"fullnat_add_toa_fail_proto_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.fullnatConnReused = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "fullnat_conn_reused_total"),
		"fullnat_conn_reused_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.fullnatConnReusedClose = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "fullnat_conn_reused_close_total"),
		"fullnat_conn_reused_close_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.fullnatConnReusedTimewait = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "fullnat_conn_reused_timewait_total"),
		"fullnat_conn_reused_timewait_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.fullnatConnReusedFinwait = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "fullnat_conn_reused_finwait_total"),
		"fullnat_conn_reused_finwait_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.fullnatConnReusedClosewait = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "fullnat_conn_reused_closewait_total"),
		"fullnat_conn_reused_closewait_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.fullnatConnReusedLastack = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "fullnat_conn_reused_lastack_total"),
		"fullnat_conn_reused_lastack_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.fullnatConnReusedEstab = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "fullnat_conn_reused_estab_total"),
		"fullnat_conn_reused_estab_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.synproxyRsError = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "synproxy_rs_error_total"),
		"synproxy_rs_error_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}

	c.synproxyNullAck = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "synproxy_null_ack_total"),
		"synproxy_null_ack_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}

	c.synproxyBadAck = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "synproxy_bad_ack_total"),
		"synproxy_bad_ack_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}

	c.synproxyOkAck = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "synproxy_ok_ack_total"),
		"synproxy_ok_ack_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.synproxySynCnt = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "synproxy_syn_cnt_total"),
		"synproxy_syn_cnt_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.synproxyAckstorm = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "synproxy_ackstorm_total"),
		"synproxy_ackstorm_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.synproxySynsendQlen = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "synproxy_synsend_qlen_total"),
		"synproxy_synsend_qlen_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.synproxyConnReused = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "synproxy_conn_reused_total"),
		"synproxy_conn_reused_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.synproxyConnReusedClose = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "synproxy_conn_reused_close_total"),
		"synproxy_conn_reused_close_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.synproxyConnReusedTimewait = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "synproxy_conn_reused_timewait_total"),
		"synproxy_conn_reused_timewait_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.synproxyConnReusedFinwait = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "synproxy_conn_reused_finwait_total"),
		"synproxy_conn_reused_finwait_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.synproxyConnReusedClosewait = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "synproxy_conn_reused_closewait_total"),
		"synproxy_conn_reused_closewait_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.synproxyConnReusedLastack = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "synproxy_conn_reused_lastack_total"),
		"synproxy_conn_reused_lastack_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.defenceIpFragDrop = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "defence_ip_frag_drop_total"),
		"defence_ip_frag_drop_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.defenceIpFragGather = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "defence_ip_frag_gather_total"),
		"defence_ip_frag_gather_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.defenceTcpDrop = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "defence_tcp_drop_total"),
		"defence_tcp_drop_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.defenceUdpDrop = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "defence_udp_drop_total"),
		"defence_udp_drop_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.fastXmitReject = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "fast_xmit_reject_total"),
		"fast_xmit_reject_total_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.fastXmitPass = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "fast_xmit_pass_total"),
		"fast_xmit_pass_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.fastXmitSkbCopy = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "fast_xmit_skb_copy_total"),
		"fast_xmit_skb_copy_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.fastXmitNoMac = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "fast_xmit_no_mac_total"),
		"fast_xmit_no_mac_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.fastXmitSynproxySave = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "fast_xmit_synproxy_save_total"),
		"fast_xmit_synproxy_save_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.fastXmitDevLost = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "fast_xmit_dev_lost_total"),
		"fast_xmit_dev_lost_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.rstInSynSent = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "rst_in_syn_sent_total"),
		"rst_in_syn_sent_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.rstOutSynSent = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "rst_out_syn_sent_total"),
		"rst_out_syn_sent_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.rstInEstablished = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "rst_in_established_total"),
		"rst_in_established_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.rstOutEstablished = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "rst_out_established_total"),
		"rst_out_established_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.groPass = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "gro_pass_total"),
		"gro_pass_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.lroReject = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "lro_reject_total"),
		"lro_reject_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.xmitUnexpectedMtu = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "xmit_unexpected_mtu_total"),
		"xmit_unexpected_mtu_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.connSchedUnreach = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "conn_sched_unreach_total"),
		"conn_sched_unreach_total.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}

	return &c, nil
}

func (c *fnatCollector) Update(ch chan<- prometheus.Metric) error {
	fnatStats, err := c.fs.NewFNATStats()
	if err != nil {
		// Cannot access ipvs metrics, report no error.
		if os.IsNotExist(err) {
			log.Debug("fnat collector metrics are not available for this system")
			return nil
		}
		return fmt.Errorf("could not get FNAT stats: %s", err)
	}
	for _, statsFnat := range fnatStats.Stat {
		statsLabelValues := []string{
			statsFnat.Cpu,
		}
		ch <- c.connections.mustNewConstMetric(float64(statsFnat.Connections), statsLabelValues...)
		ch <- c.incomingPackets.mustNewConstMetric(float64(statsFnat.IncomingPackets), statsLabelValues...)
		ch <- c.outgoingPackets.mustNewConstMetric(float64(statsFnat.OutgoingPackets), statsLabelValues...)
		ch <- c.incomingBytes.mustNewConstMetric(float64(statsFnat.IncomingBytes), statsLabelValues...)
		ch <- c.outgoingBytes.mustNewConstMetric(float64(statsFnat.OutgoingBytes), statsLabelValues...)

	}

	backendStats, err := c.fs.NewFNATBackendStatus()
	if err != nil {
		return fmt.Errorf("could not get backend status: %s", err)
	}

	for _, backend := range backendStats {
		labelValues := []string{
			backend.LocalAddress.String(),
			strconv.FormatUint(uint64(backend.LocalPort), 10),
			backend.RemoteAddress.String(),
			strconv.FormatUint(uint64(backend.RemotePort), 10),
			backend.Proto,
		}
		ch <- c.backendConnectionsActive.mustNewConstMetric(float64(backend.ActiveConn), labelValues...)
		ch <- c.backendConnectionsInact.mustNewConstMetric(float64(backend.InactConn), labelValues...)
		ch <- c.backendWeight.mustNewConstMetric(float64(backend.Weight), labelValues...)
	}

	extStats, err := c.fs.NewFNATExtStats()
	if err != nil {
		return fmt.Errorf("could not get ext  stats: %s", err)
	}

	extStatsSendData(ch, c.fullnatAddToaOk, extStats.FullnatAddToaOk)
	extStatsSendData(ch, c.fullnatAddToaHeadFull, extStats.FullnatAddToaHeadFull)
	extStatsSendData(ch, c.fullnatAddToaFailMem, extStats.FullnatAddToaFailMem)
	extStatsSendData(ch, c.fullnatAddToaFailProto, extStats.FullnatAddToaFailProto)
	extStatsSendData(ch, c.fullnatConnReused, extStats.FullnatConnReused)
	extStatsSendData(ch, c.fullnatConnReusedClose, extStats.FullnatConnReusedClose)
	extStatsSendData(ch, c.fullnatConnReusedTimewait, extStats.FullnatConnReusedTimewait)
	extStatsSendData(ch, c.fullnatConnReusedFinwait, extStats.FullnatConnReusedFinwait)
	extStatsSendData(ch, c.fullnatConnReusedClosewait, extStats.FullnatConnReusedClosewait)
	extStatsSendData(ch, c.fullnatConnReusedLastack, extStats.FullnatConnReusedLastack)
	extStatsSendData(ch, c.fullnatConnReusedEstab, extStats.FullnatConnReusedEstab)
	extStatsSendData(ch, c.synproxyRsError, extStats.SynproxyRsError)
	extStatsSendData(ch, c.synproxyNullAck, extStats.SynproxyNullAck)
	extStatsSendData(ch, c.synproxyBadAck, extStats.SynproxyBadAck)
	extStatsSendData(ch, c.synproxyOkAck, extStats.SynproxyOkAck)
	extStatsSendData(ch, c.synproxySynCnt, extStats.SynproxySynCnt)
	extStatsSendData(ch, c.synproxyAckstorm, extStats.SynproxyAckstorm)
	extStatsSendData(ch, c.synproxySynsendQlen, extStats.SynproxySynsendQlen)
	extStatsSendData(ch, c.synproxyConnReused, extStats.SynproxyConnReused)
	extStatsSendData(ch, c.synproxyConnReusedClose, extStats.SynproxyConnReusedClose)
	extStatsSendData(ch, c.synproxyConnReusedTimewait, extStats.SynproxyConnReusedTimewait)
	extStatsSendData(ch, c.synproxyConnReusedFinwait, extStats.SynproxyConnReusedFinwait)
	extStatsSendData(ch, c.synproxyConnReusedClosewait, extStats.SynproxyConnReusedClosewait)
	extStatsSendData(ch, c.synproxyConnReusedLastack, extStats.SynproxyConnReusedLastack)
	extStatsSendData(ch, c.defenceIpFragDrop, extStats.DefenceIpFragDrop)
	extStatsSendData(ch, c.defenceIpFragGather, extStats.DefenceIpFragGather)
	extStatsSendData(ch, c.defenceTcpDrop, extStats.DefenceTcpDrop)
	extStatsSendData(ch, c.defenceUdpDrop, extStats.DefenceUdpDrop)
	extStatsSendData(ch, c.fastXmitReject, extStats.FastXmitReject)
	extStatsSendData(ch, c.fastXmitPass, extStats.FastXmitPass)
	extStatsSendData(ch, c.fastXmitSkbCopy, extStats.FastXmitSkbCopy)
	extStatsSendData(ch, c.fastXmitSynproxySave, extStats.FastXmitSynproxySave)
	extStatsSendData(ch, c.fastXmitDevLost, extStats.FastXmitDevLost)
	extStatsSendData(ch, c.rstInSynSent, extStats.RstInSynSent)
	extStatsSendData(ch, c.rstOutSynSent, extStats.RstOutSynSent)
	extStatsSendData(ch, c.rstInEstablished, extStats.RstInEstablished)
	extStatsSendData(ch, c.rstOutEstablished, extStats.RstOutEstablished)
	extStatsSendData(ch, c.groPass, extStats.GroPass)
	extStatsSendData(ch, c.lroReject, extStats.LroReject)
	extStatsSendData(ch, c.xmitUnexpectedMtu, extStats.XmitUnexpectedMtu)
	extStatsSendData(ch, c.connSchedUnreach, extStats.ConnSchedUnreach)

	return nil
}

func extStatsSendData(ch chan<- prometheus.Metric, c typedDesc, ext procfs.ExtStatsPerCpu) {
	for k, v := range ext {
		labelValues := []string{
			k,
		}
		ch <- c.mustNewConstMetric(float64(v), labelValues...)
	}
}
