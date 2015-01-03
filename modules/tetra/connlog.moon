sqlite3 = require "lsqlite3"

use "strings"

export db = sqlite3.open "var/tetra.db"
export done = false

db\exec [[
  CREATE TABLE IF NOT EXISTS Connections (
    id          INTEGER PRIMARY KEY,
    date        INTEGER,
    server      TEXT,
    nick        TEXT,
    ident       TEXT,
    ip          TEXT,
    reversedns  TEXT,
    cloakedhost TEXT,
    gecos       TEXT
  );
]]

db\exec [[ BEGIN TRANSACTION; ]]

insert_stmt = assert db\prepare "INSERT INTO Connections VALUES (NULL, ?, ?, ?, ?, ?, ?, ?, ?)"
select_stmt = assert db\prepare "SELECT * FROM Connections"

Hook "SHUTDOWN", ->
  db\exec [[ COMMIT; ]]
  db\close!

Hook "CRON-HEARTBEAT", ->
  db\exec [[
    COMMIT;
    BEGIN TRANSACTION;
  ]]

Hook "NEWCLIENT", (c) ->
  insert_stmt\bind_values os.time!, c.Server.Name, c.Nick, c.User, c.Ip, c.Host, c.VHost, c.Gecos
  insert_stmt\step!
  insert_stmt\reset!

commands = {
  GREP: {"searches logs", "CONNLOG GREP <field> <matcher>", 2,
    (source, args) ->
      field = args[1]\lower!
      matcher = table.concat [i for i in *args[2,]], " "

      if field != "ip"
        if field != "nick"
          if field != "server"
            if field != "ident"
              if field != "reversedns"
                if field != "cloakedhost"
                  if field != "gecos"
                    return "invalid field"

      count = 0

      for row in db\nrows "SELECT * FROM Connections WHERE #{field}='#{matcher}';"
        count += 1
        client.Notice source, "#{os.date "%c", row.date} #{row.nick} #{row.ident} #{row.ip} #{row.reversedns} #{row.cloakedhost} #{row.gecos}"

      "Found #{count} matches"
  },
}

Command "CONNLOG", true, (source, destination, args) ->
  cmdargs = {}
  command = nil

  if #args == 0
    command = nil
  else
    command = args[1]\upper!

  usage = ->
    client.Notice source, "CONNLOG subcommands: "
    for k,v in pairs commands
      client.Notice source, strings.format("%-10s - %s", k, v[1])
    return "End of command list"

  if not command or command == ""
    return usage!

  if commands[command]
    command = commands[command]

    cmdargs = [i for i in *args[2,]]

    if #cmdargs > command[3]-1
      return command[4] source, cmdargs
    else
      return command[2]
  else
    return usage!
