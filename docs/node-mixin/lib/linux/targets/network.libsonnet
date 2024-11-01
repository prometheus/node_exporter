local g = import '../../g.libsonnet';
local prometheusQuery = g.query.prometheus;
local lokiQuery = g.query.loki;

{
  new(this): {
    local variables = this.grafana.variables.main,
    local config = this.config,
    local prometheusDatasource = '${' + variables.datasources.prometheus.name + '}',
    local lokiDatasource = '${' + variables.datasources.loki.name + '}',

    networkUp:
      prometheusQuery.new(
        prometheusDatasource,
        'node_network_up{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('{{device}}'),
    networkCarrier:
      prometheusQuery.new(
        prometheusDatasource,
        'node_network_carrier{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('{{device}}'),
    networkArpEntries:
      prometheusQuery.new(
        prometheusDatasource,
        'node_arp_entries{%(queriesSelector)s}' % variables,
      ),
    networkMtuBytes:
      prometheusQuery.new(
        prometheusDatasource,
        'node_network_mtu_bytes{%(queriesSelector)s}' % variables,
      ),
    networkSpeedBitsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'node_network_speed_bytes{%(queriesSelector)s} * 8' % variables,
      ),
    networkTransmitQueueLength:
      prometheusQuery.new(
        prometheusDatasource,
        'node_network_transmit_queue_length{%(queriesSelector)s}' % variables,
      ),
    networkInfo:
      prometheusQuery.new(
        prometheusDatasource,
        'node_network_info{%(queriesSelector)s}' % variables,
      ),

    networkOutBitPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_network_transmit_bytes_total{%(queriesSelector)s}[$__rate_interval])*8' % variables
      )
      + prometheusQuery.withLegendFormat('{{ device }} transmitted'),
    networkInBitPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_network_receive_bytes_total{%(queriesSelector)s}[$__rate_interval])*8' % variables
      )
      + prometheusQuery.withLegendFormat('{{ device }} received'),
    networkOutBitPerSecFiltered:
      prometheusQuery.new(
        prometheusDatasource,
        |||
          irate(node_network_transmit_bytes_total{%(queriesSelector)s}[$__rate_interval])*8
          # only show interfaces that had traffic change at least once during selected dashboard interval:
          and
          increase(
              node_network_transmit_bytes_total{%(queriesSelector)s}[$__range]
              ) > 0
        ||| % variables
      )
      + prometheusQuery.withLegendFormat('{{ device }} transmitted'),
    networkInBitPerSecFiltered:
      prometheusQuery.new(
        prometheusDatasource,
        |||
          irate(node_network_receive_bytes_total{%(queriesSelector)s}[$__rate_interval])*8
          # only show interfaces that had traffic change at least once during selected dashboard interval:
          and
          increase(
              node_network_receive_bytes_total{%(queriesSelector)s}[$__range]
              ) > 0
        ||| % variables
      )
      + prometheusQuery.withLegendFormat('{{ device }} received'),


    networkOutErrorsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_network_transmit_errs_total{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('{{ device }} errors transmitted'),
    networkInErrorsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_network_receive_errs_total{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('{{ device }} errors received'),
    networkOutDroppedPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_network_transmit_drop_total{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('{{ device }} transmitted dropped'),
    networkInDroppedPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_network_receive_drop_total{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('{{ device }} received dropped'),

    networkInPacketsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_network_receive_packets_total{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('{{ device }} received'),
    networkOutPacketsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_network_transmit_packets_total{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('{{ device }} transmitted'),

    networkInMulticastPacketsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_network_receive_multicast_total{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('{{ device }} received'),
    networkOutMulticastPacketsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_network_transmit_multicast_total{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('{{ device }} transmitted'),
    networkFifoInPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_network_receive_fifo_total{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('{{ device }} received'),
    networkFifoOutPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_network_transmit_fifo_total{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('{{ device }} transmitted'),

    networkCompressedInPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_network_receive_compressed_total{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('{{ device }} received'),
    networkCompressedOutPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_network_transmit_compressed_total{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('{{ device }} transmitted'),

    networkNFConntrackEntries:
      prometheusQuery.new(
        prometheusDatasource,
        'node_nf_conntrack_entries{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('NF conntrack entries'),
    networkNFConntrackLimits:
      prometheusQuery.new(
        prometheusDatasource,
        'node_nf_conntrack_entries_limit{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('NF conntrack limits'),

    networkSoftnetProcessedPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_softnet_processed_total{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('CPU {{ cpu }} processed'),
    networkSoftnetDroppedPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_softnet_dropped_total{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('CPU {{ cpu }} dropped'),
    networkSoftnetSqueezedPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_softnet_times_squeezed_total{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('CPU {{ cpu }} out of quota'),

    networkSocketsUsed:
      prometheusQuery.new(
        prometheusDatasource,
        'node_sockstat_sockets_used{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('IPv4 sockets in use'),
    networkSocketsTCPAllocated:
      prometheusQuery.new(
        prometheusDatasource,
        'node_sockstat_TCP_alloc{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('Allocated'),
    networkSocketsTCPIPv6:
      prometheusQuery.new(
        prometheusDatasource,
        'node_sockstat_TCP6_inuse{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('IPv6 in use'),
    networkSocketsTCPIPv4:
      prometheusQuery.new(
        prometheusDatasource,
        'node_sockstat_TCP_inuse{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('IPv4 in use'),
    networkSocketsTCPOrphans:
      prometheusQuery.new(
        prometheusDatasource,
        'node_sockstat_TCP_orphan{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('Orphan sockets'),
    networkSocketsTCPTimeWait:
      prometheusQuery.new(
        prometheusDatasource,
        'node_sockstat_TCP_tw{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('Time wait'),

    networkSocketsUDPLiteInUse:
      prometheusQuery.new(
        prometheusDatasource,
        'node_sockstat_UDPLITE_inuse{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('IPv4 UDPLITE in use'),
    networkSocketsUDPInUse:
      prometheusQuery.new(
        prometheusDatasource,
        'node_sockstat_UDP_inuse{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('IPv4 UDP in use'),
    networkSocketsUDPLiteIPv6InUse:
      prometheusQuery.new(
        prometheusDatasource,
        'node_sockstat_UDPLITE6_inuse{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('IPv6 UDPLITE in use'),
    networkSocketsUDPIPv6InUse:
      prometheusQuery.new(
        prometheusDatasource,
        'node_sockstat_UDP6_inuse{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('IPv6 UDP in use'),

    networkSocketsFragInUse:
      prometheusQuery.new(
        prometheusDatasource,
        'node_sockstat_FRAG_inuse{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('IPv4 Frag sockets in use'),
    networkSocketsFragIPv6InUse:
      prometheusQuery.new(
        prometheusDatasource,
        'node_sockstat_FRAG6_inuse{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('IPv6 Frag sockets in use'),
    networkSocketsRawInUse:
      prometheusQuery.new(
        prometheusDatasource,
        'node_sockstat_RAW_inuse{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('IPv4 Raw sockets in use'),
    networkSocketsIPv6RawInUse:
      prometheusQuery.new(
        prometheusDatasource,
        'node_sockstat_RAW6_inuse{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('IPv6 Raw sockets in use'),

    networkSocketsTCPMemoryPages:
      prometheusQuery.new(
        prometheusDatasource,
        'node_sockstat_TCP_mem{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('Memory pages allocated for TCP sockets'),
    networkSocketsUDPMemoryPages:
      prometheusQuery.new(
        prometheusDatasource,
        'node_sockstat_UDP_mem{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('Memory pages allocated for UDP sockets'),

    networkSocketsTCPMemoryBytes:
      prometheusQuery.new(
        prometheusDatasource,
        'node_sockstat_TCP_mem_bytes{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('Memory bytes allocated for TCP sockets'),
    networkSocketsUDPMemoryBytes:
      prometheusQuery.new(
        prometheusDatasource,
        'node_sockstat_UDP_mem_bytes{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('Memory bytes allocated for UDP sockets'),

    networkNetstatIPInOctetsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_IpExt_InOctets{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('Octets received'),
    networkNetstatIPOutOctetsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_IpExt_OutOctets{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('Octets transmitted'),

    networkNetstatTCPInSegmentsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Tcp_InSegs{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('TCP received'),
    networkNetstatTCPOutSegmentsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Tcp_OutSegs{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('TCP transmitted'),

    networkNetstatTCPOverflowPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_TcpExt_ListenOverflows{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('TCP overflow'),

    networkNetstatTCPListenDropsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_TcpExt_ListenDrops{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('TCP ListenDrops - SYNs to LISTEN sockets ignored'),

    networkNetstatTCPRetransPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_TcpExt_TCPSynRetrans{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('TCP SYN rentransmits'),

    networkNetstatTCPRetransSegPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Tcp_RetransSegs{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('TCP retransmitted segments, containing one or more previously transmitted octets'),
    networkNetstatTCPInWithErrorsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Tcp_InErrs{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('TCP received with errors'),

    networkNetstatTCPOutWithRstPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Tcp_OutRsts{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('TCP segments sent with RST flag'),

    networkNetstatIPInUDPPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Udp_InDatagrams{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('UDP received'),

    networkNetstatIPOutUDPPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Udp_OutDatagrams{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('UDP transmitted'),

    networkNetstatIPInUDP6PerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Udp6_InDatagrams{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('UDP6 received'),

    networkNetstatIPOutUDP6PerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Udp6_OutDatagrams{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('UDP6 transmitted'),

    //UDP errors
    networkNetstatUDPLiteInErrorsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_UdpLite_InErrors{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('UDPLite InErrors'),

    networkNetstatUDPInErrorsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Udp_InErrors{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('UDP InErrors'),
    networkNetstatUDP6InErrorsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Udp6_InErrors{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('UDP6 InErrors'),
    networkNetstatUDPNoPortsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Udp_NoPorts{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('UDP NoPorts'),
    networkNetstatUDP6NoPortsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Udp6_NoPorts{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('UDP6 NoPorts'),
    networkNetstatUDPRcvBufErrsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Udp_RcvbufErrors{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('UDP receive buffer errors'),
    networkNetstatUDP6RcvBufErrsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Udp6_RcvbufErrors{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('UDP6 receive buffer errors'),
    networkNetstatUDPSndBufErrsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Udp_SndbufErrors{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('UDP transmit buffer errors'),
    networkNetstatUDP6SndBufErrsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Udp6_SndbufErrors{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('UDP6 transmit buffer errors'),

    //ICMP
    networkNetstatICMPInPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Icmp_InMsgs{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('ICMP received'),
    networkNetstatICMPOutPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Icmp_OutMsgs{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('ICMP transmitted'),
    networkNetstatICMP6InPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Icmp6_InMsgs{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('ICMP6 received'),
    networkNetstatICMP6OutPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Icmp6_OutMsgs{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('ICMP6 transmitted'),

    networkNetstatICMPInErrorsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Icmp_InErrors{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('ICMP6 errors'),
    networkNetstatICM6PInErrorsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Icmp6_InErrors{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('ICMP6 errors'),
  },
}
