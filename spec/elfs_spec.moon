describe "elfs can", ->
  it "be imported", ->
    require "lib/elfs"

    assert.truthy elfs

  it "generate names", ->
    require "lib/elfs"

    name = elfs.GenName!

    assert.truthy name
    assert.True #name > 0
