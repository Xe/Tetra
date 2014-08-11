require "modules/base"

protohook("PRIVMSG") .. (line) ->
  source, destination, message = parseLine line

  if destination.IsChannel!
    destination.Name = destination.Name\upper!

    for kind, client in pairs tetra.bot.Services
      if client.Channels[destination.Target!] ~= nil
        if destination.Name == tetra.bot.Config.Server.SnoopChan\upper!
          return
        tetra.bot.RunHook "#{kind\upper!}-CHANMSG", source, destination, message
