-- Based on:
-- https://github.com/TheLinx/Juiz/blob/master/modules/ytvlookup.lua

local json = require "json"
local http = require "socket.http"

function ytlookup(id)
  local res = getjson("http://gdata.youtube.com/feeds/api/videos/"..id.."?alt=json&fields=author,title")
  local author, title = res.entry.author[1].name["$t"], res.entry.title["$t"]

  return "^ Youtube: " .. title .. " - Uploaded by: " .. author
end

Hook("CHATBOT-CHANMSG", function(source, destination, message)
  message = table.concat(luar.slice2table(message), " ")

  if message:find("youtube%.com/watch") then
    client.Privmsg(destination, ytlookup(message:match("v=(...........)")))
  elseif message:find("youtu%.be/") then
    client.Privmsg(destination, ytlookup(message:match("%.be/(...........)")))
  else return end
end)

Command("YT", function(source, destination, message)
  if #message < 1 then
    return "Params: string to search youtube for"
  end

  local search = table.concat(luar.slice2table(message, " "))

  local info = getjson("https://gdata.youtube.com/feeds/api/videos?q=" .. search .. "&v=2&alt=jsonc")
  local video = info.data.items[1]

  return "Youtube: " .. video.title .. " http://youtu.be/" .. video.id
end)
