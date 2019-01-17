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

// marshal marshals the contents of a into b and returns the number of bytes
// written to b, including attribute alignment padding.
func (a *Attribute) marshal(b []byte) (int, error) {
	if int(a.Length) < nlaHeaderLen {
		return 0, errInvalidAttribute
	}

	nlenc.PutUint16(b[0:2], a.Length)
	nlenc.PutUint16(b[2:4], a.Type)
	n := copy(b[nlaHeaderLen:], a.Data)

	return nlaHeaderLen + nlaAlign(n), nil
}

// unmarshal unmarshals the contents of a byte slice into an Attribute.
func (a *Attribute) unmarshal(b []byte) error {
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
	// Count how many bytes we should allocate to store each attribute's contents.
	var c int
	for _, a := range attrs {
		c += nlaHeaderLen + nlaAlign(len(a.Data))
	}

	// Advance through b with idx to place attribute data at the correct offset.
	var idx int
	b := make([]byte, c)
	for _, a := range attrs {
		// Infer the length of attribute if zero.
		if a.Length == 0 {
			a.Length = uint16(nlaHeaderLen + len(a.Data))
		}

		// Marshal a into b and advance idx to show many bytes are occupied.
		n, err := a.marshal(b[idx:])
		if err != nil {
			return nil, err
		}
		idx += n
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
		if err := (&a).unmarshal(b[i:]); err != nil {
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

// Bytes returns the raw bytes of the current Attribute's data.
func (ad *AttributeDecoder) Bytes() []byte {
	src := ad.data()
	dest := make([]byte, len(src))
	copy(dest, src)
	return dest
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

// An AttributeEncoder provides a safe way to encode attributes.
//
// It is recommended to use an AttributeEncoder where possible instead of
// calling MarshalAttributes or using package nlenc directly.
//
// Errors from intermediate encoding steps are returned in the call to
// Encode.
type AttributeEncoder struct {
	// ByteOrder defines a specific byte order to use when processing integer
	// attributes.  ByteOrder should be set immediately after creating the
	// AttributeEncoder: before any attributes are encoded.
	//
	// If not set, the native byte order will be used.
	ByteOrder binary.ByteOrder

	attrs []Attribute
	err   error
}

// NewAttributeEncoder creates an AttributeEncoder that encodes Attributes.
func NewAttributeEncoder() *AttributeEncoder {
	return &AttributeEncoder{
		ByteOrder: nlenc.NativeEndian(),
	}
}

// Uint8 encodes uint8 data into an Attribute specified by typ.
func (ae *AttributeEncoder) Uint8(typ uint16, v uint8) {
	if ae.err != nil {
		return
	}

	ae.attrs = append(ae.attrs, Attribute{
		Type: typ,
		Data: []byte{v},
	})
}

// Uint16 encodes uint16 data into an Attribute specified by typ.
func (ae *AttributeEncoder) Uint16(typ uint16, v uint16) {
	if ae.err != nil {
		return
	}

	b := make([]byte, 2)
	ae.ByteOrder.PutUint16(b, v)

	ae.attrs = append(ae.attrs, Attribute{
		Type: typ,
		Data: b,
	})
}

// Uint32 encodes uint32 data into an Attribute specified by typ.
func (ae *AttributeEncoder) Uint32(typ uint16, v uint32) {
	if ae.err != nil {
		return
	}

	b := make([]byte, 4)
	ae.ByteOrder.PutUint32(b, v)

	ae.attrs = append(ae.attrs, Attribute{
		Type: typ,
		Data: b,
	})
}

// Uint64 encodes uint64 data into an Attribute specified by typ.
func (ae *AttributeEncoder) Uint64(typ uint16, v uint64) {
	if ae.err != nil {
		return
	}

	b := make([]byte, 8)
	ae.ByteOrder.PutUint64(b, v)

	ae.attrs = append(ae.attrs, Attribute{
		Type: typ,
		Data: b,
	})
}

// String encodes string s as a null-terminated string into an Attribute
// specified by typ.
func (ae *AttributeEncoder) String(typ uint16, s string) {
	if ae.err != nil {
		return
	}

	ae.attrs = append(ae.attrs, Attribute{
		Type: typ,
		Data: nlenc.Bytes(s),
	})
}

// Bytes embeds raw byte data into an Attribute specified by typ.
func (ae *AttributeEncoder) Bytes(typ uint16, b []byte) {
	if ae.err != nil {
		return
	}

	ae.attrs = append(ae.attrs, Attribute{
		Type: typ,
		Data: b,
	})
}

// Do is a general purpose function to encode arbitrary data into an attribute
// specified by typ.
//
// Do is especially helpful in encoding nested attributes, attribute arrays,
// or encoding arbitrary types (such as C structures) which don't fit cleanly
// into an unsigned integer value.
func (ae *AttributeEncoder) Do(typ uint16, fn func() ([]byte, error)) {
	if ae.err != nil {
		return
	}

	b, err := fn()
	if err != nil {
		ae.err = err
		return
	}

	ae.attrs = append(ae.attrs, Attribute{
		Type: typ,
		Data: b,
	})
}

// Encode returns the encoded bytes representing the attributes.
func (ae *AttributeEncoder) Encode() ([]byte, error) {
	if ae.err != nil {
		return nil, ae.err
	}

	return MarshalAttributes(ae.attrs)
}
