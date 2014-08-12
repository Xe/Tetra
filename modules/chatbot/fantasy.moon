require "modules/base"

hook("CHATBOT-CHANMSG") .. (source, dest, msg) ->
  command = strings.rest(msg[1])\upper!
  prefix = strings.first(msg[1])

  if prefix == tetra.bot.Config.General.Prefix
    if client.Commands[command] ~= nil
      args = [i for i in *msg[2,]]

      res, err = script.Call command, source, dest, args

      if command == "HELP"
        client.Notice source, res
      else
        client.Privmsg dest, res

-- TODO: implement ping-prefixing
