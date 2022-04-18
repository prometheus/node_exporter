local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';
local row = grafana.row;
local template = grafana.template;
local logPanel = grafana.logPanel;
local loki = grafana.loki;

{
  new(dashboardSelector=null):: {
    lokiDatasourceTemplate:: {
      current:
        {
          text: 'Loki',
          value: 'Loki',
        },
      name: 'loki_datasource',
      label: 'Loki Data Source',
      options: [],
      query: 'loki',
      hide: 0,
      refresh: 1,
      regex: '',
      type: 'datasource',
    },

    unitTemplate:: template.new(
      'unit',
      '$loki_datasource',
      'label_values(unit)',
      label='Systemd Unit',
      refresh='time',
      includeAll=true,
      multi=true,
      allValues='.+',
    ),

    errlog::
      logPanel.new(
        'Errors in system logs',
        datasource='$loki_datasource',
      )
      .addTargets(
        [
          loki.target('{level=~"err|crit|alert|emerg",unit=~"$unit", %s}' % dashboardSelector),
          loki.target('{filename=~"/var/log/syslog*|/var/log/messages*", %s} |~".+(?i)error(?-i).+"' % dashboardSelector),
        ]
      ),

    authlog::
      logPanel.new(
        'authlog',
        datasource='$loki_datasource',
      )
      .addTargets(
        [
          loki.target('{unit="ssh.service", %s}' % dashboardSelector),
          loki.target('{filename=~"/var/log/auth.log|/var/log/secure", %s}' % dashboardSelector),
        ]
      ),

    kernellog::
      logPanel.new(
        'Kernel logs',
        datasource='$loki_datasource',
      )
      .addTargets(
        [
          loki.target('{transport="kernel", %s}' % dashboardSelector),
          loki.target('{filename="/var/log/kern.log", %s}' % dashboardSelector),
        ]

      ),

    alllogs::
      logPanel.new(
        'All system logs',
        datasource='$loki_datasource',
      )
      .addTargets(
        [
          loki.target('{transport!="", %s}' % dashboardSelector),
          loki.target('{filename!="", %s}' % dashboardSelector),
        ]
      ),

    lokiLogRow::
      row.new(
        'Logs'
      )
      .addPanel(self.errlog)
      .addPanel(self.authlog)
      .addPanel(self.kernellog)
      .addPanel(self.alllogs),

    rows::
      [
        self.lokiLogRow,
      ],

    templates::
      [
        self.lokiDatasourceTemplate,
        self.unitTemplate,
      ],

  },

}
