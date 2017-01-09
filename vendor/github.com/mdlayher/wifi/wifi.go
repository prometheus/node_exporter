package wifi

import (
	"net"
	"time"
)

// An InterfaceType is the operating mode of an Interface.
type InterfaceType int

const (
	// InterfaceTypeUnspecified indicates that an interface's type is unspecified
	// and the driver determines its function.
	InterfaceTypeUnspecified InterfaceType = iota

	// InterfaceTypeAdHoc indicates that an interface is part of an independent
	// basic service set (BSS) of client devices without a controlling access
	// point.
	InterfaceTypeAdHoc

	// InterfaceTypeStation indicates that an interface is part of a managed
	// basic service set (BSS) of client devices with a controlling access point.
	InterfaceTypeStation

	// InterfaceTypeAP indicates that an interface is an access point.
	InterfaceTypeAP

	// InterfaceTypeAPVLAN indicates that an interface is a VLAN interface
	// associated with an access point.
	InterfaceTypeAPVLAN

	// InterfaceTypeWDS indicates that an interface is a wireless distribution
	// interface, used as part of a network of multiple access points.
	InterfaceTypeWDS

	// InterfaceTypeMonitor indicates that an interface is a monitor interface,
	// receiving all frames from all clients in a given network.
	InterfaceTypeMonitor

	// InterfaceTypeMeshPoint indicates that an interface is part of a wireless
	// mesh network.
	InterfaceTypeMeshPoint

	// InterfaceTypeP2PClient indicates that an interface is a client within
	// a peer-to-peer network.
	InterfaceTypeP2PClient

	// InterfaceTypeP2PGroupOwner indicates that an interface is the group
	// owner within a peer-to-peer network.
	InterfaceTypeP2PGroupOwner

	// InterfaceTypeP2PDevice indicates that an interface is a device within
	// a peer-to-peer client network.
	InterfaceTypeP2PDevice

	// InterfaceTypeOCB indicates that an interface is outside the context
	// of a basic service set (BSS).
	InterfaceTypeOCB

	// InterfaceTypeNAN indicates that an interface is part of a near-me
	// area network (NAN).
	InterfaceTypeNAN
)

// An Interface is a WiFi network interface.
type Interface struct {
	// The index of the interface.
	Index int

	// The name of the interface.
	Name string

	// The hardware address of the interface.
	HardwareAddr net.HardwareAddr

	// The physical device that this interface belongs to.
	PHY int

	// The virtual device number of this interface within a PHY.
	Device int

	// The operating mode of the interface.
	Type InterfaceType

	// The interface's wireless frequency in MHz.
	Frequency int
}

// StationInfo contains statistics about a WiFi interface operating in
// station mode.
type StationInfo struct {
	// The time since the station last connected.
	Connected time.Duration

	// The time since wireless activity last occurred.
	Inactive time.Duration

	// The number of bytes received by this station.
	ReceivedBytes uint64

	// The number of bytes transmitted by this station.
	TransmittedBytes uint64

	// The number of packets received by this station.
	ReceivedPackets uint32

	// The number of packets transmitted by this station.
	TransmittedPackets uint32

	// The current data receive bitrate, in bits/second.
	ReceiveBitrate int

	// The current data transmit bitrate, in bits/second.
	TransmitBitrate int

	// The signal strength of this station's connection, in dBm.
	Signal int

	// The number of times the station has had to retry while sending a packet.
	TransmitRetries uint32

	// The number of times a packet transmission failed.
	TransmitFailed uint32

	// The number of times a beacon loss was detected.
	BeaconLoss uint32
}
