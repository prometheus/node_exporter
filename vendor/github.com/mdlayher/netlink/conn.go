package netlink

import (
	"errors"
	"io"
	"math/rand"
	"os"
	"sync/atomic"
	"syscall"
	"time"

	"golang.org/x/net/bpf"
)

// Error messages which can be returned by Validate.
var (
	errMismatchedSequence = errors.New("mismatched sequence in netlink reply")
	errMismatchedPID      = errors.New("mismatched PID in netlink reply")
	errShortErrorMessage  = errors.New("not enough data for netlink error code")
)

// Errors which can be returned by a Socket that does not implement
// all exposed methods of Conn.
var (
	errReadWriteCloserNotSupported = errors.New("raw read/write/closer not supported")
	errMulticastGroupsNotSupported = errors.New("multicast groups not supported")
	errBPFFiltersNotSupported      = errors.New("BPF filters not supported")
	errOptionsNotSupported         = errors.New("options not supported")
	errSetBufferNotSupported       = errors.New("setting buffer sizes not supported")
	errSyscallConnNotSupported     = errors.New("syscall.RawConn operation not supported")
)

// A Conn is a connection to netlink.  A Conn can be used to send and
// receives messages to and from netlink.
//
// A Conn is safe for concurrent use, but to avoid contention in
// high-throughput applications, the caller should almost certainly create a
// pool of Conns and distribute them among workers.
type Conn struct {
	// sock is the operating system-specific implementation of
	// a netlink sockets connection.
	sock Socket

	// seq is an atomically incremented integer used to provide sequence
	// numbers when Conn.Send is called.
	seq *uint32

	// pid is the PID assigned by netlink.
	pid uint32

	// d provides debugging capabilities for a Conn if not nil.
	d *debugger
}

// A Socket is an operating-system specific implementation of netlink
// sockets used by Conn.
type Socket interface {
	Close() error
	Send(m Message) error
	SendMessages(m []Message) error
	Receive() ([]Message, error)
}

// Dial dials a connection to netlink, using the specified netlink family.
// Config specifies optional configuration for Conn.  If config is nil, a default
// configuration will be used.
func Dial(family int, config *Config) (*Conn, error) {
	// Use OS-specific dial() to create Socket
	c, pid, err := dial(family, config)
	if err != nil {
		return nil, err
	}

	return NewConn(c, pid), nil
}

// NewConn creates a Conn using the specified Socket and PID for netlink
// communications.
//
// NewConn is primarily useful for tests. Most applications should use
// Dial instead.
func NewConn(sock Socket, pid uint32) *Conn {
	// Seed the sequence number using a random number generator.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	seq := r.Uint32()

	// Configure a debugger if arguments are set.
	var d *debugger
	if len(debugArgs) > 0 {
		d = newDebugger(debugArgs)
	}

	return &Conn{
		sock: sock,
		seq:  &seq,
		pid:  pid,
		d:    d,
	}
}

// debug executes fn with the debugger if the debugger is not nil.
func (c *Conn) debug(fn func(d *debugger)) {
	if c.d == nil {
		return
	}

	fn(c.d)
}

// Close closes the connection.
func (c *Conn) Close() error {
	return c.sock.Close()
}

// Execute sends a single Message to netlink using Conn.Send, receives one or more
// replies using Conn.Receive, and then checks the validity of the replies against
// the request using Validate.
//
// See the documentation of Conn.Send, Conn.Receive, and Validate for details about
// each function.
func (c *Conn) Execute(message Message) ([]Message, error) {
	req, err := c.Send(message)
	if err != nil {
		return nil, err
	}

	replies, err := c.Receive()
	if err != nil {
		return nil, err
	}

	if err := Validate(req, replies); err != nil {
		return nil, err
	}

	return replies, nil
}

func (c *Conn) fixMsg(m *Message, ml int) {
	if m.Header.Length == 0 {
		m.Header.Length = uint32(nlmsgAlign(ml))
	}

	if m.Header.Sequence == 0 {
		m.Header.Sequence = c.nextSequence()
	}

	if m.Header.PID == 0 {
		m.Header.PID = c.pid
	}
}

// SendMessages sends multiple Messages to netlink. The handling of
// a Header's Length, Sequence and PID fields is the same as when
// calling Send.
func (c *Conn) SendMessages(messages []Message) ([]Message, error) {
	for idx, m := range messages {
		ml := nlmsgLength(len(m.Data))

		// TODO(mdlayher): fine-tune this limit.
		if ml > (1024 * 32) {
			return nil, errors.New("netlink message data too large")
		}

		c.fixMsg(&messages[idx], ml)
	}

	c.debug(func(d *debugger) {
		for _, m := range messages {
			d.debugf(1, "send msgs: %+v", m)
		}
	})

	if err := c.sock.SendMessages(messages); err != nil {
		c.debug(func(d *debugger) {
			d.debugf(1, "send msgs: err: %v", err)
		})

		return nil, err
	}

	return messages, nil
}

// Send sends a single Message to netlink.  In most cases, a Header's Length,
// Sequence, and PID fields should be set to 0, so they can be populated
// automatically before the Message is sent.  On success, Send returns a copy
// of the Message with all parameters populated, for later validation.
//
// If Header.Length is 0, it will be automatically populated using the
// correct length for the Message, including its payload.
//
// If Header.Sequence is 0, it will be automatically populated using the
// next sequence number for this connection.
//
// If Header.PID is 0, it will be automatically populated using a PID
// assigned by netlink.
func (c *Conn) Send(message Message) (Message, error) {
	ml := nlmsgLength(len(message.Data))

	// TODO(mdlayher): fine-tune this limit.
	if ml > (1024 * 32) {
		return Message{}, errors.New("netlink message data too large")
	}

	c.fixMsg(&message, ml)

	c.debug(func(d *debugger) {
		d.debugf(1, "send: %+v", message)
	})

	if err := c.sock.Send(message); err != nil {
		c.debug(func(d *debugger) {
			d.debugf(1, "send: err: %v", err)
		})

		return Message{}, err
	}

	return message, nil
}

// Receive receives one or more messages from netlink.  Multi-part messages are
// handled transparently and returned as a single slice of Messages, with the
// final empty "multi-part done" message removed.
//
// If any of the messages indicate a netlink error, that error will be returned.
func (c *Conn) Receive() ([]Message, error) {
	msgs, err := c.receive()
	if err != nil {
		c.debug(func(d *debugger) {
			d.debugf(1, "recv: err: %v", err)
		})

		return nil, err
	}

	c.debug(func(d *debugger) {
		for _, m := range msgs {
			d.debugf(1, "recv: %+v", m)
		}
	})

	// When using nltest, it's possible for zero messages to be returned by receive.
	if len(msgs) == 0 {
		return msgs, nil
	}

	// Trim the final message with multi-part done indicator if
	// present.
	if m := msgs[len(msgs)-1]; m.Header.Flags&HeaderFlagsMulti != 0 && m.Header.Type == HeaderTypeDone {
		return msgs[:len(msgs)-1], nil
	}

	return msgs, nil
}

// receive is the internal implementation of Conn.Receive, which can be called
// recursively to handle multi-part messages.
func (c *Conn) receive() ([]Message, error) {
	var res []Message
	for {
		msgs, err := c.sock.Receive()
		if err != nil {
			return nil, err
		}

		// If this message is multi-part, we will need to perform an recursive call
		// to continue draining the socket
		var multi bool

		for _, m := range msgs {
			if err := checkMessage(m); err != nil {
				return nil, err
			}

			// Does this message indicate a multi-part message?
			if m.Header.Flags&HeaderFlagsMulti == 0 {
				// No, check the next messages.
				continue
			}

			// Does this message indicate the last message in a series of
			// multi-part messages from a single read?
			multi = m.Header.Type != HeaderTypeDone
		}

		res = append(res, msgs...)

		if !multi {
			// No more messages coming.
			return res, nil
		}
	}
}

// An fder is a Socket that supports retrieving its raw file descriptor.
type fder interface {
	Socket
	FD() int
}

var _ io.ReadWriteCloser = &fileReadWriteCloser{}

// A fileReadWriteCloser is a limited *os.File which only allows access to its
// Read and Write methods.
type fileReadWriteCloser struct {
	f *os.File
}

// Read implements io.ReadWriteCloser.
func (rwc *fileReadWriteCloser) Read(b []byte) (int, error) { return rwc.f.Read(b) }

// Write implements io.ReadWriteCloser.
func (rwc *fileReadWriteCloser) Write(b []byte) (int, error) { return rwc.f.Write(b) }

// Close implements io.ReadWriteCloser.
func (rwc *fileReadWriteCloser) Close() error { return rwc.f.Close() }

// ReadWriteCloser returns a raw io.ReadWriteCloser backed by the connection
// of the Conn.
//
// ReadWriteCloser is intended for advanced use cases, such as those that do
// not involve standard netlink message passing.
//
// Once invoked, it is the caller's responsibility to ensure that operations
// performed using Conn and the raw io.ReadWriteCloser do not conflict with
// each other.  In almost all scenarios, only one of the two should be used.
func (c *Conn) ReadWriteCloser() (io.ReadWriteCloser, error) {
	fc, ok := c.sock.(fder)
	if !ok {
		return nil, errReadWriteCloserNotSupported
	}

	return &fileReadWriteCloser{
		// Backing the io.ReadWriteCloser with an *os.File enables easy reading
		// and writing without more system call boilerplate.
		f: os.NewFile(uintptr(fc.FD()), "netlink"),
	}, nil
}

// A groupJoinLeaver is a Socket that supports joining and leaving
// netlink multicast groups.
type groupJoinLeaver interface {
	Socket
	JoinGroup(group uint32) error
	LeaveGroup(group uint32) error
}

// JoinGroup joins a netlink multicast group by its ID.
func (c *Conn) JoinGroup(group uint32) error {
	gc, ok := c.sock.(groupJoinLeaver)
	if !ok {
		return errMulticastGroupsNotSupported
	}

	return gc.JoinGroup(group)
}

// LeaveGroup leaves a netlink multicast group by its ID.
func (c *Conn) LeaveGroup(group uint32) error {
	gc, ok := c.sock.(groupJoinLeaver)
	if !ok {
		return errMulticastGroupsNotSupported
	}

	return gc.LeaveGroup(group)
}

// A bpfSetter is a Socket that supports setting and removing BPF filters.
type bpfSetter interface {
	Socket
	bpf.Setter
	RemoveBPF() error
}

// SetBPF attaches an assembled BPF program to a Conn.
func (c *Conn) SetBPF(filter []bpf.RawInstruction) error {
	bc, ok := c.sock.(bpfSetter)
	if !ok {
		return errBPFFiltersNotSupported
	}

	return bc.SetBPF(filter)
}

// RemoveBPF removes a BPF filter from a Conn.
func (c *Conn) RemoveBPF() error {
	s, ok := c.sock.(bpfSetter)
	if !ok {
		return errBPFFiltersNotSupported
	}

	return s.RemoveBPF()
}

// A ConnOption is a boolean option that may be set for a Conn.
type ConnOption int

// Possible ConnOption values.  These constants are equivalent to the Linux
// setsockopt boolean options for netlink sockets.
const (
	PacketInfo ConnOption = iota
	BroadcastError
	NoENOBUFS
	ListenAllNSID
	CapAcknowledge
)

// An optionSetter is a Socket that supports setting netlink options.
type optionSetter interface {
	Socket
	SetOption(option ConnOption, enable bool) error
}

// SetOption enables or disables a netlink socket option for the Conn.
func (c *Conn) SetOption(option ConnOption, enable bool) error {
	fc, ok := c.sock.(optionSetter)
	if !ok {
		return errOptionsNotSupported
	}

	return fc.SetOption(option, enable)
}

// A bufferSetter is a Socket that supports setting connection buffer sizes.
type bufferSetter interface {
	Socket
	SetReadBuffer(bytes int) error
	SetWriteBuffer(bytes int) error
}

// SetReadBuffer sets the size of the operating system's receive buffer
// associated with the Conn.
func (c *Conn) SetReadBuffer(bytes int) error {
	conn, ok := c.sock.(bufferSetter)
	if !ok {
		return errSetBufferNotSupported
	}

	return conn.SetReadBuffer(bytes)
}

// SetWriteBuffer sets the size of the operating system's transmit buffer
// associated with the Conn.
func (c *Conn) SetWriteBuffer(bytes int) error {
	conn, ok := c.sock.(bufferSetter)
	if !ok {
		return errSetBufferNotSupported
	}

	return conn.SetWriteBuffer(bytes)
}

var _ syscall.Conn = &Conn{}

// TODO(mdlayher): mutex or similar to enforce syscall.RawConn contract of
// FD remaining valid for duration of calls?

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
	conn, ok := c.sock.(fder)
	if !ok {
		return nil, errSyscallConnNotSupported
	}

	return &rawConn{
		fd: uintptr(conn.FD()),
	}, nil
}

var _ syscall.RawConn = &rawConn{}

// A rawConn is a syscall.RawConn.
type rawConn struct {
	fd uintptr
}

func (rc *rawConn) Control(f func(fd uintptr)) error {
	f(rc.fd)
	return nil
}

// TODO(mdlayher): implement Read and Write?

func (rc *rawConn) Read(_ func(fd uintptr) (done bool)) error  { return errSyscallConnNotSupported }
func (rc *rawConn) Write(_ func(fd uintptr) (done bool)) error { return errSyscallConnNotSupported }

// nextSequence atomically increments Conn's sequence number and returns
// the incremented value.
func (c *Conn) nextSequence() uint32 {
	return atomic.AddUint32(c.seq, 1)
}

// Validate validates one or more reply Messages against a request Message,
// ensuring that they contain matching sequence numbers and PIDs.
func Validate(request Message, replies []Message) error {
	for _, m := range replies {
		// Check for mismatched sequence, unless:
		//   - request had no sequence, meaning we are probably validating
		//     a multicast reply
		if m.Header.Sequence != request.Header.Sequence && request.Header.Sequence != 0 {
			return errMismatchedSequence
		}

		// Check for mismatched PID, unless:
		//   - request had no PID, meaning we are either:
		//     - validating a multicast reply
		//     - netlink has not yet assigned us a PID
		//   - response had no PID, meaning it's from the kernel as a multicast reply
		if m.Header.PID != request.Header.PID && request.Header.PID != 0 && m.Header.PID != 0 {
			return errMismatchedPID
		}
	}

	return nil
}

// Config contains options for a Conn.
type Config struct {
	// Groups is a bitmask which specifies multicast groups. If set to 0,
	// no multicast group subscriptions will be made.
	Groups uint32

	// Network namespace the Conn needs to operate in. If set to 0,
	// no network namespace will be entered.
	NetNS int
}
