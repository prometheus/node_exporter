package conntrack

import (
	"fmt"

	"github.com/mdlayher/netlink"
	"github.com/ti-mo/netfilter"
)

// Event holds information about a Conntrack event.
type Event struct {
	Type eventType

	Flow   *Flow
	Expect *Expect
}

// eventType is a custom type that describes the Conntrack event type.
type eventType uint8

// List of all types of Conntrack events. This is an internal representation
// unrelated to any message types in the kernel source.
const (
	EventUnknown eventType = iota
	EventNew
	EventUpdate
	EventDestroy
	EventExpNew
	EventExpDestroy
)

// unmarshal unmarshals a Conntrack EventType from a Netfilter header.
func (et *eventType) unmarshal(h netfilter.Header) error {

	// Fail when the message is not a conntrack message
	if h.SubsystemID == netfilter.NFSubsysCTNetlink {
		switch messageType(h.MessageType) {
		case ctNew:
			// Since the MessageType is only of kind new, get or delete,
			// the header's flags are used to distinguish between NEW and UPDATE.
			if h.Flags&(netlink.HeaderFlagsCreate|netlink.HeaderFlagsExcl) != 0 {
				*et = EventNew
			} else {
				*et = EventUpdate
			}
		case ctDelete:
			*et = EventDestroy
		default:
			return fmt.Errorf(errUnknownEventType, h.MessageType)
		}
	} else if h.SubsystemID == netfilter.NFSubsysCTNetlinkExp {
		switch expMessageType(h.MessageType) {
		case ctExpNew:
			*et = EventExpNew
		case ctExpDelete:
			*et = EventExpDestroy
		default:
			return fmt.Errorf(errUnknownEventType, h.MessageType)
		}
	} else {
		return errNotConntrack
	}

	return nil
}

// unmarshal unmarshals a Netlink message into an Event structure.
func (e *Event) unmarshal(nlmsg netlink.Message) error {

	// Make sure we don't re-use an Event structure
	if e.Expect != nil || e.Flow != nil {
		return errReusedEvent
	}

	var err error

	// Unmarshal a netlink.Message into netfilter.Attributes and Header
	h, attrs, err := netfilter.UnmarshalNetlink(nlmsg)
	if err != nil {
		return err
	}

	// Decode the header to make sure we're dealing with a Conntrack event
	err = e.Type.unmarshal(h)
	if err != nil {
		return err
	}

	// Unmarshal Netfilter attributes into the event's Flow or Expect entry
	if h.SubsystemID == netfilter.NFSubsysCTNetlink {
		e.Flow = new(Flow)
		err = e.Flow.unmarshal(attrs)
	} else if h.SubsystemID == netfilter.NFSubsysCTNetlinkExp {
		e.Expect = new(Expect)
		err = e.Expect.unmarshal(attrs)
	}

	if err != nil {
		return err
	}

	return nil
}
