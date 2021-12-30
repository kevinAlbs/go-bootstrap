for i in {1..10}; do
    go run ./investigations/help28905/main.go

    if [[ "$?" != "0" ]]; then
        echo "Test failed, see logs above!"
        break
    fi
done
evergreen notify slack --target "@kevin.albertson" --msg "Done with $MONGODB_URI"