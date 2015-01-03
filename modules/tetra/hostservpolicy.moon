require "lib/etcd"

use "strings"

-- Allow policy for vhosts to be automatically rejected or activated
--
-- table for each vhost will look like:
--
-- {
--    pattern: -- a lua pattern to match against
--    action:  -- either REJECT or ACTIVATE
--    setter:  -- oper that set that policy
--    date:    -- date the policy was in effect
--    uuid:    -- uuid of the policy
-- }
--

store = etcd.PathStore "vhostpolicy"

math.randomseed os.time!

-- Initialize the schema
if store.data.pc == nil
  store.data.pc = {}
  store.data.num = 0

Hook "HOSTSERV-REQUEST", (nick, vhost) ->
  client.OperLog "HostServ: #{nick} requested #{vhost}"

  --info = tetra.Atheme.NickServ.Info nick

  --if info == nil
  --  return

  --if info.vhost == nil
  --  client.OperLog "#{nick} doesn't have a vhost"
  --else
  --  client.OperLog info.vhost

  for uuid, policy in pairs store.data.pc
    if policy.pattern == nil
      continue

    if vhost\match policy.pattern
      if policy.action == "REJECT"
        tetra.Atheme.HostServ.Reject nick, "Your vhost failed a policy test (#{policy.uuid})"
        client.OperLog "Vhost #{vhost} for #{name} matched test #{policy.pattern} (#{policy.uuid}) and was rejected"
      else if policy.action == "ACTIVATE"
        tetra.Atheme.HostServ.Activate nick
        client.OperLog "Vhost #{vhost} for #{name} matched test #{policy.pattern} (#{policy.uuid}) and was activated"

      break

commands = {
  ADD: {"adds a vhost policy", "VHOSTPOLICY ADD <pattern> <REJECT|ACTIVATE>", 2,
    (source, args) ->
      id = store.data.num + 1
      store.data.num += 1

      action = args[2]\upper!

      if action != "REJECT"
        if action != "ACTIVATE"
          return "VHOSTPOLICY ADD <pattern> <REJECT|ACTIVATE>"

      policy =
        pattern: args[1]
        action: args[2]\upper!
        setter: source.Nick
        date: os.time!
        uuid: id

      store.data.pc[policy.uuid] = policy

      return "Added policy #{policy.uuid} with pattern #{policy.pattern}: #{policy.action}"
  }

  LIST: {"lists policies", "VHOSTPOLICY LIST", 0,
    (source) ->
      client.Notice source, "%-#{#store.data.pc}s | policy   | pattern"\format "id"

      for uuid, policy in pairs store.data.pc
        if policy.uuid
          client.Notice(source, "%-#{#store.data.pc}s | %-#{#"ACTIVATE"}s | %s"\format(policy.uuid, policy.action, policy.pattern))

      return "End list"
    }

  INFO: {"gets information on a specific policy", "VHOSTPOLICY INFO <id>", 1,
    (source, args) ->
      pid = tonumber args[1]
      if store.data.pc[pid] == nil
        return "No such policy"

      policy = store.data.pc[pid]

      client.Notice source, "Information on policy #{pid}"
      client.Notice source, "Pattern: #{policy.pattern}"
      client.Notice source, "Action:  #{policy.action}"
      client.Notice source, "Setter:  #{policy.setter}"
      client.Notice source, "Date:    #{os.date "%c", policy.date}"

      return "End of info"
  }

  DELETE: {"Deletes a policy", "VHOSTPOLICY DELETE <id>", 1,
    (source, args) ->
      pid = tonumber args[1]
      if store.data.pc[pid] == nil
        return "No such policy"

      store.data.pc[pid] = nil

      return "Deleted"
  }
}

Command "VHOSTPOLICY", true, (source, destination, args) ->
  cmdargs = {}
  command = nil

  if #args == 0
    command = nil
  else
    command = args[1]\upper!

  usage = ->
    client.Notice source, "VHOSTPOLICY subcommands: "
    for k,v in pairs commands
      client.Notice source, strings.format("%-10s - %s", k, v[1])
    return "End of command list"

  if not command or command == ""
    return usage!

  if commands[command]
    command = commands[command]

    cmdargs = [i for i in *args[2,]]

    if #cmdargs == command[3]
      return command[4] source, cmdargs
    else
      return command[2]
  else
    return usage!
