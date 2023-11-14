// This file is generated, do not manually edit.
{
  '#': { help: 'grafonnet.alerting.muteTiming', name: 'muteTiming' },
  '#withName': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
  withName(value): { name: value },
  '#withTimeIntervals': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
  withTimeIntervals(value): { time_intervals: (if std.isArray(value)
                                               then value
                                               else [value]) },
  '#withTimeIntervalsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
  withTimeIntervalsMixin(value): { time_intervals+: (if std.isArray(value)
                                                     then value
                                                     else [value]) },
  time_intervals+:
    {
      '#': { help: '', name: 'time_intervals' },
      '#withDaysOfMonth': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
      withDaysOfMonth(value): { days_of_month: (if std.isArray(value)
                                                then value
                                                else [value]) },
      '#withDaysOfMonthMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
      withDaysOfMonthMixin(value): { days_of_month+: (if std.isArray(value)
                                                      then value
                                                      else [value]) },
      '#withLocation': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
      withLocation(value): { location: value },
      '#withMonths': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
      withMonths(value): { months: (if std.isArray(value)
                                    then value
                                    else [value]) },
      '#withMonthsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
      withMonthsMixin(value): { months+: (if std.isArray(value)
                                          then value
                                          else [value]) },
      '#withTimes': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
      withTimes(value): { times: (if std.isArray(value)
                                  then value
                                  else [value]) },
      '#withTimesMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
      withTimesMixin(value): { times+: (if std.isArray(value)
                                        then value
                                        else [value]) },
      times+:
        {
          '#': { help: '', name: 'times' },
          '#withFrom': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
          withFrom(value): { from: value },
          '#withTo': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
          withTo(value): { to: value },
        },
      '#withWeekdays': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
      withWeekdays(value): { weekdays: (if std.isArray(value)
                                        then value
                                        else [value]) },
      '#withWeekdaysMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
      withWeekdaysMixin(value): { weekdays+: (if std.isArray(value)
                                              then value
                                              else [value]) },
      '#withYears': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
      withYears(value): { years: (if std.isArray(value)
                                  then value
                                  else [value]) },
      '#withYearsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
      withYearsMixin(value): { years+: (if std.isArray(value)
                                        then value
                                        else [value]) },
    },
}
