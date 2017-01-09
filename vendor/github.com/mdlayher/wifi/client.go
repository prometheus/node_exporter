package wifi

// A Client is a type which can access WiFi device actions and statistics
// using operating system-specific operations.
type Client struct {
	c osClient
}

// New creates a new Client.
func New() (*Client, error) {
	c, err := newClient()
	if err != nil {
		return nil, err
	}

	return &Client{
		c: c,
	}, nil
}

// Interfaces returns a list of the system's WiFi network interfaces.
func (c *Client) Interfaces() ([]*Interface, error) {
	return c.c.Interfaces()
}

// StationInfo retrieves statistics about a WiFi interface operating in
// station mode.
func (c *Client) StationInfo(ifi *Interface) (*StationInfo, error) {
	return c.c.StationInfo(ifi)
}

// An osClient is the operating system-specific implementation of Client.
type osClient interface {
	Interfaces() ([]*Interface, error)
	StationInfo(ifi *Interface) (*StationInfo, error)
}
