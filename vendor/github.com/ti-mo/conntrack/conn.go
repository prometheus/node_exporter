package conntrack

import (
	"fmt"

	"github.com/mdlayher/netlink"
	"github.com/pkg/errors"
	"github.com/ti-mo/netfilter"
)

// Conn represents a Netlink connection to the Netfilter
// subsystem and implements all Conntrack actions.
type Conn struct {
	conn *netfilter.Conn
}

// Dial opens a new Netfilter Netlink connection and returns it
// wrapped in a Conn structure that implements the Conntrack API.
func Dial(config *netlink.Config) (*Conn, error) {
	c, err := netfilter.Dial(config)
	if err != nil {
		return nil, err
	}

	return &Conn{c}, nil
}

// Close closes a Conn.
func (c *Conn) Close() error {
	return c.conn.Close()
}

// Listen joins the Netfilter connection to a multicast group and starts a given
// amount of Flow decoders from the Conn to the Flow channel. Returns an error channel
// the workers will return any errors on. Any errors during Flow decoding are fatal.
//
// The Conn will be marked as having listeners active, which will prevent Listen from being
// called again. For listening on other groups, open another socket.
func (c *Conn) Listen(evChan chan<- Event, numWorkers uint8, groups []netfilter.NetlinkGroup) (chan error, error) {

	if numWorkers < 1 {
		return nil, errors.Errorf(errWorkerCount, numWorkers)
	}

	// Prevent Listen() from being called twice on the same Conn.
	// This is checked again in JoinGroups(), but an early failure is preferred.
	if c.conn.IsMulticast() {
		return nil, errConnHasListeners
	}

	err := c.conn.JoinGroups(groups)
	if err != nil {
		return nil, err
	}

	errChan := make(chan error)

	// Start numWorkers amount of worker goroutines
	for id := uint8(0); id < numWorkers; id++ {
		go c.eventWorker(id, evChan, errChan)
	}

	return errChan, nil
}

// eventWorker is a worker function that decodes Netlink messages into Events.
func (c *Conn) eventWorker(workerID uint8, evChan chan<- Event, errChan chan<- error) {

	var err error
	var recv []netlink.Message
	var ev Event

	for {
		// Receive data from the Netlink socket
		recv, err = c.conn.Receive()
		if err != nil {
			errChan <- errors.Wrap(err, fmt.Sprintf(errWorkerReceive, workerID))
			return
		}

		// Receive() always returns a list of Netlink Messages, but multicast messages should never be multi-part
		if len(recv) > 1 {
			errChan <- errMultipartEvent
			return
		}

		// Decode event and send on channel
		ev = *new(Event)
		err := ev.unmarshal(recv[0])
		if err != nil {
			errChan <- err
			return
		}

		evChan <- ev
	}
}

// Dump gets all Conntrack connections from the kernel in the form of a list
// of Flow objects.
func (c *Conn) Dump() ([]Flow, error) {

	req, err := netfilter.MarshalNetlink(
		netfilter.Header{
			SubsystemID: netfilter.NFSubsysCTNetlink,
			MessageType: netfilter.MessageType(ctGet),
			Family:      netfilter.ProtoUnspec, // ProtoUnspec dumps both IPv4 and IPv6
			Flags:       netlink.HeaderFlagsRequest | netlink.HeaderFlagsDump,
		},
		nil)

	if err != nil {
		return nil, err
	}

	nlm, err := c.conn.Query(req)
	if err != nil {
		return nil, err
	}

	return unmarshalFlows(nlm)
}

// DumpFilter gets all Conntrack connections from the kernel in the form of a list
// of Flow objects, but only returns Flows matching the connmark specified in the Filter parameter.
func (c *Conn) DumpFilter(f Filter) ([]Flow, error) {

	req, err := netfilter.MarshalNetlink(
		netfilter.Header{
			SubsystemID: netfilter.NFSubsysCTNetlink,
			MessageType: netfilter.MessageType(ctGet),
			Family:      netfilter.ProtoUnspec, // ProtoUnspec dumps both IPv4 and IPv6
			Flags:       netlink.HeaderFlagsRequest | netlink.HeaderFlagsDump,
		},
		f.marshal())

	if err != nil {
		return nil, err
	}

	nlm, err := c.conn.Query(req)
	if err != nil {
		return nil, err
	}

	return unmarshalFlows(nlm)
}

// DumpExpect gets all expected Conntrack expectations from the kernel in the form
// of a list of Expect objects.
func (c *Conn) DumpExpect() ([]Expect, error) {

	req, err := netfilter.MarshalNetlink(
		netfilter.Header{
			SubsystemID: netfilter.NFSubsysCTNetlinkExp,
			MessageType: netfilter.MessageType(ctGet),
			Family:      netfilter.ProtoUnspec, // ProtoUnspec dumps both IPv4 and IPv6
			Flags:       netlink.HeaderFlagsRequest | netlink.HeaderFlagsDump | netlink.HeaderFlagsAcknowledge,
		},
		nil)

	if err != nil {
		return nil, err
	}

	nlm, err := c.conn.Query(req)
	if err != nil {
		return nil, err
	}

	return unmarshalExpects(nlm)
}

// Flush empties the Conntrack table. Deletes all IPv4 and IPv6 entries.
func (c *Conn) Flush() error {

	req, err := netfilter.MarshalNetlink(
		netfilter.Header{
			SubsystemID: netfilter.NFSubsysCTNetlink,
			MessageType: netfilter.MessageType(ctDelete),
			Family:      netfilter.ProtoUnspec, // Family is ignored for flush
			Flags:       netlink.HeaderFlagsRequest | netlink.HeaderFlagsAcknowledge,
		},
		nil)

	if err != nil {
		return err
	}

	_, err = c.conn.Query(req)
	if err != nil {
		return err
	}

	return nil
}

// FlushFilter deletes all entries from the Conntrack table matching a given Filter.
// Both IPv4 and IPv6 entries are considered for deletion.
func (c *Conn) FlushFilter(f Filter) error {

	req, err := netfilter.MarshalNetlink(
		netfilter.Header{
			SubsystemID: netfilter.NFSubsysCTNetlink,
			MessageType: netfilter.MessageType(ctDelete),
			Family:      netfilter.ProtoUnspec, // Family is ignored for flush
			Flags:       netlink.HeaderFlagsRequest | netlink.HeaderFlagsAcknowledge,
		},
		f.marshal())

	if err != nil {
		return err
	}

	_, err = c.conn.Query(req)
	if err != nil {
		return err
	}

	return nil
}

// Create creates a new Conntrack entry.
func (c *Conn) Create(f Flow) error {

	// Conntrack create requires timeout to be set.
	if f.Timeout == 0 {
		return errNeedTimeout
	}

	attrs, err := f.marshal()
	if err != nil {
		return err
	}

	pf := netfilter.ProtoIPv4
	if f.TupleOrig.IP.IsIPv6() && f.TupleReply.IP.IsIPv6() {
		pf = netfilter.ProtoIPv6
	}

	req, err := netfilter.MarshalNetlink(
		netfilter.Header{
			SubsystemID: netfilter.NFSubsysCTNetlink,
			MessageType: netfilter.MessageType(ctNew),
			Family:      pf,
			Flags: netlink.HeaderFlagsRequest | netlink.HeaderFlagsAcknowledge |
				netlink.HeaderFlagsExcl | netlink.HeaderFlagsCreate,
		}, attrs)

	if err != nil {
		return err
	}

	_, err = c.conn.Query(req)
	if err != nil {
		return err
	}

	return nil
}

// CreateExpect creates a new Conntrack Expect entry. Warning: Experimental, haven't
// got this to create an Expect correctly. Best-effort implementation based on kernel source.
func (c *Conn) CreateExpect(ex Expect) error {

	attrs, err := ex.marshal()
	if err != nil {
		return err
	}

	pf := netfilter.ProtoIPv4
	if ex.Tuple.IP.IsIPv6() && ex.Mask.IP.IsIPv6() && ex.TupleMaster.IP.IsIPv6() {
		pf = netfilter.ProtoIPv6
	}

	req, err := netfilter.MarshalNetlink(
		netfilter.Header{
			SubsystemID: netfilter.NFSubsysCTNetlinkExp,
			MessageType: netfilter.MessageType(ctExpNew),
			Family:      pf,
			Flags: netlink.HeaderFlagsRequest | netlink.HeaderFlagsAcknowledge |
				netlink.HeaderFlagsExcl | netlink.HeaderFlagsCreate,
		}, attrs)

	if err != nil {
		return err
	}

	_, err = c.conn.Query(req)
	if err != nil {
		return err
	}

	return nil
}

// Get queries the conntrack table for a connection matching some attributes of a given Flow.
// The following attributes are considered in the query: TupleOrig or TupleReply, in that order,
// and Zone. One of TupleOrig or TupleReply is required for a successful query.
func (c *Conn) Get(f Flow) (Flow, error) {

	var qf Flow

	attrs, err := f.marshal()
	if err != nil {
		return qf, err
	}

	pf := netfilter.ProtoIPv4
	if f.TupleOrig.IP.IsIPv6() && f.TupleReply.IP.IsIPv6() {
		pf = netfilter.ProtoIPv6
	}

	req, err := netfilter.MarshalNetlink(
		netfilter.Header{
			SubsystemID: netfilter.NFSubsysCTNetlink,
			MessageType: netfilter.MessageType(ctGet),
			Family:      pf,
			Flags:       netlink.HeaderFlagsRequest | netlink.HeaderFlagsAcknowledge,
		}, attrs)

	if err != nil {
		return qf, err
	}

	nlm, err := c.conn.Query(req)
	if err != nil {
		return qf, err
	}

	// Since this is not a dump (and ACK flag is set), the kernel sends a message containing
	// the flow, followed by a Netlink (non-)error message. The error is already parsed by
	// the netlink library, so we only read the first message containing the Flow.
	qf, err = unmarshalFlow(nlm[0])
	if err != nil {
		return qf, err
	}

	return qf, nil
}

// Update updates a Conntrack entry. Only the following attributes are considered
// when sending a Flow update: Helper, Timeout, Status, ProtoInfo, Mark, SeqAdj (orig/reply),
// SynProxy, Labels. All other attributes are immutable past the point of creation.
// See the ctnetlink_change_conntrack() kernel function for exact behaviour.
func (c *Conn) Update(f Flow) error {

	// Kernel rejects updates with a master tuple set
	if f.TupleMaster.filled() {
		return errUpdateMaster
	}

	attrs, err := f.marshal()
	if err != nil {
		return err
	}

	pf := netfilter.ProtoIPv4
	if f.TupleOrig.IP.IsIPv6() && f.TupleReply.IP.IsIPv6() {
		pf = netfilter.ProtoIPv6
	}

	req, err := netfilter.MarshalNetlink(
		netfilter.Header{
			SubsystemID: netfilter.NFSubsysCTNetlink,
			MessageType: netfilter.MessageType(ctNew),
			Family:      pf,
			Flags:       netlink.HeaderFlagsRequest | netlink.HeaderFlagsAcknowledge,
		}, attrs)

	if err != nil {
		return err
	}

	_, err = c.conn.Query(req)
	if err != nil {
		return err
	}

	return nil
}

// Delete removes a Conntrack entry given a Flow. Flows are looked up in the conntrack table
// based on the original and reply tuple. When the Flow's ID field is filled, it must match the
// ID on the connection returned from the tuple lookup, or the delete will fail.
func (c *Conn) Delete(f Flow) error {

	attrs, err := f.marshal()
	if err != nil {
		return err
	}

	// Default to IPv4, set netlink protocol family to IPv6 if orig/reply is IPv6.
	pf := netfilter.ProtoIPv4
	if f.TupleOrig.IP.IsIPv6() && f.TupleReply.IP.IsIPv6() {
		pf = netfilter.ProtoIPv6
	}

	req, err := netfilter.MarshalNetlink(
		netfilter.Header{
			SubsystemID: netfilter.NFSubsysCTNetlink,
			MessageType: netfilter.MessageType(ctDelete),
			Family:      pf,
			Flags:       netlink.HeaderFlagsRequest | netlink.HeaderFlagsAcknowledge,
		}, attrs)

	if err != nil {
		return err
	}

	_, err = c.conn.Query(req)
	if err != nil {
		return err
	}

	return nil
}

// Stats returns a list of Stats structures, one per CPU present in the machine.
// Each Stats structure contains performance counters of all Conntrack actions
// performed on that specific CPU.
func (c *Conn) Stats() ([]Stats, error) {

	req, err := netfilter.MarshalNetlink(
		netfilter.Header{
			SubsystemID: netfilter.NFSubsysCTNetlink,
			MessageType: netfilter.MessageType(ctGetStatsCPU),
			Family:      netfilter.ProtoUnspec,
			Flags:       netlink.HeaderFlagsRequest | netlink.HeaderFlagsDump,
		}, nil)

	if err != nil {
		return nil, err
	}

	msgs, err := c.conn.Query(req)
	if err != nil {
		return nil, err
	}

	return unmarshalStats(msgs)
}

// StatsExpect returns a list of StatsExpect structures, one per CPU present in the machine.
// Each StatsExpect structure indicates how many Expect entries were initialized,
// created or deleted on each CPU.
func (c *Conn) StatsExpect() ([]StatsExpect, error) {

	req, err := netfilter.MarshalNetlink(
		netfilter.Header{
			SubsystemID: netfilter.NFSubsysCTNetlinkExp,
			MessageType: netfilter.MessageType(ctExpGetStatsCPU),
			Family:      netfilter.ProtoUnspec,
			Flags:       netlink.HeaderFlagsRequest | netlink.HeaderFlagsDump,
		}, nil)

	if err != nil {
		return nil, err
	}

	msgs, err := c.conn.Query(req)
	if err != nil {
		return nil, err
	}

	return unmarshalStatsExpect(msgs)
}

// Count queries Conntrack for an internal global counter that describes the total amount
// of Flow entries currently in the Conntrack table. Only the main Conntrack table has this
// fast query available. To get the amount of Expect entries, execute DumpExpect() and count
// the amount of entries returned.
func (c *Conn) Count() (uint32, error) {

	req, err := netfilter.MarshalNetlink(
		netfilter.Header{
			SubsystemID: netfilter.NFSubsysCTNetlink,
			MessageType: netfilter.MessageType(ctGetStats),
			Family:      netfilter.ProtoUnspec,
			Flags:       netlink.HeaderFlagsRequest | netlink.HeaderFlagsDump | netlink.HeaderFlagsAcknowledge,
		}, nil)

	if err != nil {
		return 0, err
	}

	msgs, err := c.conn.Query(req)
	if err != nil {
		return 0, err
	}

	return unmarshalStatsGlobal(msgs[0])
}
