// Copyright 2017 The Prometheus Authors
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
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
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

	// Event statistics
	NFSEventInodeRevalidateTotal     *prometheus.Desc
	NFSEventDnodeRevalidateTotal     *prometheus.Desc
	NFSEventDataInvalidateTotal      *prometheus.Desc
	NFSEventAttributeInvalidateTotal *prometheus.Desc
	NFSEventVFSOpenTotal             *prometheus.Desc
	NFSEventVFSLookupTotal           *prometheus.Desc
	NFSEventVFSAccessTotal           *prometheus.Desc
	NFSEventVFSUpdatePageTotal       *prometheus.Desc
	NFSEventVFSReadPageTotal         *prometheus.Desc
	NFSEventVFSReadPagesTotal        *prometheus.Desc
	NFSEventVFSWritePageTotal        *prometheus.Desc
	NFSEventVFSWritePagesTotal       *prometheus.Desc
	NFSEventVFSGetdentsTotal         *prometheus.Desc
	NFSEventVFSSetattrTotal          *prometheus.Desc
	NFSEventVFSFlushTotal            *prometheus.Desc
	NFSEventVFSFsyncTotal            *prometheus.Desc
	NFSEventVFSLockTotal             *prometheus.Desc
	NFSEventVFSFileReleaseTotal      *prometheus.Desc
	NFSEventTruncationTotal          *prometheus.Desc
	NFSEventWriteExtensionTotal      *prometheus.Desc
	NFSEventSillyRenameTotal         *prometheus.Desc
	NFSEventShortReadTotal           *prometheus.Desc
	NFSEventShortWriteTotal          *prometheus.Desc
	NFSEventJukeboxDelayTotal        *prometheus.Desc
	NFSEventPNFSReadTotal            *prometheus.Desc
	NFSEventPNFSWriteTotal           *prometheus.Desc

	proc procfs.Proc
}

// used to uniquely identify an NFS mount to prevent duplicates
type nfsDeviceIdentifier struct {
	Device   string
	Protocol string
}

func init() {
	registerCollector("mountstats", defaultDisabled, NewMountStatsCollector)
}

// NewMountStatsCollector returns a new Collector exposing NFS statistics.
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
		// For the time being, only NFS statistics are available via this mechanism.
		subsystem = "mountstats_nfs"
	)

	var (
		labels   = []string{"export", "protocol"}
		opLabels = []string{"export", "protocol", "operation"}
	)

	return &mountStatsCollector{
		NFSAgeSecondsTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "age_seconds_total"),
			"The age of the NFS mount in seconds.",
			labels,
			nil,
		),

		NFSReadBytesTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "read_bytes_total"),
			"Number of bytes read using the read() syscall.",
			labels,
			nil,
		),

		NFSWriteBytesTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "write_bytes_total"),
			"Number of bytes written using the write() syscall.",
			labels,
			nil,
		),

		NFSDirectReadBytesTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "direct_read_bytes_total"),
			"Number of bytes read using the read() syscall in O_DIRECT mode.",
			labels,
			nil,
		),

		NFSDirectWriteBytesTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "direct_write_bytes_total"),
			"Number of bytes written using the write() syscall in O_DIRECT mode.",
			labels,
			nil,
		),

		NFSTotalReadBytesTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "total_read_bytes_total"),
			"Number of bytes read from the NFS server, in total.",
			labels,
			nil,
		),

		NFSTotalWriteBytesTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "total_write_bytes_total"),
			"Number of bytes written to the NFS server, in total.",
			labels,
			nil,
		),

		NFSReadPagesTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "read_pages_total"),
			"Number of pages read directly via mmap()'d files.",
			labels,
			nil,
		),

		NFSWritePagesTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "write_pages_total"),
			"Number of pages written directly via mmap()'d files.",
			labels,
			nil,
		),

		NFSTransportBindTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "transport_bind_total"),
			"Number of times the client has had to establish a connection from scratch to the NFS server.",
			labels,
			nil,
		),

		NFSTransportConnectTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "transport_connect_total"),
			"Number of times the client has made a TCP connection to the NFS server.",
			labels,
			nil,
		),

		NFSTransportIdleTimeSeconds: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "transport_idle_time_seconds"),
			"Duration since the NFS mount last saw any RPC traffic, in seconds.",
			labels,
			nil,
		),

		NFSTransportSendsTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "transport_sends_total"),
			"Number of RPC requests for this mount sent to the NFS server.",
			labels,
			nil,
		),

		NFSTransportReceivesTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "transport_receives_total"),
			"Number of RPC responses for this mount received from the NFS server.",
			labels,
			nil,
		),

		NFSTransportBadTransactionIDsTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "transport_bad_transaction_ids_total"),
			"Number of times the NFS server sent a response with a transaction ID unknown to this client.",
			labels,
			nil,
		),

		NFSTransportBacklogQueueTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "transport_backlog_queue_total"),
			"Total number of items added to the RPC backlog queue.",
			labels,
			nil,
		),

		NFSTransportMaximumRPCSlots: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "transport_maximum_rpc_slots"),
			"Maximum number of simultaneously active RPC requests ever used.",
			labels,
			nil,
		),

		NFSTransportSendingQueueTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "transport_sending_queue_total"),
			"Total number of items added to the RPC transmission sending queue.",
			labels,
			nil,
		),

		NFSTransportPendingQueueTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "transport_pending_queue_total"),
			"Total number of items added to the RPC transmission pending queue.",
			labels,
			nil,
		),

		NFSOperationsRequestsTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "operations_requests_total"),
			"Number of requests performed for a given operation.",
			opLabels,
			nil,
		),

		NFSOperationsTransmissionsTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "operations_transmissions_total"),
			"Number of times an actual RPC request has been transmitted for a given operation.",
			opLabels,
			nil,
		),

		NFSOperationsMajorTimeoutsTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "operations_major_timeouts_total"),
			"Number of times a request has had a major timeout for a given operation.",
			opLabels,
			nil,
		),

		NFSOperationsSentBytesTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "operations_sent_bytes_total"),
			"Number of bytes sent for a given operation, including RPC headers and payload.",
			opLabels,
			nil,
		),

		NFSOperationsReceivedBytesTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "operations_received_bytes_total"),
			"Number of bytes received for a given operation, including RPC headers and payload.",
			opLabels,
			nil,
		),

		NFSOperationsQueueTimeSecondsTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "operations_queue_time_seconds_total"),
			"Duration all requests spent queued for transmission for a given operation before they were sent, in seconds.",
			opLabels,
			nil,
		),

		NFSOperationsResponseTimeSecondsTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "operations_response_time_seconds_total"),
			"Duration all requests took to get a reply back after a request for a given operation was transmitted, in seconds.",
			opLabels,
			nil,
		),

		NFSOperationsRequestTimeSecondsTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "operations_request_time_seconds_total"),
			"Duration all requests took from when a request was enqueued to when it was completely handled for a given operation, in seconds.",
			opLabels,
			nil,
		),

		NFSEventInodeRevalidateTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_inode_revalidate_total"),
			"Number of times cached inode attributes are re-validated from the server.",
			labels,
			nil,
		),

		NFSEventDnodeRevalidateTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_dnode_revalidate_total"),
			"Number of times cached dentry nodes are re-validated from the server.",
			labels,
			nil,
		),

		NFSEventDataInvalidateTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_data_invalidate_total"),
			"Number of times an inode cache is cleared.",
			labels,
			nil,
		),

		NFSEventAttributeInvalidateTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_attribute_invalidate_total"),
			"Number of times cached inode attributes are invalidated.",
			labels,
			nil,
		),

		NFSEventVFSOpenTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_vfs_open_total"),
			"Number of times cached inode attributes are invalidated.",
			labels,
			nil,
		),

		NFSEventVFSLookupTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_vfs_lookup_total"),
			"Number of times a directory lookup has occurred.",
			labels,
			nil,
		),

		NFSEventVFSAccessTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_vfs_access_total"),
			"Number of times permissions have been checked.",
			labels,
			nil,
		),

		NFSEventVFSUpdatePageTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_vfs_update_page_total"),
			"Number of updates (and potential writes) to pages.",
			labels,
			nil,
		),

		NFSEventVFSReadPageTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_vfs_read_page_total"),
			"Number of pages read directly via mmap()'d files.",
			labels,
			nil,
		),

		NFSEventVFSReadPagesTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_vfs_read_pages_total"),
			"Number of times a group of pages have been read.",
			labels,
			nil,
		),

		NFSEventVFSWritePageTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_vfs_write_page_total"),
			"Number of pages written directly via mmap()'d files.",
			labels,
			nil,
		),

		NFSEventVFSWritePagesTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_vfs_write_pages_total"),
			"Number of times a group of pages have been written.",
			labels,
			nil,
		),

		NFSEventVFSGetdentsTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_vfs_getdents_total"),
			"Number of times directory entries have been read with getdents().",
			labels,
			nil,
		),

		NFSEventVFSSetattrTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_vfs_setattr_total"),
			"Number of times directory entries have been read with getdents().",
			labels,
			nil,
		),

		NFSEventVFSFlushTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_vfs_flush_total"),
			"Number of pending writes that have been forcefully flushed to the server.",
			labels,
			nil,
		),

		NFSEventVFSFsyncTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_vfs_fsync_total"),
			"Number of times fsync() has been called on directories and files.",
			labels,
			nil,
		),

		NFSEventVFSLockTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_vfs_lock_total"),
			"Number of times locking has been attempted on a file.",
			labels,
			nil,
		),

		NFSEventVFSFileReleaseTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_vfs_file_release_total"),
			"Number of times files have been closed and released.",
			labels,
			nil,
		),

		NFSEventTruncationTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_truncation_total"),
			"Number of times files have been truncated.",
			labels,
			nil,
		),

		NFSEventWriteExtensionTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_write_extension_total"),
			"Number of times a file has been grown due to writes beyond its existing end.",
			labels,
			nil,
		),

		NFSEventSillyRenameTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_silly_rename_total"),
			"Number of times a file was removed while still open by another process.",
			labels,
			nil,
		),

		NFSEventShortReadTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_short_read_total"),
			"Number of times the NFS server gave less data than expected while reading.",
			labels,
			nil,
		),

		NFSEventShortWriteTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_short_write_total"),
			"Number of times the NFS server wrote less data than expected while writing.",
			labels,
			nil,
		),

		NFSEventJukeboxDelayTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_jukebox_delay_total"),
			"Number of times the NFS server indicated EJUKEBOX; retrieving data from offline storage.",
			labels,
			nil,
		),

		NFSEventPNFSReadTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_pnfs_read_total"),
			"Number of NFS v4.1+ pNFS reads.",
			labels,
			nil,
		),

		NFSEventPNFSWriteTotal: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_pnfs_write_total"),
			"Number of NFS v4.1+ pNFS writes.",
			labels,
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

	// store all seen nfsDeviceIdentifiers for deduplication
	deviceList := make(map[nfsDeviceIdentifier]bool)

	for _, m := range mounts {
		// For the time being, only NFS statistics are available via this mechanism
		stats, ok := m.Stats.(*procfs.MountStatsNFS)
		if !ok {
			continue
		}

		deviceIdentifier := nfsDeviceIdentifier{m.Device, stats.Transport.Protocol}
		i := deviceList[deviceIdentifier]
		if i {
			log.Debugf("Skipping duplicate device entry %q", deviceIdentifier)
			continue
		}

		deviceList[deviceIdentifier] = true
		c.updateNFSStats(ch, m.Device, stats.Transport.Protocol, stats)
	}

	return nil
}

func (c *mountStatsCollector) updateNFSStats(ch chan<- prometheus.Metric, export string, protocol string, s *procfs.MountStatsNFS) {
	ch <- prometheus.MustNewConstMetric(
		c.NFSAgeSecondsTotal,
		prometheus.CounterValue,
		s.Age.Seconds(),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSReadBytesTotal,
		prometheus.CounterValue,
		float64(s.Bytes.Read),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSWriteBytesTotal,
		prometheus.CounterValue,
		float64(s.Bytes.Write),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSDirectReadBytesTotal,
		prometheus.CounterValue,
		float64(s.Bytes.DirectRead),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSDirectWriteBytesTotal,
		prometheus.CounterValue,
		float64(s.Bytes.DirectWrite),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSTotalReadBytesTotal,
		prometheus.CounterValue,
		float64(s.Bytes.ReadTotal),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSTotalWriteBytesTotal,
		prometheus.CounterValue,
		float64(s.Bytes.WriteTotal),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSReadPagesTotal,
		prometheus.CounterValue,
		float64(s.Bytes.ReadPages),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSWritePagesTotal,
		prometheus.CounterValue,
		float64(s.Bytes.WritePages),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSTransportBindTotal,
		prometheus.CounterValue,
		float64(s.Transport.Bind),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSTransportConnectTotal,
		prometheus.CounterValue,
		float64(s.Transport.Connect),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSTransportIdleTimeSeconds,
		prometheus.GaugeValue,
		s.Transport.IdleTime.Seconds(),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSTransportSendsTotal,
		prometheus.CounterValue,
		float64(s.Transport.Sends),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSTransportReceivesTotal,
		prometheus.CounterValue,
		float64(s.Transport.Receives),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSTransportBadTransactionIDsTotal,
		prometheus.CounterValue,
		float64(s.Transport.BadTransactionIDs),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSTransportBacklogQueueTotal,
		prometheus.CounterValue,
		float64(s.Transport.CumulativeBacklog),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSTransportMaximumRPCSlots,
		prometheus.GaugeValue,
		float64(s.Transport.MaximumRPCSlotsUsed),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSTransportSendingQueueTotal,
		prometheus.CounterValue,
		float64(s.Transport.CumulativeSendingQueue),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSTransportPendingQueueTotal,
		prometheus.CounterValue,
		float64(s.Transport.CumulativePendingQueue),
		export,
		protocol,
	)

	for _, op := range s.Operations {
		ch <- prometheus.MustNewConstMetric(
			c.NFSOperationsRequestsTotal,
			prometheus.CounterValue,
			float64(op.Requests),
			export,
			protocol,
			op.Operation,
		)

		ch <- prometheus.MustNewConstMetric(
			c.NFSOperationsTransmissionsTotal,
			prometheus.CounterValue,
			float64(op.Transmissions),
			export,
			protocol,
			op.Operation,
		)

		ch <- prometheus.MustNewConstMetric(
			c.NFSOperationsMajorTimeoutsTotal,
			prometheus.CounterValue,
			float64(op.MajorTimeouts),
			export,
			protocol,
			op.Operation,
		)

		ch <- prometheus.MustNewConstMetric(
			c.NFSOperationsSentBytesTotal,
			prometheus.CounterValue,
			float64(op.BytesSent),
			export,
			protocol,
			op.Operation,
		)

		ch <- prometheus.MustNewConstMetric(
			c.NFSOperationsReceivedBytesTotal,
			prometheus.CounterValue,
			float64(op.BytesReceived),
			export,
			protocol,
			op.Operation,
		)

		ch <- prometheus.MustNewConstMetric(
			c.NFSOperationsQueueTimeSecondsTotal,
			prometheus.CounterValue,
			op.CumulativeQueueTime.Seconds(),
			export,
			protocol,
			op.Operation,
		)

		ch <- prometheus.MustNewConstMetric(
			c.NFSOperationsResponseTimeSecondsTotal,
			prometheus.CounterValue,
			op.CumulativeTotalResponseTime.Seconds(),
			export,
			protocol,
			op.Operation,
		)

		ch <- prometheus.MustNewConstMetric(
			c.NFSOperationsRequestTimeSecondsTotal,
			prometheus.CounterValue,
			op.CumulativeTotalRequestTime.Seconds(),
			export,
			protocol,
			op.Operation,
		)
	}

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventInodeRevalidateTotal,
		prometheus.CounterValue,
		float64(s.Events.InodeRevalidate),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventDnodeRevalidateTotal,
		prometheus.CounterValue,
		float64(s.Events.DnodeRevalidate),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventDataInvalidateTotal,
		prometheus.CounterValue,
		float64(s.Events.DataInvalidate),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventAttributeInvalidateTotal,
		prometheus.CounterValue,
		float64(s.Events.AttributeInvalidate),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventVFSOpenTotal,
		prometheus.CounterValue,
		float64(s.Events.VFSOpen),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventVFSLookupTotal,
		prometheus.CounterValue,
		float64(s.Events.VFSLookup),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventVFSAccessTotal,
		prometheus.CounterValue,
		float64(s.Events.VFSAccess),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventVFSUpdatePageTotal,
		prometheus.CounterValue,
		float64(s.Events.VFSUpdatePage),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventVFSReadPageTotal,
		prometheus.CounterValue,
		float64(s.Events.VFSReadPage),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventVFSReadPagesTotal,
		prometheus.CounterValue,
		float64(s.Events.VFSReadPages),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventVFSWritePageTotal,
		prometheus.CounterValue,
		float64(s.Events.VFSWritePage),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventVFSWritePagesTotal,
		prometheus.CounterValue,
		float64(s.Events.VFSWritePages),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventVFSGetdentsTotal,
		prometheus.CounterValue,
		float64(s.Events.VFSGetdents),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventVFSSetattrTotal,
		prometheus.CounterValue,
		float64(s.Events.VFSSetattr),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventVFSFlushTotal,
		prometheus.CounterValue,
		float64(s.Events.VFSFlush),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventVFSFsyncTotal,
		prometheus.CounterValue,
		float64(s.Events.VFSFsync),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventVFSLockTotal,
		prometheus.CounterValue,
		float64(s.Events.VFSLock),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventVFSFileReleaseTotal,
		prometheus.CounterValue,
		float64(s.Events.VFSFileRelease),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventTruncationTotal,
		prometheus.CounterValue,
		float64(s.Events.Truncation),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventWriteExtensionTotal,
		prometheus.CounterValue,
		float64(s.Events.WriteExtension),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventSillyRenameTotal,
		prometheus.CounterValue,
		float64(s.Events.SillyRename),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventShortReadTotal,
		prometheus.CounterValue,
		float64(s.Events.ShortRead),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventShortWriteTotal,
		prometheus.CounterValue,
		float64(s.Events.ShortWrite),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventJukeboxDelayTotal,
		prometheus.CounterValue,
		float64(s.Events.JukeboxDelay),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventPNFSReadTotal,
		prometheus.CounterValue,
		float64(s.Events.PNFSRead),
		export,
		protocol,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventPNFSWriteTotal,
		prometheus.CounterValue,
		float64(s.Events.PNFSWrite),
		export,
		protocol,
	)
}
