require "lib/etcd"

db = etcd.Store "tells"

Command "TELL", (source, destination, message) ->
  if not destination.IsChannel!
    "Command must be run from inside a channel"

  if #message < 2
    return "TELL: <person> a message"

  rec = message[1]\upper!

  if rec == source.Nick\upper!
    return "Cannot send a message to yourself!"

  if db.data[destination.Target!] == nil
    db.data[destination.Target!] = {}

  if db.data[destination.Target!][rec] == nil
    db.data[destination.Target!][rec] = {}

  if #db.data[destination.Target!][rec] > 10
    return "Too many messages for that person"

  msg = ""

  for i=2, #message
    msg ..= message[i] .. " "

  note = "<#{source.Nick}> #{msg}"

  table.insert(db.data[destination.Target!][rec], 1, note)
  db\Commit!

  "I will try to let them know."

Hook "CHATBOT-CHANMSG", (source, destination, message) ->
  if db.data[destination.Target!] == nil
    return

  if db.data[destination.Target!][source.Nick\upper!] ~= nil
    client.Notice source, "You've got mail!"

    for i, message in pairs db.data[destination.Target!][source.Nick\upper!]
      client.Notice source, "#{i}: #{message}"

    db.data[destination.Target!][source.Nick\upper!] = nil
    db\Commit!
