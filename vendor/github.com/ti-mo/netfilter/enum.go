package netfilter

// Subsystem specifiers for Netfilter Netlink messages
const (
	NFSubsysNone SubsystemID = iota // NFNL_SUBSYS_NONE

	NFSubsysCTNetlink        // NFNL_SUBSYS_CTNETLINK
	NFSubsysCTNetlinkExp     // NFNL_SUBSYS_CTNETLINK_EXP
	NFSubsysQueue            // NFNL_SUBSYS_QUEUE
	NFSubsysULOG             // NFNL_SUBSYS_ULOG
	NFSubsysOSF              // NFNL_SUBSYS_OSF
	NFSubsysIPSet            // NFNL_SUBSYS_IPSET
	NFSubsysAcct             // NFNL_SUBSYS_ACCT
	NFSubsysCTNetlinkTimeout // NFNL_SUBSYS_CTNETLINK_TIMEOUT
	NFSubsysCTHelper         // NFNL_SUBSYS_CTHELPER
	NFSubsysNFTables         // NFNL_SUBSYS_NFTABLES
	NFSubsysNFTCompat        // NFNL_SUBSYS_NFT_COMPAT
	NFSubsysCount            // NFNL_SUBSYS_COUNT
)

// ProtoFamily represents a protocol family in the Netfilter header (nfgenmsg).
type ProtoFamily uint8

// anonymous enum in uapi/linux/netfilter.h
const (
	ProtoUnspec ProtoFamily = 0  // NFPROTO_UNSPEC
	ProtoInet   ProtoFamily = 1  // NFPROTO_INET
	ProtoIPv4   ProtoFamily = 2  // NFPROTO_IPV4
	ProtoARP    ProtoFamily = 3  // NFPROTO_ARP
	ProtoNetDev ProtoFamily = 5  // NFPROTO_NETDEV
	ProtoBridge ProtoFamily = 7  // NFPROTO_BRIDGE
	ProtoIPv6   ProtoFamily = 10 // NFPROTO_IPV6
	ProtoDECNet ProtoFamily = 12 // NFPROTO_DECNET
)

// NetlinkGroup represents the multicast groups that can be joined with a Netlink socket.
type NetlinkGroup uint8

// enum nfnetlink_groups
const (
	GroupNone NetlinkGroup = iota // NFNLGRP_NONE

	GroupCTNew        // NFNLGRP_CONNTRACK_NEW
	GroupCTUpdate     // NFNLGRP_CONNTRACK_UPDATE
	GroupCTDestroy    // NFNLGRP_CONNTRACK_DESTROY
	GroupCTExpNew     // NFNLGRP_CONNTRACK_EXP_NEW
	GroupCTExpUpdate  // NFNLGRP_CONNTRACK_EXP_UPDATE
	GroupCTExpDestroy // NFNLGRP_CONNTRACK_EXP_DESTROY
	GroupNFTables     // NFNLGRP_NFTABLES
	GroupAcctQuota    // NFNLGRP_ACCT_QUOTA
	GroupNFTrace      // NFNLGRP_NFTRACE
)

var (
	// GroupsCT is a list of all Conntrack multicast groups.
	GroupsCT = []NetlinkGroup{GroupCTNew, GroupCTUpdate, GroupCTDestroy}
	// GroupsCTExp is a list of all Conntrack-expect multicast groups.
	GroupsCTExp = []NetlinkGroup{GroupCTExpNew, GroupCTExpUpdate, GroupCTExpDestroy}
)
