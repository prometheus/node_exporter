package genetlink

import (
	"errors"
	"fmt"
	"math"

	"github.com/mdlayher/netlink"
	"github.com/mdlayher/netlink/nlenc"
)

// Constants used to request information from generic netlink controller.
// Reference: http://lxr.free-electrons.com/source/include/linux/genetlink.h?v=3.3#L35
const (
	ctrlVersion = 1

	ctrlCommandGetFamily = 3
)

var (
	// errInvalidFamilyVersion is returned when a family's version is greater
	// than an 8-bit integer.
	errInvalidFamilyVersion = errors.New("invalid family version attribute")

	// errInvalidMulticastGroupArray is returned when a multicast group array
	// of attributes is malformed.
	errInvalidMulticastGroupArray = errors.New("invalid multicast group attribute array")
)

// A Family is a generic netlink family.
type Family struct {
	ID      uint16
	Version uint8
	Name    string
	Groups  []MulticastGroup
}

// A MulticastGroup is a generic netlink multicast group, which can be joined
// for notifications from generic netlink families when specific events take
// place.
type MulticastGroup struct {
	ID   uint32
	Name string
}

// A FamilyService is used to retrieve generic netlink family information.
type FamilyService struct {
	c *Conn
}

// Get retrieves a generic netlink family with the specified name.  If the
// family does not exist, the error value can be checked using os.IsNotExist.
func (s *FamilyService) Get(name string) (Family, error) {
	b, err := netlink.MarshalAttributes([]netlink.Attribute{{
		Type: attrFamilyName,
		Data: nlenc.Bytes(name),
	}})
	if err != nil {
		return Family{}, err
	}

	req := Message{
		Header: Header{
			Command: ctrlCommandGetFamily,
			Version: ctrlVersion,
		},
		Data: b,
	}

	msgs, err := s.c.Execute(req, Controller, netlink.HeaderFlagsRequest)
	if err != nil {
		return Family{}, err
	}

	// TODO(mdlayher): consider interpreting generic netlink header values

	families, err := buildFamilies(msgs)
	if err != nil {
		return Family{}, err
	}
	if len(families) != 1 {
		// If this were to ever happen, netlink must be in a state where
		// its answers cannot be trusted
		panic(fmt.Sprintf("netlink returned multiple families for name: %q", name))
	}

	return families[0], nil
}

// List retrieves all registered generic netlink families.
func (s *FamilyService) List() ([]Family, error) {
	req := Message{
		Header: Header{
			Command: ctrlCommandGetFamily,
			Version: ctrlVersion,
		},
	}

	flags := netlink.HeaderFlagsRequest | netlink.HeaderFlagsDump
	msgs, err := s.c.Execute(req, Controller, flags)
	if err != nil {
		return nil, err
	}

	return buildFamilies(msgs)
}

// buildFamilies builds a slice of Families by parsing attributes from the
// input Messages.
func buildFamilies(msgs []Message) ([]Family, error) {
	families := make([]Family, 0, len(msgs))
	for _, m := range msgs {
		attrs, err := netlink.UnmarshalAttributes(m.Data)
		if err != nil {
			return nil, err
		}

		var f Family
		if err := (&f).parseAttributes(attrs); err != nil {
			return nil, err
		}

		families = append(families, f)
	}

	return families, nil
}

// Attribute IDs mapped to specific family fields.
const (
	// TODO(mdlayher): parse additional attributes

	// Family attributes
	attrUnspecified     = 0
	attrFamilyID        = 1
	attrFamilyName      = 2
	attrVersion         = 3
	attrMulticastGroups = 7

	// Multicast group-specific attributes
	attrMGName = 1
	attrMGID   = 2
)

// parseAttributes parses netlink attributes into a Family's fields.
func (f *Family) parseAttributes(attrs []netlink.Attribute) error {
	for _, a := range attrs {
		switch a.Type {
		case attrFamilyID:
			f.ID = nlenc.Uint16(a.Data)
		case attrFamilyName:
			f.Name = nlenc.String(a.Data)
		case attrVersion:
			v := nlenc.Uint32(a.Data)
			if v > math.MaxUint8 {
				return errInvalidFamilyVersion
			}

			f.Version = uint8(v)
		case attrMulticastGroups:
			groups, err := parseMulticastGroups(a.Data)
			if err != nil {
				return err
			}

			f.Groups = groups
		}
	}

	return nil
}

// parseMulticastGroups parses an array of multicast group nested attributes
// into a slice of MulticastGroups.
func parseMulticastGroups(b []byte) ([]MulticastGroup, error) {
	attrs, err := netlink.UnmarshalAttributes(b)
	if err != nil {
		return nil, err
	}

	groups := make([]MulticastGroup, 0, len(attrs))
	for i, a := range attrs {
		// The type attribute is essentially an array index here; it starts
		// at 1 and should increment for each new array element
		if int(a.Type) != i+1 {
			return nil, errInvalidMulticastGroupArray
		}

		nattrs, err := netlink.UnmarshalAttributes(a.Data)
		if err != nil {
			return nil, err
		}

		var g MulticastGroup
		for _, na := range nattrs {
			switch na.Type {
			case attrMGName:
				g.Name = nlenc.String(na.Data)
			case attrMGID:
				g.ID = nlenc.Uint32(na.Data)
			}
		}

		groups = append(groups, g)
	}

	return groups, nil
}
