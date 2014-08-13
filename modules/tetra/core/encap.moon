-- :752 ENCAP * GCAP :QS EX IE KLN UNKLN ENCAP SERVICES EUID EOPMOD
Hook "ENCAP-GCAP", (source, caps) ->
  server = tetra.bot.Servers[source]
  server.Capab = caps

-- :6YK ENCAP * METADATA SET 7RT100001 CLOAKEDHOST :yolo-swag.com
Hook "ENCAP-METADATA", (source, args) ->
  action = args[1]
  target = args[2]
  key = args[3]
  value = args[4]

  switch action
    when "SET", "ADD"
      if strings.first(target) == "#"
        target = tetra.bot.Channels[target\upper!]
      else
        target = tetra.bot.Clients.ByUID[target]

      target.Metadata[key] = value

    when "DELETE", "CLEAR"
      if strings.first(target) == "#"
        target = tetra.bot.Channels[target\upper!]
      else
        target = tetra.bot.Clients.BuUID[target]

      target.Metadata[key] = nil

-- :7RT100001 ENCAP * CERTFP :6d73b6c3-039e-40a3-a61f-db1e76d83ca2
Hook "ENCAP-CERTFP", (source, args) ->
  tetra.bot.Clients.ByUID[source].Certfp = args[1]

-- :00A ENCAP * SU 7RT100002 :Tetra
Hook "ENCAP-SU", (source, args) ->
  target = args[1]
  account = args[2] if args[2] ~= nil else "*"
  tetra.bot.Clients.ByUID[target].Account = account

-- :42F ENCAP * SNOTE s :Failed OPER attempt - host mismatch by xena (xena@0::z)
Hook "ENCAP-SNOTE", (source, args) ->
  client.ServicesLog("Server notice #{args[1]}: #{args[2]}")
