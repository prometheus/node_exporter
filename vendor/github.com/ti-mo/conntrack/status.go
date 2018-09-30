package conntrack

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/ti-mo/netfilter"
)

const (
	opUnStatus = "Status unmarshal"
)

// Status represents a snapshot of a conntrack connection's state.
type Status struct {
	Value StatusFlag
}

// unmarshal unmarshals a netfilter.Attribute into a Status structure.
func (s *Status) unmarshal(attr netfilter.Attribute) error {

	if attributeType(attr.Type) != ctaStatus {
		return fmt.Errorf(errAttributeWrongType, attr.Type, ctaStatus)
	}

	if attr.Nested {
		return errors.Wrap(errNested, opUnStatus)
	}

	if len(attr.Data) != 4 {
		return errors.Wrap(errIncorrectSize, opUnStatus)
	}

	s.Value = StatusFlag(attr.Uint32())

	return nil
}

// marshal marshals a Status into a netfilter.Attribute.
func (s Status) marshal() netfilter.Attribute {
	return netfilter.Attribute{
		Type: uint16(ctaStatus),
		Data: netfilter.Uint32Bytes(uint32(s.Value)),
	}
}

// Expected indicates that this connection is an expected connection,
// created by Conntrack helpers based on the state of another, related connection.
func (s Status) Expected() bool {
	return s.Value&StatusExpected != 0
}

// SeenReply is set when the flow has seen traffic both ways.
func (s Status) SeenReply() bool {
	return s.Value&StatusSeenReply != 0
}

// Assured is set when eg. three-way handshake is completed on a TCP flow.
func (s Status) Assured() bool {
	return s.Value&StatusAssured != 0
}

// Confirmed is set when the original packet has left the box.
func (s Status) Confirmed() bool {
	return s.Value&StatusConfirmed != 0
}

// SrcNAT means the connection needs source NAT in the original direction.
func (s Status) SrcNAT() bool {
	return s.Value&StatusSrcNAT != 0
}

// DstNAT means the connection needs destination NAT in the original direction.
func (s Status) DstNAT() bool {
	return s.Value&StatusDstNAT != 0
}

// SeqAdjust means the connection needs its TCP sequence to be adjusted.
func (s Status) SeqAdjust() bool {
	return s.Value&StatusSeqAdjust != 0
}

// SrcNATDone is set when source NAT was applied onto the connection.
func (s Status) SrcNATDone() bool {
	return s.Value&StatusSrcNATDone != 0
}

// DstNATDone is set when destination NAT was applied onto the connection.
func (s Status) DstNATDone() bool {
	return s.Value&StatusDstNATDone != 0
}

// Dying means the connection has concluded and needs to be cleaned up by GC.
func (s Status) Dying() bool {
	return s.Value&StatusDying != 0
}

// FixedTimeout means the connection's timeout value cannot be changed.
func (s Status) FixedTimeout() bool {
	return s.Value&StatusFixedTimeout != 0
}

// Template indicates if the connection is a template.
func (s Status) Template() bool {
	return s.Value&StatusTemplate != 0
}

// Helper is set when a helper was explicitly attached using a Conntrack target.
func (s Status) Helper() bool {
	return s.Value&StatusHelper != 0
}

// Offload is set when the connection was offloaded to flow table.
func (s Status) Offload() bool {
	return s.Value&StatusOffload != 0
}

// StatusFlag describes a status bit in a Status structure.
type StatusFlag uint32

// Conntrack connection's status flags, from enum ip_conntrack_status.
// uapi/linux/netfilter/nf_conntrack_common.h
const (
	StatusExpected  StatusFlag = 1      // IPS_EXPECTED
	StatusSeenReply StatusFlag = 1 << 1 // IPS_SEEN_REPLY
	StatusAssured   StatusFlag = 1 << 2 // IPS_ASSURED
	StatusConfirmed StatusFlag = 1 << 3 // IPS_CONFIRMED
	StatusSrcNAT    StatusFlag = 1 << 4 // IPS_SRC_NAT
	StatusDstNAT    StatusFlag = 1 << 5 // IPS_DST_NAT

	StatusNATMask = StatusDstNAT | StatusSrcNAT // IPS_NAT_MASK

	StatusSeqAdjust  StatusFlag = 1 << 6 // IPS_SEQ_ADJUST
	StatusSrcNATDone StatusFlag = 1 << 7 // IPS_SRC_NAT_DONE
	StatusDstNATDone StatusFlag = 1 << 8 // IPS_DST_NAT_DONE

	StatusNATDoneMask = StatusDstNATDone | StatusSrcNATDone // IPS_NAT_DONE_MASK

	StatusDying        StatusFlag = 1 << 9
	StatusFixedTimeout StatusFlag = 1 << 10 // IPS_FIXED_TIMEOUT
	StatusTemplate     StatusFlag = 1 << 11 // IPS_TEMPLATE
	StatusUntracked    StatusFlag = 1 << 12 // IPS_UNTRACKED
	StatusHelper       StatusFlag = 1 << 13 // IPS_HELPER
	StatusOffload      StatusFlag = 1 << 14 // IPS_OFFLOAD
)
