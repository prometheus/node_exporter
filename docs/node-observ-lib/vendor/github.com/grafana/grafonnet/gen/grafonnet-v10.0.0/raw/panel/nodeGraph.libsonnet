// This file is generated, do not manually edit.
{
  '#': { help: 'grafonnet.panel.nodeGraph', name: 'nodeGraph' },
  '#withArcOption': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withArcOption(value): { ArcOption: value },
  '#withArcOptionMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withArcOptionMixin(value): { ArcOption+: value },
  ArcOption+:
    {
      '#withColor': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'The color of the arc.' } },
      withColor(value): { ArcOption+: { color: value } },
      '#withField': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Field from which to get the value. Values should be less than 1, representing fraction of a circle.' } },
      withField(value): { ArcOption+: { field: value } },
    },
  '#withEdgeOptions': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withEdgeOptions(value): { EdgeOptions: value },
  '#withEdgeOptionsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withEdgeOptionsMixin(value): { EdgeOptions+: value },
  EdgeOptions+:
    {
      '#withMainStatUnit': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Unit for the main stat to override what ever is set in the data frame.' } },
      withMainStatUnit(value): { EdgeOptions+: { mainStatUnit: value } },
      '#withSecondaryStatUnit': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Unit for the secondary stat to override what ever is set in the data frame.' } },
      withSecondaryStatUnit(value): { EdgeOptions+: { secondaryStatUnit: value } },
    },
  '#withNodeOptions': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withNodeOptions(value): { NodeOptions: value },
  '#withNodeOptionsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withNodeOptionsMixin(value): { NodeOptions+: value },
  NodeOptions+:
    {
      '#withArcs': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'Define which fields are shown as part of the node arc (colored circle around the node).' } },
      withArcs(value): { NodeOptions+: { arcs: (if std.isArray(value)
                                                then value
                                                else [value]) } },
      '#withArcsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'Define which fields are shown as part of the node arc (colored circle around the node).' } },
      withArcsMixin(value): { NodeOptions+: { arcs+: (if std.isArray(value)
                                                      then value
                                                      else [value]) } },
      arcs+:
        {
          '#': { help: '', name: 'arcs' },
          '#withColor': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'The color of the arc.' } },
          withColor(value): { color: value },
          '#withField': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Field from which to get the value. Values should be less than 1, representing fraction of a circle.' } },
          withField(value): { field: value },
        },
      '#withMainStatUnit': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Unit for the main stat to override what ever is set in the data frame.' } },
      withMainStatUnit(value): { NodeOptions+: { mainStatUnit: value } },
      '#withSecondaryStatUnit': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Unit for the secondary stat to override what ever is set in the data frame.' } },
      withSecondaryStatUnit(value): { NodeOptions+: { secondaryStatUnit: value } },
    },
  '#withOptions': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withOptions(value): { options: value },
  '#withOptionsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withOptionsMixin(value): { options+: value },
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
  '#withType': { 'function': { args: [], help: '' } },
  withType(): { type: 'nodeGraph' },
}
