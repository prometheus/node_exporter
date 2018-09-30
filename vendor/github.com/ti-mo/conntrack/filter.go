package conntrack

import (
	"github.com/ti-mo/netfilter"
)

// Filter is a structure used in dump operations to filter the response
// based on a given connmark and mask. The mask is applied to the Mark field of
// all flows in the conntrack table, the result is compared to the filter's Mark.
// Each flow that matches will be returned by the kernel.
type Filter struct {
	Mark, Mask uint32
}

// marshal marshals a Filter into a list of netfilter.Attributes.
func (f Filter) marshal() []netfilter.Attribute {

	return []netfilter.Attribute{
		{
			Type: uint16(ctaMark),
			Data: netfilter.Uint32Bytes(f.Mark),
		},
		{
			Type: uint16(ctaMarkMask),
			Data: netfilter.Uint32Bytes(f.Mask),
		},
	}
}
