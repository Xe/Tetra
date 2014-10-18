require "lib/etcd"

export db = etcd.Store "joins"
export done = false

Command "JOIN", true, (source, destination, message) ->
  parc = #message

  if parc == 0
    return "Cannot join, need channel name or service and channel name"

  service = client
  local chan

  if parc == 2
    tmp = message[1]

    if tetra.bot.Services[tmp] ~= nil
      service = tetra.bot.Services[tmp]
      chan = message[2]

    else
      return "Cannot have #{tmp} join #{message[2]}, #{tmp} does not exist!"

  if parc == 1
    chan = message[1]

  if parc > 2
    return "Too many arguments"

  chan = chan\upper!

  if contains keys(service.Channels), chan
    return "#{service.Nick} is already in #{chan}, cannot join again!"

  if tetra.bot.Channels[chan] == nil
    return "Cannot join #{chan} as it does not exist."

  service.Join(chan)

  if db.data[service.Kind] == nil
    db.data[service.Kind] = {chan}
  else
    table.insert db.data[service.Kind], chan

  return "Joined #{service.Nick} to #{chan}"

Command "PART", true, (source, destination, message) ->
  parc = #message

  if parc == 0
    return "Cannot part, need channel name or service and channel name"

  service = client
  local chan

  if parc == 2
    tmp = message[1]

    if tetra.bot.Services[tmp] ~= nil
      service = tetra.bot.Services[tmp]
      chan = message[2]

    else
      return "Cannot have #{tmp} part #{message[2]}, #{tmp} does not exist!"

  if parc == 1
    chan = message[1]

  if parc > 2
    return "Too many arguments"

  chan = chan\upper!

  if not contains keys(service.Channels), chan
    return "#{service.Nick} is not in #{chan}, cannot part!"

  if tetra.bot.Channels[chan] == nil
    return "Cannot part #{chan} as it does not exist."

  service.Part(chan)

  do
    idx = find db.data[service.Kind], chan
    table.remove db.data[service.Kind], idx

  return "Parted #{service.Nick} from #{chan}"

Protohook "PING", (line) ->
  if done
    return

  for name, channels in pairs db.data
    service = tetra.bot.Services[name]

    for i, chan in pairs channels
      print "#{service.Kind} joining #{chan}"
      service.Join(chan)

  done = true
