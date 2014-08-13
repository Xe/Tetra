Command "DOGE", ->
  info = getjson "http://pubapi.cryptsy.com/api.php?method=singlemarketdata&marketid=132"
  info = info["return"]["markets"]["DOGE"]

  "One Dogecoin is worth #{info.lasttradeprice} BTC"
