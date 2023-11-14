local g = import '../../../g.libsonnet';
local base = import './base.libsonnet';
local statusHistory = g.panel.statusHistory;
base {
  new(title='NTP status', targets, description=''):
    super.new(title, targets, description),

  stylize(allLayers=true):
    (if allLayers then super.stylize() else {})
    + statusHistory.standardOptions.color.withMode('fixed')
    + statusHistory.standardOptions.withMappings(
      {
        type: 'value',
        options: {
          '0': {
            text: 'Not in sync',
            color: 'light-yellow',
            index: 1,
          },
          '1': {
            text: 'In sync',
            color: 'light-green',
            index: 0,
          },
        },
      }
    ),
}
