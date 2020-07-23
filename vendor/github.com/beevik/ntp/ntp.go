// Copyright 2015-2017 Brett Vickers.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ntp provides an implementation of a Simple NTP (SNTP) client
// capable of querying the current time from a remote NTP server.  See
// RFC5905 (https://tools.ietf.org/html/rfc5905) for more details.
//
// This approach grew out of a go-nuts post by Michael Hofmann:
// https://groups.google.com/forum/?fromgroups#!topic/golang-nuts/FlcdMU5fkLQ
package ntp

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"time"

	"golang.org/x/net/ipv4"
)

// The LeapIndicator is used to warn if a leap second should be inserted
// or deleted in the last minute of the current month.
type LeapIndicator uint8

const (
	// LeapNoWarning indicates no impending leap second.
	LeapNoWarning LeapIndicator = 0

	// LeapAddSecond indicates the last minute of the day has 61 seconds.
	LeapAddSecond = 1

	// LeapDelSecond indicates the last minute of the day has 59 seconds.
	LeapDelSecond = 2

	// LeapNotInSync indicates an unsynchronized leap second.
	LeapNotInSync = 3
)

// Internal constants
const (
	defaultNtpVersion = 4
	nanoPerSec        = 1000000000
	maxStratum        = 16
	defaultTimeout    = 5 * time.Second
	maxPollInterval   = (1 << 17) * time.Second
	maxDispersion     = 16 * time.Second
)

// Internal variables
var (
	ntpEpoch = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
)

type mode uint8

// NTP modes. This package uses only client mode.
const (
	reserved mode = 0 + iota
	symmetricActive
	symmetricPassive
	client
	server
	broadcast
	controlMessage
	reservedPrivate
)

// An ntpTime is a 64-bit fixed-point (Q32.32) representation of the number of
// seconds elapsed.
type ntpTime uint64

// Duration interprets the fixed-point ntpTime as a number of elapsed seconds
// and returns the corresponding time.Duration value.
func (t ntpTime) Duration() time.Duration {
	sec := (t >> 32) * nanoPerSec
	frac := (t & 0xffffffff) * nanoPerSec >> 32
	return time.Duration(sec + frac)
}

// Time interprets the fixed-point ntpTime as an absolute time and returns
// the corresponding time.Time value.
func (t ntpTime) Time() time.Time {
	return ntpEpoch.Add(t.Duration())
}

// toNtpTime converts the time.Time value t into its 64-bit fixed-point
// ntpTime representation.
func toNtpTime(t time.Time) ntpTime {
	nsec := uint64(t.Sub(ntpEpoch))
	sec := nsec / nanoPerSec
	// Round up the fractional component so that repeated conversions
	// between time.Time and ntpTime do not yield continually decreasing
	// results.
	frac := (((nsec - sec*nanoPerSec) << 32) + nanoPerSec - 1) / nanoPerSec
	return ntpTime(sec<<32 | frac)
}

// An ntpTimeShort is a 32-bit fixed-point (Q16.16) representation of the
// number of seconds elapsed.
type ntpTimeShort uint32

// Duration interprets the fixed-point ntpTimeShort as a number of elapsed
// seconds and returns the corresponding time.Duration value.
func (t ntpTimeShort) Duration() time.Duration {
	t64 := uint64(t)
	sec := (t64 >> 16) * nanoPerSec
	frac := (t64 & 0xffff) * nanoPerSec >> 16
	return time.Duration(sec + frac)
}

// msg is an internal representation of an NTP packet.
type msg struct {
	LiVnMode       uint8 // Leap Indicator (2) + Version (3) + Mode (3)
	Stratum        uint8
	Poll           int8
	Precision      int8
	RootDelay      ntpTimeShort
	RootDispersion ntpTimeShort
	ReferenceID    uint32
	ReferenceTime  ntpTime
	OriginTime     ntpTime
	ReceiveTime    ntpTime
	TransmitTime   ntpTime
}

// setVersion sets the NTP protocol version on the message.
func (m *msg) setVersion(v int) {
	m.LiVnMode = (m.LiVnMode & 0xc7) | uint8(v)<<3
}

// setMode sets the NTP protocol mode on the message.
func (m *msg) setMode(md mode) {
	m.LiVnMode = (m.LiVnMode & 0xf8) | uint8(md)
}

// setLeap modifies the leap indicator on the message.
func (m *msg) setLeap(li LeapIndicator) {
	m.LiVnMode = (m.LiVnMode & 0x3f) | uint8(li)<<6
}

// getVersion returns the version value in the message.
func (m *msg) getVersion() int {
	return int((m.LiVnMode >> 3) & 0x07)
}

// getMode returns the mode value in the message.
func (m *msg) getMode() mode {
	return mode(m.LiVnMode & 0x07)
}

// getLeap returns the leap indicator on the message.
func (m *msg) getLeap() LeapIndicator {
	return LeapIndicator((m.LiVnMode >> 6) & 0x03)
}

// QueryOptions contains the list of configurable options that may be used
// with the QueryWithOptions function.
type QueryOptions struct {
	Timeout      time.Duration // defaults to 5 seconds
	Version      int           // NTP protocol version, defaults to 4
	LocalAddress string        // IP address to use for the client address
	Port         int           // Server port, defaults to 123
	TTL          int           // IP TTL to use, defaults to system default
}

// A Response contains time data, some of which is returned by the NTP server
// and some of which is calculated by the client.
type Response struct {
	// Time is the transmit time reported by the server just before it
	// responded to the client's NTP query.
	Time time.Time

	// ClockOffset is the estimated offset of the client clock relative to
	// the server. Add this to the client's system clock time to obtain a
	// more accurate time.
	ClockOffset time.Duration

	// RTT is the measured round-trip-time delay estimate between the client
	// and the server.
	RTT time.Duration

	// Precision is the reported precision of the server's clock.
	Precision time.Duration

	// Stratum is the "stratum level" of the server. The smaller the number,
	// the closer the server is to the reference clock. Stratum 1 servers are
	// attached directly to the reference clock. A stratum value of 0
	// indicates the "kiss of death," which typically occurs when the client
	// issues too many requests to the server in a short period of time.
	Stratum uint8

	// ReferenceID is a 32-bit identifier identifying the server or
	// reference clock.
	ReferenceID uint32

	// ReferenceTime is the time when the server's system clock was last
	// set or corrected.
	ReferenceTime time.Time

	// RootDelay is the server's estimated aggregate round-trip-time delay to
	// the stratum 1 server.
	RootDelay time.Duration

	// RootDispersion is the server's estimated maximum measurement error
	// relative to the stratum 1 server.
	RootDispersion time.Duration

	// RootDistance is an estimate of the total synchronization distance
	// between the client and the stratum 1 server.
	RootDistance time.Duration

	// Leap indicates whether a leap second should be added or removed from
	// the current month's last minute.
	Leap LeapIndicator

	// MinError is a lower bound on the error between the client and server
	// clocks. When the client and server are not synchronized to the same
	// clock, the reported timestamps may appear to violate the principle of
	// causality. In other words, the NTP server's response may indicate
	// that a message was received before it was sent. In such cases, the
	// minimum error may be useful.
	MinError time.Duration

	// KissCode is a 4-character string describing the reason for a
	// "kiss of death" response (stratum = 0). For a list of standard kiss
	// codes, see https://tools.ietf.org/html/rfc5905#section-7.4.
	KissCode string

	// Poll is the maximum interval between successive NTP polling messages.
	// It is not relevant for simple NTP clients like this one.
	Poll time.Duration
}

// Validate checks if the response is valid for the purposes of time
// synchronization.
func (r *Response) Validate() error {
	// Handle invalid stratum values.
	if r.Stratum == 0 {
		return fmt.Errorf("kiss of death received: %s", r.KissCode)
	}
	if r.Stratum >= maxStratum {
		return errors.New("invalid stratum in response")
	}

	// Handle invalid leap second indicator.
	if r.Leap == LeapNotInSync {
		return errors.New("invalid leap second")
	}

	// Estimate the "freshness" of the time. If it exceeds the maximum
	// polling interval (~36 hours), then it cannot be considered "fresh".
	freshness := r.Time.Sub(r.ReferenceTime)
	if freshness > maxPollInterval {
		return errors.New("server clock not fresh")
	}

	// Calculate the peer synchronization distance, lambda:
	//  	lambda := RootDelay/2 + RootDispersion
	// If this value exceeds MAXDISP (16s), then the time is not suitable
	// for synchronization purposes.
	// https://tools.ietf.org/html/rfc5905#appendix-A.5.1.1.
	lambda := r.RootDelay/2 + r.RootDispersion
	if lambda > maxDispersion {
		return errors.New("invalid dispersion")
	}

	// If the server's transmit time is before its reference time, the
	// response is invalid.
	if r.Time.Before(r.ReferenceTime) {
		return errors.New("invalid time reported")
	}

	// nil means the response is valid.
	return nil
}

// Query returns a response from the remote NTP server host. It contains
// the time at which the server transmitted the response as well as other
// useful information about the time and the remote server.
func Query(host string) (*Response, error) {
	return QueryWithOptions(host, QueryOptions{})
}

// QueryWithOptions performs the same function as Query but allows for the
// customization of several query options.
func QueryWithOptions(host string, opt QueryOptions) (*Response, error) {
	m, now, err := getTime(host, opt)
	if err != nil {
		return nil, err
	}
	return parseTime(m, now), nil
}

// TimeV returns the current time using information from a remote NTP server.
// On error, it returns the local system time. The version may be 2, 3, or 4.
//
// Deprecated: TimeV is deprecated. Use QueryWithOptions instead.
func TimeV(host string, version int) (time.Time, error) {
	m, recvTime, err := getTime(host, QueryOptions{Version: version})
	if err != nil {
		return time.Now(), err
	}

	r := parseTime(m, recvTime)
	err = r.Validate()
	if err != nil {
		return time.Now(), err
	}

	// Use the clock offset to calculate the time.
	return time.Now().Add(r.ClockOffset), nil
}

// Time returns the current time using information from a remote NTP server.
// It uses version 4 of the NTP protocol. On error, it returns the local
// system time.
func Time(host string) (time.Time, error) {
	return TimeV(host, defaultNtpVersion)
}

// getTime performs the NTP server query and returns the response message
// along with the local system time it was received.
func getTime(host string, opt QueryOptions) (*msg, ntpTime, error) {
	if opt.Version == 0 {
		opt.Version = defaultNtpVersion
	}
	if opt.Version < 2 || opt.Version > 4 {
		return nil, 0, errors.New("invalid protocol version requested")
	}

	// Resolve the remote NTP server address.
	raddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(host, "123"))
	if err != nil {
		return nil, 0, err
	}

	// Resolve the local address if specified as an option.
	var laddr *net.UDPAddr
	if opt.LocalAddress != "" {
		laddr, err = net.ResolveUDPAddr("udp", net.JoinHostPort(opt.LocalAddress, "0"))
		if err != nil {
			return nil, 0, err
		}
	}

	// Override the port if requested.
	if opt.Port != 0 {
		raddr.Port = opt.Port
	}

	// Prepare a "connection" to the remote server.
	con, err := net.DialUDP("udp", laddr, raddr)
	if err != nil {
		return nil, 0, err
	}
	defer con.Close()

	// Set a TTL for the packet if requested.
	if opt.TTL != 0 {
		ipcon := ipv4.NewConn(con)
		err = ipcon.SetTTL(opt.TTL)
		if err != nil {
			return nil, 0, err
		}
	}

	// Set a timeout on the connection.
	if opt.Timeout == 0 {
		opt.Timeout = defaultTimeout
	}
	con.SetDeadline(time.Now().Add(opt.Timeout))

	// Allocate a message to hold the response.
	recvMsg := new(msg)

	// Allocate a message to hold the query.
	xmitMsg := new(msg)
	xmitMsg.setMode(client)
	xmitMsg.setVersion(opt.Version)
	xmitMsg.setLeap(LeapNotInSync)

	// To ensure privacy and prevent spoofing, try to use a random 64-bit
	// value for the TransmitTime. If crypto/rand couldn't generate a
	// random value, fall back to using the system clock. Keep track of
	// when the messsage was actually transmitted.
	bits := make([]byte, 8)
	_, err = rand.Read(bits)
	var xmitTime time.Time
	if err == nil {
		xmitMsg.TransmitTime = ntpTime(binary.BigEndian.Uint64(bits))
		xmitTime = time.Now()
	} else {
		xmitTime = time.Now()
		xmitMsg.TransmitTime = toNtpTime(xmitTime)
	}

	// Transmit the query.
	err = binary.Write(con, binary.BigEndian, xmitMsg)
	if err != nil {
		return nil, 0, err
	}

	// Receive the response.
	err = binary.Read(con, binary.BigEndian, recvMsg)
	if err != nil {
		return nil, 0, err
	}

	// Keep track of the time the response was received.
	delta := time.Since(xmitTime)
	if delta < 0 {
		// The local system may have had its clock adjusted since it
		// sent the query. In go 1.9 and later, time.Since ensures
		// that a monotonic clock is used, so delta can never be less
		// than zero. In versions before 1.9, a monotonic clock is
		// not used, so we have to check.
		return nil, 0, errors.New("client clock ticked backwards")
	}
	recvTime := toNtpTime(xmitTime.Add(delta))

	// Check for invalid fields.
	if recvMsg.getMode() != server {
		return nil, 0, errors.New("invalid mode in response")
	}
	if recvMsg.TransmitTime == ntpTime(0) {
		return nil, 0, errors.New("invalid transmit time in response")
	}
	if recvMsg.OriginTime != xmitMsg.TransmitTime {
		return nil, 0, errors.New("server response mismatch")
	}
	if recvMsg.ReceiveTime > recvMsg.TransmitTime {
		return nil, 0, errors.New("server clock ticked backwards")
	}

	// Correct the received message's origin time using the actual
	// transmit time.
	recvMsg.OriginTime = toNtpTime(xmitTime)

	return recvMsg, recvTime, nil
}

// parseTime parses the NTP packet along with the packet receive time to
// generate a Response record.
func parseTime(m *msg, recvTime ntpTime) *Response {
	r := &Response{
		Time:           m.TransmitTime.Time(),
		ClockOffset:    offset(m.OriginTime, m.ReceiveTime, m.TransmitTime, recvTime),
		RTT:            rtt(m.OriginTime, m.ReceiveTime, m.TransmitTime, recvTime),
		Precision:      toInterval(m.Precision),
		Stratum:        m.Stratum,
		ReferenceID:    m.ReferenceID,
		ReferenceTime:  m.ReferenceTime.Time(),
		RootDelay:      m.RootDelay.Duration(),
		RootDispersion: m.RootDispersion.Duration(),
		Leap:           m.getLeap(),
		MinError:       minError(m.OriginTime, m.ReceiveTime, m.TransmitTime, recvTime),
		Poll:           toInterval(m.Poll),
	}

	// Calculate values depending on other calculated values
	r.RootDistance = rootDistance(r.RTT, r.RootDelay, r.RootDispersion)

	// If a kiss of death was received, interpret the reference ID as
	// a kiss code.
	if r.Stratum == 0 {
		r.KissCode = kissCode(r.ReferenceID)
	}

	return r
}

// The following helper functions calculate additional metadata about the
// timestamps received from an NTP server.  The timestamps returned by
// the server are given the following variable names:
//
//   org = Origin Timestamp (client send time)
//   rec = Receive Timestamp (server receive time)
//   xmt = Transmit Timestamp (server reply time)
//   dst = Destination Timestamp (client receive time)

func rtt(org, rec, xmt, dst ntpTime) time.Duration {
	// round trip delay time
	//   rtt = (dst-org) - (xmt-rec)
	a := dst.Time().Sub(org.Time())
	b := xmt.Time().Sub(rec.Time())
	rtt := a - b
	if rtt < 0 {
		rtt = 0
	}
	return rtt
}

func offset(org, rec, xmt, dst ntpTime) time.Duration {
	// local clock offset
	//   offset = ((rec-org) + (xmt-dst)) / 2
	a := rec.Time().Sub(org.Time())
	b := xmt.Time().Sub(dst.Time())
	return (a + b) / time.Duration(2)
}

func minError(org, rec, xmt, dst ntpTime) time.Duration {
	// Each NTP response contains two pairs of send/receive timestamps.
	// When either pair indicates a "causality violation", we calculate the
	// error as the difference in time between them. The minimum error is
	// the greater of the two causality violations.
	var error0, error1 ntpTime
	if org >= rec {
		error0 = org - rec
	}
	if xmt >= dst {
		error1 = xmt - dst
	}
	if error0 > error1 {
		return error0.Duration()
	}
	return error1.Duration()
}

func rootDistance(rtt, rootDelay, rootDisp time.Duration) time.Duration {
	// The root distance is:
	// 	the maximum error due to all causes of the local clock
	//	relative to the primary server. It is defined as half the
	//	total delay plus total dispersion plus peer jitter.
	//	(https://tools.ietf.org/html/rfc5905#appendix-A.5.5.2)
	//
	// In the reference implementation, it is calculated as follows:
	//	rootDist = max(MINDISP, rootDelay + rtt)/2 + rootDisp
	//			+ peerDisp + PHI * (uptime - peerUptime)
	//			+ peerJitter
	// For an SNTP client which sends only a single packet, most of these
	// terms are irrelevant and become 0.
	totalDelay := rtt + rootDelay
	return totalDelay/2 + rootDisp
}

func toInterval(t int8) time.Duration {
	switch {
	case t > 0:
		return time.Duration(uint64(time.Second) << uint(t))
	case t < 0:
		return time.Duration(uint64(time.Second) >> uint(-t))
	default:
		return time.Second
	}
}

func kissCode(id uint32) string {
	isPrintable := func(ch byte) bool { return ch >= 32 && ch <= 126 }

	b := []byte{
		byte(id >> 24),
		byte(id >> 16),
		byte(id >> 8),
		byte(id),
	}
	for _, ch := range b {
		if !isPrintable(ch) {
			return ""
		}
	}
	return string(b)
}
