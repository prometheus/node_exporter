local g = import '../../../g.libsonnet';
local generic = import '../../generic/timeSeries/main.libsonnet';
local base = import './base.libsonnet';
base {

  new(
    title='Load average',
    loadTargets,
    cpuCountTarget,
    description=|||
      System load average over the previous 1, 5, and 15 minute ranges.

      A measurement of how many processes are waiting for CPU cycles. The maximum number is the number of CPU cores for the node.
    |||
  ):
    // validate inputs
    std.prune(
      {
        checks: [
          if !(std.objectHas(cpuCountTarget, 'legendFormat')) then error 'cpuCountTarget must have legendFormat"',
        ],
      }
    )
    +
    local targets = loadTargets + [cpuCountTarget];
    super.new(title, targets, description)
    // call directly threshold styler (not called from super automatically)
    + self.stylizeCpuCores(cpuCountTarget.legendFormat),

  stylizeCpuCores(cpuCountName):
    generic.threshold.stylizeByRegexp(cpuCountName),

  stylize(allLayers=true, cpuCountName=null):
    (if allLayers then super.stylize() else {})

    + g.panel.timeSeries.fieldConfig.defaults.custom.withFillOpacity(0)
    + g.panel.timeSeries.standardOptions.withMin(0)
    + g.panel.timeSeries.standardOptions.withUnit('short')
    // this is only called if cpuCountName provided
    + (if cpuCountName != null then self.stylizeCpuCores(cpuCountName) else {}),
}
