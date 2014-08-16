Protohook "PRIVMSG", (line) ->
  source, destination, message = parseLine line

  mymessage = strings.split message, " "

  if destination.IsChannel!
    destination.Name = destination.Name\upper!

    for kind, client in pairs tetra.bot.Services
      if client.Channels[destination.Target!] ~= nil
        if destination.Name == tetra.bot.Config.General.SnoopChan\upper!
          return
        tetra.bot.RunHook "#{kind\upper!}-CHANMSG", source, destination, mymessage
