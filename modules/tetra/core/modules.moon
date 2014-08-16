require "lib/etcd"

db = etcd.Store "modules"

if db.data.loads == nil
  db.data.loads = {}
else
  for _, script in pairs db.data.loads
    tetra.bot.LoadScript(script)

Command "LOAD", true, (src, dest, msg) ->
  if #msg == 0
    return "Need a script name"

  name = msg[1]
  script, err = tetra.bot.LoadScript(name)

  if err ~= nil
    tetra.log.Printf("Can't load script " .. name .. ": %#v", err)
    return "Script #{name} failed load: #{err}"

  table.insert db.data.loads, 1, name

  "#{name} loaded."

Command "UNLOAD", true, (src, dest, msg) ->
  if #msg == 0
    return "Need a script name"

  name = msg[1]

  if tetra.bot.Scripts[name] == nil
    return "#{name} is not loaded."

  if name == script.Name
    return "Cannot unload this script!"

  table.remove db.data, find(db.data name)

  sleep(0.5)

  script, err = tetra.bot.UnloadScript(name)

  "#{name} unloaded."
