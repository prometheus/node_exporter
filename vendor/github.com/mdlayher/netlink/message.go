package netlink

import (
	"errors"
	"fmt"

	"github.com/mdlayher/netlink/nlenc"
)

// Various errors which may occur when attempting to marshal or unmarshal
// a Message to and from its binary form.
var (
	errIncorrectMessageLength = errors.New("netlink message header length incorrect")
	errShortMessage           = errors.New("not enough data to create a netlink message")
	errUnalignedMessage       = errors.New("input data is not properly aligned for netlink message")
)

// HeaderFlags specify flags which may be present in a Header.
type HeaderFlags uint16

const (
	// General netlink communication flags.

	// HeaderFlagsRequest indicates a request to netlink.
	HeaderFlagsRequest HeaderFlags = 1

	// HeaderFlagsMulti indicates a multi-part message, terminated
	// by HeaderTypeDone on the last message.
	HeaderFlagsMulti HeaderFlags = 2

	// HeaderFlagsAcknowledge requests that netlink reply with
	// an acknowledgement using HeaderTypeError and, if needed,
	// an error code.
	HeaderFlagsAcknowledge HeaderFlags = 4

	// HeaderFlagsEcho requests that netlink echo this request
	// back to the sender.
	HeaderFlagsEcho HeaderFlags = 8

	// HeaderFlagsDumpInterrupted indicates that a dump was
	// inconsistent due to a sequence change.
	HeaderFlagsDumpInterrupted HeaderFlags = 16

	// HeaderFlagsDumpFiltered indicates that a dump was filtered
	// as requested.
	HeaderFlagsDumpFiltered HeaderFlags = 32

	// Flags used to retrieve data from netlink.

	// HeaderFlagsRoot requests that netlink return a complete table instead
	// of a single entry.
	HeaderFlagsRoot HeaderFlags = 0x100

	// HeaderFlagsMatch requests that netlink return a list of all matching
	// entries.
	HeaderFlagsMatch HeaderFlags = 0x200

	// HeaderFlagsAtomic requests that netlink send an atomic snapshot of
	// its entries.  Requires CAP_NET_ADMIN or an effective UID of 0.
	HeaderFlagsAtomic HeaderFlags = 0x400

	// HeaderFlagsDump requests that netlink return a complete list of
	// all entries.
	HeaderFlagsDump HeaderFlags = HeaderFlagsRoot | HeaderFlagsMatch

	// Flags used to create objects.

	// HeaderFlagsReplace indicates request replaces an existing matching object.
	HeaderFlagsReplace HeaderFlags = 0x100

	// HeaderFlagsExcl  indicates request does not replace the object if it already exists.
	HeaderFlagsExcl HeaderFlags = 0x200

	// HeaderFlagsCreate indicates request creates an object if it doesn't already exist.
	HeaderFlagsCreate HeaderFlags = 0x400

	// HeaderFlagsAppend indicates request adds to the end of the object list.
	HeaderFlagsAppend HeaderFlags = 0x800
)

// String returns the string representation of a HeaderFlags.
func (f HeaderFlags) String() string {
	names := []string{
		"request",
		"multi",
		"acknowledge",
		"echo",
		"dumpinterrupted",
		"dumpfiltered",
	}

	var s string

	left := uint(f)

	for i, name := range names {
		if f&(1<<uint(i)) != 0 {
			if s != "" {
				s += "|"
			}

			s += name

			left ^= (1 << uint(i))
		}
	}

	if s == "" && left == 0 {
		s = "0"
	}

	if left > 0 {
		if s != "" {
			s += "|"
		}
		s += fmt.Sprintf("%#x", left)
	}

	return s
}

// HeaderType specifies the type of a Header.
type HeaderType uint16

const (
	// HeaderTypeNoop indicates that no action was taken.
	HeaderTypeNoop HeaderType = 0x1

	// HeaderTypeError indicates an error code is present, which is also
	// used to indicate success when the code is 0.
	HeaderTypeError HeaderType = 0x2

	// HeaderTypeDone indicates the end of a multi-part message.
	HeaderTypeDone HeaderType = 0x3

	// HeaderTypeOverrun indicates that data was lost from this message.
	HeaderTypeOverrun HeaderType = 0x4
)

// String returns the string representation of a HeaderType.
func (t HeaderType) String() string {
	switch t {
	case HeaderTypeNoop:
		return "noop"
	case HeaderTypeError:
		return "error"
	case HeaderTypeDone:
		return "done"
	case HeaderTypeOverrun:
		return "overrun"
	default:
		return fmt.Sprintf("unknown(%d)", t)
	}
}

// NB: the memory layout of Header and Linux's syscall.NlMsgHdr must be
// exactly the same.  Cannot reorder, change data type, add, or remove fields.
// Named types of the same size (e.g. HeaderFlags is a uint16) are okay.

// A Header is a netlink header.  A Header is sent and received with each
// Message to indicate metadata regarding a Message.
type Header struct {
	// Length of a Message, including this Header.
	Length uint32

	// Contents of a Message.
	Type HeaderType

	// Flags which may be used to modify a request or response.
	Flags HeaderFlags

	// The sequence number of a Message.
	Sequence uint32

	// The process ID of the sending process.
	PID uint32
}

// A Message is a netlink message.  It contains a Header and an arbitrary
// byte payload, which may be decoded using information from the Header.
//
// Data is encoded in the native endianness of the host system.  For easier
// of encoding and decoding of integers, use package nlenc.
type Message struct {
	Header Header
	Data   []byte
}

// MarshalBinary marshals a Message into a byte slice.
func (m Message) MarshalBinary() ([]byte, error) {
	ml := nlmsgAlign(int(m.Header.Length))
	if ml < nlmsgHeaderLen || ml != int(m.Header.Length) {
		return nil, errIncorrectMessageLength
	}

	b := make([]byte, ml)

	nlenc.PutUint32(b[0:4], m.Header.Length)
	nlenc.PutUint16(b[4:6], uint16(m.Header.Type))
	nlenc.PutUint16(b[6:8], uint16(m.Header.Flags))
	nlenc.PutUint32(b[8:12], m.Header.Sequence)
	nlenc.PutUint32(b[12:16], m.Header.PID)
	copy(b[16:], m.Data)

	return b, nil
}

// UnmarshalBinary unmarshals the contents of a byte slice into a Message.
func (m *Message) UnmarshalBinary(b []byte) error {
	if len(b) < nlmsgHeaderLen {
		return errShortMessage
	}
	if len(b) != nlmsgAlign(len(b)) {
		return errUnalignedMessage
	}

	// Don't allow misleading length
	m.Header.Length = nlenc.Uint32(b[0:4])
	if int(m.Header.Length) != len(b) {
		return errShortMessage
	}

	m.Header.Type = HeaderType(nlenc.Uint16(b[4:6]))
	m.Header.Flags = HeaderFlags(nlenc.Uint16(b[6:8]))
	m.Header.Sequence = nlenc.Uint32(b[8:12])
	m.Header.PID = nlenc.Uint32(b[12:16])
	m.Data = b[16:]

	return nil
}

// checkMessage checks a single Message for netlink errors.
func checkMessage(m Message) error {
	const success = 0

	// Per libnl documentation, only messages that indicate type error can
	// contain error codes:
	// https://www.infradead.org/~tgr/libnl/doc/core.html#core_errmsg.
	//
	// However, at one point, this package checked both done and error for
	// error codes.  Because there was no issue associated with the change,
	// it is unknown whether this change was correct or not.  If you run into
	// a problem with your application because of this change, please file
	// an issue.
	if m.Header.Type != HeaderTypeError {
		return nil
	}

	if len(m.Data) < 4 {
		return errShortErrorMessage
	}

	if c := nlenc.Int32(m.Data[0:4]); c != success {
		// Error code is a negative integer, convert it into
		// an OS-specific system call error
		return newError(-1 * int(c))
	}

	return nil
}
