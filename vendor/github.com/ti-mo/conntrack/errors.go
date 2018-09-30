package conntrack

import "errors"

var (
	errNotConntrack     = errors.New("trying to decode a non-conntrack or conntrack-exp message")
	errConnHasListeners = errors.New("Conn has existing listeners, open another to listen on more groups")
	errMultipartEvent   = errors.New("received multicast event with more than one Netlink message")

	errNested          = errors.New("unexpected Nested attribute")
	errNotNested       = errors.New("need a Nested attribute to decode this structure")
	errNeedSingleChild = errors.New("need (at least) 1 child attribute")
	errNeedChildren    = errors.New("need (at least) 2 child attributes")
	errIncorrectSize   = errors.New("binary attribute data has incorrect size")

	errReusedEvent     = errors.New("cannot to unmarshal into existing Event")
	errReusedProtoInfo = errors.New("cannot to unmarshal into existing ProtoInfo")

	errBadIPTuple = errors.New("IPTuple source and destination addresses must be valid and belong to the same address family")

	errNeedTimeout = errors.New("Flow needs Timeout field set for this operation")
	errNeedTuples  = errors.New("Flow needs Original and Reply Tuple set for this operation")

	errUpdateMaster = errors.New("cannot send TupleMaster in Flow update")

	errExpectNeedTuples = errors.New("Expect needs Tuple, Mask and TupleMaster Tuples set for this operation")
)

const (
	errUnknownEventType   = "unknown event type %d"
	errWorkerCount        = "invalid worker count %d"
	errWorkerReceive      = "netlink.Receive error in listenWorker %d, exiting"
	errAttributeWrongType = "attribute type '%d' is not a %s"
	errAttributeChild     = "child Type '%d' unknown for attribute type %s"
	errAttributeUnknown   = "attribute type '%d' unknown"
	errExactChildren      = "need exactly %d child attributes for attribute type %s"
)
