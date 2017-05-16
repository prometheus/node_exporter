build:
	go fmt
	go build
	go vet
	staticcheck
	#golint -set_exit_status
	go test -v -race -tags=integration

cover:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out
