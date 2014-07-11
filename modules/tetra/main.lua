function linetest(line)
  local source, destination, message = parseLine(line)

  if is_targeted_pm(destination) then
    client.Notice(source, message)
  end
end

tetra.script.AddLuaProtohook("PRIVMSG", "linetest")
