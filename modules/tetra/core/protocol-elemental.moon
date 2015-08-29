use "strings"

Protohook "NICK", (line) ->
  source = tetra.Clients.ByUID[line.Source]
  tetra.Clients.ChangeNick source, line.Args[1]
  source.Nick = line.Args[1]

Protohook "CHGHOST", (line) ->
  client = tetra.Clients.ByUID[line.Args[1]]
  tetra.RunHook "CHGHOST", client, line.Args[2]
  client.VHost = line.Args[2]

Protohook "QUIT", (line) ->
  client = tetra.Clients.ByUID[line.Source]
  tetra.RunHook "CLIENTQUIT", client
  tetra.Clients.DelClient client

  for _, chan in pairs client.Channels
    chan.DelChanUser client

  client.Server.DelClient!

Hook "ENCAP-SU", (uid, args) ->
  accname = args[2] or "*"
  target = tetra.Clients.ByUID[args[1]]
  target.Account = accname
