Genvhost = (c) ->
  if c.Account == "*"
    return "tor/anonymous/#{crypto.hash(c.User)\sub(1,14)}/#{c.Uid}", false

  return "tor/registered/#{c.Account}}/#{c.Uid}", true

DoCloak = (c) ->
  newhost, registered = Genvhost connclient

  client.Chghost c, newhost
  client.Notice c, "Your host has been scrambled to #{newhost} to allow for accountability."
  client.Notice c, "Please use this anonymous hidden service with care."

  if registered
    client.Notice c, "If you have a VHost you wish to use, please run /msg HostServ ON"

Hook "NEWCLIENT", (connclient) ->
  if strings.hassuffix connclient.Server.Name, ".onion"
    DoCloak connclient

  if connclient.VHost == "tor.sasl.user"
    DoCloak connclient

