if [[ -z "$NAME" ]]; then
	echo "Usage: NAME=<name> create_investigation.sh";
	exit 1
fi

if [[ -d "investigations/$NAME" ]]; then
    echo "investigations/$NAME already exists"
    exit 1
fi

mkdir investigations/$NAME

cat <<EOF > investigations/$NAME/main.go
package main

import (
	"os"

	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kevinalbs.com/go-bootstrap/util"
)

func main() {
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		uri = "mongodb://localhost:27017"
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri).SetMonitor(util.CreateMonitor()))
	if err != nil {
		panic(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
}
EOF