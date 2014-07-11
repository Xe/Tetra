function linetest(line)
  local source, destination, message = parseLine(line)

  if destination.Target():sub(1,1) == "#" then
    client.Privmsg(destination, message)
  else
    if destination.Kind == client.Kind then
      client.Notice(source, message)
    end
  end
end

tetra.script.AddLuaProtohook("PRIVMSG", "linetest")
