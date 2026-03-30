// Copyright 2026 The Prometheus Authors
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

//go:build !nocifs

package collector

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

const cifsSubsystem = "cifs"

type cifsCollector struct {
	logger *slog.Logger

	sessionsDesc          *prometheus.Desc
	sharesDesc            *prometheus.Desc
	sessionReconnectsDesc *prometheus.Desc
	shareReconnectsDesc   *prometheus.Desc
	vfsOpsDesc            *prometheus.Desc
	smbsDesc              *prometheus.Desc
	readBytesDesc         *prometheus.Desc
	writeBytesDesc        *prometheus.Desc
}

func init() {
	registerCollector("cifs", defaultDisabled, NewCIFSCollector)
}

// NewCIFSCollector returns a new Collector exposing CIFS client statistics.
func NewCIFSCollector(logger *slog.Logger) (Collector, error) {
	return &cifsCollector{
		logger: logger,
		sessionsDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cifsSubsystem, "sessions"),
			"Number of active CIFS sessions.",
			nil, nil,
		),
		sharesDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cifsSubsystem, "shares"),
			"Number of unique mount targets.",
			nil, nil,
		),
		sessionReconnectsDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cifsSubsystem, "session_reconnects_total"),
			"Total number of CIFS session reconnects.",
			nil, nil,
		),
		shareReconnectsDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cifsSubsystem, "share_reconnects_total"),
			"Total number of CIFS share reconnects.",
			nil, nil,
		),
		vfsOpsDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cifsSubsystem, "vfs_operations_total"),
			"Total number of VFS operations.",
			nil, nil,
		),
		smbsDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cifsSubsystem, "smbs_total"),
			"Total number of SMBs sent for this share.",
			[]string{"share"}, nil,
		),
		readBytesDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cifsSubsystem, "read_bytes_total"),
			"Total bytes read from this share.",
			[]string{"share"}, nil,
		),
		writeBytesDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cifsSubsystem, "write_bytes_total"),
			"Total bytes written to this share.",
			[]string{"share"}, nil,
		),
	}, nil
}

func (c *cifsCollector) Update(ch chan<- prometheus.Metric) error {
	stats, err := parseCIFSStats(procFilePath("fs/cifs/Stats"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			c.logger.Debug("Not collecting CIFS metrics", "err", err)
			return ErrNoData
		}
		return fmt.Errorf("failed to read CIFS stats: %w", err)
	}

	ch <- prometheus.MustNewConstMetric(c.sessionsDesc, prometheus.GaugeValue, stats.sessions)
	ch <- prometheus.MustNewConstMetric(c.sharesDesc, prometheus.GaugeValue, stats.shares)
	ch <- prometheus.MustNewConstMetric(c.sessionReconnectsDesc, prometheus.CounterValue, stats.sessionReconnects)
	ch <- prometheus.MustNewConstMetric(c.shareReconnectsDesc, prometheus.CounterValue, stats.shareReconnects)
	ch <- prometheus.MustNewConstMetric(c.vfsOpsDesc, prometheus.CounterValue, stats.vfsOps)

	for _, share := range stats.perShare {
		ch <- prometheus.MustNewConstMetric(c.smbsDesc, prometheus.CounterValue, share.smbs, share.name)
		ch <- prometheus.MustNewConstMetric(c.readBytesDesc, prometheus.CounterValue, share.readBytes, share.name)
		ch <- prometheus.MustNewConstMetric(c.writeBytesDesc, prometheus.CounterValue, share.writeBytes, share.name)
	}

	return nil
}

type cifsStats struct {
	sessions          float64
	shares            float64
	sessionReconnects float64
	shareReconnects   float64
	vfsOps            float64
	perShare          []cifsShareStats
}

type cifsShareStats struct {
	name       string
	smbs       float64
	readBytes  float64
	writeBytes float64
}

func parseCIFSStats(path string) (*cifsStats, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	stats := &cifsStats{}
	scanner := bufio.NewScanner(f)

	var currentShare string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		fields := strings.Fields(line)

		if strings.HasPrefix(line, "CIFS Session:") && len(fields) >= 3 {
			v, err := strconv.ParseFloat(fields[2], 64)
			if err == nil {
				stats.sessions = v
			}
			continue
		}

		if strings.HasPrefix(line, "Share (unique mount targets):") && len(fields) >= 5 {
			v, err := strconv.ParseFloat(fields[4], 64)
			if err == nil {
				stats.shares = v
			}
			continue
		}

		if strings.HasSuffix(line, "share reconnects") && len(fields) >= 4 {
			v, err := strconv.ParseFloat(fields[0], 64)
			if err == nil {
				stats.sessionReconnects = v
			}
			v, err = strconv.ParseFloat(fields[2], 64)
			if err == nil {
				stats.shareReconnects = v
			}
			continue
		}

		if strings.HasPrefix(line, "Total vfs operations:") && len(fields) >= 4 {
			v, err := strconv.ParseFloat(fields[3], 64)
			if err == nil {
				stats.vfsOps = v
			}
			continue
		}

		if len(fields) >= 2 && strings.HasSuffix(fields[0], ")") {
			_, err := strconv.Atoi(strings.TrimSuffix(fields[0], ")"))
			if err == nil {
				currentShare = fields[1]
				stats.perShare = append(stats.perShare, cifsShareStats{name: currentShare})
			}
			continue
		}

		if currentShare == "" {
			continue
		}
		idx := len(stats.perShare) - 1

		if strings.HasPrefix(line, "SMBs:") && len(fields) >= 2 {
			v, err := strconv.ParseFloat(fields[1], 64)
			if err == nil {
				stats.perShare[idx].smbs = v
			}
			continue
		}

		if strings.HasPrefix(line, "Bytes read:") && len(fields) >= 3 {
			v, err := strconv.ParseFloat(fields[2], 64)
			if err == nil {
				stats.perShare[idx].readBytes = v
			}
			continue
		}

		if strings.HasPrefix(line, "Bytes written:") && len(fields) >= 3 {
			v, err := strconv.ParseFloat(fields[2], 64)
			if err == nil {
				stats.perShare[idx].writeBytes = v
			}
			continue
		}
	}

	return stats, scanner.Err()
}
