ld = nil

-- Sandbox
os = {}
io = {}

function eval_channel(line)
  local source, destination, message = parseLine(line)

  if not is_common_channel(destination) then
    return
  end

  if not source.IsOper() then
    return
  end

  ld = destination

  if message:sub(1,1) == "=" then
    tetra.log.Printf(message)
    toeval = message:sub(3)

    local func, err = loadstring(toeval)

    if err ~= nil then
      client.Privmsg(destination, string.format("error: %s",(err)))
      return
    end

    local res, err = func()

    tetra.log.Printf("%#v: %#v", res, err)

    client.Privmsg(ld, "> " .. res)
  end
end

tetra.script.AddLuaProtohook("PRIVMSG", "eval_channel")
