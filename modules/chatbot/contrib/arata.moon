Command "ARATA", (source, destination, args) ->
  if #args == 0
    return "need module name to lookup"

  if args[1]\lower! ~= "plugin"
    return "can only look up plugins"

  if not (args[2]\lower!)\match "nickserv"
    return "can only look up plugins for nickserv"

  url = "https://github.com/shockkolate/arata/blob/master/plugins/#{args[2]\gsub "%.", "%/", 1}.hs"

  _, err = geturl url

  if err ~= nil
    return "no such plugin #{args[2]}"

  return "> #{url}"
