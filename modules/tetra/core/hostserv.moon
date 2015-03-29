use "strings"

Hook "HOSTSERV-SERVICELOG", (message) ->
  -- (@HostServ) Xena REQUEST: ninjas
  if message[2]\find "REQUEST"
    tetra.RunHook "HOSTSERV-REQUEST", message[1], strings.shuck message[3]
    return

  -- (@HostServ) Heartmender (Quora) REQUEST: a.vhost.like.this
  if message[3]\find "REQUEST"
    tetra.RunHook "HOSTSERV-REQUEST", strings.shuck(message[2]), strings.shuck message[4]
    return

  client.ServicesLog strings.format "core/hostserv: I didn't understand %q", message
