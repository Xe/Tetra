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
cpath = os.getenv "TETRA_CONFIG_PATH"
if cpath == nil
  print "Using default config path"
  os.execute "cp /app/etc/config.yaml.example /app/etc/config.yaml"
  cpath = "/app/etc/config.yaml"

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

with etcd = get_or_fail "ETCD_PORT_4001_TCP_ADDR", "etcd host"
  config.etcd.machines = { "http://#{etcd}:4001" }

--with influxhost = get_or_fail "TETRA_INFLUX_HOST", "influxdb host"
--  config.stats.host = influxhost

config.stats.host = "NOCOLLECTION"

print "Using ircd at " .. host .. " port " .. port

print "Writing config"

with fout = io.open cpath, "w"
  \write yaml.dump config
  \close!

cmd = "/bin/sh -c 'TETRA_CONFIG_PATH=" .. cpath .. " cd /app; /app/Tetra'"

os.execute cmd

