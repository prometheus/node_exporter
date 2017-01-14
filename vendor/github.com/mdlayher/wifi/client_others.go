//+build !linux

package wifi

import (
	"fmt"
	"runtime"
)

var (
	// errUnimplemented is returned by all functions on platforms that
	// do not have package wifi implemented.
	errUnimplemented = fmt.Errorf("package wifi not implemented on %s/%s",
		runtime.GOOS, runtime.GOARCH)
)

var _ osClient = &client{}

// A conn is the no-op implementation of a netlink sockets connection.
type client struct{}

// newClient always returns an error.
func newClient() (*client, error) {
	return nil, errUnimplemented
}

// Close always returns an error.
func (c *client) Close() error {
	return errUnimplemented
}

// Interfaces always returns an error.
func (c *client) Interfaces() ([]*Interface, error) {
	return nil, errUnimplemented
}

// StationInfo always returns an error.
func (c *client) StationInfo(ifi *Interface) (*StationInfo, error) {
	return nil, errUnimplemented
}
