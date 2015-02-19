Hook "HOSTSERV-REQUEST", (nick, vhost) ->
  client.OperLog "HostServ: #{nick} requested #{vhost}"
