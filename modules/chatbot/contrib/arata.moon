Command "ARATA", (source, destination, args) ->
  if #args == 0
    return "need module name to lookup"

  if args[1]\lower! ~= "plugin"
    return "can only look up plugins"

  kind = args[1]\lower!

  if kind ~= "plugin" or kind ~= "src"
    return "unsupported operation"

  if kind == "plugin"
    kind = "plugins"

  path = args[2]\gsub "%.", "%/"
  print path

  url = "https://raw.githubusercontent.com/shockkolate/arata/master/#{kind}/#{path}.hs"
  print url

  res, err = geturl url

  if err ~= nil
    return "no such #{kind} #{args[2]}"

  if res == "Not Found"
    return "no such #{kind} #{args[2]}"

  return "> #{url}"
