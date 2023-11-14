local jsonpath = import '../jsonpath.libsonnet';
local test = import 'github.com/jsonnet-libs/testonnet/main.libsonnet';

test.new(std.thisFile)

// Root
+ test.case.new(
  name='root $',
  test=test.expect.eq(
    actual=jsonpath.getJSONPath({ key: 'content' }, '$'),
    expected={ key: 'content' },
  )
)
+ test.case.new(
  name='root (empty path)',
  test=test.expect.eq(
    actual=jsonpath.getJSONPath({ key: 'content' }, ''),
    expected={ key: 'content' },
  )
)
+ test.case.new(
  name='root .',
  test=test.expect.eq(
    actual=jsonpath.getJSONPath({ key: 'content' }, '.'),
    expected={ key: 'content' },
  )
)

// Single key
+ test.case.new(
  name='path without dot prefix',
  test=test.expect.eq(
    actual=jsonpath.getJSONPath({ key: 'content' }, 'key'),
    expected='content',
  )
)
+ test.case.new(
  name='single key',
  test=test.expect.eq(
    actual=jsonpath.getJSONPath({ key: 'content' }, '.key'),
    expected='content',
  )
)
+ test.case.new(
  name='single bracket key',
  test=test.expect.eq(
    actual=jsonpath.getJSONPath({ key: 'content' }, '[key]'),
    expected='content',
  )
)
+ test.case.new(
  name='single bracket key with $',
  test=test.expect.eq(
    actual=jsonpath.getJSONPath({ key: 'content' }, '$[key]'),
    expected='content',
  )
)
+ test.case.new(
  name='single array index',
  test=test.expect.eq(
    actual=jsonpath.getJSONPath(['content'], '.[0]'),
    expected='content',
  )
)
+ test.case.new(
  name='single array index without dot prefix',
  test=test.expect.eq(
    actual=jsonpath.getJSONPath(['content'], '[0]'),
    expected='content',
  )
)

// Nested
+ test.case.new(
  name='nested key',
  test=test.expect.eq(
    actual=jsonpath.getJSONPath(
      { key1: { key2: { key3: 'content' } } },
      '.key1.key2.key3'
    ),
    expected='content',
  )
)
+ test.case.new(
  name='nested bracket key',
  test=test.expect.eq(
    actual=jsonpath.getJSONPath(
      { key1: { key2: { key3: 'content' } } },
      '.key1.key2[key3]'
    ),
    expected='content',
  )
)
+ test.case.new(
  name='nested bracket key (quoted)',
  test=test.expect.eq(
    actual=jsonpath.getJSONPath(
      { key1: { key2: { key3: 'content' } } },
      ".key1.key2['key3']"
    ),
    expected='content',
  )
)
+ test.case.new(
  name='nested bracket star key',
  test=test.expect.eq(
    actual=jsonpath.getJSONPath(
      { key1: { key2: { key3: 'content' } } },
      '.key1.key2[*]'
    ),
    expected={ key3: 'content' },
  )
)
+ test.case.new(
  name='nested array index',
  test=test.expect.eq(
    actual=jsonpath.getJSONPath(
      { key1: { key2: ['content1', 'content2'] } },
      '.key1.key2[1]'
    ),
    expected='content2',
  )
)
+ test.case.new(
  name='nested array index with $',
  test=test.expect.eq(
    actual=jsonpath.getJSONPath(
      { key1: { key2: ['content1', 'content2'] } },
      '$.key1.key2[1]'
    ),
    expected='content2',
  )
)
+ test.case.new(
  name='nested array index without brackets',
  test=test.expect.eq(
    actual=jsonpath.getJSONPath(
      { key1: { key2: ['content1', 'content2'] } },
      '.key1.key2.1'
    ),
    expected='content2',
  )
)
+ test.case.new(
  name='nested array star index',
  test=test.expect.eq(
    actual=jsonpath.getJSONPath(
      { key1: { key2: ['content1', 'content2'] } },
      '.key1.key2[*]'
    ),
    expected=['content1', 'content2'],
  )
)
+ test.case.new(
  name='nested bracket keys and array index combo',
  test=test.expect.eq(
    actual=jsonpath.getJSONPath(
      { key1: { key2: ['content1', 'content2'] } },
      '$.[key1][key2][1]'
    ),
    expected='content2',
  )
)
+ test.case.new(
  name='all keys in bracket and quoted',
  test=test.expect.eq(
    actual=jsonpath.getJSONPath(
      { key1: { key2: ['content1', 'content2'] } },
      "$['key1']['key2']"
    ),
    expected=['content1', 'content2'],
  )
)

// index range/slice
+ test.case.new(
  name='array with index range (first two)',
  test=test.expect.eq(
    actual=jsonpath.getJSONPath(
      { key: ['content1', 'content2', 'content3'] },
      'key[0:2]'
    ),
    expected=['content1', 'content2'],
  )
)
+ test.case.new(
  name='array with index range (last two)',
  test=test.expect.eq(
    actual=jsonpath.getJSONPath(
      { key: ['content1', 'content2', 'content3'] },
      'key[1:3]'
    ),
    expected=['content2', 'content3'],
  )
)
+ test.case.new(
  name='array with index range (until end)',
  test=test.expect.eq(
    actual=jsonpath.getJSONPath(
      { key: ['content1', 'content2', 'content3'] },
      'key[1:]'
    ),
    expected=['content2', 'content3'],
  )
)
+ test.case.new(
  name='array with index range (from beginning)',
  test=test.expect.eq(
    actual=jsonpath.getJSONPath(
      { key: ['content1', 'content2', 'content3'] },
      'key[:2]'
    ),
    expected=['content1', 'content2'],
  )
)
+ test.case.new(
  name='array with index range (negative start)',
  test=test.expect.eq(
    actual=jsonpath.getJSONPath(
      { key: ['content1', 'content2', 'content3'] },
      'key[-2:]'
    ),
    expected=['content2', 'content3'],
  )
)
+ test.case.new(
  name='array with index range (negative end)',
  test=test.expect.eq(
    actual=jsonpath.getJSONPath(
      { key: ['content1', 'content2', 'content3'] },
      'key[:-1]'
    ),
    expected=['content1', 'content2'],
  )
)
+ test.case.new(
  name='array with index range (step)',
  test=test.expect.eq(
    actual=jsonpath.getJSONPath(
      { key: [
        'content%s' % i
        for i in std.range(1, 10)
      ] },
      'key[:5:2]'
    ),
    expected=['content1', 'content3', 'content5'],
  )
)

// filter expr
+ test.case.new(
  name='array with filter expression - string',
  test=test.expect.eq(
    actual=jsonpath.getJSONPath(
      { key: [
        {
          key: 'content%s' % i,
        }
        for i in std.range(1, 10)
      ] },
      '.key[?(@.key==content2)]'
    ),
    expected=[{
      key: 'content2',
    }],
  )
)
+ test.case.new(
  name='array with filter expression - number',
  test=test.expect.eq(
    actual=jsonpath.getJSONPath(
      { key: [
        {
          count: i,
        }
        for i in std.range(1, 10)
      ] },
      '.key[?(@.count<=2)]'
    ),
    expected=[{
      count: 1,
    }, {
      count: 2,
    }],
  )
)
+ test.case.new(
  name='array with filter expression - has key',
  test=test.expect.eq(
    actual=jsonpath.getJSONPath(
      { key: [
        {
          key1: 'value',
        },
        {
          key2: 'value',
        },
      ] },
      '.key[?(@.key1)]'
    ),
    expected=[{
      key1: 'value',
    }],
  )
)
