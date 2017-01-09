package netlink

import (
	"errors"

	"github.com/mdlayher/netlink/nlenc"
)

var (
	// errInvalidAttribute specifies if an Attribute's length is incorrect.
	errInvalidAttribute = errors.New("invalid attribute; length too short or too large")
)

// An Attribute is a netlink attribute.  Attributes are packed and unpacked
// to and from the Data field of Message for some netlink protocol families.
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
	copy(b[4:], a.Data)

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
	case a.Length < 4:
		return errInvalidAttribute
	// Data present
	case a.Length >= 4:
		a.Data = make([]byte, len(b[4:a.Length]))
		copy(a.Data, b[4:a.Length])
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
