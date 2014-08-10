SCRIPTS = elevated() .. function(source, message)
  for name, script in pairs(tetra.bot.Scripts) do
    client.Notice(source, script.Client.Nick .. ": " .. name .. " (" .. script.Uuid:sub(1,8) .. ")" .. " handlers: " .. #script.Handlers)
  end

  return "End of scripts list"
end

tetra.script.AddLuaCommand("SCRIPTS", "SCRIPTS")
