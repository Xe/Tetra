export handler = (line) ->
  source, destination, message = parseLine line

tetra.protohook "PRIVMSG", "handler"
