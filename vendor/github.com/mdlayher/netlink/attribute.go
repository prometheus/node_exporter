package netlink

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/mdlayher/netlink/nlenc"
)

var (
	// errInvalidAttribute specifies if an Attribute's length is incorrect.
	errInvalidAttribute = errors.New("invalid attribute; length too short or too large")
)

// An Attribute is a netlink attribute.  Attributes are packed and unpacked
// to and from the Data field of Message for some netlink families.
type Attribute struct {
	// Length of an Attribute, including this field and Type.
	Length uint16

	// The type of this Attribute, typically matched to a constant.
	Type uint16

	// An arbitrary payload which is specified by Type.
	Data []byte
}

// MarshalBinary marshals an Attribute into a byte slice.
func (a Attribute) MarshalBinary() ([]byte, error) {
	if int(a.Length) < nlaHeaderLen {
		return nil, errInvalidAttribute
	}

	b := make([]byte, nlaAlign(int(a.Length)))

	nlenc.PutUint16(b[0:2], a.Length)
	nlenc.PutUint16(b[2:4], a.Type)

	copy(b[nlaHeaderLen:], a.Data)

	return b, nil
}

// UnmarshalBinary unmarshals the contents of a byte slice into an Attribute.
func (a *Attribute) UnmarshalBinary(b []byte) error {
	if len(b) < nlaHeaderLen {
		return errInvalidAttribute
	}

	a.Length = nlenc.Uint16(b[0:2])
	a.Type = nlenc.Uint16(b[2:4])

	if nlaAlign(int(a.Length)) > len(b) {
		return errInvalidAttribute
	}

	switch {
	// No length, no data
	case a.Length == 0:
		a.Data = make([]byte, 0)
	// Not enough length for any data
	case int(a.Length) < nlaHeaderLen:
		return errInvalidAttribute
	// Data present
	case int(a.Length) >= nlaHeaderLen:
		a.Data = make([]byte, len(b[nlaHeaderLen:a.Length]))
		copy(a.Data, b[nlaHeaderLen:a.Length])
	}

	return nil
}

// MarshalAttributes packs a slice of Attributes into a single byte slice.
// In most cases, the Length field of each Attribute should be set to 0, so it
// can be calculated and populated automatically for each Attribute.
func MarshalAttributes(attrs []Attribute) ([]byte, error) {
	var c int
	for _, a := range attrs {
		c += nlaAlign(len(a.Data))
	}

	b := make([]byte, 0, c)
	for _, a := range attrs {
		if a.Length == 0 {
			a.Length = uint16(nlaHeaderLen + len(a.Data))
		}

		ab, err := a.MarshalBinary()
		if err != nil {
			return nil, err
		}

		b = append(b, ab...)
	}

	return b, nil
}

// UnmarshalAttributes unpacks a slice of Attributes from a single byte slice.
//
// It is recommend to use the AttributeDecoder type where possible instead of calling
// UnmarshalAttributes and using package nlenc functions directly.
func UnmarshalAttributes(b []byte) ([]Attribute, error) {
	var attrs []Attribute
	var i int
	for {
		if len(b[i:]) == 0 {
			break
		}

		var a Attribute
		if err := (&a).UnmarshalBinary(b[i:]); err != nil {
			return nil, err
		}

		if a.Length == 0 {
			i += nlaHeaderLen
			continue
		}

		i += nlaAlign(int(a.Length))

		attrs = append(attrs, a)
	}

	return attrs, nil
}

// An AttributeDecoder provides a safe, iterator-like, API around attribute
// decoding.
//
// It is recommend to use an AttributeDecoder where possible instead of calling
// UnmarshalAttributes and using package nlenc functions directly.
//
// The Err method must be called after the Next method returns false to determine
// if any errors occurred during iteration.
type AttributeDecoder struct {
	// ByteOrder defines a specific byte order to use when processing integer
	// attributes.  ByteOrder should be set immediately after creating the
	// AttributeDecoder: before any attributes are parsed.
	//
	// If not set, the native byte order will be used.
	ByteOrder binary.ByteOrder

	// The attributes being worked on, and the iterator index into the slice of
	// attributes.
	attrs []Attribute
	i     int

	// Any error encountered while decoding attributes.
	err error
}

// NewAttributeDecoder creates an AttributeDecoder that unpacks Attributes
// from b and prepares the decoder for iteration.
func NewAttributeDecoder(b []byte) (*AttributeDecoder, error) {
	attrs, err := UnmarshalAttributes(b)
	if err != nil {
		return nil, err
	}

	return &AttributeDecoder{
		// By default, use native byte order.
		ByteOrder: nlenc.NativeEndian(),

		attrs: attrs,
	}, nil
}

// Next advances the decoder to the next netlink attribute.  It returns false
// when no more attributes are present, or an error was encountered.
func (ad *AttributeDecoder) Next() bool {
	if ad.err != nil {
		// Hit an error, stop iteration.
		return false
	}

	ad.i++

	if len(ad.attrs) < ad.i {
		// No more attributes, stop iteration.
		return false
	}

	return true
}

// Type returns the Attribute.Type field of the current netlink attribute
// pointed to by the decoder.
func (ad *AttributeDecoder) Type() uint16 {
	return ad.attr().Type
}

// attr returns the current Attribute pointed to by the decoder.
func (ad *AttributeDecoder) attr() Attribute {
	return ad.attrs[ad.i-1]
}

// data returns the Data field of the current Attribute pointed to by the decoder.
func (ad *AttributeDecoder) data() []byte {
	return ad.attr().Data
}

// Err returns the first error encountered by the decoder.
func (ad *AttributeDecoder) Err() error {
	return ad.err
}

// String returns the string representation of the current Attribute's data.
func (ad *AttributeDecoder) String() string {
	if ad.err != nil {
		return ""
	}

	return nlenc.String(ad.data())
}

// Uint8 returns the uint8 representation of the current Attribute's data.
func (ad *AttributeDecoder) Uint8() uint8 {
	if ad.err != nil {
		return 0
	}

	b := ad.data()
	if len(b) != 1 {
		ad.err = fmt.Errorf("netlink: attribute %d is not a uint8; length: %d", ad.Type(), len(b))
		return 0
	}

	return uint8(b[0])
}

// Uint16 returns the uint16 representation of the current Attribute's data.
func (ad *AttributeDecoder) Uint16() uint16 {
	if ad.err != nil {
		return 0
	}

	b := ad.data()
	if len(b) != 2 {
		ad.err = fmt.Errorf("netlink: attribute %d is not a uint16; length: %d", ad.Type(), len(b))
		return 0
	}

	return ad.ByteOrder.Uint16(b)
}

// Uint32 returns the uint32 representation of the current Attribute's data.
func (ad *AttributeDecoder) Uint32() uint32 {
	if ad.err != nil {
		return 0
	}

	b := ad.data()
	if len(b) != 4 {
		ad.err = fmt.Errorf("netlink: attribute %d is not a uint32; length: %d", ad.Type(), len(b))
		return 0
	}

	return ad.ByteOrder.Uint32(b)
}

// Uint64 returns the uint64 representation of the current Attribute's data.
func (ad *AttributeDecoder) Uint64() uint64 {
	if ad.err != nil {
		return 0
	}

	b := ad.data()
	if len(b) != 8 {
		ad.err = fmt.Errorf("netlink: attribute %d is not a uint64; length: %d", ad.Type(), len(b))
		return 0
	}

	return ad.ByteOrder.Uint64(b)
}

// Do is a general purpose function which allows access to the current data
// pointed to by the AttributeDecoder.
//
// Do can be used to allow parsing arbitrary data within the context of the
// decoder.  Do is most useful when dealing with nested attributes, attribute
// arrays, or decoding arbitrary types (such as C structures) which don't fit
// cleanly into a typical unsigned integer value.
//
// The function fn should not retain any reference to the data b outside of the
// scope of the function.
func (ad *AttributeDecoder) Do(fn func(b []byte) error) {
	if ad.err != nil {
		return
	}

	b := ad.data()
	if err := fn(b); err != nil {
		ad.err = err
	}
}
