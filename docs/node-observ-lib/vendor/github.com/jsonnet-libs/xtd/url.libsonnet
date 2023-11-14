local d = import 'doc-util/main.libsonnet';

{
  '#': d.pkg(
    name='url',
    url='github.com/jsonnet-libs/xtd/url.libsonnet',
    help='`url` provides functions to deal with URLs',
  ),

  '#escapeString': d.fn(
    '`escapeString` escapes the given string so it can be safely placed inside an URL, replacing special characters with `%XX` sequences',
    args=[
      d.arg('str', d.T.string),
      d.arg('excludedChars', d.T.array, default=[]),
    ],
  ),
  escapeString(str, excludedChars=[])::
    local allowedChars = '0123456789abcdefghijklmnopqrstuvwqxyzABCDEFGHIJKLMNOPQRSTUVWQXYZ';
    local utf8(char) = std.foldl(function(a, b) a + '%%%02X' % b, std.encodeUTF8(char), '');
    local escapeChar(char) = if std.member(excludedChars, char) || std.member(allowedChars, char) then char else utf8(char);
    std.join('', std.map(escapeChar, std.stringChars(str))),

  '#encodeQuery': d.fn(
    '`encodeQuery` takes an object of query parameters and returns them as an escaped `key=value` string',
    args=[d.arg('params', d.T.object)],
  ),
  encodeQuery(params)::
    local fmtParam(p) = '%s=%s' % [self.escapeString(p), self.escapeString(params[p])];
    std.join('&', std.map(fmtParam, std.objectFields(params))),

  '#parse': d.fn(
    |||
      `parse` parses absolute and relative URLs.

      <scheme>://<netloc>/<path>;parameters?<query>#<fragment>

      Inspired by Python's urllib.urlparse, following several RFC specifications.
    |||,
    args=[d.arg('url', d.T.string)],
  ),
  parse(url):
    local hasFragment = std.member(url, '#');
    local fragmentSplit = std.splitLimit(url, '#', 1);
    local fragment = fragmentSplit[1];

    local hasQuery = std.member(fragmentSplit[0], '?');
    local querySplit = std.splitLimit(fragmentSplit[0], '?', 1);
    local query = querySplit[1];

    local hasParams = std.member(querySplit[0], ';');
    local paramsSplit = std.splitLimit(querySplit[0], ';', 1);
    local params = paramsSplit[1];

    local hasNetLoc = std.member(paramsSplit[0], '//');
    local netLocSplit = std.splitLimit(paramsSplit[0], '//', 1);
    local netLoc = std.splitLimit(netLocSplit[1], '/', 1)[0];

    local hasScheme = std.member(netLocSplit[0], ':');
    local schemeSplit = std.splitLimit(netLocSplit[0], ':', 1);
    local scheme = schemeSplit[0];

    local path =
      if hasNetLoc && std.member(netLocSplit[1], '/')
      then '/' + std.splitLimit(netLocSplit[1], '/', 1)[1]
      else if hasScheme
      then schemeSplit[1]
      else netLocSplit[0];
    local hasPath = (path != '');

    local hasUserInfo = hasNetLoc && std.member(netLoc, '@');
    local userInfoSplit = std.reverse(std.splitLimitR(netLoc, '@', 1));
    local userInfo = userInfoSplit[1];

    local hasPassword = hasUserInfo && std.member(userInfo, ':');
    local passwordSplit = std.splitLimitR(userInfo, ':', 1);
    local username = passwordSplit[0];
    local password = passwordSplit[1];

    local hasPort = hasNetLoc && std.length(std.findSubstr(':', userInfoSplit[0])) > 0;
    local portSplit = std.splitLimitR(userInfoSplit[0], ':', 1);
    local host = portSplit[0];
    local port = portSplit[1];

    {
      [if hasScheme then 'scheme']: scheme,
      [if hasNetLoc then 'netloc']: netLoc,
      [if hasPath then 'path']: path,
      [if hasParams then 'params']: params,
      [if hasQuery then 'query']: query,
      [if hasFragment then 'fragment']: fragment,

      [if hasUserInfo then 'username']: username,
      [if hasPassword then 'password']: password,
      [if hasNetLoc then 'hostname']: host,
      [if hasPort then 'port']: port,
    },

  '#join': d.fn(
    '`join` joins URLs from the object generated from `parse`',
    args=[d.arg('splitObj', d.T.object)],
  ),
  join(splitObj):
    std.join('', [
      if 'scheme' in splitObj then splitObj.scheme + ':' else '',
      if 'netloc' in splitObj then '//' + splitObj.netloc else '',
      if 'path' in splitObj then splitObj.path else '',
      if 'params' in splitObj then ';' + splitObj.params else '',
      if 'query' in splitObj then '?' + splitObj.query else '',
      if 'fragment' in splitObj then '#' + splitObj.fragment else '',
    ]),
}
