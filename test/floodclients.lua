smurfs = {}

for i=1,500 do 
  local c = tetra.bot.AddService(
    "smurf"..i,
    uuid.new():sub(1,20),
    "user",
    "host",
    "Lol spamming to get stats fed"
  )
  c.Join("#"..c.Nick)
end