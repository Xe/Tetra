export staffchan = tetra.bot.Config.General.StaffChan

lookupVhost = (vhost) ->
  if #vhost < 3
    return

  if #vhost > 7
    _, err = tetra.bot.Atheme.HostServ.ListPattern strings.format("*%s*", vhost\sub(3,#vhost-2))
    if err ~= nil
      error err
  else
    return tetra.bot.Atheme.HostServ.ListPattern strings.format("*%s*", vhost)

Hook "HOSTSERV-SERVICELOG", (message) ->
  -- (@HostServ) Xena REQUEST: ninjas
  account = message[1]
  verb = message[2]\sub 1, #message[2]-1

  print account, " ", verb

  if strings.first message[2], "("
    -- (@HostServ) Xena (Xe) REQUEST: ninjas
    account = string.sub(message[2], 2, #message[2]-1)
    verb = message[3]\sub 1, #message[2]-1

  if verb == "REQUEST"
    print "REQUEST"
    vhost = message[4]
    split = strings.split vhost, "."

    user = tetra.bot.Atheme.NickServ.Info account

    if user["vhost"] == nil
      client.Privmsg tetra.bot.Channels[staffchan], "#{account} has no vhost"
    else
      client.Privmsg tetra.bot.Channels[staffchan], "Vhost for #{account} is #{user.vhost}"

    for k, v in pairs split
      lookupVhost v
  else
    client.Privmsg tetra.bot.Channels[staffchan], "HostServ: #{message}"
