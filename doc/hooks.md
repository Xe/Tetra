## Hooks

Tetra has hooks for scripts to listen for various events. This can be used as 
a way for scripts to inter-communicate, but you should not do this.

Tetra has a few builtin hooks to ease development.

### `$CLIENTNAME-CHANMSG`

This will be called when a channel adjacent to `$CLIENTNAME` gets a `PRIVMSG`. 
The arguments are:

 - Source Client
 - Destination Channel
 - Message as an array of strings

```moonscript
Hook "CHATBOT-CHANMSG", (src, dest, msg) ->
  print "#{dest.Name} <#{src.Nick}> #{strings.join msg, " "}"
```

### `CRON-HEARTBEAT`

This is a hook that is run every 5 minutes and has no arguments.

```moonscript
Hook "CRON-HEARTBEAT", ->
  print "It's been 5 minutes!"
```

### `SHUTDOWN`

This is a hook that is run when Tetra is shutting down and has no arguments.

```moonscript
Hook "SHUTDOWN", ->
  print "Goodbye!"
```

### `ENCAP-$VERB`

This is a hook that is run when the uplink sends an `ENCAP` message. `$VERB` is 
the verb that was encapped.

```moonscript
-- :7RT100001 ENCAP * CERTFP :6d73b6c3-039e-40a3-a61f-db1e76d83ca2
Hook "ENCAP-CERTFP", (source, args) ->
  tetra.bot.Clients.ByUID[source].Certfp = args[1]
```

### `NEWCLIENT`

Called when a new client joins to the network. Argument is the client struct.

### `CLIENTQUIT`

Called when a client quits from the network. Argument is the client struct.
