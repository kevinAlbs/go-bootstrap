# Testing Serverless

I have the credentials to create a serverless instance in a local file: ~/.serverless/serverless-env.sh. Ask over slack for the contents.

To create a new serverless instance to test:
```
. ~/.serverless/serverless-env.sh
cd ~/code/drivers-evergreen-tools/.evergreen/serverless
./create-instance.sh
```

`./create-instance.sh` produces an `serverless-expansions.yml` file. Use that file to set environment variables for the Go driver:

```
export MONGODB_URI="<SERVERLESS_URI from serverless-expansions.yml>"
export MONGODB_SRV_URI="<SERVERLESS_URI from serverless-expansions.yml>"
export SSL="ssl"
export AUTH="auth"
export SERVERLESS="serverless"
export SINGLE_ATLASPROXY_SERVERLESS_URI="<SERVERLESS_URI from serverless-expansions.yml>"
export MULTI_ATLASPROXY_SERVERLESS_URI="<SERVERLESS_URI from serverless-expansions.yml>"
```

Then run tests in the Go driver:
```
go test -count=1 ./mongo/integration -run TestCrudSpec -v
```

Once finished, use the `delete-instance.sh` script to delete the instance. It is automatically reaped after a time period.
```
export SERVERLESS_INSTANCE_NAME=<SERVERLESS_INSTANCE_NAME from serverless-expansions.yml>
./delete-instance.sh
```

