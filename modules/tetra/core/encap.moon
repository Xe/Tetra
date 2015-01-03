-- :6YK ENCAP * METADATA SET 7RT100001 CLOAKEDHOST :yolo-swag.com
Hook "ENCAP-METADATA", (source, args) ->
  action = args[1]
  target = args[2]
  key = args[3]
  value = args[4]

  switch action
    when "SET", "ADD"
      if strings.first(target) == "#"
        target = tetra.Channels[target\upper!]
      else
        target = tetra.Clients.ByUID[target]

      target.Metadata[key] = value

    when "DELETE", "CLEAR"
      if strings.first(target) == "#"
        target = tetra.Channels[target\upper!]
      else
        target = tetra.Clients.BuUID[target]

      target.Metadata[key] = nil

-- :7RT100001 ENCAP * CERTFP :6d73b6c3-039e-40a3-a61f-db1e76d83ca2
Hook "ENCAP-CERTFP", (source, args) ->
  tetra.Clients.ByUID[source].Certfp = args[1]

-- :42F ENCAP * SNOTE s :Failed OPER attempt - host mismatch by xena (xena@0::z)
Hook "ENCAP-SNOTE", (source, args) ->
  tetra.RunHook("ENCAP-SNOTE-#{args[1]\upper!}", args[2])
