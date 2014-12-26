Command "BTC", ->
  info = getjson("https://www.bitstamp.net/api/ticker/")
  "Bitstamp prices (USD): $"..info.ask.." average, $"..info.low.." low, $"..info.high.." high"
