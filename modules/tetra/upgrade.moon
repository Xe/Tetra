export logOrDie = (command, kind) ->
  client.ServicesLog "$ " .. command

  proc = io.popen command
  for line in proc\lines!
    client.ServicesLog "#{kind}: #{line}"
    tetra.bot.RunHook "UPGRADE-"..kind, line

  if not proc\close!
    client.ServicesLog "Error in upgrade process, aborting!"
    return false
  true

Command "UPGRADE", true, (source) ->
  for _, group in pairs { {"git pull", "GIT"}, {"go get -u -v . 2>&1", "GOLANG"}, {"make build", "BUILD"}}
    if not logOrDie group[1], group[2]
      return "Upgrade failed"

  client.ServicesLog "#{source.Nick}: UPGRADE: Upgraded to latest version."

  return "Done."
