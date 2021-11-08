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

# Remember #
- `waitQueueTimeoutMS` is not supported (because of context deadlines).

# Q&A #

## How many connections are expected with a maxPoolSize=M ##

Assume one Client, and S servers, with maxPoolSize=M.

For Go driver 1.4.0+: max of `S * (2 + M)`.
- There are 2 * S connections for monitoring (one for ismaster, one for RTT monitoring).
- There are a max of M * S connections for application use.
- Changed in streamable ismaster https://jira.mongodb.org/browse/GODRIVER-1489. The RTT monitoring thread is spawned regardless of the server version.

For Go driver < 1.4.0: max of `S * (1 + M)`.
- Note, there was a bug fix for maxPoolSize in 1.3.5: https://jira.mongodb.org/browse/GODRIVER-1613

## Is it safe to call mongo.Client.Connect() after mongo.Client.Disconnect() ?
No, you will get a `server is closed` error. See investigations/connect_after_disconnect.

## Why aren't types like Client/Database/Collection interfaces?

The main hesitation around exposing an interface is backwards compatibility.
We haven't made the Client/Database/Collection types interfaces despite multiple user requests to do so. For struct types, the only breaking changes are removing existing functions or changing function signatures. We can add new functions whenever we want. For interfaces, though, even adding functions is a breaking change because we're technically breaking all external implementations by doing so.

## What is the default value of connectTimeoutMS?

[30 seconds](https://github.com/kevinAlbs/mongo-go-driver/blob/cdacb6473abf8f2abaac11f58b7577fbd148440e/x/mongo/driver/topology/connection_options.go#L60)

This differs from the URI options spec, which suggests 10 seconds.

## When was Go driver v1.x released?

tag                                                                               date
--------------------------------------------------------------------------------  ----------
[v1.7.0](https://api.github.com/repos/mongodb/mongo-go-driver/releases/46498265)  2021-07-20
[v1.6.0](https://api.github.com/repos/mongodb/mongo-go-driver/releases/46072524)  2021-07-12
[v1.5.0](https://api.github.com/repos/mongodb/mongo-go-driver/releases/39536658)  2021-03-09
[v1.4.0](https://api.github.com/repos/mongodb/mongo-go-driver/releases/29127200)  2020-07-30
[v1.3.0](https://api.github.com/repos/mongodb/mongo-go-driver/releases/23415933)  2020-02-05
[v1.2.0](https://api.github.com/repos/mongodb/mongo-go-driver/releases/22120596)  2019-12-10
[v1.1.0](https://api.github.com/repos/mongodb/mongo-go-driver/releases/19232847)  2019-08-13
[v1.0.0](https://api.github.com/repos/mongodb/mongo-go-driver/releases/16098456)  2019-03-13