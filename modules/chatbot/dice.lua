Command("ROLL", function(source, destination, message)
  if #message < 1 then
    return "Need a kind of dice to roll! (try 1d20)"
  end

  if message[1]:upper() == "RICK" then
    return "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
  end

  local number, kind = message[1]:match("(%d+)d(%d+)")

  kind = tonumber(kind)
  if kind == nil and number == nil then
    kind = message[1]:match("d(%d+)")
    kind = tonumber(kind)

    if kind == nil then
      return "Invalid kind of die! Try a number."
    end
  end

  number = tonumber(number)
  if number == nil then
    number = 1
  end

  result = 0

  if number > 50 then
    return "Too many dice"
  end

  if kind > 50 then
    return "Too many sides"
  end

  for i = 1,number do
    result = result + math.random(kind)
  end

  return "Dice roll of " .. number .. " " .. kind .. "-sided dice: " .. result
end)

math.randomseed(os.time())
