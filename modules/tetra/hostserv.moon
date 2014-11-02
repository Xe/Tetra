Hook "HOSTSERV-SERVICELOG", (message) ->
  -- (@HostServ) Xena REQUEST: ninjas
  client.OperLog "HostServ: #{strings.join message, " "}"

  tetra.bot.RunHook "HOSTSERV-REQUEST", message[1], strings.shuck(message[3])
