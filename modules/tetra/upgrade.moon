Command "UPGRADE", (source) ->
  with proc = io.popen "git pull"
    for line in proc\lines!
      client.ServicesLog "GIT  : " .. line

  with proc = io.popen "make build"
    for line in proc\lines!
      client.ServicesLog "BUILD: " .. line

  client.ServicesLog "#{source.Nick}: UPGRADE: Upgraded to latest version."

  return "Done."
