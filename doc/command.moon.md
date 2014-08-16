## Command Tutorial

In this file we describe an example command `TEST`. `TEST` will return some 
information about the place the command is used as well as explain the 
arguments involved.

Because Tetra is a polyglot of Lua, Moonscript and Go, the relevant Go objects 
will have their type definitions linked to on [godoc](http://godoc.org)

Declaring commands is done with the `Command` macro. It takes in two arguments.

1. The command verb
2. The command function

It also can take in 3 arguments if the command needs to be restricted to IRCops 
only.

1. The command verb
2. `true`
3. The command function

The command function can have up to 3 arguments set when it is called. These 
are:

1. The [Client](https://godoc.org/github.com/Xe/Tetra/bot#Client) that 
   originated the command call.
2. The [Destination](https://godoc.org/github.com/Xe/Tetra/bot#Targeter) or 
   where the command was sent to. This will be a Client if the target is an 
   internal client or 
   a [Channel](https://godoc.org/github.com/Xe/Tetra/bot#Channel) if the target 
   is a channel.
3. The command arguments as a string array.

```moonscript
Command "TEST", (source, destination, args) ->
```

All scripts have `client` pointing to the pseudoclient that the script is 
spawned in. If the script name is `chatbot/8ball`, the value of `client` will 
point to the `chatbot` pseudoclient.

```moonscript
  client.Notice source, "Hello there!"
```

This will send a `NOTICE` to the source of the command saying "Hello there!".

```moonscript
  client.Notice source, "You are #{source.Nick} sending this to #{destination.Target!} with #{#args} arguments"
```

All command must return a string with a message to the user. This is a good 
place to do things like summarize the output of the command or if it worked or 
not. If the command is oper-only, this will be the message logged to the 
services snoop channel.

```moonscript
  "End of TEST output"
```
