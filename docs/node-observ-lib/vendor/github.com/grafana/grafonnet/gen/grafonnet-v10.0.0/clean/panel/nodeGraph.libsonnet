// This file is generated, do not manually edit.
(import '../../clean/panel.libsonnet')
+ {
  '#': { help: 'grafonnet.panel.nodeGraph', name: 'nodeGraph' },
  panelOptions+:
    {
      '#withType': { 'function': { args: [], help: '' } },
      withType(): { type: 'nodeGraph' },
    },
  options+:
    {
      '#withEdges': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
      withEdges(value): { options+: { edges: value } },
      '#withEdgesMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
      withEdgesMixin(value): { options+: { edges+: value } },
      edges+:
        {
          '#withMainStatUnit': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Unit for the main stat to override what ever is set in the data frame.' } },
          withMainStatUnit(value): { options+: { edges+: { mainStatUnit: value } } },
          '#withSecondaryStatUnit': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Unit for the secondary stat to override what ever is set in the data frame.' } },
          withSecondaryStatUnit(value): { options+: { edges+: { secondaryStatUnit: value } } },
        },
      '#withNodes': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
      withNodes(value): { options+: { nodes: value } },
      '#withNodesMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
      withNodesMixin(value): { options+: { nodes+: value } },
      nodes+:
        {
          '#withArcs': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'Define which fields are shown as part of the node arc (colored circle around the node).' } },
          withArcs(value): { options+: { nodes+: { arcs: (if std.isArray(value)
                                                          then value
                                                          else [value]) } } },
          '#withArcsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'Define which fields are shown as part of the node arc (colored circle around the node).' } },
          withArcsMixin(value): { options+: { nodes+: { arcs+: (if std.isArray(value)
                                                                then value
                                                                else [value]) } } },
          arcs+:
            {
              '#': { help: '', name: 'arcs' },
              '#withColor': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'The color of the arc.' } },
              withColor(value): { color: value },
              '#withField': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Field from which to get the value. Values should be less than 1, representing fraction of a circle.' } },
              withField(value): { field: value },
            },
          '#withMainStatUnit': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Unit for the main stat to override what ever is set in the data frame.' } },
          withMainStatUnit(value): { options+: { nodes+: { mainStatUnit: value } } },
          '#withSecondaryStatUnit': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Unit for the secondary stat to override what ever is set in the data frame.' } },
          withSecondaryStatUnit(value): { options+: { nodes+: { secondaryStatUnit: value } } },
        },
    },
}
+ { panelOptions+: { '#withType':: {} } }
