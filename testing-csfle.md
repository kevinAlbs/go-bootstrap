# Testing CSFLE

Install [libmongocrypt](https://github.com/mongodb/libmongocrypt#installing-libmongocrypt-on-macos). On macOS, use `brew install mongodb/brew/libmongocrypt`.

Ensure `mongocryptd` is downloaded and included in the `PATH`. mongocryptd is included in the downloads from [download-mongodb.sh](https://github.com/mongodb-labs/drivers-evergreen-tools/blob/c862bb9e1ed85c589748483d6b3f8e7c17f19169/.evergreen/download-mongodb.sh).

Get credentials for Key Management Service (KMS) providers (AWS, GCP, and Azure). Ask Kevin over slack to send the `kms_providers.json` file. Store that file in the path `~/.csfle/kms_providers.json`.

Install the `aws` CLI and `jq` tool. Use the following script to load the credentials as environment variables expected by the Go driver tests:

```sh
export AWS_ACCESS_KEY_ID=$(cat ~/.csfle/kms_providers.json | jq -r ".aws.accessKeyId")
export AWS_SECRET_ACCESS_KEY=$(cat ~/.csfle/kms_providers.json | jq -r ".aws.secretAccessKey")
export AZURE_TENANT_ID=$(cat ~/.csfle/kms_providers.json | jq -r ".azure.tenantId")
export AZURE_CLIENT_ID=$(cat ~/.csfle/kms_providers.json | jq -r ".azure.clientId")
export AZURE_CLIENT_SECRET=$(cat ~/.csfle/kms_providers.json | jq -r ".azure.clientSecret")
export GCP_EMAIL=$(cat ~/.csfle/kms_providers.json | jq -r ".gcp.email")
export GCP_PRIVATE_KEY=$(cat ~/.csfle/kms_providers.json | jq -r ".gcp.privateKey")

TEMPCREDS=$(
    AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID \
    AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY \
    AWS_DEFAULT_REGION=us-east-1 \
    aws sts get-session-token
)

export CSFLE_AWS_TEMP_ACCESS_KEY_ID=$(echo $TEMPCREDS | jq -r ".Credentials.AccessKeyId")
export CSFLE_AWS_TEMP_SECRET_ACCESS_KEY=$(echo $TEMPCREDS | jq -r ".Credentials.SecretAccessKey")
export CSFLE_AWS_TEMP_SESSION_TOKEN=$(echo $TEMPCREDS | jq -r ".Credentials.SessionToken")

# TODO: update this path.
export CSFLE_TLS_CA_FILE=/Users/kevin.albertson/code/drivers-evergreen-tools/.evergreen/x509gen/ca.pem
export CSFLE_TLS_CERTIFICATE_KEY_FILE=/Users/kevin.albertson/code/drivers-evergreen-tools/.evergreen/x509gen/client.pem
```

This should enable running most CSFLE tests. This excludes a small number of mock server tests that I rarely run locally:

```
go test -tags cse -v -count=1 ./mongo/integration -run TestClientSideEncryptionSpec/basic
```


# Q & A
## Q3: How do I test with the `csfle` shared library?
libmongocrypt 1.5.0 supports the `csfle` shared library to replace the `mongocryptd` process.

At time of writing, libmongocrypt 1.5.0 is only released in alpha. See Q2 for obtaining 1.5.0-alpha1.

To get the `csfle` shared library, use the [mongodl.py script](https://github.com/mongodb-labs/drivers-evergreen-tools/blob/c862bb9e1ed85c589748483d6b3f8e7c17f19169/.evergreen/mongodl.py) as follows:

```
python mongodl.py --version 6.0.0-rc3 --component csfle --out csfle
```

## Q2: How can I get a development build of libmongocrypt?

On macOS, use `brew install mongodb/brew/libmongocrypt --HEAD` to get the latest build from the `master` branch.

Alternatively, get the built binaries attached to the [`upload-all` task in Evergreen](https://evergreen.mongodb.com/waterfall/libmongocrypt?task_filter=upload-all).

## Q1: How do I use a non-system install of libmongocrypt?

On Windows, this is not possible. The Go driver expects libmongocrypt to be in the path: C:\libmongocrypt.

On Unix, the Go driver uses `pkg-config` to determine the path to libmongocrypt. Suppose libmongocrypt is installed in the non-system directory `/install/libmongocrypt`. Specify the following environment variables to use a non-system install of libmongocrypt in the Go driver:

```
export PKG_CONFIG_PATH=/install/libmongocrypt/lib/pkgconfig
# Specify DYLD_LIBRARY_PATH for macOS
export DYLD_LIBRARY_PATH=/install/libmongocrypt/lib
# Specify LD_LIBRARY_PATH for Linux
export LD_LIBRARY_PATH=/install/libmongocrypt/lib

go test -tags cse -v -count=1 ./mongo/integration -run TestClientSideEncryptionSpec/basic
```

If headers change, use `go clean -cache -testcache`.