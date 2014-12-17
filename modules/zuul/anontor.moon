Hash = (foo, salt) ->
  crypto.hash(foo, salt)\sub(1,8)\upper!

Genvhost = (c) ->
  if c.Account == "*"
    return "tor.anonymous.#{Hash c.Uid, c.Server.Name}.#{Hash c.User, c.Server.Name}.#{Hash c.Gecos, c.Server.Name}", false

  return "tor.registered.#{Hash c.Uid, c.Server.Name}.#{Hash c.Account, c.Server.Name}.#{Hash c.Gecos, c.Server.Name}", true

DoCloak = (c) ->
  newhost, registered = Genvhost c

  client.Chghost c, "" .. newhost
  client.ServicesLog "#{c.Nick}: ANONYMOUS USER: VHOST: #{newhost}"
  client.Notice c, "Your host has been scrambled to #{newhost} to allow for accountability."
  client.Notice c, "Please use this anonymous hidden service with care."

  if registered
    client.Notice c, "If you have a VHost you wish to use, please run /msg HostServ ON"

Hook "NEWCLIENT", (connclient) ->
  if strings.hassuffix connclient.Server.Name, ".onion"
    DoCloak connclient

  if connclient.VHost == "tor.sasl.user"
    DoCloak connclient

