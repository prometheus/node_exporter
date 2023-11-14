// This file is generated, do not manually edit.
{
  '#': { help: 'grafonnet.alerting.notificationPolicy', name: 'notificationPolicy' },
  '#withContinue': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
  withContinue(value=true): { continue: value },
  '#withGroupBy': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
  withGroupBy(value): { group_by: (if std.isArray(value)
                                   then value
                                   else [value]) },
  '#withGroupByMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
  withGroupByMixin(value): { group_by+: (if std.isArray(value)
                                         then value
                                         else [value]) },
  '#withGroupInterval': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
  withGroupInterval(value): { group_interval: value },
  '#withGroupWait': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
  withGroupWait(value): { group_wait: value },
  '#withMatch': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'Deprecated. Remove before v1.0 release.' } },
  withMatch(value): { match: value },
  '#withMatchMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'Deprecated. Remove before v1.0 release.' } },
  withMatchMixin(value): { match+: value },
  '#withMatchRe': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'MatchRegexps represents a map of Regexp.' } },
  withMatchRe(value): { match_re: value },
  '#withMatchReMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'MatchRegexps represents a map of Regexp.' } },
  withMatchReMixin(value): { match_re+: value },
  '#withMatchers': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'Matchers is a slice of Matchers that is sortable, implements Stringer, and\nprovides a Matches method to match a LabelSet against all Matchers in the\nslice. Note that some users of Matchers might require it to be sorted.' } },
  withMatchers(value): { matchers: (if std.isArray(value)
                                    then value
                                    else [value]) },
  '#withMatchersMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'Matchers is a slice of Matchers that is sortable, implements Stringer, and\nprovides a Matches method to match a LabelSet against all Matchers in the\nslice. Note that some users of Matchers might require it to be sorted.' } },
  withMatchersMixin(value): { matchers+: (if std.isArray(value)
                                          then value
                                          else [value]) },
  matchers+:
    {
      '#': { help: '', name: 'matchers' },
      '#withName': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
      withName(value): { Name: value },
      '#withType': { 'function': { args: [{ default: null, enums: ['=', '!=', '=~', '!~'], name: 'value', type: 'string' }], help: 'MatchType is an enum for label matching types.' } },
      withType(value): { Type: value },
      '#withValue': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
      withValue(value): { Value: value },
    },
  '#withMuteTimeIntervals': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
  withMuteTimeIntervals(value): { mute_time_intervals: (if std.isArray(value)
                                                        then value
                                                        else [value]) },
  '#withMuteTimeIntervalsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
  withMuteTimeIntervalsMixin(value): { mute_time_intervals+: (if std.isArray(value)
                                                              then value
                                                              else [value]) },
  '#withObjectMatchers': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'Matchers is a slice of Matchers that is sortable, implements Stringer, and\nprovides a Matches method to match a LabelSet against all Matchers in the\nslice. Note that some users of Matchers might require it to be sorted.' } },
  withObjectMatchers(value): { object_matchers: (if std.isArray(value)
                                                 then value
                                                 else [value]) },
  '#withObjectMatchersMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'Matchers is a slice of Matchers that is sortable, implements Stringer, and\nprovides a Matches method to match a LabelSet against all Matchers in the\nslice. Note that some users of Matchers might require it to be sorted.' } },
  withObjectMatchersMixin(value): { object_matchers+: (if std.isArray(value)
                                                       then value
                                                       else [value]) },
  object_matchers+:
    {
      '#': { help: '', name: 'object_matchers' },
      '#withName': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
      withName(value): { Name: value },
      '#withType': { 'function': { args: [{ default: null, enums: ['=', '!=', '=~', '!~'], name: 'value', type: 'string' }], help: 'MatchType is an enum for label matching types.' } },
      withType(value): { Type: value },
      '#withValue': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
      withValue(value): { Value: value },
    },
  '#withProvenance': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
  withProvenance(value): { provenance: value },
  '#withReceiver': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
  withReceiver(value): { receiver: value },
  '#withRepeatInterval': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
  withRepeatInterval(value): { repeat_interval: value },
  '#withRoutes': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
  withRoutes(value): { routes: (if std.isArray(value)
                                then value
                                else [value]) },
  '#withRoutesMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
  withRoutesMixin(value): { routes+: (if std.isArray(value)
                                      then value
                                      else [value]) },
}
