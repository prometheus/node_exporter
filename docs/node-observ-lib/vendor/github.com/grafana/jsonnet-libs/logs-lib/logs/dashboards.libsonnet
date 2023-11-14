local g = import './g.libsonnet';

function(
  title,
  showLogsVolume,
  panels,
  variables,
)
  {
    local this = self,
    logs:
      g.dashboard.new(title)
      + g.dashboard.withUid(g.util.string.slugify(title))
      + g.dashboard.withVariables(variables.toArray)
      + g.dashboard.withPanels(
        (
          if showLogsVolume then
            [panels.logsVolume
             + g.panel.timeSeries.gridPos.withH(6)
             + g.panel.timeSeries.gridPos.withW(24)]
          else []
        )
        +
        [
          panels.logs
          + g.panel.logs.gridPos.withH(18)
          + g.panel.logs.gridPos.withW(24),
        ]
      ),
  }
