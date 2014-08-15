Command "VERSION", ->
  local commit = io.popen("git rev-parse --short HEAD")\read "*a"
  "Running Tetra 0.1-" .. commit
