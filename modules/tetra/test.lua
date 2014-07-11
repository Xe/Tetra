function test(line)
  local source, destination, message = parseLine(line)

  if is_common_channel(destination) then
    tetra.log.Printf("Common channel!")
  else
    tetra.log.Printf("Uncommon channel!")
  end
end

tetra.script.AddLuaProtohook("PRIVMSG", "test")
