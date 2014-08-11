channel_data = {}

scan_replace = protohook("PRIVMSG") .. function(line)
  local source, destination, message = parseLine(line)

  if not is_common_channel(destination) then return end

  target = destination.Target()

  if target == "#SERVICES" then return end

  if channel_data[target] == nil then
    channel_data[target] = LimitQueue(30)
    client.ServicesLog("Tracking messages for " .. destination.Target())
  end

  local cqueue = channel_data[target]

  if message:sub(1,1) == "s" and message:sub(2,2) == "/" then
    -- User is doing a search/replace
    local mymsg = message:sub(3)
    local pattern = mymsg:match("(.+)/")
    local replacement = mymsg:match("/(.+)")

    for i = #cqueue.table, 1, -1 do
      local line = cqueue.table[i]
      local rep = string.gsub(line.Line, pattern, replacement)

      if rep ~= line.Line then
        if source.Nick == line.Nick then
          line.Line = rep
        end

        client.Privmsg(destination, "<" .. line.Nick .. "> ".. rep)
        break
      end
    end

  else
    cqueue:Add {
      Nick = source.Nick,
      Line = message,
    }
  end
end

replay = protohook("JOIN") .. function(line)
  local source = tetra.bot.Clients.ByUID[line.Source]
  local channel = line.Args[2]:upper()

  if channel_data[channel] ~= nil then
    client.Notice(source, "Replaying last 30 lines of chat in " .. line.Args[2] .. " for you")
    for i, line in pairs(channel_data[channel].table) do
      client.Notice(source, i .. ": <" .. line.Nick .. "> ".. line.Line)
    end
  end
end
