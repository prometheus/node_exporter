package netfilter

import (
	"encoding/binary"
	"fmt"

	"github.com/mdlayher/netlink"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

// An Attribute is a copy of a netlink.Attribute that can be nested.
type Attribute struct {

	// The type of this Attribute, typically matched to a constant.
	Type uint16

	// An arbitrary payload which is specified by Type.
	Data []byte

	// Whether the attribute's data contains nested attributes.
	Nested   bool
	Children []Attribute

	// Whether the attribute's data is in network (true) or native (false) byte order.
	NetByteOrder bool
}

func (a Attribute) String() string {
	if a.Nested {
		return fmt.Sprintf("<Length %d, Type %d, Nested %t, %d Children (%v)>", len(a.Data), a.Type, a.Nested, len(a.Children), a.Children)
	}

	return fmt.Sprintf("<Length %d, Type %d, Nested %t, NetByteOrder %t, %v>", len(a.Data), a.Type, a.Nested, a.NetByteOrder, a.Data)

}

// Uint16 interprets a non-nested Netfilter attribute in network byte order as a uint16.
func (a Attribute) Uint16() uint16 {

	if a.Nested {
		panic("Uint16: unexpected Nested attribute")
	}

	if l := len(a.Data); l != 2 {
		panic(fmt.Sprintf("Uint16: unexpected byte slice length: %d", l))
	}

	return binary.BigEndian.Uint16(a.Data)
}

// PutUint16 sets the Attribute's data field to a Uint16 encoded in net byte order.
func (a *Attribute) PutUint16(v uint16) {

	if len(a.Data) != 2 {
		a.Data = make([]byte, 2)
	}

	binary.BigEndian.PutUint16(a.Data, v)
}

// Uint32 interprets a non-nested Netfilter attribute in network byte order as a uint32.
func (a Attribute) Uint32() uint32 {

	if a.Nested {
		panic("Uint32: unexpected Nested attribute")
	}

	if l := len(a.Data); l != 4 {
		panic(fmt.Sprintf("Uint32: unexpected byte slice length: %d", l))
	}

	return binary.BigEndian.Uint32(a.Data)
}

// PutUint32 sets the Attribute's data field to a Uint32 encoded in net byte order.
func (a *Attribute) PutUint32(v uint32) {

	if len(a.Data) != 4 {
		a.Data = make([]byte, 4)
	}

	binary.BigEndian.PutUint32(a.Data, v)
}

// Int32 converts the result of Uint16() to an int32.
func (a Attribute) Int32() int32 {
	return int32(a.Uint32())
}

// Uint64 interprets a non-nested Netfilter attribute in network byte order as a uint64.
func (a Attribute) Uint64() uint64 {

	if a.Nested {
		panic("Uint64: unexpected Nested attribute")
	}

	if l := len(a.Data); l != 8 {
		panic(fmt.Sprintf("Uint64: unexpected byte slice length: %d", l))
	}

	return binary.BigEndian.Uint64(a.Data)
}

// PutUint64 sets the Attribute's data field to a Uint64 encoded in net byte order.
func (a *Attribute) PutUint64(v uint64) {

	if len(a.Data) != 8 {
		a.Data = make([]byte, 8)
	}

	binary.BigEndian.PutUint64(a.Data, v)
}

// Int64 converts the result of Uint16() to an int64.
func (a Attribute) Int64() int64 {
	return int64(a.Uint64())
}

// Uint16Bytes gets the big-endian 2-byte representation of a uint16.
func Uint16Bytes(u uint16) []byte {
	d := make([]byte, 2)
	binary.BigEndian.PutUint16(d, u)
	return d
}

// Uint32Bytes gets the big-endian 4-byte representation of a uint32.
func Uint32Bytes(u uint32) []byte {
	d := make([]byte, 4)
	binary.BigEndian.PutUint32(d, u)
	return d
}

// Uint64Bytes gets the big-endian 8-byte representation of a uint64.
func Uint64Bytes(u uint64) []byte {
	d := make([]byte, 8)
	binary.BigEndian.PutUint64(d, u)
	return d
}

// unmarshalAttributes returns an array of netfilter.Attributes decoded from
// a byte array. This byte array should be taken from the netlink.Message's
// Data payload after the nfHeaderLen offset.
func unmarshalAttributes(b []byte) ([]Attribute, error) {

	// Obtain a list of parsed netlink attributes possibly holding
	// nested Netfilter attributes in their binary Data field.
	attrs, err := netlink.UnmarshalAttributes(b)
	if err != nil {
		return nil, errors.Wrap(err, errWrapNetlinkUnmarshalAttrs)
	}

	var ra []Attribute

	// Only allocate backing array when there are netlink attributes to decode.
	if len(attrs) != 0 {
		ra = make([]Attribute, 0, len(attrs))
	}

	// Wrap all netlink.Attributes into netfilter.Attributes to support nesting
	for _, nla := range attrs {

		// Copy the netlink attribute's fields into the netfilter attribute.
		nfa := Attribute{
			// Only consider the rightmost 14 bits for Type
			Type: nla.Type & ^(uint16(unix.NLA_F_NESTED) | uint16(unix.NLA_F_NET_BYTEORDER)),
			Data: nla.Data,
		}

		// Boolean flags extracted from the two leftmost bits of Type
		nfa.Nested = (nla.Type & uint16(unix.NLA_F_NESTED)) != 0
		nfa.NetByteOrder = (nla.Type & uint16(unix.NLA_F_NET_BYTEORDER)) != 0

		if nfa.NetByteOrder && nfa.Nested {
			return nil, errInvalidAttributeFlags
		}

		// Unmarshal recursively if the netlink Nested flag is set
		if nfa.Nested {
			if nfa.Children, err = unmarshalAttributes(nla.Data); err != nil {
				return nil, err
			}
		}

		ra = append(ra, nfa)
	}

	return ra, nil
}

// marshalAttributes marshals a nested attribute structure into a byte slice.
// This byte slice can then be copied into a netlink.Message's Data field after
// the nfHeaderLen offset.
func marshalAttributes(attrs []Attribute) ([]byte, error) {

	// netlink.Attribute to use as scratch buffer, requires a single allocation
	nla := netlink.Attribute{}

	// Output array, initialized to the length of the input array
	ra := make([]netlink.Attribute, 0, len(attrs))

	for _, nfa := range attrs {

		if nfa.NetByteOrder && nfa.Nested {
			return nil, errInvalidAttributeFlags
		}

		// Save nested or byte order flags back to the netlink.Attribute's
		// Type field to include it in the marshaling operation
		nla.Type = nfa.Type

		switch {
		case nfa.Nested:
			nla.Type = nla.Type | unix.NLA_F_NESTED
		case nfa.NetByteOrder:
			nla.Type = nla.Type | unix.NLA_F_NET_BYTEORDER
		}

		// Recursively marshal the attribute's children
		if nfa.Nested {
			nfnab, err := marshalAttributes(nfa.Children)
			if err != nil {
				return nil, err
			}

			nla.Data = nfnab
		} else {
			nla.Data = nfa.Data
		}

		ra = append(ra, nla)
	}

	// Marshal all Netfilter attributes into binary representation of Netlink attributes
	return netlink.MarshalAttributes(ra)
}
