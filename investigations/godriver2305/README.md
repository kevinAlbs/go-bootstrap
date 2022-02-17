# Reproducing steps
Start a 2 node replica set:
mlaunch init --replicaset --nodes=2 --name=rs0 --priority --binarypath bin --hostname localhost --enableMajorityReadConcern --setParameter enableTestCommands=1  --dir rs2

Reconfig to add tags:
cfg = rs.config()
cfg.members[0].tags = {t1: "a", t2: "b"}
cfg.members[1].tags = {t1: "a", t2: "b"}
rs.reconfig(cfg)

Then try to connect with secondary read preference:

MONGODB_URI="mongodb://localhost:27017/?readPreference=secondary&readPreferenceTags=foo:bar&readPreferenceTags=" go run ./count

