require "json"
require "socket"

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

-- http://lua-users.org/wiki/SplitJoin
-- Compatibility: Lua-5.1
function split(str, pat)
  local t = {}  -- NOTE: use {n = 0} in Lua-5.0
  local fpat = "(.-)" .. pat
  local last_end = 1
  local s, e, cap = str:find(fpat, 1)
  while s do
    if s ~= 1 or cap ~= "" then
      table.insert(t,cap)
    end
    last_end = e+1
    s, e, cap = str:find(fpat, last_end)
  end
  if last_end <= #str then
    cap = str:sub(last_end)
    table.insert(t, cap)
  end
  return t
end

function join(table, str)
  return table.concat(table, str)
end

-- http://lua-users.org/wiki/SimpleLuaClasses
-- class.lua
-- Compatible with Lua 5.1 (not 5.0).
function class(base, init)
  local c = {}    -- a new class instance
  if not init and type(base) == 'function' then
    init = base
    base = nil
  elseif type(base) == 'table' then
    -- our new class is a shallow copy of the base class!
    for i,v in pairs(base) do
      c[i] = v
    end
    c._base = base
  end
  -- the class will be the metatable for all its objects,
  -- and they will look up their methods in it.
  c.__index = c

  -- expose a constructor which can be called by <classname>(<args>)
  local mt = {}
  mt.__call = function(class_tbl, ...)
    local obj = {}
    setmetatable(obj,c)
    if init then
      init(obj,...)
    else
      -- make sure that any stuff from the base class is initialized!
      if base and base.init then
        base.init(obj, ...)
      end
    end
    return obj
  end
  c.init = init
  c.is_a = function(self, klass)
    local m = getmetatable(self)
    while m do
      if m == klass then return true end
      m = m._base
    end
    return false
  end
  setmetatable(c, mt)
  return c
end

--[[
A = class(function (self, x)
self.x = x
end)

function A:test()
print(self.x)
end

a = A("5")
a:test() --> 5
--]]

-- Class for a limited length table
LimitQueue = class(function(self, max)
  self.max = max
  self.table = {}
end)

function LimitQueue:Add(data)
  local ret = false

  if #self.table == self.max then
    table.remove(self.table, 1)

    ret = true
  end

  table.insert(self.table, data)

  return ret
end

function LimitQueue:Pop()
  return table.remove(self.table, 1)
end

-- Simple disk-backed table
FooDB = class(function(self, fname)
  self.fname = fname

  local fhandle = io.open(fname, "r")

  if fhandle == nil then
    self.data = {}
    return
  end

  local data = fhandle:read("*a")
  fhandle:close()

  self.data = json.decode(data)

  if self.data == nil then self.data = {} end
end)

function FooDB:Commit()
  local fhandle = io.open(self.fname, "w")
  if fhandle == nil then
    error("Cannot open "..self.fname)
  end
  local string = json.encode(self.data)

  fhandle:write(string)
  fhandle:close()
end



--- Command wraps a function to be a bot command
-- @param verb the command verb
-- @param operonly if the command should be restricted to opers or not
-- @param func the function to wrap
function Command(verb, operonly, func)
  if func == nil then
    func = operonly
    operonly = false
  end

  verb = verb:upper()

  local my_uuid = uuid.new()
  _G[my_uuid] = function(...)
    return func(...)
  end

  local _, err = tetra.script.AddLuaCommand(verb, my_uuid)

  if err ~= nil then error(err) end

  client.Commands[verb].NeedsOper = operonly
end

--- Hook wraps a function to act as a hook
-- @param verb the hook verb to listen for
-- @param func the function to wrap
function Hook(verb, func)
  verb = verb:upper()

  local my_uuid = uuid.new()
  _G[my_uuid] = function(...) func(...) end

  tetra.script.AddLuaHook(verb, my_uuid)
end

--- Protohook wraps a function to be called on a protocol verb
-- @param verb the protocol verb to be called on
-- @param func the function to call
function Protohook(verb, func)
  verb = verb:upper()

  local my_uuid = uuid.new()
  _G[my_uuid] = function(...) func(...) end

  tetra.script.AddLuaProtohook(verb, my_uuid)
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
  local source = tetra.bot.Clients.ByUID[line.Source]
  local destination = line.Args[1]
  local message = line.Args[2]

  if destination:sub(1,1) == "#" then
    destination = destination:upper()
    destination = tetra.bot.Channels[destination]
  else
    destination = tetra.bot.Clients.ByUID[destination]
  end

  return source, destination, message
end

function is_common_channel(destination)
  if not destination.IsChannel() then return false end

  if client.Channels[destination.Target()] ~= nil then
    return true
  else
    return false
  end
end

function is_targeted_pm(destination)
  return not destination.IsChannel() and destination.Nick == client.Nick
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
--  main: function()
--    io.open("filethatDoesnOtexist", "r")
--  end,
--  catch: function(e)
--    print("Caught error!", e)
--  end,
-- }
--]]

-- Golua sucks
pcall = unsafe_pcall
xpcall = unsafe_xpcall

local yaml = require("yaml")

do
  local _base_0 = {
    Commit = function(self)
      return tetra.bot.Etcd.Set(self.path, yaml.dump(self.data), 0)
    end
  }
  _base_0.__index = _base_0
  local _class_0 = setmetatable({
    __init = function(self, kind)
      self.path = "/tetra/script/" .. script.Name .. "/" .. tostring(kind)
      self.data = { }
      local etcd_value = tetra.bot.Etcd.Get(self.path, false, false)
      if etcd_value == nil then
        return 
      end
      local err
      self.data, err = yaml.load(etcd_value.Node.Value)
      if err ~= nil then
        return error(err)
      end
    end,
    __base = _base_0,
    __name = "EtcdStore"
  }, {
    __index = _base_0,
    __call = function(cls, ...)
      local _self_0 = setmetatable({}, _base_0)
      cls.__init(_self_0, ...)
      return _self_0
    end
  })
  _base_0.__class = _class_0
  EtcdStore = _class_0
  return _class_0
end
