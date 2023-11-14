local g = import '../../../g.libsonnet';
local base = import './base.libsonnet';
local statusHistory = g.panel.statusHistory;
base {
  new(title='Interface status', targets, description='Interfaces statuses'):
    super.new(title, targets, description),

  stylize(allLayers=true):
    (if allLayers then super.stylize() else {})
    + statusHistory.standardOptions.color.withMode('fixed')
    + statusHistory.options.withShowValue('never')
    + statusHistory.standardOptions.withMappings(
      {
        type: 'value',
        options: {
          '0': {
            text: 'Down',
            color: 'light-red',
            index: 0,
          },
          '1': {
            text: 'Up',
            color: 'light-green',
            index: 1,
          },
        },
      }
    ),
}
