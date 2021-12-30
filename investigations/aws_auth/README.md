Start mongod with auth enabled:

```
mlaunch init --replicaset  --nodes=1 --auth --username user --password password --name replicaSet --setParameter enableTestCommands=1 --verbose --binarypath $(m bin 5.0-ent)  --setParameter authenticationMechanisms="MONGODB-AWS,SCRAM-SHA-1" --dir menv
```