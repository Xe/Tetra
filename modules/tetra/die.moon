export die = command("DIE") .. elevated! .. (source) ->
  client.ServicesLog "#{source.Nick}: DIE"
  tetra.bot.Quit!
  return "Okay"

client.Commands.DIE.NeedsOper = true
