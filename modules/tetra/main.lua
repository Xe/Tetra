function linetest(line)
  local source, destination, message = parseLine(line)

  client.Notice(source, "Test")
end

tetra.script.AddLuaProtohook("PRIVMSG", "linetest")
