//+build linux

package netlink

import (
	"errors"
	"os"
	"runtime"
	"sync"
	"syscall"
	"unsafe"

	"golang.org/x/net/bpf"
	"golang.org/x/sys/unix"
)

var (
	errInvalidSockaddr = errors.New("expected unix.SockaddrNetlink but received different unix.Sockaddr")
	errInvalidFamily   = errors.New("received invalid netlink family")
)

var _ Socket = &conn{}

// A conn is the Linux implementation of a netlink sockets connection.
type conn struct {
	s  socket
	sa *unix.SockaddrNetlink
}

// A socket is an interface over socket system calls.
type socket interface {
	Bind(sa unix.Sockaddr) error
	Close() error
	FD() int
	Getsockname() (unix.Sockaddr, error)
	Recvmsg(p, oob []byte, flags int) (n int, oobn int, recvflags int, from unix.Sockaddr, err error)
	Sendmsg(p, oob []byte, to unix.Sockaddr, flags int) error
	SetSockopt(level, name int, v unsafe.Pointer, l uint32) error
}

// dial is the entry point for Dial.  dial opens a netlink socket using
// system calls, and returns its PID.
func dial(family int, config *Config) (*conn, uint32, error) {
	// Prepare sysSocket's internal loop and create the socket.
	//
	// The conditional is inverted because a zero value of false is desired
	// if no config, but it's easier to interpret within this code when the
	// value is inverted.
	if config == nil {
		config = &Config{}
	}

	lockThread := !config.NoLockThread
	sock := newSysSocket(lockThread)

	if err := sock.Socket(family); err != nil {
		return nil, 0, err
	}

	return bind(sock, config)
}

// bind binds a connection to netlink using the input socket, which may be
// a system call implementation or a mocked one for tests.
func bind(s socket, config *Config) (*conn, uint32, error) {
	if config == nil {
		config = &Config{}
	}

	addr := &unix.SockaddrNetlink{
		Family: unix.AF_NETLINK,
		Groups: config.Groups,
	}

	// Socket must be closed in the event of any system call errors, to avoid
	// leaking file descriptors.

	if err := s.Bind(addr); err != nil {
		_ = s.Close()
		return nil, 0, err
	}

	sa, err := s.Getsockname()
	if err != nil {
		_ = s.Close()
		return nil, 0, err
	}

	pid := sa.(*unix.SockaddrNetlink).Pid

	return &conn{
		s:  s,
		sa: addr,
	}, pid, nil
}

// Send sends a single Message to netlink.
func (c *conn) Send(m Message) error {
	b, err := m.MarshalBinary()
	if err != nil {
		return err
	}

	addr := &unix.SockaddrNetlink{
		Family: unix.AF_NETLINK,
	}

	return c.s.Sendmsg(b, nil, addr, 0)
}

// Receive receives one or more Messages from netlink.
func (c *conn) Receive() ([]Message, error) {
	b := make([]byte, os.Getpagesize())
	for {
		// Peek at the buffer to see how many bytes are available
		n, _, _, _, err := c.s.Recvmsg(b, nil, unix.MSG_PEEK)
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
	n, _, _, from, err := c.s.Recvmsg(b, nil, 0)
	if err != nil {
		return nil, err
	}

	addr, ok := from.(*unix.SockaddrNetlink)
	if !ok {
		return nil, errInvalidSockaddr
	}
	if addr.Family != unix.AF_NETLINK {
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

// FD retrieves the file descriptor of the Conn.
func (c *conn) FD() int {
	return c.s.FD()
}

// JoinGroup joins a multicast group by ID.
func (c *conn) JoinGroup(group uint32) error {
	return c.s.SetSockopt(
		unix.SOL_NETLINK,
		unix.NETLINK_ADD_MEMBERSHIP,
		unsafe.Pointer(&group),
		uint32(unsafe.Sizeof(group)),
	)
}

// LeaveGroup leaves a multicast group by ID.
func (c *conn) LeaveGroup(group uint32) error {
	return c.s.SetSockopt(
		unix.SOL_NETLINK,
		unix.NETLINK_DROP_MEMBERSHIP,
		unsafe.Pointer(&group),
		uint32(unsafe.Sizeof(group)),
	)
}

// SetBPF attaches an assembled BPF program to a conn.
func (c *conn) SetBPF(filter []bpf.RawInstruction) error {
	prog := unix.SockFprog{
		Len:    uint16(len(filter)),
		Filter: (*unix.SockFilter)(unsafe.Pointer(&filter[0])),
	}

	return c.s.SetSockopt(
		unix.SOL_SOCKET,
		unix.SO_ATTACH_FILTER,
		unsafe.Pointer(&prog),
		uint32(unsafe.Sizeof(prog)),
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

	wg    *sync.WaitGroup
	funcC chan<- func()
}

// newSysSocket creates a sysSocket that optionally locks its internal goroutine
// to a single thread.
func newSysSocket(lockThread bool) *sysSocket {
	var wg sync.WaitGroup
	wg.Add(1)

	// This system call loop strategy was inspired by:
	// https://github.com/golang/go/wiki/LockOSThread.  Thanks to squeed on
	// Gophers Slack for providing this useful link.

	funcC := make(chan func())
	go func() {
		// It is important to lock this goroutine to its OS thread for the duration
		// of the netlink socket being used, or else the kernel may end up routing
		// messages to the wrong places.
		// See: http://lists.infradead.org/pipermail/libnl/2017-February/002293.html.
		//
		// But since this is very experimental, we'll leave it as a configurable at
		// this point.
		if lockThread {
			// Never unlock the OS thread, so that the thread will terminate when
			// the goroutine exits starting in Go 1.10:
			// https://go-review.googlesource.com/c/go/+/46038.
			runtime.LockOSThread()
		}

		defer wg.Done()

		for f := range funcC {
			f()
		}
	}()

	return &sysSocket{
		wg:    &wg,
		funcC: funcC,
	}
}

// do runs f in a worker goroutine which can be locked to one thread.
func (s *sysSocket) do(f func()) {
	done := make(chan bool, 1)
	s.funcC <- func() {
		f()
		done <- true
	}
	<-done
}

func (s *sysSocket) Socket(family int) error {
	var (
		fd  int
		err error
	)

	s.do(func() {
		fd, err = unix.Socket(
			unix.AF_NETLINK,
			unix.SOCK_RAW,
			family,
		)
	})
	if err != nil {
		return err
	}

	s.fd = fd
	return nil
}

func (s *sysSocket) Bind(sa unix.Sockaddr) error {
	var err error
	s.do(func() {
		err = unix.Bind(s.fd, sa)
	})

	return err
}

func (s *sysSocket) Close() error {
	var err error
	s.do(func() {
		err = unix.Close(s.fd)
	})

	close(s.funcC)
	s.wg.Wait()

	return err
}

func (s *sysSocket) FD() int { return s.fd }

func (s *sysSocket) Getsockname() (unix.Sockaddr, error) {
	var (
		sa  unix.Sockaddr
		err error
	)

	s.do(func() {
		sa, err = unix.Getsockname(s.fd)
	})

	return sa, err
}
func (s *sysSocket) Recvmsg(p, oob []byte, flags int) (int, int, int, unix.Sockaddr, error) {
	var (
		n, oobn, recvflags int
		from               unix.Sockaddr
		err                error
	)

	s.do(func() {
		n, oobn, recvflags, from, err = unix.Recvmsg(s.fd, p, oob, flags)
	})

	return n, oobn, recvflags, from, err
}

func (s *sysSocket) Sendmsg(p, oob []byte, to unix.Sockaddr, flags int) error {
	var err error
	s.do(func() {
		err = unix.Sendmsg(s.fd, p, oob, to, flags)
	})

	return err
}

func (s *sysSocket) SetSockopt(level, name int, v unsafe.Pointer, l uint32) error {
	var err error
	s.do(func() {
		err = setsockopt(s.fd, level, name, v, l)
	})

	return err
}
