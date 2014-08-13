Command "SENDFILE", true, (source, dest, message) ->
  if #message < 2
    return "Need a destination and file to send."

  target = message[1]
  filename = strings.split(message[2], "/")[1]

  if tetra.bot.Clients.ByNick[target\upper!] ~= nil
    target = tetra.bot.Clients.ByNick[target\upper!]
  elseif tetra.bot.Channels[target\upper!] ~= nil
    target = tetra.bot.Channels[target\upper!]
  else
    return "No such channel or user #{target}"

  with fin, err = io.open "etc/sendfile/" .. filename
    if err ~= nil
      return err

    for line in fin\lines!
      if target.IsChannel!
        client.Privmsg target, line
      else
        client.Notice target, line

  client.ServicesLog "#{source.Nick}: SENDFILE:#{message[1]}: #{filename}"
  return "Complete"
