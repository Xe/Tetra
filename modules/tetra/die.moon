command("DIE", true) .. elevated! .. (source) ->
  client.ServicesLog "#{source.Nick}: DIE"
  tetra.bot.Quit!
  return "Okay"
