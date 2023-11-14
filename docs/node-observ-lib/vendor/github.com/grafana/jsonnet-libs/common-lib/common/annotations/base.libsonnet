local g = import '../g.libsonnet';
local annotation = g.dashboard.annotation;
{
  new(
    title,
    target,
  ):
    annotation.withEnable(true)
    + annotation.withName(title)
    + annotation.withDatasourceMixin(target.datasource)
      {
      titleFormat: title,
      expr: target.expr,

    }
    + (if std.objectHas(target, 'interval') then { step: target.interval } else {}),

  withTagKeys(value):
    {
      tagKeys: value,
    },
  withValueForTime(value=false):
    {
      useValueForTime: value,
    },
  withTextFormat(value=''):
    {
      textFormat: value,
    },


}
