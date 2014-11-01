Hook "ENCAP-SNOTE-S", (message) ->
  if not message\match "host mismatch"
    return

  split = strings.split message, " "
  user = split[8]

  if user == nil
    return

  target = tetra.bot.Clients.ByNick[user\upper!]

  for _, line in pairs {
    "You may not attempt to become an operator without being on staff.",
    "Staff has been notified. Depending on staff decisions, you might",
    "have additional consequences for this action.",
    " ",
    "Your connection will be closed. Have a good day.",
    }
    client.Notice target, line

  client.ServicesLog "#{target.Nick} falied OPER attempt (#{target.User}@#{target.Host} : #{target.VHost}), warned and killed"
  client.Kill target, "Failed OPER attempt"
