function test(user, message)
  return "Test!"
end

tetra.script.AddLuaCommand("TEST", "test")
