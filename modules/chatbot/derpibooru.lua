function dblookup(id)
  local url = "http://derpiboo.ru/"..id..".json?nocomments"

  local obj, err = getjson(url)

  if err ~= nil then
    return nil, err
  end

  return obj, nil
end

function summarize(info)
  local ret = "^ Derpibooru: "

  if contains(info.tag_ids, "explicit") then
    ret = ret .. "[NSFW] "
  else
    ret = ret .. "[SAFE] "
  end

  ret = ret .. "Tags: " .. info.tags

  return ret
end

Hook("CHATBOT-CHANMSG", function(source, destination, message)
  message = strings.join(message, " ")
  if message:find("derpiboo.ru/(%d+)") then
    local id = message:match("/(%d+)")

    if id == nil then
      return
    end

    local info, err = dblookup(id)

    if err ~= nil then
      client.Privmsg(destination, "Could not look up that image. Does it exist?")
      return
    end

    client.Privmsg(destination, summarize(info))
  end
end)
