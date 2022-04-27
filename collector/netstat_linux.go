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

//go:build !nonetstat
// +build !nonetstat

package collector

import (
	"fmt"
	"github.com/prometheus/procfs"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	netStatsSubsystem = "netstat"
)

type netStatCollector struct {
	proc procfs.Proc

	// Netstat
	// TcpExt
	TcpExtSyncookiesSent            *prometheus.Desc
	TcpExtSyncookiesRecv            *prometheus.Desc
	TcpExtSyncookiesFailed          *prometheus.Desc
	TcpExtEmbryonicRsts             *prometheus.Desc
	TcpExtPruneCalled               *prometheus.Desc
	TcpExtRcvPruned                 *prometheus.Desc
	TcpExtOfoPruned                 *prometheus.Desc
	TcpExtOutOfWindowIcmps          *prometheus.Desc
	TcpExtLockDroppedIcmps          *prometheus.Desc
	TcpExtArpFilter                 *prometheus.Desc
	TcpExtTW                        *prometheus.Desc
	TcpExtTWRecycled                *prometheus.Desc
	TcpExtTWKilled                  *prometheus.Desc
	TcpExtPAWSActive                *prometheus.Desc
	TcpExtPAWSEstab                 *prometheus.Desc
	TcpExtDelayedACKs               *prometheus.Desc
	TcpExtDelayedACKLocked          *prometheus.Desc
	TcpExtDelayedACKLost            *prometheus.Desc
	TcpExtListenOverflows           *prometheus.Desc
	TcpExtListenDrops               *prometheus.Desc
	TcpExtTCPHPHits                 *prometheus.Desc
	TcpExtTCPPureAcks               *prometheus.Desc
	TcpExtTCPHPAcks                 *prometheus.Desc
	TcpExtTCPRenoRecovery           *prometheus.Desc
	TcpExtTCPSackRecovery           *prometheus.Desc
	TcpExtTCPSACKReneging           *prometheus.Desc
	TcpExtTCPSACKReorder            *prometheus.Desc
	TcpExtTCPRenoReorder            *prometheus.Desc
	TcpExtTCPTSReorder              *prometheus.Desc
	TcpExtTCPFullUndo               *prometheus.Desc
	TcpExtTCPPartialUndo            *prometheus.Desc
	TcpExtTCPDSACKUndo              *prometheus.Desc
	TcpExtTCPLossUndo               *prometheus.Desc
	TcpExtTCPLostRetransmit         *prometheus.Desc
	TcpExtTCPRenoFailures           *prometheus.Desc
	TcpExtTCPSackFailures           *prometheus.Desc
	TcpExtTCPLossFailures           *prometheus.Desc
	TcpExtTCPFastRetrans            *prometheus.Desc
	TcpExtTCPSlowStartRetrans       *prometheus.Desc
	TcpExtTCPTimeouts               *prometheus.Desc
	TcpExtTCPLossProbes             *prometheus.Desc
	TcpExtTCPLossProbeRecovery      *prometheus.Desc
	TcpExtTCPRenoRecoveryFail       *prometheus.Desc
	TcpExtTCPSackRecoveryFail       *prometheus.Desc
	TcpExtTCPRcvCollapsed           *prometheus.Desc
	TcpExtTCPDSACKOldSent           *prometheus.Desc
	TcpExtTCPDSACKOfoSent           *prometheus.Desc
	TcpExtTCPDSACKRecv              *prometheus.Desc
	TcpExtTCPDSACKOfoRecv           *prometheus.Desc
	TcpExtTCPAbortOnData            *prometheus.Desc
	TcpExtTCPAbortOnClose           *prometheus.Desc
	TcpExtTCPAbortOnMemory          *prometheus.Desc
	TcpExtTCPAbortOnTimeout         *prometheus.Desc
	TcpExtTCPAbortOnLinger          *prometheus.Desc
	TcpExtTCPAbortFailed            *prometheus.Desc
	TcpExtTCPMemoryPressures        *prometheus.Desc
	TcpExtTCPMemoryPressuresChrono  *prometheus.Desc
	TcpExtTCPSACKDiscard            *prometheus.Desc
	TcpExtTCPDSACKIgnoredOld        *prometheus.Desc
	TcpExtTCPDSACKIgnoredNoUndo     *prometheus.Desc
	TcpExtTCPSpuriousRTOs           *prometheus.Desc
	TcpExtTCPMD5NotFound            *prometheus.Desc
	TcpExtTCPMD5Unexpected          *prometheus.Desc
	TcpExtTCPMD5Failure             *prometheus.Desc
	TcpExtTCPSackShifted            *prometheus.Desc
	TcpExtTCPSackMerged             *prometheus.Desc
	TcpExtTCPSackShiftFallback      *prometheus.Desc
	TcpExtTCPBacklogDrop            *prometheus.Desc
	TcpExtPFMemallocDrop            *prometheus.Desc
	TcpExtTCPMinTTLDrop             *prometheus.Desc
	TcpExtTCPDeferAcceptDrop        *prometheus.Desc
	TcpExtIPReversePathFilter       *prometheus.Desc
	TcpExtTCPTimeWaitOverflow       *prometheus.Desc
	TcpExtTCPReqQFullDoCookies      *prometheus.Desc
	TcpExtTCPReqQFullDrop           *prometheus.Desc
	TcpExtTCPRetransFail            *prometheus.Desc
	TcpExtTCPRcvCoalesce            *prometheus.Desc
	TcpExtTCPOFOQueue               *prometheus.Desc
	TcpExtTCPOFODrop                *prometheus.Desc
	TcpExtTCPOFOMerge               *prometheus.Desc
	TcpExtTCPChallengeACK           *prometheus.Desc
	TcpExtTCPSYNChallenge           *prometheus.Desc
	TcpExtTCPFastOpenActive         *prometheus.Desc
	TcpExtTCPFastOpenActiveFail     *prometheus.Desc
	TcpExtTCPFastOpenPassive        *prometheus.Desc
	TcpExtTCPFastOpenPassiveFail    *prometheus.Desc
	TcpExtTCPFastOpenListenOverflow *prometheus.Desc
	TcpExtTCPFastOpenCookieReqd     *prometheus.Desc
	TcpExtTCPFastOpenBlackhole      *prometheus.Desc
	TcpExtTCPSpuriousRtxHostQueues  *prometheus.Desc
	TcpExtBusyPollRxPackets         *prometheus.Desc
	TcpExtTCPAutoCorking            *prometheus.Desc
	TcpExtTCPFromZeroWindowAdv      *prometheus.Desc
	TcpExtTCPToZeroWindowAdv        *prometheus.Desc
	TcpExtTCPWantZeroWindowAdv      *prometheus.Desc
	TcpExtTCPSynRetrans             *prometheus.Desc
	TcpExtTCPOrigDataSent           *prometheus.Desc
	TcpExtTCPHystartTrainDetect     *prometheus.Desc
	TcpExtTCPHystartTrainCwnd       *prometheus.Desc
	TcpExtTCPHystartDelayDetect     *prometheus.Desc
	TcpExtTCPHystartDelayCwnd       *prometheus.Desc
	TcpExtTCPACKSkippedSynRecv      *prometheus.Desc
	TcpExtTCPACKSkippedPAWS         *prometheus.Desc
	TcpExtTCPACKSkippedSeq          *prometheus.Desc
	TcpExtTCPACKSkippedFinWait2     *prometheus.Desc
	TcpExtTCPACKSkippedTimeWait     *prometheus.Desc
	TcpExtTCPACKSkippedChallenge    *prometheus.Desc
	TcpExtTCPWinProbe               *prometheus.Desc
	TcpExtTCPKeepAlive              *prometheus.Desc
	TcpExtTCPMTUPFail               *prometheus.Desc
	TcpExtTCPMTUPSuccess            *prometheus.Desc
	TcpExtTCPWqueueTooBig           *prometheus.Desc

	// IpExt
	IpExtInNoRoutes      *prometheus.Desc
	IpExtInTruncatedPkts *prometheus.Desc
	IpExtInMcastPkts     *prometheus.Desc
	IpExtOutMcastPkts    *prometheus.Desc
	IpExtInBcastPkts     *prometheus.Desc
	IpExtOutBcastPkts    *prometheus.Desc
	IpExtInOctets        *prometheus.Desc
	IpExtOutOctets       *prometheus.Desc
	IpExtInMcastOctets   *prometheus.Desc
	IpExtOutMcastOctets  *prometheus.Desc
	IpExtInBcastOctets   *prometheus.Desc
	IpExtOutBcastOctets  *prometheus.Desc
	IpExtInCsumErrors    *prometheus.Desc
	IpExtInNoECTPkts     *prometheus.Desc
	IpExtInECT1Pkts      *prometheus.Desc
	IpExtInECT0Pkts      *prometheus.Desc
	IpExtInCEPkts        *prometheus.Desc
	IpExtReasmOverlaps   *prometheus.Desc

	// SNMP
	// Ip
	IpForwarding      *prometheus.Desc
	IpDefaultTTL      *prometheus.Desc
	IpInReceives      *prometheus.Desc
	IpInHdrErrors     *prometheus.Desc
	IpInAddrErrors    *prometheus.Desc
	IpForwDatagrams   *prometheus.Desc
	IpInUnknownProtos *prometheus.Desc
	IpInDiscards      *prometheus.Desc
	IpInDelivers      *prometheus.Desc
	IpOutRequests     *prometheus.Desc
	IpOutDiscards     *prometheus.Desc
	IpOutNoRoutes     *prometheus.Desc
	IpReasmTimeout    *prometheus.Desc
	IpReasmReqds      *prometheus.Desc
	IpReasmOKs        *prometheus.Desc
	IpReasmFails      *prometheus.Desc
	IpFragOKs         *prometheus.Desc
	IpFragFails       *prometheus.Desc
	IpFragCreates     *prometheus.Desc

	// Icmp
	IcmpInMsgs           *prometheus.Desc
	IcmpInErrors         *prometheus.Desc
	IcmpInCsumErrors     *prometheus.Desc
	IcmpInDestUnreachs   *prometheus.Desc
	IcmpInTimeExcds      *prometheus.Desc
	IcmpInParmProbs      *prometheus.Desc
	IcmpInSrcQuenchs     *prometheus.Desc
	IcmpInRedirects      *prometheus.Desc
	IcmpInEchos          *prometheus.Desc
	IcmpInEchoReps       *prometheus.Desc
	IcmpInTimestamps     *prometheus.Desc
	IcmpInTimestampReps  *prometheus.Desc
	IcmpInAddrMasks      *prometheus.Desc
	IcmpInAddrMaskReps   *prometheus.Desc
	IcmpOutMsgs          *prometheus.Desc
	IcmpOutErrors        *prometheus.Desc
	IcmpOutDestUnreachs  *prometheus.Desc
	IcmpOutTimeExcds     *prometheus.Desc
	IcmpOutParmProbs     *prometheus.Desc
	IcmpOutSrcQuenchs    *prometheus.Desc
	IcmpOutRedirects     *prometheus.Desc
	IcmpOutEchos         *prometheus.Desc
	IcmpOutEchoReps      *prometheus.Desc
	IcmpOutTimestamps    *prometheus.Desc
	IcmpOutTimestampReps *prometheus.Desc
	IcmpOutAddrMasks     *prometheus.Desc
	IcmpOutAddrMaskReps  *prometheus.Desc

	// IcmpMsg
	IcmpMsgInType3  *prometheus.Desc
	IcmpMsgOutType3 *prometheus.Desc

	// Tcp
	TcpRtoAlgorithm *prometheus.Desc
	TcpRtoMin       *prometheus.Desc
	TcpRtoMax       *prometheus.Desc
	TcpMaxConn      *prometheus.Desc
	TcpActiveOpens  *prometheus.Desc
	TcpPassiveOpens *prometheus.Desc
	TcpAttemptFails *prometheus.Desc
	TcpEstabResets  *prometheus.Desc
	TcpCurrEstab    *prometheus.Desc
	TcpInSegs       *prometheus.Desc
	TcpOutSegs      *prometheus.Desc
	TcpRetransSegs  *prometheus.Desc
	TcpInErrs       *prometheus.Desc
	TcpOutRsts      *prometheus.Desc
	TcpInCsumErrors *prometheus.Desc

	// Udp
	UdpInDatagrams  *prometheus.Desc
	UdpNoPorts      *prometheus.Desc
	UdpInErrors     *prometheus.Desc
	UdpOutDatagrams *prometheus.Desc
	UdpRcvbufErrors *prometheus.Desc
	UdpSndbufErrors *prometheus.Desc
	UdpInCsumErrors *prometheus.Desc
	UdpIgnoredMulti *prometheus.Desc

	// UdpLite
	UdpLiteInDatagrams  *prometheus.Desc
	UdpLiteNoPorts      *prometheus.Desc
	UdpLiteInErrors     *prometheus.Desc
	UdpLiteOutDatagrams *prometheus.Desc
	UdpLiteRcvbufErrors *prometheus.Desc
	UdpLiteSndbufErrors *prometheus.Desc
	UdpLiteInCsumErrors *prometheus.Desc
	UdpLiteIgnoredMulti *prometheus.Desc

	// Snmp6
	// Ip6
	Ip6InReceives       *prometheus.Desc
	Ip6InHdrErrors      *prometheus.Desc
	Ip6InTooBigErrors   *prometheus.Desc
	Ip6InNoRoutes       *prometheus.Desc
	Ip6InAddrErrors     *prometheus.Desc
	Ip6InUnknownProtos  *prometheus.Desc
	Ip6InTruncatedPkts  *prometheus.Desc
	Ip6InDiscards       *prometheus.Desc
	Ip6InDelivers       *prometheus.Desc
	Ip6OutForwDatagrams *prometheus.Desc
	Ip6OutRequests      *prometheus.Desc
	Ip6OutDiscards      *prometheus.Desc
	Ip6OutNoRoutes      *prometheus.Desc
	Ip6ReasmTimeout     *prometheus.Desc
	Ip6ReasmReqds       *prometheus.Desc
	Ip6ReasmOKs         *prometheus.Desc
	Ip6ReasmFails       *prometheus.Desc
	Ip6FragOKs          *prometheus.Desc
	Ip6FragFails        *prometheus.Desc
	Ip6FragCreates      *prometheus.Desc
	Ip6InMcastPkts      *prometheus.Desc
	Ip6OutMcastPkts     *prometheus.Desc
	Ip6InOctets         *prometheus.Desc
	Ip6OutOctets        *prometheus.Desc
	Ip6InMcastOctets    *prometheus.Desc
	Ip6OutMcastOctets   *prometheus.Desc
	Ip6InBcastOctets    *prometheus.Desc
	Ip6OutBcastOctets   *prometheus.Desc
	Ip6InNoECTPkts      *prometheus.Desc
	Ip6InECT1Pkts       *prometheus.Desc
	Ip6InECT0Pkts       *prometheus.Desc
	Ip6InCEPkts         *prometheus.Desc

	// Icmp6
	Icmp6InMsgs                    *prometheus.Desc
	Icmp6InErrors                  *prometheus.Desc
	Icmp6OutMsgs                   *prometheus.Desc
	Icmp6OutErrors                 *prometheus.Desc
	Icmp6InCsumErrors              *prometheus.Desc
	Icmp6InDestUnreachs            *prometheus.Desc
	Icmp6InPktTooBigs              *prometheus.Desc
	Icmp6InTimeExcds               *prometheus.Desc
	Icmp6InParmProblems            *prometheus.Desc
	Icmp6InEchos                   *prometheus.Desc
	Icmp6InEchoReplies             *prometheus.Desc
	Icmp6InGroupMembQueries        *prometheus.Desc
	Icmp6InGroupMembResponses      *prometheus.Desc
	Icmp6InGroupMembReductions     *prometheus.Desc
	Icmp6InRouterSolicits          *prometheus.Desc
	Icmp6InRouterAdvertisements    *prometheus.Desc
	Icmp6InNeighborSolicits        *prometheus.Desc
	Icmp6InNeighborAdvertisements  *prometheus.Desc
	Icmp6InRedirects               *prometheus.Desc
	Icmp6InMLDv2Reports            *prometheus.Desc
	Icmp6OutDestUnreachs           *prometheus.Desc
	Icmp6OutPktTooBigs             *prometheus.Desc
	Icmp6OutTimeExcds              *prometheus.Desc
	Icmp6OutParmProblems           *prometheus.Desc
	Icmp6OutEchos                  *prometheus.Desc
	Icmp6OutEchoReplies            *prometheus.Desc
	Icmp6OutGroupMembQueries       *prometheus.Desc
	Icmp6OutGroupMembResponses     *prometheus.Desc
	Icmp6OutGroupMembReductions    *prometheus.Desc
	Icmp6OutRouterSolicits         *prometheus.Desc
	Icmp6OutRouterAdvertisements   *prometheus.Desc
	Icmp6OutNeighborSolicits       *prometheus.Desc
	Icmp6OutNeighborAdvertisements *prometheus.Desc
	Icmp6OutRedirects              *prometheus.Desc
	Icmp6OutMLDv2Reports           *prometheus.Desc
	Icmp6InType1                   *prometheus.Desc
	Icmp6InType134                 *prometheus.Desc
	Icmp6InType135                 *prometheus.Desc
	Icmp6InType136                 *prometheus.Desc
	Icmp6InType143                 *prometheus.Desc
	Icmp6OutType133                *prometheus.Desc
	Icmp6OutType135                *prometheus.Desc
	Icmp6OutType136                *prometheus.Desc
	Icmp6OutType143                *prometheus.Desc

	// Udp6
	Udp6InDatagrams  *prometheus.Desc
	Udp6NoPorts      *prometheus.Desc
	Udp6InErrors     *prometheus.Desc
	Udp6OutDatagrams *prometheus.Desc
	Udp6RcvbufErrors *prometheus.Desc
	Udp6SndbufErrors *prometheus.Desc
	Udp6InCsumErrors *prometheus.Desc
	Udp6IgnoredMulti *prometheus.Desc

	// UdpLite6
	UdpLite6InDatagrams  *prometheus.Desc
	UdpLite6NoPorts      *prometheus.Desc
	UdpLite6InErrors     *prometheus.Desc
	UdpLite6OutDatagrams *prometheus.Desc
	UdpLite6RcvbufErrors *prometheus.Desc
	UdpLite6SndbufErrors *prometheus.Desc
	UdpLite6InCsumErrors *prometheus.Desc

	logger log.Logger
}

func init() {
	registerCollector("netstat", defaultEnabled, NewNetStatCollector)
}

// NewNetStatCollector takes and returns
// a new Collector exposing network stats.
func NewNetStatCollector(logger log.Logger) (Collector, error) {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}
	proc, err := fs.Self()
	if err != nil {
		return nil, fmt.Errorf("failed to open /proc/self: %w", err)
	}
	return &netStatCollector{
		proc,

		// TcpExt
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_SyncookiesSent"),
			"Statistic TcpExtSyncookiesSent.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_SyncookiesRecv"),
			"Statistic TcpExtSyncookiesRecv.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_SyncookiesFailed"),
			"Statistic TcpExtSyncookiesFailed.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_EmbryonicRsts"),
			"Statistic TcpExtEmbryonicRsts.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_PruneCalled"),
			"Statistic TcpExtPruneCalled.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_RcvPruned"),
			"Statistic TcpExtRcvPruned.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_OfoPruned"),
			"Statistic TcpExtOfoPruned.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_OutOfWindowIcmps"),
			"Statistic TcpExtOutOfWindowIcmps.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_LockDroppedIcmps"),
			"Statistic TcpExtLockDroppedIcmps.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_ArpFilter"),
			"Statistic TcpExtArpFilter.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TW"),
			"Statistic TcpExtTW.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TWRecycled"),
			"Statistic TcpExtTWRecycled.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TWKilled"),
			"Statistic TcpExtTWKilled.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_PAWSActive"),
			"Statistic TcpExtPAWSActive.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_PAWSEstab"),
			"Statistic TcpExtPAWSEstab.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_DelayedACKs"),
			"Statistic TcpExtDelayedACKs.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_DelayedACKLocked"),
			"Statistic TcpExtDelayedACKLocked.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_DelayedACKLost"),
			"Statistic TcpExtDelayedACKLost.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_ListenOverflows"),
			"Statistic TcpExtListenOverflows.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_ListenDrops"),
			"Statistic TcpExtListenDrops.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPHPHits"),
			"Statistic TcpExtTCPHPHits.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPPureAcks"),
			"Statistic TcpExtTCPPureAcks.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPHPAcks"),
			"Statistic TcpExtTCPHPAcks.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPRenoRecovery"),
			"Statistic TcpExtTCPRenoRecovery.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPSackRecovery"),
			"Statistic TcpExtTCPSackRecovery.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPSACKReneging"),
			"Statistic TcpExtTCPSACKReneging.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPSACKReorder"),
			"Statistic TcpExtTCPSACKReorder.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPRenoReorder"),
			"Statistic TcpExtTCPRenoReorder.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPTSReorder"),
			"Statistic TcpExtTCPTSReorder.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPFullUndo"),
			"Statistic TcpExtTCPFullUndo.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPPartialUndo"),
			"Statistic TcpExtTCPPartialUndo.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPDSACKUndo"),
			"Statistic TcpExtTCPDSACKUndo.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPLossUndo"),
			"Statistic TcpExtTCPLossUndo.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPLostRetransmit"),
			"Statistic TcpExtTCPLostRetransmit.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPRenoFailures"),
			"Statistic TcpExtTCPRenoFailures.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPSackFailures"),
			"Statistic TcpExtTCPSackFailures.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPLossFailures"),
			"Statistic TcpExtTCPLossFailures.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPFastRetrans"),
			"Statistic TcpExtTCPFastRetrans.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPSlowStartRetrans"),
			"Statistic TcpExtTCPSlowStartRetrans.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPTimeouts"),
			"Statistic TcpExtTCPTimeouts.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPLossProbes"),
			"Statistic TcpExtTCPLossProbes.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPLossProbeRecovery"),
			"Statistic TcpExtTCPLossProbeRecovery.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPRenoRecoveryFail"),
			"Statistic TcpExtTCPRenoRecoveryFail.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPSackRecoveryFail"),
			"Statistic TcpExtTCPSackRecoveryFail.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPRcvCollapsed"),
			"Statistic TcpExtTCPRcvCollapsed.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPDSACKOldSent"),
			"Statistic TcpExtTCPDSACKOldSent.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPDSACKOfoSent"),
			"Statistic TcpExtTCPDSACKOfoSent.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPDSACKRecv"),
			"Statistic TcpExtTCPDSACKRecv.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPDSACKOfoRecv"),
			"Statistic TcpExtTCPDSACKOfoRecv.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPAbortOnData"),
			"Statistic TcpExtTCPAbortOnData.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPAbortOnClose"),
			"Statistic TcpExtTCPAbortOnClose.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPAbortOnMemory"),
			"Statistic TcpExtTCPAbortOnMemory.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPAbortOnTimeout"),
			"Statistic TcpExtTCPAbortOnTimeout.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPAbortOnLinger"),
			"Statistic TcpExtTCPAbortOnLinger.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPAbortFailed"),
			"Statistic TcpExtTCPAbortFailed.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPMemoryPressures"),
			"Statistic TcpExtTCPMemoryPressures.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPMemoryPressuresChrono"),
			"Statistic TcpExtTCPMemoryPressuresChrono.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPSACKDiscard"),
			"Statistic TcpExtTCPSACKDiscard.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPDSACKIgnoredOld"),
			"Statistic TcpExtTCPDSACKIgnoredOld.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPDSACKIgnoredNoUndo"),
			"Statistic TcpExtTCPDSACKIgnoredNoUndo.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPSpuriousRTOs"),
			"Statistic TcpExtTCPSpuriousRTOs.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPMD5NotFound"),
			"Statistic TcpExtTCPMD5NotFound.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPMD5Unexpected"),
			"Statistic TcpExtTCPMD5Unexpected.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPMD5Failure"),
			"Statistic TcpExtTCPMD5Failure.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPSackShifted"),
			"Statistic TcpExtTCPSackShifted.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPSackMerged"),
			"Statistic TcpExtTCPSackMerged.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPSackShiftFallback"),
			"Statistic TcpExtTCPSackShiftFallback.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPBacklogDrop"),
			"Statistic TcpExtTCPBacklogDrop.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_PFMemallocDrop"),
			"Statistic TcpExtPFMemallocDrop.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPMinTTLDrop"),
			"Statistic TcpExtTCPMinTTLDrop.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPDeferAcceptDrop"),
			"Statistic TcpExtTCPDeferAcceptDrop.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_IPReversePathFilter"),
			"Statistic TcpExtIPReversePathFilter.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPTimeWaitOverflow"),
			"Statistic TcpExtTCPTimeWaitOverflow.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPReqQFullDoCookies"),
			"Statistic TcpExtTCPReqQFullDoCookies.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPReqQFullDrop"),
			"Statistic TcpExtTCPReqQFullDrop.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPRetransFail"),
			"Statistic TcpExtTCPRetransFail.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPRcvCoalesce"),
			"Statistic TcpExtTCPRcvCoalesce.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPOFOQueue"),
			"Statistic TcpExtTCPOFOQueue.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPOFODrop"),
			"Statistic TcpExtTCPOFODrop.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPOFOMerge"),
			"Statistic TcpExtTCPOFOMerge.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPChallengeACK"),
			"Statistic TcpExtTCPChallengeACK.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPSYNChallenge"),
			"Statistic TcpExtTCPSYNChallenge.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPFastOpenActive"),
			"Statistic TcpExtTCPFastOpenActive.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPFastOpenActiveFail"),
			"Statistic TcpExtTCPFastOpenActiveFail.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPFastOpenPassive"),
			"Statistic TcpExtTCPFastOpenPassive.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPFastOpenPassiveFail"),
			"Statistic TcpExtTCPFastOpenPassiveFail.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPFastOpenListenOverflow"),
			"Statistic TcpExtTCPFastOpenListenOverflow.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPFastOpenCookieReqd"),
			"Statistic TcpExtTCPFastOpenCookieReqd.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPFastOpenBlackhole"),
			"Statistic TcpExtTCPFastOpenBlackhole.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPSpuriousRtxHostQueues"),
			"Statistic TcpExtTCPSpuriousRtxHostQueues.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_BusyPollRxPackets"),
			"Statistic TcpExtBusyPollRxPackets.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPAutoCorking"),
			"Statistic TcpExtTCPAutoCorking.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPFromZeroWindowAdv"),
			"Statistic TcpExtTCPFromZeroWindowAdv.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPToZeroWindowAdv"),
			"Statistic TcpExtTCPToZeroWindowAdv.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPWantZeroWindowAdv"),
			"Statistic TcpExtTCPWantZeroWindowAdv.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPSynRetrans"),
			"Statistic TcpExtTCPSynRetrans.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPOrigDataSent"),
			"Statistic TcpExtTCPOrigDataSent.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPHystartTrainDetect"),
			"Statistic TcpExtTCPHystartTrainDetect.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPHystartTrainCwnd"),
			"Statistic TcpExtTCPHystartTrainCwnd.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPHystartDelayDetect"),
			"Statistic TcpExtTCPHystartDelayDetect.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPHystartDelayCwnd"),
			"Statistic TcpExtTCPHystartDelayCwnd.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPACKSkippedSynRecv"),
			"Statistic TcpExtTCPACKSkippedSynRecv.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPACKSkippedPAWS"),
			"Statistic TcpExtTCPACKSkippedPAWS.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPACKSkippedSeq"),
			"Statistic TcpExtTCPACKSkippedSeq.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPACKSkippedFinWait2"),
			"Statistic TcpExtTCPACKSkippedFinWait2.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPACKSkippedTimeWait"),
			"Statistic TcpExtTCPACKSkippedTimeWait.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPACKSkippedChallenge"),
			"Statistic TcpExtTCPACKSkippedChallenge.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPWinProbe"),
			"Statistic TcpExtTCPWinProbe.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPKeepAlive"),
			"Statistic TcpExtTCPKeepAlive.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPMTUPFail"),
			"Statistic TcpExtTCPMTUPFail.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPMTUPSuccess"),
			"Statistic TcpExtTCPMTUPSuccess.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "TcpExt_TCPWqueueTooBig"),
			"Statistic TcpExtTCPWqueueTooBig.",
			nil, nil,
		),

		// IpExt
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "IpExt_InNoRoutes"),
			"Statistic IpExtInNoRoutes.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "IpExt_InTruncatedPkts"),
			"Statistic IpExtInTruncatedPkts.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "IpExt_InMcastPkts"),
			"Statistic IpExtInMcastPkts.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "IpExt_OutMcastPkts"),
			"Statistic IpExtOutMcastPkts.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "IpExt_InBcastPkts"),
			"Statistic IpExtInBcastPkts.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "IpExt_OutBcastPkts"),
			"Statistic IpExtOutBcastPkts.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "IpExt_InOctets"),
			"Statistic IpExtInOctets.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "IpExt_OutOctets"),
			"Statistic IpExtOutOctets.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "IpExt_InMcastOctets"),
			"Statistic IpExtInMcastOctets.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "IpExt_OutMcastOctets"),
			"Statistic IpExtOutMcastOctets.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "IpExt_InBcastOctets"),
			"Statistic IpExtInBcastOctets.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "IpExt_OutBcastOctets"),
			"Statistic IpExtOutBcastOctets.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "IpExt_InCsumErrors"),
			"Statistic IpExtInCsumErrors.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "IpExt_InNoECTPkts"),
			"Statistic IpExtInNoECTPkts.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "IpExt_InECT1Pkts"),
			"Statistic IpExtInECT1Pkts.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "IpExt_InECT0Pkts"),
			"Statistic IpExtInECT0Pkts.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "IpExt_InCEPkts"),
			"Statistic IpExtInCEPkts.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "IpExt_ReasmOverlaps"),
			"Statistic IpExtReasmOverlaps.",
			nil, nil,
		),

		// Snmp
		// Ip
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip_Forwarding"),
			"Statistic IpForwarding.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip_DefaultTTL"),
			"Statistic IpDefaultTTL.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip_InReceives"),
			"Statistic IpInReceives.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip_InHdrErrors"),
			"Statistic IpInHdrErrors.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip_InAddrErrors"),
			"Statistic IpInAddrErrors.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip_ForwDatagrams"),
			"Statistic IpForwDatagrams.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip_InUnknownProtos"),
			"Statistic IpInUnknownProtos.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip_InDiscards"),
			"Statistic IpInDiscards.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip_InDelivers"),
			"Statistic IpInDelivers.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip_OutRequests"),
			"Statistic IpOutRequests.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip_OutDiscards"),
			"Statistic IpOutDiscards.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip_OutNoRoutes"),
			"Statistic IpOutNoRoutes.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip_ReasmTimeout"),
			"Statistic IpReasmTimeout.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip_ReasmReqds"),
			"Statistic IpReasmReqds.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip_ReasmOKs"),
			"Statistic IpReasmOKs.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip_ReasmFails"),
			"Statistic IpReasmFails.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip_FragOKs"),
			"Statistic IpFragOKs.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip_FragFails"),
			"Statistic IpFragFails.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip_FragCreates"),
			"Statistic IpFragCreates.",
			nil, nil,
		),

		// Icmp
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp_InMsgs"),
			"Statistic IcmpInMsgs.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp_InErrors"),
			"Statistic IcmpInErrors.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp_InCsumErrors"),
			"Statistic IcmpInCsumErrors.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp_InDestUnreachs"),
			"Statistic IcmpInDestUnreachs.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp_InTimeExcds"),
			"Statistic IcmpInTimeExcds.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp_InParmProbs"),
			"Statistic IcmpInParmProbs.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp_InSrcQuenchs"),
			"Statistic IcmpInSrcQuenchs.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp_InRedirects"),
			"Statistic IcmpInRedirects.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp_InEchos"),
			"Statistic IcmpInEchos.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp_InEchoReps"),
			"Statistic IcmpInEchoReps.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp_InTimestamps"),
			"Statistic IcmpInTimestamps.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp_InTimestampReps"),
			"Statistic IcmpInTimestampReps.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp_InAddrMasks"),
			"Statistic IcmpInAddrMasks.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp_InAddrMaskReps"),
			"Statistic IcmpInAddrMaskReps.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp_OutMsgs"),
			"Statistic IcmpOutMsgs.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp_OutErrors"),
			"Statistic IcmpOutErrors.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp_OutDestUnreachs"),
			"Statistic IcmpOutDestUnreachs.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp_OutTimeExcds"),
			"Statistic IcmpOutTimeExcds.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp_OutParmProbs"),
			"Statistic IcmpOutParmProbs.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp_OutSrcQuenchs"),
			"Statistic IcmpOutSrcQuenchs.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp_OutRedirects"),
			"Statistic IcmpOutRedirects.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp_OutEchos"),
			"Statistic IcmpOutEchos.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp_OutEchoReps"),
			"Statistic IcmpOutEchoReps.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp_OutTimestamps"),
			"Statistic IcmpOutTimestamps.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp_OutTimestampReps"),
			"Statistic IcmpOutTimestampReps.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp_OutAddrMasks"),
			"Statistic IcmpOutAddrMasks.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp_OutAddrMaskReps"),
			"Statistic IcmpOutAddrMaskReps.",
			nil, nil,
		),

		// IcmpMsg
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "IcmpMsg_InType3"),
			"Statistic IcmpMsgInType3.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "IcmpMsg_OutType3"),
			"Statistic IcmpMsgOutType3.",
			nil, nil,
		),

		// Tcp
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Tcp_RtoAlgorithm"),
			"Statistic TcpRtoAlgorithm.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Tcp_RtoMin"),
			"Statistic TcpRtoMin.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Tcp_RtoMax"),
			"Statistic TcpRtoMax.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Tcp_MaxConn"),
			"Statistic TcpMaxConn.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Tcp_ActiveOpens"),
			"Statistic TcpActiveOpens.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Tcp_PassiveOpens"),
			"Statistic TcpPassiveOpens.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Tcp_AttemptFails"),
			"Statistic TcpAttemptFails.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Tcp_EstabResets"),
			"Statistic TcpEstabResets.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Tcp_CurrEstab"),
			"Statistic TcpCurrEstab.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Tcp_InSegs"),
			"Statistic TcpInSegs.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Tcp_OutSegs"),
			"Statistic TcpOutSegs.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Tcp_RetransSegs"),
			"Statistic TcpRetransSegs.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Tcp_InErrs"),
			"Statistic TcpInErrs.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Tcp_OutRsts"),
			"Statistic TcpOutRsts.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Tcp_InCsumErrors"),
			"Statistic TcpInCsumErrors.",
			nil, nil,
		),

		// Udp
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Udp_InDatagrams"),
			"Statistic UdpInDatagrams.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Udp_NoPorts"),
			"Statistic UdpNoPorts.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Udp_InErrors"),
			"Statistic UdpInErrors.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Udp_OutDatagrams"),
			"Statistic UdpOutDatagrams.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Udp_RcvbufErrors"),
			"Statistic UdpRcvbufErrors.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Udp_SndbufErrors"),
			"Statistic UdpSndbufErrors.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Udp_InCsumErrors"),
			"Statistic UdpInCsumErrors.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Udp_IgnoredMulti"),
			"Statistic UdpIgnoredMulti.",
			nil, nil,
		),

		// UdpLite
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "UdpLite_InDatagrams"),
			"Statistic UdpLiteInDatagrams.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "UdpLite_NoPorts"),
			"Statistic UdpLiteNoPorts.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "UdpLite_InErrors"),
			"Statistic UdpLiteInErrors.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "UdpLite_OutDatagrams"),
			"Statistic UdpLiteOutDatagrams.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "UdpLite_RcvbufErrors"),
			"Statistic UdpLiteRcvbufErrors.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "UdpLite_SndbufErrors"),
			"Statistic UdpLiteSndbufErrors.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "UdpLite_InCsumErrors"),
			"Statistic UdpLiteInCsumErrors.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "UdpLite_IgnoredMulti"),
			"Statistic UdpLiteIgnoredMulti.",
			nil, nil,
		),

		// Snmp6
		// Ip6
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip6_InReceives"),
			"Statistic Ip6InReceives.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip6_InHdrErrors"),
			"Statistic Ip6InHdrErrors.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip6_InTooBigErrors"),
			"Statistic Ip6InTooBigErrors.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip6_InNoRoutes"),
			"Statistic Ip6InNoRoutes.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip6_InAddrErrors"),
			"Statistic Ip6InAddrErrors.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip6_InUnknownProtos"),
			"Statistic Ip6InUnknownProtos.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip6_InTruncatedPkts"),
			"Statistic Ip6InTruncatedPkts.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip6_InDiscards"),
			"Statistic Ip6InDiscards.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip6_InDelivers"),
			"Statistic Ip6InDelivers.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip6_OutForwDatagrams"),
			"Statistic Ip6OutForwDatagrams.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip6_OutRequests"),
			"Statistic Ip6OutRequests.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip6_OutDiscards"),
			"Statistic Ip6OutDiscards.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip6_OutNoRoutes"),
			"Statistic Ip6OutNoRoutes.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip6_ReasmTimeout"),
			"Statistic Ip6ReasmTimeout.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip6_ReasmReqds"),
			"Statistic Ip6ReasmReqds.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip6_ReasmOKs"),
			"Statistic Ip6ReasmOKs.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip6_ReasmFails"),
			"Statistic Ip6ReasmFails.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip6_FragOKs"),
			"Statistic Ip6FragOKs.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip6_FragFails"),
			"Statistic Ip6FragFails.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip6_FragCreates"),
			"Statistic Ip6FragCreates.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip6_InMcastPkts"),
			"Statistic Ip6InMcastPkts.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip6_OutMcastPkts"),
			"Statistic Ip6OutMcastPkts.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip6_InOctets"),
			"Statistic Ip6InOctets.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip6_OutOctets"),
			"Statistic Ip6OutOctets.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip6_InMcastOctets"),
			"Statistic Ip6InMcastOctets.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip6_OutMcastOctets"),
			"Statistic Ip6OutMcastOctets.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip6_InBcastOctets"),
			"Statistic Ip6InBcastOctets.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip6_OutBcastOctets"),
			"Statistic Ip6OutBcastOctets.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip6_InNoECTPkts"),
			"Statistic Ip6InNoECTPkts.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip6_InECT1Pkts"),
			"Statistic Ip6InECT1Pkts.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip6_InECT0Pkts"),
			"Statistic Ip6InECT0Pkts.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Ip6_InCEPkts"),
			"Statistic Ip6InCEPkts.",
			nil, nil,
		),

		// Icmp6
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_InMsgs"),
			"Statistic Icmp6InMsgs.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_InErrors"),
			"Statistic Icmp6InErrors.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_OutMsgs"),
			"Statistic Icmp6OutMsgs.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_OutErrors"),
			"Statistic Icmp6OutErrors.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_InCsumErrors"),
			"Statistic Icmp6InCsumErrors.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_InDestUnreachs"),
			"Statistic Icmp6InDestUnreachs.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_InPktTooBigs"),
			"Statistic Icmp6InPktTooBigs.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_InTimeExcds"),
			"Statistic Icmp6InTimeExcds.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_InParmProblems"),
			"Statistic Icmp6InParmProblems.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_InEchos"),
			"Statistic Icmp6InEchos.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_InEchoReplies"),
			"Statistic Icmp6InEchoReplies.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_InGroupMembQueries"),
			"Statistic Icmp6InGroupMembQueries.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_InGroupMembResponses"),
			"Statistic Icmp6InGroupMembResponses.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_InGroupMembReductions"),
			"Statistic Icmp6InGroupMembReductions.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_InRouterSolicits"),
			"Statistic Icmp6InRouterSolicits.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_InRouterAdvertisements"),
			"Statistic Icmp6InRouterAdvertisements.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_InNeighborSolicits"),
			"Statistic Icmp6InNeighborSolicits.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_InNeighborAdvertisements"),
			"Statistic Icmp6InNeighborAdvertisements.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_InRedirects"),
			"Statistic Icmp6InRedirects.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_InMLDv2Reports"),
			"Statistic Icmp6InMLDv2Reports.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_OutDestUnreachs"),
			"Statistic Icmp6OutDestUnreachs.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_OutPktTooBigs"),
			"Statistic Icmp6OutPktTooBigs.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_OutTimeExcds"),
			"Statistic Icmp6OutTimeExcds.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_OutParmProblems"),
			"Statistic Icmp6OutParmProblems.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_OutEchos"),
			"Statistic Icmp6OutEchos.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_OutEchoReplies"),
			"Statistic Icmp6OutEchoReplies.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_OutGroupMembQueries"),
			"Statistic Icmp6OutGroupMembQueries.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_OutGroupMembResponses"),
			"Statistic Icmp6OutGroupMembResponses.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_OutGroupMembReductions"),
			"Statistic Icmp6OutGroupMembReductions.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_OutRouterSolicits"),
			"Statistic Icmp6OutRouterSolicits.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_OutRouterAdvertisements"),
			"Statistic Icmp6OutRouterAdvertisements.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_OutNeighborSolicits"),
			"Statistic Icmp6OutNeighborSolicits.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_OutNeighborAdvertisements"),
			"Statistic Icmp6OutNeighborAdvertisements.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_OutRedirects"),
			"Statistic Icmp6OutRedirects.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_OutMLDv2Reports"),
			"Statistic Icmp6OutMLDv2Reports.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_InType1"),
			"Statistic Icmp6InType1.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_InType134"),
			"Statistic Icmp6InType134.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_InType135"),
			"Statistic Icmp6InType135.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_InType136"),
			"Statistic Icmp6InType136.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_InType143"),
			"Statistic Icmp6InType143.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_OutType133"),
			"Statistic Icmp6OutType133.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_OutType135"),
			"Statistic Icmp6OutType135.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_OutType136"),
			"Statistic Icmp6OutType136.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Icmp6_OutType143"),
			"Statistic Icmp6OutType143.",
			nil, nil,
		),

		// Udp6
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Udp6_InDatagrams"),
			"Statistic Udp6InDatagrams.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Udp6_NoPorts"),
			"Statistic Udp6NoPorts.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Udp6_InErrors"),
			"Statistic Udp6InErrors.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Udp6_OutDatagrams"),
			"Statistic Udp6OutDatagrams.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Udp6_RcvbufErrors"),
			"Statistic Udp6RcvbufErrors.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Udp6_SndbufErrors"),
			"Statistic Udp6SndbufErrors.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Udp6_InCsumErrors"),
			"Statistic Udp6InCsumErrors.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "Udp6_IgnoredMulti"),
			"Statistic Udp6IgnoredMulti.",
			nil, nil,
		),

		// UdpLite6
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "UdpLite6_InDatagrams"),
			"Statistic UdpLite6InDatagrams.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "UdpLite6_NoPorts"),
			"Statistic UdpLite6NoPorts.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "UdpLite6_InErrors"),
			"Statistic UdpLite6InErrors.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "UdpLite6_OutDatagrams"),
			"Statistic UdpLite6OutDatagrams.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "UdpLite6_RcvbufErrors"),
			"Statistic UdpLite6RcvbufErrors.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "UdpLite6_SndbufErrors"),
			"Statistic UdpLite6SndbufErrors.",
			nil, nil,
		),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netStatsSubsystem, "UdpLite6_InCsumErrors"),
			"Statistic UdpLite6InCsumErrors.",
			nil, nil,
		),
		logger,
	}, nil
}

func (c *netStatCollector) Update(ch chan<- prometheus.Metric) error {
	if err := c.updateNetstat(ch); err != nil {
		return err
	}

	if err := c.updateSnmp(ch); err != nil {
		return err
	}

	if err := c.updateSnmp6(ch); err != nil {
		return err
	}

	return nil
}

func (c *netStatCollector) updateNetstat(ch chan<- prometheus.Metric) error {
	procNetstat, err := c.proc.Netstat()
	if err != nil {
		return err
	}

	// TcpExt
	ch <- prometheus.MustNewConstMetric(c.TcpExtSyncookiesSent,
		prometheus.UntypedValue, procNetstat.TcpExt.SyncookiesSent,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtSyncookiesRecv,
		prometheus.UntypedValue, procNetstat.TcpExt.SyncookiesRecv,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtSyncookiesFailed,
		prometheus.UntypedValue, procNetstat.TcpExt.SyncookiesFailed,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtEmbryonicRsts,
		prometheus.UntypedValue, procNetstat.TcpExt.EmbryonicRsts,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtPruneCalled,
		prometheus.UntypedValue, procNetstat.TcpExt.PruneCalled,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtRcvPruned,
		prometheus.UntypedValue, procNetstat.TcpExt.RcvPruned,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtOfoPruned,
		prometheus.UntypedValue, procNetstat.TcpExt.OfoPruned,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtOutOfWindowIcmps,
		prometheus.UntypedValue, procNetstat.TcpExt.OutOfWindowIcmps,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtLockDroppedIcmps,
		prometheus.UntypedValue, procNetstat.TcpExt.LockDroppedIcmps,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtArpFilter,
		prometheus.UntypedValue, procNetstat.TcpExt.ArpFilter,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTW,
		prometheus.UntypedValue, procNetstat.TcpExt.TW,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTWRecycled,
		prometheus.UntypedValue, procNetstat.TcpExt.TWRecycled,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTWKilled,
		prometheus.UntypedValue, procNetstat.TcpExt.TWKilled,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtPAWSActive,
		prometheus.UntypedValue, procNetstat.TcpExt.PAWSActive,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtPAWSEstab,
		prometheus.UntypedValue, procNetstat.TcpExt.PAWSEstab,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtDelayedACKs,
		prometheus.UntypedValue, procNetstat.TcpExt.DelayedACKs,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtDelayedACKLocked,
		prometheus.UntypedValue, procNetstat.TcpExt.DelayedACKLocked,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtDelayedACKLost,
		prometheus.UntypedValue, procNetstat.TcpExt.DelayedACKLost,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtListenOverflows,
		prometheus.UntypedValue, procNetstat.TcpExt.ListenOverflows,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtListenDrops,
		prometheus.UntypedValue, procNetstat.TcpExt.ListenDrops,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPHPHits,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPHPHits,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPPureAcks,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPPureAcks,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPHPAcks,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPHPAcks,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPRenoRecovery,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPRenoRecovery,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPSackRecovery,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPSackRecovery,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPSACKReneging,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPSACKReneging,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPSACKReorder,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPSACKReorder,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPRenoReorder,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPRenoReorder,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPTSReorder,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPTSReorder,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPFullUndo,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPFullUndo,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPPartialUndo,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPPartialUndo,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPDSACKUndo,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPDSACKUndo,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPLossUndo,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPLossUndo,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPLostRetransmit,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPLostRetransmit,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPRenoFailures,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPRenoFailures,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPSackFailures,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPSackFailures,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPLossFailures,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPLossFailures,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPFastRetrans,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPFastRetrans,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPSlowStartRetrans,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPSlowStartRetrans,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPTimeouts,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPTimeouts,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPLossProbes,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPLossProbes,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPLossProbeRecovery,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPLossProbeRecovery,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPRenoRecoveryFail,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPRenoRecoveryFail,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPSackRecoveryFail,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPSackRecoveryFail,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPRcvCollapsed,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPRcvCollapsed,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPDSACKOldSent,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPDSACKOldSent,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPDSACKOfoSent,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPDSACKOfoSent,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPDSACKRecv,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPDSACKRecv,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPDSACKOfoRecv,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPDSACKOfoRecv,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPAbortOnData,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPAbortOnData,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPAbortOnClose,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPAbortOnClose,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPAbortOnMemory,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPAbortOnMemory,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPAbortOnTimeout,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPAbortOnTimeout,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPAbortOnLinger,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPAbortOnLinger,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPAbortFailed,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPAbortFailed,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPMemoryPressures,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPMemoryPressures,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPMemoryPressuresChrono,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPMemoryPressuresChrono,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPSACKDiscard,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPSACKDiscard,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPDSACKIgnoredOld,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPDSACKIgnoredOld,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPDSACKIgnoredNoUndo,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPDSACKIgnoredNoUndo,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPSpuriousRTOs,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPSpuriousRTOs,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPMD5NotFound,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPMD5NotFound,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPMD5Unexpected,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPMD5Unexpected,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPMD5Failure,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPMD5Failure,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPSackShifted,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPSackShifted,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPSackMerged,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPSackMerged,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPSackShiftFallback,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPSackShiftFallback,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPBacklogDrop,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPBacklogDrop,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtPFMemallocDrop,
		prometheus.UntypedValue, procNetstat.TcpExt.PFMemallocDrop,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPMinTTLDrop,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPMinTTLDrop,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPDeferAcceptDrop,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPDeferAcceptDrop,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtIPReversePathFilter,
		prometheus.UntypedValue, procNetstat.TcpExt.IPReversePathFilter,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPTimeWaitOverflow,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPTimeWaitOverflow,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPReqQFullDoCookies,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPReqQFullDoCookies,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPReqQFullDrop,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPReqQFullDrop,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPRetransFail,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPRetransFail,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPRcvCoalesce,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPRcvCoalesce,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPOFOQueue,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPOFOQueue,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPOFODrop,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPOFODrop,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPOFOMerge,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPOFOMerge,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPChallengeACK,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPChallengeACK,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPSYNChallenge,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPSYNChallenge,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPFastOpenActive,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPFastOpenActive,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPFastOpenActiveFail,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPFastOpenActiveFail,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPFastOpenPassive,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPFastOpenPassive,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPFastOpenPassiveFail,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPFastOpenPassiveFail,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPFastOpenListenOverflow,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPFastOpenListenOverflow,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPFastOpenCookieReqd,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPFastOpenCookieReqd,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPFastOpenBlackhole,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPFastOpenBlackhole,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPSpuriousRtxHostQueues,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPSpuriousRtxHostQueues,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtBusyPollRxPackets,
		prometheus.UntypedValue, procNetstat.TcpExt.BusyPollRxPackets,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPAutoCorking,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPAutoCorking,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPFromZeroWindowAdv,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPFromZeroWindowAdv,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPToZeroWindowAdv,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPToZeroWindowAdv,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPWantZeroWindowAdv,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPWantZeroWindowAdv,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPSynRetrans,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPSynRetrans,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPOrigDataSent,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPOrigDataSent,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPHystartTrainDetect,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPHystartTrainDetect,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPHystartTrainCwnd,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPHystartTrainCwnd,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPHystartDelayDetect,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPHystartDelayDetect,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPHystartDelayCwnd,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPHystartDelayCwnd,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPACKSkippedSynRecv,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPACKSkippedSynRecv,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPACKSkippedPAWS,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPACKSkippedPAWS,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPACKSkippedSeq,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPACKSkippedSeq,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPACKSkippedFinWait2,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPACKSkippedFinWait2,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPACKSkippedTimeWait,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPACKSkippedTimeWait,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPACKSkippedChallenge,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPACKSkippedChallenge,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPWinProbe,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPWinProbe,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPKeepAlive,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPKeepAlive,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPMTUPFail,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPMTUPFail,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPMTUPSuccess,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPMTUPSuccess,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpExtTCPWqueueTooBig,
		prometheus.UntypedValue, procNetstat.TcpExt.TCPWqueueTooBig,
	)

	// IpExt
	ch <- prometheus.MustNewConstMetric(c.IpExtInNoRoutes,
		prometheus.UntypedValue, procNetstat.IpExt.InNoRoutes,
	)
	ch <- prometheus.MustNewConstMetric(c.IpExtInTruncatedPkts,
		prometheus.UntypedValue, procNetstat.IpExt.InTruncatedPkts,
	)
	ch <- prometheus.MustNewConstMetric(c.IpExtInMcastPkts,
		prometheus.UntypedValue, procNetstat.IpExt.InMcastPkts,
	)
	ch <- prometheus.MustNewConstMetric(c.IpExtOutMcastPkts,
		prometheus.UntypedValue, procNetstat.IpExt.OutMcastPkts,
	)
	ch <- prometheus.MustNewConstMetric(c.IpExtInBcastPkts,
		prometheus.UntypedValue, procNetstat.IpExt.InBcastPkts,
	)
	ch <- prometheus.MustNewConstMetric(c.IpExtOutBcastPkts,
		prometheus.UntypedValue, procNetstat.IpExt.OutBcastPkts,
	)
	ch <- prometheus.MustNewConstMetric(c.IpExtInOctets,
		prometheus.UntypedValue, procNetstat.IpExt.InOctets,
	)
	ch <- prometheus.MustNewConstMetric(c.IpExtOutOctets,
		prometheus.UntypedValue, procNetstat.IpExt.OutOctets,
	)
	ch <- prometheus.MustNewConstMetric(c.IpExtInMcastOctets,
		prometheus.UntypedValue, procNetstat.IpExt.InMcastOctets,
	)
	ch <- prometheus.MustNewConstMetric(c.IpExtOutMcastOctets,
		prometheus.UntypedValue, procNetstat.IpExt.OutMcastOctets,
	)
	ch <- prometheus.MustNewConstMetric(c.IpExtInBcastOctets,
		prometheus.UntypedValue, procNetstat.IpExt.InBcastOctets,
	)
	ch <- prometheus.MustNewConstMetric(c.IpExtOutBcastOctets,
		prometheus.UntypedValue, procNetstat.IpExt.OutBcastOctets,
	)
	ch <- prometheus.MustNewConstMetric(c.IpExtInCsumErrors,
		prometheus.UntypedValue, procNetstat.IpExt.InCsumErrors,
	)
	ch <- prometheus.MustNewConstMetric(c.IpExtInNoECTPkts,
		prometheus.UntypedValue, procNetstat.IpExt.InNoECTPkts,
	)
	ch <- prometheus.MustNewConstMetric(c.IpExtInECT1Pkts,
		prometheus.UntypedValue, procNetstat.IpExt.InECT1Pkts,
	)
	ch <- prometheus.MustNewConstMetric(c.IpExtInECT0Pkts,
		prometheus.UntypedValue, procNetstat.IpExt.InECT0Pkts,
	)
	ch <- prometheus.MustNewConstMetric(c.IpExtInCEPkts,
		prometheus.UntypedValue, procNetstat.IpExt.InCEPkts,
	)
	ch <- prometheus.MustNewConstMetric(c.IpExtReasmOverlaps,
		prometheus.UntypedValue, procNetstat.IpExt.ReasmOverlaps,
	)

	return nil
}

func (c *netStatCollector) updateSnmp(ch chan<- prometheus.Metric) error {
	snmp, err := c.proc.Snmp()
	if err != nil {
		return err
	}

	// IP
	ch <- prometheus.MustNewConstMetric(c.IpForwarding,
		prometheus.UntypedValue, snmp.Ip.Forwarding,
	)
	ch <- prometheus.MustNewConstMetric(c.IpDefaultTTL,
		prometheus.UntypedValue, snmp.Ip.DefaultTTL,
	)
	ch <- prometheus.MustNewConstMetric(c.IpInReceives,
		prometheus.UntypedValue, snmp.Ip.InReceives,
	)
	ch <- prometheus.MustNewConstMetric(c.IpInHdrErrors,
		prometheus.UntypedValue, snmp.Ip.InHdrErrors,
	)
	ch <- prometheus.MustNewConstMetric(c.IpInAddrErrors,
		prometheus.UntypedValue, snmp.Ip.InAddrErrors,
	)
	ch <- prometheus.MustNewConstMetric(c.IpForwDatagrams,
		prometheus.UntypedValue, snmp.Ip.ForwDatagrams,
	)
	ch <- prometheus.MustNewConstMetric(c.IpInUnknownProtos,
		prometheus.UntypedValue, snmp.Ip.InUnknownProtos,
	)
	ch <- prometheus.MustNewConstMetric(c.IpInDiscards,
		prometheus.UntypedValue, snmp.Ip.InDiscards,
	)
	ch <- prometheus.MustNewConstMetric(c.IpInDelivers,
		prometheus.UntypedValue, snmp.Ip.InDelivers,
	)
	ch <- prometheus.MustNewConstMetric(c.IpOutRequests,
		prometheus.UntypedValue, snmp.Ip.OutRequests,
	)
	ch <- prometheus.MustNewConstMetric(c.IpOutDiscards,
		prometheus.UntypedValue, snmp.Ip.OutDiscards,
	)
	ch <- prometheus.MustNewConstMetric(c.IpOutNoRoutes,
		prometheus.UntypedValue, snmp.Ip.OutNoRoutes,
	)
	ch <- prometheus.MustNewConstMetric(c.IpReasmTimeout,
		prometheus.UntypedValue, snmp.Ip.ReasmTimeout,
	)
	ch <- prometheus.MustNewConstMetric(c.IpReasmReqds,
		prometheus.UntypedValue, snmp.Ip.ReasmReqds,
	)
	ch <- prometheus.MustNewConstMetric(c.IpReasmOKs,
		prometheus.UntypedValue, snmp.Ip.ReasmOKs,
	)
	ch <- prometheus.MustNewConstMetric(c.IpReasmFails,
		prometheus.UntypedValue, snmp.Ip.ReasmFails,
	)
	ch <- prometheus.MustNewConstMetric(c.IpFragOKs,
		prometheus.UntypedValue, snmp.Ip.FragOKs,
	)
	ch <- prometheus.MustNewConstMetric(c.IpFragFails,
		prometheus.UntypedValue, snmp.Ip.FragFails,
	)
	ch <- prometheus.MustNewConstMetric(c.IpFragCreates,
		prometheus.UntypedValue, snmp.Ip.FragCreates,
	)

	// Icmp
	ch <- prometheus.MustNewConstMetric(c.IcmpInMsgs,
		prometheus.UntypedValue, snmp.Icmp.InMsgs,
	)
	ch <- prometheus.MustNewConstMetric(c.IcmpInErrors,
		prometheus.UntypedValue, snmp.Icmp.InErrors,
	)
	ch <- prometheus.MustNewConstMetric(c.IcmpInCsumErrors,
		prometheus.UntypedValue, snmp.Icmp.InCsumErrors,
	)
	ch <- prometheus.MustNewConstMetric(c.IcmpInDestUnreachs,
		prometheus.UntypedValue, snmp.Icmp.InDestUnreachs,
	)
	ch <- prometheus.MustNewConstMetric(c.IcmpInTimeExcds,
		prometheus.UntypedValue, snmp.Icmp.InTimeExcds,
	)
	ch <- prometheus.MustNewConstMetric(c.IcmpInParmProbs,
		prometheus.UntypedValue, snmp.Icmp.InParmProbs,
	)
	ch <- prometheus.MustNewConstMetric(c.IcmpInSrcQuenchs,
		prometheus.UntypedValue, snmp.Icmp.InSrcQuenchs,
	)
	ch <- prometheus.MustNewConstMetric(c.IcmpInRedirects,
		prometheus.UntypedValue, snmp.Icmp.InRedirects,
	)
	ch <- prometheus.MustNewConstMetric(c.IcmpInEchos,
		prometheus.UntypedValue, snmp.Icmp.InEchos,
	)
	ch <- prometheus.MustNewConstMetric(c.IcmpInEchoReps,
		prometheus.UntypedValue, snmp.Icmp.InEchoReps,
	)
	ch <- prometheus.MustNewConstMetric(c.IcmpInTimestamps,
		prometheus.UntypedValue, snmp.Icmp.InTimestamps,
	)
	ch <- prometheus.MustNewConstMetric(c.IcmpInTimestampReps,
		prometheus.UntypedValue, snmp.Icmp.InTimestampReps,
	)
	ch <- prometheus.MustNewConstMetric(c.IcmpInAddrMasks,
		prometheus.UntypedValue, snmp.Icmp.InAddrMasks,
	)
	ch <- prometheus.MustNewConstMetric(c.IcmpInAddrMaskReps,
		prometheus.UntypedValue, snmp.Icmp.InAddrMaskReps,
	)
	ch <- prometheus.MustNewConstMetric(c.IcmpOutMsgs,
		prometheus.UntypedValue, snmp.Icmp.OutMsgs,
	)
	ch <- prometheus.MustNewConstMetric(c.IcmpOutErrors,
		prometheus.UntypedValue, snmp.Icmp.OutErrors,
	)
	ch <- prometheus.MustNewConstMetric(c.IcmpOutDestUnreachs,
		prometheus.UntypedValue, snmp.Icmp.OutDestUnreachs,
	)
	ch <- prometheus.MustNewConstMetric(c.IcmpOutTimeExcds,
		prometheus.UntypedValue, snmp.Icmp.OutTimeExcds,
	)
	ch <- prometheus.MustNewConstMetric(c.IcmpOutParmProbs,
		prometheus.UntypedValue, snmp.Icmp.OutParmProbs,
	)
	ch <- prometheus.MustNewConstMetric(c.IcmpOutSrcQuenchs,
		prometheus.UntypedValue, snmp.Icmp.OutSrcQuenchs,
	)
	ch <- prometheus.MustNewConstMetric(c.IcmpOutRedirects,
		prometheus.UntypedValue, snmp.Icmp.OutRedirects,
	)
	ch <- prometheus.MustNewConstMetric(c.IcmpOutEchos,
		prometheus.UntypedValue, snmp.Icmp.OutEchos,
	)
	ch <- prometheus.MustNewConstMetric(c.IcmpOutEchoReps,
		prometheus.UntypedValue, snmp.Icmp.OutEchoReps,
	)
	ch <- prometheus.MustNewConstMetric(c.IcmpOutTimestamps,
		prometheus.UntypedValue, snmp.Icmp.OutTimestamps,
	)
	ch <- prometheus.MustNewConstMetric(c.IcmpOutTimestampReps,
		prometheus.UntypedValue, snmp.Icmp.OutTimestampReps,
	)
	ch <- prometheus.MustNewConstMetric(c.IcmpOutAddrMasks,
		prometheus.UntypedValue, snmp.Icmp.OutAddrMasks,
	)
	ch <- prometheus.MustNewConstMetric(c.IcmpOutAddrMaskReps,
		prometheus.UntypedValue, snmp.Icmp.OutAddrMaskReps,
	)

	// IcmpMsg
	ch <- prometheus.MustNewConstMetric(c.IcmpMsgInType3,
		prometheus.UntypedValue, snmp.IcmpMsg.InType3,
	)
	ch <- prometheus.MustNewConstMetric(c.IcmpMsgOutType3,
		prometheus.UntypedValue, snmp.IcmpMsg.OutType3,
	)

	// Tcp
	ch <- prometheus.MustNewConstMetric(c.TcpRtoAlgorithm,
		prometheus.UntypedValue, snmp.Tcp.RtoAlgorithm,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpRtoMin,
		prometheus.UntypedValue, snmp.Tcp.RtoMin,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpRtoMax,
		prometheus.UntypedValue, snmp.Tcp.RtoMax,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpMaxConn,
		prometheus.UntypedValue, snmp.Tcp.MaxConn,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpActiveOpens,
		prometheus.UntypedValue, snmp.Tcp.ActiveOpens,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpPassiveOpens,
		prometheus.UntypedValue, snmp.Tcp.PassiveOpens,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpAttemptFails,
		prometheus.UntypedValue, snmp.Tcp.AttemptFails,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpEstabResets,
		prometheus.UntypedValue, snmp.Tcp.EstabResets,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpCurrEstab,
		prometheus.UntypedValue, snmp.Tcp.CurrEstab,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpInSegs,
		prometheus.UntypedValue, snmp.Tcp.InSegs,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpOutSegs,
		prometheus.UntypedValue, snmp.Tcp.OutSegs,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpRetransSegs,
		prometheus.UntypedValue, snmp.Tcp.RetransSegs,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpInErrs,
		prometheus.UntypedValue, snmp.Tcp.InErrs,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpOutRsts,
		prometheus.UntypedValue, snmp.Tcp.OutRsts,
	)
	ch <- prometheus.MustNewConstMetric(c.TcpInCsumErrors,
		prometheus.UntypedValue, snmp.Tcp.InCsumErrors,
	)

	// Udp
	ch <- prometheus.MustNewConstMetric(c.UdpInDatagrams,
		prometheus.UntypedValue, snmp.Udp.InDatagrams,
	)
	ch <- prometheus.MustNewConstMetric(c.UdpNoPorts,
		prometheus.UntypedValue, snmp.Udp.NoPorts,
	)
	ch <- prometheus.MustNewConstMetric(c.UdpInErrors,
		prometheus.UntypedValue, snmp.Udp.InErrors,
	)
	ch <- prometheus.MustNewConstMetric(c.UdpOutDatagrams,
		prometheus.UntypedValue, snmp.Udp.OutDatagrams,
	)
	ch <- prometheus.MustNewConstMetric(c.UdpRcvbufErrors,
		prometheus.UntypedValue, snmp.Udp.RcvbufErrors,
	)
	ch <- prometheus.MustNewConstMetric(c.UdpSndbufErrors,
		prometheus.UntypedValue, snmp.Udp.SndbufErrors,
	)
	ch <- prometheus.MustNewConstMetric(c.UdpInCsumErrors,
		prometheus.UntypedValue, snmp.Udp.InCsumErrors,
	)
	ch <- prometheus.MustNewConstMetric(c.UdpIgnoredMulti,
		prometheus.UntypedValue, snmp.Udp.IgnoredMulti,
	)

	// UdpLite
	ch <- prometheus.MustNewConstMetric(c.UdpLiteInDatagrams,
		prometheus.UntypedValue, snmp.UdpLite.InDatagrams,
	)
	ch <- prometheus.MustNewConstMetric(c.UdpLiteNoPorts,
		prometheus.UntypedValue, snmp.UdpLite.NoPorts,
	)
	ch <- prometheus.MustNewConstMetric(c.UdpLiteInErrors,
		prometheus.UntypedValue, snmp.UdpLite.InErrors,
	)
	ch <- prometheus.MustNewConstMetric(c.UdpLiteOutDatagrams,
		prometheus.UntypedValue, snmp.UdpLite.OutDatagrams,
	)
	ch <- prometheus.MustNewConstMetric(c.UdpLiteRcvbufErrors,
		prometheus.UntypedValue, snmp.UdpLite.RcvbufErrors,
	)
	ch <- prometheus.MustNewConstMetric(c.UdpLiteSndbufErrors,
		prometheus.UntypedValue, snmp.UdpLite.SndbufErrors,
	)
	ch <- prometheus.MustNewConstMetric(c.UdpLiteInCsumErrors,
		prometheus.UntypedValue, snmp.UdpLite.InCsumErrors,
	)
	ch <- prometheus.MustNewConstMetric(c.UdpLiteIgnoredMulti,
		prometheus.UntypedValue, snmp.UdpLite.IgnoredMulti,
	)

	return nil
}

func (c *netStatCollector) updateSnmp6(ch chan<- prometheus.Metric) error {
	// SNMP6
	snmp6, err := c.proc.Snmp6()
	if err != nil {
		return err
	}

	// Ip6
	ch <- prometheus.MustNewConstMetric(c.Ip6InReceives,
		prometheus.UntypedValue, snmp6.Ip6.InReceives,
	)
	ch <- prometheus.MustNewConstMetric(c.Ip6InHdrErrors,
		prometheus.UntypedValue, snmp6.Ip6.InHdrErrors,
	)
	ch <- prometheus.MustNewConstMetric(c.Ip6InTooBigErrors,
		prometheus.UntypedValue, snmp6.Ip6.InTooBigErrors,
	)
	ch <- prometheus.MustNewConstMetric(c.Ip6InNoRoutes,
		prometheus.UntypedValue, snmp6.Ip6.InNoRoutes,
	)
	ch <- prometheus.MustNewConstMetric(c.Ip6InAddrErrors,
		prometheus.UntypedValue, snmp6.Ip6.InAddrErrors,
	)
	ch <- prometheus.MustNewConstMetric(c.Ip6InUnknownProtos,
		prometheus.UntypedValue, snmp6.Ip6.InUnknownProtos,
	)
	ch <- prometheus.MustNewConstMetric(c.Ip6InTruncatedPkts,
		prometheus.UntypedValue, snmp6.Ip6.InTruncatedPkts,
	)
	ch <- prometheus.MustNewConstMetric(c.Ip6InDiscards,
		prometheus.UntypedValue, snmp6.Ip6.InDiscards,
	)
	ch <- prometheus.MustNewConstMetric(c.Ip6InDelivers,
		prometheus.UntypedValue, snmp6.Ip6.InDelivers,
	)
	ch <- prometheus.MustNewConstMetric(c.Ip6OutForwDatagrams,
		prometheus.UntypedValue, snmp6.Ip6.OutForwDatagrams,
	)
	ch <- prometheus.MustNewConstMetric(c.Ip6OutRequests,
		prometheus.UntypedValue, snmp6.Ip6.OutRequests,
	)
	ch <- prometheus.MustNewConstMetric(c.Ip6OutDiscards,
		prometheus.UntypedValue, snmp6.Ip6.OutDiscards,
	)
	ch <- prometheus.MustNewConstMetric(c.Ip6OutNoRoutes,
		prometheus.UntypedValue, snmp6.Ip6.OutNoRoutes,
	)
	ch <- prometheus.MustNewConstMetric(c.Ip6ReasmTimeout,
		prometheus.UntypedValue, snmp6.Ip6.ReasmTimeout,
	)
	ch <- prometheus.MustNewConstMetric(c.Ip6ReasmReqds,
		prometheus.UntypedValue, snmp6.Ip6.ReasmReqds,
	)
	ch <- prometheus.MustNewConstMetric(c.Ip6ReasmOKs,
		prometheus.UntypedValue, snmp6.Ip6.ReasmOKs,
	)
	ch <- prometheus.MustNewConstMetric(c.Ip6ReasmFails,
		prometheus.UntypedValue, snmp6.Ip6.ReasmFails,
	)
	ch <- prometheus.MustNewConstMetric(c.Ip6FragOKs,
		prometheus.UntypedValue, snmp6.Ip6.FragOKs,
	)
	ch <- prometheus.MustNewConstMetric(c.Ip6FragFails,
		prometheus.UntypedValue, snmp6.Ip6.FragFails,
	)
	ch <- prometheus.MustNewConstMetric(c.Ip6FragCreates,
		prometheus.UntypedValue, snmp6.Ip6.FragCreates,
	)
	ch <- prometheus.MustNewConstMetric(c.Ip6InMcastPkts,
		prometheus.UntypedValue, snmp6.Ip6.InMcastPkts,
	)
	ch <- prometheus.MustNewConstMetric(c.Ip6OutMcastPkts,
		prometheus.UntypedValue, snmp6.Ip6.OutMcastPkts,
	)
	ch <- prometheus.MustNewConstMetric(c.Ip6InOctets,
		prometheus.UntypedValue, snmp6.Ip6.InOctets,
	)
	ch <- prometheus.MustNewConstMetric(c.Ip6OutOctets,
		prometheus.UntypedValue, snmp6.Ip6.OutOctets,
	)
	ch <- prometheus.MustNewConstMetric(c.Ip6InMcastOctets,
		prometheus.UntypedValue, snmp6.Ip6.InMcastOctets,
	)
	ch <- prometheus.MustNewConstMetric(c.Ip6OutMcastOctets,
		prometheus.UntypedValue, snmp6.Ip6.OutMcastOctets,
	)
	ch <- prometheus.MustNewConstMetric(c.Ip6InBcastOctets,
		prometheus.UntypedValue, snmp6.Ip6.InBcastOctets,
	)
	ch <- prometheus.MustNewConstMetric(c.Ip6OutBcastOctets,
		prometheus.UntypedValue, snmp6.Ip6.OutBcastOctets,
	)
	ch <- prometheus.MustNewConstMetric(c.Ip6InNoECTPkts,
		prometheus.UntypedValue, snmp6.Ip6.InNoECTPkts,
	)
	ch <- prometheus.MustNewConstMetric(c.Ip6InECT1Pkts,
		prometheus.UntypedValue, snmp6.Ip6.InECT1Pkts,
	)
	ch <- prometheus.MustNewConstMetric(c.Ip6InECT0Pkts,
		prometheus.UntypedValue, snmp6.Ip6.InECT0Pkts,
	)
	ch <- prometheus.MustNewConstMetric(c.Ip6InCEPkts,
		prometheus.UntypedValue, snmp6.Ip6.InCEPkts,
	)

	// Icmp6
	ch <- prometheus.MustNewConstMetric(c.Icmp6InMsgs,
		prometheus.UntypedValue, snmp6.Icmp6.InMsgs,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6InErrors,
		prometheus.UntypedValue, snmp6.Icmp6.InErrors,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6OutMsgs,
		prometheus.UntypedValue, snmp6.Icmp6.OutMsgs,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6OutErrors,
		prometheus.UntypedValue, snmp6.Icmp6.OutErrors,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6InCsumErrors,
		prometheus.UntypedValue, snmp6.Icmp6.InCsumErrors,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6InDestUnreachs,
		prometheus.UntypedValue, snmp6.Icmp6.InDestUnreachs,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6InPktTooBigs,
		prometheus.UntypedValue, snmp6.Icmp6.InPktTooBigs,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6InTimeExcds,
		prometheus.UntypedValue, snmp6.Icmp6.InTimeExcds,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6InParmProblems,
		prometheus.UntypedValue, snmp6.Icmp6.InParmProblems,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6InEchos,
		prometheus.UntypedValue, snmp6.Icmp6.InEchos,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6InEchoReplies,
		prometheus.UntypedValue, snmp6.Icmp6.InEchoReplies,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6InGroupMembQueries,
		prometheus.UntypedValue, snmp6.Icmp6.InGroupMembQueries,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6InGroupMembResponses,
		prometheus.UntypedValue, snmp6.Icmp6.InGroupMembResponses,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6InGroupMembReductions,
		prometheus.UntypedValue, snmp6.Icmp6.InGroupMembReductions,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6InRouterSolicits,
		prometheus.UntypedValue, snmp6.Icmp6.InRouterSolicits,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6InRouterAdvertisements,
		prometheus.UntypedValue, snmp6.Icmp6.InRouterAdvertisements,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6InNeighborSolicits,
		prometheus.UntypedValue, snmp6.Icmp6.InNeighborSolicits,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6InNeighborAdvertisements,
		prometheus.UntypedValue, snmp6.Icmp6.InNeighborAdvertisements,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6InRedirects,
		prometheus.UntypedValue, snmp6.Icmp6.InRedirects,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6InMLDv2Reports,
		prometheus.UntypedValue, snmp6.Icmp6.InMLDv2Reports,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6OutDestUnreachs,
		prometheus.UntypedValue, snmp6.Icmp6.OutDestUnreachs,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6OutPktTooBigs,
		prometheus.UntypedValue, snmp6.Icmp6.OutPktTooBigs,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6OutTimeExcds,
		prometheus.UntypedValue, snmp6.Icmp6.OutTimeExcds,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6OutParmProblems,
		prometheus.UntypedValue, snmp6.Icmp6.OutParmProblems,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6OutEchos,
		prometheus.UntypedValue, snmp6.Icmp6.OutEchos,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6OutEchoReplies,
		prometheus.UntypedValue, snmp6.Icmp6.OutEchoReplies,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6OutGroupMembQueries,
		prometheus.UntypedValue, snmp6.Icmp6.OutGroupMembQueries,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6OutGroupMembResponses,
		prometheus.UntypedValue, snmp6.Icmp6.OutGroupMembResponses,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6OutGroupMembReductions,
		prometheus.UntypedValue, snmp6.Icmp6.OutGroupMembReductions,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6OutRouterSolicits,
		prometheus.UntypedValue, snmp6.Icmp6.OutRouterSolicits,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6OutRouterAdvertisements,
		prometheus.UntypedValue, snmp6.Icmp6.OutRouterAdvertisements,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6OutNeighborSolicits,
		prometheus.UntypedValue, snmp6.Icmp6.OutNeighborSolicits,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6OutNeighborAdvertisements,
		prometheus.UntypedValue, snmp6.Icmp6.OutNeighborAdvertisements,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6OutRedirects,
		prometheus.UntypedValue, snmp6.Icmp6.OutRedirects,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6OutMLDv2Reports,
		prometheus.UntypedValue, snmp6.Icmp6.OutMLDv2Reports,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6InType1,
		prometheus.UntypedValue, snmp6.Icmp6.InType1,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6InType134,
		prometheus.UntypedValue, snmp6.Icmp6.InType134,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6InType135,
		prometheus.UntypedValue, snmp6.Icmp6.InType135,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6InType136,
		prometheus.UntypedValue, snmp6.Icmp6.InType136,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6InType143,
		prometheus.UntypedValue, snmp6.Icmp6.InType143,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6OutType133,
		prometheus.UntypedValue, snmp6.Icmp6.OutType133,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6OutType135,
		prometheus.UntypedValue, snmp6.Icmp6.OutType135,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6OutType136,
		prometheus.UntypedValue, snmp6.Icmp6.OutType136,
	)
	ch <- prometheus.MustNewConstMetric(c.Icmp6OutType143,
		prometheus.UntypedValue, snmp6.Icmp6.OutType143,
	)

	// Udp6
	ch <- prometheus.MustNewConstMetric(c.Udp6InDatagrams,
		prometheus.UntypedValue, snmp6.Udp6.InDatagrams,
	)
	ch <- prometheus.MustNewConstMetric(c.Udp6NoPorts,
		prometheus.UntypedValue, snmp6.Udp6.NoPorts,
	)
	ch <- prometheus.MustNewConstMetric(c.Udp6InErrors,
		prometheus.UntypedValue, snmp6.Udp6.InErrors,
	)
	ch <- prometheus.MustNewConstMetric(c.Udp6OutDatagrams,
		prometheus.UntypedValue, snmp6.Udp6.OutDatagrams,
	)
	ch <- prometheus.MustNewConstMetric(c.Udp6RcvbufErrors,
		prometheus.UntypedValue, snmp6.Udp6.RcvbufErrors,
	)
	ch <- prometheus.MustNewConstMetric(c.Udp6SndbufErrors,
		prometheus.UntypedValue, snmp6.Udp6.SndbufErrors,
	)
	ch <- prometheus.MustNewConstMetric(c.Udp6InCsumErrors,
		prometheus.UntypedValue, snmp6.Udp6.InCsumErrors,
	)
	ch <- prometheus.MustNewConstMetric(c.Udp6IgnoredMulti,
		prometheus.UntypedValue, snmp6.Udp6.IgnoredMulti,
	)

	// UdpLite6
	ch <- prometheus.MustNewConstMetric(c.UdpLite6InDatagrams,
		prometheus.UntypedValue, snmp6.UdpLite6.InDatagrams,
	)
	ch <- prometheus.MustNewConstMetric(c.UdpLite6NoPorts,
		prometheus.UntypedValue, snmp6.UdpLite6.NoPorts,
	)
	ch <- prometheus.MustNewConstMetric(c.UdpLite6InErrors,
		prometheus.UntypedValue, snmp6.UdpLite6.InErrors,
	)
	ch <- prometheus.MustNewConstMetric(c.UdpLite6OutDatagrams,
		prometheus.UntypedValue, snmp6.UdpLite6.OutDatagrams,
	)
	ch <- prometheus.MustNewConstMetric(c.UdpLite6RcvbufErrors,
		prometheus.UntypedValue, snmp6.UdpLite6.RcvbufErrors,
	)
	ch <- prometheus.MustNewConstMetric(c.UdpLite6SndbufErrors,
		prometheus.UntypedValue, snmp6.UdpLite6.SndbufErrors,
	)
	ch <- prometheus.MustNewConstMetric(c.UdpLite6InCsumErrors,
		prometheus.UntypedValue, snmp6.UdpLite6.InCsumErrors,
	)

	return nil
}
