local json = require "json"

function dblookup(id)
  local url = "http://derpiboo.ru/"..id..".json?nocomments"

  local c, err = web.get(url)
  if err ~= nil then
    tetra.log.Printf("URL error: %#v", err)
  end

  local str, err = ioutil.readall(c.Body)
  if err ~= nil then
    tetra.log.Printf("Read error: %#v", err)
  end

  str = ioutil.byte2string(str)
  local obj = json.decode(str)
  return obj
end

function summarize(info)
  local ret = "^ Derpibooru: "
  if info.tags.explicit then ret = ret .. "[NSFW] " else ret = ret .. "[SFW] " end
  ret = ret .. "Tags: " .. info.tags

  return ret
end

function db_scrape(line)
  local source, destination, message = parseLine(line)

  if message:find("derpiboo.ru") then
    local id = message:match("/(%d+)")
    local info = dblookup(id)

    client.Privmsg(destination, summarize(info))
  end
end

tetra.script.AddLuaProtohook("PRIVMSG", "db_scrape")
