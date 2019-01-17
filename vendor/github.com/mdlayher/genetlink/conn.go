package genetlink

import (
	"syscall"

	"github.com/mdlayher/netlink"
	"golang.org/x/net/bpf"
)

// Protocol is the netlink protocol constant used to specify generic netlink.
const Protocol = 0x10 // unix.NETLINK_GENERIC

// A Conn is a generic netlink connection.  A Conn can be used to send and
// receive generic netlink messages to and from netlink.
type Conn struct {
	// Operating system-specific netlink connection.
	c *netlink.Conn
}

// Dial dials a generic netlink connection.  Config specifies optional
// configuration for the underlying netlink connection.  If config is
// nil, a default configuration will be used.
func Dial(config *netlink.Config) (*Conn, error) {
	c, err := netlink.Dial(Protocol, config)
	if err != nil {
		return nil, err
	}

	return NewConn(c), nil
}

// NewConn creates a Conn that wraps an existing *netlink.Conn for
// generic netlink communications.
//
// NewConn is primarily useful for tests. Most applications should use
// Dial instead.
func NewConn(c *netlink.Conn) *Conn {
	return &Conn{c: c}
}

// Close closes the connection.
func (c *Conn) Close() error {
	return c.c.Close()
}

// GetFamily retrieves a generic netlink family with the specified name.  If the
// family does not exist, the error value can be checked using os.IsNotExist.
func (c *Conn) GetFamily(name string) (Family, error) {
	return c.getFamily(name)
}

// ListFamilies retrieves all registered generic netlink families.
func (c *Conn) ListFamilies() ([]Family, error) {
	return c.listFamilies()
}

// JoinGroup joins a netlink multicast group by its ID.
func (c *Conn) JoinGroup(group uint32) error {
	return c.c.JoinGroup(group)
}

// LeaveGroup leaves a netlink multicast group by its ID.
func (c *Conn) LeaveGroup(group uint32) error {
	return c.c.LeaveGroup(group)
}

// SetBPF attaches an assembled BPF program to a Conn.
func (c *Conn) SetBPF(filter []bpf.RawInstruction) error {
	return c.c.SetBPF(filter)
}

// RemoveBPF removes a BPF filter from a Conn.
func (c *Conn) RemoveBPF() error {
	return c.c.RemoveBPF()
}

// SetOption enables or disables a netlink socket option for the Conn.
func (c *Conn) SetOption(option netlink.ConnOption, enable bool) error {
	return c.c.SetOption(option, enable)
}

// SetReadBuffer sets the size of the operating system's receive buffer
// associated with the Conn.
func (c *Conn) SetReadBuffer(bytes int) error {
	return c.c.SetReadBuffer(bytes)
}

// SetWriteBuffer sets the size of the operating system's transmit buffer
// associated with the Conn.
func (c *Conn) SetWriteBuffer(bytes int) error {
	return c.c.SetWriteBuffer(bytes)
}

// SyscallConn returns a raw network connection. This implements the
// syscall.Conn interface.
//
// Only the Control method of the returned syscall.RawConn is currently
// implemented.
//
// SyscallConn is intended for advanced use cases, such as getting and setting
// arbitrary socket options using the netlink socket's file descriptor.
//
// Once invoked, it is the caller's responsibility to ensure that operations
// performed using Conn and the syscall.RawConn do not conflict with
// each other.
func (c *Conn) SyscallConn() (syscall.RawConn, error) {
	return c.c.SyscallConn()
}

// Send sends a single Message to netlink, wrapping it in a netlink.Message
// using the specified generic netlink family and flags.  On success, Send
// returns a copy of the netlink.Message with all parameters populated, for
// later validation.
func (c *Conn) Send(m Message, family uint16, flags netlink.HeaderFlags) (netlink.Message, error) {
	nm := netlink.Message{
		Header: netlink.Header{
			Type:  netlink.HeaderType(family),
			Flags: flags,
		},
	}

	mb, err := m.MarshalBinary()
	if err != nil {
		return netlink.Message{}, err
	}
	nm.Data = mb

	reqnm, err := c.c.Send(nm)
	if err != nil {
		return netlink.Message{}, err
	}

	return reqnm, nil
}

// Receive receives one or more Messages from netlink.  The netlink.Messages
// used to wrap each Message are available for later validation.
func (c *Conn) Receive() ([]Message, []netlink.Message, error) {
	msgs, err := c.c.Receive()
	if err != nil {
		return nil, nil, err
	}

	gmsgs := make([]Message, 0, len(msgs))
	for _, nm := range msgs {
		var gm Message
		if err := (&gm).UnmarshalBinary(nm.Data); err != nil {
			return nil, nil, err
		}

		gmsgs = append(gmsgs, gm)
	}

	return gmsgs, msgs, nil
}

// Execute sends a single Message to netlink using Conn.Send, receives one or
// more replies using Conn.Receive, and then checks the validity of the replies
// against the request using netlink.Validate.
//
// See the documentation of Conn.Send, Conn.Receive, and netlink.Validate for
// details about each function.
func (c *Conn) Execute(m Message, family uint16, flags netlink.HeaderFlags) ([]Message, error) {
	req, err := c.Send(m, family, flags)
	if err != nil {
		return nil, err
	}

	msgs, replies, err := c.Receive()
	if err != nil {
		return nil, err
	}

	if err := netlink.Validate(req, replies); err != nil {
		return nil, err
	}

	return msgs, nil
}
