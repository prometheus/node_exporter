local g = import '../g.libsonnet';
local annotation = g.dashboard.annotation;
local base = import './base.libsonnet';

base {
  new(
    title,
    target,
    instanceLabels,
  ):
    super.new(title, target)
    + annotation.withIconColor('light-yellow')
    + annotation.withHide(true)
    + { useValueForTime: 'on' }
    + base.withTagKeys(instanceLabels),
}
