Command "INFO", true, (source, target, args) ->
  if #args == 0
    return "INFO <channel name>"

  channel = tetra.Channels[args[1]\upper!]
  if not channel
    return "Channel #{args[1]} does not exist"

  client.Notice source, "Info about #{channel.Name}:"
  client.Notice source, "#{#channel.Clients} clients:"

  for _, chanuser in pairs channel.Clients
    client.Notice source, "#{chanuser.Client.Nick}: (#{chanuser.Client.User}@#{chanuser.Client.VHost}: #{chanuser.Client.Ip}) 0x#{"%x"\format chanuser.Prefix}"

  client.Notice source, "Channel mode: 0x#{"%x"\format channel.Modes}"

  "End info on #{args[1]}"
