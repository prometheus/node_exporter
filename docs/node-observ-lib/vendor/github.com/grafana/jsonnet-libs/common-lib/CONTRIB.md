
## Panels overview

All panels in this lib should implement one of the following methods: 

- `panel.new(title,targets,description)` - creates new panel. List of arguments could vary;
- `panel.stylize(allLayers=true)` - directly applies this panel style to existing panel. By default includes all layers of styles. To apply only top layer, set allLayers=false. This mode is useful to cherry-pick style layers to create new style combination.

Some other methods could be found such as:
- `panel.stylizeByRegexp(regexp)` - attaches style as panel overrides (by regexp);
- `panel.stylizeByName(name)` - attaches style as panel overrides (by name).

## Panels common groups

This library consists of multiple common groups of panels for widely used resources such as CPU, memory, disks and so on.

All of those groups inherit `generic` group as their base.

All panels inherit `generic/base.libsonnet` via `generic/<paneltype>/base.libsonnet`.

