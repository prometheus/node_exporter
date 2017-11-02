//+build linux

package genetlink

import (
	"errors"
	"fmt"
	"math"

	"github.com/mdlayher/netlink"
	"github.com/mdlayher/netlink/nlenc"
	"golang.org/x/sys/unix"
)

var (
	// errInvalidFamilyVersion is returned when a family's version is greater
	// than an 8-bit integer.
	errInvalidFamilyVersion = errors.New("invalid family version attribute")

	// errInvalidMulticastGroupArray is returned when a multicast group array
	// of attributes is malformed.
	errInvalidMulticastGroupArray = errors.New("invalid multicast group attribute array")
)

// getFamily retrieves a generic netlink family with the specified name.
func (c *Conn) getFamily(name string) (Family, error) {
	b, err := netlink.MarshalAttributes([]netlink.Attribute{{
		Type: unix.CTRL_ATTR_FAMILY_NAME,
		Data: nlenc.Bytes(name),
	}})
	if err != nil {
		return Family{}, err
	}

	req := Message{
		Header: Header{
			Command: unix.CTRL_CMD_GETFAMILY,
			// TODO(mdlayher): grab nlctrl version?
			Version: 1,
		},
		Data: b,
	}

	msgs, err := c.Execute(req, unix.GENL_ID_CTRL, netlink.HeaderFlagsRequest)
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

// listFamilies retrieves all registered generic netlink families.
func (c *Conn) listFamilies() ([]Family, error) {
	req := Message{
		Header: Header{
			Command: unix.CTRL_CMD_GETFAMILY,
			// TODO(mdlayher): grab nlctrl version?
			Version: 1,
		},
	}

	flags := netlink.HeaderFlagsRequest | netlink.HeaderFlagsDump
	msgs, err := c.Execute(req, unix.GENL_ID_CTRL, flags)
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

// parseAttributes parses netlink attributes into a Family's fields.
func (f *Family) parseAttributes(attrs []netlink.Attribute) error {
	for _, a := range attrs {
		switch a.Type {
		case unix.CTRL_ATTR_FAMILY_ID:
			f.ID = nlenc.Uint16(a.Data)
		case unix.CTRL_ATTR_FAMILY_NAME:
			f.Name = nlenc.String(a.Data)
		case unix.CTRL_ATTR_VERSION:
			v := nlenc.Uint32(a.Data)
			if v > math.MaxUint8 {
				return errInvalidFamilyVersion
			}

			f.Version = uint8(v)
		case unix.CTRL_ATTR_MCAST_GROUPS:
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
			case unix.CTRL_ATTR_MCAST_GRP_NAME:
				g.Name = nlenc.String(na.Data)
			case unix.CTRL_ATTR_MCAST_GRP_ID:
				g.ID = nlenc.Uint32(na.Data)
			}
		}

		groups = append(groups, g)
	}

	return groups, nil
}
