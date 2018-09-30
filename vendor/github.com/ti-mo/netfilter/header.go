package netfilter

import (
	"encoding/binary"
	"fmt"

	"github.com/mdlayher/netlink"
)

// SubsystemID denotes the Netfilter Subsystem ID the message is for. It is a const that
// is defined in the kernel at uapi/linux/netfilter/nfnetlink.h.
type SubsystemID uint8

// MessageType denotes the message type specific to the subsystem. Its meaning can only
// be determined after decoding the Netfilter Subsystem type, because it only has meaning
// in that context. Possible values and string representations need to be implemented in
// a subsystem-specific package.
type MessageType uint8

// String representation of the netfilter Header/
func (h Header) String() string {
	return fmt.Sprintf("<Subsystem: %s, Message Type: %d, Family: %s, Version: %d, ResourceID: %d>",
		h.SubsystemID, h.MessageType, h.Family, h.Version, h.ResourceID)
}

// Header is an abstraction over the Netlink header's Type field and the Netfilter message header,
// also known as 'nfgenmsg'.
//
// The Netlink header's Type field is divided into two bytes by netfilter: the most significant byte
// is the subsystem ID and the least significant is the message type. The significance of the MessageType
// field fully depends on the subsystem the message is for (eg. conntrack). This package is only responsible
// for splitting the field and providing a list of known SubsystemIDs. Subpackages use the MessageType field
// to implement subsystem-specific behaviour.
//
// nfgenmsg holds the protocol family, version and resource ID of the Netfilter message.
// Family describes a protocol family that can be managed using Netfilter (eg. IPv4/6, ARP, Bridge)
// Version is a protocol version descriptor, and always set to 0 (NFNETLINK_V0)
// ResourceID is a generic field specific to the upper layer protocol (eg. CPU ID of Conntrack stats)
type Header struct {

	// Netlink header flags, to (un)marshal to a netlink Message in a single operation
	Flags netlink.HeaderFlags

	// netlink Header Type
	SubsystemID SubsystemID
	MessageType MessageType

	// nfgenmsg
	Family     ProtoFamily
	Version    uint8 // Usually NFNETLINK_V0 (Go: NFNLv0)
	ResourceID uint16
}

// Size of a Netfilter header (nfgenmsg - 4 bytes)
const nfHeaderLen = 4

// unmarshal unmarshals a netlink.Message into a Header. The message Data must be at least 4 bytes long.
// The first 4 bytes of the message's Data field and the message's Header Type/Flags are used.
func (h *Header) unmarshal(nlm netlink.Message) error {

	if len(nlm.Data) < nfHeaderLen {
		return errMessageLen
	}

	h.Flags = nlm.Header.Flags

	h.SubsystemID = SubsystemID(uint16(nlm.Header.Type) & 0xff00 >> 8)
	h.MessageType = MessageType(uint16(nlm.Header.Type) & 0x00ff)

	h.Family = ProtoFamily(nlm.Data[0])
	h.Version = nlm.Data[1]
	h.ResourceID = binary.BigEndian.Uint16(nlm.Data[2:4])

	return nil
}

// marshal marshals a Header into a netlink.Message. The message Data must be initialized with at least 4 bytes.
// The Header Type and the first 4 bytes of the Data field are overwritten.
func (h Header) marshal(nlm *netlink.Message) error {

	if len(nlm.Data) < nfHeaderLen {
		return errMessageLen
	}

	nlm.Header.Flags = h.Flags

	nlm.Header.Type = netlink.HeaderType(uint16(h.SubsystemID)<<8 | uint16(h.MessageType))

	nlm.Data[0] = uint8(h.Family)
	nlm.Data[1] = h.Version
	copy(nlm.Data[2:4], Uint16Bytes(h.ResourceID))

	return nil
}
