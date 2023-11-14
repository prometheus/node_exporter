local g = import './g.libsonnet';

local timeSeries = g.panel.timeSeries;
local logsPanel = g.panel.logs;
local defaults = timeSeries.fieldConfig.defaults;
local custom = timeSeries.fieldConfig.defaults.custom;
local options = timeSeries.options;
local fieldConfig = timeSeries.fieldConfig;
function(
  logsVolumeTarget,
  logsTarget,
  logsVolumeGroupBy,
)

  {
    local this = self,

    logsVolumeInit(targets, title='Logs volume')::
      timeSeries.new(title)
      + timeSeries.queryOptions.withTargets(targets)
      + timeSeries.panelOptions.withDescription('Logs volume grouped by "%s" label.' % logsVolumeGroupBy)
      // set type to first target's type
      + timeSeries.queryOptions.withDatasource(
        logsVolumeTarget.datasource.type, logsVolumeTarget.datasource.uid
      )
      + custom.withDrawStyle('bars')
      + custom.stacking.withMode('normal')
      + custom.withFillOpacity(50)
      // should be set, otherwise interval is around 1s by default
      + timeSeries.queryOptions.withInterval('30s')
      + options.tooltip.withMode('multi')
      + options.tooltip.withSort('desc')
      + timeSeries.standardOptions.withUnit('none')
      + timeSeries.queryOptions.withTransformationsMixin(
        {
          id: 'renameByRegex',
          options: {
            regex: 'Value',
            renamePattern: 'logs',
          },
        }
      )
      + timeSeries.standardOptions.withOverridesMixin(
        [
          {
            matcher: {
              id: 'byRegexp',
              options: o.regex,
            },
            properties: [
              {
                id: 'color',
                value: {
                  mode: 'fixed',
                  fixedColor: o.color,
                },
              },
            ],
          }
          // https://grafana.com/docs/grafana/latest/explore/logs-integration/#log-level
          for o in
            [
              { regex: '(E|e)merg|(F|f)atal|(A|a)lert|(C|c)rit.*', color: 'purple' },
              { regex: '(E|e)(rr.*|RR.*)', color: 'red' },
              { regex: '(W|w)(arn.*|ARN.*|rn|RN)', color: 'orange' },
              { regex: '(N|n)(otice|ote)|(I|i)(nf.*|NF.*)', color: 'green' },
              { regex: 'dbg.*|DBG.*|(D|d)(EBUG|ebug)', color: 'blue' },
              { regex: '(T|t)(race|RACE)', color: 'light-blue' },
              { regex: 'logs', color: 'text' },
            ]
        ]
      ),

    logsInit(targets, title='Logs')::
      logsPanel.new(title)
      + logsPanel.queryOptions.withTargets(targets)
      + logsPanel.options.withDedupStrategy('exact')  //"none", "exact", "numbers", "signature"
      + logsPanel.options.withEnableLogDetails(true)
      + logsPanel.options.withShowTime(true)
      + logsPanel.options.withWrapLogMessage(true)
      + logsPanel.options.withPrettifyLogMessage(true),

    logsVolume: self.logsVolumeInit(logsVolumeTarget),
    logs: self.logsInit(logsTarget),

  }
