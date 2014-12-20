Command("SCRIPTS", true, function(source, message)
  for name, script in pairs(tetra.Scripts) do
    local res = ""

    local res = name
    res = res .. " (" .. script.Uuid:sub(1,8) .. ")"

    if #script.Handlers > 0 then
      res = res .. " handlers: " .. #script.Handlers
    end
    if #script.Commands > 0 then
      res = res .. " commands: " .. #script.Commands
    end
    if #script.Hooks > 0 then
      res = res .. " hooks: " .. #script.Hooks
    end
    client.Notice(source, res)
  end

  return "End of scripts list"
end)
