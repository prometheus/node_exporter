local g = import '../../../g.libsonnet';
local generic = import '../../generic/stat/main.libsonnet';
local base = import './base.libsonnet';
local stat = g.panel.stat;
// Uptime panel. expects duration in seconds as input
base {
  new(title='Uptime', targets, description='The duration of time that has passed since the last reboot or system start.'):
    super.new(title, targets, description)
    + stat.options.withReduceOptions({})
    + stat.options.reduceOptions.withCalcsMixin(
      [
        'lastNotNull',
      ]
    )
    + self.stylize(),
  stylize(allLayers=true):
    (if allLayers then super.stylize() else {})
    + stat.standardOptions.withDecimals(1)
    + stat.standardOptions.withUnit('dtdurations')
    + stat.options.withColorMode('value')
    + stat.options.withGraphMode('none')
    + stat.standardOptions.thresholds.withMode('absolute')
    + stat.standardOptions.thresholds.withSteps(
      [
        // Warn with orange color when uptime resets:
        stat.thresholdStep.withColor('orange')
        + stat.thresholdStep.withValue(null),
        // clear color after 10 minutes:
        stat.thresholdStep.withColor('text')
        + stat.thresholdStep.withValue(600),
      ]
    ),
}
