# Add path to `mongocryptd` so driver can spawn mongocryptd.
export PATH=$PATH:/Users/kevin.albertson/bin/mongodl/archive/7.0.0-rc0/mongodb-macos-aarch64-enterprise-7.0.0-rc0/bin/
go run -tags cse .

# A `Library not loaded` error may be caused by go runtime caching. It may be solved by running `go clean -cache -testcache`
