sqlite3 = require "lsqlite3"

export sdb = assert sqlite3.open "var/tetra.db"

sdb\exec [[
  CREATE TABLE IF NOT EXISTS Scripts (
    id   INTEGER PRIMARY KEY,
    name TEXT UNIQUE
  );
]]

select_stmt = assert sdb\prepare "SELECT * FROM Scripts"
insert_stmt = assert sdb\prepare "INSERT INTO Scripts VALUES (NULL, ?)"
delete_stmt = assert sdb\prepare "DELETE FROM Scripts WHERE name = ?"

for row in select_stmt\nrows!
    tetra.LoadScript row.name
    log.Printf "loaded %s", row.name

Command "LOAD", true, (src, dest, msg) ->
  if #msg == 0
    return "Need a script name"

  name = msg[1]
  script, err = tetra.LoadScript(name)

  if err ~= nil
    tetra.log.Printf("Can't load script " .. name .. ": %#v", err)
    return "Script #{name} failed load: #{err}"

  do
    insert_stmt\bind_values name
    insert_stmt\step!
    insert_stmt\reset!

  "#{name} loaded."

Command "UNLOAD", true, (src, dest, msg) ->
  if #msg == 0
    return "Need a script name"

  name = msg[1]

  if tetra.Scripts[name] == nil
    return "#{name} is not loaded."

  if name == script.Name
    return "Cannot unload this script!"

  do
    delete_stmt\bind_values name
    delete_stmt\step!
    delete_stmt\reset!

  sleep(0.5)

  script, err = tetra.UnloadScript(name)

  "#{name} unloaded."
