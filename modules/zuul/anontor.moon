Hash = (foo, salt) ->
  crypto.fnv(foo)

Genvhost = (c) ->
  if c.Account == "*"
    return "tor/anonymous/#{Hash c.Uid}/#{Hash c.User}/#{Hash c.Gecos}", false

  "tor/registered#{Hash c.Uid}/#{Hash c.Account}.#{Hash c.Gecos}", true

Info = (c) ->
  client.Notice c, "Here is what I know about you: "
  client.Notice c, "   Nick:       #{c.Nick}"
  client.Notice c, "   User:       #{c.User}"
  client.Notice c, "   Gecos:      #{c.Gecos}"
  client.Notice c, "   IP Address: #{c.Host}"
  "   Account     #{c.Account}"

Command "INFO", Info

DoCloak = (c) ->
  newhost, registered = Genvhost c

  client.Chghost c, "" .. newhost
  client.ServicesLog "#{c.Nick}: ANONYMOUS USER: VHOST: #{newhost}"
  client.Notice c, "Your host has been scrambled to #{newhost}."
  client.Notice c, "Please use this anonymous hidden service with care."

  Info c

  if registered
    client.Notice c, " "
    client.Notice c, "If you have been assigned a VHost you wish to use, please run /msg HostServ ON"

Hook "NEWCLIENT", (connclient) ->
  if strings.hassuffix connclient.Server.Name, ".onion"
    DoCloak connclient

  if connclient.VHost == "tor.sasl.user"
    DoCloak connclient

