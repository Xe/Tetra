-- Shortens long URL's with bit.ly
-- Needs a bit.ly api key

function url_encode(str)
  if (str) then
    str = string.gsub (str, "\n", "\r\n")
    str = string.gsub (str, "([^%w %-%_%.%~])",
    function (c) return string.format ("%%%02X", string.byte(c)) end)
    str = string.gsub (str, " ", "+")
  end
  return str
end

function scrapeurl(str)
  return str:match("https?://[%w-_%.%?%.:/%+=&]+")
end

API_KEY = tetra.bot.Config.ApiKeys.bitly

URL = "https://api-ssl.bitly.com/v3/shorten?access_token=" .. API_KEY .. "&longUrl="

function shorten(url)
  eurl = URL .. url_encode(url)
  res = getjson(eurl)

  if res.status_txt ~= "OK" then
    client.ServicesLog("url \"" .. url .. "\" failed to shorten: " .. res.status_txt)
  end

  return res.data.url
end

shorten_if_long = hook("CHATBOT-CHANMSG") .. function(source, destination, message)
  local url = scrapeurl(message)

  if url ~= nil then
    if #url > 100 then
      client.Privmsg(destination, "^ " .. shorten(url))
    end
  end
end
