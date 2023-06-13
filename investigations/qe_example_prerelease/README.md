# Queryable Encryption v2 (QEv2) Example

Contains an example project using the Go driver with QEv2.

The `go.mod` file was created to depend on an unreleased Go driver containing support for QEv2 with:

```
go mod init
go get go.mongodb.org/mongo-driver@c7207c3735c3bc9212592a03d3a16f77e3012d09
go mod tidy
```

## Run the example

Run a 7.0+ MongoDB server.

Install libmongocrypt 1.8.0+. On macOS, do `brew install mongodb/brew/libmongocrypt`.

Ensure mongocryptd is on the path.

Run the example: `go run -tags cse .`

See the `run.sh` script for a full example.
