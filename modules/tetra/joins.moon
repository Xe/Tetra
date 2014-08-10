export db = FooDB "var/autojoin.json"
export done = false

export joincmd = elevated! .. (source, destination, message) ->
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

  db\Commit!

  return "Joined #{service.Nick} to #{chan}"

export partcmd = elevated! .. (source, destination, message) ->
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
    return "#{service.Nick} is not in #{chan}, cannot join again!"

  if tetra.bot.Channels[chan] == nil
    return "Cannot part #{chan} as it does not exist."

  service.Part(chan)

  do
    idx = find(db.data[service.Kind], chan)
    table.remove(db.data[service.Kind], idx

  db\Commit!

  return "Joined #{service.Nick} to #{chan}"

tetra.script.AddLuaCommand "JOIN", "joincmd"
tetra.script.AddLuaCommand "PART", "partcmd"
