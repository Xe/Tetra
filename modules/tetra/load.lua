LOAD = command("LOAD") .. elevated() .. function(source, destination, msg)
  print(msg)

  local message = luar.slice2table(msg)
  if message == nil then
    return "Need a script name"
  end

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
end

client.Commands.LOAD.NeedsOper = true
