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
	"sort"

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

func init() {
	Factories["mountstats"] = NewMountStatsCollector
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

		NFSEventInodeRevalidateTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "event_inode_revalidate_total"),
			"Number of times cached inode attributes are re-validated from the server.",
			labels,
			nil,
		),

		NFSEventDnodeRevalidateTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "event_dnode_revalidate_total"),
			"Number of times cached dentry nodes are re-validated from the server.",
			labels,
			nil,
		),

		NFSEventDataInvalidateTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "event_data_invalidate_total"),
			"Number of times an inode cache is cleared.",
			labels,
			nil,
		),

		NFSEventAttributeInvalidateTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "event_attribute_invalidate_total"),
			"Number of times cached inode attributes are invalidated.",
			labels,
			nil,
		),

		NFSEventVFSOpenTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "event_vfs_open_total"),
			"Number of times cached inode attributes are invalidated.",
			labels,
			nil,
		),

		NFSEventVFSLookupTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "event_vfs_lookup_total"),
			"Number of times a directory lookup has occurred.",
			labels,
			nil,
		),

		NFSEventVFSAccessTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "event_vfs_access_total"),
			"Number of times permissions have been checked.",
			labels,
			nil,
		),

		NFSEventVFSUpdatePageTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "event_vfs_update_page_total"),
			"Number of updates (and potential writes) to pages.",
			labels,
			nil,
		),

		NFSEventVFSReadPageTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "event_vfs_read_page_total"),
			"Number of pages read directly via mmap()'d files.",
			labels,
			nil,
		),

		NFSEventVFSReadPagesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "event_vfs_read_pages_total"),
			"Number of times a group of pages have been read.",
			labels,
			nil,
		),

		NFSEventVFSWritePageTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "event_vfs_write_page_total"),
			"Number of pages written directly via mmap()'d files.",
			labels,
			nil,
		),

		NFSEventVFSWritePagesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "event_vfs_write_pages_total"),
			"Number of times a group of pages have been written.",
			labels,
			nil,
		),

		NFSEventVFSGetdentsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "event_vfs_getdents_total"),
			"Number of times directory entries have been read with getdents().",
			labels,
			nil,
		),

		NFSEventVFSSetattrTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "event_vfs_setattr_total"),
			"Number of times directory entries have been read with getdents().",
			labels,
			nil,
		),

		NFSEventVFSFlushTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "event_vfs_flush_total"),
			"Number of pending writes that have been forcefully flushed to the server.",
			labels,
			nil,
		),

		NFSEventVFSFsyncTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "event_vfs_fsync_total"),
			"Number of times fsync() has been called on directories and files.",
			labels,
			nil,
		),

		NFSEventVFSLockTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "event_vfs_lock_total"),
			"Number of times locking has been attempted on a file.",
			labels,
			nil,
		),

		NFSEventVFSFileReleaseTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "event_vfs_file_release_total"),
			"Number of times files have been closed and released.",
			labels,
			nil,
		),

		NFSEventTruncationTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "event_truncation_total"),
			"Number of times files have been truncated.",
			labels,
			nil,
		),

		NFSEventWriteExtensionTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "event_write_extension_total"),
			"Number of times a file has been grown due to writes beyond its existing end.",
			labels,
			nil,
		),

		NFSEventSillyRenameTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "event_silly_rename_total"),
			"Number of times a file was removed while still open by another process.",
			labels,
			nil,
		),

		NFSEventShortReadTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "event_short_read_total"),
			"Number of times the NFS server gave less data than expected while reading.",
			labels,
			nil,
		),

		NFSEventShortWriteTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "event_short_write_total"),
			"Number of times the NFS server wrote less data than expected while writing.",
			labels,
			nil,
		),

		NFSEventJukeboxDelayTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "event_jukebox_delay_total"),
			"Number of times the NFS server indicated EJUKEBOX; retrieving data from offline storage.",
			labels,
			nil,
		),

		NFSEventPNFSReadTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "event_pnfs_read_total"),
			"Number of NFS v4.1+ pNFS reads.",
			labels,
			nil,
		),

		NFSEventPNFSWriteTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "event_pnfs_write_total"),
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
	var deviceList []string

	for _, m := range mounts {
		// For the time being, only NFS statistics are available via this mechanism
		stats, ok := m.Stats.(*procfs.MountStatsNFS)
		if !ok {
			continue
		}

		sort.Strings(deviceList)
		i := sort.SearchStrings(deviceList, m.Device)
		if i < len(deviceList) && deviceList[i] == m.Device {
			log.Debugf("Skipping duplicate device entry %q", m.Device)
		} else {
			deviceList = append(deviceList, m.Device)
			c.updateNFSStats(ch, m.Device, stats)
		}
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

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventInodeRevalidateTotal,
		prometheus.CounterValue,
		float64(s.Events.InodeRevalidate),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventDnodeRevalidateTotal,
		prometheus.CounterValue,
		float64(s.Events.DnodeRevalidate),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventDataInvalidateTotal,
		prometheus.CounterValue,
		float64(s.Events.DataInvalidate),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventAttributeInvalidateTotal,
		prometheus.CounterValue,
		float64(s.Events.AttributeInvalidate),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventVFSOpenTotal,
		prometheus.CounterValue,
		float64(s.Events.VFSOpen),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventVFSLookupTotal,
		prometheus.CounterValue,
		float64(s.Events.VFSLookup),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventVFSAccessTotal,
		prometheus.CounterValue,
		float64(s.Events.VFSAccess),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventVFSUpdatePageTotal,
		prometheus.CounterValue,
		float64(s.Events.VFSUpdatePage),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventVFSReadPageTotal,
		prometheus.CounterValue,
		float64(s.Events.VFSReadPage),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventVFSReadPagesTotal,
		prometheus.CounterValue,
		float64(s.Events.VFSReadPages),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventVFSWritePageTotal,
		prometheus.CounterValue,
		float64(s.Events.VFSWritePage),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventVFSWritePagesTotal,
		prometheus.CounterValue,
		float64(s.Events.VFSWritePages),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventVFSGetdentsTotal,
		prometheus.CounterValue,
		float64(s.Events.VFSGetdents),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventVFSSetattrTotal,
		prometheus.CounterValue,
		float64(s.Events.VFSSetattr),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventVFSFlushTotal,
		prometheus.CounterValue,
		float64(s.Events.VFSFlush),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventVFSFsyncTotal,
		prometheus.CounterValue,
		float64(s.Events.VFSFsync),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventVFSLockTotal,
		prometheus.CounterValue,
		float64(s.Events.VFSLock),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventVFSFileReleaseTotal,
		prometheus.CounterValue,
		float64(s.Events.VFSFileRelease),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventTruncationTotal,
		prometheus.CounterValue,
		float64(s.Events.Truncation),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventWriteExtensionTotal,
		prometheus.CounterValue,
		float64(s.Events.WriteExtension),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventSillyRenameTotal,
		prometheus.CounterValue,
		float64(s.Events.SillyRename),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventShortReadTotal,
		prometheus.CounterValue,
		float64(s.Events.ShortRead),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventShortWriteTotal,
		prometheus.CounterValue,
		float64(s.Events.ShortWrite),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventJukeboxDelayTotal,
		prometheus.CounterValue,
		float64(s.Events.JukeboxDelay),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventPNFSReadTotal,
		prometheus.CounterValue,
		float64(s.Events.PNFSRead),
		export,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NFSEventPNFSWriteTotal,
		prometheus.CounterValue,
		float64(s.Events.PNFSWrite),
		export,
	)
}
