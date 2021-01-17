BIN_NAME=worklog

test:
	make build DIR="./integration"
	@echo "Running tests"
	go clean -testcache ./...
	go test ./...
	@echo "Tests passed"
	rm ./integration/$(BIN_NAME)

test-repeat:
	make build DIR="./integration"
	@echo "Running tests multiple times"
	go clean -testcache ./...
	go test -count 100 ./...
	@echo "Tests passed"
	rm ./integration/$(BIN_NAME)

format:
	@echo "Running format checks"
	go fmt ./...
	go vet ./...
	@echo "Format checks passed"

build:
	@echo "Building to $(DIR)"
	go build -ldflags="-w -s" -o "$(DIR)/$(BIN_NAME)" main.go
	@echo "Built to $(DIR)"

build-local:
	make build DIR="/tmp"

compile:
	@echo "Compiling for OS and Platform"
	@echo "32 bit systems"
	GOOS=freebsd GOARCH=386 go build -ldflags='-w -s' -o bin/32bit/$(BIN_NAME)-freebsd main.go
	GOOS=linux GOARCH=386 go build -ldflags='-w -s' -o bin/32bit/$(BIN_NAME)-linux main.go
	GOOS=windows GOARCH=386 go build -ldflags='-w -s' -o bin/32bit/$(BIN_NAME)-windows main.go

	@echo "64 bit systems"
	GOOS=freebsd GOARCH=amd64 go build -ldflags='-w -s' -o bin/64bit/$(BIN_NAME)-freebsd main.go
	GOOS=darwin GOARCH=amd64 go build -ldflags='-w -s' -o bin/64bit/$(BIN_NAME)-darwin main.go
	GOOS=linux GOARCH=amd64 go build -ldflags='-w -s' -o bin/64bit/$(BIN_NAME)-linux main.go
	GOOS=windows GOARCH=amd64 go build -ldflags='-w -s' -o bin/64bit/$(BIN_NAME)-windows main.go
	@echo "Built all binaries"

deps:
	@echo "Checking and downloading dependencies"
	go mod tidy
	go mod download
	go mod verify
	@echo "Passed all dependency checks"
