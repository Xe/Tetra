function linetest(line)
  tetra.log.Printf("%#v", line)
  source = tetra.bot.Clients.ByUID[line.Source]
  tetra.log.Printf("%#v", source)
end

tetra.log.Printf("%#v", tetra.script)

tetra.script.AddLuaProtohook("PRIVMSG", "linetest")
