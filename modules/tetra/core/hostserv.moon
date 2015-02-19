use "strings"

Hook "HOSTSERV-SERVICELOG", (message) ->
  -- (@HostServ) Xena REQUEST: ninjas
  tetra.RunHook "HOSTSERV-REQUEST", message[1], strings.shuck(message[3]) if message\match "REQUEST"
