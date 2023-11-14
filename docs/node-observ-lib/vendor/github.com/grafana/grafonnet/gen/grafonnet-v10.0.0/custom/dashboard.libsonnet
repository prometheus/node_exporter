local util = import './util/main.libsonnet';
local d = import 'github.com/jsonnet-libs/docsonnet/doc-util/main.libsonnet';

{
  '#new':: d.func.new(
    'Creates a new dashboard with a title.',
    args=[d.arg('title', d.T.string)]
  ),
  new(title):
    self.withTitle(title)
    + self.withSchemaVersion()
    + self.withTimezone('utc')
    + self.time.withFrom('now-6h')
    + self.time.withTo('now'),

  withPanels(value): {
    _panels:: if std.isArray(value) then value else [value],
    panels: util.panel.setPanelIDs(self._panels),
  },
  withPanelsMixin(value): {
    _panels+:: if std.isArray(value) then value else [value],
    panels: util.panel.setPanelIDs(self._panels),
  },

  graphTooltip+: {
    // 0 - Default
    // 1 - Shared crosshair
    // 2 - Shared tooltip
    '#withSharedCrosshair':: d.func.new(
      'Share crosshair on all panels.',
    ),
    withSharedCrosshair():
      { graphTooltip: 1 },

    '#withSharedTooltip':: d.func.new(
      'Share crosshair and tooltip on all panels.',
    ),
    withSharedTooltip():
      { graphTooltip: 2 },
  },
}
+ (import './dashboard/annotation.libsonnet')
+ (import './dashboard/link.libsonnet')
+ (import './dashboard/variable.libsonnet')
