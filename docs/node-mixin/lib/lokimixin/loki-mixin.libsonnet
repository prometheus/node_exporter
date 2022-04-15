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

    syslog::
      logPanel.new(
        'syslog Errors',
        datasource='$loki_datasource',
      )
      .addTarget(
        loki.target('{filename=~"/var/log/syslog*|/var/log/messages*", %s} |~".+(?i)error(?-i).+"' % dashboardSelector)
      ),

    authlog::
      logPanel.new(
        'authlog',
        datasource='$loki_datasource',
      )
      .addTarget(
        loki.target('{filename=~"/var/log/auth.log|/var/log/secure", %s}' % dashboardSelector)
      ),

    kernellog::
      logPanel.new(
        'Kernel logs',
        datasource='$loki_datasource',
      )
      .addTarget(
        loki.target('{filename=~"/var/log/kern.log*", %s}' % dashboardSelector)
      ),

    journalsyslog::
      logPanel.new(
        'Journal syslogs',
        datasource='$loki_datasource',
      )
      .addTarget(
        loki.target('{transport="syslog", %s}' % dashboardSelector)
      ),

    journalkernel::
      logPanel.new(
        'Journal Kernel logs',
        datasource='$loki_datasource',
      )
      .addTarget(
        loki.target('{transport="kernel", %s}' % dashboardSelector)
      ),

    journalstdout::
      logPanel.new(
        'Journal stdout Errors',
        datasource='$loki_datasource',
      )
      .addTarget(
        loki.target('{transport="stdout", %s, unit=~"$unit"} |~".+(?i)error(?-i).+"' % dashboardSelector)
      ),

    lokiDirectLogRow::
      row.new(
        'Loki Direct Log Scrapes'
      )
      .addPanel(self.syslog)
      .addPanel(self.authlog)
      .addPanel(self.kernellog),

    lokiJournalLogRow::
      row.new(
        'Loki Journal Log Scrapes'
      )
      .addPanel(self.journalsyslog)
      .addPanel(self.journalkernel)
      .addPanel(self.journalstdout),

    rows::
      [
        self.lokiDirectLogRow,
        self.lokiJournalLogRow,
      ],

    templates::
      [
        self.lokiDatasourceTemplate,
        self.unitTemplate,
      ],

  },

}
