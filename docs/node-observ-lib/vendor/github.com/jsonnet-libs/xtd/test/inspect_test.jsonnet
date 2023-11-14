local xtd = import '../main.libsonnet';
local test = import 'github.com/jsonnet-libs/testonnet/main.libsonnet';

test.new(std.thisFile)

+ test.case.new(
  name='emptyobject',
  test=test.expect.eq(
    actual=xtd.inspect.inspect({}),
    expected={}
  )
)

+ test.case.new(
  name='flatObject',
  test=test.expect.eq(
    actual=xtd.inspect.inspect({
      key: 'value',
      hidden_key:: 'value',
      func(value): value,
      hidden_func(value):: value,
    }),
    expected={
      fields: ['key'],
      hidden_fields: ['hidden_key'],
      functions: ['func'],
      hidden_functions: ['hidden_func'],
    }
  )
)

+ test.case.new(
  name='nestedObject',
  test=test.expect.eq(
    actual=xtd.inspect.inspect({
      nested: {
        key: 'value',
        hidden_key:: 'value',
        func(value): value,
        hidden_func(value):: value,
      },
      key: 'value',
      hidden_func(value):: value,
    }),
    expected={
      nested: {
        fields: ['key'],
        hidden_fields: ['hidden_key'],
        functions: ['func'],
        hidden_functions: ['hidden_func'],
      },
      fields: ['key'],
      hidden_functions: ['hidden_func'],
    }
  )
)

+ test.case.new(
  name='maxRecursionDepth',
  test=test.expect.eq(
    actual=xtd.inspect.inspect({
      key: 'value',
      nested: {
        key: 'value',
        nested: {
          key: 'value',
        },
      },
    }, maxDepth=1),
    expected={
      fields: ['key'],
      nested: {
        fields: ['key', 'nested'],
      },
    }
  )
)

+ test.case.new(
  name='noDiff',
  test=test.expect.eq(
    actual=xtd.inspect.diff('', ''),
    expected=''
  )
)
+ test.case.new(
  name='typeDiff',
  test=test.expect.eq(
    actual=xtd.inspect.diff('string', true),
    expected='~[ string , true ]'
  )
)
+ (
  local input1 = {
    same: 'same',
    change: 'this',
    remove: 'removed',
  };
  local input2 = {
    same: 'same',
    change: 'changed',
    add: 'added',
  };
  test.case.new(
    name='objectDiff',
    test=test.expect.eq(
      actual=xtd.inspect.diff(input1, input2),
      expected={
        'add +': 'added',
        'change ~': '~[ this , changed ]',
        'remove -': 'removed',
      }
    )
  )
)

+ (
  local input1 = [
    'same',
    'this',
    [
      'same',
      'this',
    ],
    'remove',
  ];
  local input2 = [
    'same',
    'changed',
    [
      'same',
      'changed',
      'added',
    ],
  ];
  test.case.new(
    name='arrayDiff',
    test=test.expect.eq(
      actual=xtd.inspect.diff(input1, input2),
      expected=[
        'same',
        '~[ this , changed ]',
        [
          'same',
          '~[ this , changed ]',
          '+ added',
        ],
        '- remove',
      ]
    )
  )
)
