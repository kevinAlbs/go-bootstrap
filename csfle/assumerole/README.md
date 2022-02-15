This is an example of using CSFLE with temporary credentials taken from an AWS AssumeRole request.

To run this example:

1. [Install libmongocrypt](https://github.com/mongodb/libmongocrypt#installing-libmongocrypt-on-macos).

On macOS: `brew install mongodb/brew/libmongocrypt`

2. Get temporary credentials.

```
export ACCESS_KEY_ID="..."
export SECRET_ACCESS_KEY="..."
export CMK_ARN="..."
export CMK_REGION="..."
. ./csfle/assumerole/get_credentials.sh
```

3. Run mongod.

4. Run script.

```
go run -tags cse ./csfle/assumerole
```
