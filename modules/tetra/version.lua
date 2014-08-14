Command("VERSION", function(source, message)
  local commit = io.popen("git rev-parse --short HEAD"):read("*a")
  return "Running Tetra 0.1-" .. commit
end)
