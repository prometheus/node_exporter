package netfilter

import "github.com/mdlayher/netlink"

// UnmarshalNetlink unmarshals a netlink.Message into a Netfilter Header and Attributes.
func UnmarshalNetlink(msg netlink.Message) (Header, []Attribute, error) {

	var h Header

	err := h.unmarshal(msg)
	if err != nil {
		return Header{}, nil, err
	}

	attrs, err := unmarshalAttributes(msg.Data[nfHeaderLen:])
	if err != nil {
		return Header{}, nil, err
	}

	return h, attrs, nil
}

// MarshalNetlink takes a Netfilter Header and Attributes and returns a netlink.Message.
func MarshalNetlink(h Header, attrs []Attribute) (netlink.Message, error) {

	ba, err := marshalAttributes(attrs)
	if err != nil {
		return netlink.Message{}, err
	}

	// initialize with 4 bytes of Data before unmarshal
	nlm := netlink.Message{Data: make([]byte, 4)}

	// marshal error ignored, safe to do if msg Data is initialized
	h.marshal(&nlm)

	nlm.Data = append(nlm.Data, ba...)

	return nlm, nil
}
