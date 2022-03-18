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

## Q1 How many connections are expected with a maxPoolSize=M ##

Assume one Client, and S servers, with maxPoolSize=M.

For Go driver 1.4.0+: max of `S * (2 + M)`.
- There are 2 * S connections for monitoring (one for ismaster, one for RTT monitoring).
- There are a max of M * S connections for application use.
- Changed in streamable ismaster https://jira.mongodb.org/browse/GODRIVER-1489. The RTT monitoring thread is spawned regardless of the server version.

For Go driver < 1.4.0: max of `S * (1 + M)`.
- Note, there was a bug fix for maxPoolSize in 1.3.5: https://jira.mongodb.org/browse/GODRIVER-1613

## Q2 Is it safe to call mongo.Client.Connect() after mongo.Client.Disconnect() ?
No, you will get a `server is closed` error. See investigations/connect_after_disconnect.

## Q3 Why aren't types like Client/Database/Collection interfaces?

The main hesitation around exposing an interface is backwards compatibility.
We haven't made the Client/Database/Collection types interfaces despite multiple user requests to do so. For struct types, the only breaking changes are removing existing functions or changing function signatures. We can add new functions whenever we want. For interfaces, though, even adding functions is a breaking change because we're technically breaking all external implementations by doing so.

This is frequently requested to make mocking Go driver types easier.

The Go wiki doc https://github.com/golang/go/wiki/CodeReviewComments#interfaces advises against interfaces in the API for mocking.

> Do not define interfaces on the implementor side of an API 'for mocking'; instead, design the API so that it can be tested using the public API of the real implementation.

## Q4 What is the default value of connectTimeoutMS?

[30 seconds](https://github.com/kevinAlbs/mongo-go-driver/blob/cdacb6473abf8f2abaac11f58b7577fbd148440e/x/mongo/driver/topology/connection_options.go#L60)

This differs from the URI options spec, which suggests 10 seconds.

## Q5 When was Go driver v1.x released?

| tag                                                                      | date       |
|--------------------------------------------------------------------------|------------|
| [v1.8.0](https://github.com/mongodb/mongo-go-driver/releases/tag/v1.8.0) | 2021-11-23 |
| [v1.7.0](https://github.com/mongodb/mongo-go-driver/releases/tag/v1.7.0) | 2021-07-20 |
| [v1.6.0](https://github.com/mongodb/mongo-go-driver/releases/tag/v1.6.0) | 2021-07-12 |
| [v1.5.0](https://github.com/mongodb/mongo-go-driver/releases/tag/v1.5.0) | 2021-03-09 |
| [v1.4.0](https://github.com/mongodb/mongo-go-driver/releases/tag/v1.4.0) | 2020-07-30 |
| [v1.3.0](https://github.com/mongodb/mongo-go-driver/releases/tag/v1.3.0) | 2020-02-05 |
| [v1.2.0](https://github.com/mongodb/mongo-go-driver/releases/tag/v1.2.0) | 2019-12-10 |
| [v1.1.0](https://github.com/mongodb/mongo-go-driver/releases/tag/v1.1.0) | 2019-08-13 |
| [v1.0.0](https://github.com/mongodb/mongo-go-driver/releases/tag/v1.0.0) | 2019-03-13 |

# Q6 Which of the index APIs is implemented?

There are three index APIs in mongodb/specifications:
- Enumerate Indexes
- Standard API
- Index View API

The Go driver implements the Index View API. GODRIVER-31 decided not to implement the Enumerate Indexes API.

# Q7 What is the process for the Cloud team requesting backports?

https://jira.mongodb.org/browse/CLOUDP-96157
# Q8 Does the Go driver generate deprecated UUID (subtype 3)? #
Yes.

See https://bsonspec.org/spec.html.
See ./investigations/legacy_uuid

# Q9 Is there a guide for migrating from mgo to the mongo-go-driver? #

https://www.mongodb.com/blog/post/go-migration-guide
But it was last updated in 2019 (pending https://jira.mongodb.org/browse/WEBSITE-11912)

# Q10 What version of the MongoDB server does globalsign/mgo support? #

I would not recommend using mgo with server version 4.0+. The README from the github.com/globalsign/mgo repo says support for “MongoDB 4.0 is currently experimental”.

# Q11 What are other drivers to compare against? #

Consider these:
- https://github.com/lib/pq - Postgres driver for Go database/sql package.
- https://pkg.go.dev/database/sql - the Go SQL database API
- https://github.com/arangodb/go-driver - the Go driver for ArangoDB
- https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/ - Go AWS SDK for DynamoDB