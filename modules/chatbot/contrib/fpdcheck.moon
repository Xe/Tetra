socket = require "socket"

-- Automatically checks all servers on the network every 5 minutes.
-- There is no way to disable this behavior without unloading this script.

Hook "CRON-HEARTBEAT", ->
  for _, server in pairs tetra.Servers
    if server.Count > 50 -- don't scan hubs
      try
        main: ->
          with socket.connect server.Name, 8430
            assert \read "*a"
        except: (e) ->
          client.OperLog "#{server.Name} is not replying to flash policy requests."
