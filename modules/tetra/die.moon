export die = elevated! .. (source) ->
  client.ServicesLog "#{source.Nick}: DIE"
  tetra.bot.Quit!
  return "Okay"

tetra.script.AddLuaCommand "DIE", "die"

client.Commands.DIE.NeedsOper = true
