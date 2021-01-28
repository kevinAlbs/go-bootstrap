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

# FAQ #

## How many connections are expected with a maxPoolSize=M ##

Assume one Client, and S servers, with maxPoolSize=M.

For Go driver 1.4.0+: max of `S * (2 + M)`.
- There are 2 * S connections for monitoring (one for ismaster, one for RTT monitoring).
- There are a max of M * S connections for application use.
- Changed in streamable ismaster https://jira.mongodb.org/browse/GODRIVER-1489. The RTT monitoring thread is spawned regardless of the server version.

For Go driver < 1.4.0: max of `S * (1 + M)`.
- Note, there was a bug fix for maxPoolSize in 1.3.5: https://jira.mongodb.org/browse/GODRIVER-1613
