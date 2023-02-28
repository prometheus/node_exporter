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
        title='Network Interfaces Operational Status'
      )
      .withColor(mode='palette-classic')
      .withFillOpacity(100)
      .withLegend(mode='list')
      .addTarget(commonPromTarget(
        expr='node_network_up{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
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
      ],

    dashboard: if platform == 'Linux' then
      dashboard.new(
        '%sNode Network' % config { nodeQuerySelector: c.nodeQuerySelector }.dashboardNamePrefix,
        time_from=config.dashboardInterval,
        tags=(config.dashboardTags),
        timezone=config.dashboardTimezone,
        refresh=config.dashboardRefresh,
        graphTooltip='shared_crosshair',
        uid='node-network'
      ) { editable: true }
      .addLink(c.links.fleetDash)
      .addLink(c.links.nodeDash)
      .addLink(c.links.otherDashes)
      .addTemplates(templates)
      .addRows(rows)
    else if platform == 'Darwin' then {},
  },
}
