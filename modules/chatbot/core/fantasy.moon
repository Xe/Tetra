use "strings"

Hook "CHATBOT-CHANMSG", (source, dest, msg) ->
  command = strings.rest(msg[1])\upper!
  prefix = strings.first(msg[1])

  if prefix == tetra.Config.General.Prefix
    if client.Commands[command] ~= nil
      args = [i for i in *msg[2,]]

      if not source.IsOper! and command.NeedsOper
        client.Notice source, "Permission denied"
        return

      res, err = script.Call command, source, dest, args

      if command == "HELP"
        client.Notice source, res
      else
        client.Privmsg dest, res

Hook "CHATBOT-CHANMSG", (source, dest, msg) ->
  if msg[1]\upper!\match client.Nick\upper!
    command = msg[2]\upper!

    if client.Commands[command]
      if not source.IsOper! and command.NeedsOper
        client.Notice source, "Permission denied"
        return

      args = [i for i in *msg[3,]]

      res, err = script.Call command, source, dest, args

      if command == "HELP"
        client.Notice source, res
      else
        client.Privmsg dest, res
