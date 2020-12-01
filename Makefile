bootstrap: ./runner/bootstrap.go ./runner/util.go
	go build -o bootstrap ./runner/bootstrap.go ./runner/util.go

clean: bootstrap
	rm bootstrap