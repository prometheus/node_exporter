local d = import 'github.com/jsonnet-libs/docsonnet/doc-util/main.libsonnet';

{
  local this = self,

  '#setPanelIDs':: d.func.new(
    |||
      `setPanelIDs` ensures that all `panels` have a unique ID, this functions is used in
      `dashboard.withPanels` and `dashboard.withPanelsMixin` to provide a consistent
      experience.

      used in ../dashboard.libsonnet
    |||,
    args=[
      d.arg('panels', d.T.array),
    ]
  ),
  setPanelIDs(panels):
    local infunc(panels, start=1) =
      std.foldl(
        function(acc, panel)
          acc + {
            index:  // Track the index to ensure no duplicates exist.
              acc.index
              + 1
              + (if panel.type == 'row'
                    && 'panels' in panel
                 then std.length(panel.panels)
                 else 0),

            panels+: [
              panel + { id: acc.index }
              + (
                if panel.type == 'row'
                   && 'panels' in panel
                then {
                  panels:
                    infunc(
                      panel.panels,
                      acc.index + 1
                    ),
                }
                else {}
              ),
            ],
          },
        panels,
        { index: start, panels: [] }
      ).panels;
    infunc(panels),
}
