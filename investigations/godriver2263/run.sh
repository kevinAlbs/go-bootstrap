export ROOT=/Users/kevin.albertson/code/go-bootstrap/investigations/godriver2263
export CLIENTCERT=$ROOT/certs/Drivers-Testing-Client-Intermediate-combined.pem
export CAPATH=$ROOT/certs/Drivers-Testing-CA-combined.pem
export MONGODB_URI="mongodb://localhost:27017/?tlsCAFile=$CAPATH&tlsCertificateKeyFile=$CLIENTCERT&tls=true"
go run .