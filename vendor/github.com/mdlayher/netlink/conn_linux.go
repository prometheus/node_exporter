//+build linux

package netlink

import (
	"errors"
	"os"
	"syscall"
	"unsafe"
)

var (
	errInvalidSockaddr = errors.New("expected syscall.SockaddrNetlink but received different syscall.Sockaddr")
	errInvalidFamily   = errors.New("received invalid netlink family")
)

var _ osConn = &conn{}

// A conn is the Linux implementation of a netlink sockets connection.
type conn struct {
	s  socket
	sa *syscall.SockaddrNetlink
}

// A socket is an interface over socket system calls.
type socket interface {
	Bind(sa syscall.Sockaddr) error
	Close() error
	Recvfrom(p []byte, flags int) (int, syscall.Sockaddr, error)
	Sendto(p []byte, flags int, to syscall.Sockaddr) error
	SetSockopt(level, name int, v unsafe.Pointer, l uint32) error
}

// dial is the entry point for Dial.  dial opens a netlink socket using
// system calls.
func dial(family int, config *Config) (*conn, error) {
	fd, err := syscall.Socket(
		syscall.AF_NETLINK,
		syscall.SOCK_RAW,
		family,
	)
	if err != nil {
		return nil, err
	}

	return bind(&sysSocket{fd: fd}, config)
}

// bind binds a connection to netlink using the input socket, which may be
// a system call implementation or a mocked one for tests.
func bind(s socket, config *Config) (*conn, error) {
	if config == nil {
		config = &Config{}
	}

	addr := &syscall.SockaddrNetlink{
		Family: syscall.AF_NETLINK,
		Groups: config.Groups,
	}

	if err := s.Bind(addr); err != nil {
		return nil, err
	}

	return &conn{
		s:  s,
		sa: addr,
	}, nil
}

// Send sends a single Message to netlink.
func (c *conn) Send(m Message) error {
	b, err := m.MarshalBinary()
	if err != nil {
		return err
	}

	return c.s.Sendto(b, 0, &syscall.SockaddrNetlink{
		Family: syscall.AF_NETLINK,
	})
}

// Receive receives one or more Messages from netlink.
func (c *conn) Receive() ([]Message, error) {
	b := make([]byte, os.Getpagesize())
	for {
		// Peek at the buffer to see how many bytes are available
		n, _, err := c.s.Recvfrom(b, syscall.MSG_PEEK)
		if err != nil {
			return nil, err
		}

		// Break when we can read all messages
		if n < len(b) {
			break
		}

		// Double in size if not enough bytes
		b = make([]byte, len(b)*2)
	}

	// Read out all available messages
	n, from, err := c.s.Recvfrom(b, 0)
	if err != nil {
		return nil, err
	}

	addr, ok := from.(*syscall.SockaddrNetlink)
	if !ok {
		return nil, errInvalidSockaddr
	}
	if addr.Family != syscall.AF_NETLINK {
		return nil, errInvalidFamily
	}

	raw, err := syscall.ParseNetlinkMessage(b[:n])
	if err != nil {
		return nil, err
	}

	msgs := make([]Message, 0, len(raw))
	for _, r := range raw {
		m := Message{
			Header: sysToHeader(r.Header),
			Data:   r.Data,
		}

		msgs = append(msgs, m)
	}

	return msgs, nil
}

// Close closes the connection.
func (c *conn) Close() error {
	return c.s.Close()
}

const (
	// #define SOL_NETLINK     270
	solNetlink = 270
)

// JoinGroup joins a multicast group by ID.
func (c *conn) JoinGroup(group uint32) error {
	return c.s.SetSockopt(
		solNetlink,
		syscall.NETLINK_ADD_MEMBERSHIP,
		unsafe.Pointer(&group),
		uint32(unsafe.Sizeof(group)),
	)
}

// LeaveGroup leaves a multicast group by ID.
func (c *conn) LeaveGroup(group uint32) error {
	return c.s.SetSockopt(
		solNetlink,
		syscall.NETLINK_DROP_MEMBERSHIP,
		unsafe.Pointer(&group),
		uint32(unsafe.Sizeof(group)),
	)
}

// sysToHeader converts a syscall.NlMsghdr to a Header.
func sysToHeader(r syscall.NlMsghdr) Header {
	// NB: the memory layout of Header and syscall.NlMsgHdr must be
	// exactly the same for this unsafe cast to work
	return *(*Header)(unsafe.Pointer(&r))
}

// newError converts an error number from netlink into the appropriate
// system call error for Linux.
func newError(errno int) error {
	return syscall.Errno(errno)
}

var _ socket = &sysSocket{}

// A sysSocket is a socket which uses system calls for socket operations.
type sysSocket struct {
	fd int
}

func (s *sysSocket) Bind(sa syscall.Sockaddr) error { return syscall.Bind(s.fd, sa) }
func (s *sysSocket) Close() error                   { return syscall.Close(s.fd) }
func (s *sysSocket) Recvfrom(p []byte, flags int) (int, syscall.Sockaddr, error) {
	return syscall.Recvfrom(s.fd, p, flags)
}
func (s *sysSocket) Sendto(p []byte, flags int, to syscall.Sockaddr) error {
	return syscall.Sendto(s.fd, p, flags, to)
}
func (s *sysSocket) SetSockopt(level, name int, v unsafe.Pointer, l uint32) error {
	return setsockopt(s.fd, level, name, v, l)
}
