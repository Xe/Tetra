UNLOAD = command("UNLOAD") .. elevated() .. function(source, destination, message)
  message = luar.slice2table(message)
  if message == nil then
    return "Need a script name"
  end

  if #message == 0 then
    return "Need a script name"
  end

  local name = message[1]

  if tetra.bot.Scripts[name] == nil then
    return "Script " .. name .. " is not loaded."
  end

  if name == script.Name then
    return "Cannot unload this script!"
  end

  sleep(0.5)

  local err = tetra.bot.UnloadScript(name)

  return "Script " .. name .. " unloaded"
end

client.Commands.UNLOAD.NeedOper = true
