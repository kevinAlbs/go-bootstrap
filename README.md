# MongoDB Go Driver Bootstrap #

Provides a simple way to run independent tests against a local copy of [mongo-go-driver](git@github.com:mongodb/mongo-go-driver.git).

To use locally, update the `replace` directive in [go.mod](./go.mod) to point to your copy of mongo-go-driver.

[runner](./runner) provides a main function and simple utility for logging. To use it:
```
make bootstrap
./bootstrap
```

Alternatively, write a test file. For example:
```
cd investigations
go test -v -count=1 -run Test1777
```