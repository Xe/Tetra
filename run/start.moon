yaml = require "yaml"

local config
local cpath
local host
port = "6667"

get_or_fail = (name, message) ->
  with var = os.getenv name
    if var == nil
      error "Cannot find " .. message
    else
      print "Loaded #{name}: #{var}"
      return var

-- Figure out where the config is.
with cpath = os.getenv "TETRA_CONFIG_PATH"
  if cpath == nil
    print "Using default config path"
    cpath = "/app/etc/config.yaml.example"

print "Loading config from " .. cpath

with fin, message = io.open cpath, "r"
  if fin == nil
    error cpath .. " unreadable! " .. message

  input = \read "*a"
  config = yaml.load input
  \close!

with host = get_or_fail "IRCD_PORT_6667_TCP_ADDR", "ircd address"
  config.uplink.host = host
  config.uplink.port = port

with influxhost = get_or_fail "TETRA_INFLUX_HOST", "influxdb host"
  config.stats.host = influxhost

print "Using ircd at " .. host .. " port " .. port

print "Writing config"

with fout = io.open cpath, "w"
  \write yaml.dump config
  \close!

with indocker = os.getenv "TETRA_DOCKER"
  if indocker == nil
    os.execute "Tetra"
  else
    cmd = "/bin/sh -c 'cd /app; /app/Tetra'"
    if cpath == "/app/etc/config.yaml.example"
      cmd = "/bin/sh -c 'TETRA_CONFIG_PATH=" .. cpath .. " cd /app; /app/Tetra'"
    os.execute cmd

