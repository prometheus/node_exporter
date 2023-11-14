local xtd = import '../main.libsonnet';
local test = import 'github.com/jsonnet-libs/testonnet/main.libsonnet';

test.new(std.thisFile)

+ test.case.new(
  name='split: nostring',
  test=test.expect.eq(
    actual=xtd.camelcase.split(''),
    expected=[''],
  )
)
+ test.case.new(
  name='split: lowercase',
  test=test.expect.eq(
    actual=xtd.camelcase.split('lowercase'),
    expected=['lowercase'],
  )
)
+ test.case.new(
  name='split: Class',
  test=test.expect.eq(
    actual=xtd.camelcase.split('Class'),
    expected=['Class'],
  )
)
+ test.case.new(
  name='split: MyClass',
  test=test.expect.eq(
    actual=xtd.camelcase.split('MyClass'),
    expected=['My', 'Class'],
  )
)
+ test.case.new(
  name='split: MyC',
  test=test.expect.eq(
    actual=xtd.camelcase.split('MyC'),
    expected=['My', 'C'],
  )
)
+ test.case.new(
  name='split: HTML',
  test=test.expect.eq(
    actual=xtd.camelcase.split('HTML'),
    expected=['HTML'],
  )
)
+ test.case.new(
  name='split: PDFLoader',
  test=test.expect.eq(
    actual=xtd.camelcase.split('PDFLoader'),
    expected=['PDF', 'Loader'],
  )
)
+ test.case.new(
  name='split: AString',
  test=test.expect.eq(
    actual=xtd.camelcase.split('AString'),
    expected=['A', 'String'],
  )
)
+ test.case.new(
  name='split: SimpleXMLParser',
  test=test.expect.eq(
    actual=xtd.camelcase.split('SimpleXMLParser'),
    expected=['Simple', 'XML', 'Parser'],
  )
)
+ test.case.new(
  name='split: vimRPCPlugin',
  test=test.expect.eq(
    actual=xtd.camelcase.split('vimRPCPlugin'),
    expected=['vim', 'RPC', 'Plugin'],
  )
)
+ test.case.new(
  name='split: GL11Version',
  test=test.expect.eq(
    actual=xtd.camelcase.split('GL11Version'),
    expected=['GL', '11', 'Version'],
  )
)
+ test.case.new(
  name='split: 99Bottles',
  test=test.expect.eq(
    actual=xtd.camelcase.split('99Bottles'),
    expected=['99', 'Bottles'],
  )
)
+ test.case.new(
  name='split: May5',
  test=test.expect.eq(
    actual=xtd.camelcase.split('May5'),
    expected=['May', '5'],
  )
)
+ test.case.new(
  name='split: BFG9000',
  test=test.expect.eq(
    actual=xtd.camelcase.split('BFG9000'),
    expected=['BFG', '9000'],
  )
)
+ test.case.new(
  name='split: Two  spaces',
  test=test.expect.eq(
    actual=xtd.camelcase.split('Two  spaces'),
    expected=['Two', '  ', 'spaces'],
  )
)
+ test.case.new(
  name='split: Multiple   Random  spaces',
  test=test.expect.eq(
    actual=xtd.camelcase.split('Multiple   Random  spaces'),
    expected=['Multiple', '   ', 'Random', '  ', 'spaces'],
  )
)
+ test.case.new(
  name='toCamelCase: empty string',
  test=test.expect.eq(
    actual=xtd.camelcase.toCamelCase(''),
    expected='',
  )
)
+ test.case.new(
  name='toCamelCase: lowercase',
  test=test.expect.eq(
    actual=xtd.camelcase.toCamelCase('lowercase'),
    expected='lowercase',
  )
)
+ test.case.new(
  name='toCamelCase: underscores',
  test=test.expect.eq(
    actual=xtd.camelcase.toCamelCase('lower_case'),
    expected='lowerCase',
  )
)
+ test.case.new(
  name='toCamelCase: dashes',
  test=test.expect.eq(
    actual=xtd.camelcase.toCamelCase('lower-case'),
    expected='lowerCase',
  )
)
+ test.case.new(
  name='toCamelCase: spaces',
  test=test.expect.eq(
    actual=xtd.camelcase.toCamelCase('lower case'),
    expected='lowerCase',
  )
)
+ test.case.new(
  name='toCamelCase: mixed',
  test=test.expect.eq(
    actual=xtd.camelcase.toCamelCase('lower_case-mixed'),
    expected='lowerCaseMixed',
  )
)
+ test.case.new(
  name='toCamelCase: multiple spaces',
  test=test.expect.eq(
    actual=xtd.camelcase.toCamelCase('lower  case'),
    expected='lowerCase',
  )
)
+ test.case.new(
  name='toCamelCase: PascalCase',
  test=test.expect.eq(
    actual=xtd.camelcase.toCamelCase('PascalCase'),
    expected='pascalCase',
  )
)

// TODO: find or create is(Upper|Lower) for non-ascii characters
// Something like this for Jsonnet:
// https://cs.opensource.google/go/go/+/refs/tags/go1.17.3:src/unicode/tables.go
//+ test.case.new(
//  name='BöseÜberraschung',
//  test=test.expect.eq(
//    actual=xtd.camelcase.split('BöseÜberraschung'),
//    expected=['Böse', 'Überraschung'],
//  )
//)

// This doesn't even render in Jsonnet
//+ test.case.new(
//  name="BadUTF8\xe2\xe2\xa1",
//  test=test.expect.eq(
//    actual=xtd.camelcase.split("BadUTF8\xe2\xe2\xa1"),
//    expected=["BadUTF8\xe2\xe2\xa1"],
//  )
//)
