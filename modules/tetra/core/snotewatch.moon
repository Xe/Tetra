Hook "ENCAP-SNOTE", (source, args) ->
  client.ServicesLog "SNOTE: " .. args[1] .. " " .. args[2]
