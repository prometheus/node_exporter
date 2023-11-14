local xtd = import '../main.libsonnet';
local test = import 'github.com/jsonnet-libs/testonnet/main.libsonnet';

test.new(std.thisFile)
+ test.case.new(
  name='empty',
  test=test.expect.eq(
    actual=xtd.url.escapeString(''),
    expected='',
  )
)

+ test.case.new(
  name='abc',
  test=test.expect.eq(
    actual=xtd.url.escapeString('abc'),
    expected='abc',
  )
)

+ test.case.new(
  name='space',
  test=test.expect.eq(
    actual=xtd.url.escapeString('one two'),
    expected='one%20two',
  )
)

+ test.case.new(
  name='percent',
  test=test.expect.eq(
    actual=xtd.url.escapeString('10%'),
    expected='10%25',
  )
)

+ test.case.new(
  name='complex',
  test=test.expect.eq(
    actual=xtd.url.escapeString(" ?&=#+%!<>#\"{}|\\^[]`☺\t:/@$'()*,;"),
    expected='%20%3F%26%3D%23%2B%25%21%3C%3E%23%22%7B%7D%7C%5C%5E%5B%5D%60%E2%98%BA%09%3A%2F%40%24%27%28%29%2A%2C%3B',
  )
)

+ test.case.new(
  name='exclusions',
  test=test.expect.eq(
    actual=xtd.url.escapeString('hello, world', [',']),
    expected='hello,%20world',
  )
)

+ test.case.new(
  name='multiple exclusions',
  test=test.expect.eq(
    actual=xtd.url.escapeString('hello, world,&', [',', '&']),
    expected='hello,%20world,&',
  )
)

+ test.case.new(
  name='empty',
  test=test.expect.eq(
    actual=xtd.url.encodeQuery({}),
    expected='',
  )
)

+ test.case.new(
  name='simple',
  test=test.expect.eq(
    actual=xtd.url.encodeQuery({ q: 'puppies', oe: 'utf8' }),
    expected='oe=utf8&q=puppies',
  )
)

// url.parse
+ test.case.new(
  name='Full absolute URL',
  test=test.expect.eqJson(
    actual=xtd.url.parse('https://example.com/path/to/location;type=person?name=john#address'),
    expected={
      scheme: 'https',
      netloc: 'example.com',
      hostname: 'example.com',
      path: '/path/to/location',
      params: 'type=person',
      query: 'name=john',
      fragment: 'address',
    },
  )
)

+ test.case.new(
  name='URL with fragment before params and query',
  test=test.expect.eqJson(
    actual=xtd.url.parse('https://example.com/path/to/location#address;type=person?name=john'),
    expected={
      scheme: 'https',
      netloc: 'example.com',
      hostname: 'example.com',
      path: '/path/to/location',
      fragment: 'address;type=person?name=john',
    },
  )
)

+ test.case.new(
  name='URL without query',
  test=test.expect.eqJson(
    actual=xtd.url.parse('https://example.com/path/to/location;type=person#address'),
    expected={
      scheme: 'https',
      netloc: 'example.com',
      hostname: 'example.com',
      path: '/path/to/location',
      params: 'type=person',
      fragment: 'address',
    },
  )
)

+ test.case.new(
  name='URL without params',
  test=test.expect.eqJson(
    actual=xtd.url.parse('https://example.com/path/to/location?name=john#address'),
    expected={
      scheme: 'https',
      netloc: 'example.com',
      hostname: 'example.com',
      path: '/path/to/location',
      query: 'name=john',
      fragment: 'address',
    },
  )
)

+ test.case.new(
  name='URL with empty fragment',
  test=test.expect.eqJson(
    actual=xtd.url.parse('https://example.com/path/to/location#'),
    expected={
      scheme: 'https',
      netloc: 'example.com',
      hostname: 'example.com',
      path: '/path/to/location',
      fragment: '',
    },
  )
)

+ test.case.new(
  name='host with port',
  test=test.expect.eqJson(
    actual=xtd.url.parse('//example.com:80'),
    expected={
      netloc: 'example.com:80',
      hostname: 'example.com',
      port: '80',
    },
  )
)

+ test.case.new(
  name='mailto',
  test=test.expect.eqJson(
    actual=xtd.url.parse('mailto:john@example.com'),
    expected={
      scheme: 'mailto',
      path: 'john@example.com',
    },
  )
)

+ test.case.new(
  name='UserInfo',
  test=test.expect.eqJson(
    actual=xtd.url.parse('ftp://admin:password@example.com'),

    expected={
      hostname: 'example.com',
      netloc: 'admin:password@example.com',
      scheme: 'ftp',
      username: 'admin',
      password: 'password',
    }
    ,
  )
)

+ test.case.new(
  name='Relative URL only',
  test=test.expect.eqJson(
    actual=xtd.url.parse('/path/to/location'),
    expected={
      path: '/path/to/location',
    },
  )
)

+ test.case.new(
  name='URL fragment only',
  test=test.expect.eqJson(
    actual=xtd.url.parse('#fragment_only'),
    expected={
      fragment: 'fragment_only',
    },
  )
)
