commands = {
  LOAD = elevated() .. function(source, message)
    if #message == 0 then
      return "Need a script name"
    end

    local name = message[1]
    local script, err = tetra.bot.LoadScript(message[1])

    if err ~= nil then
      tetra.log.Printf("Can't load script " .. name .. ": %#v", err)
      return "Script " .. name .. " failed load: " .. err
    end

    return "Script " .. script.Name .. " loaded with uuid " .. script.Uuid
  end,
  UNLOAD = elevated() .. function(source, message)
    if #message == 0 then
      return "Need a script name"
    end

    local name = message[1]

    if tetra.bot.Scripts[name] == nil then
      return "Script " .. name .. " is not loaded."
    end

    tetra.log.Printf("%#v", tetra.script)
    tetra.log.Printf("%#v", client)

    local err = tetra.bot.UnloadScript(name)

    return "Script " .. name .. " unloaded"
  end,
}

function parsecommands(source, message)
  message = split(message, " ")
  local verb = string.upper(table.remove(message, 1))
  local reply = ""

  if commands[verb] ~= nil then
    reply = commands[verb](source, message)
  else
    reply = "No such command " .. verb .. ". If you are having trouble, join #help."
  end

  client.Notice(source, reply)
end

function admincommands(line)
  local source, destination, message = parseLine(line)

  if is_targeted_pm(destination) then
    parsecommands(source, message)
  end
end

tetra.script.AddLuaProtohook("PRIVMSG", "admincommands")
