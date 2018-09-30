package conntrack

import "github.com/ti-mo/netfilter"

// All enums in this file are translated from the Linux kernel source at
// include/uapi/linux/netfilter/nfnetlink_conntrack.h

// messageType is a Conntrack-specific representation of a netfilter.MessageType.
// It is used to specify the type of action to execute on the kernel's state table
// (get, create, delete, etc.).
type messageType netfilter.MessageType

// The first three members are similar to NF_NETLINK_CONNTRACK_*, which is still used
// in libnetfilter_conntrack. They can still be used to subscribe to Netlink groups with bind(),
// but subscribing using setsockopt() (like mdlayher/netlink) requires the NFNLGRP_* enum.
//
// enum cntl_msg_types (upstream typo)
const (
	ctNew            messageType = iota // IPCTNL_MSG_CT_NEW
	ctGet                               // IPCTNL_MSG_CT_GET
	ctDelete                            // IPCTNL_MSG_CT_DELETE
	ctGetCtrZero                        // IPCTNL_MSG_CT_GET_CTRZERO
	ctGetStatsCPU                       // IPCTNL_MSG_CT_GET_STATS_CPU
	ctGetStats                          // IPCTNL_MSG_CT_GET_STATS
	ctGetDying                          // IPCTNL_MSG_CT_GET_DYING
	ctGetUnconfirmed                    // IPCTNL_MSG_CT_GET_UNCONFIRMED
)

// expMessageType is a Conntrack-specific representation of a netfilter.MessageType.
// It holds information about Conntrack Expect events; state created by Conntrack helpers.
type expMessageType netfilter.MessageType

// enum ctnl_exp_msg_types
const (
	ctExpNew         expMessageType = iota // IPCTNL_MSG_EXP_NEW
	ctExpGet                               // IPCTNL_MSG_EXP_GET
	ctExpDelete                            // IPCTNL_MSG_EXP_DELETE
	ctExpGetStatsCPU                       // IPCTNL_MSG_EXP_GET_STATS_CPU
)

// attributeType defines the meaning of a root-level Type
// value of a Conntrack-specific Netfilter attribute.
type attributeType uint8

// enum ctattr_type
const (
	ctaUnspec        attributeType = iota // CTA_UNSPEC
	ctaTupleOrig                          // CTA_TUPLE_ORIG
	ctaTupleReply                         // CTA_TUPLE_REPLY
	ctaStatus                             // CTA_STATUS
	ctaProtoInfo                          // CTA_PROTOINFO
	ctaHelp                               // CTA_HELP
	ctaNatSrc                             // CTA_NAT_SRC, Deprecated
	ctaTimeout                            // CTA_TIMEOUT
	ctaMark                               // CTA_MARK
	ctaCountersOrig                       // CTA_COUNTERS_ORIG
	ctaCountersReply                      // CTA_COUNTERS_REPLY
	ctaUse                                // CTA_USE
	ctaID                                 // CTA_ID
	ctaNatDst                             // CTA_NAT_DST, Deprecated
	ctaTupleMaster                        // CTA_TUPLE_MASTER
	ctaSeqAdjOrig                         // CTA_SEQ_ADJ_ORIG
	ctaSeqAdjReply                        // CTA_SEQ_ADJ_REPLY
	ctaSecMark                            // CTA_SECMARK, Deprecated
	ctaZone                               // CTA_ZONE
	ctaSecCtx                             // CTA_SECCTX
	ctaTimestamp                          // CTA_TIMESTAMP
	ctaMarkMask                           // CTA_MARK_MASK
	ctaLabels                             // CTA_LABELS
	ctaLabelsMask                         // CTA_LABELS_MASK
	ctaSynProxy                           // CTA_SYNPROXY
)

// tupleType describes the type of tuple contained in this container.
type tupleType uint8

// enum ctattr_tuple
const (
	ctaTupleUnspec tupleType = iota //CTA_TUPLE_UNSPEC
	ctaTupleIP                      // CTA_TUPLE_IP
	ctaTupleProto                   // CTA_TUPLE_PROTO
	ctaTupleZone                    // CTA_TUPLE_ZONE
)

// protoTupleType describes the type of Layer 4 protocol metadata in this container.
type protoTupleType uint8

// enum ctattr_l4proto
const (
	ctaProtoUnspec     protoTupleType = iota // CTA_PROTO_UNSPEC
	ctaProtoNum                              // CTA_PROTO_NUM
	ctaProtoSrcPort                          // CTA_PROTO_SRC_PORT
	ctaProtoDstPort                          // CTA_PROTO_DST_PORT
	ctaProtoICMPID                           // CTA_PROTO_ICMP_ID
	ctaProtoICMPType                         // CTA_PROTO_ICMP_TYPE
	ctaProtoICMPCode                         // CTA_PROTO_ICMP_CODE
	ctaProtoICMPv6ID                         // CTA_PROTO_ICMPV6_ID
	ctaProtoICMPv6Type                       // CTA_PROTO_ICMPV6_TYPE
	ctaProtoICMPv6Code                       // CTA_PROTO_ICMPV6_CODE
)

// ipTupleType describes the type of IP address in this container.
type ipTupleType uint8

// enum ctattr_ip
const (
	ctaIPUnspec ipTupleType = iota // CTA_IP_UNSPEC
	ctaIPv4Src                     // CTA_IP_V4_SRC
	ctaIPv4Dst                     // CTA_IP_V4_DST
	ctaIPv6Src                     // CTA_IP_V6_SRC
	ctaIPv6Dst                     // CTA_IP_V6_DST
)

// helperType describes the kind of helper in this container.
type helperType uint8

// enum ctattr_help
const (
	ctaHelpUnspec helperType = iota // CTA_HELP_UNSPEC
	ctaHelpName                     // CTA_HELP_NAME
	ctaHelpInfo                     // CTA_HELP_INFO
)

// counterType describes the kind of counter in this container.
type counterType uint8

// enum ctattr_counters
const (
	ctaCountersUnspec  counterType = iota // CTA_COUNTERS_UNSPEC
	ctaCountersPackets                    // CTA_COUNTERS_PACKETS
	ctaCountersBytes                      // CTA_COUNTERS_BYTES
)

// timestampType describes the type of timestamp in this container.
type timestampType uint8

// enum ctattr_tstamp
const (
	ctaTimestampUnspec timestampType = iota // CTA_TIMESTAMP_UNSPEC
	ctaTimestampStart                       // CTA_TIMESTAMP_START
	ctaTimestampStop                        // CTA_TIMESTAMP_STOP
	ctaTimestampPad                         // CTA_TIMESTAMP_PAD
)

// securityType describes the type of SecCtx value in this container.
type securityType uint8

// enum ctattr_secctx
const (
	ctaSecCtxUnspec securityType = iota // CTA_SECCTX_UNSPEC
	ctaSecCtxName                       // CTA_SECCTX_NAME
)

// protoInfoType describes the kind of protocol info in this container.
type protoInfoType uint8

// enum ctattr_protoinfo
const (
	ctaProtoInfoUnspec protoInfoType = iota // CTA_PROTOINFO_UNSPEC
	ctaProtoInfoTCP                         // CTA_PROTOINFO_TCP
	ctaProtoInfoDCCP                        // CTA_PROTOINFO_DCCP
	ctaProtoInfoSCTP                        // CTA_PROTOINFO_SCTP
)

// protoInfoTCPType describes the kind of TCP protocol info attribute in this container.
type protoInfoTCPType uint8

// enum ctattr_protoinfo_tcp
const (
	ctaProtoInfoTCPUnspec         protoInfoTCPType = iota // CTA_PROTOINFO_TCP_UNSPEC
	ctaProtoInfoTCPState                                  // CTA_PROTOINFO_TCP_STATE
	ctaProtoInfoTCPWScaleOriginal                         // CTA_PROTOINFO_TCP_WSCALE_ORIGINAL
	ctaProtoInfoTCPWScaleReply                            // CTA_PROTOINFO_TCP_WSCALE_REPLY
	ctaProtoInfoTCPFlagsOriginal                          // CTA_PROTOINFO_TCP_FLAGS_ORIGINAL
	ctaProtoInfoTCPFlagsReply                             // CTA_PROTOINFO_TCP_FLAGS_REPLY
)

// protoInfoDCCPType describes the kind of DCCP protocol info attribute in this container.
type protoInfoDCCPType uint8

// enum ctattr_protoinfo_dccp
const (
	ctaProtoInfoDCCPUnspec       protoInfoDCCPType = iota // CTA_PROTOINFO_DCCP_UNSPEC
	ctaProtoInfoDCCPState                                 // CTA_PROTOINFO_DCCP_STATE
	ctaProtoInfoDCCPRole                                  // CTA_PROTOINFO_DCCP_ROLE
	ctaProtoInfoDCCPHandshakeSeq                          // CTA_PROTOINFO_DCCP_HANDSHAKE_SEQ
	ctaProtoInfoDCCPPad                                   // CTA_PROTOINFO_DCCP_PAD (never sent by kernel)
)

// protoInfoSCTPType describes the kind of SCTP protocol info attribute in this container.
type protoInfoSCTPType uint8

// enum ctattr_protoinfo_sctp
const (
	ctaProtoInfoSCTPUnspec       protoInfoSCTPType = iota // CTA_PROTOINFO_SCTP_UNSPEC
	ctaProtoInfoSCTPState                                 // CTA_PROTOINFO_SCTP_STATE
	ctaProtoInfoSCTPVTagOriginal                          // CTA_PROTOINFO_SCTP_VTAG_ORIGINAL
	ctaProtoInfoSCTPVtagReply                             // CTA_PROTOINFO_SCTP_VTAG_REPLY
)

// seqAdjType describes the type of sequence adjustment in this container.
type seqAdjType uint8

// enum ctattr_seqadj
const (
	ctaSeqAdjUnspec        seqAdjType = iota // CTA_SEQADJ_UNSPEC
	ctaSeqAdjCorrectionPos                   // CTA_SEQADJ_CORRECTION_POS
	ctaSeqAdjOffsetBefore                    // CTA_SEQADJ_OFFSET_BEFORE
	ctaSeqAdjOffsetAfter                     // CTA_SEQADJ_OFFSET_AFTER
)

// synProxyType describes the type of SYNproxy attribute in this container.
type synProxyType uint8

// enum ctattr_synproxy
const (
	ctaSynProxyUnspec synProxyType = iota // CTA_SYNPROXY_UNSPEC
	ctaSynProxyISN                        // CTA_SYNPROXY_ISN
	ctaSynProxyITS                        // CTA_SYNPROXY_ITS
	ctaSynProxyTSOff                      // CTA_SYNPROXY_TSOFF
)

// expectType describes the type of expect attribute in this container.
type expectType uint8

// enum ctattr_expect
const (
	ctaExpectUnspec   expectType = iota // CTA_EXPECT_UNSPEC
	ctaExpectMaster                     // CTA_EXPECT_MASTER
	ctaExpectTuple                      // CTA_EXPECT_TUPLE
	ctaExpectMask                       // CTA_EXPECT_MASK
	ctaExpectTimeout                    // CTA_EXPECT_TIMEOUT
	ctaExpectID                         // CTA_EXPECT_ID
	ctaExpectHelpName                   // CTA_EXPECT_HELP_NAME
	ctaExpectZone                       // CTA_EXPECT_ZONE
	ctaExpectFlags                      // CTA_EXPECT_FLAGS
	ctaExpectClass                      // CTA_EXPECT_CLASS
	ctaExpectNAT                        // CTA_EXPECT_NAT
	ctaExpectFN                         // CTA_EXPECT_FN
)

// expectNATType describes the type of NAT expect attribute in this container.
type expectNATType uint8

// enum ctattr_expect_nat
const (
	ctaExpectNATUnspec expectNATType = iota // CTA_EXPECT_NAT_UNSPEC
	ctaExpectNATDir                         // CTA_EXPECT_NAT_DIR
	ctaExpectNATTuple                       // CTA_EXPECT_NAT_TUPLE
)

// cpuStatsType describes the type of CPU-specific conntrack statistics attribute in this container.
type cpuStatsType uint8

// ctattr_stats_cpu
const (
	ctaStatsUnspec        cpuStatsType = iota // CTA_STATS_UNSPEC
	ctaStatsSearched                          // CTA_STATS_SEARCHED, no longer used
	ctaStatsFound                             // CTA_STATS_FOUND
	ctaStatsNew                               // CTA_STATS_NEW, no longer used
	ctaStatsInvalid                           // CTA_STATS_INVALID
	ctaStatsIgnore                            // CTA_STATS_IGNORE
	ctaStatsDelete                            // CTA_STATS_DELETE, no longer used
	ctaStatsDeleteList                        // CTA_STATS_DELETE_LIST, no longer used
	ctaStatsInsert                            // CTA_STATS_INSERT
	ctaStatsInsertFailed                      // CTA_STATS_INSERT_FAILED
	ctaStatsDrop                              // CTA_STATS_DROP
	ctaStatsEarlyDrop                         // CTA_STATS_EARLY_DROP
	ctaStatsError                             // CTA_STATS_ERROR
	ctaStatsSearchRestart                     // CTA_STATS_SEARCH_RESTART
)

// globalStatsType describes the type of global conntrack statistics attribute in this container.
type globalStatsType uint8

// enum ctattr_stats_global
const (
	ctaStatsGlobalUnspec  globalStatsType = iota // CTA_STATS_GLOBAL_UNSPEC
	ctaStatsGlobalEntries                        // CTA_STATS_GLOBAL_ENTRIES
)

// expectStatsType describes the type of expectation statistics attribute in this container.
type expectStatsType uint8

// enum ctattr_expect_stats
const (
	ctaStatsExpUnspec expectStatsType = iota // CTA_STATS_EXP_UNSPEC
	ctaStatsExpNew                           // CTA_STATS_EXP_NEW
	ctaStatsExpCreate                        // CTA_STATS_EXP_CREATE
	ctaStatsExpDelete                        // CTA_STATS_EXP_DELETE
)

// enum ctattr_natseq is unused in the kernel source
