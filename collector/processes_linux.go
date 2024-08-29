// Copyright 2018 The Prometheus Authors
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

//go:build !noprocesses
// +build !noprocesses

package collector

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"
	"strconv"
	"strings"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

type processCollector struct {
	fs           procfs.FS
	threadAlloc  *prometheus.Desc
	threadLimit  *prometheus.Desc
	threadsState *prometheus.Desc
	procsState   *prometheus.Desc
	pidUsed      *prometheus.Desc
	pidMax       *prometheus.Desc
	logger       *slog.Logger
}

func init() {
	registerCollector("processes", defaultDisabled, NewProcessStatCollector)
}

// NewProcessStatCollector returns a new Collector exposing process data read from the proc filesystem.
func NewProcessStatCollector(logger *slog.Logger) (Collector, error) {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}
	subsystem := "processes"
	return &processCollector{
		fs: fs,
		threadAlloc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "threads"),
			"Allocated threads in system",
			nil, nil,
		),
		threadLimit: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "max_threads"),
			"Limit of threads in the system",
			nil, nil,
		),
		threadsState: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "threads_state"),
			"Number of threads in each state.",
			[]string{"thread_state"}, nil,
		),
		procsState: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "state"),
			"Number of processes in each state.",
			[]string{"state"}, nil,
		),
		pidUsed: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "pids"),
			"Number of PIDs", nil, nil,
		),
		pidMax: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "max_processes"),
			"Number of max PIDs limit", nil, nil,
		),
		logger: logger,
	}, nil
}
func (c *processCollector) Update(ch chan<- prometheus.Metric) error {
	pids, states, threads, threadStates, err := c.getAllocatedThreads()
	if err != nil {
		return fmt.Errorf("unable to retrieve number of allocated threads: %w", err)
	}

	ch <- prometheus.MustNewConstMetric(c.threadAlloc, prometheus.GaugeValue, float64(threads))
	maxThreads, err := readUintFromFile(procFilePath("sys/kernel/threads-max"))
	if err != nil {
		return fmt.Errorf("unable to retrieve limit number of threads: %w", err)
	}
	ch <- prometheus.MustNewConstMetric(c.threadLimit, prometheus.GaugeValue, float64(maxThreads))

	for state := range states {
		ch <- prometheus.MustNewConstMetric(c.procsState, prometheus.GaugeValue, float64(states[state]), state)
	}

	for state := range threadStates {
		ch <- prometheus.MustNewConstMetric(c.threadsState, prometheus.GaugeValue, float64(threadStates[state]), state)
	}

	pidM, err := readUintFromFile(procFilePath("sys/kernel/pid_max"))
	if err != nil {
		return fmt.Errorf("unable to retrieve limit number of maximum pids alloved: %w", err)
	}
	ch <- prometheus.MustNewConstMetric(c.pidUsed, prometheus.GaugeValue, float64(pids))
	ch <- prometheus.MustNewConstMetric(c.pidMax, prometheus.GaugeValue, float64(pidM))

	return nil
}

func (c *processCollector) getAllocatedThreads() (int, map[string]int32, int, map[string]int32, error) {
	p, err := c.fs.AllProcs()
	if err != nil {
		return 0, nil, 0, nil, fmt.Errorf("unable to list all processes: %w", err)
	}
	pids := 0
	thread := 0
	procStates := make(map[string]int32)
	threadStates := make(map[string]int32)

	for _, pid := range p {
		stat, err := pid.Stat()
		if err != nil {
			// PIDs can vanish between getting the list and getting stats.
			if c.isIgnoredError(err) {
				c.logger.Debug("file not found when retrieving stats for pid", "pid", pid.PID, "err", err)
				continue
			}
			c.logger.Debug("error reading stat for pid", "pid", pid.PID, "err", err)
			return 0, nil, 0, nil, fmt.Errorf("error reading stat for pid %d: %w", pid.PID, err)
		}
		pids++
		procStates[stat.State]++
		thread += stat.NumThreads
		err = c.getThreadStates(pid.PID, stat, threadStates)
		if err != nil {
			return 0, nil, 0, nil, err
		}
	}
	return pids, procStates, thread, threadStates, nil
}

func (c *processCollector) getThreadStates(pid int, pidStat procfs.ProcStat, threadStates map[string]int32) error {
	fs, err := procfs.NewFS(procFilePath(path.Join(strconv.Itoa(pid), "task")))
	if err != nil {
		if c.isIgnoredError(err) {
			c.logger.Debug("file not found when retrieving tasks for pid", "pid", pid, "err", err)
			return nil
		}
		c.logger.Debug("error reading tasks for pid", "pid", pid, "err", err)
		return fmt.Errorf("error reading task for pid %d: %w", pid, err)
	}

	t, err := fs.AllProcs()
	if err != nil {
		if c.isIgnoredError(err) {
			c.logger.Debug("file not found when retrieving tasks for pid", "pid", pid, "err", err)
			return nil
		}
		return fmt.Errorf("unable to list all threads for pid: %d %w", pid, err)
	}

	for _, thread := range t {
		if pid == thread.PID {
			threadStates[pidStat.State]++
			continue
		}
		threadStat, err := thread.Stat()
		if err != nil {
			if c.isIgnoredError(err) {
				c.logger.Debug("file not found when retrieving stats for thread", "pid", pid, "threadId", thread.PID, "err", err)
				continue
			}
			c.logger.Debug("error reading stat for thread", "pid", pid, "threadId", thread.PID, "err", err)
			return fmt.Errorf("error reading stat for pid:%d thread:%d err:%w", pid, thread.PID, err)
		}
		threadStates[threadStat.State]++
	}
	return nil
}

func (c *processCollector) isIgnoredError(err error) bool {
	if errors.Is(err, os.ErrNotExist) || strings.Contains(err.Error(), syscall.ESRCH.Error()) {
		return true
	}
	return false
}
