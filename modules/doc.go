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

All modules will also have base.lua loaded.
*/
package modules
