// Copyright 2016 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build !nologind

package collector

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-kit/kit/log"
	"github.com/godbus/dbus"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	logindSubsystem = "logind"
	dbusObject      = "org.freedesktop.login1"
	dbusPath        = "/org/freedesktop/login1"
)

var (
	// Taken from logind as of systemd v229.
	// "other" is the fallback value for unknown values (in case logind gets extended in the future).
	attrRemoteValues = []string{"true", "false"}
	attrTypeValues   = []string{"other", "unspecified", "tty", "x11", "wayland", "mir", "web"}
	attrClassValues  = []string{"other", "user", "greeter", "lock-screen", "background"}

	sessionsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, logindSubsystem, "sessions"),
		"Number of sessions registered in logind.", []string{"seat", "remote", "type", "class"}, nil,
	)
)

type logindCollector struct {
	logger log.Logger
}

type logindDbus struct {
	conn   *dbus.Conn
	object dbus.BusObject
}

type logindInterface interface {
	listSeats() ([]string, error)
	listSessions() ([]logindSessionEntry, error)
	getSession(logindSessionEntry) *logindSession
}

type logindSession struct {
	seat        string
	remote      string
	sessionType string
	class       string
}

// Struct elements must be public for the reflection magic of godbus to work.
type logindSessionEntry struct {
	SessionID         string
	UserID            uint32
	UserName          string
	SeatID            string
	SessionObjectPath dbus.ObjectPath
}

type logindSeatEntry struct {
	SeatID         string
	SeatObjectPath dbus.ObjectPath
}

func init() {
	registerCollector("logind", defaultDisabled, NewLogindCollector)
}

// NewLogindCollector returns a new Collector exposing logind statistics.
func NewLogindCollector(logger log.Logger) (Collector, error) {
	return &logindCollector{logger}, nil
}

func (lc *logindCollector) Update(ch chan<- prometheus.Metric) error {
	c, err := newDbus()
	if err != nil {
		return fmt.Errorf("unable to connect to dbus: %w", err)
	}
	defer c.conn.Close()

	return collectMetrics(ch, c)
}

func collectMetrics(ch chan<- prometheus.Metric, c logindInterface) error {
	seats, err := c.listSeats()
	if err != nil {
		return fmt.Errorf("unable to get seats: %w", err)
	}

	sessionList, err := c.listSessions()
	if err != nil {
		return fmt.Errorf("unable to get sessions: %w", err)
	}

	sessions := make(map[logindSession]float64)

	for _, s := range sessionList {
		session := c.getSession(s)
		if session != nil {
			sessions[*session]++
		}
	}

	for _, remote := range attrRemoteValues {
		for _, sessionType := range attrTypeValues {
			for _, class := range attrClassValues {
				for _, seat := range seats {
					count := sessions[logindSession{seat, remote, sessionType, class}]

					ch <- prometheus.MustNewConstMetric(
						sessionsDesc, prometheus.GaugeValue, count,
						seat, remote, sessionType, class)
				}
			}
		}
	}

	return nil
}

func knownStringOrOther(value string, known []string) string {
	for i := range known {
		if value == known[i] {
			return value
		}
	}

	return "other"
}

func newDbus() (*logindDbus, error) {
	conn, err := dbus.SystemBusPrivate()
	if err != nil {
		return nil, err
	}

	methods := []dbus.Auth{dbus.AuthExternal(strconv.Itoa(os.Getuid()))}

	err = conn.Auth(methods)
	if err != nil {
		conn.Close()
		return nil, err
	}

	err = conn.Hello()
	if err != nil {
		conn.Close()
		return nil, err
	}

	object := conn.Object(dbusObject, dbus.ObjectPath(dbusPath))

	return &logindDbus{
		conn:   conn,
		object: object,
	}, nil
}

func (c *logindDbus) listSeats() ([]string, error) {
	var result [][]interface{}
	err := c.object.Call(dbusObject+".Manager.ListSeats", 0).Store(&result)
	if err != nil {
		return nil, err
	}

	resultInterface := make([]interface{}, len(result))
	for i := range result {
		resultInterface[i] = result[i]
	}

	seats := make([]logindSeatEntry, len(result))
	seatsInterface := make([]interface{}, len(seats))
	for i := range seats {
		seatsInterface[i] = &seats[i]
	}

	err = dbus.Store(resultInterface, seatsInterface...)
	if err != nil {
		return nil, err
	}

	ret := make([]string, len(seats)+1)
	for i := range seats {
		ret[i] = seats[i].SeatID
	}
	// Always add the empty seat, which is used for remote sessions like SSH
	ret[len(seats)] = ""

	return ret, nil
}

func (c *logindDbus) listSessions() ([]logindSessionEntry, error) {
	var result [][]interface{}
	err := c.object.Call(dbusObject+".Manager.ListSessions", 0).Store(&result)
	if err != nil {
		return nil, err
	}

	resultInterface := make([]interface{}, len(result))
	for i := range result {
		resultInterface[i] = result[i]
	}

	sessions := make([]logindSessionEntry, len(result))
	sessionsInterface := make([]interface{}, len(sessions))
	for i := range sessions {
		sessionsInterface[i] = &sessions[i]
	}

	err = dbus.Store(resultInterface, sessionsInterface...)
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (c *logindDbus) getSession(session logindSessionEntry) *logindSession {
	object := c.conn.Object(dbusObject, session.SessionObjectPath)

	remote, err := object.GetProperty(dbusObject + ".Session.Remote")
	if err != nil {
		return nil
	}

	sessionType, err := object.GetProperty(dbusObject + ".Session.Type")
	if err != nil {
		return nil
	}

	sessionTypeStr, ok := sessionType.Value().(string)
	if !ok {
		return nil
	}

	class, err := object.GetProperty(dbusObject + ".Session.Class")
	if err != nil {
		return nil
	}

	classStr, ok := class.Value().(string)
	if !ok {
		return nil
	}

	return &logindSession{
		seat:        session.SeatID,
		remote:      remote.String(),
		sessionType: knownStringOrOther(sessionTypeStr, attrTypeValues),
		class:       knownStringOrOther(classStr, attrClassValues),
	}
}
