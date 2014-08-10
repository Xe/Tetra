VERSION = function(source, message)
  local commit = os.capture("git rev-parse --short HEAD")
  return "Tetra 0.1-" .. commit
end

tetra.script.AddLuaCommand("VERSION", "VERSION")
