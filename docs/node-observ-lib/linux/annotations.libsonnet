local g = import '../g.libsonnet';
local commonlib = import 'common-lib/common/main.libsonnet';
{
  new(this):
    {
      reboot:
        commonlib.annotations.reboot.new(
          title='Reboot',
          target=this.grafana.targets.reboot,
          instanceLabels=std.join(',', this.config.instanceLabels),
        )
        + commonlib.annotations.base.withTagKeys(std.join(',', this.config.groupLabels + this.config.instanceLabels)),
      memoryOOM:
        commonlib.annotations.base.new(
          'OOMkill',
          this.grafana.targets.memoryOOMkiller
        )
        + commonlib.annotations.base.withTagKeys(std.join(',', this.config.groupLabels + this.config.instanceLabels))
        + commonlib.annotations.base.withTextFormat('')
          {
          hide: true,
          iconColor: 'light-purple',
        },
      kernelUpdate:
        commonlib.annotations.base.new(
          'Kernel update',
          this.grafana.targets.kernelUpdate
        )
        + commonlib.annotations.base.withTagKeys(std.join(',', this.config.groupLabels + this.config.instanceLabels))
        + commonlib.annotations.base.withTextFormat('')
          {
          hide: true,
          iconColor: 'light-blue',
          step: '5m',
        },
    }
    +
    if
      this.config.enableLokiLogs
    then
      {
        serviceFailed: commonlib.annotations.serviceFailed.new(
                         title='Service failed',
                         target=this.grafana.targets.serviceFailed,
                       )
                       + commonlib.annotations.base.withTagKeys(std.join(',', this.config.groupLabels + this.config.instanceLabels + ['level'])),
        criticalEvents: commonlib.annotations.fatal.new(
                          title='Critical system event',
                          target=this.grafana.targets.criticalEvents,
                        )
                        + commonlib.annotations.base.withTagKeys(std.join(',', this.config.groupLabels + this.config.instanceLabels + ['level'])),
        sessionOpened:
          commonlib.annotations.base.new(
            title='Session opened',
            target=this.grafana.targets.sessionOpened,
          )
          + commonlib.annotations.base.withTagKeys(std.join(',', this.config.groupLabels + this.config.instanceLabels + ['level']))
            { hide: true },
        sessionClosed:
          commonlib.annotations.base.new(
            title='Session closed',
            target=this.grafana.targets.sessionOpened,
          )
          + commonlib.annotations.base.withTagKeys(std.join(',', this.config.groupLabels + this.config.instanceLabels + ['level']))
            { hide: true },
      }
    else
      {},
}
