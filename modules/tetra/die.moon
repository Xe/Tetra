Command "DIE", true, (source) ->
  client.ServicesLog "#{source.Nick}: DIE"
  tetra.Quit!
  return "Okay"
