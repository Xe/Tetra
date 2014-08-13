Command "DIE", true, (source) ->
  client.ServicesLog "#{source.Nick}: DIE"
  tetra.bot.Quit!
  return "Okay"
