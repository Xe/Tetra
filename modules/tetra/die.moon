export die = elevated! .. ->
  tetra.bot.Quit!
  return "Okay"

tetra.script.AddLuaCommand "DIE", "die"

client.Commands.DIE.NeedsOper = true
