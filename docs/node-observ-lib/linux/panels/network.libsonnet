local g = import '../g.libsonnet';
local commonlib = import 'common-lib/common/main.libsonnet';
local utils = commonlib.utils;
{
  new(this):
    {
      local t = this.grafana.targets,
      local table = g.panel.table,
      local fieldOverride = g.panel.table.fieldOverride,
      local instanceLabel = this.config.instanceLabels[0],

      networkErrorsAndDroppedPerSec:
        commonlib.panels.network.timeSeries.errors.new(
          'Network errors and dropped packets',
          targets=std.map(
            function(t) t
                        {
              expr: t.expr + '>0',
            },
            [
              t.networkOutErrorsPerSec,
              t.networkInErrorsPerSec,
              t.networkOutDroppedPerSec,
              t.networkInDroppedPerSec,
            ]
          ),
          description=|||
            **Network errors**:

            Network errors refer to issues that occur during the transmission of data across a network. 

            These errors can result from various factors, including physical issues, jitter, collisions, noise and interference.

            Monitoring network errors is essential for diagnosing and resolving issues, as they can indicate problems with network hardware or environmental factors affecting network quality.

            **Dropped packets**:

            Dropped packets occur when data packets traveling through a network are intentionally discarded or lost due to congestion, resource limitations, or network configuration issues. 

            Common causes include network congestion, buffer overflows, QoS settings, and network errors, as corrupted or incomplete packets may be discarded by receiving devices.

            Dropped packets can impact network performance and lead to issues such as degraded voice or video quality in real-time applications.
          |||
        )
        + commonlib.panels.network.timeSeries.errors.withNegateOutPackets(),
      networkErrorsAndDroppedPerSecTopK:
        commonlib.panels.network.timeSeries.errors.new(
          'Network errors and dropped packets',
          targets=std.map(
            function(t) t
                        {
              expr: 'topk(25, ' + t.expr + ')>0',
              legendFormat: '{{' + this.config.instanceLabels[0] + '}}: ' + std.get(t, 'legendFormat', '{{ nic }}'),
            },
            [
              t.networkOutErrorsPerSec,
              t.networkInErrorsPerSec,
              t.networkOutDroppedPerSec,
              t.networkInDroppedPerSec,
            ]
          ),
          description=|||
            Top 25.

            **Network errors**:

            Network errors refer to issues that occur during the transmission of data across a network. 

            These errors can result from various factors, including physical issues, jitter, collisions, noise and interference.

            Monitoring network errors is essential for diagnosing and resolving issues, as they can indicate problems with network hardware or environmental factors affecting network quality.

            **Dropped packets**:

            Dropped packets occur when data packets traveling through a network are intentionally discarded or lost due to congestion, resource limitations, or network configuration issues. 

            Common causes include network congestion, buffer overflows, QoS settings, and network errors, as corrupted or incomplete packets may be discarded by receiving devices.

            Dropped packets can impact network performance and lead to issues such as degraded voice or video quality in real-time applications.
          |||
        )
        + g.panel.timeSeries.fieldConfig.defaults.custom.withDrawStyle('points')
        + g.panel.timeSeries.fieldConfig.defaults.custom.withPointSize(5),

      networkErrorsPerSec:
        commonlib.panels.network.timeSeries.errors.new(
          'Network errors',
          targets=[t.networkInErrorsPerSec, t.networkOutErrorsPerSec]
        )
        + commonlib.panels.network.timeSeries.errors.withNegateOutPackets(),
      networkDroppedPerSec:
        commonlib.panels.network.timeSeries.dropped.new(
          targets=[t.networkInDroppedPerSec, t.networkOutDroppedPerSec]
        )
        + commonlib.panels.network.timeSeries.errors.withNegateOutPackets(),
      networkUsagePerSec:
        commonlib.panels.network.timeSeries.traffic.new(
          targets=[t.networkInBitPerSecFiltered, t.networkOutBitPerSecFiltered]
        )
        + commonlib.panels.network.timeSeries.traffic.withNegateOutPackets(),
      networkPacketsPerSec:
        commonlib.panels.network.timeSeries.packets.new(
          targets=[t.networkInPacketsPerSec, t.networkOutPacketsPerSec]
        )
        + commonlib.panels.network.timeSeries.traffic.withNegateOutPackets(),
      networkMulticastPerSec:
        commonlib.panels.network.timeSeries.multicast.new(
          'Multicast packets',
          targets=[t.networkInMulticastPacketsPerSec, t.networkOutMulticastPacketsPerSec],
          description='Multicast packets received and transmitted.'
        )
        + commonlib.panels.network.timeSeries.traffic.withNegateOutPackets(),

      networkFifo:
        commonlib.panels.network.timeSeries.packets.new(
          'Network FIFO',
          targets=[t.networkFifoInPerSec, t.networkFifoOutPerSec],
          description=|||
            Network FIFO (First-In, First-Out) refers to a buffer used by the network stack to store packets in a queue.
            It is a mechanism used to manage network traffic and ensure that packets are delivered to their destination in the order they were received.
            Packets are stored in the FIFO buffer until they can be transmitted or processed further.
          |||
        )
        + commonlib.panels.network.timeSeries.traffic.withNegateOutPackets(),
      networkCompressedPerSec:
        commonlib.panels.network.timeSeries.packets.new(
          'Compressed packets',
          targets=[t.networkCompressedInPerSec, t.networkCompressedOutPerSec],
          description=|||
            - Compressed received: 
            Number of correctly received compressed packets. This counters is only meaningful for interfaces which support packet compression (e.g. CSLIP, PPP).

            - Compressed transmitted:
            Number of transmitted compressed packets. This counters is only meaningful for interfaces which support packet compression (e.g. CSLIP, PPP).

            https://docs.kernel.org/networking/statistics.html
          |||,
        )
        + commonlib.panels.network.timeSeries.traffic.withNegateOutPackets(),
      networkNFConntrack:
        commonlib.panels.generic.timeSeries.base.new(
          'NF conntrack',
          targets=[t.networkNFConntrackEntries, t.networkNFConntrackLimits],
          description=|||
            NF Conntrack is a component of the Linux kernel's netfilter framework that provides stateful packet inspection to track and manage network connections,
            enforce firewall rules, perform NAT, and manage network address/port translation.
          |||
        )
        + g.panel.timeSeries.fieldConfig.defaults.custom.withFillOpacity(0),

      networkSoftnet:
        commonlib.panels.network.timeSeries.packets.new(
          'Softnet packets',
          targets=[t.networkSoftnetProcessedPerSec, t.networkSoftnetDroppedPerSec],
          description=|||
            Softnet packets are received by the network and queued for processing by the kernel's networking stack.
            Softnet packets are usually generated by network traffic that is directed to the local host, and they are typically processed by the kernel's networking subsystem before being passed on to the relevant application. 
          |||
        )
        + commonlib.panels.network.timeSeries.traffic.withNegateOutPackets('/dropped/')
        + g.panel.timeSeries.fieldConfig.defaults.custom.withAxisLabel('Dropped(-) | Processed(+)'),
      networkSoftnetSqueeze:
        commonlib.panels.network.timeSeries.packets.new(
          'Softnet out of quota',
          targets=[t.networkSoftnetSqueezedPerSec],
          description=|||
            "Softnet out of quota" is a network-related metric in Linux that measures the number of times the kernel's softirq processing was unable to handle incoming network traffic due to insufficient softirq processing capacity.
            This means that the kernel has reached its processing capacity limit for incoming packets, and any additional packets will be dropped or deferred.
          |||
        ),
      networkOperStatus:
        commonlib.panels.network.statusHistory.interfaceStatus.new(
          'Network interfaces carrier status',
          targets=[t.networkCarrier],
          description='Network interfaces carrier status',
        ),
      networkOverviewTable:
        commonlib.panels.generic.table.base.new(
          'Network interfaces overview',
          targets=
          [
            t.networkUp
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
            + g.query.prometheus.withRefId('Up'),
            t.networkCarrier
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
            + g.query.prometheus.withRefId('Carrier'),
            t.networkOutBitPerSec
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(false)
            + g.query.prometheus.withRefId('Transmitted'),
            t.networkInBitPerSec
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(false)
            + g.query.prometheus.withRefId('Received'),
            t.networkArpEntries
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
            + g.query.prometheus.withRefId('ARP entries'),
            t.networkMtuBytes
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
            + g.query.prometheus.withRefId('MTU'),
            t.networkSpeedBitsPerSec
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
            + g.query.prometheus.withRefId('Speed'),
            t.networkTransmitQueueLength
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
            + g.query.prometheus.withRefId('Queue length'),
            t.networkInfo
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
            + g.query.prometheus.withRefId('Info'),
          ],
          description='Network interfaces overview.'
        )
        + g.panel.table.standardOptions.withOverridesMixin([
          fieldOverride.byName.new('Speed')
          + fieldOverride.byName.withPropertiesFromOptions(
            table.standardOptions.withUnit('bps')
          ),
        ])
        + g.panel.table.standardOptions.withOverridesMixin([
          fieldOverride.byRegexp.new('Transmitted|Received')
          + fieldOverride.byRegexp.withProperty('custom.displayMode', 'gradient-gauge')
          + fieldOverride.byRegexp.withPropertiesFromOptions(
            table.standardOptions.withUnit('bps')
            + table.standardOptions.color.withMode('continuous-BlYlRd')
            + table.standardOptions.withMax(1000 * 1000 * 100)
          ),
        ])
        + g.panel.table.standardOptions.withOverridesMixin([
          fieldOverride.byRegexp.new('Carrier|Up')
          + fieldOverride.byRegexp.withProperty('custom.displayMode', 'color-text')
          + fieldOverride.byRegexp.withPropertiesFromOptions(
            table.standardOptions.withMappings(
              {
                type: 'value',
                options: {
                  '0': {
                    text: 'Down',
                    color: 'light-red',
                    index: 0,
                  },
                  '1': {
                    text: 'Up',
                    color: 'light-green',
                    index: 1,
                  },
                },
              }
            ),
          ),
        ])
        + table.queryOptions.withTransformationsMixin(
          [
            {
              id: 'joinByField',
              options: {
                byField: 'device',
                mode: 'outer',
              },
            },
            {
              id: 'filterFieldsByName',
              options: {
                include: {
                  pattern: 'device|duplex|address|Value.+',
                },
              },
            },
            {
              id: 'renameByRegex',
              options: {
                regex: '(Value) #(.*)',
                renamePattern: '$2',
              },
            },
            {
              id: 'organize',
              options: {
                excludeByName: {
                  Info: true,
                },
                renameByName:
                  {
                    device: 'Interface',
                    duplex: 'Duplex',
                    address: 'Address',
                  },
              },
            },
            {
              id: 'organize',
              options: {
                indexByName: {
                  Interface: 0,
                  Up: 1,
                  Carrier: 2,
                  Received: 3,
                  Transmitted: 4,
                },
              },
            },
          ]
        ),
      networkSockstatAll:
        commonlib.panels.generic.timeSeries.base.new(
          'Sockets in use',
          targets=[t.networkSocketsUsed],
          description='Number of sockets currently in use.',
        ),

      networkSockstatTCP:
        commonlib.panels.generic.timeSeries.base.new(
          'Sockets TCP',
          targets=[t.networkSocketsTCPAllocated, t.networkSocketsTCPIPv4, t.networkSocketsTCPIPv6, t.networkSocketsTCPOrphans, t.networkSocketsTCPTimeWait],
          description=|||
            TCP sockets are used for establishing and managing network connections between two endpoints over the TCP/IP protocol.

            Orphan sockets: If a process terminates unexpectedly or is terminated without closing its sockets properly, the sockets may become orphaned.
          |||
        ),
      networkSockstatUDP:
        commonlib.panels.generic.timeSeries.base.new(
          'Sockets UDP',
          targets=[t.networkSocketsUDPLiteInUse, t.networkSocketsUDPInUse, t.networkSocketsUDPLiteIPv6InUse, t.networkSocketsUDPIPv6InUse],
          description=|||
            UDP (User Datagram Protocol) and UDPlite (UDP-Lite) sockets are used for transmitting and receiving data over the UDP and UDPlite protocols, respectively.
            Both UDP and UDPlite are connectionless protocols that do not provide a reliable data delivery mechanism.
          |||
        ),
      networkSockstatOther:
        commonlib.panels.generic.timeSeries.base.new(
          'Sockets other',
          targets=[t.networkSocketsFragInUse, t.networkSocketsFragIPv6InUse, t.networkSocketsRawInUse, t.networkSocketsIPv6RawInUse],
          description=|||
            FRAG (IP fragment) sockets: Used to receive and process fragmented IP packets. FRAG sockets are useful in network monitoring and analysis.

            RAW sockets: Allow applications to send and receive raw IP packets directly without the need for a transport protocol like TCP or UDP.
          |||
        ),
      networkSockstatMemory:
        local panel = g.panel.timeSeries;
        local override = g.panel.timeSeries.standardOptions.override;
        commonlib.panels.generic.timeSeries.base.new(
          title='Sockets memory',
          targets=[t.networkSocketsTCPMemoryPages, t.networkSocketsUDPMemoryPages, t.networkSocketsTCPMemoryBytes, t.networkSocketsUDPMemoryBytes],
          description=|||
            Memory currently in use for sockets.
          |||,
        )
        + panel.queryOptions.withMaxDataPoints(100)
        + panel.fieldConfig.defaults.custom.withAxisLabel('Pages')
        + panel.standardOptions.withOverridesMixin(
          panel.standardOptions.override.byRegexp.new('/bytes/')
          + override.byType.withPropertiesFromOptions(
            panel.standardOptions.withDecimals(2)
            + panel.standardOptions.withUnit('bytes')
            + panel.fieldConfig.defaults.custom.withDrawStyle('bars')
            + panel.fieldConfig.defaults.custom.withStacking(value={ mode: 'normal', group: 'A' })
            + panel.fieldConfig.defaults.custom.withAxisLabel('Bytes')
          )
        ),

      networkNetstatIP:
        local panel = g.panel.timeSeries;
        local override = g.panel.timeSeries.standardOptions.override;
        commonlib.panels.network.timeSeries.packets.new(
          'IP octets',
          targets=[t.networkNetstatIPInOctetsPerSec, t.networkNetstatIPOutOctetsPerSec],
          description='Rate of IP octets received and transmitted.'
        )
        + commonlib.panels.network.timeSeries.traffic.withNegateOutPackets()
        + panel.standardOptions.withUnit('oct/s'),

      networkNetstatTCP:
        local panel = g.panel.timeSeries;
        local override = g.panel.timeSeries.standardOptions.override;
        commonlib.panels.network.timeSeries.packets.new(
          'TCP segments',
          targets=[t.networkNetstatTCPInSegmentsPerSec, t.networkNetstatTCPOutSegmentsPerSec],
          description='Rate of TCP segments received and transmitted.'
        )
        + commonlib.panels.network.timeSeries.traffic.withNegateOutPackets()
        + panel.standardOptions.withUnit('seg/s'),

      networkNetstatTCPerrors:
        local panel = g.panel.timeSeries;
        local override = g.panel.timeSeries.standardOptions.override;
        commonlib.panels.network.timeSeries.errors.new(
          title='TCP errors rate',
          targets=[
            t.networkNetstatTCPOverflowPerSec,
            t.networkNetstatTCPListenDropsPerSec,
            t.networkNetstatTCPRetransPerSec,
            t.networkNetstatTCPRetransSegPerSec,
            t.networkNetstatTCPInWithErrorsPerSec,
            t.networkNetstatTCPOutWithRstPerSec,
          ],
          description='Rate of TCP errors.'
        )
        + panel.standardOptions.withUnit('err/s'),

      networkNetstatUDP:
        local panel = g.panel.timeSeries;
        local override = g.panel.timeSeries.standardOptions.override;
        commonlib.panels.network.timeSeries.packets.new(
          'UDP datagrams',
          targets=[
            t.networkNetstatIPInUDPPerSec,
            t.networkNetstatIPOutUDPPerSec,
            t.networkNetstatIPInUDP6PerSec,
            t.networkNetstatIPOutUDP6PerSec,
          ],
          description='Rate of UDP datagrams received and transmitted.'
        )
        + commonlib.panels.network.timeSeries.traffic.withNegateOutPackets()
        + panel.standardOptions.withUnit('dat/s'),

      networkNetstatUDPerrors:
        local panel = g.panel.timeSeries;
        local override = g.panel.timeSeries.standardOptions.override;
        commonlib.panels.network.timeSeries.errors.new(
          title='UDP errors rate',
          targets=[
            t.networkNetstatUDPLiteInErrorsPerSec,
            t.networkNetstatUDPInErrorsPerSec,
            t.networkNetstatUDP6InErrorsPerSec,
            t.networkNetstatUDPNoPortsPerSec,
            t.networkNetstatUDP6NoPortsPerSec,
            t.networkNetstatUDPRcvBufErrsPerSec,
            t.networkNetstatUDP6RcvBufErrsPerSec,
            t.networkNetstatUDPSndBufErrsPerSec,
            t.networkNetstatUDP6SndBufErrsPerSec,
          ],
          description='Rate of UDP errors.'
        )
        + panel.standardOptions.withUnit('err/s'),

      networkNetstatICMP:
        local panel = g.panel.timeSeries;
        local override = g.panel.timeSeries.standardOptions.override;
        commonlib.panels.network.timeSeries.packets.new(
          'ICMP messages',
          targets=[
            t.networkNetstatICMPInPerSec,
            t.networkNetstatICMPOutPerSec,
            t.networkNetstatICMP6InPerSec,
            t.networkNetstatICMP6OutPerSec,
          ],
          description="Rate of ICMP messages, like 'ping', received and transmitted."
        )
        + commonlib.panels.network.timeSeries.traffic.withNegateOutPackets()
        + panel.standardOptions.withUnit('msg/s'),

      networkNetstatICMPerrors:
        local panel = g.panel.timeSeries;
        local override = g.panel.timeSeries.standardOptions.override;
        commonlib.panels.network.timeSeries.errors.new(
          title='ICMP errors rate',
          targets=[
            t.networkNetstatICMPInErrorsPerSec,
            t.networkNetstatICM6PInErrorsPerSec,
          ],
          description='Rate of ICMP messages received and transmitted with errors.'
        )
        + panel.standardOptions.withUnit('err/s'),

    },
}
