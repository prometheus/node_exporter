local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';
local dashboard = grafana.dashboard;
local row = grafana.row;
local prometheus = grafana.prometheus;
local template = grafana.template;
local graphPanel = grafana.graphPanel;
local nodePanels = import '../lib/panels/panels.libsonnet';
local commonPanels = import '../lib/panels/common/panels.libsonnet';
local nodeTimeseries = nodePanels.timeseries;
local common = import '../lib/common.libsonnet';

{

  new(config=null, platform=null):: {
    local c = common.new(config=config, platform=platform),
    local commonPromTarget = c.commonPromTarget,
    local templates = c.templates,
    local q = c.queries,

    local networkTrafficPanel =
      commonPanels.networkTrafficGraph.new('Network Traffic')
      .addTarget(commonPromTarget(
        expr=q.networkReceiveBitsPerSec,
        legendFormat='{{device}} received',
      ))
      .addTarget(commonPromTarget(
        expr=q.networkTransmitBitsPerSec,
        legendFormat='{{device}} transmitted',
      )),

    local networkPacketsPanel =
      nodeTimeseries.new('Unicast Packets')
      .addTarget(commonPromTarget(
        'irate(node_network_receive_packets_total{%(nodeQuerySelector)s,}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='{{device}} received',
      ))
      .addTarget(commonPromTarget(
        'irate(node_network_transmit_packets_total{%(nodeQuerySelector)s,}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='{{device}} transmitted',
      ))
      .withDecimals(1)
      .withUnits('pps')
      .withNegativeYByRegex('transmit')
      .withAxisLabel('out(-) / in(+)'),

    local networkErrorsPanel =
      nodeTimeseries.new('Network Errors')
      .addTarget(commonPromTarget(
        expr=q.networkReceiveErrorsPerSec,
        legendFormat='{{device}} received',
      ))
      .addTarget(commonPromTarget(
        expr=q.networkTransmitErrorsPerSec,
        legendFormat='{{device}} transmitted',
      ))
      .withDecimals(1)
      .withUnits('pps')
      .withNegativeYByRegex('transmit')
      .withAxisLabel('out(-) / in(+)'),

    local networkDropsPanel =
      nodeTimeseries.new('Dropped Packets')
      .addTarget(commonPromTarget(
        expr=q.networkReceiveDropsPerSec,
        legendFormat='{{device}} received',
      ))
      .addTarget(commonPromTarget(
        expr=q.networkTransmitDropsPerSec,
        legendFormat='{{device}} transmitted',
      ))
      .withDecimals(1)
      .withUnits('pps')
      .withNegativeYByRegex('transmit')
      .withAxisLabel('out(-) / in(+)'),
    local networkCompressedPanel =
      nodeTimeseries.new('Compressed Packets')
      .addTarget(commonPromTarget(
        'irate(node_network_receive_compressed_total{%(nodeQuerySelector)s,}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='{{device}} received',
      ))
      .addTarget(commonPromTarget(
        'irate(node_network_transmit_compressed_total{%(nodeQuerySelector)s,}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='{{device}} transmitted',
      ))
      .withDecimals(1)
      .withUnits('pps')
      .withNegativeYByRegex('transmit')
      .withAxisLabel('out(-) / in(+)'),

    local networkMulticastPanel =
      nodeTimeseries.new('Multicast Packets')
      .addTarget(commonPromTarget(
        'irate(node_network_receive_multicast_total{%(nodeQuerySelector)s,}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='{{device}} received',
      ))
      .addTarget(commonPromTarget(
        'irate(node_network_transmit_multicast_total{%(nodeQuerySelector)s,}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='{{device}} transmitted',
      ))
      .withDecimals(1)
      .withUnits('pps')
      .withNegativeYByRegex('transmit'),

    local networkFifoPanel =
      nodeTimeseries.new('Network FIFO')
      .addTarget(commonPromTarget(
        'irate(node_network_receive_fifo_total{%(nodeQuerySelector)s,}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='{{device}} received',
      ))
      .addTarget(commonPromTarget(
        'irate(node_network_transmit_fifo_total{%(nodeQuerySelector)s,}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='{{device}} transmitted',
      ))
      .withDecimals(1)
      .withUnits('pps')
      .withNegativeYByRegex('transmit')
      .withAxisLabel('out(-) / in(+)'),

    local networkNFConntrack =
      nodeTimeseries.new('NF Conntrack')
      .addTarget(commonPromTarget(
        'node_nf_conntrack_entries{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='NF conntrack entries',
      ))
      .addTarget(commonPromTarget(
        'node_nf_conntrack_entries_limit{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='NF conntrack limits',
      ))
      .withFillOpacity(0),

    local networkSoftnetPanel =
      nodeTimeseries.new('Softnet Packets')
      .addTarget(commonPromTarget(
        'irate(node_softnet_processed_total{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='CPU {{cpu }} proccessed',
      ))
      .addTarget(commonPromTarget(
        'irate(node_softnet_dropped_total{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='CPU {{cpu }} dropped',
      ))
      .withDecimals(1)
      .withUnits('pps')
      .withNegativeYByRegex('dropped')
      .withAxisLabel('Dropped(-) | Processed(+)'),

    local networkSoftnetSqueezePanel =
      nodeTimeseries.new('Softnet Out of Quota')
      .addTarget(commonPromTarget(
        'irate(node_softnet_times_squeezed_total{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='CPU {{cpu}} out of quota',
      ))
      .withDecimals(1)
      .withUnits('pps'),

    local networkInterfacesTable =
      nodePanels.table.new(
        title='Network Interfaces Overview'
      )
      // "Value #A"
      .addTarget(commonPromTarget(
        expr='node_network_up{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        format='table',
        instant=true,
      ))
      // "Value #B"
      .addTarget(commonPromTarget(
        expr='node_network_carrier{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        format='table',
        instant=true,
      ))
      // "Value #C"
      .addTarget(commonPromTarget(
        expr=q.networkTransmitBitsPerSec,
        format='table',
        instant=true,
      ))
      // "Value #D"
      .addTarget(commonPromTarget(
        expr=q.networkReceiveBitsPerSec,
        format='table',
        instant=true,
      ))
      // "Value #E"
      .addTarget(commonPromTarget(
        expr='node_arp_entries{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        format='table',
        instant=true,
      ))
      // "Value #F"
      .addTarget(commonPromTarget(
        expr='node_network_mtu_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        format='table',
        instant=true,
      ))
      // "Value #G"
      .addTarget(commonPromTarget(
        expr='node_network_speed_bytes{%(nodeQuerySelector)s} * 8' % config { nodeQuerySelector: c.nodeQuerySelector },
        format='table',
        instant=true,
      ))
      // "Value #H"
      .addTarget(commonPromTarget(
        expr='node_network_transmit_queue_length{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        format='table',
        instant=true,
      ))
      // "VALUE #I"
      .addTarget(commonPromTarget(
        expr='node_network_info{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        format='table',
        instant=true,
      ))
      // "VALUE #J"
      // .addTarget(commonPromTarget(
      //   expr='node_network_protocol_type{%(nodeQuerySelector)s}' % config {nodeQuerySelector: c.nodeQuerySelector},
      //   format="table",
      //   instant=true,
      // ))
      .withTransform()
      .joinByField(field='device')
      // .merge()
      .filterFieldsByName('device|address|duplex|Value.+')
      .organize(
        excludeByName={
          'Value #I': true,
        },
        renameByName=
        {
          device: 'Interface',
          address: 'Address',
          duplex: 'Duplex',
          'Value #A': 'Up',
          'Value #B': 'Carrier',
          'Value #C': 'Transmit',
          'Value #D': 'Receive',
          'Value #E': 'ARP entries',
          'Value #F': 'MTU',
          'Value #G': 'Speed',
          'Value #H': 'Queue length',
          // "Value #J": "Type",
        }
      )
      .addOverride(
        matcher={
          id: 'byRegexp',
          options: 'Speed',
        },
        properties=[
          {
            id: 'unit',
            value: 'bps',
          },
        ]
      )
      .addOverride(
        matcher={
          id: 'byRegexp',
          options: 'Carrier|Up',
        },
        properties=[
          {
            id: 'custom.displayMode',
            value: 'color-text',
          },
          {
            id: 'mappings',
            value: [
              {
                type: 'value',
                options: {
                  '0': {
                    text: 'Down',
                    color: 'light-red',
                    index: 1,
                  },
                  '1': {
                    text: 'Up',
                    color: 'light-green',
                    index: 0,
                  },
                },
              },
            ],
          },
        ]
      )
      // TODO
      // possible values: https://github.com/torvalds/linux/blob/master/include/uapi/linux/if_arp.h
      // .addOverride(
      //   matcher={
      //     id: 'byName',
      //     options: 'Type',
      //   },
      //   properties=[
      //     {
      //       "id": "mappings",
      //       "value": [
      //         {
      //           "type": "value",
      //           "options": {
      //             "0": {
      //               "text": "NET/ROM pseudo",
      //               "index": 0
      //             },
      //             "1": {
      //               "text": "Ethernet 10Mbps",
      //               "index": 1
      //             },
      //             "2": {
      //               "text": "Experimental Ethernet",
      //               "index": 2
      //             },
      //             "3": {
      //               "text": "AX.25 Level 2",
      //               "index": 3
      //             },
      //             "4": {
      //               "text": "PROnet token ring",
      //               "index": 4
      //             },
      //             "5": {
      //               "text": "Chaosnet",
      //               "index": 5
      //             },
      //           }
      //         }
      //       ]
      //     }
      //   ]
      // )
      .addOverride(
        matcher={
          id: 'byRegexp',
          options: 'Transmit|Receive',
        },
        properties=[
          {
            id: 'unit',
            value: 'bps',
          },
          {
            id: 'custom.displayMode',
            value: 'gradient-gauge',
          },
          {
            id: 'color',
            value: {
              mode: 'continuous-BlYlRd',
            },
          },
          {
            id: 'max',
            value: 1000 * 1000 * 100,
          },
        ]
      )
    ,

    local networkOperStatus =
      nodeTimeseries.new(
        title='Network Interfaces Carrier Status'
      )
      .withColor(mode='palette-classic')
      .withFillOpacity(100)
      .withLegend(mode='list')
      .addTarget(commonPromTarget(
        expr='node_network_carrier{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='{{device}}'
      ))
      + {
        maxDataPoints: 100,
        type: 'status-history',
        fieldConfig+: {
          defaults+: {
            mappings+: [
              {
                type: 'value',
                options: {
                  '1': {
                    text: 'Up',
                    color: 'light-green',
                    index: 1,
                  },
                },
              },
              {
                type: 'value',
                options: {
                  '0': {
                    text: 'Down',
                    color: 'light-red',
                    index: 0,
                  },
                },
              },

            ],
          },
        },
      },
    // https://github.com/prometheus/node_exporter/pull/2346/files#diff-3699c850869aecf912f8e8272958b556913fc266534206833a5dcb7d6cca3610
    local networkSockstatTCP =
      nodeTimeseries.new(
        title='Sockets TCP'
      )
      .addTarget(commonPromTarget(
        expr='node_sockstat_TCP_alloc{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Allocated'
      ))
      .addTarget(commonPromTarget(
        expr='node_sockstat_TCP6_inuse{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='IPv6 In use'
      ))
      .addTarget(commonPromTarget(
        expr='node_sockstat_TCP_inuse{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='IPv4 In use'
      ))
      .addTarget(commonPromTarget(
        expr='node_sockstat_TCP_orphan{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Orphan sockets'
      ))
      .addTarget(commonPromTarget(
        expr='node_sockstat_TCP_tw{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Time wait'
      )),

    local networkSockstatUDP =
      nodeTimeseries.new(
        title='Sockets UDP'
      )
      .addTarget(commonPromTarget(
        expr='node_sockstat_UDPLITE_inuse{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='IPv4 UDPLITE in use'
      ))
      .addTarget(commonPromTarget(
        expr='node_sockstat_UDP_inuse{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='IPv4 UDP in use'
      ))
      .addTarget(commonPromTarget(
        expr='node_sockstat_UDPLITE6_inuse{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='IPv6 UDPLITE in use'
      ))
      .addTarget(commonPromTarget(
        expr='node_sockstat_UDP6_inuse{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='IPv6 UDP in use'
      )),

    local networkSockstatOther =
      nodeTimeseries.new(
        title='Sockets Other'
      )
      .addTarget(commonPromTarget(
        expr='node_sockstat_FRAG_inuse{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='IPv4 Frag sockets in use'
      ))
      .addTarget(commonPromTarget(
        expr='node_sockstat_FRAG6_inuse{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='IPv6 Frag sockets in use'
      ))
      .addTarget(commonPromTarget(
        expr='node_sockstat_RAW_inuse{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='IPv4 Raw sockets in use'
      ))
      .addTarget(commonPromTarget(
        expr='node_sockstat_RAW6_inuse{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='IPv6 Raw sockets in use'
      )),


    local networkSockstatMemory =
      nodeTimeseries.new(
        title='Sockets Memory'
      )
      .withMaxDataPoints(100)
      .addTarget(commonPromTarget(
        expr='node_sockstat_TCP_mem{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Memory pages allocated for TCP sockets'
      ))
      .addTarget(commonPromTarget(
        expr='node_sockstat_UDP_mem{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Memory pages allocated for UDP sockets'
      ))
      .addTarget(commonPromTarget(
        expr='node_sockstat_TCP_mem_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Memory bytes allocated for TCP sockets'
      ))
      .addTarget(commonPromTarget(
        expr='node_sockstat_UDP_mem_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Memory bytes allocated for UDP sockets'
      ))
      .addOverride(
        matcher={
          id: 'byRegexp',
          options: '/bytes/',
        },
        properties=[
          {
            id: 'unit',
            value: 'bytes',
          },
          {
            id: 'custom.drawStyle',
            value: 'lines',
          },
          {
            id: 'custom.drawStyle',
            value: 'bars',
          },
          {
            id: 'custom.stacking',
            value: {
              mode: 'normal',
              group: 'A',
            },
          },
        ]
      ),

    local networkSockstatAll =
      nodeTimeseries.new(
        title='Sockets in use'
      )
      .addTarget(commonPromTarget(
        expr='node_sockstat_sockets_used{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='IPv4 sockets in use'
      )),

    local networkNetstatIP =
      nodeTimeseries.new(
        title='IP octets'
      )
      .withUnits('oct/s')
      .withNegativeYByRegex('transmit')
      .withAxisLabel('out(-) / in(+)')
      .addTarget(commonPromTarget(
        expr='irate(node_netstat_IpExt_InOctets{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Octets received'
      ))
      .addTarget(commonPromTarget(
        expr='irate(node_netstat_IpExt_OutOctets{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Octets transmitted'
      )),


    local networkNetstatTCP =
      nodeTimeseries.new(
        title='TCP segments'
      )
      .withUnits('seg/s')
      .withNegativeYByRegex('transmit')
      .withAxisLabel('out(-) / in(+)')
      .addTarget(commonPromTarget(
        expr='irate(node_netstat_Tcp_InSegs{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='TCP received'
      ))
      .addTarget(commonPromTarget(
        expr='irate(node_netstat_Tcp_OutSegs{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='TCP transmitted'
      )),

    local networkNetstatTCPerrors =
      nodeTimeseries.new(
        title='TCP errors rate'
      )
      .withUnits('err/s')
      .addTarget(commonPromTarget(
        expr='irate(node_netstat_TcpExt_ListenOverflows{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='TCP overflow'
      ))
      .addTarget(commonPromTarget(
        expr='irate(node_netstat_TcpExt_ListenDrops{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='TCP ListenDrops - SYNs to LISTEN sockets ignored'
      ))
      .addTarget(commonPromTarget(
        expr='irate(node_netstat_TcpExt_TCPSynRetrans{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='TCP SYN rentransmits'
      ))
      .addTarget(commonPromTarget(
        expr='irate(node_netstat_Tcp_RetransSegs{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='TCP retransmitted segments, containing one or more previously transmitted octets'
      ))
      .addTarget(commonPromTarget(
        expr='irate(node_netstat_Tcp_InErrs{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='TCP received with errors'
      ))
      .addTarget(commonPromTarget(
        expr='irate(node_netstat_Tcp_OutRsts{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='TCP segments sent with RST flag'
      )),

    local networkNetstatUDP =
      nodeTimeseries.new(
        title='UDP datagrams'
      )
      .withUnits('dat/s')
      .withNegativeYByRegex('transmit')
      .withAxisLabel('out(-) / in(+)')
      .addTarget(commonPromTarget(
        expr='irate(node_netstat_Udp_InDatagrams{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='UDP received'
      ))
      .addTarget(commonPromTarget(
        expr='irate(node_netstat_Udp_OutDatagrams{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='UDP transmitted'
      ))
      .addTarget(commonPromTarget(
        expr='irate(node_netstat_Udp6_InDatagrams{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='UDP6 received'
      ))
      .addTarget(commonPromTarget(
        expr='irate(node_netstat_Udp6_OutDatagrams{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='UDP6 transmitted'
      )),

    local networkNetstatUDPerrors =
      nodeTimeseries.new(
        title='UDP errors rate'
      )
      .withUnits('err/s')
      .addTarget(commonPromTarget(
        expr='irate(node_netstat_UdpLite_InErrors{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='UDPLite InErrors'
      ))
      .addTarget(commonPromTarget(
        expr='irate(node_netstat_Udp_InErrors{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='UDP InErrors'
      ))
      .addTarget(commonPromTarget(
        expr='irate(node_netstat_Udp6_InErrors{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='UDP6 InErrors'
      ))
      .addTarget(commonPromTarget(
        expr='irate(node_netstat_Udp_NoPorts{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='UDP NoPorts'
      ))
      .addTarget(commonPromTarget(
        expr='irate(node_netstat_Udp6_NoPorts{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='UDP6 NoPorts'
      ))
      .addTarget(commonPromTarget(
        expr='irate(node_netstat_Udp_RcvbufErrors{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='UDP receive buffer errors'
      ))
      .addTarget(commonPromTarget(
        expr='irate(node_netstat_Udp6_RcvbufErrors{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='UDP6 receive buffer errors'
      ))
      .addTarget(commonPromTarget(
        expr='irate(node_netstat_Udp_SndbufErrors{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='UDP send buffer errors'
      ))
      .addTarget(commonPromTarget(
        expr='irate(node_netstat_Udp6_SndbufErrors{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='UDP6 send buffer errors'
      )),


    local networkNetstatICMP =
      nodeTimeseries.new(
        title='ICMP messages'
      )
      .withUnits('msg/s')
      .withNegativeYByRegex('transmit')
      .withAxisLabel('out(-) / in(+)')
      .addTarget(commonPromTarget(
        expr='irate(node_netstat_Icmp_InMsgs{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='ICMP received'
      ))
      .addTarget(commonPromTarget(
        expr='irate(node_netstat_Icmp_OutMsgs{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='ICMP transmitted'
      ))
      .addTarget(commonPromTarget(
        expr='irate(node_netstat_Icmp6_InMsgs{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='ICMP6 received'
      ))
      .addTarget(commonPromTarget(
        expr='irate(node_netstat_Icmp6_OutMsgs{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='ICMP6 transmitted'
      )),

    local networkNetstatICMPerrors =
      nodeTimeseries.new(
        title='ICMP errors rate'
      )
      .withUnits('err/s')
      .addTarget(commonPromTarget(
        expr='irate(node_netstat_Icmp_InErrors{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='ICMP Errors'
      ))
      .addTarget(commonPromTarget(
        expr='irate(node_netstat_Icmp6_InErrors{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='ICMP6 Errors'
      )),


    local rows =
      [
        row.new('Network')
        .addPanel(networkInterfacesTable { span: 12 })
        .addPanel(networkTrafficPanel { span: 6 })
        .addPanel(networkOperStatus { span: 6 })
        .addPanel(networkErrorsPanel { span: 6 })
        .addPanel(networkDropsPanel { span: 6 })
        .addPanel(networkPacketsPanel { span: 6 })
        .addPanel(networkMulticastPanel { span: 6 })
        .addPanel(networkFifoPanel { span: 6 })
        .addPanel(networkCompressedPanel { span: 6 })
        .addPanel(networkNFConntrack { span: 6 })
        .addPanel(networkSoftnetPanel { span: 6 })
        .addPanel(networkSoftnetSqueezePanel { span: 6 }),
        row.new('Network Sockets')
        .addPanel(networkSockstatAll { span: 12 })
        .addPanel(networkSockstatTCP { span: 6 })
        .addPanel(networkSockstatUDP { span: 6 })
        .addPanel(networkSockstatMemory { span: 6 })
        .addPanel(networkSockstatOther { span: 6 }),

        row.new('Network Netstat')
        .addPanel(networkNetstatIP { span: 12 })
        .addPanel(networkNetstatTCP { span: 6 })
        .addPanel(networkNetstatTCPerrors { span: 6 })
        .addPanel(networkNetstatUDP { span: 6 })
        .addPanel(networkNetstatUDPerrors { span: 6 })
        .addPanel(networkNetstatICMP { span: 6 })
        .addPanel(networkNetstatICMPerrors { span: 6 }),
      ],

    dashboard: if platform == 'Linux' then
      dashboard.new(
        '%sNode Network' % config { nodeQuerySelector: c.nodeQuerySelector }.dashboardNamePrefix,
        time_from=config.dashboardInterval,
        tags=(config.dashboardTags),
        timezone=config.dashboardTimezone,
        refresh=config.dashboardRefresh,
        graphTooltip='shared_crosshair',
        uid=config.grafanaDashboardIDs['nodes-network.json']
      )
      .addLink(c.links.fleetDash)
      .addLink(c.links.nodeDash)
      .addLink(c.links.otherDashes)
      .addAnnotations(c.annotations)
      .addTemplates(templates)
      .addRows(rows)
    else if platform == 'Darwin' then {},
  },
}
