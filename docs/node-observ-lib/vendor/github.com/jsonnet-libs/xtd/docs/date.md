---
permalink: /date/
---

# date

```jsonnet
local date = import "github.com/jsonnet-libs/xtd/date.libsonnet"
```

`time` provides various date related functions.

## Index

* [`fn dayOfWeek(year, month, day)`](#fn-dayofweek)
* [`fn dayOfYear(year, month, day)`](#fn-dayofyear)
* [`fn isLeapYear(year)`](#fn-isleapyear)
* [`fn parseRFC3339(input)`](#fn-parserfc3339)
* [`fn toUnixTimestamp(year, month, day, hour, minute, second)`](#fn-tounixtimestamp)

## Fields

### fn dayOfWeek

```ts
dayOfWeek(year, month, day)
```

`dayOfWeek` returns the day of the week for the given date. 0=Sunday, 1=Monday, etc.

### fn dayOfYear

```ts
dayOfYear(year, month, day)
```

`dayOfYear` calculates the ordinal day of the year based on the given date. The range of outputs is 1-365
for common years, and 1-366 for leap years.


### fn isLeapYear

```ts
isLeapYear(year)
```

`isLeapYear` returns true if the given year is a leap year.

### fn parseRFC3339

```ts
parseRFC3339(input)
```

`parseRFC3339` parses an RFC3339-formatted date & time string (like `2020-01-02T03:04:05Z`) into an object containing the 'year', 'month', 'day', 'hour', 'minute' and 'second fields.
This is a limited implementation that does not support timezones (so it requires an UTC input ending in 'Z' or 'z') nor sub-second precision.
The returned object has a `toUnixTimestamp()` method that can be used to obtain the unix timestamp of the parsed date.


### fn toUnixTimestamp

```ts
toUnixTimestamp(year, month, day, hour, minute, second)
```

`toUnixTimestamp` calculates the unix timestamp of a given date.
