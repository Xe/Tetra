etlua = require "etlua"

math.randomseed os.time!

local config

fin = io.open "/home/ircd/template/ircd.conf", "r"
config = fin\read "*a"
fin\close!

settings =
  hostname: os.getenv("HOSTNAME")\upper!
  sid: math.random 100, 999

template = etlua.compile config

print template settings

