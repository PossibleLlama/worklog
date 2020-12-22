NO_COLOUR=\x1b[0m
OK_COLOUR=\x1b[32;01m
WARN_COLOUR=\x1b[33;01m

BIN_NAME=worklog

test:
	@echo "$(WARN_COLOUR)Running unit tests$(NO_COLOUR)"
	go clean -testcache ./...
	go test ./...
	@echo "$(OK_COLOUR)Unit tests passed$(NO_COLOUR)"

format:
	@echo "$(WARN_COLOUR)Running format checks$(NO_COLOUR)"
	go fmt ./...
	go vet ./...
	@echo "$(OK_COLOUR)Format checks passed$(NO_COLOUR)"

build:
	@echo "$(WARN_COLOUR)Building to $(DIR)$(NO_COLOUR)"
	go build -ldflags="-w -s" -o "$(DIR)/$(BIN_NAME)" main.go
	@echo "$(OK_COLOUR)Built to $(DIR)$(NO_COLOUR)"

build-local:
	make build DIR="/tmp"

compile:
	@echo "$(WARN_COLOUR)Compiling for OS and Platform$(NO_COLOUR)"
	@echo "$(WARN_COLOUR)32 bit systems$(NO_COLOUR)"
	GOOS=freebsd GOARCH=386 go build -ldflags='-w -s' -o bin/32bit/$(BIN_NAME)-freebsd main.go
	GOOS=linux GOARCH=386 go build -ldflags='-w -s' -o bin/32bit/$(BIN_NAME)-linux main.go
	GOOS=windows GOARCH=386 go build -ldflags='-w -s' -o bin/32bit/$(BIN_NAME)-windows main.go

	@echo "$(WARN_COLOUR)64 bit systems$(NO_COLOUR)"
	GOOS=freebsd GOARCH=amd64 go build -ldflags='-w -s' -o bin/64bit/$(BIN_NAME)-freebsd main.go
	GOOS=darwin GOARCH=amd64 go build -ldflags='-w -s' -o bin/64bit/$(BIN_NAME)-darwin main.go
	GOOS=linux GOARCH=amd64 go build -ldflags='-w -s' -o bin/64bit/$(BIN_NAME)-linux main.go
	GOOS=windows GOARCH=amd64 go build -ldflags='-w -s' -o bin/64bit/$(BIN_NAME)-windows main.go
	@echo "$(OK_COLOUR)Built all binaries$(NO_COLOUR)"

deps:
	@echo "$(WARN_COLOUR)Checking and downloading dependencies$(NO_COLOUR)"
	go mod tidy
	go mod download
	go mod verify
	@echo "$(OK_COLOUR)Passed all dependency checks$(NO_COLOUR)"
