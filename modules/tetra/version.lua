VERSION = command("VERSION") .. function(source, message)
  local commit = os.capture("git rev-parse --short HEAD")
  return "Running Tetra 0.1-" .. commit
end
