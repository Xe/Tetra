require "json"
require "socket"
require "moonscript"

-- http://stackoverflow.com/a/12674376
function keys(tab)
  local keyset={}
  local n=0

  for k,v in pairs(tab) do
    n=n+1
    keyset[n]=k
  end

  return keyset
end

--- Command wraps a function to be a bot command
-- @param verb the command verb
-- @param operonly if the command should be restricted to opers or not
-- @param func the function to wrap
function Command(verb, operonly, func) --> *tetra.Command
  if func == nil then
    func = operonly
    operonly = false
  end

  verb = verb:upper()

  local my_uuid = uuid.new()
  _G[my_uuid] = function(...)
    return func(...)
  end

  local cmd, err = tetra.script.AddLuaCommand(verb, my_uuid)

  if err ~= nil then error(err) end

  client.Commands[verb].NeedsOper = operonly

  return cmd
end

--- Hook wraps a function to act as a hook
-- @param verb the hook verb to listen for
-- @param func the function to wrap
function Hook(verb, func) --> *tetra.Hook
  verb = verb:upper()

  local my_uuid = uuid.new()
  _G[my_uuid] = function(...) func(...) end

  return tetra.script.AddLuaHook(verb, my_uuid)
end

--- Protohook wraps a function to be called on a protocol verb
-- @param verb the protocol verb to be called on
-- @param func the function to call
function Protohook(verb, func) --> *tetra.Handler
  verb = verb:upper()

  local my_uuid = uuid.new()
  _G[my_uuid] = function(...) func(...) end

  return tetra.script.AddLuaProtohook(verb, my_uuid)
end

function geturl(url)
  local c, err = web.get(url)
  if err ~= nil then
    tetra.log.Printf("URL error: %#v", err)
    return nil, err
  end

  local str, err = ioutil.readall(c.Body)
  if err ~= nil then
    tetra.log.Printf("Read error: %#v", err)
    return nil, err
  end

  str = ioutil.byte2string(str)

  return str, nil
end

function getjson(url)
  obj = json.decode(geturl(url))

  return obj, nil
end

function parseLine(line)
  local source = tetra.Clients.ByUID[line.Source]
  local destination = line.Args[1]
  local message = line.Args[2]

  if destination:sub(1,1) == "#" then
    destination = destination:upper()
    destination = tetra.Channels[destination]
  else
    destination = tetra.Clients.ByUID[destination]
  end

  return source, destination, message
end

-- https://stackoverflow.com/questions/2282444/
function contains(table, element)
  for _, value in pairs(table) do
    if value == element then
      return true
    end
  end
  return false
end

function find(tab, val)
  for i=1, #tab do
    if tab[i] == val then return i end
  end

  return 0
end

function sleep(sec)
  socket.select(nil, nil, sec)
end

function try(t)
  local ok, err = pcall(t.main)
  if not ok then
    t.catch(err)
  end
  if t.finally then
    return t.finally()
  end
end

--[[
-- Usage:
--
-- try {
--  main = function()
--    io.open("filethatDoesnOtexist", "r")
--  end,
--  catch = function(e)
--    print("Caught error!", e)
--  end,
-- }
--]]

-- Golua sucks
pcall = unsafe_pcall
xpcall = unsafe_xpcall

function url_encode(str)
  if (str) then
    str = string.gsub (str, "\n", "\r\n")
    str = string.gsub (str, "([^%w %-%_%.%~])",
    function (c) return string.format ("%%%02X", string.byte(c)) end)
    str = string.gsub (str, " ", "+")
  end
  return str
end

