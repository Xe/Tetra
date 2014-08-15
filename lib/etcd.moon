yaml = require "yaml"

--- Module etcd implements useful interfaces between Tetra's etcd client and lua tables.
module "etcd", package.seeall
export ^

--- Store implements a simple etcd-backed table using yaml serialization.
-- It is pretty simple. Set data in self.table and then commit it with self.Commit.
-- If path is nil, a path in the form of `/tetra/script/scriptname/kind` will
-- be generated.
--
-- @param kind the kind of data being stored
-- @param path an arbitrary path inside etcd for the data to be stored or nil
class Store
  new: (kind, path=nil) =>
    if path == nil
      @path = "/tetra/script/" .. script.Name .. "/#{kind}"
    else
      @path = "/tetra/#{path}/#{kind}"

    @data = {}

  @Load!

  --- Load loads the table from etcd, discarding the local copy unless
  -- deserialization from the yaml document fails.
  -- @param self the instance of Store
  Load: =>
    etcd_value = tetra.bot.Etcd.Get @path, false, false

    if etcd_value == nil
      return

    data, err = yaml.load(etcd_value.Node.Value)

    if err ~= nil
      error(err)

    @data = data

  --- Commit saves the local changes to etcd. This can take longer if your
  -- table is big enough.
  -- @param self the instance of Store to commit
  Commit: =>
    tetra.bot.Etcd.Set @path, yaml.dump(@data), 0
