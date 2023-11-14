local g = import '../g.libsonnet';
local annotation = g.dashboard.annotation;
local base = import './base.libsonnet';

base {
  new(
    title,
    target,
  ):
    super.new(title, target)
    + annotation.withIconColor('light-orange')
    + annotation.withHide(true),
}
