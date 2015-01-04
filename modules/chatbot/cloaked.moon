use "charybdis"
use "crypto"

Command "CLOAKED", (source, destination, args) ->
  if #args < 1
    return "CLOAKED <nick>"

  target = tetra.Clients.ByNick[args[1]\upper!]

  if target == nil
    return "No such client #{args[1]}"

  client.Notice source, "Connection ID:       #{crypto.fnv target.Uid}"
  client.Notice source, "Cloaked IP:          #{charybdis.cloakip target.Ip}"

  if target.Ip != target.Host
    client.Notice source, "Cloaked reverse DNS: #{charybdis.cloakhost target.Host}"

  "Done"
