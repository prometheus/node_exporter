package conntrack

import (
	"fmt"
	"strconv"
)

// protoLookup translates a protocol integer into its string representation.
func protoLookup(p uint8) string {
	protos := map[uint8]string{
		1:   "icmp",
		2:   "igmp",
		6:   "tcp",
		17:  "udp",
		33:  "dccp",
		47:  "gre",
		58:  "ipv6-icmp",
		94:  "ipip",
		115: "l2tp",
		132: "sctp",
		136: "udplite",
	}

	if val, ok := protos[p]; ok {
		return val
	}

	return strconv.FormatUint(uint64(p), 10)
}

func (s Status) String() string {
	names := []string{
		"EXPECTED",
		"SEEN_REPLY",
		"ASSURED",
		"CONFIRMED",
		"SRC_NAT",
		"DST_NAT",
		"SEQ_ADJUST",
		"SRC_NAT_DONE",
		"DST_NAT_DONE",
		"DYING",
		"FIXED_TIMEOUT",
		"TEMPLATE",
		"UNTRACKED",
		"HELPER",
		"OFFLOAD",
	}

	var rs string

	// Loop over the field's bits
	for i, name := range names {
		if s.Value&(1<<uint32(i)) != 0 {
			if rs != "" {
				rs += "|"
			}
			rs += name
		}
	}

	if rs == "" {
		rs = "NONE"
	}

	return rs
}

func (e Event) String() string {

	if e.Flow != nil {

		// Status flag
		status := ""
		if !e.Flow.Status.SeenReply() {
			status = " (Unreplied)"
		}

		// Accounting information
		acct := "<No Accounting>"
		if e.Flow.CountersOrig.filled() || e.Flow.CountersReply.filled() {
			acct = fmt.Sprintf("Acct: %s %s", e.Flow.CountersOrig, e.Flow.CountersReply)
		}

		// Labels/mask
		labels := "<No Labels>"
		if len(e.Flow.Labels) != 0 && len(e.Flow.LabelsMask) != 0 {
			labels = fmt.Sprintf("Label: <%#x/%#x>", e.Flow.Labels, e.Flow.LabelsMask)
		}

		// Mark/mask
		mark := "<No Mark>"
		if e.Flow.Mark != 0 {
			mark = fmt.Sprintf("Mark: <%#x>", e.Flow.Mark)
		}

		// SeqAdj
		seqadjo := "<No SeqAdjOrig>"
		if e.Flow.SeqAdjOrig.filled() {
			seqadjo = fmt.Sprintf("SeqAdjOrig: %s", e.Flow.SeqAdjOrig)
		}
		seqadjr := "<No SeqAdjReply>"
		if e.Flow.SeqAdjReply.filled() {
			seqadjr = fmt.Sprintf("SeqAdjReply: %s", e.Flow.SeqAdjReply)
		}

		// Security Context
		secctx := "<No SecCtx>"
		if e.Flow.SecurityContext != "" {
			secctx = fmt.Sprintf("SecCtx: %s", e.Flow.SecurityContext)
		}

		return fmt.Sprintf("[%s]%s Timeout: %d, %s, Zone %d, %s, %s, %s, %s, %s, %s",
			e.Type, status,
			e.Flow.Timeout,
			e.Flow.TupleOrig,
			e.Flow.Zone,
			acct, labels, mark,
			seqadjo, seqadjr, secctx)

	} else if e.Expect != nil {

		return fmt.Sprintf("[%s] Timeout: %d, Master: %s, Tuple: %s, Mask: %s, Zone: %d, Helper: '%s', Class: %#x",
			e.Type, e.Expect.Timeout,
			e.Expect.TupleMaster, e.Expect.Tuple, e.Expect.Mask,
			e.Expect.Zone, e.Expect.HelpName, e.Expect.Class,
		)

	} else {
		return fmt.Sprintf("[%s] <Empty Event>", e.Type)
	}

}
