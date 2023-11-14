local xtd = import '../main.libsonnet';
local test = import 'github.com/jsonnet-libs/testonnet/main.libsonnet';

test.new(std.thisFile)

+ test.case.new(
  name='Leap Year commonYear',
  test=test.expect.eq(
    actual=xtd.date.isLeapYear(1995),
    expected=false,
  )
)

+ test.case.new(
  name='Leap Year fourYearCycle',
  test=test.expect.eq(
    actual=xtd.date.isLeapYear(1996),
    expected=true,
  )
)

+ test.case.new(
  name='Leap Year fourHundredYearCycle',
  test=test.expect.eq(
    actual=xtd.date.isLeapYear(2000),
    expected=true,
  )
)

+ test.case.new(
  name='Leap Year hundredYearCycle',
  test=test.expect.eq(
    actual=xtd.date.isLeapYear(2100),
    expected=false,
  )
)

+ test.case.new(
  name='Day Of Week leapYearStart',
  test=test.expect.eq(
    actual=xtd.date.dayOfWeek(2000, 1, 1),
    expected=6,
  )
)

+ test.case.new(
  name='Day Of Week leapYearEnd',
  test=test.expect.eq(
    actual=xtd.date.dayOfWeek(2000, 12, 31),
    expected=0,
  )
)

+ test.case.new(
  name='Day Of Week commonYearStart',
  test=test.expect.eq(
    actual=xtd.date.dayOfWeek(1995, 1, 1),
    expected=0,
  )
)

+ test.case.new(
  name='Day Of Week commonYearEnd',
  test=test.expect.eq(
    actual=xtd.date.dayOfWeek(2003, 12, 31),
    expected=3,
  )
)

+ test.case.new(
  name='Day Of Week leapYearMid',
  test=test.expect.eq(
    actual=xtd.date.dayOfWeek(2024, 7, 19),
    expected=5,
  )
)

+ test.case.new(
  name='Day Of Week commonYearMid',
  test=test.expect.eq(
    actual=xtd.date.dayOfWeek(2023, 6, 15),
    expected=4,
  )
)
+ test.case.new(
  name='Day Of Year leapYearStart',
  test=test.expect.eq(
    actual=xtd.date.dayOfYear(2000, 1, 1),
    expected=1,
  )
)

+ test.case.new(
  name='Day Of Year leapYearEnd',
  test=test.expect.eq(
    actual=xtd.date.dayOfYear(2000, 12, 31),
    expected=366,
  )
)

+ test.case.new(
  name='Day Of Year commonYearStart',
  test=test.expect.eq(
    actual=xtd.date.dayOfYear(1995, 1, 1),
    expected=1,
  )
)

+ test.case.new(
  name='Day Of Year commonYearEnd',
  test=test.expect.eq(
    actual=xtd.date.dayOfYear(2003, 12, 31),
    expected=365,
  )
)

+ test.case.new(
  name='Day Of Year leapYearMid',
  test=test.expect.eq(
    actual=xtd.date.dayOfYear(2024, 7, 19),
    expected=201,
  )
)

+ test.case.new(
  name='Day Of Year commonYearMid',
  test=test.expect.eq(
    actual=xtd.date.dayOfYear(2023, 6, 15),
    expected=166,
  )
)

+ test.case.new(
  name='toUnixTimestamp of 1970-01-01 00:00:00 (zero)',
  test=test.expect.eq(
    actual=xtd.date.toUnixTimestamp(1970, 1, 1, 0, 0, 0),
    expected=0,
  ),
)

+ test.case.new(
  name='toUnixTimestamp of 1970-01-02 00:00:00 (one day)',
  test=test.expect.eq(
    actual=xtd.date.toUnixTimestamp(1970, 1, 2, 0, 0, 0),
    expected=86400,
  ),
)

+ test.case.new(
  name='toUnixTimestamp of 1971-01-01 00:00:00 (one year)',
  test=test.expect.eq(
    actual=xtd.date.toUnixTimestamp(1971, 1, 1, 0, 0, 0),
    expected=365 * 24 * 3600,
  ),
)

+ test.case.new(
  name='toUnixTimestamp of 1972-03-01 00:00:00 (month of leap year)',
  test=test.expect.eq(
    actual=xtd.date.toUnixTimestamp(1972, 3, 1, 0, 0, 0),
    expected=2 * 365 * 24 * 3600 + 31 * 24 * 3600 + 29 * 24 * 3600,
  ),
)

+ test.case.new(
  name='toUnixTimestamp of 1974-01-01 00:00:00 (incl leap year)',
  test=test.expect.eq(
    actual=xtd.date.toUnixTimestamp(1974, 1, 1, 0, 0, 0),
    expected=(4 * 365 + 1) * 24 * 3600,
  ),
)

+ test.case.new(
  name='toUnixTimestamp of 2020-01-02 03:04:05 (full date)',
  test=test.expect.eq(
    actual=xtd.date.toUnixTimestamp(2020, 1, 2, 3, 4, 5),
    expected=1577934245,
  ),
)

+ test.case.new(
  name='parseRFC3339 of 1970-01-01T00:00:00Z (standard unix zero)',
  test=test.expect.eq(
    actual=xtd.date.parseRFC3339('1970-01-01T00:00:00Z'),
    expected={ year: 1970, month: 1, day: 1, hour: 0, minute: 0, second: 0 },
  ),
)

+ test.case.new(
  name='parseRFC3339 of 2020-01-02T03:04:05Z (non-zero date)',
  test=test.expect.eq(
    actual=xtd.date.parseRFC3339('2020-01-02T03:04:05Z'),
    expected={ year: 2020, month: 1, day: 2, hour: 3, minute: 4, second: 5 },
  ),
)

+ test.case.new(
  name='parseRFC3339 of 2020-01-02 03:04:05Z (space separator)',
  test=test.expect.eq(
    actual=xtd.date.parseRFC3339('2020-01-02 03:04:05Z'),
    expected={ year: 2020, month: 1, day: 2, hour: 3, minute: 4, second: 5 },
  ),
)

+ test.case.new(
  name='parseRFC3339 of 2020-01-02t03:04:05Z (lowercase t separator and lowercase z)',
  test=test.expect.eq(
    actual=xtd.date.parseRFC3339('2020-01-02t03:04:05z'),
    expected={ year: 2020, month: 1, day: 2, hour: 3, minute: 4, second: 5 },
  ),
)

+ test.case.new(
  name='parseRFC3339(..).toUnixTimestamp()',
  test=test.expect.eq(
    actual=xtd.date.parseRFC3339('2020-01-02T03:04:05Z').toUnixTimestamp(),
    expected=1577934245,
  ),
)
