local g = import '../../g.libsonnet';
local commonlib = import 'common-lib/common/main.libsonnet';
local utils = commonlib.utils;
local xtd = import 'github.com/jsonnet-libs/xtd/main.libsonnet';
{
  new(this):
    {
      local t = this.grafana.targets,
      local table = g.panel.table,
      local fieldOverride = g.panel.table.fieldOverride,
      local instanceLabel = xtd.array.slice(this.config.instanceLabels, -1)[0],

      uptime: commonlib.panels.system.stat.uptime.new(targets=[t.system.uptime]),

      systemLoad:
        commonlib.panels.system.timeSeries.loadAverage.new(
          loadTargets=[t.system.systemLoad1, t.system.systemLoad5, t.system.systemLoad15],
          cpuCountTarget=t.cpu.cpuCount,
        ),

      systemContextSwitchesAndInterrupts:
        commonlib.panels.generic.timeSeries.base.new(
          'Context switches/Interrupts',
          targets=[
            t.system.systemContextSwitches,
            t.system.systemInterrupts,
          ],
          description=|||
            Context switches occur when the operating system switches from running one process to another. Interrupts are signals sent to the CPU by external devices to request its attention.

            A high number of context switches or interrupts can indicate that the system is overloaded or that there are problems with specific devices or processes.
          |||
        ),

      timeNtpStatus:
        commonlib.panels.system.statusHistory.ntp.new(
          'NTP status',
          targets=[t.system.timeNtpStatus],
          description='Status of time synchronization.'
        )
        + g.panel.timeSeries.standardOptions.withNoValue('No data.')
        + g.panel.statusHistory.options.withLegend(false),
      timeSyncDrift:
        commonlib.panels.generic.timeSeries.base.new(
          'Time synchronized drift',
          targets=[
            t.system.timeEstimatedError,
            t.system.timeOffset,
            t.system.timeMaxError,
          ],
          description=|||
            Time synchronization is essential to ensure accurate timekeeping, which is critical for many system operations such as logging, authentication, and network communication, as well as distributed systems or clusters where data consistency is important.
          |||
        )
        + g.panel.timeSeries.standardOptions.withUnit('seconds')
        + g.panel.timeSeries.standardOptions.withNoValue('No data.'),
      osInfo: commonlib.panels.generic.stat.info.new(
        'OS',
        targets=[t.system.osInfo],
        description='Operating system'
      )
              { options+: { reduceOptions+: { fields: '/^pretty_name$/' } } },
      kernelVersion:
        commonlib.panels.generic.stat.info.new('Kernel version',
                                               targets=[t.system.unameInfo],
                                               description='Kernel version of linux host.')
        { options+: { reduceOptions+: { fields: '/^release$/' } } },
      osTimezone:
        commonlib.panels.generic.stat.info.new(
          'Timezone', targets=[t.system.osTimezone], description='Current system timezone.'
        )
        { options+: { reduceOptions+: { fields: '/^time_zone$/' } } },
      hostname:
        commonlib.panels.generic.stat.info.new(
          'Hostname',
          targets=[t.system.unameInfo],
          description="System's hostname."
        )
        { options+: { reduceOptions+: { fields: '/^nodename$/' } } },

    },
}
