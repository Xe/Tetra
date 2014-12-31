Info = (c) ->
  client.Notice c, "Here is what I know about you: "
  client.Notice c, "   Nick:       #{c.Nick}"
  client.Notice c, "   User:       #{c.User}"
  client.Notice c, "   Gecos:      #{c.Gecos}"
  client.Notice c, "   IP Address: #{c.Ip}"
  "   Account     #{c.Account}"

Command "INFO", Info

