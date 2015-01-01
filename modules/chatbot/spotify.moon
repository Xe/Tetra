Hook "CHATBOT-CHANMSG", (source, destination, message) ->
  message = table.concat luar.slice2table(message), " "

  if message\find "spotify:track"
    client.Privmsg(destination, tracklookup message\match "spotify:track:(.+)")
    return
  if message\find "open.spotify.com/track/"
    client.Privmsg(destination, tracklookup message\match "open.spotify.com/track/(.+)")
    return

export tracklookup = (id) ->
  res, err = getjson "https://api.spotify.com/v1/tracks/" .. id
  if err ~= nil
    return err.Error!

  return "^ Spotify: #{res.name} by #{res.artists[1].name}"
