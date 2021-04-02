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

// +build !nochrony

package collector

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/facebookincubator/ntp/ntpcheck/checker"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	chronySubsystem = "chrony"
)

var (
	chronySocketPath      = kingpin.Flag("collector.chrony.socket-path", "chronyd Socket Path").Default("/var/run/chrony/chronyd.sock").String()
	chronyLogResponseJson = kingpin.Flag("collector.chrony.log-response-json", "Log chrony socket response as json through the debug level").Default("false").Bool()
)

type chronyCollector struct {

	// tracking response
	trackingLI, //int
	// trackingLIDesc, //str -> str rep of LI
	trackingClockSource, //str -> always shown as 1
	trackingCorrection, //float
	// trackingEvent, //str
	// trackingEventCount, //int

	// tracking as sysvars
	// trackingVersion, //str
	// trackingProcessor, //str
	// trackingSystem, //str
	// trackingLeap, //int -> same value as LI
	trackingStratum, //int
	trackingPrecision, //int
	trackingRootDelay, //float
	trackingRootDisp, //float
	// trackingPeer, //int
	// trackingTC, //int
	// trackingMinTC, //int
	// trackingClock, //str
	trackingRefID, //str -> parsed as float
	trackingRefTime, //str -> parsed as float
	trackingOffset, //float
	// trackingSysJitter, //int
	trackingFrequency, //float
	// trackingClkWander, //int
	// trackingClkJitter, //int
	// trackingTai, //int

	// sources response
	sourcesPeerConfigured, //bool -> parsed as 1/0
	sourcesPeerAuthPossible, //bool -> parsed as 1/0
	sourcesPeerAuthentic, //bool -> parsed as 1/0
	sourcesPeerReachable, //bool -> parsed as 1/0
	sourcesPeerBroadcast, //bool -> parsed as 1/0
	sourcesPeerSelection, //int
	sourcesPeerCondition, //str -> str rep of selection
	// sourcesPeerSRCAdr, //str
	// sourcesPeerSRCPort, //int
	// sourcesPeerDSTAdr, //str
	// sourcesPeerDSTPort, //int
	sourcesPeerLeap, //int
	sourcesPeerStratum, //int
	sourcesPeerPrecision, //int
	sourcesPeerRootDelay, //float
	sourcesPeerRootDisp, //float
	sourcesPeerRefID, //str -> parsed as float
	sourcesPeerRefTime, //str -> parsed as float
	sourcesPeerReach, //int
	sourcesPeerUnreach, //int
	sourcesPeerHMode, //int
	sourcesPeerPMode, //int
	sourcesPeerHPoll, //int
	sourcesPeerPPoll, //int
	sourcesPeerHeadway, //int
	sourcesPeerFlash, //int
	sourcesPeerOffset, //float
	sourcesPeerDelay, //float
	sourcesPeerDispersion, //float
	sourcesPeerJitter, //float
	sourcesPeerXleave, //float
	// sourcesPeerRec, //str
	// sourcesPeerFiltDelay, //str
	// sourcesPeerFiltOffset, //str
	// sourcesPeerFiltDisp, //str

	// serverstats response
	serverStatsPacketsReceived, //int
	serverStatsPacketsDropped, //int

	// incomplete flag
	ntpIncomplete, //bool -> parsed as 1/0

	// manually added metrics
	sourcesPeerCount typedDesc //int

	logger log.Logger
}

func init() {
	registerCollector("chrony", defaultDisabled, NewChronyCollector)
}

func NewChronyCollector(logger log.Logger) (Collector, error) {
	const subsystem = "chrony"
	peerLabels := []string{"id", "addr"}

	return &chronyCollector{
		trackingLI:          typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "tracking_leap_indicator"), "Tracking Leap Indicator. 0 - no warning, 3 - alarm", []string{"desc"}, nil), prometheus.GaugeValue},                                                                                                                                                                                                                                         //int
		trackingClockSource: typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "tracking_clock_source"), "Clock Source, str value in 'src' label, value always 1.", []string{"src"}, nil), prometheus.GaugeValue},                                                                                                                                                                                                                                       //int
		trackingCorrection:  typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "tracking_correction"), "Current correction value.", nil, nil), prometheus.GaugeValue},                                                                                                                                                                                                                                                                                   //float
		trackingStratum:     typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "tracking_stratum"), "The stratum indicates how many hops away from a computer with an attached reference clock we are.", nil, nil), prometheus.GaugeValue},                                                                                                                                                                                                              //int
		trackingPrecision:   typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "tracking_precision"), "Current precision.", nil, nil), prometheus.GaugeValue},                                                                                                                                                                                                                                                                                           //int
		trackingRootDelay:   typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "tracking_root_delay"), "Total of the network path delays to the stratum-1 computer from which the computer is ultimately synchronized.", nil, nil), prometheus.GaugeValue},                                                                                                                                                                                              //float
		trackingRootDisp:    typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "tracking_root_disp"), "Total dispersion accumulated through all the computers back to the stratum-1 computer from which the computer is ultimately synchronized.", nil, nil), prometheus.GaugeValue},                                                                                                                                                                    //float
		trackingRefID:       typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "tracking_ref_id"), "Encoded address of connected machine (if available).", nil, nil), prometheus.GaugeValue},                                                                                                                                                                                                                                                            //float
		trackingRefTime:     typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "tracking_ref_time"), "The time (UTC) at which the last measurement from the reference source was processed.", nil, nil), prometheus.GaugeValue},                                                                                                                                                                                                                         //float
		trackingOffset:      typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "tracking_offset"), "The estimated local offset on the last clock update.", nil, nil), prometheus.GaugeValue},                                                                                                                                                                                                                                                            //float
		trackingFrequency:   typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "tracking_frequency"), "The rate by which the system’s clock would be wrong if chronyd was not correcting it. It is expressed in ppm (parts per million). For example, a value of 1 ppm would mean that when the system’s clock thinks it has advanced 1 second, it has actually advanced by 1.000001 seconds relative to true time.", nil, nil), prometheus.GaugeValue}, //float

		sourcesPeerSelection:    typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "sources_peer_selection"), "State of the source (int code of *|+|-|?|x|~). str value in 'desc' label. See: https://github.com/facebookincubator/ntp/blob/master/protocol/chrony/packet.go#L81", append(peerLabels, "desc"), nil), prometheus.GaugeValue}, //int
		sourcesPeerOffset:       typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "sources_peer_offset"), "Offset of last update.", peerLabels, nil), prometheus.GaugeValue},                                                                                                                                                               //float
		sourcesPeerDelay:        typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "sources_peer_delay"), "Delay to peer.", peerLabels, nil), prometheus.GaugeValue},                                                                                                                                                                        //float
		sourcesPeerDispersion:   typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "sources_peer_dispersion"), "Peer dispersion.", peerLabels, nil), prometheus.GaugeValue},                                                                                                                                                                 //float
		sourcesPeerJitter:       typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "sources_peer_jitter"), "Peer jitter.", peerLabels, nil), prometheus.GaugeValue},                                                                                                                                                                         //float
		sourcesPeerRefID:        typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "sources_peer_ref_id"), "Peer refId (see tracking_ref_id).", peerLabels, nil), prometheus.GaugeValue},                                                                                                                                                    //float
		sourcesPeerRefTime:      typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "sources_peer_ref_time"), "Peer refTime (see tracking_ref_time).", peerLabels, nil), prometheus.GaugeValue},                                                                                                                                              //float
		sourcesPeerRootDelay:    typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "sources_peer_root_delay"), "Peer root delay (see tracking_root_delay).", peerLabels, nil), prometheus.GaugeValue},                                                                                                                                       //float
		sourcesPeerRootDisp:     typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "sources_peer_root_disp"), "Peer root dispersion (see tracking_root_disp).", peerLabels, nil), prometheus.GaugeValue},                                                                                                                                    //float
		sourcesPeerConfigured:   typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "sources_peer_peer_configured"), "Configured flag (1|0).", peerLabels, nil), prometheus.GaugeValue},                                                                                                                                                      //bool
		sourcesPeerAuthPossible: typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "sources_peer_auth_possible"), "AuthPosible flag (1|0).", peerLabels, nil), prometheus.GaugeValue},                                                                                                                                                       //bool
		sourcesPeerAuthentic:    typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "sources_peer_authentic"), "Authenticatd flag (0|1).", peerLabels, nil), prometheus.GaugeValue},                                                                                                                                                          //bool
		sourcesPeerReachable:    typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "sources_peer_reachable"), "Reachability (1 if sourceReply Reachability flag == 255 else 0).", peerLabels, nil), prometheus.GaugeValue},                                                                                                                  //bool
		sourcesPeerBroadcast:    typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "sources_peer_broadcast"), "Broadcast flag (1|0).", peerLabels, nil), prometheus.GaugeValue},                                                                                                                                                             //bool
		sourcesPeerLeap:         typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "sources_peer_leap"), "Peer leap value (see tracking_leap_indicator).", peerLabels, nil), prometheus.GaugeValue},                                                                                                                                         //int
		sourcesPeerStratum:      typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "sources_peer_stratum"), "Peer startum (see tracking_stratum).", peerLabels, nil), prometheus.GaugeValue},                                                                                                                                                //int
		sourcesPeerPrecision:    typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "sources_peer_precision"), "Peer precision (see tracking_precision).", peerLabels, nil), prometheus.GaugeValue},                                                                                                                                          //int
		sourcesPeerReach:        typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "sources_peer_reach"), "Int value of sourceReply Reachability flag (see sources_peer_reachable).", peerLabels, nil), prometheus.GaugeValue},                                                                                                              //int
		sourcesPeerUnreach:      typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "sources_peer_unreach"), "Unreach flag from .", peerLabels, nil), prometheus.GaugeValue},                                                                                                                                                                 //int
		sourcesPeerHMode:        typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "sources_peer_hmode"), "Hmode value from source data.", peerLabels, nil), prometheus.GaugeValue},                                                                                                                                                         //int
		sourcesPeerPMode:        typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "sources_peer_pmode"), "Pmode value from source data.", peerLabels, nil), prometheus.GaugeValue},                                                                                                                                                         //int
		sourcesPeerHPoll:        typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "sources_peer_hpoll"), "Hpoll value from source data.", peerLabels, nil), prometheus.GaugeValue},                                                                                                                                                         //int
		sourcesPeerPPoll:        typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "sources_peer_ppoll"), "Ppoll value from source data.", peerLabels, nil), prometheus.GaugeValue},                                                                                                                                                         //int
		sourcesPeerHeadway:      typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "sources_peer_headway"), "Headway value from source data.", peerLabels, nil), prometheus.GaugeValue},                                                                                                                                                     //int
		sourcesPeerXleave:       typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "sources_peer_xleave"), "Xleave value from source data.", peerLabels, nil), prometheus.GaugeValue},                                                                                                                                                       //float

		serverStatsPacketsReceived: typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "server_stats_packets_received"), "Packets received (only available if serverstats request was succeeded).", nil, nil), prometheus.GaugeValue}, //int
		serverStatsPacketsDropped:  typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "server_stats_packets_dropped"), "Packets dropped (only available if serverstats request was succeeded).", nil, nil), prometheus.GaugeValue},   //int

		ntpIncomplete: typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "ntp_incomplete"), "0 if some of ntpData requests was failed (1|0).", nil, nil), prometheus.GaugeValue}, //bool

		sourcesPeerCount: typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "sources_peer_count"), "Numbers of total peers.", nil, nil), prometheus.GaugeValue}, //int

		logger: logger,
	}, nil
}

func (c *chronyCollector) Update(ch chan<- prometheus.Metric) error {
	var b2i = map[bool]float64{false: 0, true: 1}

	resp, err := checker.RunCheck(*chronySocketPath)
	// error handling
	if err != nil {
		level.Debug(c.logger).Log("msg", "request to chronyd failed", "err", err)
		return ErrNoData
	}

	// response logging if needed
	if *chronyLogResponseJson {
		j, err := json.Marshal(resp)
		if err == nil {
			level.Debug(c.logger).Log("msg", j)
		}
	}

	// tracking value subset
	ch <- c.trackingLI.mustNewConstMetric(float64(resp.LI), resp.LIDesc)
	ch <- c.trackingClockSource.mustNewConstMetric(1, resp.ClockSource)
	ch <- c.trackingCorrection.mustNewConstMetric(resp.Correction)
	ch <- c.trackingStratum.mustNewConstMetric(float64(resp.SysVars.Stratum))
	ch <- c.trackingPrecision.mustNewConstMetric(float64(resp.SysVars.Precision))
	ch <- c.trackingRootDelay.mustNewConstMetric(resp.SysVars.RootDelay / 1e3) //initially in ms, see: https://github.com/facebookincubator/ntp/blob/143e098237b0161198f3057998fdf8773c42d612/ntpcheck/checker/system.go#L66
	ch <- c.trackingRootDisp.mustNewConstMetric(resp.SysVars.RootDisp / 1e3)
	ch <- c.trackingOffset.mustNewConstMetric(resp.SysVars.Offset / 1e3)
	ch <- c.trackingFrequency.mustNewConstMetric(resp.SysVars.Frequency)

	// refid parsing (hex as float)
	refId, err := strconv.ParseInt(resp.SysVars.RefID, 16, 64)
	if err == nil {
		ch <- c.trackingRefID.mustNewConstMetric(float64(refId))
	}

	// refTime parsing (str as float)
	t, _ := time.Parse("2006-01-02 15:04:05.99 -0700 MST", resp.SysVars.RefTime)
	if t.Unix() > 0 {
		// Go Zero is   0001-01-01 00:00:00 UTC
		// NTP Zero is  1900-01-01 00:00:00 UTC
		// UNIX Zero is 1970-01-01 00:00:00 UTC
		// so let's keep ALL ancient `reftime` values as zero
		ch <- c.trackingRefTime.mustNewConstMetric(float64(t.UnixNano()) / 1e9)
	} else {
		ch <- c.trackingRefTime.mustNewConstMetric(0)
	}

	// peers value subset
	for id, peer := range resp.Peers {
		peerLabelValues := []string{fmt.Sprintf("%d", id), peer.SRCAdr} //peer address appears as SRCAdr

		//floats
		ch <- c.sourcesPeerOffset.mustNewConstMetric(peer.Offset/1e3, peerLabelValues...) //initially in ms, see: https://github.com/facebookincubator/ntp/blob/81cb02c05f82f8c9cdf32e16f4ee02a3b05bfaf1/ntpcheck/checker/peer.go#L196
		ch <- c.sourcesPeerDelay.mustNewConstMetric(peer.Delay/1e3, peerLabelValues...)
		ch <- c.sourcesPeerDispersion.mustNewConstMetric(peer.Dispersion/1e3, peerLabelValues...)
		ch <- c.sourcesPeerJitter.mustNewConstMetric(peer.Jitter/1e3, peerLabelValues...)
		ch <- c.sourcesPeerRootDelay.mustNewConstMetric(peer.RootDelay/1e3, peerLabelValues...)
		ch <- c.sourcesPeerRootDisp.mustNewConstMetric(peer.RootDisp/1e3, peerLabelValues...)

		//booleans
		ch <- c.sourcesPeerConfigured.mustNewConstMetric(b2i[peer.Configured], peerLabelValues...)
		ch <- c.sourcesPeerAuthPossible.mustNewConstMetric(b2i[peer.AuthPossible], peerLabelValues...)
		ch <- c.sourcesPeerAuthentic.mustNewConstMetric(b2i[peer.Authentic], peerLabelValues...)
		ch <- c.sourcesPeerReachable.mustNewConstMetric(b2i[peer.Reachable], peerLabelValues...)
		ch <- c.sourcesPeerBroadcast.mustNewConstMetric(b2i[peer.Broadcast], peerLabelValues...)

		//integers
		ch <- c.sourcesPeerSelection.mustNewConstMetric(float64(peer.Selection), append(peerLabelValues, peer.Condition)...)
		ch <- c.sourcesPeerLeap.mustNewConstMetric(float64(peer.Leap), peerLabelValues...)
		ch <- c.sourcesPeerStratum.mustNewConstMetric(float64(peer.Stratum), peerLabelValues...)
		ch <- c.sourcesPeerPrecision.mustNewConstMetric(float64(peer.Precision), peerLabelValues...)
		ch <- c.sourcesPeerReach.mustNewConstMetric(float64(peer.Reach), peerLabelValues...)
		ch <- c.sourcesPeerUnreach.mustNewConstMetric(float64(peer.Unreach), peerLabelValues...)
		ch <- c.sourcesPeerHMode.mustNewConstMetric(float64(peer.HMode), peerLabelValues...)
		ch <- c.sourcesPeerPMode.mustNewConstMetric(float64(peer.PMode), peerLabelValues...)
		ch <- c.sourcesPeerHPoll.mustNewConstMetric(float64(peer.HPoll), peerLabelValues...)
		ch <- c.sourcesPeerPPoll.mustNewConstMetric(float64(peer.PPoll), peerLabelValues...)
		ch <- c.sourcesPeerHeadway.mustNewConstMetric(float64(peer.Headway), peerLabelValues...)
		ch <- c.sourcesPeerXleave.mustNewConstMetric(float64(peer.Xleave), peerLabelValues...)

		// refid parsing (hex as float)
		refId, err := strconv.ParseInt(peer.RefID, 16, 64)
		if err == nil {
			ch <- c.sourcesPeerRefID.mustNewConstMetric(float64(refId), peerLabelValues...)
		}

		// refTime parsing (str as float)
		t, _ := time.Parse("2006-01-02 15:04:05.99 -0700 MST", peer.RefTime)
		if t.Unix() > 0 {
			ch <- c.sourcesPeerRefTime.mustNewConstMetric(float64(t.UnixNano())/1e9, peerLabelValues...)
		} else {
			ch <- c.sourcesPeerRefTime.mustNewConstMetric(0, peerLabelValues...)
		}
	}

	// server stats value subset
	if resp.ServerStats != nil {
		ch <- c.serverStatsPacketsReceived.mustNewConstMetric(float64(resp.ServerStats.PacketsReceived))
		ch <- c.serverStatsPacketsDropped.mustNewConstMetric(float64(resp.ServerStats.PacketsDropped))
	}
	// incomplete flag
	ch <- c.ntpIncomplete.mustNewConstMetric(b2i[resp.Incomplete])

	//manually added metrics
	ch <- c.sourcesPeerCount.mustNewConstMetric(float64(len(resp.Peers)))

	return nil
}

/* chronyChecker response sample
{
  "LI": 0,
  "LIDesc": "none",
  "ClockSource": "ntp",
  "Correction": -9.862084880296607e-06,
  "Event": "clock_sync",
  "EventCount": 0,
  "SysVars": {
    "Version": "",
    "Processor": "",
    "System": "",
    "Leap": 0,
    "Stratum": 2,
    "Precision": 0,
    "RootDelay": 11.23967207968235,
    "RootDisp": 21883.459091186523,
    "Peer": 0,
    "TC": 0,
    "MinTC": 0,
    "Clock": "",
    "RefID": "51D32512",
    "RefTime": "2021-03-26 14:55:33.784183547 +0000 UTC",
    "Offset": 6.212211214005947,
    "SysJitter": 0,
    "Frequency": -1923.4349365234375,
    "ClkWander": 0,
    "ClkJitter": 0,
    "Tai": 0
  },
  "Peers": {
    "0": {
      "Configured": true,
      "AuthPossible": false,
      "Authentic": false,
      "Reachable": false,
      "Broadcast": false,
      "Selection": 0,
      "Condition": "unreach",
      "SRCAdr": "2001:67c:380:120::33",
      "SRCPort": 0,
      "DSTAdr": "::",
      "DSTPort": 0,
      "Leap": 0,
      "Stratum": 0,
      "Precision": 0,
      "RootDelay": 0,
      "RootDisp": 0,
      "RefID": "00000000",
      "RefTime": "1970-01-01 00:00:00 +0000 UTC",
      "Reach": 0,
      "Unreach": 0,
      "HMode": 0,
      "PMode": 0,
      "HPoll": 0,
      "PPoll": 0,
      "Headway": 0,
      "Flash": 0,
      "Flashers": [
        "pkt_auth",
        "tst_max_delay",
        "tst_delay_ratio",
        "pkt_dup",
        "pkt_invalid",
        "pkt_stratum",
        "pkt_header",
        "tst_delay_dev_ration",
        "tst_sync_loop",
        "pkt_bogus"
      ],
      "Offset": 0,
      "Delay": 0,
      "Dispersion": 0,
      "Jitter": 0,
      "Xleave": 0,
      "Rec": "",
      "FiltDelay": "",
      "FiltOffset": "",
      "FiltDisp": ""
    }
	},
  "ServerStats": {
    "ntp.server.packets_received": 0,
    "ntp.server.packets_dropped": 0
  },
  "Incomplete": false
}
*/
