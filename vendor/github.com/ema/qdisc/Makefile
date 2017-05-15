build:
	go fmt
	go build
	go vet
	staticcheck
	#golint -set_exit_status
	go test -v -race -tags=integration
