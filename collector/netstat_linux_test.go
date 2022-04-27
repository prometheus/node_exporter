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
//
package collector

import (
	"fmt"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"os"
	"strings"
	"testing"
)

type testNetStatCollector struct {
	ntc Collector
}

func (c testNetStatCollector) Collect(ch chan<- prometheus.Metric) {
	c.ntc.Update(ch)
}

func (c testNetStatCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func newTestNetStatCollector(logger log.Logger) (prometheus.Collector, error) {
	ntc, err := NewNetStatCollector(logger)
	if err != nil {
		return testNetStatCollector{}, err
	}
	return testNetStatCollector{
		ntc: ntc,
	}, nil
}

func TestNetStatCollector_Update(t *testing.T) {
	*sysPath = "fixtures/sys"
	*procPath = "fixtures/proc"
	*ignoredDevices = "^(ram|loop|fd|(h|s|v|xv)d[a-z]|nvme\\d+n\\d+p)\\d+$"

	testcase := `
	# HELP node_netstat_TcpExt_SyncookiesSent Statistic TcpExtSyncookiesSent.
	# TYPE node_netstat_TcpExt_SyncookiesSent untyped
	node_netstat_TcpExt_SyncookiesSent 0
	# HELP node_netstat_TcpExt_SyncookiesRecv Statistic TcpExtSyncookiesRecv.
	# TYPE node_netstat_TcpExt_SyncookiesRecv untyped
	node_netstat_TcpExt_SyncookiesRecv 0
	# HELP node_netstat_TcpExt_SyncookiesFailed Statistic TcpExtSyncookiesFailed.
	# TYPE node_netstat_TcpExt_SyncookiesFailed untyped
	node_netstat_TcpExt_SyncookiesFailed 2
	# HELP node_netstat_TcpExt_EmbryonicRsts Statistic TcpExtEmbryonicRsts.
	# TYPE node_netstat_TcpExt_EmbryonicRsts untyped
	node_netstat_TcpExt_EmbryonicRsts 0
	# HELP node_netstat_TcpExt_PruneCalled Statistic TcpExtPruneCalled.
	# TYPE node_netstat_TcpExt_PruneCalled untyped
	node_netstat_TcpExt_PruneCalled 0
	# HELP node_netstat_TcpExt_RcvPruned Statistic TcpExtRcvPruned.
	# TYPE node_netstat_TcpExt_RcvPruned untyped
	node_netstat_TcpExt_RcvPruned 0
	# HELP node_netstat_TcpExt_OfoPruned Statistic TcpExtOfoPruned.
	# TYPE node_netstat_TcpExt_OfoPruned untyped
	node_netstat_TcpExt_OfoPruned 0
	# HELP node_netstat_TcpExt_OutOfWindowIcmps Statistic TcpExtOutOfWindowIcmps.
	# TYPE node_netstat_TcpExt_OutOfWindowIcmps untyped
	node_netstat_TcpExt_OutOfWindowIcmps 0
	# HELP node_netstat_TcpExt_LockDroppedIcmps Statistic TcpExtLockDroppedIcmps.
	# TYPE node_netstat_TcpExt_LockDroppedIcmps untyped
	node_netstat_TcpExt_LockDroppedIcmps 0
	# HELP node_netstat_TcpExt_ArpFilter Statistic TcpExtArpFilter.
	# TYPE node_netstat_TcpExt_ArpFilter untyped
	node_netstat_TcpExt_ArpFilter 0
	# HELP node_netstat_TcpExt_TW Statistic TcpExtTW.
	# TYPE node_netstat_TcpExt_TW untyped
	node_netstat_TcpExt_TW 388812
	# HELP node_netstat_TcpExt_TWRecycled Statistic TcpExtTWRecycled.
	# TYPE node_netstat_TcpExt_TWRecycled untyped
	node_netstat_TcpExt_TWRecycled 0
	# HELP node_netstat_TcpExt_TWKilled Statistic TcpExtTWKilled.
	# TYPE node_netstat_TcpExt_TWKilled untyped
	node_netstat_TcpExt_TWKilled 0
	# HELP node_netstat_TcpExt_PAWSPassive Statistic TcpExtPAWSPassive.
	# TYPE node_netstat_TcpExt_PAWSPassive untyped
	node_netstat_TcpExt_PAWSPassive 0
	# HELP node_netstat_TcpExt_PAWSActive Statistic TcpExtPAWSActive.
	# TYPE node_netstat_TcpExt_PAWSActive untyped
	node_netstat_TcpExt_PAWSActive 0
	# HELP node_netstat_TcpExt_PAWSEstab Statistic TcpExtPAWSEstab.
	# TYPE node_netstat_TcpExt_PAWSEstab untyped
	node_netstat_TcpExt_PAWSEstab 6
	# HELP node_netstat_TcpExt_DelayedACKs Statistic TcpExtDelayedACKs.
	# TYPE node_netstat_TcpExt_DelayedACKs untyped
	node_netstat_TcpExt_DelayedACKs 102471
	# HELP node_netstat_TcpExt_DelayedACKLocked Statistic TcpExtDelayedACKLocked.
	# TYPE node_netstat_TcpExt_DelayedACKLocked untyped
	node_netstat_TcpExt_DelayedACKLocked 17
	# HELP node_netstat_TcpExt_DelayedACKLost Statistic TcpExtDelayedACKLost.
	# TYPE node_netstat_TcpExt_DelayedACKLost untyped
	node_netstat_TcpExt_DelayedACKLost 9
	# HELP node_netstat_TcpExt_ListenOverflows Statistic TcpExtListenOverflows.
	# TYPE node_netstat_TcpExt_ListenOverflows untyped
	node_netstat_TcpExt_ListenOverflows 0
	# HELP node_netstat_TcpExt_ListenDrops Statistic TcpExtListenDrops.
	# TYPE node_netstat_TcpExt_ListenDrops untyped
	node_netstat_TcpExt_ListenDrops 0
	# HELP node_netstat_TcpExt_TCPPrequeued Statistic TcpExtTCPPrequeued.
	# TYPE node_netstat_TcpExt_TCPPrequeued untyped
	node_netstat_TcpExt_TCPPrequeued 80568
	# HELP node_netstat_TcpExt_TCPDirectCopyFromBacklog Statistic TcpExtTCPDirectCopyFromBacklog.
	# TYPE node_netstat_TcpExt_TCPDirectCopyFromBacklog untyped
	node_netstat_TcpExt_TCPDirectCopyFromBacklog 0
	# HELP node_netstat_TcpExt_TCPDirectCopyFromPrequeue Statistic TcpExtTCPDirectCopyFromPrequeue.
	# TYPE node_netstat_TcpExt_TCPDirectCopyFromPrequeue untyped
	node_netstat_TcpExt_TCPDirectCopyFromPrequeue 168808
	# HELP node_netstat_TcpExt_TCPPrequeueDropped Statistic TcpExtTCPPrequeueDropped.
	# TYPE node_netstat_TcpExt_TCPPrequeueDropped untyped
	node_netstat_TcpExt_TCPPrequeueDropped 0
	# HELP node_netstat_TcpExt_TCPHPHits Statistic TcpExtTCPHPHits.
	# TYPE node_netstat_TcpExt_TCPHPHits untyped
	node_netstat_TcpExt_TCPHPHits 4471289
	# HELP node_netstat_TcpExt_TCPHPHitsToUser Statistic TcpExtTCPHPHitsToUser.
	# TYPE node_netstat_TcpExt_TCPHPHitsToUser untyped
	node_netstat_TcpExt_TCPHPHitsToUser 26
	# HELP node_netstat_TcpExt_TCPPureAcks Statistic TcpExtTCPPureAcks.
	# TYPE node_netstat_TcpExt_TCPPureAcks untyped
	node_netstat_TcpExt_TCPPureAcks 1433940
	# HELP node_netstat_TcpExt_TCPHPAcks Statistic TcpExtTCPHPAcks.
	# TYPE node_netstat_TcpExt_TCPHPAcks untyped
	node_netstat_TcpExt_TCPHPAcks 3744565
	# HELP node_netstat_TcpExt_TCPRenoRecovery Statistic TcpExtTCPRenoRecovery.
	# TYPE node_netstat_TcpExt_TCPRenoRecovery untyped
	node_netstat_TcpExt_TCPRenoRecovery 0
	# HELP node_netstat_TcpExt_TCPSackRecovery Statistic TcpExtTCPSackRecovery.
	# TYPE node_netstat_TcpExt_TCPSackRecovery untyped
	node_netstat_TcpExt_TCPSackRecovery 1
	# HELP node_netstat_TcpExt_TCPSACKReneging Statistic TcpExtTCPSACKReneging.
	# TYPE node_netstat_TcpExt_TCPSACKReneging untyped
	node_netstat_TcpExt_TCPSACKReneging 0
	# HELP node_netstat_TcpExt_TCPFACKReorder Statistic TcpExtTCPFACKReorder.
	# TYPE node_netstat_TcpExt_TCPFACKReorder untyped
	node_netstat_TcpExt_TCPFACKReorder 0
	# HELP node_netstat_TcpExt_TCPSACKReorder Statistic TcpExtTCPSACKReorder.
	# TYPE node_netstat_TcpExt_TCPSACKReorder untyped
	node_netstat_TcpExt_TCPSACKReorder 0
	# HELP node_netstat_TcpExt_TCPRenoReorder Statistic TcpExtTCPRenoReorder.
	# TYPE node_netstat_TcpExt_TCPRenoReorder untyped
	node_netstat_TcpExt_TCPRenoReorder 0
	# HELP node_netstat_TcpExt_TCPTSReorder Statistic TcpExtTCPTSReorder.
	# TYPE node_netstat_TcpExt_TCPTSReorder untyped
	node_netstat_TcpExt_TCPTSReorder 0
	# HELP node_netstat_TcpExt_TCPFullUndo Statistic TcpExtTCPFullUndo.
	# TYPE node_netstat_TcpExt_TCPFullUndo untyped
	node_netstat_TcpExt_TCPFullUndo 0
	# HELP node_netstat_TcpExt_TCPPartialUndo Statistic TcpExtTCPPartialUndo.
	# TYPE node_netstat_TcpExt_TCPPartialUndo untyped
	node_netstat_TcpExt_TCPPartialUndo 0
	# HELP node_netstat_TcpExt_TCPDSACKUndo Statistic TcpExtTCPDSACKUndo.
	# TYPE node_netstat_TcpExt_TCPDSACKUndo untyped
	node_netstat_TcpExt_TCPDSACKUndo 0
	# HELP node_netstat_TcpExt_TCPLossUndo Statistic TcpExtTCPLossUndo.
	# TYPE node_netstat_TcpExt_TCPLossUndo untyped
	node_netstat_TcpExt_TCPLossUndo 48
	# HELP node_netstat_TcpExt_TCPLoss Statistic TcpExtTCPLoss.
	# TYPE node_netstat_TcpExt_TCPLoss untyped
	node_netstat_TcpExt_TCPLoss 0
	# HELP node_netstat_TcpExt_TCPLostRetransmit Statistic TcpExtTCPLostRetransmit.
	# TYPE node_netstat_TcpExt_TCPLostRetransmit untyped
	node_netstat_TcpExt_TCPLostRetransmit 0
	# HELP node_netstat_TcpExt_TCPRenoFailures Statistic TcpExtTCPRenoFailures.
	# TYPE node_netstat_TcpExt_TCPRenoFailures untyped
	node_netstat_TcpExt_TCPRenoFailures 0
	# HELP node_netstat_TcpExt_TCPSackFailures Statistic TcpExtTCPSackFailures.
	# TYPE node_netstat_TcpExt_TCPSackFailures untyped
	node_netstat_TcpExt_TCPSackFailures 1
	# HELP node_netstat_TcpExt_TCPLossFailures Statistic TcpExtTCPLossFailures.
	# TYPE node_netstat_TcpExt_TCPLossFailures untyped
	node_netstat_TcpExt_TCPLossFailures 0
	# HELP node_netstat_TcpExt_TCPFastRetrans Statistic TcpExtTCPFastRetrans.
	# TYPE node_netstat_TcpExt_TCPFastRetrans untyped
	node_netstat_TcpExt_TCPFastRetrans 1
	# HELP node_netstat_TcpExt_TCPForwardRetrans Statistic TcpExtTCPForwardRetrans.
	# TYPE node_netstat_TcpExt_TCPForwardRetrans untyped
	node_netstat_TcpExt_TCPForwardRetrans 0
	# HELP node_netstat_TcpExt_TCPSlowStartRetrans Statistic TcpExtTCPSlowStartRetrans.
	# TYPE node_netstat_TcpExt_TCPSlowStartRetrans untyped
	node_netstat_TcpExt_TCPSlowStartRetrans 1
	# HELP node_netstat_TcpExt_TCPTimeouts Statistic TcpExtTCPTimeouts.
	# TYPE node_netstat_TcpExt_TCPTimeouts untyped
	node_netstat_TcpExt_TCPTimeouts 115
	# HELP node_netstat_TcpExt_TCPRenoRecoveryFail Statistic TcpExtTCPRenoRecoveryFail.
	# TYPE node_netstat_TcpExt_TCPRenoRecoveryFail untyped
	node_netstat_TcpExt_TCPRenoRecoveryFail 0
	# HELP node_netstat_TcpExt_TCPSackRecoveryFail Statistic TcpExtTCPSackRecoveryFail.
	# TYPE node_netstat_TcpExt_TCPSackRecoveryFail untyped
	node_netstat_TcpExt_TCPSackRecoveryFail 0
	# HELP node_netstat_TcpExt_TCPSchedulerFailed Statistic TcpExtTCPSchedulerFailed.
	# TYPE node_netstat_TcpExt_TCPSchedulerFailed untyped
	node_netstat_TcpExt_TCPSchedulerFailed 0
	# HELP node_netstat_TcpExt_TCPRcvCollapsed Statistic TcpExtTCPRcvCollapsed.
	# TYPE node_netstat_TcpExt_TCPRcvCollapsed untyped
	node_netstat_TcpExt_TCPRcvCollapsed 0
	# HELP node_netstat_TcpExt_TCPDSACKOldSent Statistic TcpExtTCPDSACKOldSent.
	# TYPE node_netstat_TcpExt_TCPDSACKOldSent untyped
	node_netstat_TcpExt_TCPDSACKOldSent 9
	# HELP node_netstat_TcpExt_TCPDSACKOfoSent Statistic TcpExtTCPDSACKOfoSent.
	# TYPE node_netstat_TcpExt_TCPDSACKOfoSent untyped
	node_netstat_TcpExt_TCPDSACKOfoSent 0
	# HELP node_netstat_TcpExt_TCPDSACKRecv Statistic TcpExtTCPDSACKRecv.
	# TYPE node_netstat_TcpExt_TCPDSACKRecv untyped
	node_netstat_TcpExt_TCPDSACKRecv 5
	# HELP node_netstat_TcpExt_TCPDSACKOfoRecv Statistic TcpExtTCPDSACKOfoRecv.
	# TYPE node_netstat_TcpExt_TCPDSACKOfoRecv untyped
	node_netstat_TcpExt_TCPDSACKOfoRecv 0
	# HELP node_netstat_TcpExt_TCPAbortOnData Statistic TcpExtTCPAbortOnData.
	# TYPE node_netstat_TcpExt_TCPAbortOnData untyped
	node_netstat_TcpExt_TCPAbortOnData 41
	# HELP node_netstat_TcpExt_TCPAbortOnClose Statistic TcpExtTCPAbortOnClose.
	# TYPE node_netstat_TcpExt_TCPAbortOnClose untyped
	node_netstat_TcpExt_TCPAbortOnClose 4
	# HELP node_netstat_TcpExt_TCPAbortOnMemory Statistic TcpExtTCPAbortOnMemory.
	# TYPE node_netstat_TcpExt_TCPAbortOnMemory untyped
	node_netstat_TcpExt_TCPAbortOnMemory 0
	# HELP node_netstat_TcpExt_TCPAbortOnTimeout Statistic TcpExtTCPAbortOnTimeout.
	# TYPE node_netstat_TcpExt_TCPAbortOnTimeout untyped
	node_netstat_TcpExt_TCPAbortOnTimeout 0
	# HELP node_netstat_TcpExt_TCPAbortOnLinger Statistic TcpExtTCPAbortOnLinger.
	# TYPE node_netstat_TcpExt_TCPAbortOnLinger untyped
	node_netstat_TcpExt_TCPAbortOnLinger 0
	# HELP node_netstat_TcpExt_TCPAbortFailed Statistic TcpExtTCPAbortFailed.
	# TYPE node_netstat_TcpExt_TCPAbortFailed untyped
	node_netstat_TcpExt_TCPAbortFailed 0
	# HELP node_netstat_TcpExt_TCPMemoryPressures Statistic TcpExtTCPMemoryPressures.
	# TYPE node_netstat_TcpExt_TCPMemoryPressures untyped
	node_netstat_TcpExt_TCPMemoryPressures 0
	# HELP node_netstat_TcpExt_TCPSACKDiscard Statistic TcpExtTCPSACKDiscard.
	# TYPE node_netstat_TcpExt_TCPSACKDiscard untyped
	node_netstat_TcpExt_TCPSACKDiscard 0
	# HELP node_netstat_TcpExt_TCPDSACKIgnoredOld Statistic TcpExtTCPDSACKIgnoredOld.
	# TYPE node_netstat_TcpExt_TCPDSACKIgnoredOld untyped
	node_netstat_TcpExt_TCPDSACKIgnoredOld 0
	# HELP node_netstat_TcpExt_TCPDSACKIgnoredNoUndo Statistic TcpExtTCPDSACKIgnoredNoUndo.
	# TYPE node_netstat_TcpExt_TCPDSACKIgnoredNoUndo untyped
	node_netstat_TcpExt_TCPDSACKIgnoredNoUndo 1
	# HELP node_netstat_TcpExt_TCPSpuriousRTOs Statistic TcpExtTCPSpuriousRTOs.
	# TYPE node_netstat_TcpExt_TCPSpuriousRTOs untyped
	node_netstat_TcpExt_TCPSpuriousRTOs 0
	# HELP node_netstat_TcpExt_TCPMD5NotFound Statistic TcpExtTCPMD5NotFound.
	# TYPE node_netstat_TcpExt_TCPMD5NotFound untyped
	node_netstat_TcpExt_TCPMD5NotFound 0
	# HELP node_netstat_TcpExt_TCPMD5Unexpected Statistic TcpExtTCPMD5Unexpected.
	# TYPE node_netstat_TcpExt_TCPMD5Unexpected untyped
	node_netstat_TcpExt_TCPMD5Unexpected 0
	# HELP node_netstat_TcpExt_TCPSackShifted Statistic TcpExtTCPSackShifted.
	# TYPE node_netstat_TcpExt_TCPSackShifted untyped
	node_netstat_TcpExt_TCPSackShifted 0
	# HELP node_netstat_TcpExt_TCPSackMerged Statistic TcpExtTCPSackMerged.
	# TYPE node_netstat_TcpExt_TCPSackMerged untyped
	node_netstat_TcpExt_TCPSackMerged 2
	# HELP node_netstat_TcpExt_TCPSackShiftFallback Statistic TcpExtTCPSackShiftFallback.
	# TYPE node_netstat_TcpExt_TCPSackShiftFallback untyped
	node_netstat_TcpExt_TCPSackShiftFallback 5
	# HELP node_netstat_TcpExt_TCPBacklogDrop Statistic TcpExtTCPBacklogDrop.
	# TYPE node_netstat_TcpExt_TCPBacklogDrop untyped
	node_netstat_TcpExt_TCPBacklogDrop 0
	# HELP node_netstat_TcpExt_TCPMinTTLDrop Statistic TcpExtTCPMinTTLDrop.
	# TYPE node_netstat_TcpExt_TCPMinTTLDrop untyped
	node_netstat_TcpExt_TCPMinTTLDrop 0
	# HELP node_netstat_TcpExt_TCPDeferAcceptDrop Statistic TcpExtTCPDeferAcceptDrop.
	# TYPE node_netstat_TcpExt_TCPDeferAcceptDrop untyped
	node_netstat_TcpExt_TCPDeferAcceptDrop 0
	# HELP node_netstat_TcpExt_IPReversePathFilter Statistic TcpExtIPReversePathFilter.
	# TYPE node_netstat_TcpExt_IPReversePathFilter untyped
	node_netstat_TcpExt_IPReversePathFilter 0
	# HELP node_netstat_TcpExt_TCPTimeWaitOverflow Statistic TcpExtTCPTimeWaitOverflow.
	# TYPE node_netstat_TcpExt_TCPTimeWaitOverflow untyped
	node_netstat_TcpExt_TCPTimeWaitOverflow 0
	# HELP node_netstat_TcpExt_TCPReqQFullDoCookies Statistic TcpExtTCPReqQFullDoCookies.
	# TYPE node_netstat_TcpExt_TCPReqQFullDoCookies untyped
	node_netstat_TcpExt_TCPReqQFullDoCookies 0
	# HELP node_netstat_TcpExt_TCPReqQFullDrop Statistic TcpExtTCPReqQFullDrop.
	# TYPE node_netstat_TcpExt_TCPReqQFullDrop untyped
	node_netstat_TcpExt_TCPReqQFullDrop 0
	# HELP node_netstat_TcpExt_TCPChallengeACK Statistic TcpExtTCPChallengeACK.
	# TYPE node_netstat_TcpExt_TCPChallengeACK untyped
	node_netstat_TcpExt_TCPChallengeACK 2
	# HELP node_netstat_TcpExt_TCPSYNChallenge Statistic TcpExtTCPSYNChallenge.
	# TYPE node_netstat_TcpExt_TCPSYNChallenge untyped
	node_netstat_TcpExt_TCPSYNChallenge 2
	# HELP node_netstat_IpExt_InNoRoutes Statistic IpExtInNoRoutes.
	# TYPE node_netstat_IpExt_InNoRoutes untyped
	node_netstat_IpExt_InNoRoutes 0
	# HELP node_netstat_IpExt_InTruncatedPkts Statistic IpExtInTruncatedPkts.
	# TYPE node_netstat_IpExt_InTruncatedPkts untyped
	node_netstat_IpExt_InTruncatedPkts 0
	# HELP node_netstat_IpExt_InMcastPkts Statistic IpExtInMcastPkts.
	# TYPE node_netstat_IpExt_InMcastPkts untyped
	node_netstat_IpExt_InMcastPkts 0
	# HELP node_netstat_IpExt_OutMcastPkts Statistic IpExtOutMcastPkts.
	# TYPE node_netstat_IpExt_OutMcastPkts untyped
	node_netstat_IpExt_OutMcastPkts 0
	# HELP node_netstat_IpExt_InBcastPkts Statistic IpExtInBcastPkts.
	# TYPE node_netstat_IpExt_InBcastPkts untyped
	node_netstat_IpExt_InBcastPkts 0
	# HELP node_netstat_IpExt_OutBcastPkts Statistic IpExtOutBcastPkts.
	# TYPE node_netstat_IpExt_OutBcastPkts untyped
	node_netstat_IpExt_OutBcastPkts 0
	# HELP node_netstat_IpExt_InOctets Statistic IpExtInOctets.
	# TYPE node_netstat_IpExt_InOctets untyped
	node_netstat_IpExt_InOctets 6286396970
	# HELP node_netstat_IpExt_OutOctets Statistic IpExtOutOctets.
	# TYPE node_netstat_IpExt_OutOctets untyped
	node_netstat_IpExt_OutOctets 2786264347
	# HELP node_netstat_IpExt_InMcastOctets Statistic IpExtInMcastOctets.
	# TYPE node_netstat_IpExt_InMcastOctets untyped
	node_netstat_IpExt_InMcastOctets 0
	# HELP node_netstat_IpExt_OutMcastOctets Statistic IpExtOutMcastOctets.
	# TYPE node_netstat_IpExt_OutMcastOctets untyped
	node_netstat_IpExt_OutMcastOctets 0
	# HELP node_netstat_IpExt_InBcastOctets Statistic IpExtInBcastOctets.
	# TYPE node_netstat_IpExt_InBcastOctets untyped
	node_netstat_IpExt_InBcastOctets 0
	# HELP node_netstat_IpExt_OutBcastOctets Statistic IpExtOutBcastOctets.
	# TYPE node_netstat_IpExt_OutBcastOctets untyped
	node_netstat_IpExt_OutBcastOctets 0
	# HELP node_netstat_Ip_Forwarding Statistic IpForwarding.
	# TYPE node_netstat_Ip_Forwarding untyped
	node_netstat_Ip_Forwarding 1
	# HELP node_netstat_Ip_DefaultTTL Statistic IpDefaultTTL.
	# TYPE node_netstat_Ip_DefaultTTL untyped
	node_netstat_Ip_DefaultTTL 64
	# HELP node_netstat_Ip_InReceives Statistic IpInReceives.
	# TYPE node_netstat_Ip_InReceives untyped
	node_netstat_Ip_InReceives 57740232
	# HELP node_netstat_Ip_InHdrErrors Statistic IpInHdrErrors.
	# TYPE node_netstat_Ip_InHdrErrors untyped
	node_netstat_Ip_InHdrErrors 0
	# HELP node_netstat_Ip_InAddrErrors Statistic IpInAddrErrors.
	# TYPE node_netstat_Ip_InAddrErrors untyped
	node_netstat_Ip_InAddrErrors 25
	# HELP node_netstat_Ip_ForwDatagrams Statistic IpForwDatagrams.
	# TYPE node_netstat_Ip_ForwDatagrams untyped
	node_netstat_Ip_ForwDatagrams 397750
	# HELP node_netstat_Ip_InUnknownProtos Statistic IpInUnknownProtos.
	# TYPE node_netstat_Ip_InUnknownProtos untyped
	node_netstat_Ip_InUnknownProtos 0
	# HELP node_netstat_Ip_InDiscards Statistic IpInDiscards.
	# TYPE node_netstat_Ip_InDiscards untyped
	node_netstat_Ip_InDiscards 0
	# HELP node_netstat_Ip_InDelivers Statistic IpInDelivers.
	# TYPE node_netstat_Ip_InDelivers untyped
	node_netstat_Ip_InDelivers 57340175
	# HELP node_netstat_Ip_OutRequests Statistic IpOutRequests.
	# TYPE node_netstat_Ip_OutRequests untyped
	node_netstat_Ip_OutRequests 55365537
	# HELP node_netstat_Ip_OutDiscards Statistic IpOutDiscards.
	# TYPE node_netstat_Ip_OutDiscards untyped
	node_netstat_Ip_OutDiscards 0
	# HELP node_netstat_Ip_OutNoRoutes Statistic IpOutNoRoutes.
	# TYPE node_netstat_Ip_OutNoRoutes untyped
	node_netstat_Ip_OutNoRoutes 54
	# HELP node_netstat_Ip_ReasmTimeout Statistic IpReasmTimeout.
	# TYPE node_netstat_Ip_ReasmTimeout untyped
	node_netstat_Ip_ReasmTimeout 0
	# HELP node_netstat_Ip_ReasmReqds Statistic IpReasmReqds.
	# TYPE node_netstat_Ip_ReasmReqds untyped
	node_netstat_Ip_ReasmReqds 0
	# HELP node_netstat_Ip_ReasmOKs Statistic IpReasmOKs.
	# TYPE node_netstat_Ip_ReasmOKs untyped
	node_netstat_Ip_ReasmOKs 0
	# HELP node_netstat_Ip_ReasmFails Statistic IpReasmFails.
	# TYPE node_netstat_Ip_ReasmFails untyped
	node_netstat_Ip_ReasmFails 0
	# HELP node_netstat_Ip_FragOKs Statistic IpFragOKs.
	# TYPE node_netstat_Ip_FragOKs untyped
	node_netstat_Ip_FragOKs 0
	# HELP node_netstat_Ip_FragFails Statistic IpFragFails.
	# TYPE node_netstat_Ip_FragFails untyped
	node_netstat_Ip_FragFails 0
	# HELP node_netstat_Ip_FragCreates Statistic IpFragCreates.
	# TYPE node_netstat_Ip_FragCreates untyped
	node_netstat_Ip_FragCreates 0
	# HELP node_netstat_Icmp_InMsgs Statistic IcmpInMsgs.
	# TYPE node_netstat_Icmp_InMsgs untyped
	node_netstat_Icmp_InMsgs 104
	# HELP node_netstat_Icmp_InErrors Statistic IcmpInErrors.
	# TYPE node_netstat_Icmp_InErrors untyped
	node_netstat_Icmp_InErrors 0
	# HELP node_netstat_Icmp_InCsumErrors Statistic IcmpInCsumErrors.
	# TYPE node_netstat_Icmp_InCsumErrors untyped
	node_netstat_Icmp_InCsumErrors 0
	# HELP node_netstat_Icmp_InDestUnreachs Statistic IcmpInDestUnreachs.
	# TYPE node_netstat_Icmp_InDestUnreachs untyped
	node_netstat_Icmp_InDestUnreachs 104
	# HELP node_netstat_Icmp_InTimeExcds Statistic IcmpInTimeExcds.
	# TYPE node_netstat_Icmp_InTimeExcds untyped
	node_netstat_Icmp_InTimeExcds 0
	# HELP node_netstat_Icmp_InParmProbs Statistic IcmpInParmProbs.
	# TYPE node_netstat_Icmp_InParmProbs untyped
	node_netstat_Icmp_InParmProbs 0
	# HELP node_netstat_Icmp_InSrcQuenchs Statistic IcmpInSrcQuenchs.
	# TYPE node_netstat_Icmp_InSrcQuenchs untyped
	node_netstat_Icmp_InSrcQuenchs 0
	# HELP node_netstat_Icmp_InRedirects Statistic IcmpInRedirects.
	# TYPE node_netstat_Icmp_InRedirects untyped
	node_netstat_Icmp_InRedirects 0
	# HELP node_netstat_Icmp_InEchos Statistic IcmpInEchos.
	# TYPE node_netstat_Icmp_InEchos untyped
	node_netstat_Icmp_InEchos 0
	# HELP node_netstat_Icmp_InEchoReps Statistic IcmpInEchoReps.
	# TYPE node_netstat_Icmp_InEchoReps untyped
	node_netstat_Icmp_InEchoReps 0
	# HELP node_netstat_Icmp_InTimestamps Statistic IcmpInTimestamps.
	# TYPE node_netstat_Icmp_InTimestamps untyped
	node_netstat_Icmp_InTimestamps 0
	# HELP node_netstat_Icmp_InTimestampReps Statistic IcmpInTimestampReps.
	# TYPE node_netstat_Icmp_InTimestampReps untyped
	node_netstat_Icmp_InTimestampReps 0
	# HELP node_netstat_Icmp_InAddrMasks Statistic IcmpInAddrMasks.
	# TYPE node_netstat_Icmp_InAddrMasks untyped
	node_netstat_Icmp_InAddrMasks 0
	# HELP node_netstat_Icmp_InAddrMaskReps Statistic IcmpInAddrMaskReps.
	# TYPE node_netstat_Icmp_InAddrMaskReps untyped
	node_netstat_Icmp_InAddrMaskReps 0
	# HELP node_netstat_Icmp_OutMsgs Statistic IcmpOutMsgs.
	# TYPE node_netstat_Icmp_OutMsgs untyped
	node_netstat_Icmp_OutMsgs 120
	# HELP node_netstat_Icmp_OutErrors Statistic IcmpOutErrors.
	# TYPE node_netstat_Icmp_OutErrors untyped
	node_netstat_Icmp_OutErrors 0
	# HELP node_netstat_Icmp_OutDestUnreachs Statistic IcmpOutDestUnreachs.
	# TYPE node_netstat_Icmp_OutDestUnreachs untyped
	node_netstat_Icmp_OutDestUnreachs 120
	# HELP node_netstat_Icmp_OutTimeExcds Statistic IcmpOutTimeExcds.
	# TYPE node_netstat_Icmp_OutTimeExcds untyped
	node_netstat_Icmp_OutTimeExcds 0
	# HELP node_netstat_Icmp_OutParmProbs Statistic IcmpOutParmProbs.
	# TYPE node_netstat_Icmp_OutParmProbs untyped
	node_netstat_Icmp_OutParmProbs 0
	# HELP node_netstat_Icmp_OutSrcQuenchs Statistic IcmpOutSrcQuenchs.
	# TYPE node_netstat_Icmp_OutSrcQuenchs untyped
	node_netstat_Icmp_OutSrcQuenchs 0
	# HELP node_netstat_Icmp_OutRedirects Statistic IcmpOutRedirects.
	# TYPE node_netstat_Icmp_OutRedirects untyped
	node_netstat_Icmp_OutRedirects 0
	# HELP node_netstat_Icmp_OutEchos Statistic IcmpOutEchos.
	# TYPE node_netstat_Icmp_OutEchos untyped
	node_netstat_Icmp_OutEchos 0
	# HELP node_netstat_Icmp_OutEchoReps Statistic IcmpOutEchoReps.
	# TYPE node_netstat_Icmp_OutEchoReps untyped
	node_netstat_Icmp_OutEchoReps 0
	# HELP node_netstat_Icmp_OutTimestamps Statistic IcmpOutTimestamps.
	# TYPE node_netstat_Icmp_OutTimestamps untyped
	node_netstat_Icmp_OutTimestamps 0
	# HELP node_netstat_Icmp_OutTimestampReps Statistic IcmpOutTimestampReps.
	# TYPE node_netstat_Icmp_OutTimestampReps untyped
	node_netstat_Icmp_OutTimestampReps 0
	# HELP node_netstat_Icmp_OutAddrMasks Statistic IcmpOutAddrMasks.
	# TYPE node_netstat_Icmp_OutAddrMasks untyped
	node_netstat_Icmp_OutAddrMasks 0
	# HELP node_netstat_Icmp_OutAddrMaskReps Statistic IcmpOutAddrMaskReps.
	# TYPE node_netstat_Icmp_OutAddrMaskReps untyped
	node_netstat_Icmp_OutAddrMaskReps 0
	# HELP node_netstat_IcmpMsg_InType3 Statistic IcmpMsgInType3.
	# TYPE node_netstat_IcmpMsg_InType3 untyped
	node_netstat_IcmpMsg_InType3 104
	# HELP node_netstat_IcmpMsg_OutType3 Statistic IcmpMsgOutType3.
	# TYPE node_netstat_IcmpMsg_OutType3 untyped
	node_netstat_IcmpMsg_OutType3 120
	# HELP node_netstat_Tcp_RtoAlgorithm Statistic TcpRtoAlgorithm.
	# TYPE node_netstat_Tcp_RtoAlgorithm untyped
	node_netstat_Tcp_RtoAlgorithm 1
	# HELP node_netstat_Tcp_RtoMin Statistic TcpRtoMin.
	# TYPE node_netstat_Tcp_RtoMin untyped
	node_netstat_Tcp_RtoMin 200
	# HELP node_netstat_Tcp_RtoMax Statistic TcpRtoMax.
	# TYPE node_netstat_Tcp_RtoMax untyped
	node_netstat_Tcp_RtoMax 120000
	# HELP node_netstat_Tcp_MaxConn Statistic TcpMaxConn.
	# TYPE node_netstat_Tcp_MaxConn untyped
	node_netstat_Tcp_MaxConn -1
	# HELP node_netstat_Tcp_ActiveOpens Statistic TcpActiveOpens.
	# TYPE node_netstat_Tcp_ActiveOpens untyped
	node_netstat_Tcp_ActiveOpens 3556
	# HELP node_netstat_Tcp_PassiveOpens Statistic TcpPassiveOpens.
	# TYPE node_netstat_Tcp_PassiveOpens untyped
	node_netstat_Tcp_PassiveOpens 230
	# HELP node_netstat_Tcp_AttemptFails Statistic TcpAttemptFails.
	# TYPE node_netstat_Tcp_AttemptFails untyped
	node_netstat_Tcp_AttemptFails 341
	# HELP node_netstat_Tcp_EstabResets Statistic TcpEstabResets.
	# TYPE node_netstat_Tcp_EstabResets untyped
	node_netstat_Tcp_EstabResets 161
	# HELP node_netstat_Tcp_CurrEstab Statistic TcpCurrEstab.
	# TYPE node_netstat_Tcp_CurrEstab untyped
	node_netstat_Tcp_CurrEstab 0
	# HELP node_netstat_Tcp_InSegs Statistic TcpInSegs.
	# TYPE node_netstat_Tcp_InSegs untyped
	node_netstat_Tcp_InSegs 57252008
	# HELP node_netstat_Tcp_OutSegs Statistic TcpOutSegs.
	# TYPE node_netstat_Tcp_OutSegs untyped
	node_netstat_Tcp_OutSegs 54915039
	# HELP node_netstat_Tcp_RetransSegs Statistic TcpRetransSegs.
	# TYPE node_netstat_Tcp_RetransSegs untyped
	node_netstat_Tcp_RetransSegs 227
	# HELP node_netstat_Tcp_InErrs Statistic TcpInErrs.
	# TYPE node_netstat_Tcp_InErrs untyped
	node_netstat_Tcp_InErrs 5
	# HELP node_netstat_Tcp_OutRsts Statistic TcpOutRsts.
	# TYPE node_netstat_Tcp_OutRsts untyped
	node_netstat_Tcp_OutRsts 1003
	# HELP node_netstat_Tcp_InCsumErrors Statistic TcpInCsumErrors.
	# TYPE node_netstat_Tcp_InCsumErrors untyped
	node_netstat_Tcp_InCsumErrors 0
	# HELP node_netstat_Udp_InDatagrams Statistic UdpInDatagrams.
	# TYPE node_netstat_Udp_InDatagrams untyped
	node_netstat_Udp_InDatagrams 88542
	# HELP node_netstat_Udp_NoPorts Statistic UdpNoPorts.
	# TYPE node_netstat_Udp_NoPorts untyped
	node_netstat_Udp_NoPorts 120
	# HELP node_netstat_Udp_InErrors Statistic UdpInErrors.
	# TYPE node_netstat_Udp_InErrors untyped
	node_netstat_Udp_InErrors 0
	# HELP node_netstat_Udp_OutDatagrams Statistic UdpOutDatagrams.
	# TYPE node_netstat_Udp_OutDatagrams untyped
	node_netstat_Udp_OutDatagrams 53028
	# HELP node_netstat_Udp_RcvbufErrors Statistic UdpRcvbufErrors.
	# TYPE node_netstat_Udp_RcvbufErrors untyped
	node_netstat_Udp_RcvbufErrors 9
	# HELP node_netstat_Udp_SndbufErrors Statistic UdpSndbufErrors.
	# TYPE node_netstat_Udp_SndbufErrors untyped
	node_netstat_Udp_SndbufErrors 8
	# HELP node_netstat_Udp_InCsumErrors Statistic UdpInCsumErrors.
	# TYPE node_netstat_Udp_InCsumErrors untyped
	node_netstat_Udp_InCsumErrors 0
	# HELP node_netstat_UdpLite_InDatagrams Statistic UdpLiteInDatagrams.
	# TYPE node_netstat_UdpLite_InDatagrams untyped
	node_netstat_UdpLite_InDatagrams 0
	# HELP node_netstat_UdpLite_NoPorts Statistic UdpLiteNoPorts.
	# TYPE node_netstat_UdpLite_NoPorts untyped
	node_netstat_UdpLite_NoPorts 0
	# HELP node_netstat_UdpLite_InErrors Statistic UdpLiteInErrors.
	# TYPE node_netstat_UdpLite_InErrors untyped
	node_netstat_UdpLite_InErrors 0
	# HELP node_netstat_UdpLite_OutDatagrams Statistic UdpLiteOutDatagrams.
	# TYPE node_netstat_UdpLite_OutDatagrams untyped
	node_netstat_UdpLite_OutDatagrams 0
	# HELP node_netstat_UdpLite_RcvbufErrors Statistic UdpLiteRcvbufErrors.
	# TYPE node_netstat_UdpLite_RcvbufErrors untyped
	node_netstat_UdpLite_RcvbufErrors 0
	# HELP node_netstat_UdpLite_SndbufErrors Statistic UdpLiteSndbufErrors.
	# TYPE node_netstat_UdpLite_SndbufErrors untyped
	node_netstat_UdpLite_SndbufErrors 0
	# HELP node_netstat_UdpLite_InCsumErrors Statistic UdpLiteInCsumErrors.
	# TYPE node_netstat_UdpLite_InCsumErrors untyped
	node_netstat_UdpLite_InCsumErrors 0
	# HELP node_netstat_Ip6_InReceives Statistic Ip6InReceives.
	# TYPE node_netstat_Ip6_InReceives untyped
	node_netstat_Ip6_InReceives 7
	# HELP node_netstat_Ip6_InHdrErrors Statistic Ip6InHdrErrors.
	# TYPE node_netstat_Ip6_InHdrErrors untyped
	node_netstat_Ip6_InHdrErrors 0
	# HELP node_netstat_Ip6_InTooBigErrors Statistic Ip6InTooBigErrors.
	# TYPE node_netstat_Ip6_InTooBigErrors untyped
	node_netstat_Ip6_InTooBigErrors 0
	# HELP node_netstat_Ip6_InNoRoutes Statistic Ip6InNoRoutes.
	# TYPE node_netstat_Ip6_InNoRoutes untyped
	node_netstat_Ip6_InNoRoutes 5
	# HELP node_netstat_Ip6_InAddrErrors Statistic Ip6InAddrErrors.
	# TYPE node_netstat_Ip6_InAddrErrors untyped
	node_netstat_Ip6_InAddrErrors 0
	# HELP node_netstat_Ip6_InUnknownProtos Statistic Ip6InUnknownProtos.
	# TYPE node_netstat_Ip6_InUnknownProtos untyped
	node_netstat_Ip6_InUnknownProtos 0
	# HELP node_netstat_Ip6_InTruncatedPkts Statistic Ip6InTruncatedPkts.
	# TYPE node_netstat_Ip6_InTruncatedPkts untyped
	node_netstat_Ip6_InTruncatedPkts 0
	# HELP node_netstat_Ip6_InDiscards Statistic Ip6InDiscards.
	# TYPE node_netstat_Ip6_InDiscards untyped
	node_netstat_Ip6_InDiscards 0
	# HELP node_netstat_Ip6_InDelivers Statistic Ip6InDelivers.
	# TYPE node_netstat_Ip6_InDelivers untyped
	node_netstat_Ip6_InDelivers 0
	# HELP node_netstat_Ip6_OutForwDatagrams Statistic Ip6OutForwDatagrams.
	# TYPE node_netstat_Ip6_OutForwDatagrams untyped
	node_netstat_Ip6_OutForwDatagrams 0
	# HELP node_netstat_Ip6_OutRequests Statistic Ip6OutRequests.
	# TYPE node_netstat_Ip6_OutRequests untyped
	node_netstat_Ip6_OutRequests 8
	# HELP node_netstat_Ip6_OutDiscards Statistic Ip6OutDiscards.
	# TYPE node_netstat_Ip6_OutDiscards untyped
	node_netstat_Ip6_OutDiscards 0
	# HELP node_netstat_Ip6_OutNoRoutes Statistic Ip6OutNoRoutes.
	# TYPE node_netstat_Ip6_OutNoRoutes untyped
	node_netstat_Ip6_OutNoRoutes 3003
	# HELP node_netstat_Ip6_ReasmTimeout Statistic Ip6ReasmTimeout.
	# TYPE node_netstat_Ip6_ReasmTimeout untyped
	node_netstat_Ip6_ReasmTimeout 0
	# HELP node_netstat_Ip6_ReasmReqds Statistic Ip6ReasmReqds.
	# TYPE node_netstat_Ip6_ReasmReqds untyped
	node_netstat_Ip6_ReasmReqds 0
	# HELP node_netstat_Ip6_ReasmOKs Statistic Ip6ReasmOKs.
	# TYPE node_netstat_Ip6_ReasmOKs untyped
	node_netstat_Ip6_ReasmOKs 0
	# HELP node_netstat_Ip6_ReasmFails Statistic Ip6ReasmFails.
	# TYPE node_netstat_Ip6_ReasmFails untyped
	node_netstat_Ip6_ReasmFails 0
	# HELP node_netstat_Ip6_FragOKs Statistic Ip6FragOKs.
	# TYPE node_netstat_Ip6_FragOKs untyped
	node_netstat_Ip6_FragOKs 0
	# HELP node_netstat_Ip6_FragFails Statistic Ip6FragFails.
	# TYPE node_netstat_Ip6_FragFails untyped
	node_netstat_Ip6_FragFails 0
	# HELP node_netstat_Ip6_FragCreates Statistic Ip6FragCreates.
	# TYPE node_netstat_Ip6_FragCreates untyped
	node_netstat_Ip6_FragCreates 0
	# HELP node_netstat_Ip6_InMcastPkts Statistic Ip6InMcastPkts.
	# TYPE node_netstat_Ip6_InMcastPkts untyped
	node_netstat_Ip6_InMcastPkts 2
	# HELP node_netstat_Ip6_OutMcastPkts Statistic Ip6OutMcastPkts.
	# TYPE node_netstat_Ip6_OutMcastPkts untyped
	node_netstat_Ip6_OutMcastPkts 12
	# HELP node_netstat_Ip6_InOctets Statistic Ip6InOctets.
	# TYPE node_netstat_Ip6_InOctets untyped
	node_netstat_Ip6_InOctets 460
	# HELP node_netstat_Ip6_OutOctets Statistic Ip6OutOctets.
	# TYPE node_netstat_Ip6_OutOctets untyped
	node_netstat_Ip6_OutOctets 536
	# HELP node_netstat_Ip6_InMcastOctets Statistic Ip6InMcastOctets.
	# TYPE node_netstat_Ip6_InMcastOctets untyped
	node_netstat_Ip6_InMcastOctets 112
	# HELP node_netstat_Ip6_OutMcastOctets Statistic Ip6OutMcastOctets.
	# TYPE node_netstat_Ip6_OutMcastOctets untyped
	node_netstat_Ip6_OutMcastOctets 840
	# HELP node_netstat_Ip6_InBcastOctets Statistic Ip6InBcastOctets.
	# TYPE node_netstat_Ip6_InBcastOctets untyped
	node_netstat_Ip6_InBcastOctets 0
	# HELP node_netstat_Ip6_OutBcastOctets Statistic Ip6OutBcastOctets.
	# TYPE node_netstat_Ip6_OutBcastOctets untyped
	node_netstat_Ip6_OutBcastOctets 0
	# HELP node_netstat_Ip6_InNoECTPkts Statistic Ip6InNoECTPkts.
	# TYPE node_netstat_Ip6_InNoECTPkts untyped
	node_netstat_Ip6_InNoECTPkts 7
	# HELP node_netstat_Ip6_InECT1Pkts Statistic Ip6InECT1Pkts.
	# TYPE node_netstat_Ip6_InECT1Pkts untyped
	node_netstat_Ip6_InECT1Pkts 0
	# HELP node_netstat_Ip6_InECT0Pkts Statistic Ip6InECT0Pkts.
	# TYPE node_netstat_Ip6_InECT0Pkts untyped
	node_netstat_Ip6_InECT0Pkts 0
	# HELP node_netstat_Ip6_InCEPkts Statistic Ip6InCEPkts.
	# TYPE node_netstat_Ip6_InCEPkts untyped
	node_netstat_Ip6_InCEPkts 0
	# HELP node_netstat_Icmp6_InMsgs Statistic Icmp6InMsgs.
	# TYPE node_netstat_Icmp6_InMsgs untyped
	node_netstat_Icmp6_InMsgs 0
	# HELP node_netstat_Icmp6_InErrors Statistic Icmp6InErrors.
	# TYPE node_netstat_Icmp6_InErrors untyped
	node_netstat_Icmp6_InErrors 0
	# HELP node_netstat_Icmp6_OutMsgs Statistic Icmp6OutMsgs.
	# TYPE node_netstat_Icmp6_OutMsgs untyped
	node_netstat_Icmp6_OutMsgs 8
	# HELP node_netstat_Icmp6_OutErrors Statistic Icmp6OutErrors.
	# TYPE node_netstat_Icmp6_OutErrors untyped
	node_netstat_Icmp6_OutErrors 0
	# HELP node_netstat_Icmp6_InCsumErrors Statistic Icmp6InCsumErrors.
	# TYPE node_netstat_Icmp6_InCsumErrors untyped
	node_netstat_Icmp6_InCsumErrors 0
	# HELP node_netstat_Icmp6_InDestUnreachs Statistic Icmp6InDestUnreachs.
	# TYPE node_netstat_Icmp6_InDestUnreachs untyped
	node_netstat_Icmp6_InDestUnreachs 0
	# HELP node_netstat_Icmp6_InPktTooBigs Statistic Icmp6InPktTooBigs.
	# TYPE node_netstat_Icmp6_InPktTooBigs untyped
	node_netstat_Icmp6_InPktTooBigs 0
	# HELP node_netstat_Icmp6_InTimeExcds Statistic Icmp6InTimeExcds.
	# TYPE node_netstat_Icmp6_InTimeExcds untyped
	node_netstat_Icmp6_InTimeExcds 0
	# HELP node_netstat_Icmp6_InParmProblems Statistic Icmp6InParmProblems.
	# TYPE node_netstat_Icmp6_InParmProblems untyped
	node_netstat_Icmp6_InParmProblems 0
	# HELP node_netstat_Icmp6_InEchos Statistic Icmp6InEchos.
	# TYPE node_netstat_Icmp6_InEchos untyped
	node_netstat_Icmp6_InEchos 0
	# HELP node_netstat_Icmp6_InEchoReplies Statistic Icmp6InEchoReplies.
	# TYPE node_netstat_Icmp6_InEchoReplies untyped
	node_netstat_Icmp6_InEchoReplies 0
	# HELP node_netstat_Icmp6_InGroupMembQueries Statistic Icmp6InGroupMembQueries.
	# TYPE node_netstat_Icmp6_InGroupMembQueries untyped
	node_netstat_Icmp6_InGroupMembQueries 0
	# HELP node_netstat_Icmp6_InGroupMembResponses Statistic Icmp6InGroupMembResponses.
	# TYPE node_netstat_Icmp6_InGroupMembResponses untyped
	node_netstat_Icmp6_InGroupMembResponses 0
	# HELP node_netstat_Icmp6_InGroupMembReductions Statistic Icmp6InGroupMembReductions.
	# TYPE node_netstat_Icmp6_InGroupMembReductions untyped
	node_netstat_Icmp6_InGroupMembReductions 0
	# HELP node_netstat_Icmp6_InRouterSolicits Statistic Icmp6InRouterSolicits.
	# TYPE node_netstat_Icmp6_InRouterSolicits untyped
	node_netstat_Icmp6_InRouterSolicits 0
	# HELP node_netstat_Icmp6_InRouterAdvertisements Statistic Icmp6InRouterAdvertisements.
	# TYPE node_netstat_Icmp6_InRouterAdvertisements untyped
	node_netstat_Icmp6_InRouterAdvertisements 0
	# HELP node_netstat_Icmp6_InNeighborSolicits Statistic Icmp6InNeighborSolicits.
	# TYPE node_netstat_Icmp6_InNeighborSolicits untyped
	node_netstat_Icmp6_InNeighborSolicits 0
	# HELP node_netstat_Icmp6_InNeighborAdvertisements Statistic Icmp6InNeighborAdvertisements.
	# TYPE node_netstat_Icmp6_InNeighborAdvertisements untyped
	node_netstat_Icmp6_InNeighborAdvertisements 0
	# HELP node_netstat_Icmp6_InRedirects Statistic Icmp6InRedirects.
	# TYPE node_netstat_Icmp6_InRedirects untyped
	node_netstat_Icmp6_InRedirects 0
	# HELP node_netstat_Icmp6_InMLDv2Reports Statistic Icmp6InMLDv2Reports.
	# TYPE node_netstat_Icmp6_InMLDv2Reports untyped
	node_netstat_Icmp6_InMLDv2Reports 0
	# HELP node_netstat_Icmp6_OutDestUnreachs Statistic Icmp6OutDestUnreachs.
	# TYPE node_netstat_Icmp6_OutDestUnreachs untyped
	node_netstat_Icmp6_OutDestUnreachs 0
	# HELP node_netstat_Icmp6_OutPktTooBigs Statistic Icmp6OutPktTooBigs.
	# TYPE node_netstat_Icmp6_OutPktTooBigs untyped
	node_netstat_Icmp6_OutPktTooBigs 0
	# HELP node_netstat_Icmp6_OutTimeExcds Statistic Icmp6OutTimeExcds.
	# TYPE node_netstat_Icmp6_OutTimeExcds untyped
	node_netstat_Icmp6_OutTimeExcds 0
	# HELP node_netstat_Icmp6_OutParmProblems Statistic Icmp6OutParmProblems.
	# TYPE node_netstat_Icmp6_OutParmProblems untyped
	node_netstat_Icmp6_OutParmProblems 0
	# HELP node_netstat_Icmp6_OutEchos Statistic Icmp6OutEchos.
	# TYPE node_netstat_Icmp6_OutEchos untyped
	node_netstat_Icmp6_OutEchos 0
	# HELP node_netstat_Icmp6_OutEchoReplies Statistic Icmp6OutEchoReplies.
	# TYPE node_netstat_Icmp6_OutEchoReplies untyped
	node_netstat_Icmp6_OutEchoReplies 0
	# HELP node_netstat_Icmp6_OutGroupMembQueries Statistic Icmp6OutGroupMembQueries.
	# TYPE node_netstat_Icmp6_OutGroupMembQueries untyped
	node_netstat_Icmp6_OutGroupMembQueries 0
	# HELP node_netstat_Icmp6_OutGroupMembResponses Statistic Icmp6OutGroupMembResponses.
	# TYPE node_netstat_Icmp6_OutGroupMembResponses untyped
	node_netstat_Icmp6_OutGroupMembResponses 0
	# HELP node_netstat_Icmp6_OutGroupMembReductions Statistic Icmp6OutGroupMembReductions.
	# TYPE node_netstat_Icmp6_OutGroupMembReductions untyped
	node_netstat_Icmp6_OutGroupMembReductions 0
	# HELP node_netstat_Icmp6_OutRouterSolicits Statistic Icmp6OutRouterSolicits.
	# TYPE node_netstat_Icmp6_OutRouterSolicits untyped
	node_netstat_Icmp6_OutRouterSolicits 3
	# HELP node_netstat_Icmp6_OutRouterAdvertisements Statistic Icmp6OutRouterAdvertisements.
	# TYPE node_netstat_Icmp6_OutRouterAdvertisements untyped
	node_netstat_Icmp6_OutRouterAdvertisements 0
	# HELP node_netstat_Icmp6_OutNeighborSolicits Statistic Icmp6OutNeighborSolicits.
	# TYPE node_netstat_Icmp6_OutNeighborSolicits untyped
	node_netstat_Icmp6_OutNeighborSolicits 1
	# HELP node_netstat_Icmp6_OutNeighborAdvertisements Statistic Icmp6OutNeighborAdvertisements.
	# TYPE node_netstat_Icmp6_OutNeighborAdvertisements untyped
	node_netstat_Icmp6_OutNeighborAdvertisements 0
	# HELP node_netstat_Icmp6_OutRedirects Statistic Icmp6OutRedirects.
	# TYPE node_netstat_Icmp6_OutRedirects untyped
	node_netstat_Icmp6_OutRedirects 0
	# HELP node_netstat_Icmp6_OutMLDv2Reports Statistic Icmp6OutMLDv2Reports.
	# TYPE node_netstat_Icmp6_OutMLDv2Reports untyped
	node_netstat_Icmp6_OutMLDv2Reports 4
	# HELP node_netstat_Icmp6_OutType133 Statistic Icmp6OutType133.
	# TYPE node_netstat_Icmp6_OutType133 untyped
	node_netstat_Icmp6_OutType133 3
	# HELP node_netstat_Icmp6_OutType135 Statistic Icmp6OutType135.
	# TYPE node_netstat_Icmp6_OutType135 untyped
	node_netstat_Icmp6_OutType135 1
	# HELP node_netstat_Icmp6_OutType143 Statistic Icmp6OutType143.
	# TYPE node_netstat_Icmp6_OutType143 untyped
	node_netstat_Icmp6_OutType143 4
	# HELP node_netstat_Udp6_InDatagrams Statistic Udp6InDatagrams.
	# TYPE node_netstat_Udp6_InDatagrams untyped
	node_netstat_Udp6_InDatagrams 0
	# HELP node_netstat_Udp6_NoPorts Statistic Udp6NoPorts.
	# TYPE node_netstat_Udp6_NoPorts untyped
	node_netstat_Udp6_NoPorts 0
	# HELP node_netstat_Udp6_InErrors Statistic Udp6InErrors.
	# TYPE node_netstat_Udp6_InErrors untyped
	node_netstat_Udp6_InErrors 0
	# HELP node_netstat_Udp6_OutDatagrams Statistic Udp6OutDatagrams.
	# TYPE node_netstat_Udp6_OutDatagrams untyped
	node_netstat_Udp6_OutDatagrams 0
	# HELP node_netstat_Udp6_RcvbufErrors Statistic Udp6RcvbufErrors.
	# TYPE node_netstat_Udp6_RcvbufErrors untyped
	node_netstat_Udp6_RcvbufErrors 9
	# HELP node_netstat_Udp6_SndbufErrors Statistic Udp6SndbufErrors.
	# TYPE node_netstat_Udp6_SndbufErrors untyped
	node_netstat_Udp6_SndbufErrors 8
	# HELP node_netstat_Udp6_InCsumErrors Statistic Udp6InCsumErrors.
	# TYPE node_netstat_Udp6_InCsumErrors untyped
	node_netstat_Udp6_InCsumErrors 0
	# HELP node_netstat_Udp6_IgnoredMulti Statistic Udp6IgnoredMulti.
	# TYPE node_netstat_Udp6_IgnoredMulti untyped
	node_netstat_Udp6_IgnoredMulti 0
	# HELP node_netstat_UdpLite6_InDatagrams Statistic UdpLite6InDatagrams.
	# TYPE node_netstat_UdpLite6_InDatagrams untyped
	node_netstat_UdpLite6_InDatagrams 0
	# HELP node_netstat_UdpLite6_NoPorts Statistic UdpLite6NoPorts.
	# TYPE node_netstat_UdpLite6_NoPorts untyped
	node_netstat_UdpLite6_NoPorts 0
	# HELP node_netstat_UdpLite6_InErrors Statistic UdpLite6InErrors.
	# TYPE node_netstat_UdpLite6_InErrors untyped
	node_netstat_UdpLite6_InErrors 0
	# HELP node_netstat_UdpLite6_OutDatagrams Statistic UdpLite6OutDatagrams.
	# TYPE node_netstat_UdpLite6_OutDatagrams untyped
	node_netstat_UdpLite6_OutDatagrams 0
	# HELP node_netstat_UdpLite6_RcvbufErrors Statistic UdpLite6RcvbufErrors.
	# TYPE node_netstat_UdpLite6_RcvbufErrors untyped
	node_netstat_UdpLite6_RcvbufErrors 0
	# HELP node_netstat_UdpLite6_SndbufErrors Statistic UdpLite6SndbufErrors.
	# TYPE node_netstat_UdpLite6_SndbufErrors untyped
	node_netstat_UdpLite6_SndbufErrors 0
	# HELP node_netstat_UdpLite6_InCsumErrors Statistic UdpLite6InCsumErrors.
	# TYPE node_netstat_UdpLite6_InCsumErrors untyped
	node_netstat_UdpLite6_InCsumErrors 0
	`

	logger := log.NewLogfmtLogger(os.Stderr)
	collector, err := NewNetStatCollector(logger)
	if err != nil {
		panic(err)
	}
	c, err := newTestNetStatCollector(logger)
	if err != nil {
		t.Fatal(err)
	}
	reg := prometheus.NewRegistry()
	reg.MustRegister(c)

	sink := make(chan prometheus.Metric)
	go func() {
		err = collector.Update(sink)
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
