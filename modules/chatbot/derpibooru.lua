local json = require "json"

function dblookup(id)
  local url = "http://derpiboo.ru/"..id..".json?nocomments"

  local c, err = web.get(url)
  if err ~= nil then
    tetra.log.Printf("URL error: %#v", err)
    return nil, err
  end

  local str, err = ioutil.readall(c.Body)
  if err ~= nil then
    tetra.log.Printf("Read error: %#v", err)
    return nil, err
  end

  str = ioutil.byte2string(str)
  local obj = json.decode(str)
  return obj, nil
end

function summarize(info)
  local ret = "^ Derpibooru: "
  ret = ret .. "Tags: " .. info.tags

  return ret
end

function db_scrape(line)
  local source, destination, message = parseLine(line)

  if message:find("derpiboo.ru") then
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
end

tetra.script.AddLuaProtohook("PRIVMSG", "db_scrape")
