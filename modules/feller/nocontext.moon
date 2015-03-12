use "strings"

sqlite3 = require "lsqlite3"

export db = sqlite3.open "var/feller.db"

Hook "JOINCHANNEL", (cu) ->
  chan = cu.Channel
  myc  = cu.Client

  return if myc.Account == "*"

  if client.Channels[chan.Target!]
    print "Found matching channel for #{myc.Nick} (#{myc.Account}) at #{chan.Target!}"

    for line in db\nrows "SELECT * FROM Chatlines WHERE account='#{myc.Account}' AND channel='#{chan.Target!}' ORDER BY RANDOM() LIMIT 1"
      print strings.format "%#v", line
      client.Privmsg chan, "[#{myc.Nick}] #{line.message}"
