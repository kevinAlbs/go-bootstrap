bootstrap: ./runner/*.go
	go run -tags cse ./runner/bootstrap.go ./runner/util.go

test: ./investigations/spec1768_test.go
	go test ./investigations -tags cse -v -run TestSpec1768Case5

clean: bootstrap
	rm bootstrap