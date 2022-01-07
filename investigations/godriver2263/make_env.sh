ROOT=$(pwd)
rm -rf .menv
mlaunch init \
    --single \
    --binarypath $(m bin 5.0-ent) \
    --tlsCAFile=$ROOT/certs/Drivers-Testing-CA-combined.pem \
    --tlsCertificateKeyFile=$ROOT/certs/Drivers-Testing-Server-combined.pem \
    --tlsMode=requireTLS \
    --dir .menv