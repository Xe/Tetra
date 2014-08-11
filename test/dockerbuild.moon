-- Test docker building and running

make = (directory, rule) ->
  with retcode = os.execute "make -C #{directory} #{rule}"
    if retcode ~= 0
      error "retcode: #{retcode}"
    else
      return retcode

build = (directory, rules) ->
  for _, rule in pairs rules
    make directory, rule .. " || true"

capture = (command) ->
  print "> " .. command
  proc = io.popen command
  proc\read("*a")

build "..", {"clean", "build", "docker-build"}
build "testnet/ircd", {"build", "kill"}

ircd_id = capture "make --no-print-directory -C testnet/ircd run 2>/dev/null"
tetra_id = capture "docker run -dit --link tetra-ircd:ircd --name tetra xena/tetra"

print "ircd id:  #{ircd_id\sub 1,10}"
print "tetra id: #{tetra_id\sub 1,10}"

with proc = io.popen "docker logs -f #{tetra_id\sub 1,10}"
  pingcount = 0

  for line in \lines!
    print line

    if line\match "ERROR"
      print "There was an error :("
      capture "docker rm -f tetra tetra-ircd"

      os.exit(2)

    if line\match "PING"
      pingcount += 1
      print "Ping count: #{pingcount}"

      if pingcount > 1
        print "Tests passed"
        capture "docker rm -f tetra tetra-ircd"
        os.exit(0)

