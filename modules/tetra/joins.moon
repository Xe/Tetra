sqlite3 = require "lsqlite3"

export db = sqlite3.open "var/tetra.db"
export done = false

db\exec [[
  CREATE TABLE IF NOT EXISTS Joins (
    id       INTEGER PRIMARY KEY,
    name     TEXT,
    service  TEXT
  );
]]

db\trace (ud, sql) ->
  log.Printf "SQL: %s", sql

insert_stmt = assert db\prepare "INSERT INTO Joins VALUES (NULL, ?, ?)"
select_stmt = assert db\prepare "SELECT * FROM Joins"

delete = (name, service) ->
  db\exec "DELETE FROM Joins WHERE name='#{name}' AND service='#{service}';"

Command "JOIN", true, (source, destination, message) ->
  parc = #message

  if parc == 0
    return "Cannot join, need channel name or service and channel name"

  service = client
  local chan

  if parc == 2
    tmp = message[1]

    if tetra.Services[tmp] ~= nil
      service = tetra.Services[tmp]
      chan = message[2]

    else
      return "Cannot have #{tmp} join #{message[2]}, #{tmp} does not exist!"

  if parc == 1
    chan = message[1]

  if parc > 2
    return "Too many arguments"

  chan = chan\upper!

  if contains keys(service.Channels), chan
    return "#{service.Nick} is already in #{chan}, cannot join again!"

  if tetra.Channels[chan] == nil
    return "Cannot join #{chan} as it does not exist."

  service.Join(chan)

  do
    insert_stmt\bind_values chan, service.Kind
    insert_stmt\step!
    insert_stmt\reset!

  return "Joined #{service.Nick} to #{chan}"

Command "PART", true, (source, destination, message) ->
  parc = #message

  if parc == 0
    return "Cannot part, need channel name or service and channel name"

  service = client
  local chan

  if parc == 2
    tmp = message[1]

    if tetra.Services[tmp] ~= nil
      service = tetra.Services[tmp]
      chan = message[2]

    else
      return "Cannot have #{tmp} part #{message[2]}, #{tmp} does not exist!"

  if parc == 1
    chan = message[1]

  if parc > 2
    return "Too many arguments"

  chan = chan\upper!

  if not contains keys(service.Channels), chan
    return "#{service.Nick} is not in #{chan}, cannot part!"

  if tetra.Channels[chan] == nil
    return "Cannot part #{chan} as it does not exist."

  service.Part(chan)

  delete chan, service.Kind

  return "Parted #{service.Nick} from #{chan}"

Protohook "PING", (line) ->
  if done
    return

  for row in select_stmt\nrows!
    svc = tetra.Services[row.service]

    print "#{svc.Nick} is joining #{row.name}"
    svc.Join row.name

  done = true

Hook "SHUTDOWN", ->
  db\close!
