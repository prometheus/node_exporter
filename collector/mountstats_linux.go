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

package collector

import (
	"flag"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

type mountStatsCollector struct {
	// General statistics
	NFSAgeSecondsTotal *prometheus.Desc

	// Byte statistics
	NFSReadBytesTotal        *prometheus.Desc
	NFSWriteBytesTotal       *prometheus.Desc
	NFSDirectReadBytesTotal  *prometheus.Desc
	NFSDirectWriteBytesTotal *prometheus.Desc
	NFSTotalReadBytesTotal   *prometheus.Desc
	NFSTotalWriteBytesTotal  *prometheus.Desc
	NFSReadPagesTotal        *prometheus.Desc
	NFSWritePagesTotal       *prometheus.Desc

	// Per-operation statistics
	NFSOperationsRequestsTotal            *prometheus.Desc
	NFSOperationsTransmissionsTotal       *prometheus.Desc
	NFSOperationsMajorTimeoutsTotal       *prometheus.Desc
	NFSOperationsSentBytesTotal           *prometheus.Desc
	NFSOperationsReceivedBytesTotal       *prometheus.Desc
	NFSOperationsQueueTimeSecondsTotal    *prometheus.Desc
	NFSOperationsResponseTimeSecondsTotal *prometheus.Desc
	NFSOperationsRequestTimeSecondsTotal  *prometheus.Desc

	// Transport statistics
	NFSTransportBindTotal              *prometheus.Desc
	NFSTransportConnectTotal           *prometheus.Desc
	NFSTransportIdleTimeSeconds        *prometheus.Desc
	NFSTransportSendsTotal             *prometheus.Desc
	NFSTransportReceivesTotal          *prometheus.Desc
	NFSTransportBadTransactionIDsTotal *prometheus.Desc
	NFSTransportBacklogQueueTotal      *prometheus.Desc
	NFSTransportMaximumRPCSlots        *prometheus.Desc
	NFSTransportSendingQueueTotal      *prometheus.Desc
	NFSTransportPendingQueueTotal      *prometheus.Desc

	proc procfs.Proc
}

func init() {
	Factories["mountstats"] = NewMountStatsCollector
	CollectorsEnabledState["mountstats"] = flag.Bool("collectors.mountstats.enabled", false, "enable mountstats-collectors")
}

func NewMountStatsCollector() (Collector, error) {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %v", err)
	}

	proc, err := fs.Self()
	if err != nil {
		return nil, fmt.Errorf("failed to open /proc/self: %v", err)
	}

	const (
		// For the time being, only NFS statistics are available via this mechanism
		subsystem = "mountstats_nfs"
	)

	var (
		labels   = []string{"export"}
		opLabels = []string{"export", "operation"}
	)

	return &mountStatsCollector{
		NFSAgeSecondsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "age_seconds_total"),
			"The age of the NFS mount in seconds.",
			labels,
			nil,
		),

		NFSReadBytesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "read_bytes_total"),
			"Number of bytes read using the read() syscall.",
			labels,
			nil,
		),

		NFSWriteBytesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "write_bytes_total"),
			"Number of bytes written using the write() syscall.",
			labels,
			nil,
		),

		NFSDirectReadBytesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "direct_read_bytes_total"),
			"Number of bytes read using the read() syscall in O_DIRECT mode.",
			labels,
			nil,
		),

		NFSDirectWriteBytesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "direct_write_bytes_total"),
			"Number of bytes written using the write() syscall in O_DIRECT mode.",
			labels,
			nil,
		),

		NFSTotalReadBytesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "total_read_bytes_total"),
			"Number of bytes read from the NFS server, in total.",
			labels,
			nil,
		),

		NFSTotalWriteBytesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "total_write_bytes_total"),
			"Number of bytes written to the NFS server, in total.",
			labels,
			nil,
		),

		NFSReadPagesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "read_pages_total"),
			"Number of pages read directly via mmap()'d files.",
			labels,
			nil,
		),

		NFSWritePagesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "write_pages_total"),
			"Number of pages written directly via mmap()'d files.",
			labels,
			nil,
		),

		NFSTransportBindTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "transport_bind_total"),
			"Number of times the client has had to establish a connection from scratch to the NFS server.",
			labels,
			nil,
		),

		NFSTransportConnectTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "transport_connect_total"),
			"Number of times the client has made a TCP connection to the NFS server.",
			labels,
			nil,
		),

		NFSTransportIdleTimeSeconds: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "transport_idle_time_seconds"),
			"Duration since the NFS mount last saw any RPC traffic, in seconds.",
			labels,
			nil,
		),

		NFSTransportSendsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "transport_sends_total"),
			"Number of RPC requests for this mount sent to the NFS server.",
			labels,
			nil,
		),

		NFSTransportReceivesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "transport_receives_total"),
			"Number of RPC responses for this mount received from the NFS server.",
			labels,
			nil,
		),

		NFSTransportBadTransactionIDsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "transport_bad_transaction_ids_total"),
			"Number of times the NFS server sent a response with a transaction ID unknown to this client.",
			labels,
			nil,
		),

		NFSTransportBacklogQueueTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "transport_backlog_queue_total"),
			"Total number of items added to the RPC backlog queue.",
			labels,
			nil,
		),

		NFSTransportMaximumRPCSlots: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "transport_maximum_rpc_slots"),
			"Maximum number of simultaneously active RPC requests ever used.",
			labels,
			nil,
		),

		NFSTransportSendingQueueTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "transport_sending_queue_total"),
			"Total number of items added to the RPC transmission sending queue.",
			labels,
			nil,
		),

		NFSTransportPendingQueueTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "transport_pending_queue_total"),
			"Total number of items added to the RPC transmission pending queue.",
			labels,
			nil,
		),

		NFSOperationsRequestsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "operations_requests_total"),
			"Number of requests performed for a given operation.",
			opLabels,
			nil,
		),

		NFSOperationsTransmissionsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "operations_transmissions_total"),
			"Number of times an actual RPC request has been transmitted for a given operation.",
			opLabels,
			nil,
		),

		NFSOperationsMajorTimeoutsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "operations_major_timeouts_total"),
			"Number of times a request has had a major timeout for a given operation.",
			opLabels,
			nil,
		),

		NFSOperationsSentBytesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "operations_sent_bytes_total"),
			"Number of bytes sent for a given operation, including RPC headers and payload.",
			opLabels,
			nil,
		),

		NFSOperationsReceivedBytesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "operations_received_bytes_total"),
			"Number of bytes received for a given operation, including RPC headers and payload.",
			opLabels,
			nil,
		),

		NFSOperationsQueueTimeSecondsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "operations_queue_time_seconds_total"),
			"Duration all requests spent queued for transmission for a given operation before they were sent, in seconds.",
			opLabels,
			nil,
		),

		NFSOperationsResponseTimeSecondsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "operations_response_time_seconds_total"),
			"Duration all requests took to get a reply back after a request for a given operation was transmitted, in seconds.",
			opLabels,
			nil,
		),

		NFSOperationsRequestTimeSecondsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "operations_request_time_seconds_total"),
			"Duration all requests took from when a request was enqueued to when it was completely handled for a given operation, in seconds.",
			opLabels,
			nil,
		),

		proc: proc,
	}, nil
}

func (c *mountStatsCollector) Update(ch chan<- prometheus.Metric) error {
	mounts, err := c.proc.MountStats()
	if err != nil {
		return fmt.Errorf("failed to parse mountstats: %v", err)
	}

	for _, m := range mounts {
		// For the time being, only NFS statistics are available via this mechanism
		stats, ok := m.Stats.(*procfs.MountStatsNFS)
		if !ok {
			continue
		}

		c.updateNFSStats(ch, m.Device, stats)
	}

	return nil
}

func (c *mountStatsCollector) updateNFSStats(ch chan<- prometheus.Metric, export string, s *procfs.MountStatsNFS) {
	ch <- prometheus.MustNewConstMetric(
		c.NFSAgeSecondsTotal,
		prometheus.CounterValue,
		s.Age.Seconds(),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSReadBytesTotal,
		prometheus.CounterValue,
		float64(s.Bytes.Read),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSWriteBytesTotal,
		prometheus.CounterValue,
		float64(s.Bytes.Write),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSDirectReadBytesTotal,
		prometheus.CounterValue,
		float64(s.Bytes.DirectRead),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSDirectWriteBytesTotal,
		prometheus.CounterValue,
		float64(s.Bytes.DirectWrite),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSTotalReadBytesTotal,
		prometheus.CounterValue,
		float64(s.Bytes.ReadTotal),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSTotalWriteBytesTotal,
		prometheus.CounterValue,
		float64(s.Bytes.WriteTotal),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSReadPagesTotal,
		prometheus.CounterValue,
		float64(s.Bytes.ReadPages),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSWritePagesTotal,
		prometheus.CounterValue,
		float64(s.Bytes.WritePages),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSTransportBindTotal,
		prometheus.CounterValue,
		float64(s.Transport.Bind),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSTransportConnectTotal,
		prometheus.CounterValue,
		float64(s.Transport.Connect),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSTransportIdleTimeSeconds,
		prometheus.GaugeValue,
		s.Transport.IdleTime.Seconds(),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSTransportSendsTotal,
		prometheus.CounterValue,
		float64(s.Transport.Sends),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSTransportReceivesTotal,
		prometheus.CounterValue,
		float64(s.Transport.Receives),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSTransportBadTransactionIDsTotal,
		prometheus.CounterValue,
		float64(s.Transport.BadTransactionIDs),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSTransportBacklogQueueTotal,
		prometheus.CounterValue,
		float64(s.Transport.CumulativeBacklog),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSTransportMaximumRPCSlots,
		prometheus.GaugeValue,
		float64(s.Transport.MaximumRPCSlotsUsed),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSTransportSendingQueueTotal,
		prometheus.CounterValue,
		float64(s.Transport.CumulativeSendingQueue),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSTransportPendingQueueTotal,
		prometheus.CounterValue,
		float64(s.Transport.CumulativePendingQueue),
		export,
	)

	for _, op := range s.Operations {
		ch <- prometheus.MustNewConstMetric(
			c.NFSOperationsRequestsTotal,
			prometheus.CounterValue,
			float64(op.Requests),
			export,
			op.Operation,
		)

		ch <- prometheus.MustNewConstMetric(
			c.NFSOperationsTransmissionsTotal,
			prometheus.CounterValue,
			float64(op.Transmissions),
			export,
			op.Operation,
		)

		ch <- prometheus.MustNewConstMetric(
			c.NFSOperationsMajorTimeoutsTotal,
			prometheus.CounterValue,
			float64(op.MajorTimeouts),
			export,
			op.Operation,
		)

		ch <- prometheus.MustNewConstMetric(
			c.NFSOperationsSentBytesTotal,
			prometheus.CounterValue,
			float64(op.BytesSent),
			export,
			op.Operation,
		)

		ch <- prometheus.MustNewConstMetric(
			c.NFSOperationsReceivedBytesTotal,
			prometheus.CounterValue,
			float64(op.BytesReceived),
			export,
			op.Operation,
		)

		ch <- prometheus.MustNewConstMetric(
			c.NFSOperationsQueueTimeSecondsTotal,
			prometheus.CounterValue,
			op.CumulativeQueueTime.Seconds(),
			export,
			op.Operation,
		)

		ch <- prometheus.MustNewConstMetric(
			c.NFSOperationsResponseTimeSecondsTotal,
			prometheus.CounterValue,
			op.CumulativeTotalResponseTime.Seconds(),
			export,
			op.Operation,
		)

		ch <- prometheus.MustNewConstMetric(
			c.NFSOperationsRequestTimeSecondsTotal,
			prometheus.CounterValue,
			op.CumulativeTotalRequestTime.Seconds(),
			export,
			op.Operation,
		)
	}
}
