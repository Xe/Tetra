/*
Package modules is a collection of lua scripts that represent loadable modules for
Tetra. Modules currently are not required to have any functions, but will have
the following globals available:

    client:         the pseudoservice client that the script is running under
    tetra.script:   a reference to the script object
    tetra.log:      a reference to the script logger
    tetra.bot:      a reference to the Tetra god object
    uuid.new:       a UUID generator for convenience
    web.get:        Go's http.Get
    web.post:       Go's http.Post
    ioutil.readall: convenience wrapper
    ".byte2sting:   converts C strings to Go strings

All modules will also have base.lua loaded. Moonscript modules need to load this
manually for now.

Modules may be written in either lua or moonscript. If there is a name conflict
the lua file will be preferred over the moonscript one.

An example moonscript module is as follows:

    Protohook "PRIVMSG", (line) ->
      source, destination, message = parseLine line
      print "#{destination.Target!} <#{source.Nick}> #{table.concat message, " "}"

Please note that handler/command functions myst be exported for Tetra to be able
to use them. This is a moonscript-specific problem.

An example lua module is as follows:

    Command("PING" .. function(client, target, message)
      return "PONG"
    end)

This package will compile but is not useful for anything but documentation.
*/
package modules

import (
	"github.com/Xe/Tetra/bot"
)

// Command initializes and returns a new bot command in Lua space. It is
// referenced by a type-4 UUID. Parameters are the command verb, if the command
// is oper-only or not, and the lua function that represents the command.
func Command(verb string, operonly bool, function func()) *tetra.Command { return nil }

// Hook wraps a function to act as a named, event-like hook. It takes in the hook
// verb and the function to run as a hook.
func Hook(verb string, function func()) *tetra.Hook { return nil }

// Protohook wraps a function to be called on a protocol verb. It takes in
// the protocol verb and the function to call.
func Protohook(verb string, function func()) *tetra.Handler { return nil }
