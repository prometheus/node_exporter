// Copyright 2015 Brett Vickers.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ntp provides a simple mechanism for querying the current time from
// a remote NTP server. See RFC 5905. Approach inspired by go-nuts post by
// Michael Hofmann:
//
// https://groups.google.com/forum/?fromgroups#!topic/golang-nuts/FlcdMU5fkLQ
package ntp

import (
	"encoding/binary"
	"net"
	"time"
)

type mode uint8

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

// An ntpTime is a 64-bit fixed-point (Q32.32) representation of the number of
// seconds elapsed since the NTP epoch.
type ntpTime uint64

// Duration interprets the fixed-point ntpTime as a number of elapsed seconds
// and returns the corresponding time.Duration value.
func (t ntpTime) Duration() time.Duration {
	sec := (t >> 32) * nanoPerSec
	frac := (t & 0xffffffff) * nanoPerSec >> 32
	return time.Duration(sec + frac)
}

// Time interprets the fixed-point ntpTime as a an absolute time and returns
// the corresponding time.Time value.
func (t ntpTime) Time() time.Time {
	return ntpEpoch.Add(t.Duration())
}

// toNtpTime converts the time.Time value t into its 64-bit fixed-point
// ntpTime representation.
func toNtpTime(t time.Time) ntpTime {
	nsec := uint64(t.Sub(ntpEpoch))
	sec := nsec / nanoPerSec
	frac := (nsec - sec*nanoPerSec) << 32 / nanoPerSec
	return ntpTime(sec<<32 | frac)
}

// An ntpTimeShort is a 32-bit fixed-point (Q16.16) representation of the
// number of seconds elapsed since the NTP epoch.
type ntpTimeShort uint32

// Duration interprets the fixed-point ntpTimeShort as a number of elapsed
// seconds and returns the corresponding time.Duration value.
func (t ntpTimeShort) Duration() time.Duration {
	sec := (t >> 16) * nanoPerSec
	frac := (t & 0xffff) * nanoPerSec >> 16
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

// A Response contains time data, some of which is returned by the NTP server
// and some of which is calculated by the client.
type Response struct {
	Time           time.Time     // receive time reported by the server
	RTT            time.Duration // round-trip time between client and server
	ClockOffset    time.Duration // local clock offset relative to server
	Poll           time.Duration // maximum polling interval
	Precision      time.Duration // precision of server's system clock
	Stratum        uint8         // stratum level of NTP server's clock
	ReferenceID    uint32        // server's reference ID
	RootDelay      time.Duration // server's RTT to the reference clock
	RootDispersion time.Duration // server's dispersion to the reference clock
}

// Query returns the current time from the remote server host using the
// requested version of the NTP protocol. It also returns additional
// information about the exchanged time information. The version may be 2, 3,
// or 4; although 4 is most typically used.
func Query(host string, version int) (*Response, error) {
	m, err := getTime(host, version)
	now := toNtpTime(time.Now())
	if err != nil {
		return nil, err
	}

	r := &Response{
		Time:           m.ReceiveTime.Time(),
		RTT:            rtt(m.OriginTime, m.ReceiveTime, m.TransmitTime, now),
		ClockOffset:    offset(m.OriginTime, m.ReceiveTime, m.TransmitTime, now),
		Poll:           toInterval(m.Poll),
		Precision:      toInterval(m.Precision),
		Stratum:        m.Stratum,
		ReferenceID:    m.ReferenceID,
		RootDelay:      m.RootDelay.Duration(),
		RootDispersion: m.RootDispersion.Duration(),
	}

	// https://tools.ietf.org/html/rfc5905#section-7.3
	if r.Stratum == 0 {
		r.Stratum = maxStratum
	}

	return r, nil
}

// getTime returns the "receive time" from the remote NTP server host.
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

// TimeV returns the current time from the remote server host using the
// requested version of the NTP protocol. The version may be 2, 3, or 4;
// although 4 is most typically used.
func TimeV(host string, version int) (time.Time, error) {
	m, err := getTime(host, version)
	if err != nil {
		return time.Now(), err
	}
	return m.ReceiveTime.Time().Local(), nil
}

// Time returns the current time from the remote server host using version 4
// of the NTP protocol.
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
