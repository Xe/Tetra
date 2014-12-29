ld = nil

Hook("CHATBOT-CHANMSG", function(source, destination, message)
  if not source.IsOper() then
    return
  end

  message = strings.join(message, " ")

  ld = destination

  if message:sub(1,1) == "=" then
    tetra.log.Printf(message)
    toeval = message:sub(3)

    local func, err = loadstring(toeval)

    if err ~= nil then
      client.Privmsg(destination, string.format("error: %s", err))
      return
    end

    client.ServicesLog(source.Nick .. ": EVAL: " .. toeval)

    local res, err = func()

    tetra.log.Printf("%#v: %#v", res, err)

    if res ~= nil then
      client.Privmsg(ld, "> " .. res)
    else
      client.Privmsg(ld, "> nil")
    end
  end
end)  
