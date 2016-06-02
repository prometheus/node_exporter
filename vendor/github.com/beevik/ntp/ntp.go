// Copyright 2015 Brett Vickers.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ntp provides a simple mechanism for querying the current time from
// a remote NTP server.  This package only supports NTP client mode behavior
// and version 4 of the NTP protocol.  See RFC 5905. Approach inspired by go-
// nuts post by Michael Hofmann:
//
// https://groups.google.com/forum/?fromgroups#!topic/golang-nuts/FlcdMU5fkLQ
package ntp

import (
	"encoding/binary"
	"net"
	"time"
)

type mode byte

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

const (
	maxStratum = 16
	nanoPerSec = 1000000000
)

var (
	timeout  = 5 * time.Second
	ntpEpoch = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
)

type ntpTime struct {
	Seconds  uint32
	Fraction uint32
}

func (t ntpTime) Time() time.Time {
	return ntpEpoch.Add(t.sinceEpoch())
}

// sinceEpoch converts the ntpTime record t into a duration since the NTP
// epoch time (Jan 1, 1900).
func (t ntpTime) sinceEpoch() time.Duration {
	sec := time.Duration(t.Seconds) * time.Second
	frac := time.Duration(uint64(t.Fraction) * nanoPerSec >> 32)
	return sec + frac
}

// toNtpTime converts the time value t into an ntpTime representation.
func toNtpTime(t time.Time) ntpTime {
	nsec := uint64(t.Sub(ntpEpoch))
	return ntpTime{
		Seconds:  uint32(nsec / nanoPerSec),
		Fraction: uint32((nsec % nanoPerSec) << 32 / nanoPerSec),
	}
}

// msg is an internal representation of an NTP packet.
type msg struct {
	LiVnMode       byte // Leap Indicator (2) + Version (3) + Mode (3)
	Stratum        byte
	Poll           byte
	Precision      byte
	RootDelay      uint32
	RootDispersion uint32
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
	m.LiVnMode = (m.LiVnMode & 0xf8) | byte(md)
}

// A Response contains time data, some of which is returned by the NTP server
// and some of which is calculated by the client.
type Response struct {
	Time        time.Time     // receive time reported by the server
	RTT         time.Duration // round-trip time between client and server
	ClockOffset time.Duration // local clock offset relative to server
	Stratum     uint8         // stratum level of NTP server's clock
}

// Query returns information from the remote NTP server specifed as host.  NTP
// client mode is used.
func Query(host string, version int) (*Response, error) {
	m, err := getTime(host, version)
	now := toNtpTime(time.Now())
	if err != nil {
		return nil, err
	}

	r := &Response{
		Time:        m.ReceiveTime.Time(),
		RTT:         rtt(m.OriginTime, m.ReceiveTime, m.TransmitTime, now),
		ClockOffset: offset(m.OriginTime, m.ReceiveTime, m.TransmitTime, now),
		Stratum:     m.Stratum,
	}

	// https://tools.ietf.org/html/rfc5905#section-7.3
	if r.Stratum == 0 {
		r.Stratum = maxStratum
	}

	return r, nil
}

// Time returns the "receive time" from the remote NTP server specifed as
// host.  NTP client mode is used.
func getTime(host string, version int) (*msg, error) {
	if version < 2 || version > 4 {
		panic("ntp: invalid version number")
	}

	raddr, err := net.ResolveUDPAddr("udp", host+":123")
	if err != nil {
		return nil, err
	}

	con, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return nil, err
	}
	defer con.Close()
	con.SetDeadline(time.Now().Add(timeout))

	m := new(msg)
	m.setMode(client)
	m.setVersion(version)
	m.TransmitTime = toNtpTime(time.Now())

	err = binary.Write(con, binary.BigEndian, m)
	if err != nil {
		return nil, err
	}

	err = binary.Read(con, binary.BigEndian, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// TimeV returns the "receive time" from the remote NTP server specifed as
// host.  Use the NTP client mode with the requested version number (2, 3, or
// 4).
func TimeV(host string, version int) (time.Time, error) {
	m, err := getTime(host, version)
	if err != nil {
		return time.Now(), err
	}
	return m.ReceiveTime.Time().Local(), nil
}

// Time returns the "receive time" from the remote NTP server specifed as
// host.  NTP client mode version 4 is used.
func Time(host string) (time.Time, error) {
	return TimeV(host, 4)
}

func rtt(t1, t2, t3, t4 ntpTime) time.Duration {
	// round trip delay time (https://tools.ietf.org/html/rfc5905#section-8)
	//   T1 = client send time
	//   T2 = server receive time
	//   T3 = server reply time
	//   T4 = client receive time
	//
	// RTT d:
	//   d = (T4-T1) - (T3-T2)
	a := t4.Time().Sub(t1.Time())
	b := t3.Time().Sub(t2.Time())
	return a - b
}

func offset(t1, t2, t3, t4 ntpTime) time.Duration {
	// local offset equation (https://tools.ietf.org/html/rfc5905#section-8)
	//   T1 = client send time
	//   T2 = server receive time
	//   T3 = server reply time
	//   T4 = client receive time
	//
	// Local clock offset t:
	//   t = ((T2-T1) + (T3-T4)) / 2
	a := t2.Time().Sub(t1.Time())
	b := t3.Time().Sub(t4.Time())
	return (a + b) / time.Duration(2)
}
