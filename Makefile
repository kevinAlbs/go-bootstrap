bootstrap: ./runner/*.go
	go run ./runner/bootstrap.go 

test: ./investigations/spec1768_test.go
	go test ./investigations -tags cse -v -run TestSpec1768Case6

clean: bootstrap
	rm bootstrap