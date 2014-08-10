SCRIPTS = elevated() .. function(source, message)
  for name, script in pairs(tetra.bot.Scripts) do
    client.Notice(source, script.Client.Nick .. ": " .. name .. " (" .. script.Kind .. ")" .. " (" .. script.Uuid:sub(1,8) .. ")" .. " handlers: " .. #script.Handlers .. " commands: " .. #script.Commands .. " hooks: " .. #script.Hooks)
  end

  return "End of scripts list"
end

tetra.script.AddLuaCommand("SCRIPTS", "SCRIPTS")
client.Commands.SCRIPTS.NeedOper = true
