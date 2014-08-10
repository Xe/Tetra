export yo = (source, dest) ->
  client.ServicesLog "#{dest} got a yo from #{source}!"

tetra.script.AddLuaHook "YO", "yo"
