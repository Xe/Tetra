-- Based on:
-- https://github.com/TheLinx/Juiz/blob/master/modules/ytvlookup.lua

local json = require "json"
local http = require "socket.http"

function ytlookup(id)
  local res = getjson("http://gdata.youtube.com/feeds/api/videos/"..id.."?alt=json&fields=author,title")
  local author, title = res.entry.author[1].name["$t"], res.entry.title["$t"]

  return "^ Youtube - " .. title .. " Posted by: " .. author
end

Hook("CHATBOT-CHANMSG", function(source, destination, message)
  message = table.concat(luar.slice2table(message), " ")

  if message:find("youtube%.com/watch") then
    client.Privmsg(destination, ytlookup(message:match("v=(...........)")))
  elseif message:find("youtu%.be/") then
    client.Privmsg(destination, ytlookup(message:match("%.be/(...........)")))
  else return end
end)
