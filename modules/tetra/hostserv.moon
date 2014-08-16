export staffchan = tetra.bot.Config.General.StaffChan

lookupVhost = (vhost) ->
  if #vhost < 3
    return

  return tetra.bot.Atheme.HostServ.ListPattern strings.format("*%s*", vhost)

Hook "HOSTSERV-SERVICELOG", (message) ->
  -- (@HostServ) Xena REQUEST: ninjas
  client.OperLog "HostServ: #{strings.join message, " "}"
