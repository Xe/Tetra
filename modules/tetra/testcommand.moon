export testcommand = (client, target, message) ->
  return "PONG"

tetra.script.AddLuaCommand "PING", "testcommand"
