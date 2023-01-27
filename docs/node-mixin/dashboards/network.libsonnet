local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';
local dashboard = grafana.dashboard;
local row = grafana.row;
local prometheus = grafana.prometheus;
local template = grafana.template;
local graphPanel = grafana.graphPanel;
local grafana70 = import 'github.com/grafana/grafonnet-lib/grafonnet-7.0/grafana.libsonnet';
local gaugePanel = grafana70.panel.gauge;
local table = grafana70.panel.table;

local nodePanels = import 'panels-lib/panels.libsonnet';
local commonPanels = import 'panels-lib/common/panels.libsonnet';
local nodeTimeseries = nodePanels.timeseries;
local nodeTemplates = import '../lib/templates.libsonnet';
{

  new(config=null, platform=null):: {


//

    local networkTrafficPanel =
      nodeTimeseries.new(
        graphPanel.new(
          'Traffic',
          datasource='$datasource',
        )
      )
      .addTarget(prometheus.target(
        'irate(node_network_receive_bytes_total{%(nodeExporterSelector)s, instance="$instance",}[$__rate_interval])*8' % config,
        legendFormat='{{device}} received',
      ))
      .addTarget(prometheus.target(
        'irate(node_network_transmit_bytes_total{%(nodeExporterSelector)s, instance="$instance",}[$__rate_interval])*8' % config,
        legendFormat='{{device}} transmitted',
      ))
      .withDecimals(1)
      .withUnits("bps")
      .withNegativeYByRegex("transmit")
      .withAxisLabel("Transmit(-) | Receive(+)"),

    local networkPacketsPanel =
      nodeTimeseries.new(
        graphPanel.new(
          'Unicast Packets',
          datasource='$datasource',
        )
      )
      .addTarget(prometheus.target(
        'irate(node_network_receive_packets_total{%(nodeExporterSelector)s, instance="$instance",}[$__rate_interval])' % config,
        legendFormat='{{device}} received',
      ))
      .addTarget(prometheus.target(
        'irate(node_network_transmit_packets_total{%(nodeExporterSelector)s, instance="$instance",}[$__rate_interval])' % config,
        legendFormat='{{device}} transmitted',
      ))
      .withDecimals(1)
      .withUnits("pps")
      .withNegativeYByRegex("transmit")
      .withAxisLabel("Transmit(-) | Receive(+)"),

    local networkErrorsPanel =
      nodeTimeseries.new(
        graphPanel.new(
          'Errors',
          datasource='$datasource',
        )
      )
      .addTarget(prometheus.target(
        'irate(node_network_receive_errs_total{%(nodeExporterSelector)s, instance="$instance",}[$__rate_interval])' % config,
        legendFormat='{{device}} received',
      ))
      .addTarget(prometheus.target(
        'irate(node_network_transmit_errs_total{%(nodeExporterSelector)s, instance="$instance",}[$__rate_interval])' % config,
        legendFormat='{{device}} transmitted',
      ))
      .withDecimals(1)
      .withUnits("pps")
      .withNegativeYByRegex("transmit")
      .withAxisLabel("Transmit(-) | Receive(+)"),

    local networkDropsPanel =
      nodeTimeseries.new(
        graphPanel.new(
          'Dropped Packets',
          datasource='$datasource',
        )
      )
      .addTarget(prometheus.target(
        'irate(node_network_receive_drop_total{%(nodeExporterSelector)s, instance="$instance",}[$__rate_interval])' % config,
        legendFormat='{{device}} received',
      ))
      .addTarget(prometheus.target(
        'irate(node_network_transmit_drop_total{%(nodeExporterSelector)s, instance="$instance",}[$__rate_interval])' % config,
        legendFormat='{{device}} transmitted',
      ))
      .withDecimals(1)
      .withUnits("pps")
      .withNegativeYByRegex("transmit")
      .withAxisLabel("Transmit(-) | Receive(+)"),
    local networkCompressedPanel =
      nodeTimeseries.new(
        graphPanel.new(
          'Compressed Packets',
          datasource='$datasource',
        )
      )
      .addTarget(prometheus.target(
        'irate(node_network_receive_compressed_total{%(nodeExporterSelector)s, instance="$instance",}[$__rate_interval])' % config,
        legendFormat='{{device}} received',
      ))
      .addTarget(prometheus.target(
        'irate(node_network_transmit_compressed_total{%(nodeExporterSelector)s, instance="$instance",}[$__rate_interval])' % config,
        legendFormat='{{device}} transmitted',
      ))
      .withDecimals(1)
      .withUnits("pps")
      .withNegativeYByRegex("transmit")
      .withAxisLabel("Transmit(-) | Receive(+)"),

    local networkMulticastPanel =
      nodeTimeseries.new(
        graphPanel.new(
          'Multicast Packets',
          datasource='$datasource',
        )
      )
      .addTarget(prometheus.target(
        'irate(node_network_receive_multicast_total{%(nodeExporterSelector)s, instance="$instance",}[$__rate_interval])' % config,
        legendFormat='{{device}} received',
      ))
      .addTarget(prometheus.target(
        'irate(node_network_transmit_multicast_total{%(nodeExporterSelector)s, instance="$instance",}[$__rate_interval])' % config,
        legendFormat='{{device}} transmitted',
      ))
      .withDecimals(1)
      .withUnits("pps")
      .withNegativeYByRegex("transmit"),

    local networkFifoPanel =
      nodeTimeseries.new(
        graphPanel.new(
          'Network FIFO',
          datasource='$datasource',
        )
      )
      .addTarget(prometheus.target(
        'irate(node_network_receive_fifo_total{%(nodeExporterSelector)s, instance="$instance",}[$__rate_interval])' % config,
        legendFormat='{{device}} received',
      ))
      .addTarget(prometheus.target(
        'irate(node_network_transmit_fifo_total{%(nodeExporterSelector)s, instance="$instance",}[$__rate_interval])' % config,
        legendFormat='{{device}} transmitted',
      ))
      .withDecimals(1)
      .withUnits("pps")
      .withNegativeYByRegex("transmit")
      .withAxisLabel("Transmit(-) | Receive(+)"),

      local networkNFConntrack =
      nodeTimeseries.new(
        graphPanel.new(
          'NF Conntrack',
          datasource='$datasource',
        )
      )
      .addTarget(prometheus.target(
        'node_nf_conntrack_entries{%(nodeExporterSelector)s, instance="$instance"}' % config,
        legendFormat='NF conntrack entries',
      ))
      .addTarget(prometheus.target(
        'node_nf_conntrack_entries_limit{%(nodeExporterSelector)s, instance="$instance"}' % config,
        legendFormat='NF conntrack limits',
      ))
      .withFillOpacity(0),

      local networkSoftnetPanel =
      nodeTimeseries.new(
        graphPanel.new(
          'Softnet packets',
          datasource='$datasource',
        )
      )
      .addTarget(prometheus.target(
        'irate(node_softnet_processed_total{%(nodeExporterSelector)s, instance="$instance"}[$__rate_interval])' % config,
        legendFormat='CPU {{cpu }} proccessed',
      ))
      .addTarget(prometheus.target(
        'irate(node_softnet_dropped_total{%(nodeExporterSelector)s, instance="$instance"}[$__rate_interval])' % config,
        legendFormat='CPU {{cpu }} dropped',
      ))
      .withDecimals(1)
      .withUnits("pps")
      .withNegativeYByRegex("dropped")
      .withAxisLabel("Dropped(-) | Processed(+)"),

      local networkSoftnetSqueezePanel =
      nodeTimeseries.new(
        graphPanel.new(
          'Softnet Out of Quota',
          datasource='$datasource',
        )
      )
      .addTarget(prometheus.target(
        'irate(node_softnet_times_squeezed_total{%(nodeExporterSelector)s, instance="$instance"}[$__rate_interval])' % config,
        legendFormat='CPU {{cpu}} out of quota',
      ))
      .withDecimals(1)
      .withUnits("pps"),


    //softnet squueze
    //node_softnet_times_squeezed_total
    // CPU {{cpu}} out of quota

     local networkRow =
      row.new('Network')
      .addPanel(networkTrafficPanel + {span: 6})
      .addPanel(networkPacketsPanel+ {span: 6})

      
      .addPanel(networkErrorsPanel+ {span: 6})
      .addPanel(networkDropsPanel+ {span: 6})

      .addPanel(networkMulticastPanel+ {span: 6})
      .addPanel(networkFifoPanel+ {span: 6})

      .addPanel(networkCompressedPanel+ {span: 6})
      .addPanel(networkNFConntrack+ {span: 6})

      .addPanel(networkSoftnetPanel + {span: 6})
      .addPanel(networkSoftnetSqueezePanel + {span: 6})
      
      ,


    local rows =
      [
        networkRow,
      ],

    local templates = nodeTemplates.new(config=config, platform=platform).templates,

    dashboard: if platform == 'Linux' then
      dashboard.new(
        '%sNode Network' % config.dashboardNamePrefix,
        time_from=config.dashboardInterval,
        tags=(config.dashboardTags),
        timezone=config.dashboardTimezone,
        refresh=config.dashboardRefresh,
        graphTooltip='shared_crosshair',
        uid='node-network'
      ) { editable: true }
      .addTemplates(templates)
      .addRows(rows)
    else if platform == 'Darwin' then {},
  },
}
