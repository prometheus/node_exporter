// This file is generated, do not manually edit.
{
  '#': { help: 'grafonnet.playlist', name: 'playlist' },
  '#withInterval': { 'function': { args: [{ default: '5m', enums: null, name: 'value', type: 'string' }], help: 'Interval sets the time between switching views in a playlist.\nFIXME: Is this based on a standardized format or what options are available? Can datemath be used?' } },
  withInterval(value='5m'): { interval: value },
  '#withItems': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'The ordered list of items that the playlist will iterate over.\nFIXME! This should not be optional, but changing it makes the godegen awkward' } },
  withItems(value): { items: (if std.isArray(value)
                              then value
                              else [value]) },
  '#withItemsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'The ordered list of items that the playlist will iterate over.\nFIXME! This should not be optional, but changing it makes the godegen awkward' } },
  withItemsMixin(value): { items+: (if std.isArray(value)
                                    then value
                                    else [value]) },
  items+:
    {
      '#': { help: '', name: 'items' },
      '#withTitle': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Title is an unused property -- it will be removed in the future' } },
      withTitle(value): { title: value },
      '#withType': { 'function': { args: [{ default: null, enums: ['dashboard_by_uid', 'dashboard_by_id', 'dashboard_by_tag'], name: 'value', type: 'string' }], help: 'Type of the item.' } },
      withType(value): { type: value },
      '#withValue': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Value depends on type and describes the playlist item.\n\n - dashboard_by_id: The value is an internal numerical identifier set by Grafana. This\n is not portable as the numerical identifier is non-deterministic between different instances.\n Will be replaced by dashboard_by_uid in the future. (deprecated)\n - dashboard_by_tag: The value is a tag which is set on any number of dashboards. All\n dashboards behind the tag will be added to the playlist.\n - dashboard_by_uid: The value is the dashboard UID' } },
      withValue(value): { value: value },
    },
  '#withName': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Name of the playlist.' } },
  withName(value): { name: value },
  '#withUid': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Unique playlist identifier. Generated on creation, either by the\ncreator of the playlist of by the application.' } },
  withUid(value): { uid: value },
}
