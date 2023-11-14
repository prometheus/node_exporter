local d = import 'doc-util/main.libsonnet';

{
  '#': d.pkg(
    name='date',
    url='github.com/jsonnet-libs/xtd/date.libsonnet',
    help='`time` provides various date related functions.',
  ),

  // Lookup tables for calendar calculations
  local commonYearMonthLength = [31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31],
  local commonYearMonthOffset = [0, 3, 3, 6, 1, 4, 6, 2, 5, 0, 3, 5],
  local leapYearMonthOffset = [0, 3, 4, 0, 2, 5, 0, 3, 6, 1, 4, 6],

  // monthOffset looks up the offset to apply in day of week calculations based on the year and month
  local monthOffset(year, month) =
    if self.isLeapYear(year)
    then leapYearMonthOffset[month - 1]
    else commonYearMonthOffset[month - 1],

  '#isLeapYear': d.fn(
    '`isLeapYear` returns true if the given year is a leap year.',
    [d.arg('year', d.T.number)],
  ),
  isLeapYear(year):: year % 4 == 0 && (year % 100 != 0 || year % 400 == 0),

  '#dayOfWeek': d.fn(
    '`dayOfWeek` returns the day of the week for the given date. 0=Sunday, 1=Monday, etc.',
    [
      d.arg('year', d.T.number),
      d.arg('month', d.T.number),
      d.arg('day', d.T.number),
    ],
  ),
  dayOfWeek(year, month, day)::
    (day + monthOffset(year, month) + 5 * ((year - 1) % 4) + 4 * ((year - 1) % 100) + 6 * ((year - 1) % 400)) % 7,

  '#dayOfYear': d.fn(
    |||
      `dayOfYear` calculates the ordinal day of the year based on the given date. The range of outputs is 1-365
      for common years, and 1-366 for leap years.
    |||,
    [
      d.arg('year', d.T.number),
      d.arg('month', d.T.number),
      d.arg('day', d.T.number),
    ],
  ),
  dayOfYear(year, month, day)::
    std.foldl(
      function(a, b) a + b,
      std.slice(commonYearMonthLength, 0, month - 1, 1),
      0
    ) + day +
    if month > 2 && self.isLeapYear(year)
    then 1
    else 0,

  // yearSeconds returns the number of seconds in the given year.
  local yearSeconds(year) = (
    if $.isLeapYear(year)
    then 366 * 24 * 3600
    else 365 * 24 * 3600
  ),

  // monthSeconds returns the number of seconds in the given month of a given year.
  local monthSeconds(year, month) = (
    commonYearMonthLength[month - 1] * 24 * 3600
    + if month == 2 && $.isLeapYear(year) then 86400 else 0
  ),

  // sumYearsSeconds returns the number of seconds in all years since 1970 up to year-1.
  local sumYearsSeconds(year) = std.foldl(
    function(acc, y) acc + yearSeconds(y),
    std.range(1970, year - 1),
    0,
  ),

  // sumMonthsSeconds returns the number of seconds in all months up to month-1 of the given year.
  local sumMonthsSeconds(year, month) = std.foldl(
    function(acc, m) acc + monthSeconds(year, m),
    std.range(1, month - 1),
    0,
  ),

  // sumDaysSeconds returns the number of seconds in all days up to day-1.
  local sumDaysSeconds(day) = (day - 1) * 24 * 3600,

  '#toUnixTimestamp': d.fn(
    |||
      `toUnixTimestamp` calculates the unix timestamp of a given date.
    |||,
    [
      d.arg('year', d.T.number),
      d.arg('month', d.T.number),
      d.arg('day', d.T.number),
      d.arg('hour', d.T.number),
      d.arg('minute', d.T.number),
      d.arg('second', d.T.number),
    ],
  ),
  toUnixTimestamp(year, month, day, hour, minute, second)::
    sumYearsSeconds(year) + sumMonthsSeconds(year, month) + sumDaysSeconds(day) + hour * 3600 + minute * 60 + second,

  // isNumeric checks that the input is a non-empty string containing only digit characters.
  local isNumeric(input) =
    assert std.type(input) == 'string' : 'isNumeric() only operates on string inputs, got %s' % std.type(input);
    std.foldl(
      function(acc, char) acc && std.codepoint('0') <= std.codepoint(char) && std.codepoint(char) <= std.codepoint('9'),
      std.stringChars(input),
      std.length(input) > 0,
    ),

  // parseSeparatedNumbers parses input which has part `names` separated by `sep`.
  // Returns an object which has one field for each name in `names` with its integer value.
  local parseSeparatedNumbers(input, sep, names) = (
    assert std.type(input) == 'string' : 'parseSeparatedNumbers() only operates on string inputs, got %s' % std.type(input);
    assert std.type(sep) == 'string' : 'parseSeparatedNumbers() only operates on string separators, got %s' % std.type(sep);
    assert std.type(names) == 'array' : 'parseSeparatedNumbers() only operates on arrays of names, got input %s' % std.type(names);

    local parts = std.split(input, sep);
    assert std.length(parts) == std.length(names) : 'expected %(expected)d parts separated by %(sep)s in %(format)s formatted input "%(input)s", but got %(got)d' % {
      expected: std.length(names),
      sep: sep,
      format: std.join(sep, names),
      input: input,
      got: std.length(parts),
    };

    {
      [names[i]]:
        // Fail with meaningful message if not numeric, otherwise it will be a hell to debug.
        assert isNumeric(parts[i]) : '%(name)%s part "%(part)s" of %(format)s of input "%(input)s" is not numeric' % {
          name: names[i],
          part: parts[i],
          format: std.join(sep, names),
          input: input,
        };
        std.parseInt(parts[i])
      for i in std.range(0, std.length(parts) - 1)
    }
  ),

  // stringContains is a helper function to check whether a string contains a given substring.
  local stringContains(haystack, needle) = std.length(std.findSubstr(needle, haystack)) > 0,

  '#parseRFC3339': d.fn(
    |||
      `parseRFC3339` parses an RFC3339-formatted date & time string (like `2020-01-02T03:04:05Z`) into an object containing the 'year', 'month', 'day', 'hour', 'minute' and 'second fields.
      This is a limited implementation that does not support timezones (so it requires an UTC input ending in 'Z' or 'z') nor sub-second precision.
      The returned object has a `toUnixTimestamp()` method that can be used to obtain the unix timestamp of the parsed date.
    |||,
    [
      d.arg('input', d.T.string),
    ],
  ),
  parseRFC3339(input)::
    // Basic input type check.
    assert std.type(input) == 'string' : 'parseRFC3339() only operates on string inputs, got %s' % std.type(input);

    // Sub-second precision isn't implemented yet, warn the user about that instead of returning wrong results.
    assert !stringContains(input, '.') : 'the provided RFC3339 input "%s" has a dot, most likely representing a sub-second precision, but this function does not support that' % input;

    // We don't support timezones, so string should end with 'Z' or 'z'.
    assert std.endsWith(input, 'Z') || std.endsWith(input, 'z') : 'the provided RFC3339 "%s" should end with "Z" or "z". This implementation does not currently support timezones' % input;

    // RFC3339 can separate date and time using 'T', 't' or ' '.
    // Find out which one it is and use it.
    local sep =
      if stringContains(input, 'T') then 'T'
      else if stringContains(input, 't') then 't'
      else if stringContains(input, ' ') then ' '
      else error 'the provided RFC3339 input "%s" should contain either a "T", or a "t" or space " " as a separator for date and time parts' % input;

    // Split date and time using the selected separator.
    // Remove the last character as we know it's 'Z' or 'z' and it's not useful to us.
    local datetime = std.split(std.substr(input, 0, std.length(input) - 1), sep);
    assert std.length(datetime) == 2 : 'the provided RFC3339 timestamp "%(input)s" does not have date and time parts separated by the character "%(sep)s"' % { input: input, sep: sep };

    local date = parseSeparatedNumbers(datetime[0], '-', ['year', 'month', 'day']);
    local time = parseSeparatedNumbers(datetime[1], ':', ['hour', 'minute', 'second']);
    date + time + {
      toUnixTimestamp():: $.toUnixTimestamp(self.year, self.month, self.day, self.hour, self.minute, self.second),
    },
}
