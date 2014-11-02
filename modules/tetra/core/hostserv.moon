Hook "HOSTSERV-SERVICELOG", (message) ->
  -- (@HostServ) Xena REQUEST: ninjas
  tetra.bot.RunHook "HOSTSERV-REQUEST", message[1], strings.shuck(message[3])
