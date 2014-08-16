Protohook "PRIVMSG", (line) ->
  source, destination, message = parseLine line

  mymessage = strings.split message, " "

  if destination.IsChannel!
    destination.Name = destination.Name\upper!
  else
    return

  if destination.Name == tetra.bot.Config.General.SnoopChan\upper!
    tetra.bot.RunHook "#{source.Nick\upper!}-SERVICELOG", mymessage
    return

  for kind, client in pairs tetra.bot.Services
    if client.Channels[destination.Target!] ~= nil
      tetra.bot.RunHook "#{kind\upper!}-CHANMSG", source, destination, mymessage
