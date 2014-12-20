-- Takes in a string and outputs a trimmed string
trim = (str) ->
  str\gsub "^%s*(.-)%s*$", "%1"

--- Upgrade modules automatically
Hook "UPGRADE-GIT", (line) ->
  line = trim line

  --  modules/tetra/upgrade.moon | 12 ++++++++++++
  if not line\find "|"
    return

  line = trim strings.split(line, "|")[1]

  -- modules/tetra/upgrade.moon
  if line\sub(1, 8) ~= "modules/"
    return -- We can't upgrade non-scripts

  line = line\sub 9

  -- tetra/upgrade.moon
  line = strings.split(line, ".")

  -- tetra/upgrade
  modname = line[1]

  if tetra.Scripts[modname] == nil
    return -- Can't upgrade a script that is not loaded

  if modname == script.Name -- Don't upgrade self, currently has issues.
    return

  if modname == "tetra/upgrade"
    client.ServicesLog "UPDATER: Skipping tetra/upgrade"
    return

  err = tetra.UnloadScript modname
  if err ~= nil
    client.ServicesLog "UPDATER: unload error: #{err}"
    return

  sleep 0.5

  s, err = tetra.LoadScript modname
  if err ~= nil
    client.ServicesLog "UPDATER: load error: #{err}"
    return

  client.ServicesLog "UPDATER: updated #{modname}"
