-- Based on:
-- https://github.com/TheLinx/Juiz/blob/master/modules/ytvlookup.lua

local json = require "json"
local http = require "socket.http"

function ytlookup(id)
  local c = http.request("http://gdata.youtube.com/feeds/api/videos/"..id.."?alt=json&fields=author,title")
  local res = json.decode(c)
  local author, title = res.entry.author[1].name["$t"], res.entry.title["$t"]

  return "^ Youtube - " .. title .. " Posted by: " .. author
end

youtube_scrape = hook("CHATBOT-CHANMSG") .. function(source, destination, message)
  message = strings.join(message, " ")

  if message:find("youtube%.com/watch") then
    client.Privmsg(destination, ytlookup(message:match("v=(...........)")))
  elseif message:find("youtu%.be/") then
    client.Privmsg(destination, ytlookup(message:match("%.be/(...........)")))
  else return end
end

tetra.script.AddLuaProtohook("PRIVMSG", "youtube_scrape")

