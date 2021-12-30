FLAGS=$@
echo "Running test 10 times with flags: $FLAGS"
for i in $(seq 1 10); do
    go run main.go -quiet $FLAGS
done