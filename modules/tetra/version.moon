Command "VERSION", ->
  commit = io.popen("git rev-parse --short HEAD")\read("*a")
  return "Running Tetra 0.1-#{commit}"
