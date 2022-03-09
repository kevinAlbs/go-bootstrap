// Use mongod 5.3.0 or newer to support collectionUUID argument in SERVER-62445
// Start mongod as follows:
// mongod --dbpath ./tmp --setParameter featureFlagCommandsAcceptCollectionUUID=1

// in mongosh:
// create two collections
db.foo.insert({x:1})
db.foo2.insert({x:1})
// get the UUID of foo2
let foo2uuid = db.getCollectionInfos({name: "foo2"})[0]["info"]["uuid"]
// use the wrong UUID in foo2
db.runCommand({"insert": "foo", "documents": [{x:1}], "collectionUUID": foo2uuid})
/*
{
  n: 0,
  writeErrors: [
    {
      index: 0,
      code: 361,
      collectionUUID: UUID("ef8ba370-ed22-493f-b7bd-9c260ad00b02"),
      actualNamespace: 'test.foo2',
      errmsg: 'Collection UUID does not match that specified'
    }
  ],
  ok: 1
}
*/