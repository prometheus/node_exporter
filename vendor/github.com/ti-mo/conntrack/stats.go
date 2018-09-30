package conntrack

import (
	"fmt"

	"github.com/mdlayher/netlink"

	"github.com/ti-mo/netfilter"
)

// Stats represents the Conntrack performance counters of a single CPU (core).
// It indicates which and how many Flow operations took place on each CPU.
type Stats struct {
	CPUID         uint16
	Searched      uint32
	Found         uint32
	New           uint32
	Invalid       uint32
	Ignore        uint32
	Delete        uint32
	DeleteList    uint32
	Insert        uint32
	InsertFailed  uint32
	Drop          uint32
	EarlyDrop     uint32
	Error         uint32
	SearchRestart uint32
}

func (s Stats) String() string {
	return fmt.Sprintf(
		"<CPU %d - Searched: %d, Found: %d, New: %d, Invalid: %d, Ignore: %d, Delete: %d, DeleteList: %d, Insert: %d, InsertFailed: %d, Drop: %d, EarlyDrop: %d, Error: %d, SearchRestart: %d>",
		s.CPUID, s.Searched, s.Found, s.New, s.Invalid, s.Ignore, s.Delete, s.DeleteList, s.Insert, s.InsertFailed, s.Drop, s.EarlyDrop, s.Error, s.SearchRestart,
	)
}

// unmarshal unmarshals a list of netfilter.Attributes into a Stats structure.
func (s *Stats) unmarshal(attrs []netfilter.Attribute) error {

	for _, attr := range attrs {
		switch at := cpuStatsType(attr.Type); at {
		case ctaStatsSearched:
			s.Searched = attr.Uint32()
		case ctaStatsFound:
			s.Found = attr.Uint32()
		case ctaStatsNew:
			s.New = attr.Uint32()
		case ctaStatsInvalid:
			s.Invalid = attr.Uint32()
		case ctaStatsIgnore:
			s.Ignore = attr.Uint32()
		case ctaStatsDelete:
			s.Delete = attr.Uint32()
		case ctaStatsDeleteList:
			s.DeleteList = attr.Uint32()
		case ctaStatsInsert:
			s.Insert = attr.Uint32()
		case ctaStatsInsertFailed:
			s.InsertFailed = attr.Uint32()
		case ctaStatsDrop:
			s.Drop = attr.Uint32()
		case ctaStatsEarlyDrop:
			s.EarlyDrop = attr.Uint32()
		case ctaStatsError:
			s.Error = attr.Uint32()
		case ctaStatsSearchRestart:
			s.SearchRestart = attr.Uint32()
		default:
			return fmt.Errorf(errAttributeUnknown, at)
		}
	}

	return nil
}

// StatsExpect represents the Conntrack Expect performance counters of a single CPU (core).
// It indicates how many Expect entries were initialized, created or deleted on each CPU.
type StatsExpect struct {
	CPUID               uint16
	New, Create, Delete uint32
}

// unmarshal unmarshals a list of netfilter.Attributes into a StatsExpect structure.
func (se *StatsExpect) unmarshal(attrs []netfilter.Attribute) error {

	for _, attr := range attrs {
		switch at := expectStatsType(attr.Type); at {
		case ctaStatsExpNew:
			se.New = attr.Uint32()
		case ctaStatsExpCreate:
			se.Create = attr.Uint32()
		case ctaStatsExpDelete:
			se.Delete = attr.Uint32()
		default:
			return fmt.Errorf(errAttributeUnknown, at)
		}
	}

	return nil
}

// unmarshalStats unmarshals a list of Stats from a list of netlink.Messages.
func unmarshalStats(nlm []netlink.Message) ([]Stats, error) {

	stats := make([]Stats, len(nlm))

	for idx, m := range nlm {

		hdr, nfa, err := netfilter.UnmarshalNetlink(m)
		if err != nil {
			return nil, err
		}

		s := Stats{CPUID: hdr.ResourceID}

		err = s.unmarshal(nfa)
		if err != nil {
			return nil, err
		}

		stats[idx] = s
	}

	return stats, nil
}

// unmarshalStatsExpect unmarshals a list of StatsExpect from a list of netlink.Messages.
func unmarshalStatsExpect(nlm []netlink.Message) ([]StatsExpect, error) {

	stats := make([]StatsExpect, len(nlm))

	for idx, m := range nlm {

		hdr, nfa, err := netfilter.UnmarshalNetlink(m)
		if err != nil {
			return nil, err
		}

		se := StatsExpect{CPUID: hdr.ResourceID}

		err = se.unmarshal(nfa)
		if err != nil {
			return nil, err
		}

		stats[idx] = se
	}

	return stats, nil
}

// unmarshalStatsGlobal unmarshals the global Conntrack counter from a netlink.Message.
func unmarshalStatsGlobal(nlm netlink.Message) (uint32, error) {

	_, nfa, err := netfilter.UnmarshalNetlink(nlm)
	if err != nil {
		return 0, err
	}

	// Assert the first (and only) attribute to be a GlobalEntries
	if at := nfa[0].Type; globalStatsType(at) != ctaStatsGlobalEntries {
		return 0, fmt.Errorf(errAttributeUnknown, at)
	}

	return nfa[0].Uint32(), nil
}
