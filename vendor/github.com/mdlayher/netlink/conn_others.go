//+build !linux

package netlink

import (
	"fmt"
	"runtime"
)

var (
	// errUnimplemented is returned by all functions on platforms that
	// cannot make use of netlink sockets.
	errUnimplemented = fmt.Errorf("netlink sockets not implemented on %s/%s",
		runtime.GOOS, runtime.GOARCH)
)

var _ osConn = &conn{}

// A conn is the no-op implementation of a netlink sockets connection.
type conn struct{}

// dial is the entry point for Dial.  dial always returns an error.
func dial(family int, config *Config) (*conn, error) {
	return nil, errUnimplemented
}

// Send always returns an error.
func (c *conn) Send(m Message) error {
	return errUnimplemented
}

// Receive always returns an error.
func (c *conn) Receive() ([]Message, error) {
	return nil, errUnimplemented
}

// Close always returns an error.
func (c *conn) Close() error {
	return errUnimplemented
}

// JoinGroup always returns an error.
func (c *conn) JoinGroup(group uint32) error {
	return errUnimplemented
}

// LeaveGroup always returns an error.
func (c *conn) LeaveGroup(group uint32) error {
	return errUnimplemented
}

// newError always returns an error.
func newError(errno int) error {
	return errUnimplemented
}
