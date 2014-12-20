-- Example hook via YO

Hook "YO", (source, dest) ->
  yo, err = tetra.GetYo dest

  client.ServicesLog "#{dest} got a yo from #{source}!"

  if err ~= nil
    print err
    client.ServicesLog "Could not get Yo client for #{dest}"
  else
    yo.YoUser source
    client.ServicesLog "Yo'd #{source} back!"
